package repository

import (
	"github.com/zibianqu/eino_study/internal/model"
	"gorm.io/gorm"
)

type ChunkRepository interface {
	Create(chunk *model.DocumentChunk) error
	BatchCreate(chunks []*model.DocumentChunk) error
	GetByDocID(docID string) ([]*model.DocumentChunk, error)
	DeleteByDocID(docID string) error
	SearchSimilar(embedding string, topK int, threshold float64) ([]*model.DocumentChunk, error)
}

type chunkRepository struct {
	db *gorm.DB
}

func NewChunkRepository(db *gorm.DB) ChunkRepository {
	return &chunkRepository{db: db}
}

func (r *chunkRepository) Create(chunk *model.DocumentChunk) error {
	return r.db.Create(chunk).Error
}

func (r *chunkRepository) BatchCreate(chunks []*model.DocumentChunk) error {
	return r.db.CreateInBatches(chunks, 100).Error
}

func (r *chunkRepository) GetByDocID(docID string) ([]*model.DocumentChunk, error) {
	var chunks []*model.DocumentChunk
	err := r.db.Where("doc_id = ?", docID).Order("chunk_index").Find(&chunks).Error
	return chunks, err
}

func (r *chunkRepository) DeleteByDocID(docID string) error {
	return r.db.Where("doc_id = ?", docID).Delete(&model.DocumentChunk{}).Error
}

func (r *chunkRepository) SearchSimilar(embedding string, topK int, threshold float64) ([]*model.DocumentChunk, error) {
	var chunks []*model.DocumentChunk
	// Using pgvector cosine similarity search
	query := `
		SELECT id, doc_id, chunk_index, content, metadata, ctime,
		       1 - (embedding <=> ?::vector) as similarity
		FROM document_chunks
		WHERE 1 - (embedding <=> ?::vector) > ?
		ORDER BY embedding <=> ?::vector
		LIMIT ?
	`
	err := r.db.Raw(query, embedding, embedding, threshold, embedding, topK).Scan(&chunks).Error
	return chunks, err
}