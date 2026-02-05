# Architecture Design

## Overview

Eino Study is a knowledge base application built with Go, utilizing the CloudWeGo Eino framework for RAG (Retrieval-Augmented Generation) capabilities.

## Architecture Layers

```
┌─────────────────────────────────────────┐
│          HTTP Layer (Gin)               │
│  ┌──────────┬──────────┬──────────┐    │
│  │ Health   │ Document │  Query   │    │
│  │ Handler  │ Handler  │ Handler  │    │
│  └──────────┴──────────┴──────────┘    │
└─────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│         Service Layer                   │
│  ┌──────────────┬──────────────┐       │
│  │   Document   │     RAG      │       │
│  │   Service    │   Service    │       │
│  └──────────────┴──────────────┘       │
└─────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│      Repository Layer (GORM)            │
│  ┌──────┬────────┬──────────┐          │
│  │ Doc  │ Chunk  │ Entity   │          │
│  │ Repo │ Repo   │ Repo     │          │
│  └──────┴────────┴──────────┘          │
└─────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│         Data Layer                      │
│  ┌─────────────┬──────────────┐        │
│  │ PostgreSQL  │  Eino        │        │
│  │ (pgvector)  │  Components  │        │
│  └─────────────┴──────────────┘        │
└─────────────────────────────────────────┘
```

## Components

### 1. HTTP Layer (Gin)

Handles HTTP requests and responses using the Gin framework.

- **Health Handler**: System health checks
- **Document Handler**: Document CRUD operations
- **Query Handler**: Knowledge base queries

### 2. Service Layer

Contains business logic.

- **Document Service**: Document management, upload, processing
- **RAG Service**: Query processing with context retrieval and LLM generation

### 3. Repository Layer (GORM)

Data access layer using GORM.

- **Document Repository**: Document metadata operations
- **Chunk Repository**: Document chunks and vector similarity search
- **Entity Repository**: Entity extraction results

### 4. Data Layer

#### PostgreSQL with pgvector
- Stores document metadata
- Stores document chunks with embeddings
- Stores extracted entities
- Provides vector similarity search

#### Eino Components
- **Loader**: Load documents from various sources
- **Splitter**: Split documents into chunks
- **Indexer**: Generate embeddings and store in vector DB
- **Retriever**: Retrieve relevant chunks
- **ChatModel**: LLM integration
- **Graph**: Orchestrate RAG workflow

## Data Flow

### Document Upload Flow

```
1. User uploads document via API
2. Document Handler validates request
3. Document Service:
   - Checks file existence
   - Calculates file hash
   - Creates document record
4. Returns document metadata
```

### Document Processing Flow

```
1. User triggers document processing
2. Eino Loader reads document
3. Eino Splitter splits into chunks
4. Eino Indexer generates embeddings
5. Store chunks with embeddings in PostgreSQL
6. Extract entities (optional)
7. Update sync state
```

### Query Flow

```
1. User submits query
2. Generate query embedding
3. Retrieve similar chunks from vector DB
4. Build prompt with context
5. Call LLM through Eino ChatModel
6. Return answer with sources
```

## Database Schema

### documents
- Primary table for document metadata
- Tracks sync states for RAG and entity extraction

### document_chunks
- Stores document chunks with embeddings
- Uses pgvector for similarity search
- Foreign key to documents

### entities
- Stores extracted entities
- Foreign key to documents

## Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL 14+ with pgvector
- **AI Framework**: CloudWeGo Eino
- **Configuration**: Viper
- **Logging**: Zap (planned)

## Future Enhancements

1. **Authentication & Authorization**: JWT-based auth
2. **Async Processing**: Background job queue for document processing
3. **Caching**: Redis for query caching
4. **Monitoring**: Prometheus metrics
5. **Multiple Vector Stores**: Support for Milvus, Qdrant
6. **File Upload**: Direct file upload instead of file path
7. **Streaming Responses**: Stream LLM responses
8. **Advanced RAG**: Hybrid search, re-ranking