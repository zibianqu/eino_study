package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zibianqu/eino_study/internal/model"
	"github.com/zibianqu/eino_study/internal/pkg/database"
)

// RelationshipGraphRepository defines the interface for relationship operations
type RelationshipGraphRepository interface {
	Create(ctx context.Context, rel *model.Relationship) error
	GetByID(ctx context.Context, id string) (*model.Relationship, error)
	Delete(ctx context.Context, id string) error
	GetRelationshipsBetween(ctx context.Context, fromID, toID string) ([]*model.Relationship, error)
	GetOutgoingRelationships(ctx context.Context, nodeID string, relType string) ([]*model.Relationship, error)
	GetIncomingRelationships(ctx context.Context, nodeID string, relType string) ([]*model.Relationship, error)
	CreateDocumentContainsEntity(ctx context.Context, docID, entityID string, properties map[string]interface{}) error
	CreateDocumentSimilarity(ctx context.Context, doc1ID, doc2ID string, score float64) error
	CreateDocumentReference(ctx context.Context, fromDocID, toDocID string) error
}

type relationshipGraphRepository struct {
	driver neo4j.DriverWithContext
}

// NewRelationshipGraphRepository creates a new RelationshipGraphRepository instance
func NewRelationshipGraphRepository() RelationshipGraphRepository {
	return &relationshipGraphRepository{
		driver: database.GetNeo4jDriver(),
	}
}

// Create creates a new relationship between two nodes
func (r *relationshipGraphRepository) Create(ctx context.Context, rel *model.Relationship) error {
	propertiesJSON, err := json.Marshal(rel.Properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	query := fmt.Sprintf(`
		MATCH (from {id: $from_id})
		MATCH (to {id: $to_id})
		CREATE (from)-[r:%s {properties: $properties, created_at: datetime($created_at)}]->(to)
		RETURN id(r) as rel_id
	`, rel.Type)

	params := map[string]interface{}{
		"from_id":    rel.FromNodeID,
		"to_id":      rel.ToNodeID,
		"properties": string(propertiesJSON),
		"created_at": rel.CreatedAt.Format(time.RFC3339),
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		run, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		if run.Next(ctx) {
			return run.Record().Values[0], nil
		}
		return nil, fmt.Errorf("failed to create relationship")
	})

	if err != nil {
		return fmt.Errorf("failed to create relationship: %w", err)
	}

	rel.ID = fmt.Sprintf("%v", result)
	return nil
}

// GetByID retrieves a relationship by its ID
func (r *relationshipGraphRepository) GetByID(ctx context.Context, id string) (*model.Relationship, error) {
	query := `
		MATCH (from)-[r]->(to)
		WHERE id(r) = $id
		RETURN id(r) as rel_id, type(r) as rel_type, id(from) as from_id, id(to) as to_id,
		       r.properties as properties, r.created_at as created_at
	`

	params := map[string]interface{}{"id": id}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		run, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if run.Next(ctx) {
			record := run.Record()
			rel := &model.Relationship{
				ID:         fmt.Sprintf("%v", record.Values[0]),
				Type:       record.Values[1].(string),
				FromNodeID: fmt.Sprintf("%v", record.Values[2]),
				ToNodeID:   fmt.Sprintf("%v", record.Values[3]),
			}

			if propsStr, ok := record.Values[4].(string); ok && propsStr != "" {
				var props map[string]interface{}
				if err := json.Unmarshal([]byte(propsStr), &props); err == nil {
					rel.Properties = props
				}
			}

			if createdAt, ok := record.Values[5].(time.Time); ok {
				rel.CreatedAt = createdAt
			}

			return rel, nil
		}

		return nil, fmt.Errorf("relationship not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Relationship), nil
}

// Delete deletes a relationship by its ID
func (r *relationshipGraphRepository) Delete(ctx context.Context, id string) error {
	query := `
		MATCH ()-[r]-()
		WHERE id(r) = $id
		DELETE r
	`

	params := map[string]interface{}{"id": id}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return fmt.Errorf("failed to delete relationship: %w", err)
	}

	return nil
}

// GetRelationshipsBetween retrieves all relationships between two nodes
func (r *relationshipGraphRepository) GetRelationshipsBetween(ctx context.Context, fromID, toID string) ([]*model.Relationship, error) {
	query := `
		MATCH (from {id: $from_id})-[r]->(to {id: $to_id})
		RETURN id(r) as rel_id, type(r) as rel_type, from.id as from_id, to.id as to_id,
		       r.properties as properties, r.created_at as created_at
	`

	params := map[string]interface{}{
		"from_id": fromID,
		"to_id":   toID,
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		run, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var relationships []*model.Relationship
		for run.Next(ctx) {
			record := run.Record()
			rel := &model.Relationship{
				ID:         fmt.Sprintf("%v", record.Values[0]),
				Type:       record.Values[1].(string),
				FromNodeID: record.Values[2].(string),
				ToNodeID:   record.Values[3].(string),
			}

			if propsStr, ok := record.Values[4].(string); ok && propsStr != "" {
				var props map[string]interface{}
				if err := json.Unmarshal([]byte(propsStr), &props); err == nil {
					rel.Properties = props
				}
			}

			if createdAt, ok := record.Values[5].(time.Time); ok {
				rel.CreatedAt = createdAt
			}

			relationships = append(relationships, rel)
		}

		return relationships, run.Err()
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.Relationship), nil
}

