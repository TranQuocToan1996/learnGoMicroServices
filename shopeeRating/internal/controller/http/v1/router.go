// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/TranQuocToan1996/shopeerating/internal/usecase"
	"github.com/TranQuocToan1996/shopeerating/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(handler *gin.Engine, l logger.Interface, u usecase.Rating) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		newRatingRoutes(h, u, l)
	}
}
