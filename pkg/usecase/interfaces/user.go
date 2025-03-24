package interfaces

import (
	"context"

	"github.com/joejosephvarghese/message/server/pkg/api/handler/request"
	"github.com/joejosephvarghese/message/server/pkg/api/handler/response"
)

type UserUseCase interface {
	FindAllUsers(ctx context.Context, pagination request.Pagination) ([]response.User, error)
}
