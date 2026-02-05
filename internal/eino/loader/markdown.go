package loader

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
)

// MarkdownLoader loads markdown files
type MarkdownLoader struct {
	BaseLoader
}

// NewMarkdownLoader creates a new markdown loader
func NewMarkdownLoader() *MarkdownLoader {
	return &MarkdownLoader{}
}

// Load reads a markdown file and returns a document
func (l *MarkdownLoader) Load(ctx context.Context, filePath string) ([]*schema.Document, error) {
	content, err := l.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("file is empty: %s", filePath)
	}

	doc := l.CreateDocument(string(content), filePath)
	doc.MetaData["format"] = "markdown"

	return []*schema.Document{doc}, nil
}