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
	"time"

	"github.com/betchi/metrictor"
	"github.com/betchi/tracer"
	logger "github.com/betchi/zapper"
	elasticapmLogger "github.com/betchi/zapper/elasticapm"
	jaegerLogger "github.com/betchi/zapper/jaeger"
	"github.com/kylelemons/godebug/pretty"
	"github.com/swagchat/rtm-api/config"
	"github.com/swagchat/rtm-api/consumer"
	"github.com/swagchat/rtm-api/rest"
	"github.com/swagchat/rtm-api/rtm"
)

func main() {
	if config.StopRun {
		os.Exit(0)
	}

	cfg := config.Config()

	logger.InitGlobalLogger(&logger.Config{
		EnableConsole:  cfg.Logger.EnableConsole,
		ConsoleFormat:  cfg.Logger.ConsoleFormat,
		ConsoleLevel:   cfg.Logger.ConsoleLevel,
		EnableFile:     cfg.Logger.EnableFile,
		FileFormat:     cfg.Logger.FileFormat,
		FileLevel:      cfg.Logger.FileLevel,
		FilePath:       cfg.Logger.FilePath,
		FileMaxSize:    cfg.Logger.FileMaxSize,
		FileMaxAge:     cfg.Logger.FileMaxAge,
		FileMaxBackups: cfg.Logger.FileMaxBackups,
		FileLocalTime:  cfg.Logger.FileLocalTime,
		FileCompress:   cfg.Logger.FileCompress,
	})

	compact := &pretty.Config{
		Compact: true,
	}
	logger.Info(fmt.Sprintf("Config: %s", compact.Sprint(cfg)))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg.Profiling {
		go func() {
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
		metrictor.SetInt(metrictor.EachTime, "allConns", func() int64 {
			srv := rtm.Server()
			return int64(srv.Connection.ConnectionCount())
		})
		metrictor.SetInt(metrictor.EachTime, "userCount", func() int64 {
			srv := rtm.Server()
			return int64(len(srv.Connection.Users()))
		})
		metrictor.SetInt(metrictor.EachTime, "eventCount", func() int64 {
			srv := rtm.Server()
			return int64(len(srv.Connection.Events()))
		})
		metrictor.Run(ctx, time.Second*5)
	}

	go consumer.Provider(ctx).SubscribeMessage()

	jaegerLogger.InitGlobalLogger(&jaegerLogger.Config{Noop: !cfg.Tracer.Logging})
	elasticapmLogger.InitGlobalLogger(&elasticapmLogger.Config{Noop: !cfg.Tracer.Logging})
	err := tracer.InitGlobalTracer(&tracer.Config{
		Provider:       cfg.Tracer.Provider,
		ServiceName:    config.AppName,
		ServiceVersion: config.BuildVersion,
		Jaeger: &tracer.Jaeger{
			Logger: jaegerLogger.GlobalLogger(),
		},
		Zipkin: &tracer.Zipkin{
			Logger:    jaegerLogger.GlobalLogger(),
			Endpoint:  cfg.Tracer.Zipkin.Endpoint,
			BatchSize: cfg.Tracer.Zipkin.BatchSize,
			Timeout:   cfg.Tracer.Zipkin.Timeout,
		},
		ElasticAPM: &tracer.ElasticAPM{
			Logger: elasticapmLogger.GlobalLogger(),
		},
	})
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer tracer.Close()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTSTP, syscall.SIGKILL, syscall.SIGSTOP)
	go func() {
		<-sigChan
		cancel()
	}()

	go rtm.Server().Run()

	rest.Run(ctx)
}
