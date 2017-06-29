package services

import (
	"encoding/json"
	"log"

	"github.com/fairway-corp/swagchat-realtime/models"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	connections map[*Conn]bool
	rooms       map[string]RoomConnections
	Broadcast   chan []byte
	Register    chan *RegisterData
	Unregister  chan *RegisterData
	Close       chan *Conn
}

var Manager = Hub{
	connections: make(map[*Conn]bool),
	rooms:       make(map[string]RoomConnections),
	Broadcast:   make(chan []byte),
	Register:    make(chan *RegisterData),
	Unregister:  make(chan *RegisterData),
	Close:       make(chan *Conn),
}

type RoomConnections struct {
	connections map[string]EventConnections
}

type EventConnections struct {
	connections map[*Conn]bool
}

func (h *Hub) Run() {
	for {
		select {
		case registerData := <-h.Register:
			log.Printf("<-h.Register roomId[%s] eventName[%s]", registerData.roomId, registerData.eventName)
			h.connections[registerData.conn] = true

			if registerData.roomId == "" || registerData.eventName == "" {
				h.connectionInfo(registerData)
				continue
			}

			if _, ok := h.rooms[registerData.roomId]; !ok {
				roomConnections := RoomConnections{
					connections: make(map[string]EventConnections),
				}
				eventConnections := EventConnections{
					connections: make(map[*Conn]bool),
				}
				eventConnections.connections[registerData.conn] = true
				roomConnections.connections[registerData.eventName] = eventConnections
				h.rooms[registerData.roomId] = roomConnections
			} else {
				if _, ok := h.rooms[registerData.roomId].connections[registerData.eventName]; !ok {
					eventConnections := EventConnections{
						connections: make(map[*Conn]bool),
					}
					eventConnections.connections[registerData.conn] = true
					h.rooms[registerData.roomId].connections[registerData.eventName] = eventConnections
				}
				h.rooms[registerData.roomId].connections[registerData.eventName].connections[registerData.conn] = true
			}

			h.connectionInfo(registerData)
		case registerData := <-h.Unregister:
			log.Println("<-h.Unregister", registerData.roomId, registerData.eventName)
			if _, ok := h.rooms[registerData.roomId].connections[registerData.eventName].connections[registerData.conn]; ok {
				delete(h.rooms[registerData.roomId].connections[registerData.eventName].connections, registerData.conn)
			}
			h.connectionInfo(registerData)
		case conn := <-h.Close:
			log.Println("<-h.Close")
			for _, roomId := range conn.RoomIds {
				if _, ok := h.rooms[roomId]; ok {
					delete(h.rooms, roomId)
				}
			}
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
				close(conn.Send)
			}
		case message := <-h.Broadcast:
			log.Println("<-h.broadcast")
			var messageMap models.Message
			newMessage := message[1 : len(message)-1] // Delete preceding and following double quotes
			json.Unmarshal(newMessage, &messageMap)
			log.Printf("%#v", messageMap)

			if messageMap.Type == "text" {
				var payloadText models.PayloadText
				json.Unmarshal(messageMap.Payload, &payloadText)
			}

			for conn := range h.rooms[messageMap.RoomId].connections[messageMap.EventName].connections {
				select {
				case conn.Send <- newMessage:
				default:
					close(conn.Send)
					delete(Manager.connections, conn)
				}
			}
		}
	}
}

func (h *Hub) connectionInfo(registerData *RegisterData) {
	log.Println("----------- connection info start ------------------------------------------")
	log.Println("-  conn", registerData.conn)
	log.Println("-  connections count", len(h.connections))
	log.Printf("-  h.room[%s].connections[%s].connections count %d", registerData.roomId, registerData.eventName, len(h.rooms[registerData.roomId].connections[registerData.eventName].connections))
	log.Println("----------- connection info end --------------------------------------------")
}
