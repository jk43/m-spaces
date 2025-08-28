package main

import (
	"fmt"
	"net/http"

	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/utils"
	"github.com/rs/zerolog"
)

var port string = ":80"

type application struct {
	Logger        *zerolog.Logger
	CasbinManager *utils.CasbinManager
}

func main() {

	var app application
	//logger
	app.Logger = molylibs.Logger()
	app.Logger.Info().Msg(fmt.Sprintf("Service starting at port %s", port))

	app.CasbinManager = &utils.CasbinManager{}

	srv := &http.Server{
		Addr:    port,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		app.Logger.Error().Err(err).Msg("")
	}
}
