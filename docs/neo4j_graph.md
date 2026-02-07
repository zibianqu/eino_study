# Neo4j Graph Database Integration

Neo4j 图数据库集成文档

## 概述

本项目集成了 Neo4j 图数据库，用于存储和查询文档、实体及其之间的复杂关系。提供了完整的增删改查（CRUD）操作和关系管理功能。

## 数据模型

### 节点类型（Node Types）

#### 1. Document（文档节点）
```go
type DocumentNode struct {
    ID       string    // 文档唯一标识
    DocName  string    // 文档名称
    DocHash  string    // 文档哈希值
    FilePath string    // 文件路径
    FileType string    // 文件类型
    CTime    time.Time // 创建时间
}
```

#### 2. Entity（实体节点）
```go
type EntityNode struct {
    ID          string                 // 实体唯一标识
    EntityType  string                 // 实体类型
    EntityName  string                 // 实体名称
    EntityValue string                 // 实体值
    Metadata    map[string]interface{} // 元数据
    CTime       time.Time              // 创建时间
}
```

#### 3. ChatMessage（聊天消息节点）
```go
type ChatMessageNode struct {
    ID         string                 // 消息唯一标识
    Role       string                 // 角色（user/assistant/system）
    Content    string                 // 消息内容
    ChunkIndex int                    // 消息序号
    Metadata   map[string]interface{} // 元数据
    CTime      time.Time              // 创建时间
}
```

### 关系类型（Relationship Types）

| 关系类型 | 说明 | 示例 |
|---------|------|------|
| `CONTAINS` | 文档包含实体 | Document -[CONTAINS]-> Entity |
| `REFERENCES` | 文档引用文档 | Document -[REFERENCES]-> Document |
| `SIMILAR_TO` | 文档相似度 | Document -[SIMILAR_TO]- Document |
| `MENTIONED_IN` | 实体被提及 | Entity -[MENTIONED_IN]-> ChatMessage |
| `RELATED_TO` | 通用关联 | Node -[RELATED_TO]-> Node |
| `DERIVED_FROM` | 派生关系 | Entity -[DERIVED_FROM]-> Document |
| `PART_OF` | 包含关系 | Entity -[PART_OF]-> Entity |

## Repository 层 API

### 1. DocumentGraphRepository（文档图谱仓储）

#### 创建文档节点
```go
docRepo := repository.NewDocumentGraphRepository()

doc := &model.DocumentNode{
    ID:       "doc_123",
    DocName:  "示例文档.pdf",
    DocHash:  "abc123",
    FilePath: "/path/to/doc.pdf",
    FileType: "pdf",
    CTime:    time.Now(),
}

err := docRepo.Create(ctx, doc)
```

#### 查询文档
```go
// 按ID查询
doc, err := docRepo.GetByID(ctx, "doc_123")

// 分页列表
docs, err := docRepo.List(ctx, 20, 0) // limit=20, offset=0

// 查找相似文档
similarDocs, err := docRepo.FindSimilar(ctx, "doc_123", 5)

// 获取文档关联的实体
entities, err := docRepo.GetRelatedEntities(ctx, "doc_123")
```

#### 更新文档
```go
doc.DocName = "更新后的名称"
err := docRepo.Update(ctx, doc)
```

#### 删除文档
```go
// 删除节点及其所有关系
err := docRepo.Delete(ctx, "doc_123")
```

### 2. EntityGraphRepository（实体图谱仓储）

#### 创建实体节点
```go
entityRepo := repository.NewEntityGraphRepository()

entity := &model.EntityNode{
    ID:          "entity_456",
    EntityType:  "PERSON",
    EntityName:  "张三",
    EntityValue: "中国科学家",
    Metadata: map[string]interface{}{
        "source": "doc_123",
        "confidence": 0.95,
    },
    CTime: time.Now(),
}

err := entityRepo.Create(ctx, entity)
```

#### 查询实体
```go
// 按ID查询
entity, err := entityRepo.GetByID(ctx, "entity_456")

// 分页列表
entities, err := entityRepo.List(ctx, 20, 0)

// 按类型查询
persons, err := entityRepo.FindByType(ctx, "PERSON", 10)

// 按名称模糊查询
entities, err := entityRepo.FindByName(ctx, "张")

// 获取关联的文档
docs, err := entityRepo.GetRelatedDocuments(ctx, "entity_456")
```

#### 更新实体
```go
entity.EntityValue = "更新后的描述"
err := entityRepo.Update(ctx, entity)
```

#### 删除实体
```go
err := entityRepo.Delete(ctx, "entity_456")
```

### 3. RelationshipGraphRepository（关系仓储）

#### 创建通用关系
```go
relRepo := repository.NewRelationshipGraphRepository()

rel := &model.Relationship{
    Type:       string(model.RelContains),
    FromNodeID: "doc_123",
    ToNodeID:   "entity_456",
    Properties: map[string]interface{}{
        "confidence": 0.95,
        "position": 100,
    },
    CreatedAt: time.Now(),
}

err := relRepo.Create(ctx, rel)
```

