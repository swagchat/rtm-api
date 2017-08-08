package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fairway-corp/swagchat-realtime/messaging"
	"github.com/fairway-corp/swagchat-realtime/services"
	"github.com/fairway-corp/swagchat-realtime/utils"
	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-zoo/bone"
	"github.com/gorilla/websocket"
	"github.com/shogo82148/go-gracedown"
)

func StartServer() {
	mux := bone.New().Prefix(utils.AppendStrings("/", utils.API_VERSION))
	mux.GetFunc("", websocketHandler)
	mux.GetFunc("/", websocketHandler)
	mux.PostFunc("/message", messageHandler)
	mux.GetFunc("/stats", stats_api.Handler)
	mux.OptionsFunc("/*", optionsHandler)
	mux.NotFoundFunc(notFoundHandler)
	log.Printf("port %s", utils.Realtime.Port)

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			s := <-signalChan
			if s == syscall.SIGTERM || s == syscall.SIGINT {
				log.Println("msg", "Swagchat Realtime Shutdown start!")
				messaging.Unsubscribe()
				gracedown.Close()
			}
		}
	}()
	log.Println("Swagchat realtime server start!")
	if err := gracedown.ListenAndServe(utils.AppendStrings(":", utils.Realtime.Port), mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- messageHandler ---")
	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	services.Srv.Broadcast <- message
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &services.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	params, _ := url.ParseQuery(r.URL.RawQuery)
	if userIdSlice, ok := params["userId"]; ok {
		c.UserId = userIdSlice[0]
	}

	services.Srv.Connection.AddClient(c)
	go c.WritePump()
	go c.ReadPump()
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
