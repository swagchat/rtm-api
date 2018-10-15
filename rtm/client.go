package rtm

import (
	"strings"
	"time"

	logger "github.com/betchi/zapper"
	"github.com/gorilla/websocket"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
	"github.com/swagchat/rtm-api/config"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 4) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	Conn      *websocket.Conn
	Send      chan []byte
	RoomIds   []string
	UserId    string
	Useragent string
	IPAddress string
	Language  string
}

type RcvData struct {
	UserId    string
	RoomId    string
	EventType scpb.EventType
	Action    string
	Client    *Client
}

// ReadPump (Client -> Server)
func (c *Client) ReadPump() {
	defer func() {
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(config.Config().MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		rcvData := &RcvData{}
		err := c.Conn.ReadJSON(&rcvData)
		if err != nil {
			pos := strings.LastIndex(err.Error(), "normal")
			if pos == 0 {
				logger.Error(err.Error())
			}

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				c.Conn.Close()
				srv.Close <- c
			}
			break
		}
		rcvData.Client = c
		rcvData.UserId = c.UserId

		c.RoomIds = append(c.RoomIds, rcvData.RoomId)
		rcvData.UserId = c.UserId
		switch rcvData.Action {
		case "subscribe":
			srv.Register <- rcvData
		case "unsubscribe":
			srv.Unregister <- rcvData
		}
	}
}

// WritePump (Server -> Client)
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.Conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				logger.Error(err.Error())
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			err = w.Close()
			if err != nil {
				logger.Error(err.Error())
				return
			}
		case <-ticker.C:
			err := c.write(websocket.PingMessage, []byte{})
			if err != nil {
				logger.Error(err.Error())
				return
			}
		}
	}
}

func (c *Client) write(mt int, payload []byte) error {
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Conn.WriteMessage(mt, payload)
}
