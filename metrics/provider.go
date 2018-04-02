package metrics

import (
	"os"
	"time"

	stats "github.com/fukata/golang-stats-api-handler"
	"github.com/swagchat/rtm-api/rtm"
	"github.com/swagchat/rtm-api/utils"
)

type Metrics struct {
	Hostname          string                    `json:"hostname"`
	Stats             *stats.Stats              `json:"stats"`
	AllCount          int                       `json:"allCount"`
	UserCount         int                       `json:"userCount"`
	RoomCount         int                       `json:"roomCount"`
	EachUserCount     map[string]int            `json:"eachUserCount,omitempty"`
	EachRoomCount     map[string]int            `json:"eachRoomCount,omitempty"`
	EachRoomUserCount map[string]map[string]int `json:"eachRoomUserCount,omitempty"`
	Timestamp         string                    `json:"timestamp"`
}

type provider interface {
	Run()
}

func Provider() provider {
	c := utils.Config()

	var p provider
	switch c.Metrics.Provider {
	case "":
		p = &notuseProvider{}
	case "stdout":
		p = &stdoutProvider{}
	case "elasticsearch":
		p = &elasticsearchProvider{}
	}
	return p
}

func makeMetrics(nowTime time.Time) *Metrics {
	c := utils.Config()
	m := &Metrics{}

	hostname, _ := os.Hostname()
	m.Hostname = hostname
	m.Stats = stats.GetStats()

	srv := rtm.Server()
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
