#!/bin/bash

# Setup script for eino_study project

set -e

echo "Setting up eino_study project..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "✓ Go version: $(go version)"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Warning: Docker is not installed. You'll need to set up PostgreSQL manually."
else
    echo "✓ Docker is installed"
fi

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download
go mod tidy
echo "✓ Dependencies installed"

# Create necessary directories
echo "Creating directories..."
mkdir -p logs
mkdir -p tmp
mkdir -p build/bin
echo "✓ Directories created"

# Check if config file exists
if [ ! -f "configs/config.yaml" ]; then
    echo "Creating config file from example..."
    cp configs/config.example.yaml configs/config.yaml
    echo "✓ Config file created"
    echo "⚠️  Please edit configs/config.yaml to set your API keys and database credentials"
else
    echo "✓ Config file already exists"
fi

# Start Docker containers if Docker is available
if command -v docker &> /dev/null && command -v docker-compose &> /dev/null; then
    echo "Starting PostgreSQL with Docker..."
    docker-compose -f build/docker-compose.yaml up -d
    echo "✓ PostgreSQL started"
    
    # Wait for PostgreSQL to be ready
    echo "Waiting for PostgreSQL to be ready..."
    sleep 5
    
    # Run migrations
    echo "Running database migrations..."
    PGPASSWORD=postgres psql -h localhost -U postgres -d eino_study -f scripts/init_db.sql 2>/dev/null || echo "⚠️  Database may already be initialized"
    echo "✓ Database setup complete"
fi

echo ""
echo "========================================"
echo "Setup complete! 🎉"
echo "========================================"
echo ""
echo "Next steps:"
echo "1. Edit configs/config.yaml to set your API keys"
echo "2. Run 'make run' to start the server"
echo "3. Visit http://localhost:8080/api/v1/health to check the server"
echo ""
echo "For more information, see the README.md file."
echo ""