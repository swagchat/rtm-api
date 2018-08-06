// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kylelemons/godebug/pretty"
	"github.com/swagchat/rtm-api/config"
	"github.com/swagchat/rtm-api/logger"
	"github.com/swagchat/rtm-api/metrics"
	"github.com/swagchat/rtm-api/rest"
	"github.com/swagchat/rtm-api/rtm"
	"github.com/swagchat/rtm-api/sbroker"
	"github.com/swagchat/rtm-api/tracer"
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

	go sbroker.Provider(ctx).SubscribeMessage()

	err := tracer.Provider(ctx).NewTracer()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer tracer.Provider(ctx).Close()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTSTP, syscall.SIGKILL, syscall.SIGSTOP)
	go func() {
		<-sigChan
		cancel()
	}()

	go metrics.Provider().Run()
	go rtm.Server().Run()

	rest.Run(ctx)
}
