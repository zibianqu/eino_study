package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zibianqu/eino_study/internal/model"
	"github.com/zibianqu/eino_study/internal/pkg/database"
)

// DocumentGraphRepository defines the interface for document graph operations
type DocumentGraphRepository interface {
	Create(ctx context.Context, doc *model.DocumentNode) error
	GetByID(ctx context.Context, id string) (*model.DocumentNode, error)
	Update(ctx context.Context, doc *model.DocumentNode) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*model.DocumentNode, error)
	FindSimilar(ctx context.Context, docID string, limit int) ([]*model.DocumentNode, error)
	GetRelatedEntities(ctx context.Context, docID string) ([]*model.EntityNode, error)
}

type documentGraphRepository struct {
	driver neo4j.DriverWithContext
}

// NewDocumentGraphRepository creates a new DocumentGraphRepository instance
func NewDocumentGraphRepository() DocumentGraphRepository {
	return &documentGraphRepository{
		driver: database.GetNeo4jDriver(),
	}
}

// Create creates a new document node
func (r *documentGraphRepository) Create(ctx context.Context, doc *model.DocumentNode) error {
	query := `
		CREATE (d:Document {
			id: $id,
			doc_name: $doc_name,
			doc_hash: $doc_hash,
			file_path: $file_path,
			file_type: $file_type,
			ctime: datetime($ctime)
		})
		RETURN d
	`

	params := map[string]interface{}{
		"id":        doc.ID,
		"doc_name":  doc.DocName,
		"doc_hash":  doc.DocHash,
		"file_path": doc.FilePath,
		"file_type": doc.FileType,
		"ctime":     doc.CTime.Format(time.RFC3339),
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return fmt.Errorf("failed to create document node: %w", err)
	}

	return nil
}

// GetByID retrieves a document node by its ID
func (r *documentGraphRepository) GetByID(ctx context.Context, id string) (*model.DocumentNode, error) {
	query := `
		MATCH (d:Document {id: $id})
		RETURN d.id as id, d.doc_name as doc_name, d.doc_hash as doc_hash,
		       d.file_path as file_path, d.file_type as file_type, d.ctime as ctime
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

			return doc, nil
		}

		return nil, fmt.Errorf("document not found")
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.DocumentNode), nil
}

// Update updates a document node
func (r *documentGraphRepository) Update(ctx context.Context, doc *model.DocumentNode) error {
	query := `
		MATCH (d:Document {id: $id})
		SET d.doc_name = $doc_name,
		    d.doc_hash = $doc_hash,
		    d.file_path = $file_path,
		    d.file_type = $file_type
		RETURN d
	`

	params := map[string]interface{}{
		"id":        doc.ID,
		"doc_name":  doc.DocName,
		"doc_hash":  doc.DocHash,
		"file_path": doc.FilePath,
		"file_type": doc.FileType,
	}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return fmt.Errorf("failed to update document node: %w", err)
	}

	return nil
}

// Delete deletes a document node and all its relationships
func (r *documentGraphRepository) Delete(ctx context.Context, id string) error {
	query := `
		MATCH (d:Document {id: $id})
		DETACH DELETE d
	`

	params := map[string]interface{}{"id": id}

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return fmt.Errorf("failed to delete document node: %w", err)
	}

	return nil
}

// List retrieves a list of document nodes with pagination
func (r *documentGraphRepository) List(ctx context.Context, limit, offset int) ([]*model.DocumentNode, error) {
	query := `
		MATCH (d:Document)
		RETURN d.id as id, d.doc_name as doc_name, d.doc_hash as doc_hash,
		       d.file_path as file_path, d.file_type as file_type, d.ctime as ctime
		ORDER BY d.ctime DESC
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

// FindSimilar finds similar documents based on SIMILAR_TO relationships
func (r *documentGraphRepository) FindSimilar(ctx context.Context, docID string, limit int) ([]*model.DocumentNode, error) {
	query := `
		MATCH (d1:Document {id: $id})-[r:SIMILAR_TO]-(d2:Document)
		RETURN d2.id as id, d2.doc_name as doc_name, d2.doc_hash as doc_hash,
		       d2.file_path as file_path, d2.file_type as file_type, d2.ctime as ctime,
		       r.score as similarity_score
		ORDER BY r.score DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"id":    docID,
		"limit": limit,
	}

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

// GetRelatedEntities retrieves entities related to a document via CONTAINS relationship
func (r *documentGraphRepository) GetRelatedEntities(ctx context.Context, docID string) ([]*model.EntityNode, error) {
	query := `
		MATCH (d:Document {id: $id})-[:CONTAINS]->(e:Entity)
		RETURN e.id as id, e.entity_type as entity_type, e.entity_name as entity_name,
		       e.entity_value as entity_value, e.ctime as ctime
		ORDER BY e.entity_name
	`

	params := map[string]interface{}{"id": docID}

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

			if ctimeVal, ok := record.Values[4].(time.Time); ok {
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
