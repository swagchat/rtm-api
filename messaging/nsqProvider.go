package messaging

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"unsafe"

	nsq "github.com/nsqio/go-nsq"
	"github.com/swagchat/rtm-api/logging"
	"github.com/swagchat/rtm-api/rtm"
	"github.com/swagchat/rtm-api/utils"
	"go.uber.org/zap/zapcore"
)

var NSQConsumer *nsq.Consumer

type nsqProvider struct{}

func b2s(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func (np *nsqProvider) Subscribe() {
	c := utils.Config()
	if c.Messaging.NSQ.NsqlookupdHost != "" {
		config := nsq.NewConfig()
		channel := c.Messaging.NSQ.Channel
		hostname, err := os.Hostname()
		if err == nil {
			config.Hostname = hostname
			channel = hostname
		}
		NSQConsumer, err = nsq.NewConsumer(c.Messaging.NSQ.Topic, channel, config)
		if err != nil {
			logging.Log(zapcore.ErrorLevel, &logging.AppLog{
				Kind:     "messaging-subscribe",
				Provider: "nsq",
				Message:  err.Error(),
			})
		} else {
			logging.Log(zapcore.InfoLevel, &logging.AppLog{
				Kind:     "messaging-subscribe",
				Provider: "nsq",
				Message:  fmt.Sprintf("%p", NSQConsumer),
			})
		}
		NSQConsumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
			rtm.Server().Broadcast <- message.Body
			logging.Log(zapcore.InfoLevel, &logging.AppLog{
				Kind:     "messaging-subscribe-receive",
				Provider: "nsq",
				Message:  string(message.Body),
			})
			return nil
		}))
		err = NSQConsumer.ConnectToNSQLookupd(c.Messaging.NSQ.NsqlookupdHost + ":" + c.Messaging.NSQ.NsqlookupdPort)
		if err != nil {
			logging.Log(zapcore.ErrorLevel, &logging.AppLog{
				Kind:     "messaging-subscribe",
				Provider: "nsq",
				Message:  err.Error(),
			})
		}
	}
}

func (np *nsqProvider) Unsubscribe() {
	if NSQConsumer != nil {
		c := utils.Config()
		hostname, err := os.Hostname()
		_, err = http.Post("http://"+c.Messaging.NSQ.NsqdHost+":"+c.Messaging.NSQ.NsqdPort+"/channel/delete?topic="+c.Messaging.NSQ.Topic+"&channel="+hostname, "text/plain", nil)
		if err != nil {
			logging.Log(zapcore.ErrorLevel, &logging.AppLog{
				Kind:     "messaging-error",
				Provider: "nsq",
				Message:  err.Error(),
			})
		}
	}
}
