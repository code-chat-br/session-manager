package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"worker-session/internal"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.NewEntry(logrus.New()).WithField("desc", "cmd/server/main")

	provider := internal.App()

	srv := &http.Server{
		Addr:    ":5656",
		Handler: provider.Routers,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.WithFields(logrus.Fields{"error": err.Error()}).Fatal("servidor parado")
		}
	}()

	logger.WithFields(logrus.Fields{"port": "5656"}).Log(logrus.WarnLevel, "servidor iniciado")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	os.Stdout.Sync()

	logger.WithFields(logrus.Fields{"signal": sig}).Log(logrus.WarnLevel, "Servidor desligando")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithFields(logrus.Fields{"error": err.Error()}).Fatal("Servidor forçado a desligar")
	} else {
		logger.Info("Servidor interrompido corretamente")
	}
}
