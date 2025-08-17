# URL Shortener

A simple, modular, and production-ready URL shortener service built with Go, Gin, Bun ORM, and PostgreSQL.

## Features
- Shorten long URLs to unique short codes
- User authentication via API key
- Track click counts and last clicked time
- Soft delete for links
- RESTful API with Gin
- PostgreSQL for persistent storage
- Dockerized for easy setup
- Admin-friendly with pgAdmin

## Architecture
- **Gin**: HTTP server and routing
- **Bun ORM**: Database access and migrations
- **PostgreSQL**: Data storage
- **Fx**: Dependency injection and lifecycle management
- **Copier**: Struct mapping between DB and domain models

## Project Structure
```
cmd/app/main.go                # Application entrypoint
internal/domain/               # Domain models
internal/repo/                 # Database models and repositories
internal/transport/http/       # HTTP handlers
internal/transport/middleware/ # Gin middleware (API key auth)
internal/usecase/              # Business logic
internal/seeder/               # DB seeding utilities
docker-compose.yml             # Docker setup for Postgres and pgAdmin
go.mod, go.sum                 # Go dependencies
```

## Getting Started

### Prerequisites
- Go 1.25+
- Docker & Docker Compose

### Quick Start (Docker)
1. Clone the repo:
   ```sh
   git clone <your-repo-url>
   cd url-shortener
   ```
2. Start Postgres and pgAdmin:
   ```sh
   docker-compose up -d
   ```
3. Run the app:
   ```sh
   go run ./cmd/app/main.go
   ```
4. Access pgAdmin at [http://localhost:8081](http://localhost:8081) (user: admin@admin.com, pass: admin)

### Local Development
- Update the Postgres DSN in `cmd/app/main.go` if needed.
- Run the app as above.

## API Usage

### Authentication
- All `/api` endpoints require an `X-API-KEY` header.
- Seeded users and API keys are created on startup (see `internal/seeder/seeder.go`).

### Endpoints

#### Create Short Link
```
POST /api/links
Headers: X-API-KEY: <your-api-key>
Body: { "long_url": "https://example.com" }
Response: { "shortened_url": "http://localhost:8080/abc123" }
```

#### List User Links
```
GET /api/links
Headers: X-API-KEY: <your-api-key>
Response: [
  {
    "shortURL": "http://localhost:8080/abc123",
    "longURL": "https://example.com",
    "clickCount": 0,
    "lastClicked": null,
    "createdAt": "2024-06-01T12:00:00Z"
  }
]
```

#### Redirect Short Link
```
GET /:shortCode
Response: 302 Redirect to original URL
```

## Development Notes
- Uses Uber Fx for dependency injection and lifecycle.
- Bun ORM auto-creates tables and indexes on startup.
- API key authentication middleware is in `internal/transport/middleware/apikey.go`.
- Struct mapping uses [jinzhu/copier](https://github.com/jinzhu/copier).
- See `internal/transport/http/handler/link_http_handler.go` for main API logic.

## License
MIT
