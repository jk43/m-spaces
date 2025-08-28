package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"user/models"

	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/utils"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

var port string = ":80"

type application struct {
	DB         models.DatabaseRepo
	Logger     *zerolog.Logger
	HttpClient *utils.SimpleHttpClient
	GRPCClient GRPCClientInterface
}

type grpcServer struct {
	pb.UserServiceServer
	app *application
}

type grpcFileClientServer struct {
	pb.FileClientServer
	app *application
}

func (serv *grpcFileClientServer) GetAddr() string {
	return utils.GetGRPCAddr(utils.FileClientGRPC)
}

var grpcAddr string = utils.GetGRPCAddr(utils.GeneralGRPC)
var grpcFileClientAddr string = utils.GetGRPCAddr(utils.FileClientGRPC)

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
	if err != nil {
		app.Logger.Panic().Err(err).Msg("")
	}

	go func() {
		var server grpcServer
		fmt.Println("running grpc on ", grpcAddr)
		server = grpcServer{}
		server.app = &app
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("Failed to listen on : %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterUserServiceServer(s, &server)
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("Failed to serve on : %v", err)
		}
		fmt.Println("ended")
	}()
	//file client server
	go func() {
		var server grpcFileClientServer
		fmt.Println("running grpc file client on ", grpcFileClientAddr)
		server = grpcFileClientServer{}
		server.app = &app
		lis, err := net.Listen("tcp", grpcFileClientAddr)
		if err != nil {
			log.Fatalf("Failed to listen on : %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterFileClientServer(s, &server)
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("Failed to serve on : %v", err)
		}
		fmt.Println("ended")
	}()

	//http client for testing
	app.HttpClient = &utils.SimpleHttpClient{}

	//GRPC client for testing
	app.GRPCClient = &GRPCClient{
		app: &app,
	}

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
