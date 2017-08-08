package services

import (
	"log"
	"os"
)

type Connection struct {
	clients map[*Client]bool
	users   map[string]UserClients
	rooms   map[string]RoomClients
}

type UserClients struct {
	clients map[*Client]bool
}

type RoomClients struct {
	clients map[string]RoomUserClients
}

type RoomUserClients struct {
	clients map[string]EventClients
}

type EventClients struct {
	clients map[*Client]bool
}

func (con *Connection) AddClient(c *Client) {
	if c == nil {
		return
	}

	hostname, _ := os.Hostname()
	log.Printf("[WS-INFO][%s] ADD CLIENT %p", hostname, c)

	con.clients[c] = true

	var userClients UserClients
	if _, ok := con.users[c.UserId]; ok {
		con.users[c.UserId].clients[c] = true
	} else {
		userClients = UserClients{
			clients: make(map[*Client]bool),
		}
		userClients.clients[c] = true
		con.users[c.UserId] = userClients
	}
}

func (con *Connection) Info() {
	hostname, _ := os.Hostname()
	log.Printf("[WS-INFO][%s] All Clients %d", hostname, len(con.clients))
	for userId, _ := range con.users {
		log.Printf("[WS-INFO][%s] UserId[%s] %d", hostname, userId, len(con.users[userId].clients))
	}
	for roomId, _ := range con.rooms {
		log.Printf("[WS-INFO][%s] RoomId[%s] %d", hostname, roomId, len(con.rooms[roomId].clients))
		for userId, _ := range con.rooms[roomId].clients {
			log.Printf("[WS-INFO][%s] RoomId[%s][%s] %d", hostname, roomId, userId, len(con.rooms[roomId].clients[userId].clients))
			for eventName, _ := range con.rooms[roomId].clients[userId].clients {
				log.Printf("[WS-INFO][%s] RoomId[%s][%s][%s] %d", hostname, roomId, userId, eventName, len(con.rooms[roomId].clients[userId].clients[eventName].clients))
			}
		}
	}
}

func (con *Connection) RemoveClient(c *Client) {
	if c == nil {
		return
	}

	delete(con.clients, c)

	for client, _ := range con.users[c.UserId].clients {
		if client == c {
			delete(con.users[c.UserId].clients, c)
		}
	}
	if len(con.users[c.UserId].clients) == 0 {
		delete(con.users, c.UserId)
	}

	for roomId, roomClients := range con.rooms {
		for userId, roomUserClients := range roomClients.clients {
			for eventName, eventClients := range roomUserClients.clients {
				for client, _ := range eventClients.clients {
					if client == c {
						delete(eventClients.clients, c)
					}
				}
				if len(roomUserClients.clients[eventName].clients) == 0 {
					delete(roomUserClients.clients, eventName)
				}
			}
			if len(roomClients.clients[userId].clients) == 0 {
				delete(roomClients.clients, userId)
			}
		}
		if len(con.rooms[roomId].clients) == 0 {
			delete(con.rooms, roomId)
		}
	}

}

func (con *Connection) AddEvent(userId, roomId, eventName string, c *Client) {
	if userId == "" || roomId == "" || eventName == "" || c == nil {
		return
	}

	if _, ok := con.rooms[roomId]; !ok {
		rc := RoomClients{
			clients: make(map[string]RoomUserClients),
		}
		ruc := RoomUserClients{
			clients: make(map[string]EventClients),
		}
		ec := EventClients{
			clients: make(map[*Client]bool),
		}
		ec.clients[c] = true
		ruc.clients[eventName] = ec
		rc.clients[userId] = ruc
		con.rooms[roomId] = rc
	} else if _, ok := con.rooms[roomId].clients[userId]; !ok {
		ruc := RoomUserClients{
			clients: make(map[string]EventClients),
		}
		ec := EventClients{
			clients: make(map[*Client]bool),
		}
		ec.clients[c] = true
		ruc.clients[eventName] = ec
		con.rooms[roomId].clients[userId] = ruc
	} else if _, ok := con.rooms[roomId].clients[userId].clients[eventName]; !ok {
		ec := EventClients{
			clients: make(map[*Client]bool),
		}
		ec.clients[c] = true
		con.rooms[roomId].clients[userId].clients[eventName] = ec
	} else {
		con.rooms[roomId].clients[userId].clients[eventName].clients[c] = true
	}
}

func (con *Connection) RemoveEvent(userId, roomId, eventName string, c *Client) {
	if userId == "" || roomId == "" || eventName == "" || c == nil {
		return
	}

	if _, ok := con.rooms[roomId].clients[userId].clients[eventName].clients[c]; ok {
		delete(con.rooms[roomId].clients[userId].clients[eventName].clients, c)
	}
	if len(con.rooms[roomId].clients[userId].clients[eventName].clients) == 0 {
		delete(con.rooms[roomId].clients[userId].clients, eventName)
	}
	if len(con.rooms[roomId].clients[userId].clients) == 0 {
		delete(con.rooms[roomId].clients, userId)
	}
	if len(con.rooms[roomId].clients) == 0 {
		delete(con.rooms, roomId)
	}
}
