package kafkahadlerinterfaces

import "github.com/segmentio/kafka-go"

type KafkaHandlerInterface interface {
	HandleMessage(msg kafka.Message)
}
