package main

import (
	"context"
	"log"

	"github.com/joejosephvarghese/message/server/pkg/config"
	"github.com/joejosephvarghese/message/server/pkg/di"
	"github.com/joejosephvarghese/message/server/pkg/kafka"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	server, err := di.InitializeAPI(config)
	if err != nil {
		log.Fatal("failed to initialize server: ", err)
	}

	// ✅ Manually create Kafka consumer instance
	consumer := kafka.NewConsumer(config)
	// ✅ Start consuming in a goroutine

	go consumer.Consume(context.Background(), kafka.HandleKafkaMessage)

	if err := server.Start(); err != nil {
		log.Fatal("failed to start server: ", err)
	}

}
