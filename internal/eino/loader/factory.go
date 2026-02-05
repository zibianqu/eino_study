package loader

import (
	"fmt"
	"path/filepath"
	"strings"
)

// LoaderFactory creates appropriate loader based on file type
type LoaderFactory struct{}

// NewLoaderFactory creates a new loader factory
func NewLoaderFactory() *LoaderFactory {
	return &LoaderFactory{}
}

// GetLoader returns appropriate loader for the given file path
func (f *LoaderFactory) GetLoader(filePath string) (DocumentLoader, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".txt":
		return NewTextLoader(), nil
	case ".md", ".markdown":
		return NewMarkdownLoader(), nil
	case ".pdf":
		return nil, fmt.Errorf("PDF loader not implemented yet")
	default:
		// Default to text loader for unknown types
		return NewTextLoader(), nil
	}
}