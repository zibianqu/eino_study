package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zibianqu/eino_study/internal/model"
	"github.com/zibianqu/eino_study/internal/pkg/database"
)

// KnowledgeGraphRepository handles general knowledge graph operations
type KnowledgeGraphRepository interface {
	// Document operations
	CreateDocument(ctx context.Context, doc *model.KnowledgeDocumentNode) error
	GetDocument(ctx context.Context, id string) (*model.KnowledgeDocumentNode, error)
	UpdateDocument(ctx context.Context, doc *model.KnowledgeDocumentNode) error
	DeleteDocument(ctx context.Context, id string) error
	SearchDocuments(ctx context.Context, keyword string, limit int) ([]*model.KnowledgeDocumentNode, error)

	// Topic operations
	CreateTopic(ctx context.Context, topic *model.TopicNode) error
	GetTopic(ctx context.Context, id string) (*model.TopicNode, error)
	UpdateTopic(ctx context.Context, topic *model.TopicNode) error
	DeleteTopic(ctx context.Context, id string) error
	GetTopicHierarchy(ctx context.Context, rootID string) ([]*model.TopicNode, error)

	// Concept operations
	CreateConcept(ctx context.Context, concept *model.ConceptNode) error
	GetConcept(ctx context.Context, id string) (*model.ConceptNode, error)
	UpdateConcept(ctx context.Context, concept *model.ConceptNode) error
	DeleteConcept(ctx context.Context, id string) error

	// Entity operations
	CreateEntity(ctx context.Context, entity *model.KnowledgeEntityNode) error
	GetEntity(ctx context.Context, id string) (*model.KnowledgeEntityNode, error)
	UpdateEntity(ctx context.Context, entity *model.KnowledgeEntityNode) error
	DeleteEntity(ctx context.Context, id string) error

	// Relationship operations
	LinkDocumentToTopic(ctx context.Context, docID, topicID string) error
	LinkConceptToTopic(ctx context.Context, conceptID, topicID string) error
	GetRelatedDocuments(ctx context.Context, docID string, depth int) ([]*model.KnowledgeDocumentNode, error)
}

type knowledgeGraphRepository struct {
	driver neo4j.DriverWithContext
}

// NewKnowledgeGraphRepository creates a new knowledge graph repository
func NewKnowledgeGraphRepository() KnowledgeGraphRepository {
	return &knowledgeGraphRepository{
		driver: database.GetNeo4jDriver(),
	}
}

// CreateDocument creates a knowledge document node
func (r *knowledgeGraphRepository) CreateDocument(ctx context.Context, doc *model.KnowledgeDocumentNode) error {
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = time.Now()
	}

	query := `
		CREATE (d:KnowledgeDocument {
			id: $id,
			title: $title,
			category: $category,
			content: $content,
			tags: $tags,
			created_at: datetime($created_at)
		})
		RETURN d
	`

	params := map[string]interface{}{
		"id":         doc.ID,
		"title":      doc.Title,
		"category":   doc.Category,
		"content":    doc.Content,
		"tags":       doc.Tags,
		"created_at": doc.CreatedAt.Format(time.RFC3339),
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetDocument retrieves a document by ID
func (r *knowledgeGraphRepository) GetDocument(ctx context.Context, id string) (*model.KnowledgeDocumentNode, error) {
	query := "MATCH (d:KnowledgeDocument {id: $id}) RETURN d"
	params := map[string]interface{}{"id": id}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("document not found")
	}

	doc := &model.KnowledgeDocumentNode{ID: id}
	return doc, nil
}

