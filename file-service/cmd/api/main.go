package main

import (
	"file/models"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/pb"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

var port string = ":80"
var fileDir string = os.Getenv("FILE_DIR")
var s3Bucket string = os.Getenv("AWS_S3_BUCKET")
var s3Region string = os.Getenv("AWS_S3_REGION")
var awsAccessKey string = os.Getenv("AWS_ACCESS_KEY_ID")
var awsSecretKey string = os.Getenv("AWS_SECRET_ACCESS_KEY")
var awsToken string = os.Getenv("AWS_TOKEN")

type application struct {
	Logger *zerolog.Logger
	DB     models.DatabaseRepo
}

type grpcServer struct {
	pb.FileServiceServer
	app *application
}

var grpcAddr string = "0.0.0.0:5000"

func main() {
	var app application
	//logger
	app.Logger = molylibs.Logger()

	//gorm
	mysql, err := molylibs.Mysql()
	mysql.AutoMigrate(
		models.FileInfo{},
	)
	if err != nil {
		app.Logger.Panic().Err(err).Msg("")
	}
	app.DB = &models.DBRepo{
		Mysql: mysql,
	}

	app.Logger.Info().Msg(fmt.Sprintf("Service starting at port %s", port))

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
		pb.RegisterFileServiceServer(s, &server)
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
