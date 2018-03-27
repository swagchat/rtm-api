package utils

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

const (
	// AppName is Application name
	AppName = "rtm-api"
	// APIVersion is API version
	APIVersion = "v0"
	// BuildVersion is API build version
	BuildVersion = "v0.3.0"

	MAX_MESSAGE_SIZE = 8192
)

var (
	cfg           *config = NewConfig()
	IsShowVersion bool
)

type config struct {
	Version      string
	HttpPort     string `yaml:"httpPort"`
	Profiling    bool
	ErrorLogging bool `yaml:"errorLogging"`
	Logging      *Logging
	Realtime     *RealtimeSetting

	MessagingProvider string `yaml:"messagingProvider"`
	NSQ               *NSQ
	Kafka             *Kafka

	Metrics *Metrics
}

type Logging struct {
	Level string
}

type RealtimeSetting struct {
	IsDisplayConnectionInfo bool
}

type NSQ struct {
	Port           string
	NsqlookupdHost string
	NsqlookupdPort string
	NsqdHost       string
	NsqdPort       string
	Topic          string
	Channel        string
}

type Kafka struct {
	Host    string
	Port    string
	GroupID string `yaml:"groupId"`
	Topic   string
}

type Metrics struct {
	Provider string
	Interval int
	Verbose  bool
	Stdout   struct {
		Interval int
	}
	Elasticsearch struct {
		URL      string
		UserID   string `yaml:"userId"`
		Password string
	}
}

func NewConfig() *config {
	log.SetFlags(log.Llongfile)

	logging := &Logging{
		Level: "development",
	}

	realtimeSetting := &RealtimeSetting{
		IsDisplayConnectionInfo: true,
	}

	nsq := &NSQ{}
	kafka := &Kafka{}

	metrics := &Metrics{}

	c := &config{
		Version:           "0",
		HttpPort:          "8102",
		Profiling:         false,
		ErrorLogging:      false,
		Logging:           logging,
		Realtime:          realtimeSetting,
		MessagingProvider: "",
		NSQ:               nsq,
		Kafka:             kafka,
		Metrics:           metrics,
	}

	c.LoadYaml()
	c.LoadEnvironment()
	c.ParseFlag()

	return c
}

func GetConfig() *config {
	return cfg
}

func (c *config) LoadYaml() {
	buf, _ := ioutil.ReadFile("config/app.yaml")
	yaml.Unmarshal(buf, c)
}

func (c *config) LoadEnvironment() {
	var v string

	if v = os.Getenv("HTTP_PORT"); v != "" {
		c.HttpPort = v
	}
	if v = os.Getenv("RTM_PORT"); v != "" {
		c.HttpPort = v
	}
	if v = os.Getenv("RTM_PROFILING"); v != "" {
		if v == "true" {
			c.Profiling = true
		} else if v == "false" {
			c.Profiling = false
		}
	}
	if v = os.Getenv("RTM_ERROR_LOGGING"); v != "" {
		if v == "true" {
			c.ErrorLogging = true
		} else if v == "false" {
			c.ErrorLogging = false
		}
	}

	// Logging
	if v = os.Getenv("RTM_LOGGING_LEVEL"); v != "" {
		c.Logging.Level = v
	}

	if v = os.Getenv("RTM_MESSAGING_PROVIDER"); v != "" {
		c.MessagingProvider = v
	}

	// kafka
	if v = os.Getenv("RTM_KAFKA_HOST"); v != "" {
		c.Kafka.Host = v
	}
	if v = os.Getenv("RTM_KAFKA_PORT"); v != "" {
		c.Kafka.Port = v
	}
	if v = os.Getenv("RTM_KAFKA_GROUPID"); v != "" {
		c.Kafka.GroupID = v
	}
	if v = os.Getenv("RTM_KAFKA_TOPIC"); v != "" {
		c.Kafka.Topic = v
	}

	// metrics
	if v = os.Getenv("RTM_METRICS_PROVIDER"); v != "" {
		c.Metrics.Provider = v
	}

	if v = os.Getenv("RTM_METRICS_INTERVAL"); v != "" {
		interval, _ := strconv.Atoi(v)
		c.Metrics.Interval = interval
	}

	if v = os.Getenv("RTM_METRICS_VERBOSE"); v != "" {
		if v == "true" {
			c.Metrics.Verbose = true
		}
	}

	// elasticsearch
	if v = os.Getenv("RTM_ELASTICSEARCH_URL"); v != "" {
		c.Metrics.Elasticsearch.URL = v
	}
	if v = os.Getenv("RTM_ELASTICSEARCH_USERID"); v != "" {
		c.Metrics.Elasticsearch.UserID = v
	}
	if v = os.Getenv("RTM_ELASTICSEARCH_PASSWORD"); v != "" {
		c.Metrics.Elasticsearch.Password = v
	}
}

