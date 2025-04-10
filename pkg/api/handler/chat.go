package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joejosephvarghese/message/server/pkg/api/handler/interfaces"
	"github.com/joejosephvarghese/message/server/pkg/api/handler/request"
	"github.com/joejosephvarghese/message/server/pkg/api/handler/response"
	"github.com/joejosephvarghese/message/server/pkg/domain"
	producerinterface "github.com/joejosephvarghese/message/server/pkg/kafka/producerInterface"
	usecase "github.com/joejosephvarghese/message/server/pkg/usecase/interfaces"
)

type chatHandler struct {
	usecase  usecase.ChatUseCase
	producer producerinterface.ProdInterInterface
}

func NewChatHandler(usecase usecase.ChatUseCase, producer producerinterface.ProdInterInterface) interfaces.ChatHandler {
	return &chatHandler{
		usecase:  usecase,
		producer: producer,
	}
}

// GetRecentChats godoc
// @Summary Get user chats (User)
// @Description API for user to get all recent chats of user with others
// @Security ApiKeyAuth
// @Id GetRecentChats
// @Tags Users Chats
// @Param page_number query int false "Page Number"
// @Param count query int false "Count"
// @Router /chats [get]
// @Success 200 {object} response.Response{data=[]response.Chat} "Successfully retrieved recent chats of user"
// @Success 204 {object} response.Response{} "There is no chats recent chats for users"
// @Failure 500 {object} response.Response{} "Failed to retrieved recent chats of user"
func (c *chatHandler) GetRecentChats(ctx *gin.Context) {

	userID := request.GetUserIdFromContext(ctx)
	pagination := request.GetPagination(ctx)

	chats, err := c.usecase.FindAllRecentChatsOfUser(ctx, userID, pagination)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieved recent chats of user", err, nil)
		return
	}

	if len(chats) == 0 {
		response.SuccessResponse(ctx, http.StatusNoContent, "There is no chats recent chats for users", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved recent chats of user", chats)
}

// SaveChat godoc
// @Summary Save New chat (User)
// @Description API for user to create a new chat with another user. If a chat already exists, it returns the existing chat ID.
// @Security ApiKeyAuth
// @Id SaveChat
// @Tags Users Chats
// @Accept json
// @Produce json
// @Param input body request.Chat true "Input fields"
// @Router /chats [post]
// @Success 200 {object} response.Response{data=uint} "Successfully chat saved"
// @Failure 500 {object} response.Response{} "Failed to save chat for user"
func (c *chatHandler) SaveChat(ctx *gin.Context) {

	var body request.Chat

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}
	userID := request.GetUserIdFromContext(ctx)
	fmt.Println("user id:", userID)
	chatID, err := c.usecase.SaveChat(ctx, userID, body.OtherUserID)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to chat for user", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "Successfully chat saved", chatID)
}

// GetAllMessages godoc
// @Summary Get messages (User)
// @Description API for user to get all messages in a specific chat
// @Security ApiKeyAuth
// @Id GetAllMessages
// @Tags Users Message
// @Param chat_id path int true "Chat ID"
// @Param page_number query int false "Page Number"
// @Param count query int false "Count"
// @Router /chats/{chat_id}/messages [get]
// @Success 200 {object} response.Response{data=[]response.Message} "Successfully retrieved message for the chat"
// @Success 204 {object} response.Response{} "There is no message between users"
// @Failure 500 {object} response.Response{} "Failed to retrieve message for this chat"
func (c *chatHandler) GetAllMessages(ctx *gin.Context) {

	chatID, err := request.GetParamAsUint(ctx, "chat_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}
	pagination := request.GetPagination(ctx)
	userID := request.GetUserIdFromContext(ctx)

	messages, err := c.usecase.FindAllMessagesOfUserForAChat(ctx, chatID, userID, pagination)
	fmt.Print(messages)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve message for this chat", err, nil)
	}

	if len(messages) == 0 {
		response.ErrorResponse(ctx, http.StatusNoContent, "There is no message between users", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved message for the chat", messages)
}

// SaveMessage godoc
// @Summary Save message (User)
// @Description API for user to save a new message
// @Security ApiKeyAuth
// @Id SaveMessage
// @Tags Users Message
// @Param chat_id path int true "Chat ID"
// @Param input body request.Message true "Message field"
// @Router /chats/{chat_id}/messages [post]
// @Success 200 {object} response.Response{data=uint} "Successfully message saved"
// @Failure 500 {object} response.Response{} "Failed to save message"
func (c *chatHandler) SaveMessage(ctx *gin.Context) {

	chatID, err := request.GetParamAsUint(ctx, "chat_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	var body request.Message
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	message := domain.Message{
		ChatID:   chatID,
		SenderID: request.GetUserIdFromContext(ctx),
		Content:  body.Content,
	}
	// 🔥 Marshal the struct to JSON bytes
	msgBytes, err := json.Marshal(message)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to marshal message", err, nil)
		return
	}

	key := "message"
	if err := c.producer.Send(ctx, key, msgBytes); err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to send Kafka message", err, nil)
		return
	}
	// go c.goroute(ctx, message)
	// _, err = c.usecase.SaveMessage(ctx, message)
	// if err != nil {
	// 	response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to save message", err, nil)
	// 	return
	// }

	// // send the message to the receiver
	// received := c.socketService.SendMessage(receiverID, service.Message{
	// 	ChatID:  chatID,
	// 	Content: body.Content,
	// })

	// if received {
	// 	response.SuccessResponse(ctx, http.StatusCreated, "Successfully message saved and received on other side", received)
	// 	return
	// }
	response.SuccessResponse(ctx, http.StatusCreated, "Successfully message saved", nil)
}

// func (c *chatHandler) goroute(ctx *gin.Context, mesage domain.Message) {
// 	c.usecase.SaveMessage(ctx, mesage)
// }
