# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Luban is a Go-based chat application that provides a web interface for conversing with AI models. The application uses a clean architecture with分层分层 structure, featuring MySQL for data persistence, Redis for caching, and MinIO for file storage.

## Common Commands

### Building and Running
```bash
# Build the application
go build -o luban ./cmd

# Run the server directly
go run ./cmd/server.go

# Run with custom config
go run ./cmd/server.go --config ./config.yaml

# Using the CLI
./luban server --config ./config.yaml
```

### Development Setup
```bash
# Start dependencies with Docker Compose
docker-compose up -d

# The application will be available at:
# - HTTP: http://localhost:8080
# - MinIO Console: http://localhost:9001
```

### Testing
```bash
# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test file
go test ./internal/handler
```

## Architecture

### Core Structure
- **cmd/**: Entry points and CLI commands (root.go, server.go)
- **internal/**: Private application code
  - **handler/**: HTTP handlers and route registration
  - **model/**: Data models and business logic
  - **config/**: Configuration management
  - **cache/**: Redis caching implementation
  - **storage/**: MinIO file storage
  - **logger/**: Logging configuration
- **web/**: Frontend static files

### Key Components

#### AppDeps Pattern
Dependencies are injected via `AppDeps` struct in `internal/handler/app.go`:
```go
type AppDeps struct {
    Cfg   *config.Config
    Log   *zap.Logger
    DB    *gorm.DB
    Redis *redis.Client
    MinIO *storage.MinIO
}
```

#### Handler Registration
Routes are registered in `registerRoutes()` method:
- `/api/chat`: Main chat endpoint (POST)
- `/api/models`: List available models
- `/api/sessions`: Session management
- `/api/upload`: File upload to MinIO

#### Data Models
The application uses GORM models:
- `Session`: Chat sessions with model_id and title
- `Message`: Individual messages with role (user/assistant)
- `User`: User profiles (not fully implemented)
- `Model`: Available AI models
- `Attachment`: File metadata for uploads

### Configuration

Configuration is loaded via Viper with support for:
- YAML files (config.yaml)
- Environment variables (prefixed with `LUBAN_`)
- Default values for missing settings

Key config sections:
- `server`: HTTP server settings
- `mysql`: Database connection
- `redis`: Cache configuration
- `minio`: File storage settings
- `web`: Static file serving

### Current Implementation Notes

1. **Chat Endpoint**: Currently returns mock responses. Needs integration with actual AI model APIs.
2. **Model Management**: Basic model listing exists but no real model adapters implemented.
3. **Authentication**: No user authentication system implemented.
4. **Frontend**: Simple HTML/JS interface without build step.

## Database Schema

The application auto-migrates these tables on startup:
- sessions (id, title, model_id, created_at, updated_at)
- messages (id, session_id, role, content, created_at)
- users (id, name, avatar_url, created_at, updated_at)
- models (id, name, created_at, updated_at)
- attachments (id, bucket, object_key, url, type, created_at)

## Important Patterns

1. **Dependency Injection**: All dependencies are injected, not created locally.
2. **Graceful Shutdown**: Uses context with timeout for clean shutdown.
3. **Structured Logging**: Uses Zap logger with structured fields.
4. **Configuration Binding**: Uses Viper for config with mapstructure binding.