package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/moly-space/molylibs/utils"
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
	mux.Use(utils.ParseJWTClaims)
	mux.Get("/settings", app.GetSettingsForEndUser)
	mux.Get("/items", app.GetItems)
	mux.Get("/info", app.GetInfo)
	mux.Put("/info", app.UpdateInfo)
	mux.Put("/form-order", app.UpdateFormInputOrder)
	mux.Delete("/form", app.DeleteFormInput)
	mux.Put("/form", app.UpdateFormInput)
	mux.Post("/form", app.CreateFormInput)
	mux.Get("/form/{form}", app.GetForm)
	return mux
}
