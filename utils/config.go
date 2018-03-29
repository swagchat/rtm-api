package utils

import (
	"flag"
	"fmt"
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
	APIVersion = "0"
	// BuildVersion is API build version
	BuildVersion = "0.3.0"
)

var (
	cfg           = NewConfig()
	IsShowVersion bool
)

type config struct {
	Version         string
	HTTPPort        string `yaml:"httpPort"`
	Profiling       bool
	ErrorLogging    bool `yaml:"errorLogging"`
	ReadBufferSize  int
	WriteBufferSize int
	MaxMessageSize  int64

	Logging   *Logging
	Messaging *Messaging
	Metrics   *Metrics
}

type Logging struct {
	Level string
}

type Messaging struct {
	Provider string
	NSQ      struct {
		Port           string
		NsqlookupdHost string
		NsqlookupdPort string
		NsqdHost       string
		NsqdPort       string
		Topic          string
		Channel        string
	}
	Kafka struct {
		Host    string
		Port    string
		GroupID string `yaml:"groupId"`
		Topic   string
	}
}

type Metrics struct {
	Provider string
	Interval int
	Verbose  bool
	Stdout   struct {
		Interval int
	}
	Elasticsearch struct {
		URL             string
		UserID          string `yaml:"userId"`
		Password        string
		Index           string
		IndexTimeFormat string
		Type            string
	}
}

func NewConfig() *config {
	log.SetFlags(log.Llongfile)

	logging := &Logging{
		Level: "development",
	}
	messaging := &Messaging{}
	metrics := &Metrics{}

	c := &config{
		Version:         "0",
		HTTPPort:        "8102",
		Profiling:       false,
		ErrorLogging:    false,
		ReadBufferSize:  8192,
		WriteBufferSize: 1024,
		MaxMessageSize:  8192,
		Logging:         logging,
		Messaging:       messaging,
		Metrics:         metrics,
	}

	c.LoadYaml()
	c.LoadEnvironment()
	c.ParseFlag()
	c.after()

	return c
}

func Config() *config {
	return cfg
}

func (c *config) LoadYaml() {
	buf, _ := ioutil.ReadFile("config/app.yaml")
	yaml.Unmarshal(buf, c)
}

