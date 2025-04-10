package kafka

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

func HandleKafkaMessage(msg kafka.Message) {
	// Deserialize and handle the message
	fmt.Printf("Processing message: %s\n", string(msg.Value))
}
