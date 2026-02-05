package service

import (
	"context"
	"fmt"

	"github.com/zibianqu/eino_study/internal/app/repository"
	"github.com/zibianqu/eino_study/internal/eino/graph"
	"github.com/zibianqu/eino_study/pkg/api"
)

type RAGService interface {
	Query(query string, topK int) (*api.QueryResponse, error)
}

type ragService struct {
	chain   *graph.RAGChain
	docRepo repository.DocumentRepository
}

func NewRAGService(
	chain *graph.RAGChain,
	docRepo repository.DocumentRepository,
) RAGService {
	return &ragService{
		chain:   chain,
		docRepo: docRepo,
	}
}

func (s *ragService) Query(query string, topK int) (*api.QueryResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("query is empty")
	}

	// Execute RAG chain
	ctx := context.Background()
	result, err := s.chain.Run(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("RAG query failed: %w", err)
	}

	// Build response with sources
	sources := make([]api.SourceInfo, 0, len(result.Sources))
	for _, doc := range result.Sources {
		// Get document metadata
		docID := ""
		if id, ok := doc.MetaData["doc_id"].(string); ok {
			docID = id
		}

		// Try to get document name from database
		docName := ""
		if docID != "" {
			if dbDoc, err := s.docRepo.GetByID(docID); err == nil {
				docName = dbDoc.DocName
			}
		}

		sources = append(sources, api.SourceInfo{
			DocID:   docID,
			DocName: docName,
			Content: doc.Content,
		})
	}

	// Build usage info
	var usage *api.UsageInfo
	if result.Usage != nil {
		usage = &api.UsageInfo{
			PromptTokens:     result.Usage.PromptTokens,
			CompletionTokens: result.Usage.CompletionTokens,
			TotalTokens:      result.Usage.TotalTokens,
		}
	}

	return &api.QueryResponse{
		Answer:  result.Answer,
		Sources: sources,
		Usage:   usage,
	}, nil
}