package service

import (
	"encoding/json"
	"fmt"

	"github.com/zibianqu/eino_study/internal/app/repository"
	"github.com/zibianqu/eino_study/internal/model"
)

// ChatService defines the interface for chat message operations
type ChatService interface {
	CreateMessage(role, content string, metadata map[string]interface{}) (*model.ChatChunk, error)
	GetMessage(id int) (*model.ChatChunk, error)
	ListMessages(limit, offset int) ([]*model.ChatChunk, error)
	GetMessagesByRole(role string, limit, offset int) ([]*model.ChatChunk, error)
	DeleteMessage(id int) error
	SearchSimilarMessages(query string, topK int) ([]*model.ChatChunk, error)
}

type chatService struct {
	chatRepo repository.ChatRepository
	// embeddingService can be added here for generating embeddings
}

// NewChatService creates a new ChatService instance
func NewChatService(chatRepo repository.ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

// CreateMessage creates a new chat message
func (s *chatService) CreateMessage(role, content string, metadata map[string]interface{}) (*model.ChatChunk, error) {
	if role == "" {
		return nil, fmt.Errorf("role is required")
	}
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}

	// Validate role
	validRoles := map[string]bool{"user": true, "assistant": true, "system": true}
	if !validRoles[role] {
		return nil, fmt.Errorf("invalid role: must be user, assistant, or system")
	}

	// Convert metadata to JSON string
	var metadataJSON string
	if metadata != nil {
		data, err := json.Marshal(metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataJSON = string(data)
	}

	chunk := &model.ChatChunk{
		Role:     role,
		Content:  content,
		Metadata: metadataJSON,
		// ChunkIndex will be auto-managed or set externally
		// Embedding can be generated using embedding service
	}

	if err := s.chatRepo.Create(chunk); err != nil {
		return nil, fmt.Errorf("failed to create chat message: %w", err)
	}

	return chunk, nil
}

// GetMessage retrieves a chat message by ID
func (s *chatService) GetMessage(id int) (*model.ChatChunk, error) {
	return s.chatRepo.GetByID(id)
}

// ListMessages retrieves paginated chat messages
func (s *chatService) ListMessages(limit, offset int) ([]*model.ChatChunk, error) {
	if limit <= 0 {
		limit = 20 // default limit
	}
	if offset < 0 {
		offset = 0
	}
	return s.chatRepo.List(limit, offset)
}

// GetMessagesByRole retrieves chat messages filtered by role
func (s *chatService) GetMessagesByRole(role string, limit, offset int) ([]*model.ChatChunk, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.chatRepo.GetByRole(role, limit, offset)
}

// DeleteMessage deletes a chat message by ID
func (s *chatService) DeleteMessage(id int) error {
	return s.chatRepo.Delete(id)
}

// SearchSimilarMessages searches for similar messages using vector similarity
// Note: This requires embeddings to be generated and stored
func (s *chatService) SearchSimilarMessages(query string, topK int) ([]*model.ChatChunk, error) {
	// TODO: Generate embedding for the query using embedding service
	// For now, this is a placeholder
	return nil, fmt.Errorf("semantic search not implemented: embedding service required")
}
