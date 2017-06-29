package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fairway-corp/swagchat-realtime/models"
	"github.com/fairway-corp/swagchat-realtime/services"
	"github.com/fairway-corp/swagchat-realtime/utils"
	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-zoo/bone"
	"github.com/gorilla/websocket"
	"github.com/shogo82148/go-gracedown"
)

func StartServer(port string) {
	mux := bone.New().Prefix(utils.AppendStrings("/", utils.API_VERSION))
	mux.GetFunc("", websocketHandler)
	mux.GetFunc("/", websocketHandler)
	mux.PostFunc("/message", messageHandler)
	mux.GetFunc("/stats", stats_api.Handler)
	mux.OptionsFunc("/*", optionsHandler)
	mux.NotFoundFunc(notFoundHandler)
	log.Printf("port %s", port)

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			s := <-signalChan
			if s == syscall.SIGTERM || s == syscall.SIGINT {
				log.Println("msg", "Swagchat Realtime Shutdown start!")
				gracedown.Close()
			}
		}
	}()
	log.Println("Swagchat realtime server start!")
	if err := gracedown.ListenAndServe(utils.AppendStrings(":", port), mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- messageHandler ---")
	var sendMessage []byte

	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println("\n", string(message))

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
	sendMessage = bytes.Replace(message, []byte("\n"), []byte(" "), -1)
	sendMessage = bytes.Replace(sendMessage, []byte("\\"), []byte(""), -1)
	sendMessage = bytes.TrimSpace(sendMessage)
	services.Manager.Broadcast <- sendMessage
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	addrs, err := net.InterfaceAddrs()
	log.Println("-------------------------------------------------")
	for _, addrs := range addrs {
		log.Println(addrs.String())
	}
	log.Println("-------------------------------------------------")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn := &services.Conn{
		Send: make(chan []byte, 256),
		Ws:   ws,
	}
	//services.Manager.Register <- conn
	go conn.WritePump()
	conn.ReadPump()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("options")
	log.Println(r.Header.Get("Access-Control-Request-Headers"))

	origin := r.Header.Get("Origin")
	if origin != "" {
		log.Println(origin)
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

var allowedMethods []string = []string{
	"POST",
	"GET",
	"OPTIONS",
	"PUT",
	"PATCH",
	"DELETE",
}
