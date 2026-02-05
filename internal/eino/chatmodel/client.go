package chatmodel

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/zibianqu/eino_study/internal/config"
)

// ChatModelClient wraps Eino chat model component
type ChatModelClient struct {
	model       model.ChatModel
	temperature float64
	maxTokens   int
}

// NewChatModelClient creates a new chat model client
func NewChatModelClient(cfg *config.LLMConfig) (*ChatModelClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("LLM config is nil")
	}

	var chatModel model.ChatModel
	var err error

	switch cfg.Provider {
	case "openai":
		chatModel, err = openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
			APIKey:  cfg.APIKey,
			BaseURL: cfg.BaseURL,
			Model:   cfg.Model,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create OpenAI chat model: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", cfg.Provider)
	}

	return &ChatModelClient{
		model:       chatModel,
		temperature: cfg.Temperature,
		maxTokens:   cfg.MaxTokens,
	}, nil
}

// Generate generates a response from the chat model
func (c *ChatModelClient) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("messages is empty")
	}

	// Add default options
	options := []model.Option{
		model.WithTemperature(c.temperature),
		model.WithMaxTokens(c.maxTokens),
	}
	options = append(options, opts...)

	response, err := c.model.Generate(ctx, messages, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	return response, nil
}

// GenerateStream generates a streaming response from the chat model
func (c *ChatModelClient) GenerateStream(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("messages is empty")
	}

	// Add default options
	options := []model.Option{
		model.WithTemperature(c.temperature),
		model.WithMaxTokens(c.maxTokens),
	}
	options = append(options, opts...)

	stream, err := c.model.Stream(ctx, messages, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to generate streaming response: %w", err)
	}

	return stream, nil
}