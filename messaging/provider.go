package messaging

import (
	"os"

	"github.com/swagchat/rtm-api/utils"
	"go.uber.org/zap"
)

type MessagingInfo struct {
	Message string
}

type Provider interface {
	Init() error
	Subscribe()
	Unsubscribe()
}

func GetMessagingProvider() Provider {
	c := utils.GetConfig()

	var provider Provider
	switch c.MessagingProvider {
	case "":
		provider = &NotUseProvider{}
	case "kafka":
		provider = &KafkaProvider{}
	case "nsq":
		provider = &NsqProvider{}
	default:
		utils.AppLogger.Error("",
			zap.String("msg", "utils.Cfg.Rtm.Provider is incorrect"),
		)
		os.Exit(0)
	}
	return provider
}
