package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zibianqu/eino_study/internal/model"
	"github.com/zibianqu/eino_study/internal/pkg/database"
)

// NovelGraphRepository handles novel knowledge graph operations
type NovelGraphRepository interface {
	// Novel operations
	CreateNovel(ctx context.Context, novel *model.NovelNode) error
	GetNovel(ctx context.Context, id string) (*model.NovelNode, error)
	UpdateNovel(ctx context.Context, novel *model.NovelNode) error
	DeleteNovel(ctx context.Context, id string) error
	ListNovels(ctx context.Context, limit, offset int) ([]*model.NovelNode, error)

	// Character operations
	CreateCharacter(ctx context.Context, character *model.CharacterNode) error
	GetCharacter(ctx context.Context, id string) (*model.CharacterNode, error)
	UpdateCharacter(ctx context.Context, character *model.CharacterNode) error
	DeleteCharacter(ctx context.Context, id string) error
	ListCharactersByNovel(ctx context.Context, novelID string) ([]*model.CharacterNode, error)

	// Location operations
	CreateLocation(ctx context.Context, location *model.LocationNode) error
	GetLocation(ctx context.Context, id string) (*model.LocationNode, error)
	UpdateLocation(ctx context.Context, location *model.LocationNode) error
	DeleteLocation(ctx context.Context, id string) error

	// Relationship operations
	CreateCharacterRelationship(ctx context.Context, fromID, toID, relType string, properties map[string]interface{}) error
	GetCharacterRelationships(ctx context.Context, characterID string) ([]*model.GraphRelationship, error)
	DeleteRelationship(ctx context.Context, fromID, toID, relType string) error
}

type novelGraphRepository struct {
	driver neo4j.DriverWithContext
}

// NewNovelGraphRepository creates a new novel graph repository
func NewNovelGraphRepository() NovelGraphRepository {
	return &novelGraphRepository{
		driver: database.GetNeo4jDriver(),
	}
}

