package config

import (
	"errors"
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

	CtxSubscription ctxKey = iota
	CtxTracerTransaction
	CtxTracerSpan
)

var (
	cfg         = NewConfig()
	showVersion = false
	showHelp    = false
	// StopRun is a flag for stop run server
	StopRun = false
)

type ctxKey int

type config struct {
	Version         string
	HTTPPort        string `yaml:"httpPort"`
	Profiling       bool
	ErrorLogging    bool `yaml:"errorLogging"`
	ReadBufferSize  int
	WriteBufferSize int
	MaxMessageSize  int64
	Logging         *Logging
	Logger          *Logger
	Tracer          *Tracer
	Metrics         *Metrics
	SBroker         *SBroker `yaml:"sBroker"`
}

// Logger is settings of logger
type Logging struct {
	Level string
}

// Logger is settings of logger
type Logger struct {
	// EnableConsole is a flag for enable console log.
	EnableConsole bool `yaml:"enableConsole"`
	// ConsoleFormat is a format for console log.
	ConsoleFormat string `yaml:"consoleFormat"`
	// ConsoleLevel is a level for console log.
	ConsoleLevel string `yaml:"consoleLevel"`
	// EnableFile is a flag for enable file log.
	EnableFile bool `yaml:"enableFile"`
	// FileFormat is a format for file log.
	FileFormat string `yaml:"fileFormat"`
	// FileLevel is a log level for file log.
	FileLevel string `yaml:"fileLevel"`
	// FilePath is a file path for file log.
	FilePath string `yaml:"filePath"`
}

// Tracer is settings of tracer
type Tracer struct {
	Provider string
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
	Provider      string
	Interval      int
	Verbose       bool
	Elasticsearch struct {
		URL             string
		UserID          string `yaml:"userId"`
		Password        string
		Index           string
		IndexTimeFormat string
		Type            string
	}
}

type SBroker struct {
	Provider string

	Kafka struct {
		Host    string
		Port    string
		GroupID string `yaml:"groupId"`
		Topic   string
	}

	NSQ struct {
		Port           string
		NsqlookupdHost string `yaml:"nsqLookupdHost"`
		NsqlookupdPort string `yaml:"nsqLookupdPort"`
		NsqdHost       string `yaml:"nsqdHost"`
		NsqdPort       string `yaml:"nsqdPort"`
		Topic          string
		Channel        string
	}
}

func NewConfig() *config {
	log.SetFlags(log.Llongfile)

	c := defaultSetting()
	c.loadEnv()

	c.loadYaml()
	c.parseFlag()

	err := c.validate()
	if err != nil {
		log.Fatalf("Invalid setting. %v", err)
	}

	err = c.after()
	if err != nil {
		log.Fatalf("%v", err)
	}

	return c
}

func Config() *config {
	return cfg
}

func defaultSetting() *config {
	logging := &Logging{
		Level: "development",
	}

	c := &config{
		Version:         "0",
		HTTPPort:        "8102",
		Profiling:       false,
		ErrorLogging:    false,
		ReadBufferSize:  8192,
		WriteBufferSize: 1024,
		MaxMessageSize:  8192,
		Logging:         logging,
		Logger: &Logger{
			EnableConsole: true,
			ConsoleFormat: "text",
			ConsoleLevel:  "debug",
			EnableFile:    false,
		},
		Tracer:  &Tracer{},
		Metrics: &Metrics{},
		SBroker: &SBroker{},
	}

	return c
}

func (c *config) loadYaml() {
	buf, _ := ioutil.ReadFile("config/app.yaml")
	yaml.Unmarshal(buf, c)
}

