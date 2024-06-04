package routers

import (
	"errors"
	"net/http"
	"strings"
	handler "worker-session/api/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type HandlerFunc func(r *http.Request) *handler.Response

func ResponseRequest(handle HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := handle(r)

		render.Status(r, response.GetCode())

		if response.GetCode() >= 200 && response.GetCode() < 300 {
			render.JSON(w, r, response.GetData())
			return
		}

		render.JSON(w, r, response.ResponseError())
	}
}

func Ping(r *chi.Mux) {
	r.Options("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Mimetype", "text/plain")

		w.Write([]byte("pong"))
	})
}

func NotFound(r *chi.Mux) {
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		method := strings.ToUpper(r.Method)
		path := r.URL.Path

		response := handler.NewResponse(http.StatusNotFound)
		response.SetData("Cannot " + method + " " + path)
		response.SetError(errors.New("not_found"))

		render.Status(r, response.GetCode())
		render.JSON(w, r, response.ResponseError())
	})
}
