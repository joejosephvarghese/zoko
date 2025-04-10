//go:build wireinject
// +build wireinject

package di

import (
	http "github.com/joejosephvarghese/message/server/pkg/api"
	"github.com/joejosephvarghese/message/server/pkg/api/handler"
	"github.com/joejosephvarghese/message/server/pkg/api/middleware"
	socket "github.com/joejosephvarghese/message/server/pkg/api/service"
	"github.com/joejosephvarghese/message/server/pkg/config"
	"github.com/joejosephvarghese/message/server/pkg/db"
	"github.com/joejosephvarghese/message/server/pkg/kafka"
	"github.com/joejosephvarghese/message/server/pkg/repository"
	"github.com/joejosephvarghese/message/server/pkg/service/google"
	"github.com/joejosephvarghese/message/server/pkg/service/token"
	"github.com/joejosephvarghese/message/server/pkg/usecase"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.Server, error) {

	wire.Build(
		db.ConnectDatabase,
		kafka.NewProducer,
		token.NewTokenService,
		google.NewGoogleAuth,
		middleware.NewMiddleware,

		repository.NewUserRepository,
		repository.NewAuthRepository,
		repository.NewChatRepository,

		usecase.NewAuthUseCase,
		usecase.NewChatUseCase,
		usecase.NewUserUseCase,

		socket.NewWebSocketService,
		handler.NewAuthHandler,
		handler.NewChatHandler,
		handler.NewUserHandler,

		http.NewServerHTTP,
	)

	return &http.Server{}, nil
}
