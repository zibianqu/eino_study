package embedding

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/zibianqu/eino_study/internal/config"
)

// EmbeddingClient wraps Eino embedding component
type EmbeddingClient struct {
	embedder embedding.Embedder
	dimension int
}

// NewEmbeddingClient creates a new embedding client
func NewEmbeddingClient(cfg *config.EmbeddingConfig) (*EmbeddingClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("embedding config is nil")
	}

	var embedder embedding.Embedder
	var err error

	switch cfg.Provider {
	case "openai":
		embedder, err = openai.NewEmbedder(context.Background(), &openai.EmbedderConfig{
			APIKey: cfg.APIKey,
			Model:  cfg.Model,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create OpenAI embedder: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported embedding provider: %s", cfg.Provider)
	}

	return &EmbeddingClient{
		embedder:  embedder,
		dimension: cfg.Dimension,
	}, nil
}

// EmbedText generates embedding for a single text
func (c *EmbeddingClient) EmbedText(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text is empty")
	}

	vectors, err := c.embedder.EmbedStrings(ctx, []string{text})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(vectors) == 0 {
		return nil, fmt.Errorf("no embedding generated")
	}

	return vectors[0], nil
}

// EmbedTexts generates embeddings for multiple texts
func (c *EmbeddingClient) EmbedTexts(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("texts is empty")
	}

	vectors, err := c.embedder.EmbedStrings(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	return vectors, nil
}

// GetDimension returns the embedding dimension
func (c *EmbeddingClient) GetDimension() int {
	return c.dimension
}