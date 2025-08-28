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

	mux.Use(utils.ParseJWTClaims)
	mux.Get("/tree/{slug}", app.GetTree)
	mux.Get("/admin/tree", app.GetAdminTree)
	mux.Post("/admin/tree", app.PostAdminTree)
	mux.Delete("/admin/tree", app.DeleteAdminTree)
	mux.Put("/admin/tree", app.PutAdminTree)
	mux.Get("/admin/trees", app.GetAdminTrees)
	mux.Put("/admin/reorder", app.PutAdminTreeReorder)
	mux.Post("/admin/trees", app.PostAdminTrees)
	return mux
}