func (c *config) LoadEnvironment() {
	var v string

	if v = os.Getenv("HTTP_PORT"); v != "" {
		c.HTTPPort = v
	}
	if v = os.Getenv("RTM_PORT"); v != "" {
		c.HTTPPort = v
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

	if v = os.Getenv("RTM_READ_BUFFER_SIZE"); v != "" {
		size, _ := strconv.Atoi(v)
		c.ReadBufferSize = size
	}

	if v = os.Getenv("RTM_WRITE_BUFFER_SIZE"); v != "" {
		size, _ := strconv.Atoi(v)
		c.WriteBufferSize = size
	}

	if v = os.Getenv("RTM_MAX_MESSAGE_SIZE"); v != "" {
		size, _ := strconv.ParseInt(v, 10, 64)
		c.MaxMessageSize = size
	}

	// Logging
	if v = os.Getenv("RTM_LOGGING_LEVEL"); v != "" {
		c.Logging.Level = v
	}

	// messaging
	if v = os.Getenv("RTM_MESSAGING_PROVIDER"); v != "" {
		c.Messaging.Provider = v
	}

	// messaging NSQ
	if v = os.Getenv("RTM_MESSAGING_NSQ_PORT"); v != "" {
		c.Messaging.NSQ.Port = v
	}
	if v = os.Getenv("RTM_MESSAGING_NSQ_NSQLOOKUPDHOST"); v != "" {
		c.Messaging.NSQ.NsqlookupdHost = v
	}
	if v = os.Getenv("RTM_MESSAGING_NSQ_NSQLOOKUPDPORT"); v != "" {
		c.Messaging.NSQ.NsqlookupdPort = v
	}
	if v = os.Getenv("RTM_MESSAGING_NSQ_NSQDHOST"); v != "" {
		c.Messaging.NSQ.NsqdHost = v
	}
	if v = os.Getenv("RTM_MESSAGING_NSQ_NSQDPORT"); v != "" {
		c.Messaging.NSQ.NsqdPort = v
	}
	if v = os.Getenv("RTM_MESSAGING_NSQ_TOPIC"); v != "" {
		c.Messaging.NSQ.Topic = v
	}
	if v = os.Getenv("RTM_MESSAGING_NSQ_CHANNEL"); v != "" {
		c.Messaging.NSQ.Channel = v
	}

	// messaging kafka
	if v = os.Getenv("RTM_MESSAGING_KAFKA_HOST"); v != "" {
		c.Messaging.Kafka.Host = v
	}
	if v = os.Getenv("RTM_MESSAGING_KAFKA_PORT"); v != "" {
		c.Messaging.Kafka.Port = v
	}
	if v = os.Getenv("RTM_MESSAGING_KAFKA_GROUPID"); v != "" {
		c.Messaging.Kafka.GroupID = v
	}
	if v = os.Getenv("RTM_MESSAGING_KAFKA_TOPIC"); v != "" {
		c.Messaging.Kafka.Topic = v
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

	// metrics elasticsearch
	if v = os.Getenv("RTM_METRICS_ELASTICSEARCH_URL"); v != "" {
		c.Metrics.Elasticsearch.URL = v
	}
	if v = os.Getenv("RTM_METRICS_ELASTICSEARCH_USERID"); v != "" {
		c.Metrics.Elasticsearch.UserID = v
	}
	if v = os.Getenv("RTM_METRICS_ELASTICSEARCH_PASSWORD"); v != "" {
		c.Metrics.Elasticsearch.Password = v
	}
	if v = os.Getenv("RTM_METRICS_ELASTICSEARCH_INDEX"); v != "" {
		c.Metrics.Elasticsearch.Index = v
	}
	if v = os.Getenv("RTM_METRICS_ELASTICSEARCH_INDEXTIMEFORMAT"); v != "" {
		c.Metrics.Elasticsearch.IndexTimeFormat = v
	}
	if v = os.Getenv("RTM_METRICS_ELASTICSEARCH_TYPE"); v != "" {
		c.Metrics.Elasticsearch.Type = v
	}
}

func (c *config) ParseFlag() {
	flag.BoolVar(&IsShowVersion, "v", false, "")
	flag.BoolVar(&IsShowVersion, "version", false, "show version")

	flag.StringVar(&c.HTTPPort, "httpPort", c.HTTPPort, "")

	var profiling string
	flag.StringVar(&profiling, "profiling", "", "")

	var errorLogging string
	flag.StringVar(&errorLogging, "errorLogging", "", "false")

	flag.IntVar(&c.ReadBufferSize, "readBufferSize", c.ReadBufferSize, "")
	flag.IntVar(&c.WriteBufferSize, "writeBufferSize", c.WriteBufferSize, "")
	flag.Int64Var(&c.MaxMessageSize, "maxMessageSize", c.MaxMessageSize, "")

	// Logging
	flag.StringVar(&c.Logging.Level, "logging.level", c.Logging.Level, "")

	// messaging
	flag.StringVar(&c.Messaging.Provider, "messaging.provider", c.Messaging.Provider, "")

	// messaging NSQ
	flag.StringVar(&c.Messaging.NSQ.NsqlookupdHost, "messaging.nsq.nsqlookupdHost", c.Messaging.NSQ.NsqlookupdHost, "Host name of nsqlookupd")
	flag.StringVar(&c.Messaging.NSQ.NsqlookupdPort, "messaging.nsq.nsqlookupdPort", c.Messaging.NSQ.NsqlookupdPort, "Port no of nsqlookupd")
	flag.StringVar(&c.Messaging.NSQ.NsqdHost, "messaging.nsq.nsqdHost", c.Messaging.NSQ.NsqdHost, "Host name of nsqd")
	flag.StringVar(&c.Messaging.NSQ.NsqdPort, "messaging.nsq.nsqdPort", c.Messaging.NSQ.NsqdPort, "Port no of nsqd")
	flag.StringVar(&c.Messaging.NSQ.Topic, "messaging.nsq.topic", c.Messaging.NSQ.Topic, "Topic name")
	flag.StringVar(&c.Messaging.NSQ.Channel, "messaging.nsq.channel", c.Messaging.NSQ.Channel, "Channel name. If it's not set, channel is hostname.")

	// messaging kafka
	flag.StringVar(&c.Messaging.Kafka.Host, "messaging.kafka.host", c.Messaging.Kafka.Host, "")
	flag.StringVar(&c.Messaging.Kafka.Port, "messaging.kafka.port", c.Messaging.Kafka.Port, "")
	flag.StringVar(&c.Messaging.Kafka.GroupID, "messaging.kafka.groupId", c.Messaging.Kafka.GroupID, "")
	flag.StringVar(&c.Messaging.Kafka.Topic, "messaging.kafka.topic", c.Messaging.Kafka.Topic, "")

	// metrics
	flag.StringVar(&c.Metrics.Provider, "metrics.provider", c.Metrics.Provider, "")
	flag.IntVar(&c.Metrics.Interval, "metrics.interval", c.Metrics.Interval, "")
	flag.BoolVar(&c.Metrics.Verbose, "metrics.verbose", c.Metrics.Verbose, "")

	// metrics elasticsearch
	flag.StringVar(&c.Metrics.Elasticsearch.URL, "metrics.elasticsearch.url", c.Metrics.Elasticsearch.URL, "")
	flag.StringVar(&c.Metrics.Elasticsearch.UserID, "metrics.elasticsearch.userId", c.Metrics.Elasticsearch.UserID, "")
	flag.StringVar(&c.Metrics.Elasticsearch.Password, "metrics.elasticsearch.password", c.Metrics.Elasticsearch.Password, "")
	flag.StringVar(&c.Metrics.Elasticsearch.Index, "metrics.elasticsearch.index", c.Metrics.Elasticsearch.Index, "")
	flag.StringVar(&c.Metrics.Elasticsearch.IndexTimeFormat, "metrics.elasticsearch.indexTimeFormat", c.Metrics.Elasticsearch.IndexTimeFormat, "")
	flag.StringVar(&c.Metrics.Elasticsearch.Type, "metrics.elasticsearch.type", c.Metrics.Elasticsearch.Type, "")

	flag.Parse()
}

func (c *config) after() {
	if c.Metrics.Provider == "elasticsearch" {
		if c.Metrics.Elasticsearch.Index == "" {
			c.Metrics.Elasticsearch.Index = fmt.Sprintf("%s-%s", AppName, "metrics")
		}
		if c.Metrics.Elasticsearch.IndexTimeFormat == "" {
			c.Metrics.Elasticsearch.IndexTimeFormat = "2006.01.02"
		}
		if c.Metrics.Elasticsearch.Type == "" {
			c.Metrics.Elasticsearch.Type = "_doc"
		}
	}
}