func (c *config) loadEnv() {
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

	// Logging
	flag.BoolVar(&c.Logger.EnableConsole, "logger.enableConsole", c.Logger.EnableConsole, "")
	flag.StringVar(&c.Logger.ConsoleFormat, "logger.consoleFormat", c.Logger.ConsoleFormat, "")
	flag.StringVar(&c.Logger.ConsoleLevel, "logger.consoleLevel", c.Logger.ConsoleLevel, "")
	flag.BoolVar(&c.Logger.EnableFile, "logger.enableFile", c.Logger.EnableFile, "")
	flag.StringVar(&c.Logger.FileFormat, "logger.fileFormat", c.Logger.FileFormat, "")
	flag.StringVar(&c.Logger.FileLevel, "logger.fileLevel", c.Logger.FileLevel, "")
	flag.StringVar(&c.Logger.FilePath, "logger.filePath", c.Logger.FilePath, "")

	// Logger
	if v = os.Getenv("SWAG_LOGGER_ENABLE_CONSOLE"); v == "true" {
		c.Logger.EnableConsole = true
	}
	if v = os.Getenv("SWAG_LOGGER_CONSOLE_FORMAT"); v != "" {
		c.Logger.ConsoleFormat = v
	}
	if v = os.Getenv("SWAG_LOGGER_CONSOLE_LEVEL"); v != "" {
		c.Logger.ConsoleLevel = v
	}
	if v = os.Getenv("SWAG_LOGGER_ENABLE_FILE"); v == "true" {
		c.Logger.EnableFile = true
	}
	if v = os.Getenv("SWAG_LOGGER_FILE_FORMAT"); v != "" {
		c.Logger.FileFormat = v
	}
	if v = os.Getenv("SWAG_LOGGER_FILE_LEVEL"); v != "" {
		c.Logger.FileLevel = v
	}
	if v = os.Getenv("SWAG_LOGGER_FILE_PATH"); v != "" {
		c.Logger.FilePath = v
	}

	// Tracer
	if v = os.Getenv("SWAG_TRACER_PROVIDER"); v != "" {
		c.Tracer.Provider = v
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

	// SBroker
	if v = os.Getenv("SWAG_SBROKER_PROVIDER"); v != "" {
		c.SBroker.Provider = v
	}

	// SBroker - Kafka
	if v = os.Getenv("SWAG_SBROKER_KAFKA_HOST"); v != "" {
		c.SBroker.Kafka.Host = v
	}
	if v = os.Getenv("SWAG_SBROKER_KAFKA_PORT"); v != "" {
		c.SBroker.Kafka.Port = v
	}
	if v = os.Getenv("SWAG_SBROKER_KAFKA_GROUPID"); v != "" {
		c.SBroker.Kafka.GroupID = v
	}
	if v = os.Getenv("SWAG_SBROKER_KAFKA_TOPIC"); v != "" {
		c.SBroker.Kafka.Topic = v
	}

	// SBroker - NSQ
	if v = os.Getenv("SWAG_SBROKER_NSQ_PORT"); v != "" {
		c.SBroker.NSQ.Port = v
	}
	if v = os.Getenv("SWAG_SBROKER_NSQ_NSQLOOKUPDHOST"); v != "" {
		c.SBroker.NSQ.NsqlookupdHost = v
	}
	if v = os.Getenv("SWAG_SBROKER_NSQ_NSQLOOKUPDPORT"); v != "" {
		c.SBroker.NSQ.NsqlookupdPort = v
	}
	if v = os.Getenv("SWAG_SBROKER_NSQ_NSQDHOST"); v != "" {
		c.SBroker.NSQ.NsqdHost = v
	}
	if v = os.Getenv("SWAG_SBROKER_NSQ_NSQDPORT"); v != "" {
		c.SBroker.NSQ.NsqdPort = v
	}
	if v = os.Getenv("SWAG_SBROKER_NSQ_TOPIC"); v != "" {
		c.SBroker.NSQ.Topic = v
	}
	if v = os.Getenv("SWAG_SBROKER_NSQ_CHANNEL"); v != "" {
		c.SBroker.NSQ.Channel = v
	}
}

func (c *config) parseFlag() {
	flag.BoolVar(&showVersion, "v", false, "")
	flag.BoolVar(&showVersion, "version", false, "show version")

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

	// Tracer
	flag.StringVar(&c.Tracer.Provider, "tracer.provider", c.Tracer.Provider, "")

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

	// SBroker
	flag.StringVar(&c.SBroker.Provider, "sbroker.provider", c.SBroker.Provider, "")

	// SBroker - kafka
	flag.StringVar(&c.SBroker.Kafka.Host, "sbroker.kafka.host", c.SBroker.Kafka.Host, "")
	flag.StringVar(&c.SBroker.Kafka.Port, "sbroker.kafka.port", c.SBroker.Kafka.Port, "")
	flag.StringVar(&c.SBroker.Kafka.GroupID, "sbroker.kafka.groupId", c.SBroker.Kafka.GroupID, "")
	flag.StringVar(&c.SBroker.Kafka.Topic, "sbroker.kafka.topic", c.SBroker.Kafka.Topic, "")

	// SBroker - NSQ
	flag.StringVar(&c.SBroker.NSQ.NsqlookupdHost, "sbroker.nsq.nsqlookupdHost", c.SBroker.NSQ.NsqlookupdHost, "Host name of nsqlookupd")
	flag.StringVar(&c.SBroker.NSQ.NsqlookupdPort, "sbroker.nsq.nsqlookupdPort", c.SBroker.NSQ.NsqlookupdPort, "Port no of nsqlookupd")
	flag.StringVar(&c.SBroker.NSQ.NsqdHost, "sbroker.nsq.nsqdHost", c.SBroker.NSQ.NsqdHost, "Host name of nsqd")
	flag.StringVar(&c.SBroker.NSQ.NsqdPort, "sbroker.nsq.nsqdPort", c.SBroker.NSQ.NsqdPort, "Port no of nsqd")
	flag.StringVar(&c.SBroker.NSQ.Topic, "sbroker.nsq.topic", c.SBroker.NSQ.Topic, "Topic name")
	flag.StringVar(&c.SBroker.NSQ.Channel, "sbroker.nsq.channel", c.SBroker.NSQ.Channel, "Channel name. If it's not set, channel is hostname.")

	flag.Parse()
}

func (c *config) validate() error {
	// Logger
	if c.Logger.EnableConsole {
		f := c.Logger.ConsoleFormat
		if f == "" || !(f == "text" || f == "json") {
			return errors.New("Please set logger.consoleFormat to \"text\" or \"json\"")
		}
		l := c.Logger.ConsoleLevel
		if l == "" || !(l == "debug" || l == "info" || l == "warn" || l == "error") {
			return errors.New("Please set logger.consoleLevel to \"debug\" or \"info\" or \"warn\" or \"error\"")
		}
	}
	if c.Logger.EnableFile {
		f := c.Logger.FileFormat
		if f == "" || !(f == "text" || f == "json") {
			return errors.New("Please set logger.fileFormat to \"text\" or \"json\"")
		}
		l := c.Logger.FileLevel
		if l == "" || !(l == "debug" || l == "info" || l == "warn" || l == "error") {
			return errors.New("Please set logger.fileLevel to \"debug\" or \"info\" or \"warn\" or \"error\"")
		}
		if c.Logger.FilePath == "" {
			return errors.New("Please set logger.filePath")
		}
	}

	return nil
}

func (c *config) after() error {
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
	if c.Metrics.Interval == 0 {
		c.Metrics.Interval = 5
	}

	return nil
}
