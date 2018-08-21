package rtm

import (
	//"bytes"
	"context"
	"time"

	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/gorilla/websocket"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
	"github.com/swagchat/rtm-api/config"
	"github.com/swagchat/rtm-api/utils"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

type SpeechClient struct {
	Conn   *websocket.Conn
	Send   chan []byte
	Stream *speechpb.Speech_StreamingRecognizeClient
	Ctx    context.Context
}

// Client -> Server
func (c *SpeechClient) ReadPump() {
	defer func() {
		stream := *c.Stream
		stream.CloseSend()
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(config.Config().MaxMessageSize)
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	//log.Println("ReadPumpSpeech-------start")
	for {
		//log.Println("    IN ReadPumpSpeech-------start")
		messageType, data, err := c.Conn.ReadMessage()
		//log.Println("messageType", messageType)
		if err != nil {
			log.Println(err)
			break
		} else {
			if messageType == websocket.BinaryMessage {
				//buf := make([]byte, 1024)
				//buffer := bytes.NewBuffer(data)
				//for {
				//	n, err := buffer.Read(buf)
				//	if n == 0 || err == io.EOF {
				//		log.Println("io.EOF")
				//		break
				//	}
				//	go PostSpeech(c, buf[:n])
				//}
				go PostSpeech(c, data)
			} else {
				log.Println("Not Binary")
				break
			}
		}
		//log.Println("    IN ReadPumpSpeech-------out", c.Stream)
	}
	//log.Println("ReadPumpSpeech-------end")
}

func PostSpeech(c *SpeechClient, data []byte) {
	log.Println("=== PostSpeech ===")
	stream := *c.Stream
	//for {
	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
			AudioContent: data,
		},
	}); err != nil {
		log.Printf("Could not send audio: %v", err)
	}
	//}

	go func() {
		for {
			resp, err := stream.Recv()
			if resp == nil {
				log.Printf("[go func]resp is nil")
				stream.CloseSend()
				bgCtx := context.Background()
				ctx, _ := context.WithDeadline(bgCtx, time.Now().Add(1000*time.Second))
				newStream := utils.SpeechClient(ctx)
				c.Stream = &newStream
				break
			}
			if err == io.EOF {
				log.Printf("[go func]io.EOF")
				break
			}
			if err != nil {
				log.Printf("[go func]Cannot stream results: %v", err)
			}
			if err := resp.Error; err != nil {
				log.Printf("[go func]Could not recognize: %v", err)
			}
			for _, result := range resp.Results {
				log.Printf("%#v\n", result)
				if result.IsFinal {
					alts := result.GetAlternatives()
					transcript := alts[0].Transcript
					for i, alt := range alts {
						log.Println("[go func]", i, alt.Transcript)
					}

					payload := &utils.JSONText{}
					ps := fmt.Sprintf(`{\"text\":\"%s\"}`, transcript)
					payload.UnmarshalJSON([]byte(ps))
					m := &scpb.Message{
						Payload: *payload,
					}
					mBytes, err := json.Marshal(m)
					if err != nil {
						fmt.Println(err)
						return
					}

					c.Send <- mBytes
				}
			}
		}
	}()
	//log.Println("------- isFinish ------")
}

// Server -> Client
func (c *SpeechClient) WritePump() {
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
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *SpeechClient) write(mt int, payload []byte) error {
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Conn.WriteMessage(mt, payload)
}
