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
	//mux.Use(middleware.Heartbeat("/ping"))

	// mux.Post("/password", app.SaveCredentials)
	mux.Get("/admin/ping", app.Ping)
	mux.Get("/admin/form", app.GetForm)
	mux.Get("/admin/boards", app.GetBoards)
	mux.Get("/admin/board", app.GetBoard)
	mux.Post("/admin/board", app.PostBoard)
	mux.Put("/admin/board", app.PutBoard)
	mux.Delete("/admin/board", app.DeleteBoard)

	mux.Get("/board", app.GetBoard)
	mux.Get("/post", app.GetPost)
	mux.Post("/posts", app.PostPosts)
	mux.Put("/posts", app.PutPosts)
	mux.Delete("/posts", app.DeletePosts)
	mux.Get("/posts", app.GetPosts)
	mux.Get("/comments", app.GetComments)
	mux.Get("/files", app.GetFiles)
	mux.Delete("/file", app.DeleteFile)
	return mux
}
