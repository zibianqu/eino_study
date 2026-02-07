package model

import "time"

// GraphNode represents a base graph node
type GraphNode struct {
	ID         string                 `json:"id"`
	Label      string                 `json:"label"`
	Properties map[string]interface{} `json:"properties"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// DocumentNode represents a document node in Neo4j
type DocumentNode struct {
	ID       string    `json:"id"`
	DocName  string    `json:"doc_name"`
	DocHash  string    `json:"doc_hash"`
	FilePath string    `json:"file_path"`
	FileType string    `json:"file_type"`
	CTime    time.Time `json:"ctime"`
}

// EntityNode represents an entity node in Neo4j
type EntityNode struct {
	ID          string                 `json:"id"`
	EntityType  string                 `json:"entity_type"`
	EntityName  string                 `json:"entity_name"`
	EntityValue string                 `json:"entity_value"`
	Metadata    map[string]interface{} `json:"metadata"`
	CTime       time.Time              `json:"ctime"`
}

// ChatMessageNode represents a chat message node in Neo4j
type ChatMessageNode struct {
	ID         string                 `json:"id"`
	Role       string                 `json:"role"`
	Content    string                 `json:"content"`
	ChunkIndex int                    `json:"chunk_index"`
	Metadata   map[string]interface{} `json:"metadata"`
	CTime      time.Time              `json:"ctime"`
}

// Relationship represents a relationship between nodes
type Relationship struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	FromNodeID string                 `json:"from_node_id"`
	ToNodeID   string                 `json:"to_node_id"`
	Properties map[string]interface{} `json:"properties"`
	CreatedAt  time.Time              `json:"created_at"`
}

// RelationshipType defines common relationship types
type RelationshipType string

const (
	RelContains     RelationshipType = "CONTAINS"      // Document contains Entity
	RelReferences   RelationshipType = "REFERENCES"    // Document references Document
	RelSimilarTo    RelationshipType = "SIMILAR_TO"    // Document similar to Document
	RelMentionedIn  RelationshipType = "MENTIONED_IN"  // Entity mentioned in ChatMessage
	RelRelatedTo    RelationshipType = "RELATED_TO"    // Generic relation
	RelDerivedFrom  RelationshipType = "DERIVED_FROM"  // Entity derived from Document
	RelPartOf       RelationshipType = "PART_OF"       // Entity part of another Entity
)

// GraphQuery represents a graph query result
type GraphQuery struct {
	Nodes         []GraphNode    `json:"nodes"`
	Relationships []Relationship `json:"relationships"`
}

// PathQuery represents a path query result
type PathQuery struct {
	Paths [][]GraphNode `json:"paths"`
}
