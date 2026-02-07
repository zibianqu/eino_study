# Chat Message Records Feature

聊天消息记录功能文档

## 概述

本功能为 eino_study 项目提供了完整的聊天消息记录和管理能力，支持：

- 消息的创建、查询和删除
- 基于角色（user/assistant/system）的消息过滤
- 向量语义搜索（通过 pgvector）
- 批量消息处理
- 元数据存储（session_id、user_id、conversation_id 等）

## 数据模型

### ChatChunk 结构

```go
type ChatChunk struct {
    ID         int       `json:"id"`          // 主键ID
    Role       string    `json:"role"`        // 消息角色：user, assistant, system
    ChunkIndex int       `json:"chunk_index"` // 消息序号
    Content    string    `json:"content"`     // 消息内容
    Embedding  string    `json:"-"`           // 1536维向量嵌入（用于语义搜索）
    Metadata   string    `json:"metadata"`    // JSON格式元数据
    CTime      time.Time `json:"ctime"`       // 创建时间
}
```

### 数据库表结构

表名：`chat_chunk`

| 字段 | 类型 | 说明 |
|------|------|------|
| id | SERIAL | 主键，自增 |
| role | VARCHAR(20) | 消息角色 |
| chunk_index | INTEGER | 消息在会话中的顺序索引 |
| content | TEXT | 消息内容 |
| embedding | vector(1536) | 向量嵌入 |
| metadata | JSONB | 元数据 |
| ctime | TIMESTAMP | 创建时间 |

### 索引

- `idx_chat_chunk_role`: 角色索引
- `idx_chat_chunk_index`: 顺序索引
- `idx_chat_chunk_ctime`: 时间索引
- `idx_chat_chunk_metadata`: JSONB GIN 索引
- `idx_chat_chunk_embedding`: IVFFLAT 向量索引（用于相似度搜索）

## API 接口

### 1. 创建消息

**POST** `/api/v1/chat/messages`

**请求体：**
```json
{
  "role": "user",
  "content": "What is RAG?",
  "metadata": {
    "session_id": "sess_123",
    "user_id": "user_456"
  }
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "role": "user",
    "chunk_index": 0,
    "content": "What is RAG?",
    "metadata": "{\"session_id\":\"sess_123\",\"user_id\":\"user_456\"}",
    "ctime": "2026-02-07T08:00:00Z"
  }
}
```

### 2. 查询消息列表

**GET** `/api/v1/chat/messages?limit=20&offset=0`

**查询参数：**
- `limit`: 每页数量（默认 20）
- `offset`: 偏移量（默认 0）
- `role`: 过滤角色（可选）

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "role": "user",
      "chunk_index": 0,
      "content": "What is RAG?",
      "metadata": "{\"session_id\":\"sess_123\"}",
      "ctime": "2026-02-07T08:00:00Z"
    }
  ]
}
```

### 3. 语义搜索

**POST** `/api/v1/chat/search`

**请求体：**
```json
{
  "query": "How to use RAG?",
  "top_k": 5,
  "threshold": 0.7
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "role": "assistant",
      "content": "RAG stands for Retrieval-Augmented Generation...",
      "similarity": 0.85
    }
  ]
}
```

### 4. 删除消息

**DELETE** `/api/v1/chat/messages/:id`

**响应：**
```json
{
  "code": 0,
  "message": "Message deleted successfully"
}
```

## 使用示例

### Repository 层使用

```go
import (
    "github.com/zibianqu/eino_study/internal/app/repository"
    "github.com/zibianqu/eino_study/internal/model"
)

// 创建 repository
chatRepo := repository.NewChatRepository(db)

// 创建消息
chunk := &model.ChatChunk{
    Role:       "user",
    ChunkIndex: 0,
    Content:    "Hello, world!",
    Metadata:   `{"session_id": "sess_123"}`,
}
err := chatRepo.Create(chunk)

// 查询消息列表
chunks, err := chatRepo.List(20, 0)

