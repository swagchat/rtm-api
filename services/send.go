package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/fairway-corp/swagchat-realtime/models"
)

func Send(message []byte) {
	var receivedMessage models.ReceivedMessage
	if err := json.Unmarshal(message, &receivedMessage); err == nil {
		log.Println("==========================>")
		log.Printf("%#v\n", receivedMessage)

		if receivedMessage.Message != nil {
			log.Printf("%#v\n", receivedMessage.Message)
			log.Printf("%#v\n", receivedMessage.Message.Data)
			message, err = base64.StdEncoding.DecodeString(receivedMessage.Message.Data)
			if err != nil {
				log.Println(err.Error())
			}
		}
		log.Println("<==========================")
	}

	var sendMessage []byte
	sendMessage = bytes.Replace(message, []byte("\n"), []byte(" "), -1)
	sendMessage = bytes.Replace(sendMessage, []byte("\\"), []byte(""), -1)
	sendMessage = bytes.TrimSpace(sendMessage)
	Srv.Broadcast <- sendMessage
}
