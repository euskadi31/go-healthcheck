// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

//go:generate go run -mod=mod github.com/vektra/mockery/v2 --inpackage --case underscore --name=Handler
type Handler interface {
	Check() bool
}

type HandlerFunc func() bool

// Check calls f().
func (f HandlerFunc) Check() bool {
	return f()
}

// Response struct.
type Response struct {
	Status   bool            `json:"status"`
	Services map[string]bool `json:"services"`
}

var _ http.Handler = (*Healthcheck)(nil)

type Healthcheck struct {
	handlers map[string]Handler
}

func New() *Healthcheck {
	return &Healthcheck{
		handlers: make(map[string]Handler),
	}
}

// Add HealthCheck handler.
func (h *Healthcheck) Add(name string, handle Handler) error {
	if _, ok := h.handlers[name]; ok {
		return fmt.Errorf("the %s healthcheck handler already exists", name)
	}

	h.handlers[name] = handle

	return nil
}

func (h *Healthcheck) check() *Response {
	size := len(h.handlers)

	response := &Response{
		Status:   true,
		Services: make(map[string]bool, size),
	}

	var wg = &sync.WaitGroup{}

	wg.Add(size)

	var mutex = &sync.Mutex{}

	for name, handler := range h.handlers {
		go func(n string, h Handler) {
			defer wg.Done()

			s := h.Check()

			mutex.Lock()
			response.Services[n] = s
			defer mutex.Unlock()

			if !s {
				response.Status = false
			}
		}(name, handler)
	}

	wg.Wait()

	return response
}

// ServeHTTP implements http.Handler.
func (h *Healthcheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := http.StatusOK

	resp := h.check()

	if !resp.Status {
		code = http.StatusServiceUnavailable
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