#### 创建特定关系

**文档包含实体**
```go
err := relRepo.CreateDocumentContainsEntity(
    ctx,
    "doc_123",      // 文档ID
    "entity_456",   // 实体ID
    map[string]interface{}{"page": 5}, // 属性
)
```

**文档相似度**
```go
err := relRepo.CreateDocumentSimilarity(
    ctx,
    "doc_123",  // 文档1 ID
    "doc_789",  // 文档2 ID
    0.85,       // 相似度分数
)
```

**文档引用**
```go
err := relRepo.CreateDocumentReference(
    ctx,
    "doc_123", // 源文档ID
    "doc_789", // 目标文档ID
)
```

#### 查询关系
```go
// 按ID查询关系
rel, err := relRepo.GetByID(ctx, "rel_id")

// 查询两个节点之间的所有关系
rels, err := relRepo.GetRelationshipsBetween(ctx, "doc_123", "entity_456")

// 查询节点的出边关系
outgoingRels, err := relRepo.GetOutgoingRelationships(
    ctx,
    "doc_123",    // 节点ID
    "CONTAINS",   // 关系类型（可选，空字符串表示所有类型）
)

// 查询节点的入边关系
incomingRels, err := relRepo.GetIncomingRelationships(
    ctx,
    "entity_456", // 节点ID
    "",           // 所有类型
)
```

#### 删除关系
```go
err := relRepo.Delete(ctx, "rel_id")
```

## 使用示例

### 示例1：构建文档-实体知识图谱

```go
package main

import (
    "context"
    "time"
    "github.com/zibianqu/eino_study/internal/app/repository"
    "github.com/zibianqu/eino_study/internal/model"
)

func BuildKnowledgeGraph(ctx context.Context) error {
    // 1. 创建 repositories
    docRepo := repository.NewDocumentGraphRepository()
    entityRepo := repository.NewEntityGraphRepository()
    relRepo := repository.NewRelationshipGraphRepository()

    // 2. 创建文档节点
    doc := &model.DocumentNode{
        ID:       "doc_ai_paper",
        DocName:  "人工智能综述.pdf",
        FilePath: "/papers/ai_survey.pdf",
        FileType: "pdf",
        CTime:    time.Now(),
    }
    if err := docRepo.Create(ctx, doc); err != nil {
        return err
    }

    // 3. 创建实体节点
    entities := []*model.EntityNode{
        {
            ID:          "entity_neural_network",
            EntityType:  "CONCEPT",
            EntityName:  "神经网络",
            EntityValue: "深度学习的基础架构",
            CTime:       time.Now(),
        },
        {
            ID:          "entity_transformer",
            EntityType:  "CONCEPT",
            EntityName:  "Transformer",
            EntityValue: "注意力机制模型",
            CTime:       time.Now(),
        },
    }

    for _, entity := range entities {
        if err := entityRepo.Create(ctx, entity); err != nil {
            return err
        }

        // 4. 创建文档-实体关系
        err := relRepo.CreateDocumentContainsEntity(
            ctx,
            doc.ID,
            entity.ID,
            map[string]interface{}{
                "mention_count": 10,
            },
        )
        if err != nil {
            return err
        }
    }

    return nil
}
```

### 示例2：查询文档关联的知识网络

```go
func QueryDocumentNetwork(ctx context.Context, docID string) error {
    docRepo := repository.NewDocumentGraphRepository()
    relRepo := repository.NewRelationshipGraphRepository()

    // 1. 获取文档信息
    doc, err := docRepo.GetByID(ctx, docID)
    if err != nil {
        return err
    }
    fmt.Printf("文档: %s\n", doc.DocName)

    // 2. 获取文档包含的实体
    entities, err := docRepo.GetRelatedEntities(ctx, docID)
    if err != nil {
        return err
    }
    fmt.Printf("包含 %d 个实体\n", len(entities))

    // 3. 查找相似文档
    similarDocs, err := docRepo.FindSimilar(ctx, docID, 5)
    if err != nil {
        return err
    }
    fmt.Printf("找到 %d 个相似文档\n", len(similarDocs))

    // 4. 获取所有出边关系
    outgoingRels, err := relRepo.GetOutgoingRelationships(ctx, docID, "")
    if err != nil {
        return err
    }
    fmt.Printf("有 %d 个出边关系\n", len(outgoingRels))

    return nil
}
```

### 示例3：实体关系网络分析

```go
func AnalyzeEntityNetwork(ctx context.Context, entityID string) error {
    entityRepo := repository.NewEntityGraphRepository()
    relRepo := repository.NewRelationshipGraphRepository()

    // 1. 获取实体信息
    entity, err := entityRepo.GetByID(ctx, entityID)
    if err != nil {
        return err
    }
    fmt.Printf("实体: %s (%s)\n", entity.EntityName, entity.EntityType)

    // 2. 查找关联的文档
    docs, err := entityRepo.GetRelatedDocuments(ctx, entityID)
    if err != nil {
        return err
    }
    fmt.Printf("出现在 %d 个文档中\n", len(docs))

    // 3. 查找同类型的其他实体
    relatedEntities, err := entityRepo.FindByType(ctx, entity.EntityType, 10)
    if err != nil {
        return err
    }
    fmt.Printf("同类型实体: %d 个\n", len(relatedEntities))

    return nil
}
```

