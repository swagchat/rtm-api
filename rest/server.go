package rest

import (
	"context"
	"encoding/json"
	"expvar"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/betchi/tracer"
	logger "github.com/betchi/zapper"
	"github.com/fukata/golang-stats-api-handler"
	"github.com/go-zoo/bone"
	"github.com/gorilla/websocket"
	"github.com/shogo82148/go-gracedown"
	"github.com/swagchat/rtm-api/config"
	"github.com/swagchat/rtm-api/rtm"
	"github.com/swagchat/rtm-api/utils"
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

// Run runs start REST API server
func Run(ctx context.Context) {
	cfg := config.Config()

	mux := bone.New()
	mux.GetFunc("/", indexHandler)
	mux.GetFunc("/stats", stats_api.Handler)
	mux.GetFunc("/ws", websocketHandler)
	mux.GetFunc("/speech", commonHandler(speechHandler))
	mux.PostFunc("/message", commonHandler(messageHandler))
	mux.OptionsFunc("/*", optionsHandler)
	mux.NotFoundFunc(notFoundHandler)

	if cfg.Profiling {
		mux.GetFunc("/debug/vars", metricsHandler)
	}

	logger.Info(fmt.Sprintf("Starting %s server[REST] on listen tcp :%s", config.AppName, cfg.HTTPPort))
	errCh := make(chan error)
	go func() {
		errCh <- gracedown.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), mux)
	}()

	select {
	case <-ctx.Done():
		logger.Info(fmt.Sprintf("Stopping %s server[REST]", config.AppName))
		gracedown.Close()
	case err := <-errCh:
		logger.Error(fmt.Sprintf("Failed to serve %s server[REST]. %v", config.AppName, err))
	}
}

type customResponseWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *customResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, r, http.StatusOK, "text/plain", fmt.Sprintf("%s [API Version]%s [Build Version]%s", config.AppName, config.APIVersion, config.BuildVersion))
}

func commonHandler(fn http.HandlerFunc) http.HandlerFunc {
	return (colsHandler(
		traceHandler(
			func(w http.ResponseWriter, r *http.Request) {
				for i, v := range r.Header {
					log.Printf("%s=%s\n", i, v)
				}
				defer r.Body.Close()
				fn(w, r)
			})))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	first := true
	report := func(key string, value interface{}) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		if str, ok := value.(string); ok {
			fmt.Fprintf(w, "%q: %q", key, str)
		} else {
			fmt.Fprintf(w, "%q: %v", key, value)
		}
	}

	fmt.Fprintf(w, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		report(kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}

func traceHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.StartTransaction(
			r.Context(),
			fmt.Sprintf("%s:%s", r.Method, r.RequestURI), "REST",
			tracer.StartTransactionOptionWithHTTPRequest(r),
		)
		defer tracer.CloseTransaction(ctx)

		sw := &customResponseWriter{ResponseWriter: w}
		fn(sw, r.WithContext(ctx))

		tracer.SetHTTPStatusCode(span, sw.status)
		tracer.SetTag(span, "http.method", r.Method)
		tracer.SetTag(span, "http.content_length", sw.length)
		tracer.SetTag(span, "http.referer", r.Referer())
	}
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "messageHandler", "rest")
	defer tracer.Finish(span)

	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err.Error())
	}
	rtm.Server().Broadcast <- message
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "websocketHandler", "rest")
	defer tracer.Finish(span)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		tracer.SetError(span, err)
		logger.Error(err.Error())
		return
	}

	c := &rtm.Client{
		Conn:      conn,
		Send:      make(chan []byte, 256),
		Useragent: r.Header.Get("User-Agent"),
		IPAddress: r.Header.Get("Cf-Connecting-Ip"),
		Language:  r.Header.Get("Accept-Language"),
	}

	params, _ := url.ParseQuery(r.URL.RawQuery)
	if userIDSlice, ok := params["userId"]; ok {
		c.UserId = userIDSlice[0]
	}

	rtm.Server().Connection.AddClient(c)
	go c.WritePump()
	go c.ReadPump()
}

func speechHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := speechUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	ctx := r.Context()
	stream := utils.SpeechClient(ctx)

	c := &rtm.SpeechClient{
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Stream: &stream,
		Ctx:    ctx,
	}

	go c.WritePump()
	go c.ReadPump()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  config.Config().ReadBufferSize,
	WriteBufferSize: config.Config().WriteBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var speechUpgrader = websocket.Upgrader{
	ReadBufferSize:  config.Config().ReadBufferSize,
	WriteBufferSize: config.Config().WriteBufferSize,
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
