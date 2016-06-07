## Synopsis
Custom http server in GO language with graceful closing.

## Code Example
```go
package main

import (
    "net/http"

    "github.com/ElectronicsExtreme/exehttp"
)

func main() {
    exehttp.StartLogger("log")
    server := exehttp.NewServer(":9500")
    handler := exehttp.NewHandler(&Handler{}, "test")
    server.Handle("/", handler)
    
    servers := make([]exehttp.Server, 0, 0)
    servers = append(servers, server)
    exehttp.ListenAndServe(serverss)
}

type Handler struct {
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
    Info string `json:"info"`
}

func (self *Handler) ServeHTTP(resp *exehttp.ResponseWriter, req *http.Request) {
	Response := &exehttp.Results{
		Success: true,
		HTTPStatus: http.StatusOK,
		Data: &Data{Info: "Hello World"},
	}
		
    resp.WriteResults(Response)
}
```
