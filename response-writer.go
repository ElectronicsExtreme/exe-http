package exehttp

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	responseLogInfo *LogInfo
}

func NewResponseWriter(resp http.ResponseWriter, logInfo *LogInfo) *ResponseWriter {
	return &ResponseWriter{resp, logInfo}
}

func (self *ResponseWriter) WriteResults(data interface{}) {
	results := Results{}
	var httpStatus int = 0
	switch data := data.(type) {
	case *ErrorResponse:
		results.Success = false
		if data.HTTPStatus == 0 {
			log.Println("http status is not defined")
			self.WriteHeader(http.StatusInternalServerError)
			return
		}
		httpStatus = data.HTTPStatus
	case OtherResponse:
		results.Success = data.Success()
		httpStatus = http.StatusOK
	default:
		log.Printf("unknown response type %T\n", data)
		self.WriteHeader(http.StatusInternalServerError)
		return
	}
	results.Data = data
	resultsByte, err := json.Marshal(&results)
	if err != nil {
		self.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	self.WriteHeader(httpStatus)
	self.Write(resultsByte)

	//logWriter()
	self.responseLogInfo.HTTPStatus = httpStatus
	self.responseLogInfo.Body = string(resultsByte)
	self.responseLogInfo.Write()
}

type ErrorResponse struct {
	ErrorTag         string `json:"error"`
	ErrorDescription string `json:"error_description"`
	HTTPStatus       int    `json:"-"`
}

type OtherResponse interface {
	Success() bool
}

type Results struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

var ErrorStatusInternalServerError ErrorResponse = ErrorResponse{
	ErrorTag:         "internal_server_error",
	ErrorDescription: "internal server error",
	HTTPStatus:       500,
}
