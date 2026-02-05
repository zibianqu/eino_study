package service

import (
	"fmt"

	"github.com/zibianqu/eino_study/internal/app/repository"
	"github.com/zibianqu/eino_study/internal/config"
	"github.com/zibianqu/eino_study/internal/eino/chatmodel"
	"github.com/zibianqu/eino_study/internal/eino/embedding"
	"github.com/zibianqu/eino_study/internal/eino/graph"
	"github.com/zibianqu/eino_study/internal/eino/loader"
	"github.com/zibianqu/eino_study/internal/eino/retriever"
	"github.com/zibianqu/eino_study/internal/eino/splitter"
)

// ServiceContainer holds all services and their dependencies
type ServiceContainer struct {
	DocumentService DocumentService
	RAGService      RAGService
}

// InitServices initializes all services with their dependencies
func InitServices(
	cfg *config.Config,
	docRepo repository.DocumentRepository,
	chunkRepo repository.ChunkRepository,
	entityRepo repository.EntityRepository,
) (*ServiceContainer, error) {
	// Initialize Eino components
	embeddingClient, err := embedding.NewEmbeddingClient(&cfg.Eino.Embedding)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding client: %w", err)
	}

	chatModelClient, err := chatmodel.NewChatModelClient(&cfg.Eino.LLM)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat model client: %w", err)
	}

	// Initialize loader and splitter
	loaderFactory := loader.NewLoaderFactory()
	textSplitter := splitter.NewTextSplitter(
		cfg.Eino.Splitter.ChunkSize,
		cfg.Eino.Splitter.ChunkOverlap,
	)

	// Initialize document processor
	docProcessor := graph.NewDocumentProcessor(
		loaderFactory,
		textSplitter,
		embeddingClient,
		chunkRepo,
	)

	// Initialize retriever
	vectorRetriever := retriever.NewVectorRetriever(
		chunkRepo,
		embeddingClient,
		cfg.Eino.Retriever.TopK,
		cfg.Eino.Retriever.SimilarityThreshold,
	)

	// Initialize RAG chain
	ragChain := graph.NewRAGChain(
		vectorRetriever,
		chatModelClient,
	)

	// Initialize services
	documentService := NewDocumentService(
		docRepo,
		chunkRepo,
		entityRepo,
		docProcessor,
	)

	ragService := NewRAGService(
		ragChain,
		docRepo,
	)

	return &ServiceContainer{
		DocumentService: documentService,
		RAGService:      ragService,
	}, nil
}