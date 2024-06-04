package handler

import (
	"net/http"
	"strings"
)

type Response struct {
	Code int    `json:"statusCode"`
	Data any    `json:"message,omitempty"`
	Err  string `json:"error"`
}

func NewResponse(code int) *Response {
	return &Response{Code: code}
}

func (r *Response) SetCode(code int) {
	r.Code = code
}

func (r *Response) SetData(data any) {
	r.Data = data
}

func (r *Response) SetError(err error) {
	r.Err = err.Error()
}

func (r *Response) GetCode() int {
	return r.Code
}

func (r *Response) GetData() any {
	return r.Data
}

func extractErrorDetails(errMsg string) (fieldName string, dataType string) {
	fieldIndex := strings.Index(errMsg, "field ")
	ofIndex := strings.Index(errMsg, " of")
	if fieldIndex == -1 || ofIndex == -1 {
		return "", ""
	}
	fieldName = errMsg[fieldIndex+6 : ofIndex]

	typeIndex := strings.Index(errMsg, "type ")
	if typeIndex == -1 {
		return fieldName, ""
	}
	dataType = errMsg[typeIndex+5:]

	return fieldName, dataType
}

func UnmarshalDescriptionError(e error) *Response {
	if e != nil {
		if e.Error() == "EOF" {
			return &Response{
				Code: http.StatusBadRequest,
				Data: []string{"body is empty"},
				Err:  "Bad Request",
			}
		}
		fieldName, dataType := extractErrorDetails(e.Error())
		return &Response{
			Code: http.StatusBadRequest,
			Data: []string{fieldName + " must be of type " + dataType + "."},
			Err:  "Bad Request",
		}
	}

	return nil
}
