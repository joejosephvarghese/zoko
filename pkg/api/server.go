package http

import (
	"log"

	"github.com/gin-gonic/gin"
	docs "github.com/joejosephvarghese/message/server/cmd/api/docs"
	"github.com/joejosephvarghese/message/server/pkg/api/handler/interfaces"
	"github.com/joejosephvarghese/message/server/pkg/api/middleware"
	"github.com/joejosephvarghese/message/server/pkg/api/routes"
	socket "github.com/joejosephvarghese/message/server/pkg/api/service"
	"github.com/joejosephvarghese/message/server/pkg/config"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine *gin.Engine
	port   string
}

func NewServerHTTP(cfg config.Config, authHandler interfaces.AuthHandler,
	middleware middleware.Middleware, userHandler interfaces.UserHandler,
	chatHandler interfaces.ChatHandler, socketService socket.WebSocketService) *Server {

	engine := gin.New()

	engine.Use(middleware.Cors())
	engine.Use(gin.Logger())

	docs.SwaggerInfo.BasePath = routes.BaseURL
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	routes.SetupRoutes(engine, authHandler, middleware, userHandler, chatHandler, socketService)

	return &Server{
		engine: engine,
		port:   cfg.Port,
	}
}

func (c *Server) Start() error {
	// Ensure a valid port
	port := c.port
	if port == "" || port == "5432" { // Prevent using the DB port
		port = "8080"
		log.Println("Defaulting to port 8080")
	}
	log.Println("Starting server on port:", port)
	return c.engine.Run(":" + port)
}
