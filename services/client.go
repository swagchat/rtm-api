package services

import (
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/swagchat/rtm-api/logging"
	"github.com/swagchat/rtm-api/utils"
	"go.uber.org/zap/zapcore"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	Conn    *websocket.Conn
	Send    chan []byte
	RoomIds []string
	UserId  string
}

type RcvData struct {
	UserId    string
	RoomId    string `json:roomId`
	EventName string `json:eventName`
	Action    string `json:action,omitempty`
	Client    *Client
}

// Client -> Server
func (c *Client) ReadPump() {
	defer func() {
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(utils.Config().MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		rcvData := &RcvData{}
		err := c.Conn.ReadJSON(&rcvData)
		if err != nil {
			pos := strings.LastIndex(err.Error(), "normal")
			if pos == 0 {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Kind:    "websocket-error",
					Message: err.Error(),
				})
			}

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				c.Conn.Close()
				Srv.Close <- c
			}
			break
		}
		rcvData.Client = c
		rcvData.UserId = c.UserId

		c.RoomIds = append(c.RoomIds, rcvData.RoomId)
		rcvData.UserId = c.UserId
		switch rcvData.Action {
		case "bind":
			Srv.Register <- rcvData
		case "unbind":
			Srv.Unregister <- rcvData
		}
	}
}

// Server -> Client
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
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Kind:    "websocket-error",
					Message: err.Error(),
				})
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Kind:    "websocket-error",
					Message: err.Error(),
				})
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Kind:    "websocket-error",
					Message: err.Error(),
				})
				return
			}
		}
	}
}

func (c *Client) write(mt int, payload []byte) error {
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Conn.WriteMessage(mt, payload)
}
