package rtm

import (
	"encoding/json"
	"fmt"

	logger "github.com/betchi/zapper"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

var srv = NewServer()

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
		events:  make(map[scpb.EventType]*EventClients),
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
			logger.Info(fmt.Sprintf("Register event userId[%s] roomId[%s] eventType[%d] client[%p]", rcvData.UserId, rcvData.RoomId, rcvData.EventType, rcvData.Client.Conn))
			s.Connection.RemoveEvent(rcvData.UserId, rcvData.EventType, rcvData.Client)

			s.Connection.clients[rcvData.Client] = true
			s.Connection.AddEvent(rcvData.UserId, rcvData.EventType, rcvData.Client)

		case rcvData := <-s.Unregister:
			logger.Info(fmt.Sprintf("Unregister event userId[%s] roomId[%s] eventType[%d] client[%p]", rcvData.UserId, rcvData.RoomId, rcvData.EventType, rcvData.Client.Conn))
			s.Connection.RemoveEvent(rcvData.UserId, rcvData.EventType, rcvData.Client)

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

func (s *server) broadcast(event []byte) {
	var pbEventData scpb.EventData
	json.Unmarshal(event, &pbEventData)

	if _, ok := s.Connection.events[pbEventData.Type]; !ok {
		return
	}

	for _, userID := range pbEventData.UserIDs {
		if user, ok := s.Connection.events[pbEventData.Type].users[userID]; ok {
			for conn := range user.clients {
				if conn == nil {
					continue
				}
				select {
				case conn.Send <- event:
				default:
					close(conn.Send)
					delete(srv.Connection.clients, conn)
				}
			}
		}
	}
}
