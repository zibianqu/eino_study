# API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Endpoints

### Health Check

#### GET /health

Check the health status of the service.

**Response:**
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

---

### Document Management

#### POST /documents

Upload a new document.

**Request Body:**
```json
{
  "file_path": "/path/to/document.txt",
  "doc_name": "My Document"
}
```

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "doc_id": "abc123...",
    "doc_name": "My Document",
    "doc_hash": "def456...",
    "file_path": "/path/to/document.txt",
    "file_type": ".txt",
    "sync_rag_state": 0,
    "sync_entity_state": 0,
    "ctime": "2026-02-05T12:00:00Z"
  }
}
```

#### GET /documents

List all documents with pagination.

**Query Parameters:**
- `page` (optional): Page number, default 1
- `per_page` (optional): Items per page, default 20, max 100

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "doc_id": "abc123...",
      "doc_name": "My Document",
      "file_path": "/path/to/document.txt",
      "file_type": ".txt",
      "sync_rag_state": 1,
      "sync_entity_state": 1,
      "ctime": "2026-02-05T12:00:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "per_page": 20
}
```

#### GET /documents/:id

Get a specific document by ID.

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "doc_id": "abc123...",
    "doc_name": "My Document",
    "doc_hash": "def456...",
    "file_path": "/path/to/document.txt",
    "file_type": ".txt",
    "sync_rag_state": 1,
    "sync_entity_state": 1,
    "ctime": "2026-02-05T12:00:00Z"
  }
}
```

#### DELETE /documents/:id

Delete a document and all related data.

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "document deleted successfully"
  }
}
```

#### POST /documents/:id/process

Process a document (split into chunks, generate embeddings, extract entities).

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "document processing started"
  }
}
```

---

### Query

#### POST /query

Query the knowledge base using RAG.

**Request Body:**
```json
{
  "query": "What is the main topic of the documents?",
  "top_k": 5,
  "stream": false
}
```

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "answer": "Based on the documents, the main topics are...",
    "sources": [
      {
        "doc_id": "abc123...",
        "doc_name": "My Document",
        "content": "Relevant chunk content...",
        "similarity": 0.95
      }
    ],
    "usage": {
      "prompt_tokens": 150,
      "completion_tokens": 200,
      "total_tokens": 350
    }
  }
}
```

---

## Error Response Format

All error responses follow this format:

```json
{
  "code": -1,
  "message": "error description"
}
```

### HTTP Status Codes

- `200`: Success
- `400`: Bad Request
- `404`: Not Found
- `500`: Internal Server Error