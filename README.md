# Eino Study - Knowledge Base with RAG

A knowledge base application built with Go, Gin, GORM, and CloudWeGo Eino framework, featuring RAG (Retrieval-Augmented Generation) capabilities.

## Features

- Document management with PostgreSQL
- RAG-based question answering
- Entity extraction from documents
- Vector database integration
- RESTful API with Gin
- Database operations with GORM

## Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- **ORM**: [GORM](https://gorm.io/) - Database toolkit
- **Database**: PostgreSQL with pgvector
- **AI Framework**: [CloudWeGo Eino](https://www.cloudwego.io/docs/eino/) - LLM application framework
- **Configuration**: Viper
- **Logging**: Zap

## Project Structure

```
eino_study/
├── cmd/                    # Application entrypoints
│   ├── server/            # Web server
│   └── cli/               # CLI tools
├── internal/              # Private application code
│   ├── app/              # Application layer
│   ├── eino/             # Eino components
│   ├── model/            # Data models
│   ├── config/           # Configuration
│   └── pkg/              # Internal packages
├── configs/               # Configuration files
├── scripts/               # Scripts and migrations
└── docs/                  # Documentation
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Docker & Docker Compose (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/zibianqu/eino_study.git
cd eino_study
```

2. Install dependencies:
```bash
make tidy
```

3. Start PostgreSQL (using Docker):
```bash
make docker-up
```

4. Run database migrations:
```bash
make migrate
```

5. Copy and configure the config file:
```bash
cp configs/config.example.yaml configs/config.yaml
# Edit configs/config.yaml with your settings
```

6. Run the application:
```bash
make run
```

The server will start at `http://localhost:8080`

## Development

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Run with Hot Reload

```bash
make install-tools  # Install air
make dev
```

## API Documentation

See [docs/api.md](docs/api.md) for detailed API documentation.

### Quick Examples

- Health Check: `GET /api/v1/health`
- Upload Document: `POST /api/v1/documents`
- Query Knowledge Base: `POST /api/v1/query`
- List Documents: `GET /api/v1/documents`

## Configuration

Configuration is managed through YAML files in the `configs/` directory. See `configs/config.example.yaml` for available options.

## License

MIT