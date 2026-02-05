# Deployment Guide

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+ with pgvector extension
- Docker and Docker Compose (optional)

## Local Development Setup

### 1. Clone Repository

```bash
git clone https://github.com/zibianqu/eino_study.git
cd eino_study
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Start PostgreSQL with Docker

```bash
make docker-up
```

This will start PostgreSQL with pgvector extension on port 5432.

### 4. Run Database Migrations

```bash
make migrate
```

Or manually:

```bash
psql -U postgres -d eino_study -f scripts/init_db.sql
```

### 5. Configure Application

Copy the example config and modify as needed:

```bash
cp configs/config.example.yaml configs/config.yaml
```

Edit `configs/config.yaml` to set:
- Database connection details
- LLM API keys
- Embedding model configuration

### 6. Run the Application

```bash
make run
```

The server will start on `http://localhost:8080`

### 7. Verify Installation

Check the health endpoint:

```bash
curl http://localhost:8080/api/v1/health
```

## Production Deployment

### Build Binary

```bash
make build
```

This creates binaries in `build/bin/`:
- `server`: Web server
- `cli`: Command-line tools

### Using Docker

#### Build Docker Image

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs
EXPOSE 8080
CMD ["./server"]
```

```bash
docker build -t eino_study:latest .
docker run -p 8080:8080 eino_study:latest
```

### Environment Variables

You can override config values using environment variables:

```bash
export SERVER_PORT=8080
export DATABASE_HOST=localhost
export DATABASE_PASSWORD=your-password
export EINO_LLM_API_KEY=your-api-key
```

### Database Setup

1. Install PostgreSQL 14+
2. Install pgvector extension:

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

3. Create database:

```sql
CREATE DATABASE eino_study;
```

4. Run migrations:

```bash
psql -U postgres -d eino_study -f scripts/init_db.sql
```

### Systemd Service

Create `/etc/systemd/system/eino-study.service`:

```ini
[Unit]
Description=Eino Study Service
After=network.target postgresql.service

[Service]
Type=simple
User=eino
WorkingDirectory=/opt/eino_study
ExecStart=/opt/eino_study/server
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable eino-study
sudo systemctl start eino-study
sudo systemctl status eino-study
```

### Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Monitoring

### Health Check

Regularly check the health endpoint:

```bash
curl http://localhost:8080/api/v1/health
```

### Logs

Logs are written to:
- `logs/app.log`: Application logs
- `logs/error.log`: Error logs

### Database Backups

```bash
pg_dump -U postgres eino_study > backup_$(date +%Y%m%d).sql
```

## Troubleshooting

### Database Connection Failed

- Check PostgreSQL is running
- Verify connection details in config
- Check firewall rules

### API Key Errors

- Verify API keys in config
- Check API key permissions
- Verify network connectivity to LLM provider

### Vector Search Not Working

- Ensure pgvector extension is installed
- Check if embeddings are generated
- Verify vector dimensions match configuration