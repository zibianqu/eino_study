package model

import (
	"time"
)

// ChatChunk represents the chat_chunk table for storing chat message records
// This table stores individual chat messages with optional vector embeddings for semantic search
type ChatChunk struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Role       string    `gorm:"column:role;type:varchar(20);not null" json:"role"`         // Role of the message sender (e.g., "user", "assistant", "system")
	ChunkIndex int       `gorm:"column:chunk_index;not null" json:"chunk_index"`             // Index for ordering chunks within a conversation
	Content    string    `gorm:"column:content;type:text;not null" json:"content"`           // The actual message content
	Embedding  string    `gorm:"column:embedding;type:vector(1536)" json:"-"`                // Vector embedding for semantic search (1536 dimensions for OpenAI embeddings)
	Metadata   string    `gorm:"column:metadata;type:jsonb" json:"metadata"`                 // Additional metadata in JSON format (e.g., session_id, user_id, etc.)
	CTime      time.Time `gorm:"column:ctime;default:CURRENT_TIMESTAMP" json:"ctime"`        // Creation timestamp
}

// TableName specifies the table name
func (ChatChunk) TableName() string {
	return "chat_chunk"
}
