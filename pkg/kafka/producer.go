package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/joejosephvarghese/message/server/pkg/config"
	producerinterface "github.com/joejosephvarghese/message/server/pkg/kafka/producerInterface"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

// NewProducer initializes the Kafka writer
func NewProducer(cfg config.Config) producerinterface.ProdInterInterface {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.BrokerAddress),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}
}

// Send sends a Kafka message
func (p *Producer) Send(ctx context.Context, key string, value []byte) error {

	err := p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: value,
		},
	)

	if err != nil {
		log.Printf("Kafka send error: %v", err)
		return err
	}

	fmt.Println("Message sent to Kafka!")
	return nil
}

// Close should be called when shutting down the app
func (p *Producer) Close() error {
	return p.writer.Close()
}
