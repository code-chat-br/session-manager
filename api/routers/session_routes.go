package routers

import (
	handler "worker-session/api/handlers"
	"worker-session/api/middlewares"
	"worker-session/pkg/config"

	"github.com/go-chi/chi/v5"
)

type Session struct {
	handler *handler.Session
	env     *config.Env
}

func NewSession(handler *handler.Session, env *config.Env) *Session {
	return &Session{
		handler: handler,
		env:     env,
	}
}

func (x *Session) Routes() *chi.Mux {
	group_router := chi.NewRouter()
	group_router.Route("/", func(r chi.Router) {
		r.Use(middlewares.AuthGuard(x.env.GlobalToken))

		// Criando a pasta que conterá o grupo de instâncias
		r.Post("/", ResponseRequest(x.handler.POST_Group))
		// Removendo o grupo com todas as instâncias
		r.Delete("/", ResponseRequest(x.handler.DELETE_Group))
	})

	instance_router := chi.NewRouter()
	instance_router.Route("/", func(r chi.Router) {
		r.Use(middlewares.LimitBodySize(x.env.LimitBody))

		r.Post("/", ResponseRequest(x.handler.POST_InstanceDB))
		r.Delete("/{instance}", ResponseRequest(x.handler.DELETE_InstanceDB))

		r.Post("/{instance}/{key}", ResponseRequest(x.handler.POST_Credentials))
		r.Get("/{instance}/{key}", ResponseRequest(x.handler.GET_Credentials))
		r.Delete("/{instance}/{key}", ResponseRequest(x.handler.DELETE_Credentials))

		r.Get("/list-instances", ResponseRequest(x.handler.GET_ListInstances))
	})

	group_router.Mount("/{group}", instance_router)

	return group_router
}
