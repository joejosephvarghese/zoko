package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/joejosephvarghese/message/server/pkg/domain"
	usecase "github.com/joejosephvarghese/message/server/pkg/usecase/interfaces"
	"github.com/segmentio/kafka-go"
)

// KafkaHandler handles Kafka messages using injected dependencies
type KafkaHandler struct {
	chatUC usecase.ChatUseCase
}

// NewKafkaHandler returns a new instance of KafkaHandler
func NewKafkaHandler(chatUC usecase.ChatUseCase) *KafkaHandler {
	return &KafkaHandler{chatUC: chatUC}
}

// HandleMessage processes an incoming Kafka message
func (h *KafkaHandler) HandleMessage(msg kafka.Message) {
	var message domain.Message

	if err := json.Unmarshal(msg.Value, &message); err != nil {
		log.Printf("Failed to unmarshal Kafka message: %v", err)
		return
	}

	receiverID, err := h.chatUC.SaveMessage(context.Background(), message)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return
	}

	log.Printf("Saved message. Receiver ID: %d\n", receiverID)
}
