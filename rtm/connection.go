package rtm

import (
	"fmt"

	"github.com/swagchat/rtm-api/logging"
	"go.uber.org/zap/zapcore"
)

type Connection struct {
	clients map[*Client]bool
	users   map[string]*UserClients // index is userId
	rooms   map[string]*RoomClients // index is roomId
}

type UserClients struct {
	clients map[*Client]bool
}

type RoomClients struct {
	roomUsers map[string]*RoomUserClients // index is userId
}

type RoomUserClients struct {
	events map[string]*EventClients // index is eventName
}

type EventClients struct {
	clients map[*Client]bool
}

func (con *Connection) AddClient(c *Client) {
	if c == nil {
		return
	}

	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Kind:      "add-client",
		UserID:    c.UserId,
		Client:    fmt.Sprintf("%p", c.Conn),
		Useragent: c.Useragent,
		IPAddress: c.IPAddress,
	})

	con.clients[c] = true

	var userClients *UserClients
	if _, ok := con.users[c.UserId]; ok {
		con.users[c.UserId].clients[c] = true
	} else {
		userClients = &UserClients{
			clients: make(map[*Client]bool),
		}
		userClients.clients[c] = true
		con.users[c.UserId] = userClients
	}
}

func (con *Connection) RemoveClient(c *Client) {
	if c == nil {
		return
	}

	delete(con.clients, c)

	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Kind:   "delete-client",
		UserID: c.UserId,
		Client: fmt.Sprintf("%p", c.Conn),
	})

	for client, _ := range con.users[c.UserId].clients {
		if client == c {
			delete(con.users[c.UserId].clients, c)
		}
	}
	if len(con.users[c.UserId].clients) == 0 {
		delete(con.users, c.UserId)
	}

	for roomId, roomClients := range con.rooms {
		for userId, roomUserClients := range roomClients.roomUsers {
			for eventName, eventClients := range roomUserClients.events {
				for client, _ := range eventClients.clients {
					if client == c {
						delete(eventClients.clients, c)
					}
				}
				if len(roomUserClients.events[eventName].clients) == 0 {
					delete(roomUserClients.events, eventName)
				}
			}
			if len(roomClients.roomUsers[userId].events) == 0 {
				delete(roomClients.roomUsers, userId)
			}
		}
		if len(con.rooms[roomId].roomUsers) == 0 {
			delete(con.rooms, roomId)
		}
	}
}

func (con *Connection) AddEvent(userId, roomId, eventName string, c *Client) {
	if userId == "" || roomId == "" || eventName == "" || c == nil {
		return
	}

	if _, ok := con.rooms[roomId]; !ok {
		rc := &RoomClients{
			roomUsers: make(map[string]*RoomUserClients),
		}
		ruc := &RoomUserClients{
			events: make(map[string]*EventClients),
		}
		ec := &EventClients{
			clients: make(map[*Client]bool),
		}
		ec.clients[c] = true
		ruc.events[eventName] = ec
		rc.roomUsers[userId] = ruc
		con.rooms[roomId] = rc
	} else if _, ok := con.rooms[roomId].roomUsers[userId]; !ok {
		ruc := &RoomUserClients{
			events: make(map[string]*EventClients),
		}
		ec := &EventClients{
			clients: make(map[*Client]bool),
		}
		ec.clients[c] = true
		ruc.events[eventName] = ec
		con.rooms[roomId].roomUsers[userId] = ruc
	} else if _, ok := con.rooms[roomId].roomUsers[userId].events[eventName]; !ok {
		ec := &EventClients{
			clients: make(map[*Client]bool),
		}
		ec.clients[c] = true
		con.rooms[roomId].roomUsers[userId].events[eventName] = ec
	} else {
		con.rooms[roomId].roomUsers[userId].events[eventName].clients[c] = true
	}
}

func (con *Connection) RemoveEvent(userId, roomId, eventName string, c *Client) {
	if userId == "" || roomId == "" || eventName == "" || c == nil {
		return
	}

	if _, ok := con.rooms[roomId].roomUsers[userId].events[eventName].clients[c]; ok {
		delete(con.rooms[roomId].roomUsers[userId].events[eventName].clients, c)
	}
	if len(con.rooms[roomId].roomUsers[userId].events[eventName].clients) == 0 {
		delete(con.rooms[roomId].roomUsers[userId].events, eventName)
	}
	if len(con.rooms[roomId].roomUsers[userId].events) == 0 {
		delete(con.rooms[roomId].roomUsers, userId)
	}
	if len(con.rooms[roomId].roomUsers) == 0 {
		delete(con.rooms, roomId)
	}
}

func (conn *Connection) ConnectionCount() int {
	return len(conn.clients)
}

func (conn *Connection) Users() map[string]*UserClients {
	return conn.users
}

func (conn *Connection) Rooms() map[string]*RoomClients {
	return conn.rooms
}

func (conn *Connection) EachUserCount() map[string]int {
	euc := make(map[string]int, len(conn.users))
	for userID, userClients := range conn.users {
		euc[userID] = len(userClients.clients)
	}
	return euc
}

// func (con *Connection) Info() {
// hostname, _ := os.Hostname()
// log.Printf("[WS-INFO][%s] All Clients %d", hostname, len(con.clients))
// for userId, _ := range con.users {
// 	log.Printf("[WS-INFO][%s] UserId[%s] %d", hostname, userId, len(con.users[userId].clients))
// }
// for roomId, _ := range con.rooms {
// 	log.Printf("[WS-INFO][%s] RoomId[%s] %d", hostname, roomId, len(con.rooms[roomId].roomUsers))
// 	for userId, _ := range con.rooms[roomId].roomUsers {
// 		log.Printf("[WS-INFO][%s] RoomId[%s][%s] %d", hostname, roomId, userId, len(con.rooms[roomId].roomUsers[userId].events))
// 		for eventName, _ := range con.rooms[roomId].roomUsers[userId].events {
// 			log.Printf("[WS-INFO][%s] RoomId[%s][%s][%s] %d", hostname, roomId, userId, eventName, len(con.rooms[roomId].roomUsers[userId].events[eventName].clients))
// 		}
// 	}
// }
// }
