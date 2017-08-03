package services

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Conn struct {
	Ws      *websocket.Conn
	Send    chan []byte
	RoomIds []string
}

type socketData struct {
	RoomId    string `json:roomId`
	EventName string `json:eventName`
	Action    string `json:eventName,omitempty`
}

type RegisterData struct {
	roomId    string
	eventName string
	conn      *Conn
}

func (c *Conn) ReadPump() {
	defer func() {
		//Manager.Unregister <- registerData
		c.Ws.Close()
	}()
	c.Ws.SetReadLimit(maxMessageSize)
	c.Ws.SetReadDeadline(time.Now().Add(pongWait))
	c.Ws.SetPongHandler(func(string) error { c.Ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		log.Println("---- ReadPump ----")
		socketData := socketData{}
		err := c.Ws.ReadJSON(&socketData)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
				c.Ws.Close()
				Srv.Close <- c
			}
			break
		}

		log.Printf("Get message: %#v\n", socketData)
		c.RoomIds = append(c.RoomIds, socketData.RoomId)
		registerData := &RegisterData{
			roomId:    socketData.RoomId,
			eventName: socketData.EventName,
			conn:      c,
		}
		switch socketData.Action {
		case "bind":
			Srv.Register <- registerData
		case "unbind":
			Srv.Unregister <- registerData
		}
	}
}

func (c *Conn) write(mt int, payload []byte) error {
	c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Ws.WriteMessage(mt, payload)
}

func (c *Conn) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()
	for {
		log.Println("---- WritePump ----")
		select {
		case message, ok := <-c.Send:
			log.Println("<-c.Send", ok)
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.Ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			log.Println("<-ticker.C")
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
