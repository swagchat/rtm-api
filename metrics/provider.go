package metrics

import (
	"os"
	"time"

	stats "github.com/fukata/golang-stats-api-handler"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
	"github.com/swagchat/rtm-api/config"
	"github.com/swagchat/rtm-api/rtm"
)

type Metrics struct {
	Hostname       string                            `json:"hostname"`
	Stats          *stats.Stats                      `json:"stats"`
	AllCount       int                               `json:"allCount"`
	UserCount      int                               `json:"userCount"`
	EventCount     int                               `json:"eventCount"`
	EachUserCount  map[string]int                    `json:"eachUserCount,omitempty"`
	EachEventCount map[scpb.EventType]map[string]int `json:"eachEventCount,omitempty"`
	Timestamp      string                            `json:"timestamp"`
}

type provider interface {
	Run()
}

func Provider() provider {
	c := config.Config()

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
	c := config.Config()
	m := &Metrics{}

	hostname, _ := os.Hostname()
	m.Hostname = hostname
	// m.Stats = stats.GetStats()

	srv := rtm.Server()
	users := srv.Connection.Users()
	events := srv.Connection.Events()
	m.AllCount = srv.Connection.ConnectionCount()
	m.UserCount = len(users)
	m.EventCount = len(events)
	m.EachUserCount = srv.Connection.EachUserCount()
	m.EachEventCount = srv.Connection.EachEventCount()
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
