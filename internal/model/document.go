package model

import (
	"time"
)

// Document represents the documents table
type Document struct {
	DocID           string    `gorm:"column:doc_id;primaryKey;type:varchar(32)" json:"doc_id"`
	DocName         string    `gorm:"column:doc_name;type:varchar(255);not null" json:"doc_name"`
	DocHash         string    `gorm:"column:doc_hash;type:varchar(32);not null" json:"doc_hash"`
	FilePath        string    `gorm:"column:file_path;type:text;not null;unique" json:"file_path"`
	FileType        string    `gorm:"column:file_type;type:varchar(50);not null" json:"file_type"`
	SyncRagState    int       `gorm:"column:sync_rag_state;default:0" json:"sync_rag_state"`
	SyncEntityState int       `gorm:"column:sync_enity_state;default:0" json:"sync_entity_state"`
	CTime           time.Time `gorm:"column:ctime;default:CURRENT_TIMESTAMP" json:"ctime"`
}

// TableName specifies the table name
func (Document) TableName() string {
	return "documents"
}

// DocumentChunk represents the document_chunks table
type DocumentChunk struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DocID      string    `gorm:"column:doc_id;type:varchar(32);not null" json:"doc_id"`
	ChunkIndex int       `gorm:"column:chunk_index;not null" json:"chunk_index"`
	Content    string    `gorm:"column:content;type:text;not null" json:"content"`
	Embedding  string    `gorm:"column:embedding;type:vector(1536)" json:"-"`
	Metadata   string    `gorm:"column:metadata;type:jsonb" json:"metadata"`
	CTime      time.Time `gorm:"column:ctime;default:CURRENT_TIMESTAMP" json:"ctime"`
}

// TableName specifies the table name
func (DocumentChunk) TableName() string {
	return "document_chunks"
}

// Entity represents the entities table
type Entity struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DocID       string    `gorm:"column:doc_id;type:varchar(32);not null" json:"doc_id"`
	EntityType  string    `gorm:"column:entity_type;type:varchar(50);not null" json:"entity_type"`
	EntityName  string    `gorm:"column:entity_name;type:varchar(255);not null" json:"entity_name"`
	EntityValue string    `gorm:"column:entity_value;type:text" json:"entity_value"`
	Metadata    string    `gorm:"column:metadata;type:jsonb" json:"metadata"`
	CTime       time.Time `gorm:"column:ctime;default:CURRENT_TIMESTAMP" json:"ctime"`
}

// TableName specifies the table name
func (Entity) TableName() string {
	return "entities"
}