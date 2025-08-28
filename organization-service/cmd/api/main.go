package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"organization/models"

	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/pb"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

var port string = ":80"
var grpcAddr string = "0.0.0.0:5000"

type application struct {
	Logger *zerolog.Logger
	DB     models.DatabaseRepo
}

type grpcServer struct {
	pb.OrganizationServiceServer
	app  *application
	name string
}

func main() {
	var app application
	app.Logger = molylibs.Logger()
	mongoClient, err := molylibs.Mongo()
	if err != nil {
		app.Logger.Panic().Err(err).Msg("")
	}
	app.DB = &models.DBRepo{
		Mongo: mongoClient,
	}
	go func() {
		server := grpcServer{}
		server.app = &app
		app.Logger.Info().Msg(fmt.Sprintf("GRPC starting at address %s", grpcAddr))
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("Failed to listen on : %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterOrganizationServiceServer(s, &server)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve on : %v", err)
		}
		app.Logger.Info().Msg(fmt.Sprintf("GRPC ending at address %s", grpcAddr))
	}()
	app.Logger.Info().Msg(fmt.Sprintf("Service starting at port %s", port))
	srv := &http.Server{
		Addr:    port,
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
