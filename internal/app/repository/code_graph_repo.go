package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zibianqu/eino_study/internal/model"
	"github.com/zibianqu/eino_study/internal/pkg/database"
)

// CodeGraphRepository handles code knowledge graph operations
type CodeGraphRepository interface {
	// CodeFile operations
	CreateCodeFile(ctx context.Context, file *model.CodeFileNode) error
	GetCodeFile(ctx context.Context, id string) (*model.CodeFileNode, error)
	UpdateCodeFile(ctx context.Context, file *model.CodeFileNode) error
	DeleteCodeFile(ctx context.Context, id string) error
	ListCodeFilesByProject(ctx context.Context, projectID string) ([]*model.CodeFileNode, error)

	// Class operations
	CreateClass(ctx context.Context, class *model.ClassNode) error
	GetClass(ctx context.Context, id string) (*model.ClassNode, error)
	UpdateClass(ctx context.Context, class *model.ClassNode) error
	DeleteClass(ctx context.Context, id string) error

	// Function operations
	CreateFunction(ctx context.Context, function *model.FunctionNode) error
	GetFunction(ctx context.Context, id string) (*model.FunctionNode, error)
	UpdateFunction(ctx context.Context, function *model.FunctionNode) error
	DeleteFunction(ctx context.Context, id string) error

	// Relationship operations
	CreateInheritance(ctx context.Context, childID, parentID string) error
	CreateFunctionCall(ctx context.Context, callerID, calleeID string) error
	GetClassDependencies(ctx context.Context, classID string) ([]*model.ClassNode, error)
}

type codeGraphRepository struct {
	driver neo4j.DriverWithContext
}

// NewCodeGraphRepository creates a new code graph repository
func NewCodeGraphRepository() CodeGraphRepository {
	return &codeGraphRepository{
		driver: database.GetNeo4jDriver(),
	}
}

