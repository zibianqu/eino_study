package model

import (
	"time"
)

// GraphNode represents a generic node in Neo4j
type GraphNode struct {
	ID         string                 `json:"id"`
	Label      string                 `json:"label"`       // Node label/type
	Properties map[string]interface{} `json:"properties"`  // Flexible properties
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// GraphRelationship represents a relationship between nodes
type GraphRelationship struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`        // Relationship type
	FromNode   string                 `json:"from_node"`   // Source node ID
	ToNode     string                 `json:"to_node"`     // Target node ID
	Properties map[string]interface{} `json:"properties"`  // Flexible properties
	CreatedAt  time.Time              `json:"created_at"`
}

// ========== 场景1: 小说知识图谱 ==========

// NovelNode represents a novel
type NovelNode struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Genre       string    `json:"genre"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// WorldSettingNode represents world building elements
type WorldSettingNode struct {
	ID          string    `json:"id"`
	NovelID     string    `json:"novel_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`         // e.g., "magic_system", "technology", "culture"
	Description string    `json:"description"`
	Rules       string    `json:"rules"`        // JSON string of rules
	CreatedAt   time.Time `json:"created_at"`
}

// LocationNode represents a place in the novel world
type LocationNode struct {
	ID          string    `json:"id"`
	NovelID     string    `json:"novel_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`         // e.g., "city", "kingdom", "dungeon"
	Description string    `json:"description"`
	Coordinates string    `json:"coordinates"`  // Map coordinates if available
	CreatedAt   time.Time `json:"created_at"`
}

// CharacterNode represents a character in the novel
type CharacterNode struct {
	ID          string    `json:"id"`
	NovelID     string    `json:"novel_id"`
	Name        string    `json:"name"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	Role        string    `json:"role"`         // e.g., "protagonist", "antagonist", "supporting"
	Personality string    `json:"personality"`
	Backstory   string    `json:"backstory"`
	Attributes  string    `json:"attributes"`   // JSON string of attributes (strength, intelligence, etc.)
	CreatedAt   time.Time `json:"created_at"`
}

// FactionNode represents a group, organization, or faction
type FactionNode struct {
	ID          string    `json:"id"`
	NovelID     string    `json:"novel_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`         // e.g., "guild", "kingdom", "sect"
	Description string    `json:"description"`
	Power       int       `json:"power"`        // Power level or influence
	CreatedAt   time.Time `json:"created_at"`
}

// ========== 场景2: 代码编程知识库 ==========

// CodeFileNode represents a source code file
type CodeFileNode struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	FilePath    string    `json:"file_path"`
	Language    string    `json:"language"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// ClassNode represents a class or struct
type ClassNode struct {
	ID          string    `json:"id"`
	FileID      string    `json:"file_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`         // e.g., "class", "struct", "interface"
	Description string    `json:"description"`
	Modifiers   []string  `json:"modifiers"`    // e.g., ["public", "abstract"]
	CreatedAt   time.Time `json:"created_at"`
}

// FunctionNode represents a function or method
type FunctionNode struct {
	ID          string    `json:"id"`
	ClassID     string    `json:"class_id"`     // Optional, for methods
	FileID      string    `json:"file_id"`
	Name        string    `json:"name"`
	Parameters  string    `json:"parameters"`   // JSON string of parameters
	ReturnType  string    `json:"return_type"`
	Description string    `json:"description"`
	Complexity  int       `json:"complexity"`   // Cyclomatic complexity
	CreatedAt   time.Time `json:"created_at"`
}

// PackageNode represents a package or module
type PackageNode struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Repository  string    `json:"repository"`
	CreatedAt   time.Time `json:"created_at"`
}

// ========== 场景3: 普通知识库 ==========

// KnowledgeDocumentNode represents a document in knowledge base
type KnowledgeDocumentNode struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Content     string    `json:"content"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
}

// TopicNode represents a topic or subject
type TopicNode struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Level       int       `json:"level"`        // Hierarchy level
	CreatedAt   time.Time `json:"created_at"`
}

// ConceptNode represents a concept or idea
type ConceptNode struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Definition  string    `json:"definition"`
	Examples    []string  `json:"examples"`
	CreatedAt   time.Time `json:"created_at"`
}

// KnowledgeEntityNode represents an entity in knowledge base
type KnowledgeEntityNode struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`         // e.g., "person", "organization", "event"
	Description string    `json:"description"`
	Attributes  string    `json:"attributes"`   // JSON string of attributes
	CreatedAt   time.Time `json:"created_at"`
}

// ========== 关系类型常量 ==========

const (
	// Novel relationships
	RelNovelHasWorld      = "HAS_WORLD_SETTING"
	RelNovelHasLocation   = "HAS_LOCATION"
	RelNovelHasCharacter  = "HAS_CHARACTER"
	RelNovelHasFaction    = "HAS_FACTION"
	RelCharacterLocatedIn = "LOCATED_IN"
	RelCharacterKnows     = "KNOWS"
	RelCharacterEnemyOf   = "ENEMY_OF"
	RelCharacterMemberOf  = "MEMBER_OF"
	RelLocationContains   = "CONTAINS"
	RelFactionControls    = "CONTROLS"

	// Code relationships
	RelFileContainsClass    = "CONTAINS_CLASS"
	RelFileContainsFunction = "CONTAINS_FUNCTION"
	RelClassInherits        = "INHERITS"
	RelClassImplements      = "IMPLEMENTS"
	RelFunctionCalls        = "CALLS"
	RelPackageImports       = "IMPORTS"
	RelPackageDependsOn     = "DEPENDS_ON"

	// Knowledge relationships
	RelDocumentCovers       = "COVERS"
	RelDocumentReferences   = "REFERENCES"
	RelTopicContains        = "CONTAINS"
	RelTopicRelatedTo       = "RELATED_TO"
	RelConceptBelongsTo     = "BELONGS_TO"
	RelConceptDerivedFrom   = "DERIVED_FROM"
	RelEntityMentionedIn    = "MENTIONED_IN"
	RelEntityRelatedTo      = "RELATED_TO"
)
