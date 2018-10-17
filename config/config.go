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
	ReadBufferSize  int
	WriteBufferSize int
	MaxMessageSize  int64
	Logger          *Logger
	Tracer          *Tracer
	Metrics         *Metrics
	Consumer        *Consumer
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
	// FileMaxSize is a max size(megabytes) for file log. The Default is 100 megabytes.
	FileMaxSize int `yaml:"fileMaxSize"`
	// FileMaxAge is a max age(days) for file log. The default is 1 days.
	FileMaxAge int `yaml:"fileMaxAge"`
	// FileMaxBackups is a max backup file number for file log. The default is 10 files.
	FileMaxBackups int `yaml:"fileMaxBackups"`
	// FileLocalTime is a formatting the timestamps for backup files. The default is to use UTC time.
	FileLocalTime bool `yaml:"fileLocalTime"`
	// FileCompress is a compressed flag for file log. The default is false.
	FileCompress bool `yaml:"fileCompress"`
}

// Tracer is settings of tracer
type Tracer struct {
	Provider string
	Logging  bool

	Zipkin struct {
		Endpoint  string
		BatchSize int `yaml:"batchSize"`
		Timeout   int
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

type Consumer struct {
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

	err := c.parseFlag(os.Args[1:])
	if err != nil {
		log.Fatalf("Failed to load setting. %v", err)
	}

	err = c.validate()
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
	c := &config{
		Version:   "0",
		HTTPPort:  "8102",
		Profiling: false,
		// ErrorLogging:    false,
		ReadBufferSize:  8192,
		WriteBufferSize: 1024,
		MaxMessageSize:  8192,
		// Logging:         logging,
		Logger: &Logger{
			EnableConsole:  true,
			ConsoleFormat:  "text",
			ConsoleLevel:   "debug",
			EnableFile:     false,
			FileMaxSize:    100,
			FileMaxAge:     1,
			FileMaxBackups: 10,
			FileLocalTime:  false,
			FileCompress:   false,
		},
		Tracer: &Tracer{
			Logging: false,
		},
		Metrics:  &Metrics{},
		Consumer: &Consumer{},
	}

	return c
}

func (c *config) loadYaml(buf []byte) {
	err := yaml.Unmarshal(buf, c)
	if err != nil {
		log.Fatalf("Failed to load yaml file. %v", err)
	}
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
	if v = os.Getenv("SWAG_LOGGER_FILE_MAX_SIZE"); v != "" {
		fileMaxSize, err := strconv.Atoi(v)
		if err == nil {
			c.Logger.FileMaxSize = fileMaxSize
		}
	}
	if v = os.Getenv("SWAG_LOGGER_FILE_MAX_AGE"); v != "" {
		fileMaxAge, err := strconv.Atoi(v)
		if err == nil {
			c.Logger.FileMaxAge = fileMaxAge
		}
	}
	if v = os.Getenv("SWAG_LOGGER_FILE_MAX_BACKUPS"); v != "" {
		fileMaxBackups, err := strconv.Atoi(v)
		if err == nil {
			c.Logger.FileMaxBackups = fileMaxBackups
		}
	}
	if v = os.Getenv("SWAG_LOGGER_FILE_LOCAL_TIME"); v == "true" {
		c.Logger.FileLocalTime = true
	}
	if v = os.Getenv("SWAG_LOGGER_FILE_COMPRESS"); v == "true" {
		c.Logger.FileCompress = true
	}

	// Tracer
	if v = os.Getenv("SWAG_TRACER_PROVIDER"); v != "" {
		c.Tracer.Provider = v
	}
	if v = os.Getenv("SWAG_TRACER_LOGGING"); v == "true" {
		c.Tracer.Logging = true
	}
	if v = os.Getenv("SWAG_TRACER_ZIPKIN_ENDPOINT"); v != "" {
		c.Tracer.Zipkin.Endpoint = v
	}
	if v = os.Getenv("SWAG_TRACER_ZIPKIN_BATCHSIZE"); v != "" {
		batchSize, err := strconv.Atoi(v)
		if err == nil {
			c.Tracer.Zipkin.BatchSize = batchSize
		}
	}
	if v = os.Getenv("SWAG_TRACER_ZIPKIN_TIMEOUT"); v != "" {
		timeout, err := strconv.Atoi(v)
		if err == nil {
			c.Tracer.Zipkin.Timeout = timeout
		}
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

	// Consumer
	if v = os.Getenv("SWAG_CONSUMER_PROVIDER"); v != "" {
		c.Consumer.Provider = v
	}

	// Consumer - Kafka
	if v = os.Getenv("SWAG_CONSUMER_KAFKA_HOST"); v != "" {
		c.Consumer.Kafka.Host = v
	}
	if v = os.Getenv("SWAG_CONSUMER_KAFKA_PORT"); v != "" {
		c.Consumer.Kafka.Port = v
	}
	if v = os.Getenv("SWAG_CONSUMER_KAFKA_GROUPID"); v != "" {
		c.Consumer.Kafka.GroupID = v
	}
	if v = os.Getenv("SWAG_CONSUMER_KAFKA_TOPIC"); v != "" {
		c.Consumer.Kafka.Topic = v
	}

	// Consumer - NSQ
	if v = os.Getenv("SWAG_CONSUMER_NSQ_PORT"); v != "" {
		c.Consumer.NSQ.Port = v
	}
	if v = os.Getenv("SWAG_CONSUMER_NSQ_NSQLOOKUPDHOST"); v != "" {
		c.Consumer.NSQ.NsqlookupdHost = v
	}
	if v = os.Getenv("SWAG_CONSUMER_NSQ_NSQLOOKUPDPORT"); v != "" {
		c.Consumer.NSQ.NsqlookupdPort = v
	}
	if v = os.Getenv("SWAG_CONSUMER_NSQ_NSQDHOST"); v != "" {
		c.Consumer.NSQ.NsqdHost = v
	}
	if v = os.Getenv("SWAG_CONSUMER_NSQ_NSQDPORT"); v != "" {
		c.Consumer.NSQ.NsqdPort = v
	}
	if v = os.Getenv("SWAG_CONSUMER_NSQ_TOPIC"); v != "" {
		c.Consumer.NSQ.Topic = v
	}
	if v = os.Getenv("SWAG_CONSUMER_NSQ_CHANNEL"); v != "" {
		c.Consumer.NSQ.Channel = v
	}
}

func (c *config) parseFlag(args []string) error {
	if len(args) == 0 {
		return nil
	}

	flags := flag.NewFlagSet(fmt.Sprintf("%s Flags", AppName), flag.ContinueOnError)

	flags.BoolVar(&showVersion, "v", false, "")
	flags.BoolVar(&showVersion, "version", false, "show version")

	flags.StringVar(&c.HTTPPort, "httpPort", c.HTTPPort, "")

	var profiling string
	flags.StringVar(&profiling, "profiling", "", "")

	flag.IntVar(&c.ReadBufferSize, "readBufferSize", c.ReadBufferSize, "")
	flag.IntVar(&c.WriteBufferSize, "writeBufferSize", c.WriteBufferSize, "")
	flag.Int64Var(&c.MaxMessageSize, "maxMessageSize", c.MaxMessageSize, "")

	// Logging
	flags.BoolVar(&c.Logger.EnableConsole, "logger.enableConsole", c.Logger.EnableConsole, "")
	flags.StringVar(&c.Logger.ConsoleFormat, "logger.consoleFormat", c.Logger.ConsoleFormat, "")
	flags.StringVar(&c.Logger.ConsoleLevel, "logger.consoleLevel", c.Logger.ConsoleLevel, "")
	flags.BoolVar(&c.Logger.EnableFile, "logger.enableFile", c.Logger.EnableFile, "")
	flags.StringVar(&c.Logger.FileFormat, "logger.fileFormat", c.Logger.FileFormat, "")
	flags.StringVar(&c.Logger.FileLevel, "logger.fileLevel", c.Logger.FileLevel, "")
	flags.StringVar(&c.Logger.FilePath, "logger.filePath", c.Logger.FilePath, "")
	flags.IntVar(&c.Logger.FileMaxSize, "logger.fileMaxSize", c.Logger.FileMaxSize, "")
	flags.IntVar(&c.Logger.FileMaxAge, "logger.fileMaxAge", c.Logger.FileMaxAge, "")
	flags.IntVar(&c.Logger.FileMaxBackups, "logger.fileMaxBackups", c.Logger.FileMaxBackups, "")
	flags.BoolVar(&c.Logger.FileLocalTime, "logger.fileLocalTime", c.Logger.FileLocalTime, "")
	flags.BoolVar(&c.Logger.FileCompress, "logger.fileCompress", c.Logger.FileCompress, "")

	// Tracer
	flags.StringVar(&c.Tracer.Provider, "tracer.provider", c.Tracer.Provider, "")
	flags.BoolVar(&c.Tracer.Logging, "tracer.logging", c.Tracer.Logging, "")
	flags.StringVar(&c.Tracer.Zipkin.Endpoint, "tracer.zipkin.endpoint", c.Tracer.Zipkin.Endpoint, "")
	flags.IntVar(&c.Tracer.Zipkin.BatchSize, "tracer.zipkin.batchsize", c.Tracer.Zipkin.BatchSize, "")
	flags.IntVar(&c.Tracer.Zipkin.Timeout, "tracer.zipkin.timeout", c.Tracer.Zipkin.Timeout, "")

	// metrics
	flags.StringVar(&c.Metrics.Provider, "metrics.provider", c.Metrics.Provider, "")
	flag.IntVar(&c.Metrics.Interval, "metrics.interval", c.Metrics.Interval, "")
	flags.BoolVar(&c.Metrics.Verbose, "metrics.verbose", c.Metrics.Verbose, "")

	// metrics elasticsearch
	flags.StringVar(&c.Metrics.Elasticsearch.URL, "metrics.elasticsearch.url", c.Metrics.Elasticsearch.URL, "")
	flags.StringVar(&c.Metrics.Elasticsearch.UserID, "metrics.elasticsearch.userId", c.Metrics.Elasticsearch.UserID, "")
	flags.StringVar(&c.Metrics.Elasticsearch.Password, "metrics.elasticsearch.password", c.Metrics.Elasticsearch.Password, "")
	flags.StringVar(&c.Metrics.Elasticsearch.Index, "metrics.elasticsearch.index", c.Metrics.Elasticsearch.Index, "")
	flags.StringVar(&c.Metrics.Elasticsearch.IndexTimeFormat, "metrics.elasticsearch.indexTimeFormat", c.Metrics.Elasticsearch.IndexTimeFormat, "")
	flags.StringVar(&c.Metrics.Elasticsearch.Type, "metrics.elasticsearch.type", c.Metrics.Elasticsearch.Type, "")

	// Consumer
	flags.StringVar(&c.Consumer.Provider, "sbroker.provider", c.Consumer.Provider, "")

	// Consumer - kafka
	flags.StringVar(&c.Consumer.Kafka.Host, "sbroker.kafka.host", c.Consumer.Kafka.Host, "")
	flags.StringVar(&c.Consumer.Kafka.Port, "sbroker.kafka.port", c.Consumer.Kafka.Port, "")
	flags.StringVar(&c.Consumer.Kafka.GroupID, "sbroker.kafka.groupId", c.Consumer.Kafka.GroupID, "")
	flags.StringVar(&c.Consumer.Kafka.Topic, "sbroker.kafka.topic", c.Consumer.Kafka.Topic, "")

	// Consumer - NSQ
	flags.StringVar(&c.Consumer.NSQ.NsqlookupdHost, "sbroker.nsq.nsqlookupdHost", c.Consumer.NSQ.NsqlookupdHost, "Host name of nsqlookupd")
	flags.StringVar(&c.Consumer.NSQ.NsqlookupdPort, "sbroker.nsq.nsqlookupdPort", c.Consumer.NSQ.NsqlookupdPort, "Port no of nsqlookupd")
	flags.StringVar(&c.Consumer.NSQ.NsqdHost, "sbroker.nsq.nsqdHost", c.Consumer.NSQ.NsqdHost, "Host name of nsqd")
	flags.StringVar(&c.Consumer.NSQ.NsqdPort, "sbroker.nsq.nsqdPort", c.Consumer.NSQ.NsqdPort, "Port no of nsqd")
	flags.StringVar(&c.Consumer.NSQ.Topic, "sbroker.nsq.topic", c.Consumer.NSQ.Topic, "Topic name")
	flags.StringVar(&c.Consumer.NSQ.Channel, "sbroker.nsq.channel", c.Consumer.NSQ.Channel, "Channel name. If it's not set, channel is hostname.")

	configPath := ""
	flags.StringVar(&configPath, "config", "", "config file(yaml format)")

	if flag.Lookup("test.run") != nil { // for testing
		return nil
	}

	err := flags.Parse(args)
	if err != nil {
		return nil
	}

	if showHelp {
		flags.PrintDefaults()
		StopRun = true
		return nil
	}

	if showVersion {
		fmt.Printf("API Version %s\nBuild Version %s\n", APIVersion, BuildVersion)
		StopRun = true
		return nil
	}

	if profiling == "true" {
		c.Profiling = true
	}

	if configPath != "" {
		if !isExists(configPath) {
			return fmt.Errorf("File not found [%s]", configPath)
		}
		buf, _ := ioutil.ReadFile(configPath)
		c.loadYaml(buf)
	}

	return nil
}

func isExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
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

	// Tracer
	if c.Tracer.Provider == "zipkin" {
		if c.Tracer.Zipkin.Endpoint == "" {
			return errors.New("Please set tracer.zipkin.endpoint")
		}
	}

	return nil
}
