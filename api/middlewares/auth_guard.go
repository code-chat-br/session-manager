package middlewares

import (
	"errors"
	"net/http"
	handler "worker-session/api/handlers"

	"github.com/go-chi/render"
)

type NextFunc func(next http.Handler) http.Handler

func AuthGuard(token string) NextFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			api_key := r.Header.Get("apikey")

			if api_key != token {
				response := handler.NewResponse(http.StatusUnauthorized)
				response.SetError(errors.New("Unauthorized"))

				render.Status(r, response.GetCode())
				render.JSON(w, r, response.ResponseError())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
