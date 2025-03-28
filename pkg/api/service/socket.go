package socket

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joejosephvarghese/message/server/pkg/service/token"
)

type WebSocketService interface {
	ServeWebSocket(ctx *gin.Context)
}

type webSocketService struct {
	upgrader     websocket.Upgrader
	connections  map[uint]*websocket.Conn
	mu           sync.Mutex
	tokenService token.TokenService
}

type Message struct {
	SenderID   uint   `json:"sender_id"` // will remove if its not necessary
	ChatID     uint   `json:"chat_id"`
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
}

func NewWebSocketService(tokenService token.TokenService) WebSocketService {

	return &webSocketService{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		connections:  make(map[uint]*websocket.Conn),
		mu:           sync.Mutex{},
		tokenService: tokenService,
	}
}

// ServeWebSocket godoc
// @Summary Serve WebSocket Connection (User)
// @Description API for users to establish a WebSocket connection
// @Security ApiKeyAuth
// @ID ServeWebSocket
// @Tags Users Socket
// @Param token query string true "JWT Token"
// @Router /ws [get]
func (c *webSocketService) ServeWebSocket(ctx *gin.Context) {
	socketConn, err := c.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	fmt.Println("kjj :", err)
	if err != nil {
		log.Println("failed to upgrade connection")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	closeHand := socketConn.CloseHandler()

	userID, err := c.verifyConnection(ctx, socketConn)
	if err != nil {
		log.Println("failed to verify token", err)
		closeHand(websocket.ClosePolicyViolation, err.Error())
		return
	}

	c.mu.Lock()
	c.connections[userID] = socketConn
	c.mu.Unlock()

	{ // wait until the connection close
		fmt.Println("successfully connected for user_id : ", userID)
		go c.readMessages(ctx, socketConn) // read messages
		<-ctx.Done()
		log.Println("connection closed from request")
	}

	c.mu.Lock()
	delete(c.connections, userID)
	c.mu.Unlock()
}

type TokenRequest struct {
	Token string `json:"token"`
}

func (c *webSocketService) verifyConnection(ctx context.Context, sc *websocket.Conn) (userID uint, err error) {

	tokenChan := make(chan TokenRequest)
	errChan := make(chan error)

	// wait for the token to send within 5 second of connection established
	go func() {
		var body TokenRequest
		err := sc.ReadJSON(&body)
		if err != nil {
			errChan <- err
		}
		tokenChan <- body
	}()

	select {
	case err := <-errChan:
		return 0, err
	case body := <-tokenChan:
		tokenRes, err := c.tokenService.VerifyToken(token.VerifyTokenRequest{TokenString: body.Token, UsedFor: token.User})
		if err != nil {
			return 0, err
		}
		return tokenRes.UserID, nil
	case <-time.After(5 * time.Second):
		return 0, errors.New("time exceed for waiting token send")
	}

}

func (c *webSocketService) readMessages(ctx context.Context, sc *websocket.Conn) {

	messageChan := make(chan Message)
	go func() {
		for {
			var body Message
			if err := sc.ReadJSON(&body); err != nil {
				log.Println("failed to read message: ", err)
				return
			}
			messageChan <- body
		}

	}()

	// get each message and if the whole connection lost then return
	for {
		select {
		case message := <-messageChan:
			go c.sendMessage(message)
		case <-ctx.Done():
			return
		}
	}
}

func (c *webSocketService) sendMessage(message Message) (received bool, err error) {

	c.mu.Lock()
	defer c.mu.Unlock()
	// find other connection
	conn, ok := c.connections[message.ReceiverID]
	if !ok {
		return false, nil
	}

	// send to other connection
	if err := conn.WriteJSON(message); err != nil {
		log.Println("failed to write message: ", message)
		return false, err
	}
	return true, nil
}
