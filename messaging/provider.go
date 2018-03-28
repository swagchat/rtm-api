package messaging

import (
	"github.com/swagchat/rtm-api/utils"
)

type MessagingInfo struct {
	Message string
}

type Provider interface {
	Subscribe()
	Unsubscribe()
}

func MessagingProvider() Provider {
	c := utils.Config()

	var provider Provider
	switch c.Messaging.Provider {
	case "":
		provider = &NotuseProvider{}
	case "kafka":
		provider = &KafkaProvider{}
	case "nsq":
		provider = &NsqProvider{}
	}
	return provider
}
