package rtm

import (
	"encoding/json"
	"fmt"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
	"github.com/swagchat/rtm-api/logger"
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
		events:  make(map[string]*EventClients),
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
			logger.Info(fmt.Sprintf("Register event userId[%s] roomId[%s] eventName[%s] client[%p]", rcvData.UserId, rcvData.RoomId, rcvData.EventName, rcvData.Client.Conn))
			s.Connection.RemoveEvent(rcvData.UserId, rcvData.EventName, rcvData.Client)

			s.Connection.clients[rcvData.Client] = true
			s.Connection.AddEvent(rcvData.UserId, rcvData.EventName, rcvData.Client)

		case rcvData := <-s.Unregister:
			logger.Info(fmt.Sprintf("Unregister event userId[%s] roomId[%s] eventName[%s] client[%p]", rcvData.UserId, rcvData.RoomId, rcvData.EventName, rcvData.Client.Conn))
			s.Connection.RemoveEvent(rcvData.UserId, rcvData.EventName, rcvData.Client)

		case c := <-s.Close:
			logger.Info("Closing socket")
			s.Connection.RemoveClient(c)
			close(c.Send)

		case message := <-s.Broadcast:
			logger.Info("Broadcasting message")
			s.broadcast(message)
		}
	}
}

func (s *server) broadcast(message []byte) {
	var rtmEvent scpb.Message
	json.Unmarshal(message, &rtmEvent)

	if _, ok := s.Connection.events[rtmEvent.Type]; !ok {
		return
	}

	for _, userID := range rtmEvent.UserIDs {
		if user, ok := s.Connection.events[rtmEvent.Type].users[userID]; ok {
			for conn := range user.clients {
				if conn == nil {
					continue
				}
				select {
				case conn.Send <- rtmEvent.Payload:
				default:
					close(conn.Send)
					delete(srv.Connection.clients, conn)
				}
			}
		}
	}
}
