package main

import (
	"fmt"
	"log"
	"net"
	"rbac/models"

	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/pb"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type application struct {
	DB     models.DatabaseRepo
	Logger *zerolog.Logger
}

type grpcServer struct {
	pb.RBACServiceServer
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
	server := grpcServer{}
	server.app = &app
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen on : %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRBACServiceServer(s, &server)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve on : %v", err)
	}
	fmt.Println("ended")
}
