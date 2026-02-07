-- 聊天消息记录表
CREATE TABLE IF NOT EXISTS chat_chunk (
    id SERIAL PRIMARY KEY,
    role VARCHAR(20) NOT NULL,
    chunk_index INTEGER NOT NULL,
    content TEXT NOT NULL,
    embedding vector(1536),  -- 需要安装 pgvector 扩展
    metadata JSONB,
    ctime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_chat_chunk_role ON chat_chunk(role);
CREATE INDEX IF NOT EXISTS idx_chat_chunk_index ON chat_chunk(chunk_index);
CREATE INDEX IF NOT EXISTS idx_chat_chunk_ctime ON chat_chunk(ctime);

-- 为向量搜索创建索引（需要 pgvector 扩展）
CREATE INDEX IF NOT EXISTS idx_chat_embedding ON chat_chunk USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

-- 添加注释
COMMENT ON TABLE chat_chunk IS '聊天消息记录表';
COMMENT ON COLUMN chat_chunk.id IS '消息ID';
COMMENT ON COLUMN chat_chunk.role IS '消息角色：user-用户, assistant-助手, system-系统';
COMMENT ON COLUMN chat_chunk.chunk_index IS '消息索引，用于排序';
COMMENT ON COLUMN chat_chunk.content IS '消息内容';
COMMENT ON COLUMN chat_chunk.embedding IS '消息的向量嵌入（1536维）';
COMMENT ON COLUMN chat_chunk.metadata IS '消息元数据（JSON格式），可包含 session_id, user_id 等信息';
COMMENT ON COLUMN chat_chunk.ctime IS '创建时间';
