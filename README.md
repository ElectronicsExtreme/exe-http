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
    handler := exehttp.NewHandler(&Handler{}, "test")
    serveMux.Handle("/", handler)
    
    configs := make([]exehttp.ServerConfig, 0, 0)
    configs = append(configs, exehttp.ServerConfig{
        ServeMux:  serveMux,
        Address:   ":8080",
        TlsEnable: false,
    })
    exehttp.ListenAndServe(configs)
}

type Handler1 struct {
    errorLogInfo *exehttp.LogInfo
    transLogInfo *exehttp.LogInfo
}

func (self *Handler) SetErrorLogInfo(logInfo *exehttp.LogInfo) {
    self.errorLogInfo = logInfo
}

func (self *Handler) SetTransLogInfo(logInfo *exehttp.LogInfo) {
    self.transLogInfo = logInfo
}

type Data struct {
    Data string `json:"data"`
}

func (self *Data) Success() bool {
    return true
}

func (self *Handler) ServeHTTP(resp *exehttp.ResponseWriter, req *http.Request) {
    resp.WriteResults(&Data{Data: "Hellow World"})
}
```
