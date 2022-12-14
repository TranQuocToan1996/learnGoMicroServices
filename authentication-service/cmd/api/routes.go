package main

import (
	"authentication/model"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func routes(cfg *model.Config) http.Handler {
	mux := chi.NewRouter()

	// Auth who can connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: Need fix after dev period
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", model.ContentType, "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		Debug:            true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/authenticate", cfg.Authenticate)
	

	return mux
}
