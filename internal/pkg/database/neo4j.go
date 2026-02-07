package database

import (
	"context"
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/zibianqu/eino_study/internal/config"
)

var (
	neoDriver neo4j.DriverWithContext
)

// InitNeo4j initializes Neo4j driver
func InitNeo4j(cfg *config.Neo4jConfig) error {
	if cfg == nil {
		return fmt.Errorf("neo4j config is nil")
	}

	auth := neo4j.BasicAuth(cfg.Username, cfg.Password, "")

	driver, err := neo4j.NewDriverWithContext(
		cfg.URI,
		auth,
		func(config *neo4j.Config) {
			config.MaxConnectionPoolSize = cfg.MaxPoolSize
			if cfg.Encrypted {
				config.Encrypted = true
			}
		},
	)

	if err != nil {
		return fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	// Verify connectivity
	ctx := context.Background()
	if err := driver.VerifyConnectivity(ctx); err != nil {
		return fmt.Errorf("failed to verify neo4j connectivity: %w", err)
	}

	neoDriver = driver
	log.Println("✓ Neo4j connected successfully")

	return nil
}

// GetNeo4jDriver returns the Neo4j driver instance
func GetNeo4jDriver() neo4j.DriverWithContext {
	return neoDriver
}

// CloseNeo4j closes the Neo4j driver connection
func CloseNeo4j(ctx context.Context) error {
	if neoDriver != nil {
		return neoDriver.Close(ctx)
	}
	return nil
}

// ExecuteWrite executes a write transaction
func ExecuteWrite(ctx context.Context, query string, params map[string]interface{}) (interface{}, error) {
	session := neoDriver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return nil, fmt.Errorf("write transaction failed: %w", err)
	}

	return result, nil
}

// ExecuteRead executes a read transaction
func ExecuteRead(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	session := neoDriver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		run, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		var records []map[string]interface{}
		for run.Next(ctx) {
			record := run.Record()
			data := make(map[string]interface{})
			for _, key := range record.Keys {
				data[key] = record.Values[0]
			}
			records = append(records, data)
		}

		if err := run.Err(); err != nil {
			return nil, err
		}

		return records, nil
	})

	if err != nil {
		return nil, fmt.Errorf("read transaction failed: %w", err)
	}

	return result.([]map[string]interface{}), nil
}

// CreateConstraints creates necessary constraints and indexes
func CreateConstraints(ctx context.Context) error {
	constraints := []string{
		// Novel domain
		"CREATE CONSTRAINT novel_id IF NOT EXISTS FOR (n:Novel) REQUIRE n.id IS UNIQUE",
		"CREATE CONSTRAINT character_id IF NOT EXISTS FOR (c:Character) REQUIRE c.id IS UNIQUE",
		"CREATE CONSTRAINT location_id IF NOT EXISTS FOR (l:Location) REQUIRE l.id IS UNIQUE",
		"CREATE CONSTRAINT faction_id IF NOT EXISTS FOR (f:Faction) REQUIRE f.id IS UNIQUE",
		"CREATE CONSTRAINT world_setting_id IF NOT EXISTS FOR (w:WorldSetting) REQUIRE w.id IS UNIQUE",

		// Code domain
		"CREATE CONSTRAINT code_file_id IF NOT EXISTS FOR (cf:CodeFile) REQUIRE cf.id IS UNIQUE",
		"CREATE CONSTRAINT class_id IF NOT EXISTS FOR (c:Class) REQUIRE c.id IS UNIQUE",
		"CREATE CONSTRAINT function_id IF NOT EXISTS FOR (f:Function) REQUIRE f.id IS UNIQUE",
		"CREATE CONSTRAINT package_id IF NOT EXISTS FOR (p:Package) REQUIRE p.id IS UNIQUE",

		// Knowledge domain
		"CREATE CONSTRAINT knowledge_doc_id IF NOT EXISTS FOR (kd:KnowledgeDocument) REQUIRE kd.id IS UNIQUE",
		"CREATE CONSTRAINT topic_id IF NOT EXISTS FOR (t:Topic) REQUIRE t.id IS UNIQUE",
		"CREATE CONSTRAINT concept_id IF NOT EXISTS FOR (c:Concept) REQUIRE c.id IS UNIQUE",
		"CREATE CONSTRAINT entity_id IF NOT EXISTS FOR (e:KnowledgeEntity) REQUIRE e.id IS UNIQUE",
	}

	session := neoDriver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	for _, constraint := range constraints {
		_, err := session.Run(ctx, constraint, nil)
		if err != nil {
			log.Printf("Warning: failed to create constraint: %v", err)
			// Continue with other constraints even if one fails
		}
	}

	log.Println("✓ Neo4j constraints created successfully")
	return nil
}

// CreateIndexes creates indexes for better query performance
func CreateIndexes(ctx context.Context) error {
	indexes := []string{
		// Novel domain indexes
		"CREATE INDEX novel_title IF NOT EXISTS FOR (n:Novel) ON (n.title)",
		"CREATE INDEX character_name IF NOT EXISTS FOR (c:Character) ON (c.name)",
		"CREATE INDEX location_name IF NOT EXISTS FOR (l:Location) ON (l.name)",

		// Code domain indexes
		"CREATE INDEX code_file_path IF NOT EXISTS FOR (cf:CodeFile) ON (cf.file_path)",
		"CREATE INDEX class_name IF NOT EXISTS FOR (c:Class) ON (c.name)",
		"CREATE INDEX function_name IF NOT EXISTS FOR (f:Function) ON (f.name)",

		// Knowledge domain indexes
		"CREATE INDEX knowledge_doc_title IF NOT EXISTS FOR (kd:KnowledgeDocument) ON (kd.title)",
		"CREATE INDEX topic_name IF NOT EXISTS FOR (t:Topic) ON (t.name)",
		"CREATE INDEX concept_name IF NOT EXISTS FOR (c:Concept) ON (c.name)",
	}

	session := neoDriver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	for _, index := range indexes {
		_, err := session.Run(ctx, index, nil)
		if err != nil {
			log.Printf("Warning: failed to create index: %v", err)
		}
	}

	log.Println("✓ Neo4j indexes created successfully")
	return nil
}
