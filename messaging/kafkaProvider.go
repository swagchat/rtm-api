package messaging

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/swagchat/rtm-api/services"
	"github.com/swagchat/rtm-api/utils"
)

var KafkaConsumer *kafka.Consumer

type KafkaProvider struct{}

func (provider KafkaProvider) Init() error {
	return nil
}

func (provider KafkaProvider) Subscribe() {
	cfg := utils.GetConfig()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	var hostname string
	hostname, err := os.Hostname()
	if err != nil {
		hostname = utils.CreateUuid()
	}
	KafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", cfg.Kafka.Host, cfg.Kafka.Port),
		"group.id":          hostname,
		// "session.timeout.ms":   6000,
		// "default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"}
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", KafkaConsumer)

	err = KafkaConsumer.SubscribeTopics([]string{cfg.Kafka.Topic}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe topics: %s\n", err)
		os.Exit(1)
	}

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := KafkaConsumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				services.Srv.Broadcast <- e.Value
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	KafkaConsumer.Close()
}

func (provider KafkaProvider) Unsubscribe() {
	log.Println("kafka Unsubscribe")
	KafkaConsumer.Unsubscribe()
}
