package main

// import (
// 	"auth/models"
// 	"os"
// 	"testing"

// 	"github.com/moly-space/molylibs"
// )

// var app application
// var server *grpcServer

// func TestMain(m *testing.M) {
// 	os.Setenv("JWT_TOKEN_EXPIRY", "1")
// 	os.Setenv("JWT_REFRESH_TOKEN_EXPIRY", "1")
// 	os.Setenv("REFRESH_TOKEN_COOKIE_DOMAIN", "localhost")
// 	os.Setenv("REFRESH_TOKEN_COOKIE_NAME", "_host_refresh_token")

// 	app.DB = &models.TestDBRepo{}
// 	app.Logger = molylibs.Logger()
// 	server = &grpcServer{
// 		app: &app,
// 	}
// 	os.Exit(m.Run())
// }
