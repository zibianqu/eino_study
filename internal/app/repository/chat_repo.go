package repository

import (
	"github.com/zibianqu/eino_study/internal/model"
	"gorm.io/gorm"
)

// ChatRepository defines the interface for chat chunk operations
type ChatRepository interface {
	Create(chunk *model.ChatChunk) error
	BatchCreate(chunks []*model.ChatChunk) error
	GetByID(id int) (*model.ChatChunk, error)
	List(limit, offset int) ([]*model.ChatChunk, error)
	GetByRole(role string, limit, offset int) ([]*model.ChatChunk, error)
	Delete(id int) error
	SearchSimilar(embedding string, topK int, threshold float64) ([]*model.ChatChunk, error)
}

type chatRepository struct {
	db *gorm.DB
}

// NewChatRepository creates a new ChatRepository instance
func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

// Create inserts a new chat chunk
func (r *chatRepository) Create(chunk *model.ChatChunk) error {
	return r.db.Create(chunk).Error
}

// BatchCreate inserts multiple chat chunks in batches
func (r *chatRepository) BatchCreate(chunks []*model.ChatChunk) error {
	return r.db.CreateInBatches(chunks, 100).Error
}

// GetByID retrieves a chat chunk by its ID
func (r *chatRepository) GetByID(id int) (*model.ChatChunk, error) {
	var chunk model.ChatChunk
	err := r.db.First(&chunk, id).Error
	return &chunk, err
}

// List retrieves chat chunks with pagination, ordered by chunk_index
func (r *chatRepository) List(limit, offset int) ([]*model.ChatChunk, error) {
	var chunks []*model.ChatChunk
	err := r.db.Order("chunk_index ASC").Limit(limit).Offset(offset).Find(&chunks).Error
	return chunks, err
}

// GetByRole retrieves chat chunks filtered by role with pagination
func (r *chatRepository) GetByRole(role string, limit, offset int) ([]*model.ChatChunk, error) {
	var chunks []*model.ChatChunk
	err := r.db.Where("role = ?", role).Order("chunk_index ASC").Limit(limit).Offset(offset).Find(&chunks).Error
	return chunks, err
}

// Delete removes a chat chunk by ID
func (r *chatRepository) Delete(id int) error {
	return r.db.Delete(&model.ChatChunk{}, id).Error
}

// SearchSimilar performs vector similarity search using pgvector
func (r *chatRepository) SearchSimilar(embedding string, topK int, threshold float64) ([]*model.ChatChunk, error) {
	var chunks []*model.ChatChunk
	// Using pgvector cosine similarity search
	query := `
		SELECT id, role, chunk_index, content, metadata, ctime,
		       1 - (embedding <=> ?::vector) as similarity
		FROM chat_chunk
		WHERE 1 - (embedding <=> ?::vector) > ?
		ORDER BY embedding <=> ?::vector
		LIMIT ?
	`
	err := r.db.Raw(query, embedding, embedding, threshold, embedding, topK).Scan(&chunks).Error
	return chunks, err
}
