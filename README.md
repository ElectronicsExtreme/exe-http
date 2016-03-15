## Synopsis
Custom http server in GO language with graceful closing.

## Code Example
```go
package main

import (
    "net/http"

    "github.com/ElectronicsExtreme/exe-http"
)

func main() {
    serveMux := http.NewServeMux()
    serveMux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
        resp.Write([]byte("Hello, world"))
    })
    
    configs := make([]exeserver.ServerConfig, 0, 0)
    configs = append(configs, exeserver.ServerConfig{
        ServeMux:  serveMux,
        Address:   ":8080",
        TlsEnable: false,
    })
    exeserver.ListenAndServe(configs)
}
```
