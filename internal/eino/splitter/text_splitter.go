package splitter

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/schema"
)

// TextSplitter splits documents into chunks
type TextSplitter struct {
	ChunkSize    int
	ChunkOverlap int
}

// NewTextSplitter creates a new text splitter
func NewTextSplitter(chunkSize, chunkOverlap int) *TextSplitter {
	if chunkSize <= 0 {
		chunkSize = 1000
	}
	if chunkOverlap < 0 {
		chunkOverlap = 0
	}
	if chunkOverlap >= chunkSize {
		chunkOverlap = chunkSize / 4
	}

	return &TextSplitter{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
	}
}

// Transform splits a document into multiple chunks
func (s *TextSplitter) Transform(ctx context.Context, doc *schema.Document) ([]*schema.Document, error) {
	if doc == nil || len(doc.Content) == 0 {
		return []*schema.Document{}, nil
	}

	chunks := s.splitText(doc.Content)
	result := make([]*schema.Document, 0, len(chunks))

	for i, chunk := range chunks {
		if len(strings.TrimSpace(chunk)) == 0 {
			continue
		}

		// Create metadata for chunk
		metadata := make(map[string]any)
		for k, v := range doc.MetaData {
			metadata[k] = v
		}
		metadata["chunk_index"] = i
		metadata["chunk_size"] = len(chunk)

		result = append(result, &schema.Document{
			Content:  chunk,
			MetaData: metadata,
		})
	}

	return result, nil
}

// splitText splits text into chunks with overlap
func (s *TextSplitter) splitText(text string) []string {
	if len(text) <= s.ChunkSize {
		return []string{text}
	}

	var chunks []string
	start := 0

	for start < len(text) {
		end := start + s.ChunkSize
		if end > len(text) {
			end = len(text)
		}

		// Try to break at sentence or word boundary
		if end < len(text) {
			// Look for sentence ending
			for i := end; i > start && i > end-100; i-- {
				if text[i] == '.' || text[i] == '!' || text[i] == '?' || text[i] == '\n' {
					end = i + 1
					break
				}
			}

			// If no sentence boundary, look for space
			if end == start+s.ChunkSize {
				for i := end; i > start && i > end-50; i-- {
					if text[i] == ' ' {
						end = i
						break
					}
				}
			}
		}

		chunks = append(chunks, text[start:end])

		// Move start position with overlap
		if end >= len(text) {
			break
		}
		start = end - s.ChunkOverlap
		if start < 0 {
			start = 0
		}
	}

	return chunks
}