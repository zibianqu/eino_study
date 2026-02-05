package graph

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/zibianqu/eino_study/internal/app/repository"
	"github.com/zibianqu/eino_study/internal/eino/embedding"
	"github.com/zibianqu/eino_study/internal/eino/loader"
	"github.com/zibianqu/eino_study/internal/eino/splitter"
	"github.com/zibianqu/eino_study/internal/model"
)

// DocumentProcessor processes documents for RAG
type DocumentProcessor struct {
	loaderFactory *loader.LoaderFactory
	splitter      *splitter.TextSplitter
	embedding     *embedding.EmbeddingClient
	chunkRepo     repository.ChunkRepository
}

// NewDocumentProcessor creates a new document processor
func NewDocumentProcessor(
	loaderFactory *loader.LoaderFactory,
	splitter *splitter.TextSplitter,
	embedding *embedding.EmbeddingClient,
	chunkRepo repository.ChunkRepository,
) *DocumentProcessor {
	return &DocumentProcessor{
		loaderFactory: loaderFactory,
		splitter:      splitter,
		embedding:     embedding,
		chunkRepo:     chunkRepo,
	}
}

// Process loads, splits, embeds and stores a document
func (p *DocumentProcessor) Process(ctx context.Context, docID, filePath string) error {
	// Step 1: Load document
	loader, err := p.loaderFactory.GetLoader(filePath)
	if err != nil {
		return fmt.Errorf("failed to get loader: %w", err)
	}

	docs, err := loader.Load(ctx, filePath)
	if err != nil {
		return fmt.Errorf("failed to load document: %w", err)
	}

	if len(docs) == 0 {
		return fmt.Errorf("no documents loaded")
	}

	// Step 2: Split document into chunks
	var allChunks []*schema.Document
	for _, doc := range docs {
		chunks, err := p.splitter.Transform(ctx, doc)
		if err != nil {
			return fmt.Errorf("failed to split document: %w", err)
		}
		allChunks = append(allChunks, chunks...)
	}

	if len(allChunks) == 0 {
		return fmt.Errorf("no chunks generated")
	}

	// Step 3: Generate embeddings
	texts := make([]string, len(allChunks))
	for i, chunk := range allChunks {
		texts[i] = chunk.Content
	}

	vectors, err := p.embedding.EmbedTexts(ctx, texts)
	if err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Step 4: Store chunks with embeddings
	dbChunks := make([]*model.DocumentChunk, len(allChunks))
	for i, chunk := range allChunks {
		chunkIndex := 0
		if idx, ok := chunk.MetaData["chunk_index"].(int); ok {
			chunkIndex = idx
		}

		dbChunks[i] = &model.DocumentChunk{
			DocID:      docID,
			ChunkIndex: chunkIndex,
			Content:    chunk.Content,
			Embedding:  vectorToString(vectors[i]),
			Metadata:   "", // Could store chunk.MetaData as JSON if needed
		}
	}

	if err := p.chunkRepo.BatchCreate(dbChunks); err != nil {
		return fmt.Errorf("failed to store chunks: %w", err)
	}

	return nil
}

// vectorToString converts vector to pgvector string format
func vectorToString(vector []float32) string {
	if len(vector) == 0 {
		return "[]"
	}

	result := "["
	for i, v := range vector {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%f", v)
	}
	result += "]"
	return result
}