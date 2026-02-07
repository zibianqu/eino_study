package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zibianqu/eino_study/internal/app/service"
)

type ChatHandler struct {
	service service.ChatService
}

func NewChatHandler(service service.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

// CreateMessageRequest represents the request body for creating a chat message
type CreateMessageRequest struct {
	Role     string                 `json:"role" binding:"required"`
	Content  string                 `json:"content" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

// CreateMessage handles POST /api/v1/chat/messages
func (h *ChatHandler) CreateMessage(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	chunk, err := h.service.CreateMessage(req.Role, req.Content, req.Metadata)
	if err != nil {
		InternalError(c, "Failed to create message: "+err.Error())
		return
	}

	Success(c, chunk)
}

// GetMessage handles GET /api/v1/chat/messages/:id
func (h *ChatHandler) GetMessage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(c, "Invalid message ID")
		return
	}

	chunk, err := h.service.GetMessage(id)
	if err != nil {
		NotFound(c, "Message not found")
		return
	}

	Success(c, chunk)
}

// ListMessages handles GET /api/v1/chat/messages
func (h *ChatHandler) ListMessages(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	role := c.Query("role")

	var chunks interface{}
	var err error

	if role != "" {
		chunks, err = h.service.GetMessagesByRole(role, limit, offset)
	} else {
		chunks, err = h.service.ListMessages(limit, offset)
	}

	if err != nil {
		InternalError(c, "Failed to list messages: "+err.Error())
		return
	}

	Success(c, chunks)
}

// DeleteMessage handles DELETE /api/v1/chat/messages/:id
func (h *ChatHandler) DeleteMessage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(c, "Invalid message ID")
		return
	}

	if err := h.service.DeleteMessage(id); err != nil {
		InternalError(c, "Failed to delete message: "+err.Error())
		return
	}

	Success(c, gin.H{"message": "Message deleted successfully"})
}
