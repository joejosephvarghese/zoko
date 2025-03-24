package interfaces

import (
	"context"

	"github.com/joejosephvarghese/message/server/pkg/api/handler/request"
	"github.com/joejosephvarghese/message/server/pkg/api/handler/response"
	"github.com/joejosephvarghese/message/server/pkg/domain"
)

type UserRepository interface {
	FindUserByUserID(ctx context.Context, userID uint) (user domain.User, err error)
	FindUserByEmail(ctx context.Context, email string) (user domain.User, err error)
	FindUserByUserName(ctx context.Context, userName string) (user domain.User, err error)
	FindUserByUserNameEmailNotID(ctx context.Context, user domain.User) (domain.User, error)
	SaveUser(ctx context.Context, user domain.User) (userID uint, err error)
	FindAllUsers(ctx context.Context, pagination request.Pagination) ([]response.User, error)
}