// GetOutgoingRelationships retrieves all outgoing relationships from a node
func (r *relationshipGraphRepository) GetOutgoingRelationships(ctx context.Context, nodeID string, relType string) ([]*model.Relationship, error) {
	var query string
	if relType != "" {
		query = fmt.Sprintf(`
			MATCH (from {id: $node_id})-[r:%s]->(to)
			RETURN id(r) as rel_id, type(r) as rel_type, from.id as from_id, to.id as to_id,
			       r.properties as properties, r.created_at as created_at
		`, relType)
	} else {
		query = `
			MATCH (from {id: $node_id})-[r]->(to)
			RETURN id(r) as rel_id, type(r) as rel_type, from.id as from_id, to.id as to_id,
			       r.properties as properties, r.created_at as created_at
		`
	}

	params := map[string]interface{}{"node_id": nodeID}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		run, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var relationships []*model.Relationship
		for run.Next(ctx) {
			record := run.Record()
			rel := &model.Relationship{
				ID:         fmt.Sprintf("%v", record.Values[0]),
				Type:       record.Values[1].(string),
				FromNodeID: record.Values[2].(string),
				ToNodeID:   record.Values[3].(string),
			}

			if propsStr, ok := record.Values[4].(string); ok && propsStr != "" {
				var props map[string]interface{}
				if err := json.Unmarshal([]byte(propsStr), &props); err == nil {
					rel.Properties = props
				}
			}

			if createdAt, ok := record.Values[5].(time.Time); ok {
				rel.CreatedAt = createdAt
			}

			relationships = append(relationships, rel)
		}

		return relationships, run.Err()
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.Relationship), nil
}

// GetIncomingRelationships retrieves all incoming relationships to a node
func (r *relationshipGraphRepository) GetIncomingRelationships(ctx context.Context, nodeID string, relType string) ([]*model.Relationship, error) {
	var query string
	if relType != "" {
		query = fmt.Sprintf(`
			MATCH (from)-[r:%s]->(to {id: $node_id})
			RETURN id(r) as rel_id, type(r) as rel_type, from.id as from_id, to.id as to_id,
			       r.properties as properties, r.created_at as created_at
		`, relType)
	} else {
		query = `
			MATCH (from)-[r]->(to {id: $node_id})
			RETURN id(r) as rel_id, type(r) as rel_type, from.id as from_id, to.id as to_id,
			       r.properties as properties, r.created_at as created_at
		`
	}

	params := map[string]interface{}{"node_id": nodeID}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		run, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var relationships []*model.Relationship
		for run.Next(ctx) {
			record := run.Record()
			rel := &model.Relationship{
				ID:         fmt.Sprintf("%v", record.Values[0]),
				Type:       record.Values[1].(string),
				FromNodeID: record.Values[2].(string),
				ToNodeID:   record.Values[3].(string),
			}

			if propsStr, ok := record.Values[4].(string); ok && propsStr != "" {
				var props map[string]interface{}
				if err := json.Unmarshal([]byte(propsStr), &props); err == nil {
					rel.Properties = props
				}
			}

			if createdAt, ok := record.Values[5].(time.Time); ok {
				rel.CreatedAt = createdAt
			}

			relationships = append(relationships, rel)
		}

		return relationships, run.Err()
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.Relationship), nil
}

// CreateDocumentContainsEntity creates a CONTAINS relationship from document to entity
func (r *relationshipGraphRepository) CreateDocumentContainsEntity(ctx context.Context, docID, entityID string, properties map[string]interface{}) error {
	propertiesJSON, _ := json.Marshal(properties)

	query := `
		MATCH (d:Document {id: $doc_id})
		MATCH (e:Entity {id: $entity_id})
		MERGE (d)-[r:CONTAINS {properties: $properties, created_at: datetime()}]->(e)
		RETURN r
	`

	params := map[string]interface{}{
		"doc_id":     docID,
		"entity_id": entityID,
		"properties": string(propertiesJSON),
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	return err
}

// CreateDocumentSimilarity creates a SIMILAR_TO relationship between two documents
func (r *relationshipGraphRepository) CreateDocumentSimilarity(ctx context.Context, doc1ID, doc2ID string, score float64) error {
	query := `
		MATCH (d1:Document {id: $doc1_id})
		MATCH (d2:Document {id: $doc2_id})
		MERGE (d1)-[r:SIMILAR_TO {score: $score, created_at: datetime()}]-(d2)
		RETURN r
	`

	params := map[string]interface{}{
		"doc1_id": doc1ID,
		"doc2_id": doc2ID,
		"score":   score,
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	return err
}

// CreateDocumentReference creates a REFERENCES relationship from one document to another
func (r *relationshipGraphRepository) CreateDocumentReference(ctx context.Context, fromDocID, toDocID string) error {
	query := `
		MATCH (d1:Document {id: $from_doc_id})
		MATCH (d2:Document {id: $to_doc_id})
		MERGE (d1)-[r:REFERENCES {created_at: datetime()}]->(d2)
		RETURN r
	`

	params := map[string]interface{}{
		"from_doc_id": fromDocID,
		"to_doc_id":   toDocID,
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	return err
}
