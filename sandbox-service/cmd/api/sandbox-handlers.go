package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/kafka"
	"github.com/moly-space/molylibs/utils"
)

func (a *application) Test(w http.ResponseWriter, r *http.Request) {
	a.Logger.Debug().Msg("Gogog")

	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: os.Getenv("REDIS_ADDR")})
	return
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
}

func (app *application) PostRedis(w http.ResponseWriter, r *http.Request) {
	// p := Person{
	// 	"Alex Kim",
	// 	12,
	// 	"Male",
	// }
	//app.Redis.Client.Set(context.Background(), "person", "YOYO", 0)
	// err := app.Redis.Set(context.Background(), utils.GetRedisKey("orgId", "service", "method", "userId"), p, 0)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: os.Getenv("REDIS_ADDR")})
}

func (app *application) GetRedis(w http.ResponseWriter, r *http.Request) {
	// v, err := app.Redis.Get(context.Background(), "My-Redis").Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	p := Person{}
	//ok, err := app.Redis.Get(context.Background(), utils.GetRedisKey("orgId", "service", "method", "userId"), &p)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// if !ok {
	// 	fmt.Println("no data")
	// }
	fmt.Println(p.Name)
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: p})
}

func (app *application) DeleteRedis(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: os.Getenv("REDIS_ADDR")})
}

type Person2 struct {
	LastName  string
	FirstName string
	Email     string
}

func (app *application) Sandbox(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r)
	ef := utils.NewCasbin(req.Host, utils.ServiceBroker, "", "")
	ef.GetCasbinEnforcer()
	//utils.TermDebugging(`app.Enforcer.GetCasbinEnforcer()`, e)

	// rds := molylibs.Redis[Person2]{}
	person := Person2{
		"Kim",
		"Alex",
		"jk@jktech.net",
	}
	// //rds.Data = &person
	// rds.SetKey("person", "alex", "gogo")
	// err := rds.SetKey("person", "alex", "gogo11")
	// if err != nil {
	// 	fmt.Println("error", err)
	// }
	// err = rds.Set(r.Context(), &person, 1)
	// if err != nil {
	// 	fmt.Println("error", err)
	// }

	rds := molylibs.Redis[Person2]{DBNumber: 1}
	//rds.Data = &person
	rds.SetKey("person", "alex", "gogo")
	cached, err := rds.Get(r.Context())
	if err == nil {
		fmt.Println("cachedcachedcached", cached)
		return
	}
	err = rds.Set(r.Context(), &person)
	if err != nil {
		fmt.Println("error from set", err)
	}
	// rds := molylibs.MolyRedis[Person2]{DBNumber: 0, Key: "person", Data: &person}
	// rds.Set(r.Context(), 10)

	// test()

	// return

	// fromCache, err := rds.Get(r.Context())

	// if err != nil {
	// 	fmt.Println("error", err)
	// }

	// fmt.Println("fromCachefromCachefromCache", fromCache)
	molylibs.GetRedisClientInfo("")
}

func test() {
	// 	rds := molylibs.MolyRedis[Person2]{DBNumber: 0, Key: "person", Data: nil}
	// 	fromCache, err := rds.Get(context.Background())

	// 	if err != nil {
	// 		fmt.Println("error", err)
	// 	}

	//		fmt.Println("fromCachefromCachefromCache", fromCache)
	//	}
}

func (app *application) RbacManager(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: err.Error()})
		return
	}
	app.CasbinManager.GetEnforcer(req.Host)
	app.CasbinManager.GetEnforcer(req.Host)
	app.CasbinManager.GetEnforcer(req.Host)
	e, _ := app.CasbinManager.GetEnforcer(req.Host)
	utils.TermDebugging(`app.CasbinManager.GetEnforcer(req.Host)`, e.GetAllRoles())

	// e, err := es.GetEnforcer(req.Host)
	// if err != nil {
	// 	utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: err.Error()})
	// 	return
	// }
	// roles := e.GetAllRoles()
	// if err != nil {
	// 	utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: err.Error()})
	// 	return
	// }
	// ewh := es.GetCasbin(req.Host)
	// p, err := ewh.GetCasbinModeAndPolicy()
	// if err != nil {
	// 	utils.TermDebugging(`err`, err)
	// 	utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: err})
	// 	// return
	// }
	// utils.TermDebugging(`p`, p)
	//utils.TermDebugging(`es.GetCasbin(req.Host)`, es.GetCasbin(req.Host))
	//utils.TermDebugging(`roles`, e.GetAllRoles())
	// es.GetEnforcer(req.Host)
	// utils.TermDebugging(`es.Casbins[req.Host].Enforcer`, &es.Casbins[req.Host].Enforcer)
	// utils.TermDebugging(`es.Casbins[req.Host].Enforcer`, &es.Casbins[req.Host].Enforcer)
	// utils.TermDebugging(`es.Casbins[req.Host].Enforcer`, es.Casbins[req.Host].Enforcer.GetAllRoles())

	// c := utils.NewCasbin(req.Host, "", "", "")
	// e, err := c.GetCasbinEnforcer()
	// if err != nil {
	// 	utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: err.Error()})
	// 	return
	// }
	// utils.TermDebugging(`e`, e.GetAllRoles())
	// c.GetCasbinEnforcer()
	// c.GetCasbinEnforcer()
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

type CMSHistory struct {
	BoardSlug string `json:"boardSlug"`
	PostSlug  string `json:"postSlug"`
	Action    string `json:"action"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
}

func (app *application) KafkaHistory(w http.ResponseWriter, r *http.Request) {
	history := kafka.NewMessage(CMSHistory{
		BoardSlug: "board.Slug",
		PostSlug:  "post.Slug",
		Action:    "PutPosts",
		UserID:    "req.UserID",
	}, kafka.TopicHistory)
	kafka.SendMessage(history)
}
