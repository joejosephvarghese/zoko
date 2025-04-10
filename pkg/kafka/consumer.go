package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/joejosephvarghese/message/server/pkg/config"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer initializes a new Kafka reader
func NewConsumer(cfg config.Config) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.BrokerAddress},
		Topic:    cfg.Topic,
		GroupID:  cfg.GroupID, // used for consumer group coordination
		MinBytes: 1,           // smallest batch of data we’ll fetch
		MaxBytes: 10e6,        // max batch size (10MB)
	})

	return &Consumer{reader: reader}
}

// Consume reads messages from Kafka topic
func (c *Consumer) Consume(ctx context.Context, handler func(msg kafka.Message)) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		fmt.Printf("Received message: %s => %s\n", string(m.Key), string(m.Value))
		handler(m)
	}
}

// Close closes the Kafka reader connection
func (c *Consumer) Close() error {
	return c.reader.Close()
}
