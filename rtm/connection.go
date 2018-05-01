package rtm

import (
	"fmt"

	"github.com/swagchat/rtm-api/logging"
	"go.uber.org/zap/zapcore"
)

type Connection struct {
	clients map[*Client]bool
	users   map[string]*UserClients  // index is userId
	events  map[string]*EventClients // index is eventName
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

	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Kind:   "delete-client",
		UserID: c.UserId,
		Client: fmt.Sprintf("%p", c.Conn),
	})

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

func (con *Connection) AddEvent(userId, eventName string, c *Client) {
	if userId == "" || eventName == "" || c == nil {
		return
	}

	if _, ok := con.events[eventName]; !ok {
		uc := &UserClients{
			clients: make(map[*Client]bool),
		}
		uc.clients[c] = true

		ec := &EventClients{
			users: make(map[string]*UserClients),
		}
		ec.users[userId] = uc
		con.events[eventName] = ec
	} else if _, ok := con.events[eventName].users[userId]; !ok {
		uc := &UserClients{
			clients: make(map[*Client]bool),
		}
		uc.clients[c] = true
		con.events[eventName].users[userId] = uc
	} else if _, ok := con.events[eventName].users[userId].clients[c]; !ok {
		con.events[eventName].users[userId].clients[c] = true
	}
}

func (con *Connection) RemoveEvent(userId, eventName string, c *Client) {
	if userId == "" || eventName == "" || c == nil {
		return
	}

	if _, ok := con.events[eventName]; !ok {
		return
	}

	if _, ok := con.events[eventName].users[userId]; !ok {
		return
	}

	if _, ok := con.events[eventName].users[userId].clients[c]; ok {
		delete(con.events[eventName].users[userId].clients, c)
	}
	if len(con.events[eventName].users[userId].clients) == 0 {
		delete(con.events[eventName].users, userId)
	}
	if len(con.events) == 0 {
		delete(con.events, eventName)
	}
}

func (conn *Connection) ConnectionCount() int {
	return len(conn.clients)
}

func (conn *Connection) Users() map[string]*UserClients {
	return conn.users
}

func (conn *Connection) Events() map[string]*EventClients {
	return conn.events
}

func (conn *Connection) EachUserCount() map[string]int {
	euc := make(map[string]int, len(conn.users))
	for userID, userClients := range conn.users {
		euc[userID] = len(userClients.clients)
	}
	return euc
}

func (conn *Connection) EachEventCount() map[string]map[string]int {
	eec := make(map[string]map[string]int, len(conn.events))
	for eventName, eventClients := range conn.events {
		euc := make(map[string]int, len(eventClients.users))
		for userID, userClients := range eventClients.users {
			euc[userID] = len(userClients.clients)
		}
		eec[eventName] = euc
	}
	return eec
}
