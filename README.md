# Healthcheck HTTP Handler [![Last release](https://img.shields.io/github/release/euskadi31/go-healthcheck.svg)](https://github.com/euskadi31/go-healthcheck/releases/latest) [![Documentation](https://godoc.org/github.com/euskadi31/go-healthcheck?status.svg)](https://godoc.org/github.com/euskadi31/go-healthcheck)

[![Go Report Card](https://goreportcard.com/badge/github.com/euskadi31/go-healthcheck)](https://goreportcard.com/report/github.com/euskadi31/go-healthcheck)

| Branch | Status                                                                                                                                                        | Coverage                                                                                                                                               |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| main   | [![Go](https://github.com/euskadi31/go-healthcheck/actions/workflows/go.yml/badge.svg)](https://github.com/euskadi31/go-healthcheck/actions/workflows/go.yml) | [![Coveralls](https://img.shields.io/coveralls/euskadi31/go-healthcheck/master.svg)](https://coveralls.io/github/euskadi31/go-healthcheck?branch=main) |

## Example

```go
package main

import (
    "github.com/euskadi31/go-healthcheck"
)

type RedisHealthcheck struct {

}

func (c *RedisHealthcheck) Check() bool {
    return true
}

func main() {
    hc := healthcheck.New()

    // use healthcheck.Handler
    hc.Add("redis", &RedisHealthcheck{})

    // use healthcheck.HandlerFunc
    err := hc.Add("mysql", healthcheck.HandlerFunc(func() bool {
        return true
    }))


    http.Handle("/health", hc)

    http.ListenAndServe(":8090", nil)
}
```
