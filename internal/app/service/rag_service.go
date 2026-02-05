package service

import (
	"fmt"

	"github.com/zibianqu/eino_study/internal/app/repository"
	"github.com/zibianqu/eino_study/pkg/api"
)

type RAGService interface {
	Query(query string, topK int) (*api.QueryResponse, error)
}

type ragService struct {
	chunkRepo repository.ChunkRepository
	docRepo   repository.DocumentRepository
}

func NewRAGService(
	chunkRepo repository.ChunkRepository,
	docRepo repository.DocumentRepository,
) RAGService {
	return &ragService{
		chunkRepo: chunkRepo,
		docRepo:   docRepo,
	}
}

func (s *ragService) Query(query string, topK int) (*api.QueryResponse, error) {
	// TODO: Implement RAG query with Eino
	// 1. Generate query embedding
	// 2. Search similar chunks
	// 3. Build prompt with context
	// 4. Call LLM
	// 5. Return response with sources
	return nil, fmt.Errorf("not implemented yet")
}