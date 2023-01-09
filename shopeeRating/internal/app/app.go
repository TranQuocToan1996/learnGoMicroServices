package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/TranQuocToan1996/shopeerating/config"
	v1 "github.com/TranQuocToan1996/shopeerating/internal/controller/http/v1"
	"github.com/TranQuocToan1996/shopeerating/internal/usecase"
	"github.com/TranQuocToan1996/shopeerating/pkg/httpserver"
	"github.com/TranQuocToan1996/shopeerating/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Use case
	ratingUseCase := usecase.New(*cfg)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, ratingUseCase)
	// v2.NewRouter(handler, l, translationUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	var err error
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
