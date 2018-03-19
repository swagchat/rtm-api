package messaging

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"unsafe"

	nsq "github.com/nsqio/go-nsq"
	"github.com/swagchat/rtm-api/services"
	"github.com/swagchat/rtm-api/utils"
)

var Con *nsq.Consumer

type NsqProvider struct{}

func (provider NsqProvider) Init() error {
	return nil
}

func b2s(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func hello(s string) {
	fmt.Println(s)
}

func (provider NsqProvider) Subscribe() {
	c := utils.GetConfig()
	if c.NSQ.NsqlookupdHost != "" {
		config := nsq.NewConfig()
		channel := c.NSQ.Channel
		hostname, err := os.Hostname()
		if err == nil {
			config.Hostname = hostname
			channel = hostname
		}
		Con, _ = nsq.NewConsumer(c.NSQ.Topic, channel, config)
		Con.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
			log.Printf("[NSQ]Got a message: %v", message)
			services.Srv.Broadcast <- message.Body
			return nil
		}))
		err = Con.ConnectToNSQLookupd(c.NSQ.NsqlookupdHost + ":" + c.NSQ.NsqlookupdPort)
		if err != nil {
			log.Panic("Could not connect")
		}
	}
}

func (provider NsqProvider) Unsubscribe() {
	if Con != nil {
		c := utils.GetConfig()
		hostname, err := os.Hostname()
		resp, err := http.Post("http://"+c.NSQ.NsqdHost+":"+c.NSQ.NsqdPort+"/channel/delete?topic="+c.NSQ.Topic+"&channel="+hostname, "text/plain", nil)
		if err != nil {
			log.Println(err)
		}
		log.Println(resp)
	}
}