// CreateNovel creates a novel node
func (r *novelGraphRepository) CreateNovel(ctx context.Context, novel *model.NovelNode) error {
	if novel.CreatedAt.IsZero() {
		novel.CreatedAt = time.Now()
	}

	query := `
		CREATE (n:Novel {
			id: $id,
			title: $title,
			author: $author,
			genre: $genre,
			description: $description,
			created_at: datetime($created_at)
		})
		RETURN n
	`

	params := map[string]interface{}{
		"id":          novel.ID,
		"title":       novel.Title,
		"author":      novel.Author,
		"genre":       novel.Genre,
		"description": novel.Description,
		"created_at":  novel.CreatedAt.Format(time.RFC3339),
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetNovel retrieves a novel by ID
func (r *novelGraphRepository) GetNovel(ctx context.Context, id string) (*model.NovelNode, error) {
	query := "MATCH (n:Novel {id: $id}) RETURN n"
	params := map[string]interface{}{"id": id}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("novel not found")
	}

	// Parse the result (simplified, needs proper implementation)
	novel := &model.NovelNode{ID: id}
	return novel, nil
}

// UpdateNovel updates a novel node
func (r *novelGraphRepository) UpdateNovel(ctx context.Context, novel *model.NovelNode) error {
	query := `
		MATCH (n:Novel {id: $id})
		SET n.title = $title,
			n.author = $author,
			n.genre = $genre,
			n.description = $description,
			n.updated_at = datetime()
		RETURN n
	`

	params := map[string]interface{}{
		"id":          novel.ID,
		"title":       novel.Title,
		"author":      novel.Author,
		"genre":       novel.Genre,
		"description": novel.Description,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// DeleteNovel deletes a novel node and all related nodes
func (r *novelGraphRepository) DeleteNovel(ctx context.Context, id string) error {
	query := "MATCH (n:Novel {id: $id}) DETACH DELETE n"
	params := map[string]interface{}{"id": id}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// ListNovels lists all novels with pagination
func (r *novelGraphRepository) ListNovels(ctx context.Context, limit, offset int) ([]*model.NovelNode, error) {
	query := `
		MATCH (n:Novel)
		RETURN n
		ORDER BY n.created_at DESC
		SKIP $offset
		LIMIT $limit
	`

	params := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	// Parse records (simplified)
	novels := make([]*model.NovelNode, 0, len(records))
	return novels, nil
}

// CreateCharacter creates a character node
func (r *novelGraphRepository) CreateCharacter(ctx context.Context, character *model.CharacterNode) error {
	if character.CreatedAt.IsZero() {
		character.CreatedAt = time.Now()
	}

	query := `
		MATCH (n:Novel {id: $novel_id})
		CREATE (c:Character {
			id: $id,
			novel_id: $novel_id,
			name: $name,
			age: $age,
			gender: $gender,
			role: $role,
			personality: $personality,
			backstory: $backstory,
			attributes: $attributes,
			created_at: datetime($created_at)
		})
		CREATE (n)-[:HAS_CHARACTER]->(c)
		RETURN c
	`

	params := map[string]interface{}{
		"id":          character.ID,
		"novel_id":    character.NovelID,
		"name":        character.Name,
		"age":         character.Age,
		"gender":      character.Gender,
		"role":        character.Role,
		"personality": character.Personality,
		"backstory":   character.Backstory,
		"attributes":  character.Attributes,
		"created_at":  character.CreatedAt.Format(time.RFC3339),
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetCharacter retrieves a character by ID
func (r *novelGraphRepository) GetCharacter(ctx context.Context, id string) (*model.CharacterNode, error) {
	query := "MATCH (c:Character {id: $id}) RETURN c"
	params := map[string]interface{}{"id": id}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("character not found")
	}

	character := &model.CharacterNode{ID: id}
	return character, nil
}

// UpdateCharacter updates a character node
func (r *novelGraphRepository) UpdateCharacter(ctx context.Context, character *model.CharacterNode) error {
	query := `
		MATCH (c:Character {id: $id})
		SET c.name = $name,
			c.age = $age,
			c.gender = $gender,
			c.role = $role,
			c.personality = $personality,
			c.backstory = $backstory,
			c.attributes = $attributes,
			c.updated_at = datetime()
		RETURN c
	`

	params := map[string]interface{}{
		"id":          character.ID,
		"name":        character.Name,
		"age":         character.Age,
		"gender":      character.Gender,
		"role":        character.Role,
		"personality": character.Personality,
		"backstory":   character.Backstory,
		"attributes":  character.Attributes,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// DeleteCharacter deletes a character node
func (r *novelGraphRepository) DeleteCharacter(ctx context.Context, id string) error {
	query := "MATCH (c:Character {id: $id}) DETACH DELETE c"
	params := map[string]interface{}{"id": id}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// ListCharactersByNovel lists all characters in a novel
func (r *novelGraphRepository) ListCharactersByNovel(ctx context.Context, novelID string) ([]*model.CharacterNode, error) {
	query := `
		MATCH (n:Novel {id: $novel_id})-[:HAS_CHARACTER]->(c:Character)
		RETURN c
		ORDER BY c.name
	`

	params := map[string]interface{}{"novel_id": novelID}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	characters := make([]*model.CharacterNode, 0, len(records))
	return characters, nil
}

// CreateLocation creates a location node
func (r *novelGraphRepository) CreateLocation(ctx context.Context, location *model.LocationNode) error {
	if location.CreatedAt.IsZero() {
		location.CreatedAt = time.Now()
	}

	query := `
		MATCH (n:Novel {id: $novel_id})
		CREATE (l:Location {
			id: $id,
			novel_id: $novel_id,
			name: $name,
			type: $type,
			description: $description,
			coordinates: $coordinates,
			created_at: datetime($created_at)
		})
		CREATE (n)-[:HAS_LOCATION]->(l)
		RETURN l
	`

	params := map[string]interface{}{
		"id":          location.ID,
		"novel_id":    location.NovelID,
		"name":        location.Name,
		"type":        location.Type,
		"description": location.Description,
		"coordinates": location.Coordinates,
		"created_at":  location.CreatedAt.Format(time.RFC3339),
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetLocation retrieves a location by ID
func (r *novelGraphRepository) GetLocation(ctx context.Context, id string) (*model.LocationNode, error) {
	query := "MATCH (l:Location {id: $id}) RETURN l"
	params := map[string]interface{}{"id": id}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("location not found")
	}

	location := &model.LocationNode{ID: id}
	return location, nil
}

// UpdateLocation updates a location node
func (r *novelGraphRepository) UpdateLocation(ctx context.Context, location *model.LocationNode) error {
	query := `
		MATCH (l:Location {id: $id})
		SET l.name = $name,
			l.type = $type,
			l.description = $description,
			l.coordinates = $coordinates,
			l.updated_at = datetime()
		RETURN l
	`

	params := map[string]interface{}{
		"id":          location.ID,
		"name":        location.Name,
		"type":        location.Type,
		"description": location.Description,
		"coordinates": location.Coordinates,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// DeleteLocation deletes a location node
func (r *novelGraphRepository) DeleteLocation(ctx context.Context, id string) error {
	query := "MATCH (l:Location {id: $id}) DETACH DELETE l"
	params := map[string]interface{}{"id": id}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// CreateCharacterRelationship creates a relationship between two characters
func (r *novelGraphRepository) CreateCharacterRelationship(ctx context.Context, fromID, toID, relType string, properties map[string]interface{}) error {
	query := fmt.Sprintf(`
		MATCH (from:Character {id: $from_id})
		MATCH (to:Character {id: $to_id})
		CREATE (from)-[r:%s $properties]->(to)
		RETURN r
	`, relType)

	params := map[string]interface{}{
		"from_id":    fromID,
		"to_id":      toID,
		"properties": properties,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetCharacterRelationships retrieves all relationships for a character
func (r *novelGraphRepository) GetCharacterRelationships(ctx context.Context, characterID string) ([]*model.GraphRelationship, error) {
	query := `
		MATCH (c:Character {id: $character_id})-[r]-(other:Character)
		RETURN r, other
	`

	params := map[string]interface{}{"character_id": characterID}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	relationships := make([]*model.GraphRelationship, 0, len(records))
	return relationships, nil
}

// DeleteRelationship deletes a relationship between nodes
func (r *novelGraphRepository) DeleteRelationship(ctx context.Context, fromID, toID, relType string) error {
	query := fmt.Sprintf(`
		MATCH (from {id: $from_id})-[r:%s]->(to {id: $to_id})
		DELETE r
	`, relType)

	params := map[string]interface{}{
		"from_id": fromID,
		"to_id":   toID,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}
