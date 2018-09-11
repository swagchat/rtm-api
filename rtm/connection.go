package rtm

import (
	"fmt"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
	"github.com/swagchat/rtm-api/logger"
)

type Connection struct {
	clients map[*Client]bool
	users   map[string]*UserClients          // index is userId
	events  map[scpb.EventType]*EventClients // index is eventName
}

type UserClients struct {
	clients map[*Client]bool
}

type EventClients struct {
	users map[string]*UserClients
}

func (con *Connection) AddClient(c *Client) {
	if c == nil {
		return
	}

	logger.Info(fmt.Sprintf("add-client userId[%s] ip[%s]", c.UserId, c.IPAddress))
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

	logger.Info(fmt.Sprintf("add-client userId[%s] ip[%s]", c.UserId, c.IPAddress))
	delete(con.clients, c)

	for client := range con.users[c.UserId].clients {
		if client == c {
			delete(con.users[c.UserId].clients, c)
		}
	}
	if len(con.users[c.UserId].clients) == 0 {
		delete(con.users, c.UserId)
	}

	for eventName, eventClients := range con.events {
		for userID := range eventClients.users {
			for client := range eventClients.users[userID].clients {
				if client == c {
					delete(con.events[eventName].users[userID].clients, c)
				}
			}
			if len(con.events[eventName].users[userID].clients) == 0 {
				delete(con.events[eventName].users, userID)
			}
		}
		if len(con.events[eventName].users) == 0 {
			delete(con.events, eventName)
		}
	}
}

func (con *Connection) AddEvent(userId string, eventType scpb.EventType, c *Client) {
	if userId == "" || c == nil {
		return
	}

	if eventType != scpb.EventType_MessageEvent && eventType != scpb.EventType_RoomEvent {
		return
	}

	if _, ok := con.events[eventType]; !ok {
		uc := &UserClients{
			clients: make(map[*Client]bool),
		}
		uc.clients[c] = true

		ec := &EventClients{
			users: make(map[string]*UserClients),
		}
		ec.users[userId] = uc
		con.events[eventType] = ec
	} else if _, ok := con.events[eventType].users[userId]; !ok {
		uc := &UserClients{
			clients: make(map[*Client]bool),
		}
		uc.clients[c] = true
		con.events[eventType].users[userId] = uc
	} else if _, ok := con.events[eventType].users[userId].clients[c]; !ok {
		con.events[eventType].users[userId].clients[c] = true
	}
}

func (con *Connection) RemoveEvent(userId string, eventType scpb.EventType, c *Client) {
	if userId == "" || c == nil {
		return
	}

	if eventType != scpb.EventType_MessageEvent && eventType != scpb.EventType_RoomEvent {
		return
	}

	if _, ok := con.events[eventType]; !ok {
		return
	}

	if _, ok := con.events[eventType].users[userId]; !ok {
		return
	}

	if _, ok := con.events[eventType].users[userId].clients[c]; ok {
		delete(con.events[eventType].users[userId].clients, c)
	}
	if len(con.events[eventType].users[userId].clients) == 0 {
		delete(con.events[eventType].users, userId)
	}
	if len(con.events) == 0 {
		delete(con.events, eventType)
	}
}

func (conn *Connection) ConnectionCount() int {
	return len(conn.clients)
}

func (conn *Connection) Users() map[string]*UserClients {
	return conn.users
}

func (conn *Connection) Events() map[scpb.EventType]*EventClients {
	return conn.events
}

func (conn *Connection) EachUserCount() map[string]int {
	euc := make(map[string]int, len(conn.users))
	for userID, userClients := range conn.users {
		euc[userID] = len(userClients.clients)
	}
	return euc
}

func (conn *Connection) EachEventCount() map[scpb.EventType]map[string]int {
	eec := make(map[scpb.EventType]map[string]int, len(conn.events))
	for event, eventClients := range conn.events {
		euc := make(map[string]int, len(eventClients.users))
		for userID, userClients := range eventClients.users {
			euc[userID] = len(userClients.clients)
		}
		eec[event] = euc
	}
	return eec
}
