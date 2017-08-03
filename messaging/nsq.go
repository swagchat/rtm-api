package messaging

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"unsafe"

	"github.com/fairway-corp/swagchat-realtime/services"
	"github.com/fairway-corp/swagchat-realtime/utils"
	nsq "github.com/nsqio/go-nsq"
)

func b2s(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}
func hello(s string) {
	fmt.Println(s)
}

func InitConsumer() {
	config := nsq.NewConfig()
	channel := utils.Que.Channel
	hostname, err := os.Hostname()
	if err == nil {
		config.Hostname = hostname
		channel = hostname
	}
	c, _ := nsq.NewConsumer(utils.Que.Topic, channel, config)
	c.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %v", message)
		services.Send(message.Body)
		return nil
	}))
	//err = c.ConnectToNSQDs([]string{"127.0.0.1:4150", "127.0.0.1:5150"})
	err = c.ConnectToNSQLookupd(utils.Que.Host + ":" + utils.Que.Port)
	if err != nil {
		log.Panic("Could not connect")
	}
}
