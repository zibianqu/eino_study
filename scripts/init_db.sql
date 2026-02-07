-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- 文档管理表
CREATE TABLE IF NOT EXISTS documents (
    doc_id VARCHAR(32) PRIMARY KEY,
    doc_name VARCHAR(255) NOT NULL,
    doc_hash VARCHAR(32) NOT NULL,
    file_path TEXT NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    sync_rag_state INTEGER DEFAULT 0,
    sync_enity_state INTEGER DEFAULT 0,
    ctime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(file_path)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_doc_name ON documents(doc_name);
CREATE INDEX IF NOT EXISTS idx_sync_rag_state ON documents(sync_rag_state);
CREATE INDEX IF NOT EXISTS idx_sync_enity_state ON documents(sync_enity_state);
CREATE INDEX IF NOT EXISTS idx_ctime ON documents(ctime);

-- 添加注释
COMMENT ON TABLE documents IS '文档信息管理表';
COMMENT ON COLUMN documents.doc_id IS '文档ID（文件绝对路径的MD5）';
COMMENT ON COLUMN documents.doc_name IS '文档名称';
COMMENT ON COLUMN documents.doc_hash IS '文档内容的MD5哈希';
COMMENT ON COLUMN documents.file_path IS '文档路径';
COMMENT ON COLUMN documents.file_type IS '文件类型';
COMMENT ON COLUMN documents.sync_rag_state IS '同步向量库状态：0-未同步，1-已同步';
COMMENT ON COLUMN documents.sync_enity_state IS '同步实体库状态：0-未同步，1-已同步';
COMMENT ON COLUMN documents.ctime IS '创建时间';

-- 文档块表（用于RAG）
CREATE TABLE IF NOT EXISTS document_chunks (
    id SERIAL PRIMARY KEY,
    doc_id VARCHAR(32) NOT NULL,
    chunk_index INTEGER NOT NULL,
    content TEXT NOT NULL,
    embedding vector(1536),  -- 需要安装 pgvector 扩展
    metadata JSONB,
    ctime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (doc_id) REFERENCES documents(doc_id) ON DELETE CASCADE,
    UNIQUE(doc_id, chunk_index)
);

-- 为向量搜索创建索引
CREATE INDEX IF NOT EXISTS idx_chunk_embedding ON document_chunks USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
CREATE INDEX IF NOT EXISTS idx_chunk_doc_id ON document_chunks(doc_id);

-- 实体表
CREATE TABLE IF NOT EXISTS entities (
    id SERIAL PRIMARY KEY,
    doc_id VARCHAR(32) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_name VARCHAR(255) NOT NULL,
    entity_value TEXT,
    metadata JSONB,
    ctime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (doc_id) REFERENCES documents(doc_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_entity_doc_id ON entities(doc_id);
CREATE INDEX IF NOT EXISTS idx_entity_type ON entities(entity_type);
CREATE INDEX IF NOT EXISTS idx_entity_name ON entities(entity_name);

-- 聊天消息记录表
CREATE TABLE IF NOT EXISTS chat_chunk (
    id SERIAL PRIMARY KEY,
    role VARCHAR(20) NOT NULL,
    chunk_index INTEGER NOT NULL,
    content TEXT NOT NULL,
    embedding vector(1536),  -- Vector embedding for semantic search (1536 dimensions for OpenAI embeddings)
    metadata JSONB,          -- Store session_id, user_id, conversation_id, etc.
    ctime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引以优化查询性能
CREATE INDEX IF NOT EXISTS idx_chat_chunk_role ON chat_chunk(role);
CREATE INDEX IF NOT EXISTS idx_chat_chunk_index ON chat_chunk(chunk_index);
CREATE INDEX IF NOT EXISTS idx_chat_chunk_ctime ON chat_chunk(ctime);
CREATE INDEX IF NOT EXISTS idx_chat_chunk_metadata ON chat_chunk USING gin(metadata);

-- 为向量搜索创建 IVFFLAT 索引
CREATE INDEX IF NOT EXISTS idx_chat_chunk_embedding ON chat_chunk 
USING ivfflat (embedding vector_cosine_ops) 
WITH (lists = 100);

-- 添加表和字段注释
COMMENT ON TABLE chat_chunk IS '聊天消息记录表，用于存储对话历史和支持语义搜索';
COMMENT ON COLUMN chat_chunk.id IS '主键ID';
COMMENT ON COLUMN chat_chunk.role IS '消息角色：user（用户）、assistant（助手）、system（系统）';
COMMENT ON COLUMN chat_chunk.chunk_index IS '消息在会话中的顺序索引';
COMMENT ON COLUMN chat_chunk.content IS '消息内容';
COMMENT ON COLUMN chat_chunk.embedding IS '消息内容的向量嵌入，用于语义搜索';
COMMENT ON COLUMN chat_chunk.metadata IS 'JSON格式的元数据，可包含session_id、user_id、conversation_id等';
COMMENT ON COLUMN chat_chunk.ctime IS '创建时间';
