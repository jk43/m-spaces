package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/moly-space/molylibs/utils"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(utils.ParseJWTClaims)

	mux.Post("/save", app.RegisterUser)
	//mux.Post("/login", app.Login)

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(utils.AuthRequired)
		mux.Get("/", app.GetUser)
		mux.Put("/metadata", app.UpdateUserMetadata)
		mux.Put("/password", app.UpdateUserPassword)
		//update password with verification code
		mux.Post("/password", app.UpdatePasswordWithVerificationCode)
		mux.Put("/account", app.UpdateUserAccount)
		mux.Post("/email", app.UpdateEmailWithVerificationCode)
		mux.Post("/store", app.CreateStore)
		mux.Put("/store", app.UpdateStore)
		mux.Delete("/store", app.DeleteStore)
		mux.Get("/store/{ctx}", app.GetStore)
	})

	mux.Route("/admin", func(mux chi.Router) {
		mux.Get("/users", app.GetOrgUsers)
		mux.Post("/user", app.CreateOrgUser)
		mux.Put("/user", app.UpdateOrgUser)
		mux.Delete("/user", app.DeleteOrgUser)
	})

	return mux
}
