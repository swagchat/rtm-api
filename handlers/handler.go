package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-zoo/bone"
	"github.com/gorilla/websocket"
	"github.com/shogo82148/go-gracedown"
	"github.com/swagchat/rtm-api/logging"
	"github.com/swagchat/rtm-api/messaging"
	"github.com/swagchat/rtm-api/services"
	"github.com/swagchat/rtm-api/utils"
	"go.uber.org/zap/zapcore"
)

var (
	allowedMethods = []string{
		"POST",
		"GET",
		"OPTIONS",
		"PUT",
		"PATCH",
		"DELETE",
	}

	noBodyStatusCodes = []int{
		http.StatusNotFound,
		http.StatusConflict,
	}
)

// StartServer is HTTP Server
func StartServer() {
	mux := bone.New()
	mux.GetFunc("/", indexHandler)
	mux.GetFunc("/stats", stats_api.Handler)
	mux.GetFunc("/ws", websocketHandler)
	mux.GetFunc("/speech", speechHandler)
	mux.PostFunc("/message", messageHandler)
	mux.OptionsFunc("/*", optionsHandler)
	mux.NotFoundFunc(notFoundHandler)

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			s := <-signalChan
			if s == syscall.SIGTERM || s == syscall.SIGINT {
				logging.Log(zapcore.InfoLevel, &logging.AppLog{
					Kind:    "handler",
					Message: fmt.Sprintf("%s graceful down by signal[%s]", utils.AppName, s.String()),
				})
				messaging.MessagingProvider().Unsubscribe()
				gracedown.Close()
			}
		}
	}()

	c := utils.Config()
	sb := utils.NewStringBuilder()
	cfgStr := sb.PrintStruct("config", c)
	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Kind:    "handler",
		Message: fmt.Sprintf("%s start", utils.AppName),
		Config:  cfgStr,
	})

	gracedown.ListenAndServe(fmt.Sprintf(":%s", c.HTTPPort), mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusOK, "text/plain", fmt.Sprintf("%s [API Version]%s [Build Version]%s", utils.AppName, utils.APIVersion, utils.BuildVersion))
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Kind:    "handler-error",
			Message: err.Error(),
		})
	}
	services.Srv.Broadcast <- message
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Kind:    "handler-error",
			Message: err.Error(),
		})
		return
	}

	c := &services.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	params, _ := url.ParseQuery(r.URL.RawQuery)
	if userIDSlice, ok := params["userId"]; ok {
		c.UserId = userIDSlice[0]
	}

	services.Srv.Connection.AddClient(c)
	go c.WritePump()
	go c.ReadPump()
}

func speechHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := speechUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Kind:    "handler-error",
			Message: err.Error(),
		})
		return
	}

	ctx := r.Context()
	stream := utils.SpeechClient(ctx)

	c := &services.SpeechClient{
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Stream: &stream,
		Ctx:    ctx,
	}

	go c.WritePump()
	go c.ReadPump()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  utils.Config().ReadBufferSize,
	WriteBufferSize: utils.Config().WriteBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var speechUpgrader = websocket.Upgrader{
	ReadBufferSize:  utils.Config().ReadBufferSize,
	WriteBufferSize: utils.Config().WriteBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func colsHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		optionsHandler(w, r)
		fn(w, r)
	}
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	rHeaders := make([]string, 0, len(r.Header))
	for k, v := range r.Header {
		rHeaders = append(rHeaders, k)
		if k == "Access-Control-Request-Headers" {
			rHeaders = append(rHeaders, strings.Join(v, ", "))
		}
	}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(rHeaders, ", "))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func encodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func respond(w http.ResponseWriter, r *http.Request, status int, contentType string, data interface{}) {
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(status)
	for _, v := range noBodyStatusCodes {
		if status == v {
			data = nil
		}
	}
	if data != nil {
		encodeBody(w, r, data)
	}
}
