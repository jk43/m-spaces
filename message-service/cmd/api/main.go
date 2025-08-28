package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/pb"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

var port string = ":80"
var grpcAddr string = "0.0.0.0:5000"

type application struct {
	Logger *zerolog.Logger
}

type grpcServer struct {
	pb.MessageServiceServer
	app *application
}

func main() {
	var err error
	var app application
	//logger
	app.Logger = molylibs.Logger()

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
		pb.RegisterMessageServiceServer(s, &server)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve on : %v", err)
		}
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
