package loader

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
)

// TextLoader loads plain text files
type TextLoader struct {
	BaseLoader
}

// NewTextLoader creates a new text loader
func NewTextLoader() *TextLoader {
	return &TextLoader{}
}

// Load reads a text file and returns a document
func (l *TextLoader) Load(ctx context.Context, filePath string) ([]*schema.Document, error) {
	content, err := l.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("file is empty: %s", filePath)
	}

	doc := l.CreateDocument(string(content), filePath)
	return []*schema.Document{doc}, nil
}