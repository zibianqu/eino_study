package retriever

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/zibianqu/eino_study/internal/app/repository"
	"github.com/zibianqu/eino_study/internal/eino/embedding"
)

// VectorRetriever retrieves relevant documents using vector similarity
type VectorRetriever struct {
	chunkRepo repository.ChunkRepository
	embedding *embedding.EmbeddingClient
	topK      int
	threshold float64
}

// NewVectorRetriever creates a new vector retriever
func NewVectorRetriever(
	chunkRepo repository.ChunkRepository,
	embedding *embedding.EmbeddingClient,
	topK int,
	threshold float64,
) *VectorRetriever {
	if topK <= 0 {
		topK = 5
	}
	if threshold <= 0 {
		threshold = 0.7
	}

	return &VectorRetriever{
		chunkRepo: chunkRepo,
		embedding: embedding,
		topK:      topK,
		threshold: threshold,
	}
}

// Retrieve retrieves relevant documents for a query
func (r *VectorRetriever) Retrieve(ctx context.Context, query string) ([]*schema.Document, error) {
	if query == "" {
		return nil, fmt.Errorf("query is empty")
	}

	// Generate embedding for query
	queryVector, err := r.embedding.EmbedText(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Convert vector to string format for pgvector
	vectorStr := vectorToString(queryVector)

	// Search similar chunks
	chunks, err := r.chunkRepo.SearchSimilar(vectorStr, r.topK, r.threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to search similar chunks: %w", err)
	}

	// Convert chunks to documents
	docs := make([]*schema.Document, 0, len(chunks))
	for _, chunk := range chunks {
		doc := &schema.Document{
			Content: chunk.Content,
			MetaData: map[string]any{
				"doc_id":      chunk.DocID,
				"chunk_index": chunk.ChunkIndex,
				"chunk_id":    chunk.ID,
			},
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

// vectorToString converts float32 slice to string format for pgvector
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