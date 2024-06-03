package internal

import (
	"fmt"
	handler "worker-session/api/handlers"
	"worker-session/api/routers"
	"worker-session/internal/session"
	"worker-session/pkg/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type Provider struct {
	Routers *chi.Mux
}

func App() *Provider {
	logger := logrus.NewEntry(logrus.New()).WithField("desc", "internal-app")

	// Instanciando serviço
	service := session.New()
	logger.Info("serviço instanciado")
	// Criando pasta padrão para todas as instâncias
	err := service.Create(session.INSTANCE_PATH)
	if err != nil {
		// Matando a aplicação
		logger.Panic(err)
	}
	logger.Info(fmt.Sprintf("pasta %s criada com sucesso", session.INSTANCE_PATH))

	// Iniciando variáveis de ambiente
	env, err := config.LoadEnv()
	if err != nil {
		// Matando a aplicação
		logger.Panic(err)
	}
	logger.WithField("status", "sucesso").Info("variáveis de ambiente carregadas")

	// Iniciando manipulador de rotas
	session_handler := handler.NewSession(service)
	logger.WithField("status", "sucesso").Info("manipulador de rotas iniciado")

	// Iniciando roteador
	session_router := routers.NewSession(
		session_handler, env,
	)
	logger.WithField("status", "sucesso").Info("inicialização do roteador")

	// Roteador Global
	global_router := chi.NewRouter()

	// Middlewares globais
	global_router.Use(middleware.RealIP)
	global_router.Use(middleware.RequestID)
	global_router.Use(middleware.Recoverer)

	if env.HttpLogs {
		global_router.Use(middleware.Logger)
	}

	// Funções auxiliares
	routers.Ping(global_router)
	routers.NotFound(global_router)

	global_router.Mount("/session", session_router.Routes())
	logger.WithField("status", "sucesso").Info("roteador global inicializado")

	return &Provider{Routers: global_router}
}
