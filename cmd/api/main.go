package main

import (
	"context"
	"log"

	"github.com/joejosephvarghese/message/server/pkg/config"
	"github.com/joejosephvarghese/message/server/pkg/db"
	"github.com/joejosephvarghese/message/server/pkg/di"
	"github.com/joejosephvarghese/message/server/pkg/kafka"
	"github.com/joejosephvarghese/message/server/pkg/repository"
	"github.com/joejosephvarghese/message/server/pkg/usecase"
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
	gormDB, err := db.ConnectDatabase(config)
	if err != nil {
		return
	}
	chatRepository := repository.NewChatRepository(gormDB)
	chatUseCase := usecase.NewChatUseCase(chatRepository)
	kafkaHadl := kafka.NewKafkaHandler(chatUseCase)
	consumer := kafka.NewConsumer(config)

	// ✅ Start consuming in a goroutine

	go consumer.Consume(context.Background(), kafkaHadl.HandleMessage)

	if err := server.Start(); err != nil {
		log.Fatal("failed to start server: ", err)
	}

}
