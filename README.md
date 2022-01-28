# Healthcheck HTTP Handler

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
