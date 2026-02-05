package graph

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/schema"
	"github.com/zibianqu/eino_study/internal/eino/chatmodel"
	"github.com/zibianqu/eino_study/internal/eino/retriever"
)

// RAGChain orchestrates the RAG workflow
type RAGChain struct {
	retriever *retriever.VectorRetriever
	chatModel *chatmodel.ChatModelClient
}

// NewRAGChain creates a new RAG chain
func NewRAGChain(
	retriever *retriever.VectorRetriever,
	chatModel *chatmodel.ChatModelClient,
) *RAGChain {
	return &RAGChain{
		retriever: retriever,
		chatModel: chatModel,
	}
}

// RAGResponse represents the response from RAG chain
type RAGResponse struct {
	Answer  string
	Sources []*schema.Document
	Usage   *UsageInfo
}

// UsageInfo represents token usage information
type UsageInfo struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// Run executes the RAG workflow
func (c *RAGChain) Run(ctx context.Context, query string) (*RAGResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("query is empty")
	}

	// Step 1: Retrieve relevant documents
	docs, err := c.retriever.Retrieve(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("retrieval failed: %w", err)
	}

	if len(docs) == 0 {
		return &RAGResponse{
			Answer:  "抱歉，我没有找到相关的文档来回答您的问题。",
			Sources: []*schema.Document{},
		}, nil
	}

	// Step 2: Build context from retrieved documents
	context := c.buildContext(docs)

	// Step 3: Build prompt
	prompt := c.buildPrompt(query, context)

	// Step 4: Generate answer using LLM
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: "你是一个专业的知识库助手。请根据提供的上下文信息回答用户的问题。如果上下文中没有相关信息，请明确说明。",
		},
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	response, err := c.chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	// Extract usage info if available
	usage := &UsageInfo{}
	if response.ResponseMeta != nil {
		if promptTokens, ok := response.ResponseMeta["prompt_tokens"].(int); ok {
			usage.PromptTokens = promptTokens
		}
		if completionTokens, ok := response.ResponseMeta["completion_tokens"].(int); ok {
			usage.CompletionTokens = completionTokens
		}
		if totalTokens, ok := response.ResponseMeta["total_tokens"].(int); ok {
			usage.TotalTokens = totalTokens
		}
	}

	return &RAGResponse{
		Answer:  response.Content,
		Sources: docs,
		Usage:   usage,
	}, nil
}

// buildContext builds context string from documents
func (c *RAGChain) buildContext(docs []*schema.Document) string {
	var builder strings.Builder

	for i, doc := range docs {
		builder.WriteString(fmt.Sprintf("[文档 %d]\n", i+1))
		builder.WriteString(doc.Content)
		builder.WriteString("\n\n")
	}

	return builder.String()
}

// buildPrompt builds the final prompt with context and query
func (c *RAGChain) buildPrompt(query, context string) string {
	return fmt.Sprintf(`基于以下上下文信息回答问题。请确保答案准确、完整，并尽可能引用上下文中的具体内容。

上下文信息：
%s

问题：%s

回答：`, context, query)
}