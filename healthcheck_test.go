// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package healthcheck

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	hc := New()

	mysql := &MockHandler{}
	mysql.On("Check").Return(true)

	err := hc.Add("mysql", mysql)
	assert.NoError(t, err)

	err = hc.Add("mysql", mysql)
	assert.EqualError(t, err, "the mysql healthcheck handler already exists")

	hc.Add("redis", HandlerFunc(func() bool {
		return true
	}))

	response := hc.check()

	assert.True(t, response.Status)

	assert.Contains(t, response.Services, "mysql")
	assert.True(t, response.Services["mysql"])

	assert.Contains(t, response.Services, "redis")
	assert.True(t, response.Services["redis"])
}

func TestHealthcheckProcessorWithFailedCheck(t *testing.T) {
	hc := New()

	mysql := &MockHandler{}
	mysql.On("Check").Return(true)

	err := hc.Add("mysql", mysql)
	assert.NoError(t, err)

	redis := &MockHandler{}
	redis.On("Check").Return(false)

	err = hc.Add("redis", redis)
	assert.NoError(t, err)

	response := hc.check()

	assert.False(t, response.Status)

	assert.Contains(t, response.Services, "mysql")
	assert.True(t, response.Services["mysql"])

	assert.Contains(t, response.Services, "redis")
	assert.False(t, response.Services["redis"])
}

func TestHealthCheckServeHTTPWithSuccess(t *testing.T) {
	hc := New()

	mysql := &MockHandler{}
	mysql.On("Check").Return(true)

	err := hc.Add("mysql", mysql)
	assert.NoError(t, err)

	redis := &MockHandler{}
	redis.On("Check").Return(true)

	err = hc.Add("redis", redis)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)

	hc.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

type failWriter struct {
}

func (failWriter) Header() http.Header {
	return http.Header{}
}

func (failWriter) WriteHeader(statusCode int) {}

func (failWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("fail")
}

func TestHealthCheckServeHTTPWithFail(t *testing.T) {
	hc := New()

	mysql := &MockHandler{}
	mysql.On("Check").Return(true)

	err := hc.Add("mysql", mysql)
	assert.NoError(t, err)

	redis := &MockHandler{}
	redis.On("Check").Return(false)

	err = hc.Add("redis", redis)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)

	hc.ServeHTTP(w, r)

	assert.Equal(t, http.StatusServiceUnavailable, w.Result().StatusCode)

	w1 := &failWriter{}

	hc.ServeHTTP(w1, r)
}

func BenchmarkServeHTTP(b *testing.B) {
	hc := New()

	err := hc.Add("mysql", HandlerFunc(func() bool {
		return true
	}))
	assert.NoError(b, err)

	err = hc.Add("redis", HandlerFunc(func() bool {
		return true
	}))
	assert.NoError(b, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)

	for i := 0; i < b.N; i++ {
		hc.ServeHTTP(w, r)
	}
}
