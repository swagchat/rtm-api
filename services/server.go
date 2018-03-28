package services

import (
	"encoding/json"
	"fmt"

	"github.com/swagchat/rtm-api/logging"
	"github.com/swagchat/rtm-api/models"
	"go.uber.org/zap/zapcore"
)

var Srv *Server = NewServer()

type Server struct {
	Connection *Connection
	Broadcast  chan []byte
	Register   chan *RcvData
	Unregister chan *RcvData
	Close      chan *Client
}

func NewServer() *Server {
	conn := &Connection{
		clients: make(map[*Client]bool),
		users:   make(map[string]*UserClients),
		rooms:   make(map[string]*RoomClients),
	}

	srv := &Server{
		Connection: conn,
		Broadcast:  make(chan []byte),
		Register:   make(chan *RcvData),
		Unregister: make(chan *RcvData),
		Close:      make(chan *Client),
	}
	return srv
}

func GetServer() *Server {
	return Srv
}

func (s *Server) Run() {
	for {
		select {
		case rcvData := <-s.Register:
			// Register event
			logging.Log(zapcore.InfoLevel, &logging.AppLog{
				Kind:   "register",
				UserID: rcvData.UserId,
				RoomID: rcvData.RoomId,
				Event:  rcvData.EventName,
				Client: fmt.Sprintf("%p", rcvData.Client.Conn),
			})

			s.Connection.clients[rcvData.Client] = true
			s.Connection.AddEvent(rcvData.UserId, rcvData.RoomId, rcvData.EventName, rcvData.Client)

		case rcvData := <-s.Unregister:
			// Unregister event
			logging.Log(zapcore.InfoLevel, &logging.AppLog{
				Kind:   "unregister",
				UserID: rcvData.UserId,
				RoomID: rcvData.RoomId,
				Event:  rcvData.EventName,
				Client: fmt.Sprintf("%p", rcvData.Client.Conn),
			})
			s.Connection.RemoveEvent(rcvData.UserId, rcvData.RoomId, rcvData.EventName, rcvData.Client)

		case c := <-s.Close:
			// Socket close
			s.Connection.RemoveClient(c)
			close(c.Send)

		case message := <-s.Broadcast:
			// Broadcast message
			logging.Log(zapcore.InfoLevel, &logging.AppLog{
				Kind:    "bloadcast",
				Message: string(message),
			})
			s.broadcast(message)
		}
	}
}

func (s *Server) broadcast(message []byte) {
	var messageMap models.Message
	json.Unmarshal(message, &messageMap)
	if messageMap.Type == "text" {
		var payloadText models.PayloadText
		json.Unmarshal(messageMap.Payload, &payloadText)
	}

	for _, roomUser := range s.Connection.rooms[messageMap.RoomId].roomUsers {
		for conn, _ := range roomUser.events[messageMap.EventName].clients {
			select {
			case conn.Send <- message:
			default:
				close(conn.Send)
				delete(Srv.Connection.clients, conn)
			}
		}
	}
}
