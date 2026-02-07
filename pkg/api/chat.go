package api

import "time"

// ChatMessageResponse represents the API response for a chat message
type ChatMessageResponse struct {
	ID         int                    `json:"id"`
	Role       string                 `json:"role"`
	ChunkIndex int                    `json:"chunk_index"`
	Content    string                 `json:"content"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CTime      time.Time              `json:"ctime"`
}

// ChatMessagesListResponse represents the API response for listing chat messages
type ChatMessagesListResponse struct {
	Messages []ChatMessageResponse `json:"messages"`
	Total    int                    `json:"total"`
}
