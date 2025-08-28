package main

import (
	"broker/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/utils"
	"github.com/rs/zerolog"
)

var port string = ":80"
var routerRules = make(map[string]*models.Rule)
var routerRulesWithStar = make(map[string]*models.Rule)

type application struct {
	DB     *models.DBRepo
	Logger *zerolog.Logger
}

var rbacPolicies *pb.RBACResponse
var authEnforcer *casbin.Enforcer
var hooks *utils.Hooks

func main() {

	var app application
	//logger
	app.Logger = molylibs.Logger()
	//mongodb
	mongoClient, err := molylibs.Mongo()
	//load rbac policies
	if err != nil {
		app.Logger.Panic().Err(err).Msg("")
	}
	app.DB = &models.DBRepo{
		Mongo: mongoClient,
	}
	retry := 5
	tried := 0
	//will prevent service not start
	csb := utils.Casbin{
		Service:   "broker",
		Host:      "",
		PolicyCtx: "",
		ModelCtx:  "",
	}
	for {
		authEnforcer, err = csb.GetCasbinEnforcer()
		if retry < tried {
			break
		}
		if err != nil {
			app.Logger.Error().Err(err).Msg("Failed to connect to rbac-service")
			fmt.Println("Trying to reconnect rbac-service")
			time.Sleep(5 * time.Second)
			tried++
			continue
		}
		break
	}
	tried = 0
	for {
		hooks, err = utils.NewHookList(os.Getenv("BROKER_SERVICE_ADDR"))
		if retry < tried {
			break
		}
		if err != nil {
			app.Logger.Error().Err(err).Msg("Failed to connect to hook-service")
			fmt.Println("Trying to reconnect hook-service")
			time.Sleep(5 * time.Second)
			tried++
			continue
		}
		break
	}

	app.Logger.Info().Msg(fmt.Sprintf("Service starting at port %s", port))

	//loading broker rules
	err = app.DB.GetRules(routerRules, routerRulesWithStar)
	if err != nil {
		app.Logger.Panic().Err(err).Msg("")
	}
	app.Logger.Info().Msg(fmt.Sprintln(routerRules))

	srv := &http.Server{
		Addr:    port,
		Handler: app.router(),
	}

	err = srv.ListenAndServeTLS("certs/hotdev.com.crt", "certs/hotdev.com.key")
	if err != nil {
		app.Logger.Error().Err(err).Msg("")
	}
}
