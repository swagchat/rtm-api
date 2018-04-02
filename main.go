// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/swagchat/rtm-api/handlers"
	"github.com/swagchat/rtm-api/messaging"
	"github.com/swagchat/rtm-api/metrics"
	"github.com/swagchat/rtm-api/services"
	"github.com/swagchat/rtm-api/utils"
)

func main() {
	if utils.IsShowVersion {
		fmt.Printf("API Version %s\nBuild Version %s\n", utils.APIVersion, utils.BuildVersion)
		return
	}

	if utils.Config().Profiling {
		go func() {
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
	}

	go metrics.MetricsProvider().Run()
	go messaging.MessagingProvider().Subscribe()
	go services.Srv.Run()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handlers.StartServer(ctx)
}
