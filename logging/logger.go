package logging

import (
	"os"

	"github.com/swagchat/rtm-api/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var appLogger = NewAppLogger()

type AppLog struct {
	Kind     string `json:"kind"`
	UserID   string `json:"userId,omitempty"`
	RoomID   string `json:"roomId,omitempty"`
	Event    string `json:"event,omitempty"`
	Client   string `json:"client,omitempty"`
	Provider string `json:"provider,omitempty"`
	Config   string `json:"config,omitempty"`
	Message  string `json:"message,omitempty"`
}

func NewAppLogger() *zap.Logger {
	c := utils.Config()

	var err error
	var logger *zap.Logger
	if c.Logging.Level == "production" {
		logger, err = zap.NewProduction()
	} else if c.Logging.Level == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		os.Exit(0)
	}
	if err != nil {
		os.Exit(0)
	}
	hostname, _ := os.Hostname()
	appLogger := logger.WithOptions(zap.Fields(
		zap.String("appName", utils.AppName),
		zap.String("apiVersion", utils.APIVersion),
		zap.String("buildVersion", utils.BuildVersion),
		zap.String("hostname", hostname),
	))

	return appLogger
}

// func AppLogger() *zap.Logger {
// 	return appLogger
// }

func Log(level zapcore.Level, al *AppLog) {
	fields := make([]zapcore.Field, 0)
	if al.Kind != "" {
		fields = append(fields, zap.String("kind", al.Kind))
	}
	if al.UserID != "" {
		fields = append(fields, zap.String("userId", al.UserID))
	}
	if al.RoomID != "" {
		fields = append(fields, zap.String("roomId", al.RoomID))
	}
	if al.Event != "" {
		fields = append(fields, zap.String("event", al.Event))
	}
	if al.Client != "" {
		fields = append(fields, zap.String("client", al.Client))
	}
	if al.Provider != "" {
		fields = append(fields, zap.String("provider", al.Provider))
	}
	if al.Config != "" {
		fields = append(fields, zap.String("config", al.Config))
	}
	if al.Message != "" {
		fields = append(fields, zap.String("message", al.Message))
	}

	switch level {
	case zapcore.DebugLevel:
		appLogger.Debug("", fields...)
	case zapcore.InfoLevel:
		appLogger.Info("", fields...)
	case zapcore.WarnLevel:
		appLogger.Warn("", fields...)
	case zapcore.ErrorLevel:
		appLogger.Error("", fields...)
	case zapcore.DPanicLevel:
		appLogger.DPanic("", fields...)
	case zapcore.PanicLevel:
		appLogger.Panic("", fields...)
	case zapcore.FatalLevel:
		appLogger.Fatal("", fields...)
	}
}
