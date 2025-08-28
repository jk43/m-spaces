package main

import (
	"fmt"
	"hook/models"
	"log"
	"net"
	"net/http"

	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/utils"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

var port string = ":80"

type application struct {
	Logger       *zerolog.Logger
	DB           models.DatabaseRepo
	HookHandlers *utils.HookHandlerList
}

type grpcServer struct {
	pb.HookServiceServer
	app *application
}

var grpcAddr string = "0.0.0.0:5000"

func main() {
	var app application
	//logger
	app.Logger = molylibs.Logger()

	//mongodb
	mongoClient, err := molylibs.Mongo()
	if err != nil {
		app.Logger.Panic().Err(err).Msg("")
	}
	app.DB = &models.DBRepo{
		Mongo: mongoClient,
	}
	app.Logger.Info().Msg(fmt.Sprintf("Service starting at port %s", port))

	//hook handlers
	app.HookHandlers = utils.NewHookHandlerList()
	app.HookHandlers.Setup(app.setHookHandlers)

	//grpc
	go func() {
		server := grpcServer{}
		server.app = &app
		app.Logger.Info().Msg(fmt.Sprintf("GRPC starting at address %s", grpcAddr))
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("Failed to listen on : %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterHookServiceServer(s, &server)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve on : %v", err)
		}
	}()

	srv := &http.Server{
		Addr:    port,
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		app.Logger.Error().Err(err).Msg("")
	}
}
