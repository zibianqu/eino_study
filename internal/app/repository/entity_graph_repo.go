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

// EntityGraphRepository defines the interface for entity graph operations
type EntityGraphRepository interface {
	Create(ctx context.Context, entity *model.EntityNode) error
	GetByID(ctx context.Context, id string) (*model.EntityNode, error)
	Update(ctx context.Context, entity *model.EntityNode) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*model.EntityNode, error)
	FindByType(ctx context.Context, entityType string, limit int) ([]*model.EntityNode, error)
	FindByName(ctx context.Context, entityName string) ([]*model.EntityNode, error)
	GetRelatedDocuments(ctx context.Context, entityID string) ([]*model.DocumentNode, error)
}

type entityGraphRepository struct {
	driver neo4j.DriverWithContext
}

// NewEntityGraphRepository creates a new EntityGraphRepository instance
func NewEntityGraphRepository() EntityGraphRepository {
	return &entityGraphRepository{
		driver: database.GetNeo4jDriver(),
	}
}

// Create creates a new entity node
func (r *entityGraphRepository) Create(ctx context.Context, entity *model.EntityNode) error {
	metadataJSON, err := json.Marshal(entity.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		CREATE (e:Entity {
			id: $id,
			entity_type: $entity_type,
			entity_name: $entity_name,
			entity_value: $entity_value,
			metadata: $metadata,
			ctime: datetime($ctime)
		})
		RETURN e
	`

	params := map[string]interface{}{
		"id":           entity.ID,
		"entity_type":  entity.EntityType,
		"entity_name":  entity.EntityName,
		"entity_value": entity.EntityValue,
		"metadata":     string(metadataJSON),
		"ctime":        entity.CTime.Format(time.RFC3339),
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return fmt.Errorf("failed to create entity node: %w", err)
	}

	return nil
}

// GetByID retrieves an entity node by its ID
func (r *entityGraphRepository) GetByID(ctx context.Context, id string) (*model.EntityNode, error) {
	query := `
		MATCH (e:Entity {id: $id})
		RETURN e.id as id, e.entity_type as entity_type, e.entity_name as entity_name,
		       e.entity_value as entity_value, e.metadata as metadata, e.ctime as ctime
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
			entity := &model.EntityNode{
				ID:          record.Values[0].(string),
				EntityType:  record.Values[1].(string),
				EntityName:  record.Values[2].(string),
				EntityValue: record.Values[3].(string),
			}

			// Parse metadata JSON
			if metadataStr, ok := record.Values[4].(string); ok && metadataStr != "" {
				var metadata map[string]interface{}
				if err := json.Unmarshal([]byte(metadataStr), &metadata); err == nil {
					entity.Metadata = metadata
				}
			}

			if ctimeVal, ok := record.Values[5].(time.Time); ok {
				entity.CTime = ctimeVal
			}

			return entity, nil
		}

		return nil, fmt.Errorf("entity not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.EntityNode), nil
}

// Update updates an entity node
func (r *entityGraphRepository) Update(ctx context.Context, entity *model.EntityNode) error {
	metadataJSON, err := json.Marshal(entity.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		MATCH (e:Entity {id: $id})
		SET e.entity_type = $entity_type,
		    e.entity_name = $entity_name,
		    e.entity_value = $entity_value,
		    e.metadata = $metadata
		RETURN e
	`

	params := map[string]interface{}{
		"id":           entity.ID,
		"entity_type":  entity.EntityType,
		"entity_name":  entity.EntityName,
		"entity_value": entity.EntityValue,
		"metadata":     string(metadataJSON),
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return fmt.Errorf("failed to update entity node: %w", err)
	}

	return nil
}

// Delete deletes an entity node and all its relationships
func (r *entityGraphRepository) Delete(ctx context.Context, id string) error {
	query := `
		MATCH (e:Entity {id: $id})
		DETACH DELETE e
	`

	params := map[string]interface{}{"id": id}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return fmt.Errorf("failed to delete entity node: %w", err)
	}

	return nil
}

// List retrieves a list of entity nodes with pagination
func (r *entityGraphRepository) List(ctx context.Context, limit, offset int) ([]*model.EntityNode, error) {
	query := `
		MATCH (e:Entity)
		RETURN e.id as id, e.entity_type as entity_type, e.entity_name as entity_name,
		       e.entity_value as entity_value, e.metadata as metadata, e.ctime as ctime
		ORDER BY e.ctime DESC
		SKIP $offset LIMIT $limit
	`

	params := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		run, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var entities []*model.EntityNode
		for run.Next(ctx) {
			record := run.Record()
			entity := &model.EntityNode{
				ID:          record.Values[0].(string),
				EntityType:  record.Values[1].(string),
				EntityName:  record.Values[2].(string),
				EntityValue: record.Values[3].(string),
			}

			// Parse metadata JSON
			if metadataStr, ok := record.Values[4].(string); ok && metadataStr != "" {
				var metadata map[string]interface{}
				if err := json.Unmarshal([]byte(metadataStr), &metadata); err == nil {
					entity.Metadata = metadata
				}
			}

			if ctimeVal, ok := record.Values[5].(time.Time); ok {
				entity.CTime = ctimeVal
			}

			entities = append(entities, entity)
		}

		return entities, run.Err()
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.EntityNode), nil
}

// FindByType finds entities by their type
func (r *entityGraphRepository) FindByType(ctx context.Context, entityType string, limit int) ([]*model.EntityNode, error) {
	query := `
		MATCH (e:Entity {entity_type: $entity_type})
		RETURN e.id as id, e.entity_type as entity_type, e.entity_name as entity_name,
		       e.entity_value as entity_value, e.metadata as metadata, e.ctime as ctime
		ORDER BY e.entity_name
		LIMIT $limit
	`

	params := map[string]interface{}{
		"entity_type": entityType,
		"limit":       limit,
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		run, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var entities []*model.EntityNode
		for run.Next(ctx) {
			record := run.Record()
			entity := &model.EntityNode{
				ID:          record.Values[0].(string),
				EntityType:  record.Values[1].(string),
				EntityName:  record.Values[2].(string),
				EntityValue: record.Values[3].(string),
			}

			if metadataStr, ok := record.Values[4].(string); ok && metadataStr != "" {
				var metadata map[string]interface{}
				if err := json.Unmarshal([]byte(metadataStr), &metadata); err == nil {
					entity.Metadata = metadata
				}
			}

			if ctimeVal, ok := record.Values[5].(time.Time); ok {
				entity.CTime = ctimeVal
			}

			entities = append(entities, entity)
		}

		return entities, run.Err()
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.EntityNode), nil
}

// FindByName finds entities by their name
func (r *entityGraphRepository) FindByName(ctx context.Context, entityName string) ([]*model.EntityNode, error) {
	query := `
		MATCH (e:Entity)
		WHERE e.entity_name CONTAINS $entity_name
		RETURN e.id as id, e.entity_type as entity_type, e.entity_name as entity_name,
		       e.entity_value as entity_value, e.metadata as metadata, e.ctime as ctime
		ORDER BY e.entity_name
		LIMIT 50
	`

	params := map[string]interface{}{"entity_name": entityName}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		run, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var entities []*model.EntityNode
		for run.Next(ctx) {
			record := run.Record()
			entity := &model.EntityNode{
				ID:          record.Values[0].(string),
				EntityType:  record.Values[1].(string),
				EntityName:  record.Values[2].(string),
				EntityValue: record.Values[3].(string),
			}

			if metadataStr, ok := record.Values[4].(string); ok && metadataStr != "" {
				var metadata map[string]interface{}
				if err := json.Unmarshal([]byte(metadataStr), &metadata); err == nil {
					entity.Metadata = metadata
				}
			}

			if ctimeVal, ok := record.Values[5].(time.Time); ok {
				entity.CTime = ctimeVal
			}

			entities = append(entities, entity)
		}

		return entities, run.Err()
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.EntityNode), nil
}

// GetRelatedDocuments retrieves documents related to an entity
func (r *entityGraphRepository) GetRelatedDocuments(ctx context.Context, entityID string) ([]*model.DocumentNode, error) {
	query := `
		MATCH (d:Document)-[:CONTAINS]->(e:Entity {id: $id})
		RETURN d.id as id, d.doc_name as doc_name, d.doc_hash as doc_hash,
		       d.file_path as file_path, d.file_type as file_type, d.ctime as ctime
		ORDER BY d.ctime DESC
	`

	params := map[string]interface{}{"id": entityID}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		run, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var docs []*model.DocumentNode
		for run.Next(ctx) {
			record := run.Record()
			doc := &model.DocumentNode{
				ID:       record.Values[0].(string),
				DocName:  record.Values[1].(string),
				DocHash:  record.Values[2].(string),
				FilePath: record.Values[3].(string),
				FileType: record.Values[4].(string),
			}

			if ctimeVal, ok := record.Values[5].(time.Time); ok {
				doc.CTime = ctimeVal
			}

			docs = append(docs, doc)
		}

		return docs, run.Err()
	})

	if err != nil {
		return nil, err
	}

	return result.([]*model.DocumentNode), nil
}