// 按角色查询
userMessages, err := chatRepo.GetByRole("user", 10, 0)

// 语义搜索
embedding := "[0.1, 0.2, ...]" // 1536维向量
similar, err := chatRepo.SearchSimilar(embedding, 5, 0.7)
```

### Service 层使用

```go
import "github.com/zibianqu/eino_study/internal/app/service"

// 通过 service 创建消息
req := &api.CreateChatRequest{
    Role:    "user",
    Content: "What is Eino?",
    Metadata: map[string]interface{}{
        "session_id": "sess_123",
    },
}

chunk, err := chatService.CreateMessage(req)
```

## 数据库迁移

### 运行迁移

```bash
# 使用 psql
psql -U your_user -d your_database -f scripts/migrations/002_create_chat_chunk_table.sql

# 或使用项目的 Makefile（如果有定义）
make migrate
```

### 手动创建表

```sql
-- 启用 pgvector 扩展
CREATE EXTENSION IF NOT EXISTS vector;

-- 创建表
CREATE TABLE chat_chunk (
    id SERIAL PRIMARY KEY,
    role VARCHAR(20) NOT NULL,
    chunk_index INTEGER NOT NULL,
    content TEXT NOT NULL,
    embedding vector(1536),
    metadata JSONB,
    ctime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_chat_chunk_embedding ON chat_chunk 
USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
```

## Metadata 最佳实践

metadata 字段使用 JSONB 格式，建议存储以下信息：

```json
{
  "session_id": "会话ID",
  "conversation_id": "对话ID",
  "user_id": "用户ID",
  "timestamp": "时间戳",
  "source": "消息来源",
  "model": "使用的模型",
  "tokens": 123,
  "custom_fields": {}
}
```

## 向量搜索说明

### 相似度计算

使用 pgvector 的余弦相似度（cosine similarity）：
- 值范围：0 到 1
- 1 表示完全相似
- 0 表示完全不相似

### 阈值建议

- `threshold >= 0.8`: 高度相似
- `0.7 <= threshold < 0.8`: 中等相似
- `threshold < 0.7`: 低相似度

### 性能优化

- IVFFLAT 索引的 `lists` 参数影响查询性能和准确性
- 建议值：`lists = sqrt(total_rows)`
- 数据量大时可调整 `lists` 值以平衡性能

## 注意事项

1. **pgvector 扩展**：确保 PostgreSQL 已安装 pgvector 扩展
2. **向量维度**：embedding 固定为 1536 维（OpenAI 默认维度）
3. **索引创建**：向量索引需要在有数据后创建才能生效
4. **元数据查询**：使用 JSONB 操作符查询 metadata 字段

```sql
-- 查询特定 session_id 的消息
SELECT * FROM chat_chunk 
WHERE metadata @> '{"session_id": "sess_123"}';
```

## 文件清单

| 文件路径 | 说明 |
|---------|------|
| `internal/model/chat.go` | 数据模型定义 |
| `internal/app/repository/chat_repo.go` | 数据访问层 |
| `internal/app/service/chat_service.go` | 业务逻辑层 |
| `internal/app/handler/chat.go` | HTTP 处理器 |
| `pkg/api/chat.go` | API 请求/响应定义 |
| `scripts/migrations/002_create_chat_chunk_table.sql` | 数据库迁移脚本 |
| `docs/chat_feature.md` | 本文档 |

## 后续扩展

可能的功能扩展方向：

1. 添加会话管理（Session/Conversation 表）
2. 支持消息编辑和版本控制
3. 添加消息标签和分类
4. 实现消息统计和分析
5. 支持多模态消息（图片、文件等）
6. 添加消息导出功能

## 相关资源

- [pgvector 文档](https://github.com/pgvector/pgvector)
- [GORM 文档](https://gorm.io/)
- [Gin 框架](https://github.com/gin-gonic/gin)
- [CloudWeGo Eino](https://www.cloudwego.io/docs/eino/)