// UpdateDocument updates a document node
func (r *knowledgeGraphRepository) UpdateDocument(ctx context.Context, doc *model.KnowledgeDocumentNode) error {
	query := `
		MATCH (d:KnowledgeDocument {id: $id})
		SET d.title = $title,
			d.category = $category,
			d.content = $content,
			d.tags = $tags,
			d.updated_at = datetime()
		RETURN d
	`

	params := map[string]interface{}{
		"id":       doc.ID,
		"title":    doc.Title,
		"category": doc.Category,
		"content":  doc.Content,
		"tags":     doc.Tags,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// DeleteDocument deletes a document node
func (r *knowledgeGraphRepository) DeleteDocument(ctx context.Context, id string) error {
	query := "MATCH (d:KnowledgeDocument {id: $id}) DETACH DELETE d"
	params := map[string]interface{}{"id": id}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// SearchDocuments searches for documents by keyword
func (r *knowledgeGraphRepository) SearchDocuments(ctx context.Context, keyword string, limit int) ([]*model.KnowledgeDocumentNode, error) {
	query := `
		MATCH (d:KnowledgeDocument)
		WHERE d.title CONTAINS $keyword OR d.content CONTAINS $keyword
		RETURN d
		LIMIT $limit
	`

	params := map[string]interface{}{
		"keyword": keyword,
		"limit":   limit,
	}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	docs := make([]*model.KnowledgeDocumentNode, 0, len(records))
	return docs, nil
}

// CreateTopic creates a topic node
func (r *knowledgeGraphRepository) CreateTopic(ctx context.Context, topic *model.TopicNode) error {
	if topic.CreatedAt.IsZero() {
		topic.CreatedAt = time.Now()
	}

	query := `
		CREATE (t:Topic {
			id: $id,
			name: $name,
			description: $description,
			level: $level,
			created_at: datetime($created_at)
		})
		RETURN t
	`

	params := map[string]interface{}{
		"id":          topic.ID,
		"name":        topic.Name,
		"description": topic.Description,
		"level":       topic.Level,
		"created_at":  topic.CreatedAt.Format(time.RFC3339),
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetTopic retrieves a topic by ID
func (r *knowledgeGraphRepository) GetTopic(ctx context.Context, id string) (*model.TopicNode, error) {
	query := "MATCH (t:Topic {id: $id}) RETURN t"
	params := map[string]interface{}{"id": id}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("topic not found")
	}

	topic := &model.TopicNode{ID: id}
	return topic, nil
}

// UpdateTopic updates a topic node
func (r *knowledgeGraphRepository) UpdateTopic(ctx context.Context, topic *model.TopicNode) error {
	query := `
		MATCH (t:Topic {id: $id})
		SET t.name = $name,
			t.description = $description,
			t.level = $level,
			t.updated_at = datetime()
		RETURN t
	`

	params := map[string]interface{}{
		"id":          topic.ID,
		"name":        topic.Name,
		"description": topic.Description,
		"level":       topic.Level,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// DeleteTopic deletes a topic node
func (r *knowledgeGraphRepository) DeleteTopic(ctx context.Context, id string) error {
	query := "MATCH (t:Topic {id: $id}) DETACH DELETE t"
	params := map[string]interface{}{"id": id}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetTopicHierarchy retrieves the topic hierarchy starting from a root topic
func (r *knowledgeGraphRepository) GetTopicHierarchy(ctx context.Context, rootID string) ([]*model.TopicNode, error) {
	query := `
		MATCH path = (root:Topic {id: $root_id})-[:CONTAINS*0..5]->(child:Topic)
		RETURN DISTINCT child
		ORDER BY child.level
	`

	params := map[string]interface{}{"root_id": rootID}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	topics := make([]*model.TopicNode, 0, len(records))
	return topics, nil
}

// CreateConcept creates a concept node
func (r *knowledgeGraphRepository) CreateConcept(ctx context.Context, concept *model.ConceptNode) error {
	if concept.CreatedAt.IsZero() {
		concept.CreatedAt = time.Now()
	}

	query := `
		CREATE (c:Concept {
			id: $id,
			name: $name,
			definition: $definition,
			examples: $examples,
			created_at: datetime($created_at)
		})
		RETURN c
	`

	params := map[string]interface{}{
		"id":         concept.ID,
		"name":       concept.Name,
		"definition": concept.Definition,
		"examples":   concept.Examples,
		"created_at": concept.CreatedAt.Format(time.RFC3339),
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetConcept retrieves a concept by ID
func (r *knowledgeGraphRepository) GetConcept(ctx context.Context, id string) (*model.ConceptNode, error) {
	query := "MATCH (c:Concept {id: $id}) RETURN c"
	params := map[string]interface{}{"id": id}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("concept not found")
	}

	concept := &model.ConceptNode{ID: id}
	return concept, nil
}

// UpdateConcept updates a concept node
func (r *knowledgeGraphRepository) UpdateConcept(ctx context.Context, concept *model.ConceptNode) error {
	query := `
		MATCH (c:Concept {id: $id})
		SET c.name = $name,
			c.definition = $definition,
			c.examples = $examples,
			c.updated_at = datetime()
		RETURN c
	`

	params := map[string]interface{}{
		"id":         concept.ID,
		"name":       concept.Name,
		"definition": concept.Definition,
		"examples":   concept.Examples,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// DeleteConcept deletes a concept node
func (r *knowledgeGraphRepository) DeleteConcept(ctx context.Context, id string) error {
	query := "MATCH (c:Concept {id: $id}) DETACH DELETE c"
	params := map[string]interface{}{"id": id}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// CreateEntity creates a knowledge entity node
func (r *knowledgeGraphRepository) CreateEntity(ctx context.Context, entity *model.KnowledgeEntityNode) error {
	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}

	query := `
		CREATE (e:KnowledgeEntity {
			id: $id,
			name: $name,
			type: $type,
			description: $description,
			attributes: $attributes,
			created_at: datetime($created_at)
		})
		RETURN e
	`

	params := map[string]interface{}{
		"id":          entity.ID,
		"name":        entity.Name,
		"type":        entity.Type,
		"description": entity.Description,
		"attributes":  entity.Attributes,
		"created_at":  entity.CreatedAt.Format(time.RFC3339),
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetEntity retrieves an entity by ID
func (r *knowledgeGraphRepository) GetEntity(ctx context.Context, id string) (*model.KnowledgeEntityNode, error) {
	query := "MATCH (e:KnowledgeEntity {id: $id}) RETURN e"
	params := map[string]interface{}{"id": id}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("entity not found")
	}

	entity := &model.KnowledgeEntityNode{ID: id}
	return entity, nil
}

// UpdateEntity updates an entity node
func (r *knowledgeGraphRepository) UpdateEntity(ctx context.Context, entity *model.KnowledgeEntityNode) error {
	query := `
		MATCH (e:KnowledgeEntity {id: $id})
		SET e.name = $name,
			e.type = $type,
			e.description = $description,
			e.attributes = $attributes,
			e.updated_at = datetime()
		RETURN e
	`

	params := map[string]interface{}{
		"id":          entity.ID,
		"name":        entity.Name,
		"type":        entity.Type,
		"description": entity.Description,
		"attributes":  entity.Attributes,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// DeleteEntity deletes an entity node
func (r *knowledgeGraphRepository) DeleteEntity(ctx context.Context, id string) error {
	query := "MATCH (e:KnowledgeEntity {id: $id}) DETACH DELETE e"
	params := map[string]interface{}{"id": id}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// LinkDocumentToTopic links a document to a topic
func (r *knowledgeGraphRepository) LinkDocumentToTopic(ctx context.Context, docID, topicID string) error {
	query := `
		MATCH (d:KnowledgeDocument {id: $doc_id})
		MATCH (t:Topic {id: $topic_id})
		CREATE (d)-[:COVERS]->(t)
		RETURN d, t
	`

	params := map[string]interface{}{
		"doc_id":   docID,
		"topic_id": topicID,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// LinkConceptToTopic links a concept to a topic
func (r *knowledgeGraphRepository) LinkConceptToTopic(ctx context.Context, conceptID, topicID string) error {
	query := `
		MATCH (c:Concept {id: $concept_id})
		MATCH (t:Topic {id: $topic_id})
		CREATE (c)-[:BELONGS_TO]->(t)
		RETURN c, t
	`

	params := map[string]interface{}{
		"concept_id": conceptID,
		"topic_id":   topicID,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetRelatedDocuments retrieves documents related to a given document
func (r *knowledgeGraphRepository) GetRelatedDocuments(ctx context.Context, docID string, depth int) ([]*model.KnowledgeDocumentNode, error) {
	query := fmt.Sprintf(`
		MATCH path = (d:KnowledgeDocument {id: $doc_id})-[:REFERENCES|RELATED_TO*1..%d]-(related:KnowledgeDocument)
		RETURN DISTINCT related
	`, depth)

	params := map[string]interface{}{"doc_id": docID}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	docs := make([]*model.KnowledgeDocumentNode, 0, len(records))
	return docs, nil
}
