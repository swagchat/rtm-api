// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/kylelemons/godebug/pretty"
	"github.com/swagchat/rtm-api/config"
	"github.com/swagchat/rtm-api/handlers"
	"github.com/swagchat/rtm-api/logger"
	"github.com/swagchat/rtm-api/messaging"
	"github.com/swagchat/rtm-api/metrics"
	"github.com/swagchat/rtm-api/rtm"
)

func main() {
	if config.StopRun {
		os.Exit(0)
	}

	cfg := config.Config()

	logger.InitLogger(cfg.Logger)

	compact := &pretty.Config{
		Compact: true,
	}
	logger.Info(fmt.Sprintf("Config: %s", compact.Sprint(cfg)))

	if cfg.Profiling {
		go func() {
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go metrics.Provider().Run()
	go messaging.Provider().Subscribe()
	go rtm.Server().Run()

	handlers.StartServer(ctx)
}
