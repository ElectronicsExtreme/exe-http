package exeserver

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"time"
)

const (
	letterBytes   = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # ofletter indices fitting in 63 bits
)

var (
	src = rand.NewSource(time.Now().UnixNano())
)

type CentralHandler struct {
	Handler         Handler
	ErrorLogInfo    *LogInfo
	TransLogInfo    *LogInfo
	RequestLogInfo  *LogInfo
	ResponseLogInfo *LogInfo
}

type Handler interface {
	ServeHTTP(*ResponseWriter, *http.Request)
	SetTransLogInfo(*LogInfo)
	SetErrorLogInfo(*LogInfo)
	SetResponseLogInfo(*LogInfo)
}

func NewHandler(handler Handler, api string) *CentralHandler {
	refCode := randomString(6)
	newHandler := &CentralHandler{}
	newHandler.ErrorLogInfo = &LogInfo{
		API:     api,
		RefCode: refCode,
		ch:      errorChannel,
	}
	newHandler.TransLogInfo = &LogInfo{
		API:     api,
		RefCode: refCode,
		ch:      transactionChannel,
	}
	newHandler.RequestLogInfo = &LogInfo{
		API:     api,
		RefCode: refCode,
		ch:      requestChannel,
	}
	newHandler.ResponseLogInfo = &LogInfo{
		API:     api,
		RefCode: refCode,
		ch:      requestChannel,
	}
	handler.SetErrorLogInfo(newHandler.ErrorLogInfo)
	handler.SetTransLogInfo(newHandler.TransLogInfo)
	newHandler.Handler = handler
	return newHandler
}

func (self *CentralHandler) ServeHTTP(httpResp http.ResponseWriter, req *http.Request) {
	fmt.Println("start package")
	//request log
	data, err := dumpRequestBody(req)
	if err != nil {
		self.ErrorLogInfo.Body = err.Error()
		self.ErrorLogInfo.Write()
	}
	self.RequestLogInfo.Path = req.URL.Path
	self.RequestLogInfo.QueryString = req.URL.RawQuery
	self.RequestLogInfo.Method = req.Method
	self.RequestLogInfo.Body = string(data)
	self.RequestLogInfo.Write()
	resp := NewResponseWriter(httpResp, self.ResponseLogInfo)
	self.Handler.ServeHTTP(resp, req)
}

/*func (self *Handler) WriteError(err error) {
	self.ErrorLogInfo.Body = err.Error()
	self.ErrorLogInfo.Write()
}*/

func randomString(length int) string {
	outString := make([]byte, length)
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			outString[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(outString)
}

func dumpRequestBody(req *http.Request) ([]byte, error) {
	body, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	lastIndex := bytes.LastIndex(body, []byte("\r\n\r\n")) + 4
	return body[lastIndex:], nil
}
