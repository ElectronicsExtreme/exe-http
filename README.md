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
    handler := exeserver.NewHandler(&Handler{}, "test")
    serveMux.Handle("/", handler)
    
    configs := make([]exeserver.ServerConfig, 0, 0)
    configs = append(configs, exeserver.ServerConfig{
        ServeMux:  serveMux,
        Address:   ":8080",
        TlsEnable: false,
    })
    exeserver.ListenAndServe(configs)
}

type Handler1 struct {
    errorLogInfo *exeserver.LogInfo
    transLogInfo *exeserver.LogInfo
}

func (self *Handler) SetErrorLogInfo(logInfo *exeserver.LogInfo) {
    self.errorLogInfo = logInfo
}

func (self *Handler) SetTransLogInfo(logInfo *exeserver.LogInfo) {
    self.transLogInfo = logInfo
}

type Data struct {
    Data string `json:"data"`
}

func (self *Data) Success() bool {
    return true
}

func (self *Handler) ServeHTTP(resp *exeserver.ResponseWriter, req *http.Request) {
    resp.WriteResults(&Data{Data: "Hellow World"})
}
```