func (c *config) ParseFlag() {
	flag.BoolVar(&IsShowVersion, "v", false, "")
	flag.BoolVar(&IsShowVersion, "version", false, "show version")

	flag.StringVar(&c.HttpPort, "httpPort", c.HttpPort, "")

	var profiling string
	flag.StringVar(&profiling, "profiling", "", "")

	var demoPage string
	flag.StringVar(&demoPage, "demoPage", "", "false")

	var errorLogging string
	flag.StringVar(&errorLogging, "errorLogging", "", "false")

	// Logging
	flag.StringVar(&c.Logging.Level, "logging.level", c.Logging.Level, "")

	// kafka
	flag.StringVar(&c.Kafka.Host, "kafka.host", c.Kafka.Host, "")
	flag.StringVar(&c.Kafka.Port, "kafka.port", c.Kafka.Port, "")
	flag.StringVar(&c.Kafka.GroupID, "kafka.groupId", c.Kafka.GroupID, "")
	flag.StringVar(&c.Kafka.Topic, "kafka.topic", c.Kafka.Topic, "")

	// metrics
	flag.StringVar(&c.Metrics.Provider, "metrics.provider", c.Metrics.Provider, "")
	flag.IntVar(&c.Metrics.Interval, "metrics.interval", c.Metrics.Interval, "")
	flag.BoolVar(&c.Metrics.Verbose, "metrics.verbose", c.Metrics.Verbose, "")

	// elasticsearch
	flag.StringVar(&c.Metrics.Elasticsearch.URL, "elasticsearch.url", c.Metrics.Elasticsearch.URL, "")
	flag.StringVar(&c.Metrics.Elasticsearch.UserID, "elasticsearch.userId", c.Metrics.Elasticsearch.UserID, "")
	flag.StringVar(&c.Metrics.Elasticsearch.Password, "elasticsearch.password", c.Metrics.Elasticsearch.Password, "")

	var isDisplayConnectionInfo string
	flag.StringVar(&isDisplayConnectionInfo, "isDisplayConnectionInfo", "", "Display connection info.")
	if profiling == "true" {
		c.Realtime.IsDisplayConnectionInfo = true
	} else if profiling == "false" {
		c.Realtime.IsDisplayConnectionInfo = false
	}

	flag.StringVar(&c.MessagingProvider, "messagingProvider", c.MessagingProvider, "")

	flag.StringVar(&c.NSQ.NsqlookupdHost, "nsqlookupdHost", c.NSQ.NsqlookupdHost, "Host name of nsqlookupd")
	flag.StringVar(&c.NSQ.NsqlookupdPort, "nsqlookupdPort", c.NSQ.NsqlookupdPort, "Port no of nsqlookupd")
	flag.StringVar(&c.NSQ.NsqdHost, "nsqdHost", c.NSQ.NsqdHost, "Host name of nsqd")
	flag.StringVar(&c.NSQ.NsqdPort, "nsqdPort", c.NSQ.NsqdPort, "Port no of nsqd")
	flag.StringVar(&c.NSQ.Topic, "topic", c.NSQ.Topic, "Topic name")
	flag.StringVar(&c.NSQ.Channel, "channel", c.NSQ.Channel, "Channel name. If it's not set, channel is hostname.")
	flag.Parse()
}
