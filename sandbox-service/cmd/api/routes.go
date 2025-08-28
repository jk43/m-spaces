package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	//mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/test", app.Test)
	mux.Get("/sandbox", app.Sandbox)
	mux.Get("/rbac-manager", app.RbacManager)
	mux.Get("/kafka-history", app.KafkaHistory)

	mux.Route("/redis", func(mux chi.Router) {
		mux.Get("/", app.GetRedis)
		mux.Post("/", app.PostRedis)
		mux.Delete("/", app.DeleteRedis)

	})

	return mux
}
