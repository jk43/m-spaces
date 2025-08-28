package main

import (
	forms "cms/forms"
	"cms/models"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

var port string = ":80"
var grpcAddr string = "0.0.0.0:5000"

type application struct {
	DB            models.DatabaseRepo
	Logger        *zerolog.Logger
	CasbinManager *utils.CasbinManager
	Forms         *forms.Forms
}

type grpcServer struct {
	pb.AuthServiceServer
	app *application
}

type grpcFileClientServer struct {
	pb.FileClientServer
	app *application
}

type CMSHistory struct {
	BoardSlug string `json:"boardSlug"`
	PostSlug  string `json:"postSlug"`
	OrgID     string `json:"organizationId"`
	Action    string `json:"action"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
	IP        string `json:"ip"`
}

func (serv *grpcFileClientServer) GetAddr() string {
	return utils.GetGRPCAddr(utils.FileClientGRPC)
}

// var grpcAddr string = utils.GetGRPCAddr(utils.GeneralGRPC)
var grpcFileClientAddr string = utils.GetGRPCAddr(utils.FileClientGRPC)

func main() {
	var app application

	//logger
	app.Logger = molylibs.Logger()

	//casbin
	app.CasbinManager = &utils.CasbinManager{
		Service:   "",
		PolicyCtx: "",
		ModelCtx:  "",
	}

	ffm := forms.CMSFormElems{
		DB:     app.DB,
		Logger: app.Logger,
		Casbin: app.CasbinManager,
	}

	//forms
	app.Forms = service.NewForms(ffm.FormFuncMap)
	utils.TermDebugging(`app.Forms`, app.Forms)

	//gorm
	mysql, err := molylibs.Mysql()
	mysql.AutoMigrate(
		&models.Board{},
		&models.Setting{},
		&models.Post{},
		&models.File{},
		&models.PostAuthor{},
		&models.PostAttribute{},
		&models.GuestVerificationCode{},
	)
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
