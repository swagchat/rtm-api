package rtm

import (
	"encoding/json"
	"fmt"

	"github.com/swagchat/rtm-api/logging"
	"github.com/swagchat/rtm-api/models"
	"go.uber.org/zap/zapcore"
)

var srv *server = NewServer()

type server struct {
	Connection *Connection
	Broadcast  chan []byte
	Register   chan *RcvData
	Unregister chan *RcvData
	Close      chan *Client
}

func NewServer() *server {
	conn := &Connection{
		clients: make(map[*Client]bool),
		users:   make(map[string]*UserClients),
		rooms:   make(map[string]*RoomClients),
	}

	srv := &server{
		Connection: conn,
		Broadcast:  make(chan []byte),
		Register:   make(chan *RcvData),
		Unregister: make(chan *RcvData),
		Close:      make(chan *Client),
	}
	return srv
}

func Server() *server {
	return srv
}

func (s *server) Run() {
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

func (s *server) broadcast(message []byte) {
	var messageMap models.Message
	json.Unmarshal(message, &messageMap)
	if messageMap.Type == "text" {
		var payloadText models.PayloadText
		json.Unmarshal(messageMap.Payload, &payloadText)
	}

	if roomClient, ok := s.Connection.rooms[messageMap.RoomId]; ok {
		if roomClient.roomUsers != nil {
			for _, roomUser := range s.Connection.rooms[messageMap.RoomId].roomUsers {
				for conn := range roomUser.events[messageMap.EventName].clients {
					select {
					case conn.Send <- message:
					default:
						close(conn.Send)
						delete(srv.Connection.clients, conn)
					}
				}
			}
		}
	}
}
