package main

import (
	"fmt"
	"log"
	"net/http"
	"tree/models"

	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/pb"
	"github.com/rs/zerolog"
)

var port string = ":80"
var grpcAddr string = "0.0.0.0:5000"

type application struct {
	DB     models.DatabaseRepo
	Logger *zerolog.Logger
}

type grpcServer struct {
	pb.AuthServiceServer
	app *application
}

func main() {
	var app application

	//logger
	app.Logger = molylibs.Logger()

	//gorm
	mysql, err := molylibs.Mysql()
	mysql.AutoMigrate()
	if err != nil {
		app.Logger.Panic().Err(err).Msg("")
	}
	app.DB = &models.DBRepo{
		Mysql: mysql,
	}

	//grpc
	// go func() {
	// 	server := grpcServer{}
	// 	server.app = &app
	// 	app.Logger.Info().Msg(fmt.Sprintf("GRPC starting at address %s", grpcAddr))
	// 	lis, err := net.Listen("tcp", grpcAddr)
	// 	if err != nil {
	// 		log.Fatalf("Failed to listen on : %v", err)
	// 	}
	// 	s := grpc.NewServer()
	// 	pb.RegisterTreeServiceServer(s, &server)
	// 	if err := s.Serve(lis); err != nil {
	// 		log.Fatalf("Failed to serve on : %v", err)
	// 	}
	// }()

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
