package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/moly-space/molylibs/utils"
)

type Channel = string

//type Socket = map[*websocket.Conn]*utils.JWTClaims

type WebsocketInfo struct {
	Claims *utils.JWTClaims
	Id     string // to be access to index
}

// Index of socket
type WebSocketIndexInfo struct {
	Channel Channel
	Socket  *websocket.Conn
}

type SocketManager struct {
	sockets     map[Channel]map[*websocket.Conn]WebsocketInfo
	socketIndex map[string]WebSocketIndexInfo
}

func (sm *SocketManager) AddSocket(channel Channel, conn *websocket.Conn, claims *utils.JWTClaims) string {
	id := uuid.New().String()
	if _, ok := sm.sockets[channel]; !ok {
		sm.sockets[channel] = make(map[*websocket.Conn]WebsocketInfo)
	}
	sm.sockets[channel][conn] = WebsocketInfo{
		Claims: claims,
		Id:     id,
	}
	sm.socketIndex[id] = WebSocketIndexInfo{
		Channel: channel,
		Socket:  conn,
	}
	return id
}

func (sm *SocketManager) RemoveSocket(channel Channel, conn *websocket.Conn) {
	if s, ok := sm.sockets[channel]; ok {
		id := s[conn].Id
		delete(sm.socketIndex, id)
		delete(sm.sockets[channel], conn)
	}
}

func (sm *SocketManager) GetChannel(channel Channel) map[*websocket.Conn]WebsocketInfo {
	return sm.sockets[channel]
}

func NewWebsocketManager() *SocketManager {
	return &SocketManager{
		sockets:     make(map[Channel]map[*websocket.Conn]WebsocketInfo),
		socketIndex: make(map[string]WebSocketIndexInfo),
	}
}
