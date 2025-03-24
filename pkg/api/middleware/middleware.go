package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/joejosephvarghese/message/server/pkg/service/token"
)

type Middleware interface {
	AuthenticateUser() gin.HandlerFunc
	Cors() gin.HandlerFunc
}

type middleware struct {
	tokenService token.TokenService
}

func NewMiddleware(tokenService token.TokenService) Middleware {
	return &middleware{
		tokenService: tokenService,
	}
}
