// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"

	"github.com/fairway-corp/swagchat-realtime/handlers"
	"github.com/fairway-corp/swagchat-realtime/services"
)

var port = flag.String("port", "9100", "service port")

func main() {
	flag.Parse()
	log.SetFlags(log.Llongfile)

	go services.Manager.Run()
	handlers.StartServer(*port)
}
