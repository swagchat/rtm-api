// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"

	"github.com/fairway-corp/swagchat-realtime/handlers"
	"github.com/fairway-corp/swagchat-realtime/services"
	"github.com/fairway-corp/swagchat-realtime/utils"
)

func main() {
	var (
		port       string
		queHost    string
		quePort    string
		queTopic   string
		queChannel string
	)
	log.SetFlags(log.Llongfile)

	flag.StringVar(&port, "port", "9100", "service port")
	flag.StringVar(&queHost, "queHost", "", "Host name of nsqlookupd")
	flag.StringVar(&quePort, "quePort", "4161", "Port no of nsqlookupd")
	flag.StringVar(&queTopic, "queTopic", "websocket", "Topic name")
	flag.StringVar(&queChannel, "queChannel", "", "Channel name. If it's not set, channel is hostname.")
	flag.Parse()

	utils.Realtime.Port = port
	utils.Que.Host = queHost
	utils.Que.Port = quePort
	utils.Que.Topic = queTopic
	utils.Que.Channel = queChannel

	go services.Srv.Run()
	handlers.StartServer()
}