## Cypher 查询示例

如果需要执行自定义 Cypher 查询，可以使用底层的数据库连接：

```go
import "github.com/zibianqu/eino_study/internal/pkg/database"

// 执行写操作
result, err := database.ExecuteWrite(ctx, `
    CREATE (d:Document {id: $id, name: $name})
    RETURN d
`, map[string]interface{}{
    "id": "doc_123",
    "name": "示例文档",
})

// 执行读操作
records, err := database.ExecuteRead(ctx, `
    MATCH (d:Document)-[r:CONTAINS]->(e:Entity)
    WHERE d.id = $doc_id
    RETURN e.entity_name as name, e.entity_type as type
    ORDER BY e.entity_name
`, map[string]interface{}{
    "doc_id": "doc_123",
})
```

## 性能优化建议

### 1. 索引和约束

项目已预定义了约束和索引，在数据库初始化时自动创建：

```go
// 初始化时调用
ctx := context.Background()
err := database.CreateConstraints(ctx)
err = database.CreateIndexes(ctx)
```

### 2. 批量操作

对于大量数据的导入，建议使用批量操作：

```go
// 批量创建实体
for _, entity := range entities {
    if err := entityRepo.Create(ctx, entity); err != nil {
        log.Printf("Failed to create entity: %v", err)
        continue
    }
}
```

### 3. 分页查询

对于大结果集，始终使用分页：

```go
limit := 100
offset := 0

for {
    entities, err := entityRepo.List(ctx, limit, offset)
    if err != nil {
        return err
    }
    
    if len(entities) == 0 {
        break
    }
    
    // 处理数据
    processEntities(entities)
    
    offset += limit
}
```

### 4. 连接池配置

在配置文件中调整连接池大小：

```yaml
neo4j:
  uri: bolt://localhost:7687
  username: neo4j
  password: password
  max_pool_size: 50
  encrypted: false
```

## 配置说明

### 配置文件示例（config.yaml）

```yaml
neo4j:
  uri: "bolt://localhost:7687"     # Neo4j 连接地址
  username: "neo4j"                 # 用户名
  password: "your_password"         # 密码
  max_pool_size: 50                 # 最大连接池大小
  encrypted: false                  # 是否加密连接
```

### 初始化 Neo4j 连接

在 `cmd/server/main.go` 中初始化：

```go
import (
    "github.com/zibianqu/eino_study/internal/config"
    "github.com/zibianqu/eino_study/internal/pkg/database"
)

func main() {
    // 加载配置
    cfg, err := config.LoadConfig("configs/config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // 初始化 Neo4j
    if err := database.InitNeo4j(&cfg.Neo4j); err != nil {
        log.Fatal(err)
    }
    defer database.CloseNeo4j(context.Background())

    // 创建约束和索引
    ctx := context.Background()
    if err := database.CreateConstraints(ctx); err != nil {
        log.Printf("Warning: failed to create constraints: %v", err)
    }
    if err := database.CreateIndexes(ctx); err != nil {
        log.Printf("Warning: failed to create indexes: %v", err)
    }

    // 继续应用初始化...
}
```

## 故障排查

### 连接失败

```bash
# 检查 Neo4j 是否运行
sudo systemctl status neo4j

# 查看 Neo4j 日志
sudo journalctl -u neo4j -f

# 测试连接
cypher-shell -u neo4j -p password
```

### 性能问题

1. 检查是否创建了索引：
```cypher
SHOW INDEXES
```

2. 分析查询性能：
```cypher
PROFILE MATCH (d:Document)-[r:CONTAINS]->(e:Entity)
WHERE d.id = 'doc_123'
RETURN e
```

3. 查看数据库统计：
```cypher
CALL db.stats.retrieve('GRAPH COUNTS')
```

## 相关资源

- [Neo4j 官方文档](https://neo4j.com/docs/)
- [Neo4j Go Driver](https://neo4j.com/docs/go-manual/current/)
- [Cypher 查询语言](https://neo4j.com/docs/cypher-manual/current/)
- [图数据库最佳实践](https://neo4j.com/developer/guide-data-modeling/)

## 文件清单

| 文件路径 | 说明 |
|---------|------|
| `internal/model/graph.go` | 图数据模型定义 |
| `internal/pkg/database/neo4j.go` | Neo4j 连接层 |
| `internal/app/repository/document_graph_repo.go` | 文档图谱仓储 |
| `internal/app/repository/entity_graph_repo.go` | 实体图谱仓储 |
| `internal/app/repository/relationship_graph_repo.go` | 关系仓储 |
| `docs/neo4j_graph.md` | 本文档 |
