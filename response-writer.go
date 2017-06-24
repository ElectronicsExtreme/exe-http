package exehttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	ResponseLogInfo *LogInfo
}

func NewResponseWriter(resp http.ResponseWriter, logInfo *LogInfo) *ResponseWriter {
	return &ResponseWriter{resp, logInfo}
}

func (self *ResponseWriter) WriteResults(response interface{}) error {
	results := Results{}
	var httpStatus int = 0
	switch response := response.(type) {
	case *ErrorResponse:
		results.Success = false
		if response.HTTPStatus == 0 {
			self.WriteResults(&ErrorStatusInternalServerError)
			return fmt.Errorf("http status is not defined")
		} else {
			httpStatus = response.HTTPStatus
		}
		results.Data = response
	case *Results:
		results = *response
		if response.HTTPStatus == 0 {
			httpStatus = http.StatusOK
		} else {
			httpStatus = response.HTTPStatus
		}
	default:
		self.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("unknown response type %T\n", response)
	}
	resultsByte, err := json.Marshal(&results)
	if err != nil {
		return err
		self.WriteResults(&ErrorStatusInternalServerError)
	}
	self.WriteHeader(httpStatus)
	self.Header().Set("Content-Type", "application/json")
	self.Write(resultsByte)

	//logWriter()
	if self.ResponseLogInfo != nil {
		self.ResponseLogInfo.HTTPStatus = httpStatus
		self.ResponseLogInfo.Body = string(resultsByte)
		self.ResponseLogInfo.Write()
	}
	return nil
}

func (self *ResponseWriter) WriteError(resp *ErrorResponse, description string) {
	if description != "" {
		self.WriteResults(&ErrorResponse{
			ErrorTag:         resp.ErrorTag,
			ErrorDescription: description,
			HTTPStatus:       resp.HTTPStatus,
		})
	} else {
		self.WriteResults(resp)
	}
}

type ErrorResponse struct {
	ErrorTag         string `json:"error"`
	ErrorDescription string `json:"error_description"`
	HTTPStatus       int    `json:"-"`
}

type Results struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	HTTPStatus int         `json:"-"`
}

var ErrorStatusInternalServerError ErrorResponse = ErrorResponse{
	ErrorTag:         "internal_server_error",
	ErrorDescription: "internal server error",
	HTTPStatus:       http.StatusInternalServerError,
}

var ErrorStatusNotFound ErrorResponse = ErrorResponse{
	ErrorTag:         "not_found",
	ErrorDescription: "not found",
	HTTPStatus:       http.StatusNotFound,
}

var ErrorStatusMethodNotAllowed ErrorResponse = ErrorResponse{
	ErrorTag:         "method_not_allowed",
	ErrorDescription: "method not allowed",
	HTTPStatus:       http.StatusMethodNotAllowed,
}
