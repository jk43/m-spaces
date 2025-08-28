// import (
// 	"errors"
// 	"net/http"

// 	"github.com/casbin/casbin/v2"
// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/cors"
// 	"github.com/moly-space/molylibs/utils"
// )

// func (app *application) routes(mux *chi.Mux) *chi.Mux {
// 	// mux.Post("/user/*", app.SimpleHttpRelay)
// 	// mux.Get("/user/*", app.SimpleHttpRelay)
// 	// mux.Post("/sandbox/*", app.SimpleHttpRelay)
// 	// mux.Get("/sandbox/*", app.SimpleHttpRelay)
// 	mux.Get("/list", app.List)
// 	mux.Get("/*", app.SimpleHttpRelay)
// 	mux.Post("/*", app.SimpleHttpRelay)
// 	mux.Put("/*", app.SimpleHttpRelay)
// 	mux.Delete("/*", app.SimpleHttpRelay)
// 	return mux
// }

// func (app *application) router() http.Handler {
// 	mux := chi.NewRouter()
// 	mux.Use(cors.Handler(cors.Options{
// 		AllowedOrigins: []string{"https://*", "http://*"},
// 		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
// 		//AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
// 		AllowedHeaders:   []string{"*"},
// 		ExposedHeaders:   []string{"Link"},
// 		AllowCredentials: true,
// 		MaxAge:           300,
// 	}))
// 	mux.Use(utils.ParseJWTClaims)
// 	mux.Use(Authorizer(authEnforcer))

// 	return app.routes(mux)
// }

// func Authorizer(e *casbin.Enforcer) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			claims := utils.GetClaims(r)
// 			var role utils.UserRole
// 			if claims.Role == "" {
// 				role = utils.RoleGuest
// 			} else {
// 				role = claims.Role
// 			}
// 			method := r.Method
// 			path := r.URL.Path
// 			res, err := e.Enforce(string(role), path, method)
// 			if err != nil {
// 				utils.NewPreDefinedHttpError(utils.InternalServerErr, http.StatusInternalServerError, w, err)
// 				return
// 			}
// 			if !res {
// 				err := errors.New("res is empty on Casbin")
// 				utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, http.StatusForbidden, w, err)
// 				return
// 			}
// 			next.ServeHTTP(w, r)
// 			return
// 		})
// 	}
// }

package main

import (
	"broker/models"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/moly-space/molylibs/utils"
)

func (app *application) routes(mux *chi.Mux) *chi.Mux {
	mux.Get("/routes", app.Routes)
	mux.Get("/ext/*", app.RawHttpRelay)
	mux.Get("/*", app.SimpleHttpRelay)
	mux.Post("/*", app.SimpleHttpRelay)
	mux.Put("/*", app.SimpleHttpRelay)
	mux.Delete("/*", app.SimpleHttpRelay)
	// websocket
	mux.Get("/ws/*", app.WebSocketHandler)
	return mux
}

func (app *application) router() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		//AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(utils.ParseJWTClaims)
	mux.Use(Authorizer(authEnforcer))

	return app.routes(mux)
}

func Authorizer(e *casbin.Enforcer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := utils.GetClaims(r)
			var role utils.UserRole
			if claims == nil {
				role = utils.RoleGuest
			} else {
				role = claims.Role
			}
			method := r.Method
			path := r.URL.Path
			res, err := e.Enforce(string(role), path, method)
			if err != nil {
				utils.NewPreDefinedHttpError(utils.InternalServerErr, http.StatusInternalServerError, w, err)
				return
			}
			if res {
				next.ServeHTTP(w, r)
				return
			}
			//try to find with star
			rule := findRuleWithStar(method, path)
			if rule == nil {
				fmt.Println("method", method, "path", path)
				err = errors.New("failed to enforce the star rule")
				utils.NewPreDefinedHttpError(utils.UnauthorizedRequestFromBroker, http.StatusForbidden, w, err)
				return
			}
			res, err = e.Enforce(string(role), rule.Path, method)
			if err != nil {
				utils.NewPreDefinedHttpError(utils.InternalServerErr, http.StatusInternalServerError, w, err)
				return
			}
			if !res {
				err := errors.New("res is empty on Casbin")
				utils.NewPreDefinedHttpError(utils.UnauthorizedRequestFromBroker, http.StatusForbidden, w, err)
				return
			}
			next.ServeHTTP(w, r)
			return
		})
	}
}

func findRuleWithStar(method string, path string) *models.Rule {
	key := method + "|" + path
	for k, v := range routerRulesWithStar {
		if strings.HasPrefix(key, k) {
			newRule := &models.Rule{
				Method: method,
				Path:   strings.Replace(k, method+"|", "", 1),
				URL:    v.URL + "/" + strings.Replace(path, strings.TrimSuffix(v.Path, "*"), "", 1),
			}
			return newRule
		}
	}
	return nil
}

func findRule(path string, method string) (*models.Rule, error) {
	key := method + "|" + path
	rule, ok := routerRules[key]
	if ok {
		return rule, nil
	}
	rule = findRuleWithStar(method, path)
	if rule != nil {
		return rule, nil
	}
	methodsToSearch := []string{
		"*",
		method,
	}
	uriParts := strings.Split(path, "/")
	found := false
	for i, _ := range uriParts {
		for _, m := range methodsToSearch {
			keyToFind := m + "|" + strings.Join(uriParts[0:i+1], "/") + "/*"
			rule, ok = routerRules[keyToFind]
			if ok {
				ruleURI := strings.Replace(rule.Path, "*", "", 1)
				newURL := rule.URL + strings.Replace(path, ruleURI, "", 1)
				found = true
				rule = &models.Rule{
					Method: method,
					Path:   path,
					URL:    newURL,
				}
				routerRules[method+"|"+path] = rule
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("%s is not found on rules", key)
	}
	utils.TermDebugging("rule", rule)
	return rule, nil
}
