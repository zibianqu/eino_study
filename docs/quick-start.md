# Quick Start Guide

## Prerequisites

Before you begin, ensure you have the following installed:

- Go 1.21 or higher
- PostgreSQL 14+ with pgvector extension
- Git
- (Optional) Docker and Docker Compose

## Step 1: Clone the Repository

```bash
git clone https://github.com/zibianqu/eino_study.git
cd eino_study
```

## Step 2: Automatic Setup

Run the setup script:

```bash
chmod +x scripts/setup.sh
./scripts/setup.sh
```

This will:
- Install Go dependencies
- Create necessary directories
- Copy configuration template
- Start PostgreSQL (if Docker is available)
- Run database migrations

## Step 3: Configure the Application

Edit `configs/config.yaml` and set the following:

### Database Configuration

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: eino_study
```

### LLM Configuration (OpenAI)

```yaml
eino:
  llm:
    provider: openai
    api_key: sk-your-api-key-here
    model: gpt-4
    temperature: 0.7
```

### Embedding Configuration

```yaml
eino:
  embedding:
    provider: openai
    api_key: sk-your-api-key-here
    model: text-embedding-3-small
    dimension: 1536
```

## Step 4: Start the Server

```bash
make run
```

Or build and run:

```bash
make build
./build/bin/server
```

The server will start on `http://localhost:8080`

## Step 5: Verify Installation

### Check Health Status

```bash
curl http://localhost:8080/api/v1/health
```

Expected response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "healthy",
    "database": "connected"
  }
}
```

## Step 6: Upload Your First Document

### Create a test document

```bash
echo "This is a test document about artificial intelligence and machine learning." > /tmp/test.txt
```

### Upload the document

```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -d '{
    "file_path": "/tmp/test.txt",
    "doc_name": "Test Document"
  }'
```

Expected response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "doc_id": "abc123...",
    "doc_name": "Test Document",
    "file_path": "/tmp/test.txt",
    "sync_rag_state": 0
  }
}
```

## Step 7: Process the Document

Process the document to generate embeddings:

```bash
curl -X POST http://localhost:8080/api/v1/documents/{doc_id}/process
```

Replace `{doc_id}` with the actual document ID from the previous response.

## Step 8: Query the Knowledge Base

```bash
curl -X POST http://localhost:8080/api/v1/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "What is this document about?",
    "top_k": 3
  }'
```

Expected response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "answer": "This document is about artificial intelligence and machine learning.",
    "sources": [
      {
        "doc_id": "abc123...",
        "doc_name": "Test Document",
        "content": "This is a test document about artificial intelligence..."
      }
    ],
    "usage": {
      "prompt_tokens": 150,
      "completion_tokens": 20,
      "total_tokens": 170
    }
  }
}
```

## Common Commands

### List all documents

```bash
curl http://localhost:8080/api/v1/documents?page=1&per_page=20
```

### Get document details

```bash
curl http://localhost:8080/api/v1/documents/{doc_id}
```

### Delete a document

```bash
curl -X DELETE http://localhost:8080/api/v1/documents/{doc_id}
```

## Development Mode

For development with hot reload:

```bash
make install-tools  # Install air for hot reload
make dev            # Run with hot reload
```

## Troubleshooting

### Database Connection Failed

- Check if PostgreSQL is running: `docker ps`
- Verify database credentials in `configs/config.yaml`
- Check PostgreSQL logs: `docker logs eino_study_postgres`

### API Key Errors

- Ensure your OpenAI API key is valid
- Check API key permissions
- Verify the API key is correctly set in `configs/config.yaml`

### Document Processing Fails

- Ensure the file path is absolute and accessible
- Check if the file type is supported (.txt, .md)
- Review server logs for detailed error messages

### Vector Search Not Working

- Ensure pgvector extension is installed in PostgreSQL
- Check if embeddings are generated (sync_rag_state = 1)
- Verify embedding dimension matches configuration (default 1536)

## Next Steps

- Read the [API Documentation](api.md) for complete API reference
- Check [Architecture Documentation](architecture.md) to understand the system design
- See [Deployment Guide](deployment.md) for production deployment

## Getting Help

If you encounter any issues:

1. Check the logs in `logs/` directory
2. Review the [documentation](../docs/)
3. Create an issue on GitHub with detailed error information

## Resources

- [CloudWeGo Eino Documentation](https://www.cloudwego.io/docs/eino/)
- [PostgreSQL pgvector](https://github.com/pgvector/pgvector)
- [Gin Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)