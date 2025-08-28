package main

import (
	"broker/models"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/moly-space/molylibs/utils"
)

func (app *application) SimpleHttpRelay(w http.ResponseWriter, r *http.Request) {
	// Start pre hook
	host, err := utils.GetXHost(r)
	if err != nil {
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("")
		return
	}
	tag := strings.ToLower(host + "/" + r.Method + r.RequestURI)
	err = hooks.RunPreHooks(tag, r)
	if err != nil {
		handlerErro, ok := err.(*utils.HookHandlerError)
		if ok {
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write(handlerErro.HandlerError)
			return
		}
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("")
		utils.WriteJSON(w, http.StatusExpectationFailed,
			utils.Response{
				Result: utils.ERROR,
				Data: []utils.ErrorDetails{{
					Error: "Unable to process request",
				}},
			})
		return
	}
	// End pre hook
	rule, err := findRule(r.URL.Path, r.Method)
	if err != nil {
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("")
		return
	}

	app.Logger.Info().Str("method", rule.Method).Str("url", rule.URL).Str("requestURI", r.RequestURI).Msg("SimpleHttpRelay was invoked")
	req, err := http.NewRequest(rule.Method, rule.URL, r.Body)

	if err != nil {
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("")
		return
	}
	req.Header = r.Header
	req.URL.RawQuery = r.URL.RawQuery

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("Error request http")
		w.WriteHeader(http.StatusBadGateway)
		w.Write(nil)
		return
	}

	// post hook
	err = hooks.RunPostHooks(tag, r, response)
	if err != nil {
		handlerErro, ok := err.(*utils.HookHandlerError)
		if ok {
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write(handlerErro.HandlerError)
			return
		}
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("")
		utils.WriteJSON(w, http.StatusExpectationFailed,
			utils.Response{
				Result: utils.ERROR,
				Data: []utils.ErrorDetails{{
					Error: "Unable to process request",
				}, {
					Error: "not good",
				}},
			})
		return
	}
	// post hook
	body, err := io.ReadAll(response.Body)
	if err != nil {
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("Error parsing body")
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write(nil)
		return
	}

	for _, cookie := range response.Cookies() {
		http.SetCookie(w, cookie)
	}

	// todo: get the actual header from the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	w.Write(body)
	return
}

func (app *application) RawHttpRelay(w http.ResponseWriter, r *http.Request) {
	rule, err := findRule(r.URL.Path, r.Method)
	if err != nil {
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("")
		return
	}

	// Create new request
	req, err := http.NewRequest(r.Method, rule.URL, r.Body)
	if err != nil {
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("Error creating request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Copy headers and query parameters from original request
	req.Header = r.Header
	req.URL.RawQuery = r.URL.RawQuery

	// Create HTTP client and send request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("Error forwarding request")
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer response.Body.Close()

	// Copy response headers
	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy response cookies
	for _, cookie := range response.Cookies() {
		http.SetCookie(w, cookie)
	}

	// Set status code
	w.WriteHeader(response.StatusCode)

	// Copy response body
	_, err = io.Copy(w, response.Body)
	if err != nil {
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("Error copying response")
		return
	}
}

// WebSocketHandler is a handler for websocket connections

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func relayMessages(ctx context.Context, clientConn, relayConn *websocket.Conn) {
	defer clientConn.Close()
	defer relayConn.Close()

	// Relay messages from clientConn to relayConn
	go func() {
		defer clientConn.Close()
		defer relayConn.Close()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				messageType, message, err := clientConn.ReadMessage()
				if err != nil {
					log.Printf("error reading from client: %v", err)
					return
				}
				log.Printf("Received from client: %s", message)
				err = relayConn.WriteMessage(messageType, message)
				if err != nil {
					log.Printf("error writing to Python server: %v", err)
					return
				}
			}
		}
	}()

	// Relay messages from relayConn to clientConn
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done foe")
			return
		default:
			messageType, message, err := relayConn.ReadMessage()
			if err != nil {
				log.Printf("error reading from Python server: %v", err)
				return
			}
			log.Printf("Received from Python server: %s", message)
			err = clientConn.WriteMessage(messageType, message)
			if err != nil {
				log.Printf("error writing to client: %v", err)
				return
			}
		}
	}
}

func (app *application) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	client, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("upgrader.Upgrade")
	}
	rule, err := findRule(r.URL.Path, r.Method)
	if err != nil {
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("")
		return
	}
	// Send jwt token to the relay server
	relayHeader := http.Header{}
	relayHeader.Add("Authorization", r.Header.Get("Authorization"))

	server, _, err := websocket.DefaultDialer.Dial(rule.URL+"?"+r.URL.RawQuery, relayHeader)
	app.Logger.Info().Str("rule.URL", rule.URL).Msg("WebSocketHandler was invoked")
	if err != nil {
		client.Close()
		app.Logger.Error().Err(err).Str("uri", r.RequestURI).Msg("websocket.DefaultDialer.Dial")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer server.Close()
	defer client.Close()
	defer cancel()

	relayMessages(ctx, client, server)
}

// list broker list for debugging
func (a *application) Routes(w http.ResponseWriter, r *http.Request) {
	utils.TermDebugging(`rbacPolicies`, rbacPolicies)
	rules := map[string]map[string]*models.Rule{
		"rules":         routerRules,
		"rulesWithStar": routerRulesWithStar,
		"rbac":          nil,
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: rules})
	return
}
