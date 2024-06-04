package handler

import (
	"encoding/json"
	"net/http"
	"worker-session/internal/session"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Session struct {
	service *session.Service
}

func NewSession(service *session.Service) *Session {
	return &Session{service: service}
}

func get_params(r *http.Request) (string, string, string) {
	group := chi.URLParam(r, "group")
	instance := chi.URLParam(r, "instance")
	key := chi.URLParam(r, "key")

	return group, instance, key
}

func (h *Session) POST_GroupFolder(r *http.Request) *Response {
	var body map[string]string
	e := render.DecodeJSON(r.Body, &body)
	if err := UnmarshalDescriptionError(e); err != nil {
		return err
	}

	status, err := h.service.CreateFolder(body["group"], "")

	response := NewResponse(status)

	if err != nil {
		response.SetError(err)
		return response
	}

	return response
}

func (h *Session) POST_InstanceFolder(r *http.Request) *Response {
	var body map[string]string
	e := render.DecodeJSON(r.Body, &body)
	if err := UnmarshalDescriptionError(e); err != nil {
		return err
	}

	params_group, _, _ := get_params(r)

	status, err := h.service.CreateFolder(params_group, body["instance"])

	response := NewResponse(status)

	if err != nil {
		response.SetError(err)
		return response
	}

	return response
}

func (h *Session) DELETE_InstanceFolder(r *http.Request) *Response {
	params_group, params_instance, _ := get_params(r)

	status, err := h.service.RemoveFolder(params_group, params_instance)

	response := NewResponse(status)

	if err != nil {
		response.SetError(err)
		return response
	}

	response.SetCode(status)

	return response
}

func (h *Session) POST_Credentials(r *http.Request) *Response {
	var body map[string]string
	e := render.DecodeJSON(r.Body, &body)
	if err := UnmarshalDescriptionError(e); err != nil {
		return err
	}

	params_group, params_instance, params_key := get_params(r)

	status, err := h.service.WriterCredentials(
		params_group, params_instance, params_key, body,
	)

	response := NewResponse(status)

	if err != nil {
		response.SetError(err)
		return response
	}

	return response
}

func (h *Session) GET_Credentials(r *http.Request) *Response {
	params_group, params_instance, params_key := get_params(r)

	status, binary, err := h.service.ReadCredentials(
		params_group, params_instance, params_key,
	)

	response := NewResponse(status)

	if err != nil {
		response.SetError(err)
		return response
	}

	var data map[string]any
	err = json.Unmarshal(binary, &data)
	if err != nil {
		response.SetError(err)
		return response
	}

	response.SetData(data)

	return response
}

func (h *Session) DELETE_Credentials(r *http.Request) *Response {
	params_group, params_instance, params_key := get_params(r)

	status, err := h.service.RemoveCredential(
		params_group, params_instance, params_key,
	)

	response := NewResponse(status)

	if err != nil {
		response.SetError(err)
		return response
	}

	return response
}

func (h *Session) GET_ListInstances(r *http.Request) *Response {
	params_group, _, _ := get_params(r)

	status, list, err := h.service.ListInstances(params_group)

	response := NewResponse(status)

	if err != nil {
		response.SetError(err)
		return response
	}

	response.SetData(list)

	return response
}
