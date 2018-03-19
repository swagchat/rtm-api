// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"

	"github.com/swagchat/rtm-api/handlers"
	"github.com/swagchat/rtm-api/messaging"
	"github.com/swagchat/rtm-api/services"
	"github.com/swagchat/rtm-api/utils"
)

func main() {
	utils.SetupLogger()

	if utils.IsShowVersion {
		fmt.Printf("API Version %s\nBuild Version %s\n", utils.APIVersion, utils.BuildVersion)
		return
	}

	if utils.GetConfig().Profiling {
		go func() {
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
	}

	messaging.GetMessagingProvider().Subscribe()
	go services.Srv.Run()
	handlers.StartServer()
}
