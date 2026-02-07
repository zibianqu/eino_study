-- Migration: Create chat_chunk table for storing chat message records
-- Author: System
-- Date: 2026-02-07

-- Enable pgvector extension if not already enabled
CREATE EXTENSION IF NOT EXISTS vector;

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
