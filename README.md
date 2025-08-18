# URL Shortener

[GitHub Repository](https://github.com/longtqtqn/url-shortener.git)

A simple, modular, and production-ready URL shortener service built with Go, Gin, Bun ORM, and PostgreSQL.

## Features
- Shorten long URLs to unique short codes
- User authentication via API key
- Track click counts and last clicked time
- Soft delete for links and users
- Timestamps for creation and updates
- RESTful API with Gin
- PostgreSQL for persistent storage
- Dockerized for easy setup
- Admin-friendly with pgAdmin
- Schema managed by SQL migrations

## Architecture
- **Gin**: HTTP server and routing
- **Bun ORM**: Database access
- **PostgreSQL**: Data storage
- **Fx**: Dependency injection and lifecycle management
- **Copier**: Struct mapping between DB and domain models
- **Migrations**: Schema managed via SQL files in the `migrations/` folder

## Project Structure
```
cmd/app/main.go                # Application entrypoint
internal/domain/               # Domain models
internal/repo/                 # Database models and repositories
internal/transport/http/       # HTTP handlers
internal/transport/middleware/ # Gin middleware (API key auth)
internal/usecase/              # Business logic
internal/seeder/               # DB seeding utilities
migrations/                    # SQL migration files (schema management)
docker-compose.yml             # Docker setup for Postgres and pgAdmin
go.mod, go.sum                 # Go dependencies
```

## Getting Started

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) (for DB migrations)

### Quick Start (Docker)
1. Clone the repo:
   ```sh
   git clone https://github.com/longtqtqn/url-shortener.git
   cd url-shortener
   ```
2. Start Postgres and pgAdmin:
   ```sh
   docker-compose up -d
   ```
3. Run database migrations:
   ```sh
   migrate -path ./migrations -database "postgres://user:passhihihi@localhost:5433/urlshortener?sslmode=disable" up
   ```
4. Run the app:
   ```sh
   go run ./cmd/app/main.go
   ```
5. Access pgAdmin at [http://localhost:8081](http://localhost:8081) (user: admin@admin.com, pass: admin)

### Local Development
- Update the Postgres DSN in `cmd/app/main.go` if needed.
- Run migrations before starting the app.

## Configuration (Environment Variables)
- `SEED_ENABLED` (default: false): enable seeding on startup.
- `SEED_MODE` (default: `enforce`): seeding behavior.
  - `enforce`: upsert and restore soft-deleted users to match seeder data.
  - `exist-only`: insert only if missing; never update existing rows.
- `FREE_PLAN_MAX_LINKS` (default: 10): maximum number of links for `free` plan.

Examples (zsh/bash):
```sh
export SEED_ENABLED=true
export SEED_MODE=enforce
export FREE_PLAN_MAX_LINKS=10
```

Environment loading:
- The app auto-loads `.env` and overlays `.env.{ENV}` (default `ENV=development`).
- Create and edit `.env.development` and/or `.env.production` as needed.
- An example file is provided: `.env.example`. Keep it in the repo and use it as a template.

From the example file:
```sh
# Create environment files from the example
cp .env.example .env.development
cp .env.example .env.production   # then edit for production-safe values
```

Run per environment:
```sh
# Development
ENV=development go run ./cmd/app/main.go

# Production (local prod-like run)
ENV=production go run ./cmd/app/main.go
```

## Database Migrations
- All schema changes are managed via SQL files in the `migrations/` folder.
- **Do not rely on the app to create or update tables automatically.**
- Use [golang-migrate](https://github.com/golang-migrate/migrate) or a similar tool to apply migrations.
- Example commands:
  ```sh
  migrate -path ./migrations -database "postgres://user:passhihihi@localhost:5433/urlshortener?sslmode=disable" up
  migrate -path ./migrations -database "postgres://user:passhihihi@localhost:5433/urlshortener?sslmode=disable" down
  ```
- Create a new timestamp-based migration (generates `YYYYMMDDHHMMSS_name.up.sql`):
  ```sh
  migrate create -ext sql -dir ./migrations add_feature
  ```
  Then edit the generated `.up.sql` and `.down.sql` files.

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

#### Soft Delete Link
```
DELETE /api/links/:shortCode
Headers: X-API-KEY: <your-api-key>
Response: 204 No Content
```

#### Redirect Short Link
```
GET /:shortCode
Response: 302 Redirect to original URL
```

## Development Notes
- Uses Uber Fx for dependency injection and lifecycle.
- Bun ORM models use soft delete and timestamps.
- All schema changes are managed by SQL migrations.
- API key authentication middleware is in `internal/transport/middleware/apikey.go`.
- Struct mapping uses [jinzhu/copier](https://github.com/jinzhu/copier).
- See `internal/transport/http/handler/link_http_handler.go` for main API logic.

## License
MIT
