package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins
		return true
	},
}

var sm = NewWebsocketManager()

func (serv *grpcServer) Broadcast(ctx context.Context, in *pb.MessageRequest) (*pb.MessageResponse, error) {
	utils.TermDebugging(`MessageRequest`, in)
	ch := sm.GetChannel(in.Channel)
	for c, _ := range ch {
		err := c.WriteMessage(websocket.TextMessage, in.Message)
		if err != nil {
			fmt.Println("Error writing message:", err)
		}
	}
	return &pb.MessageResponse{Success: true}, nil
}

type ProgressCTX = string
type ProgressID = string
type ProgressMessage struct {
	Context  string `json:"context"`
	Progress int32  `json:"progress"`
	Since    string `json:"since"`
}

var Progresses = make(map[ProgressCTX]map[ProgressID]ProgressMessage)

func (serv *grpcServer) SendProgress(ctx context.Context, in *pb.ProgressMessageRequest) (*pb.MessageResponse, error) {
	ch := sm.GetChannel(in.Channel)
	_, ok := Progresses[in.ProgressCtx]
	if !ok {
		Progresses[in.ProgressCtx] = make(map[ProgressID]ProgressMessage)
	}
	Progresses[in.ProgressCtx][in.Id] = ProgressMessage{
		Context:  in.Progress.Context,
		Progress: in.Progress.Progress,
		Since:    Progresses[in.ProgressCtx][in.Id].Since,
	}
	//add time
	if in.Command == pb.ProgressCommand_ADD {
		Progresses[in.ProgressCtx][in.Id] = ProgressMessage{
			Context:  in.Progress.Context,
			Progress: in.Progress.Progress,
			Since:    time.Now().String(),
		}
	} else {
		Progresses[in.ProgressCtx][in.Id] = ProgressMessage{
			Context:  in.Progress.Context,
			Progress: in.Progress.Progress,
			Since:    Progresses[in.ProgressCtx][in.Id].Since,
		}
	}
	// delete progress
	done := false
	if in.Command == pb.ProgressCommand_DONE {
		delete(Progresses[in.ProgressCtx], in.Id)
		done = true
	}
	doneMessage, _ := json.Marshal(map[string]bool{"reload": true})
	message, _ := json.Marshal(map[string]any{"progress": Progresses[in.ProgressCtx]})
	for c, _ := range ch {
		if done {
			c.WriteMessage(websocket.TextMessage, doneMessage)
		}
		err := c.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			// remove socket if can't send message
			sm.RemoveSocket(in.Channel, c)
		}
	}
	utils.TermDebugging(`Progresses`, Progresses)
	return &pb.MessageResponse{Success: true}, nil
}

func (app *application) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello World")
}

func (app *application) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r, "channel")
	if err != nil {
		app.Logger.Err(err).Send()
		return
	}
	channel := req.Payload.Data["channel"][0]
	var claims *utils.JWTClaims
	token := r.URL.Query().Get("t")
	host, _ := utils.GetXHost(r)
	if token != "" {
		claims, err = utils.GetJWTClaims(token, host)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedJWT, w, err, r)
			return
		}
	} else {
		claims = req.Claims
	}
	var role utils.UserRole
	if claims == nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedJWT, w, err, r)
		return
		role = utils.RoleGuest
	} else {
		role = claims.Role
	}

	csb := utils.NewCasbin(claims.Audience[0], utils.ServiceMessage, "", "")
	perm, _ := csb.GetPermissions(role)
	found := false
	for _, p := range perm {
		if channel == p[2] {
			found = true
			break
		}
	}
	if !found {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedJWT, w, err, r)
		return
	}
	// ok, err := csb.HasLink(role, utils.UserRole(channel))

	// e, _ := csb.GetCasbinEnforcer()
	// e.AddRoleForUser(string(role), string(role))
	// m, _ := e.Enforce(string(role), channel, channel)
	// // m, _ := e.GetRoleManager().HasLink("admin", "member")
	// utils.TermDebugging(`channel`, channel)
	// utils.TermDebugging(`m`, m)
	// utils.TermDebugging(`perm`, perm)

	// if !ok || err != nil {
	// 	utils.TermDebugging(`err`, ok)
	// 	utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedJWT, w, err, r)
	// 	return
	// }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	id := sm.AddSocket(channel, conn, claims)
	// Send id to client
	conn.WriteJSON(utils.Response{Result: utils.SUCCESS, Data: map[string]string{"id": id}})

	for {
		// Read message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			sm.RemoveSocket(channel, conn)
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received: %s\n", message)

		// Write message
		if err := conn.WriteMessage(messageType, message); err != nil {
			sm.RemoveSocket(channel, conn)
			fmt.Println("Error writing message:", err)
			break
		}
	}
}

func (app *application) GetSockets(w http.ResponseWriter, r *http.Request) {

	ch := sm.GetChannel("wowowoowowowowow")
	utils.TermDebugging(`ch`, ch)
	// for c, v := range ch {
	// 	// utils.TermDebugging(`c`, c)
	// 	// utils.TermDebugging(`v`, v)
	// 	err := c.WriteMessage(websocket.TextMessage, []byte("Hello"))
	// 	if err != nil {
	// 		fmt.Println("Error writing message:", err)
	// 	}
	// }
	fmt.Println("GetSockets")
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: ch})
}
