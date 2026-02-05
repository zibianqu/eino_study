package loader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudwego/eino/schema"
)

// DocumentLoader interface defines methods for loading documents
type DocumentLoader interface {
	Load(ctx context.Context, filePath string) ([]*schema.Document, error)
}

// BaseLoader provides common functionality for all loaders
type BaseLoader struct{}

// ReadFile reads file content from given path
func (l *BaseLoader) ReadFile(filePath string) ([]byte, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return content, nil
}

// CreateDocument creates a schema.Document with metadata
func (l *BaseLoader) CreateDocument(content, filePath string) *schema.Document {
	return &schema.Document{
		Content: content,
		MetaData: map[string]any{
			"source":    filePath,
			"file_name": filepath.Base(filePath),
			"file_type": filepath.Ext(filePath),
		},
	}
}