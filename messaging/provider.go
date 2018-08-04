package messaging

import (
	"github.com/swagchat/rtm-api/config"
)

type MessagingInfo struct {
	Message string
}

type provider interface {
	Subscribe()
	Unsubscribe()
}

func Provider() provider {
	c := config.Config()

	var p provider
	switch c.Messaging.Provider {
	case "":
		p = &notuseProvider{}
	case "kafka":
		p = &kafkaProvider{}
	case "nsq":
		p = &nsqProvider{}
	}
	return p
}
