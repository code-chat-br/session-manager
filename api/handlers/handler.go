package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Response struct {
	code int
	data any
	err  string
}

func NewResponse(code int) *Response {
	return &Response{code: code}
}

func (r *Response) SetCode(code int) {
	r.code = code
}

func (r *Response) SetData(data any) {
	r.data = data
}

func (r *Response) SetError(err error) {
	r.err = err.Error()
}

func (r *Response) GetCode() int {
	return r.code
}

func (r *Response) GetData() any {
	return r.data
}

func (r *Response) ResponseError() map[string]any {
	return map[string]any{
		"statusCode": r.code,
		"message":    r.data,
		"error":      r.err,
	}
}

func extractErrorDetails(err_msg string) (fieldName string, dataType string) {
	field_index := strings.Index(err_msg, "field ")
	of_index := strings.Index(err_msg, " of")
	if field_index == -1 || of_index == -1 {
		return "", ""
	}
	fieldName = err_msg[field_index+6 : of_index]

	typeIndex := strings.Index(err_msg, "type ")
	if typeIndex == -1 {
		return fieldName, ""
	}
	dataType = err_msg[typeIndex+5:]

	return fieldName, dataType
}

func UnmarshalDescriptionError(e error) *Response {
	if e != nil {
		response := NewResponse(http.StatusBadRequest)
		if e.Error() == "EOF" {
			response.SetData([]string{"body is empty"})
			response.SetError(errors.New("bad_request"))
			return response
		}
		field_name, data_type := extractErrorDetails(e.Error())
		response.SetData([]string{fmt.Sprintf("%s must be of type %s", field_name, data_type)})
		response.SetError(errors.New("bad_request"))
		return response
	}

	return nil
}
