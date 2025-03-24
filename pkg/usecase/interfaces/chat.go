package interfaces

import (
	"context"

	"github.com/joejosephvarghese/message/server/pkg/api/handler/request"
	"github.com/joejosephvarghese/message/server/pkg/api/handler/response"
	"github.com/joejosephvarghese/message/server/pkg/domain"
)

type ChatUseCase interface {
	FindAllRecentChatsOfUser(ctx context.Context, userID uint, pagination request.Pagination) ([]response.Chat, error)
	SaveChat(ctx context.Context, user1ID, user2ID uint) (chatID uint, err error)
	FindAllMessagesOfUserForAChat(ctx context.Context, chatID, userID uint, pagination request.Pagination) ([]response.Message, error)

	SaveMessage(ctx context.Context, message domain.Message) (receiverID uint, err error)
}
