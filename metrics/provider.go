package metrics

import (
	"os"
	"time"

	"github.com/swagchat/rtm-api/services"
	"github.com/swagchat/rtm-api/utils"
	"go.uber.org/zap"
)

type Metrics struct {
	Hostname          string                    `json:"hostname"`
	AllCount          int                       `json:"allCount"`
	UserCount         int                       `json:"userCount"`
	RoomCount         int                       `json:"roomCount"`
	EachUserCount     map[string]int            `json:"eachUserCount,omitempty"`
	EachRoomCount     map[string]int            `json:"eachRoomCount,omitempty"`
	EachRoomUserCount map[string]map[string]int `json:"eachRoomUserCount,omitempty"`
	Timestamp         string                    `json:"timestamp"`
}

type Provider interface {
	Run()
}

func GetMetricsProvider() Provider {
	c := utils.GetConfig()

	var provider Provider
	switch c.Metrics.Provider {
	case "":
		provider = &NotuseProvider{}
	case "stdout":
		provider = &StdoutProvider{}
	case "elasticsearch":
		provider = &ElasticsearchProvider{}
	default:
		utils.AppLogger.Error("",
			zap.String("msg", "utils.Cfg.MetricsProvider is incorrect"),
		)
		os.Exit(0)
	}
	return provider
}

func makeMetrics(nowTime time.Time) *Metrics {
	c := utils.GetConfig()
	m := &Metrics{}

	hostname, _ := os.Hostname()
	m.Hostname = hostname

	srv := services.GetServer()
	users := srv.Connection.Users()
	rooms := srv.Connection.Rooms()
	m.AllCount = srv.Connection.ConnectionCount()
	m.UserCount = len(users)
	m.RoomCount = len(rooms)

	m.Timestamp = nowTime.Format(time.RFC3339)

	if c.Metrics.Verbose {
		// TODO
	}

	return m
}

func exec(fn func(), interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fn()
		}
	}
}