// CreateCodeFile creates a code file node
func (r *codeGraphRepository) CreateCodeFile(ctx context.Context, file *model.CodeFileNode) error {
	if file.CreatedAt.IsZero() {
		file.CreatedAt = time.Now()
	}

	query := `
		CREATE (f:CodeFile {
			id: $id,
			project_id: $project_id,
			file_path: $file_path,
			language: $language,
			description: $description,
			created_at: datetime($created_at)
		})
		RETURN f
	`

	params := map[string]interface{}{
		"id":          file.ID,
		"project_id":  file.ProjectID,
		"file_path":   file.FilePath,
		"language":    file.Language,
		"description": file.Description,
		"created_at":  file.CreatedAt.Format(time.RFC3339),
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetCodeFile retrieves a code file by ID
func (r *codeGraphRepository) GetCodeFile(ctx context.Context, id string) (*model.CodeFileNode, error) {
	query := "MATCH (f:CodeFile {id: $id}) RETURN f"
	params := map[string]interface{}{"id": id}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("code file not found")
	}

	file := &model.CodeFileNode{ID: id}
	return file, nil
}

// UpdateCodeFile updates a code file node
func (r *codeGraphRepository) UpdateCodeFile(ctx context.Context, file *model.CodeFileNode) error {
	query := `
		MATCH (f:CodeFile {id: $id})
		SET f.file_path = $file_path,
			f.language = $language,
			f.description = $description,
			f.updated_at = datetime()
		RETURN f
	`

	params := map[string]interface{}{
		"id":          file.ID,
		"file_path":   file.FilePath,
		"language":    file.Language,
		"description": file.Description,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// DeleteCodeFile deletes a code file node
func (r *codeGraphRepository) DeleteCodeFile(ctx context.Context, id string) error {
	query := "MATCH (f:CodeFile {id: $id}) DETACH DELETE f"
	params := map[string]interface{}{"id": id}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// ListCodeFilesByProject lists all code files in a project
func (r *codeGraphRepository) ListCodeFilesByProject(ctx context.Context, projectID string) ([]*model.CodeFileNode, error) {
	query := `
		MATCH (f:CodeFile {project_id: $project_id})
		RETURN f
		ORDER BY f.file_path
	`

	params := map[string]interface{}{"project_id": projectID}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	files := make([]*model.CodeFileNode, 0, len(records))
	return files, nil
}

// CreateClass creates a class node
func (r *codeGraphRepository) CreateClass(ctx context.Context, class *model.ClassNode) error {
	if class.CreatedAt.IsZero() {
		class.CreatedAt = time.Now()
	}

	query := `
		MATCH (f:CodeFile {id: $file_id})
		CREATE (c:Class {
			id: $id,
			file_id: $file_id,
			name: $name,
			type: $type,
			description: $description,
			created_at: datetime($created_at)
		})
		CREATE (f)-[:CONTAINS_CLASS]->(c)
		RETURN c
	`

	params := map[string]interface{}{
		"id":          class.ID,
		"file_id":     class.FileID,
		"name":        class.Name,
		"type":        class.Type,
		"description": class.Description,
		"created_at":  class.CreatedAt.Format(time.RFC3339),
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetClass retrieves a class by ID
func (r *codeGraphRepository) GetClass(ctx context.Context, id string) (*model.ClassNode, error) {
	query := "MATCH (c:Class {id: $id}) RETURN c"
	params := map[string]interface{}{"id": id}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("class not found")
	}

	class := &model.ClassNode{ID: id}
	return class, nil
}

// UpdateClass updates a class node
func (r *codeGraphRepository) UpdateClass(ctx context.Context, class *model.ClassNode) error {
	query := `
		MATCH (c:Class {id: $id})
		SET c.name = $name,
			c.type = $type,
			c.description = $description,
			c.updated_at = datetime()
		RETURN c
	`

	params := map[string]interface{}{
		"id":          class.ID,
		"name":        class.Name,
		"type":        class.Type,
		"description": class.Description,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// DeleteClass deletes a class node
func (r *codeGraphRepository) DeleteClass(ctx context.Context, id string) error {
	query := "MATCH (c:Class {id: $id}) DETACH DELETE c"
	params := map[string]interface{}{"id": id}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// CreateFunction creates a function node
func (r *codeGraphRepository) CreateFunction(ctx context.Context, function *model.FunctionNode) error {
	if function.CreatedAt.IsZero() {
		function.CreatedAt = time.Now()
	}

	query := `
		MATCH (f:CodeFile {id: $file_id})
		CREATE (fn:Function {
			id: $id,
			file_id: $file_id,
			class_id: $class_id,
			name: $name,
			parameters: $parameters,
			return_type: $return_type,
			description: $description,
			complexity: $complexity,
			created_at: datetime($created_at)
		})
		CREATE (f)-[:CONTAINS_FUNCTION]->(fn)
		RETURN fn
	`

	params := map[string]interface{}{
		"id":          function.ID,
		"file_id":     function.FileID,
		"class_id":    function.ClassID,
		"name":        function.Name,
		"parameters":  function.Parameters,
		"return_type": function.ReturnType,
		"description": function.Description,
		"complexity":  function.Complexity,
		"created_at":  function.CreatedAt.Format(time.RFC3339),
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetFunction retrieves a function by ID
func (r *codeGraphRepository) GetFunction(ctx context.Context, id string) (*model.FunctionNode, error) {
	query := "MATCH (fn:Function {id: $id}) RETURN fn"
	params := map[string]interface{}{"id": id}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("function not found")
	}

	function := &model.FunctionNode{ID: id}
	return function, nil
}

// UpdateFunction updates a function node
func (r *codeGraphRepository) UpdateFunction(ctx context.Context, function *model.FunctionNode) error {
	query := `
		MATCH (fn:Function {id: $id})
		SET fn.name = $name,
			fn.parameters = $parameters,
			fn.return_type = $return_type,
			fn.description = $description,
			fn.complexity = $complexity,
			fn.updated_at = datetime()
		RETURN fn
	`

	params := map[string]interface{}{
		"id":          function.ID,
		"name":        function.Name,
		"parameters":  function.Parameters,
		"return_type": function.ReturnType,
		"description": function.Description,
		"complexity":  function.Complexity,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// DeleteFunction deletes a function node
func (r *codeGraphRepository) DeleteFunction(ctx context.Context, id string) error {
	query := "MATCH (fn:Function {id: $id}) DETACH DELETE fn"
	params := map[string]interface{}{"id": id}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// CreateInheritance creates an inheritance relationship between classes
func (r *codeGraphRepository) CreateInheritance(ctx context.Context, childID, parentID string) error {
	query := `
		MATCH (child:Class {id: $child_id})
		MATCH (parent:Class {id: $parent_id})
		CREATE (child)-[:INHERITS]->(parent)
		RETURN child, parent
	`

	params := map[string]interface{}{
		"child_id":  childID,
		"parent_id": parentID,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// CreateFunctionCall creates a function call relationship
func (r *codeGraphRepository) CreateFunctionCall(ctx context.Context, callerID, calleeID string) error {
	query := `
		MATCH (caller:Function {id: $caller_id})
		MATCH (callee:Function {id: $callee_id})
		CREATE (caller)-[:CALLS]->(callee)
		RETURN caller, callee
	`

	params := map[string]interface{}{
		"caller_id": callerID,
		"callee_id": calleeID,
	}

	_, err := database.ExecuteWrite(ctx, query, params)
	return err
}

// GetClassDependencies retrieves all classes that a class depends on
func (r *codeGraphRepository) GetClassDependencies(ctx context.Context, classID string) ([]*model.ClassNode, error) {
	query := `
		MATCH (c:Class {id: $class_id})-[:INHERITS|IMPLEMENTS*1..3]->(dep:Class)
		RETURN DISTINCT dep
	`

	params := map[string]interface{}{"class_id": classID}

	records, err := database.ExecuteRead(ctx, query, params)
	if err != nil {
		return nil, err
	}

	dependencies := make([]*model.ClassNode, 0, len(records))
	return dependencies, nil
}
