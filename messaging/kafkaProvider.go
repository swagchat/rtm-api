package messaging

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/swagchat/rtm-api/logging"
	"github.com/swagchat/rtm-api/services"
	"github.com/swagchat/rtm-api/utils"
	"go.uber.org/zap/zapcore"
)

var KafkaConsumer *kafka.Consumer

type KafkaProvider struct{}

func (provider *KafkaProvider) Subscribe() {
	logging.Log(zapcore.ErrorLevel, &logging.AppLog{
		Kind:     "messaging-subscribe",
		Provider: "kafka",
	})

	cfg := utils.Config()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	var hostname string
	hostname, err := os.Hostname()
	if err != nil {
		hostname = utils.CreateUuid()
	}
	KafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", cfg.Messaging.Kafka.Host, cfg.Messaging.Kafka.Port),
		"group.id":          hostname,
		// "session.timeout.ms":   6000,
		// "default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"}
	})

	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Kind:     "messaging-error",
			Provider: "kafka",
			Message:  err.Error(),
		})
	}

	fmt.Printf("Created Consumer %v\n", KafkaConsumer)

	err = KafkaConsumer.SubscribeTopics([]string{cfg.Messaging.Kafka.Topic}, nil)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Kind:     "messaging-error",
			Provider: "kafka",
			Message:  err.Error(),
		})
	}

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			run = false
			logging.Log(zapcore.InfoLevel, &logging.AppLog{
				Kind:     "messaging-terminate",
				Provider: "kafka",
				Message:  sig.String(),
			})
		default:
			ev := KafkaConsumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				services.Srv.Broadcast <- e.Value
				logging.Log(zapcore.InfoLevel, &logging.AppLog{
					Kind:     "messaging-receive",
					Provider: "kafka",
					Message:  string(e.Value),
				})
			case kafka.PartitionEOF:
				logging.Log(zapcore.InfoLevel, &logging.AppLog{
					Kind:     "messaging-reached",
					Provider: "kafka",
					Message:  e.String(),
				})
			case kafka.Error:
				run = false
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Kind:     "messaging-error",
					Provider: "kafka",
					Message:  e.String(),
				})
			default:
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					Kind:     "messaging-ignored",
					Provider: "kafka",
					Message:  e.String(),
				})
			}
		}
	}

	KafkaConsumer.Close()
	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Kind:     "messaging-close",
		Provider: "kafka",
	})
}

func (provider *KafkaProvider) Unsubscribe() {
	logging.Log(zapcore.InfoLevel, &logging.AppLog{
		Kind:     "messaging-unsubscribe",
		Provider: "kafka",
	})
	KafkaConsumer.Unsubscribe()
}
