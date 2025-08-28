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

	// mux.Post("/password", app.SaveCredentials)
	mux.Post("/login", app.Login)
	mux.Post("/logout", app.Logout)
	mux.Get("/refreshtoken", app.RefreshToken)
	mux.Put("/password", app.UpdatePasswordWithVerificationCode)
	mux.Post("/forgotpassword", app.PasswordResetRequest)
	mux.Post("/verifyemail", app.VerifyEmail)
	mux.Post("/setpassword", app.SetPassword)
	mux.Get("/oauth/{provider}", app.OAuthLogin)
	mux.Get("/oauth/callback/{provider}", app.OAuthCallback)
	mux.Post("/resend-mfa-code", app.ResendMFACode)
	mux.Post("/verify-mfa-code", app.VerifyMFACode)
	mux.Post("/verify-otp", app.VerifyOTP)
	mux.Post("/send-otp", app.SendOTP)
	return mux
}
