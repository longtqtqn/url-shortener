# URL Shortener

[GitHub Repository](https://github.com/longtqtqn/url-shortener.git)

A simple, modular, and production-ready URL shortener service built with Go, Gin, Bun ORM, and PostgreSQL.

## Features
- **Simple API Key Authentication**: Every request requires an API key
- **Easy Onboarding**: Users without keys can create one with just their email
- **Flexible Short Codes**: Auto-generated or custom user-defined short codes
- **Link Management**: Create, view, and delete short URLs
- **Click Tracking**: Monitor click counts and last clicked times
- **Fair Usage**: All users have the same configurable link limits
- **RESTful API**: Clean, simple endpoints with Gin
- **PostgreSQL Storage**: Reliable data persistence
- **Docker Ready**: Easy setup with Docker Compose
- **Migration Based**: Schema managed via SQL migrations

## Architecture
- **Gin**: HTTP server and routing
- **Bun ORM**: Database access
- **PostgreSQL**: Data storage
- **Fx**: Dependency injection and lifecycle management
- **Copier**: Struct mapping between DB and domain models
- **Migrations**: Schema managed via SQL files in the `migrations/` folder

## API Overview

The service provides a simple, API key-based authentication system:

- **Public Endpoints**: Create API key, resolve short URLs
- **Protected Endpoints**: All link management operations require API key
- **Simple Onboarding**: Users create API keys with just their email
- **Flexible Short Codes**: Auto-generated or custom user-defined codes

### Key Endpoints
- `POST /create-api-key` - Create new API key (no auth required)
- `POST /api/links` - Create short URL (auth required)
- `GET /api/links` - List user's short URLs (auth required)
- `DELETE /api/links/:code` - Delete short URL (auth required)
- `GET /:code` - Resolve short URL (public)



## Project Structure
```
├── backend/                   # Go backend service
│   ├── cmd/app/main.go       # Application entrypoint
│   ├── internal/
│   │   ├── domain/           # Domain models + repository interfaces
│   │   ├── repo/             # Database models and repositories
│   │   ├── usecase/          # Business logic
│   │   ├── transport/        # HTTP handlers, router, middleware
│   │   │   ├── http/
│   │   │   │   ├── handler/  # Link and User handlers
│   │   │   │   └── router/   # Route registration
│   │   │   └── middleware/   # API key authentication
│   │   └── seeder/           # Database seeding utilities
│   ├── migrations/           # SQL migration files
│   └── go.mod, go.sum       # Go dependencies
├── frontend/                  # Frontend application (if any)
├── docker-compose.yml         # Docker setup for Postgres and pgAdmin
└── README.md                 # This file
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

2. Start PostgreSQL services:
   ```sh
   docker-compose up -d
   ```
   This starts:
   - **Primary**: PostgreSQL master on port 5433
   - **Replica**: PostgreSQL slave for read scaling
   - **pgpool**: Connection pooler on port 5432 (app connects here)
   - **pgAdmin**: Database management on port 8081

3. Wait for database to be ready, then run migrations:
   ```sh
   migrate -path ./backend/migrations -database "postgres://user:passhihihi@localhost:5433/urlshortener?sslmode=disable" up
   ```

4. Run the app:
   ```sh
   cd backend
   ENV=development go run ./cmd/app/main.go
   ```

#### **Database Architecture:**
- **Primary (Port 5433)**: Write operations, migrations, admin tasks
- **Replica (Port 5434)**: Read operations for scaling
- **pgpool (Port 5432)**: Smart connection routing (app connects here)
- **pgAdmin (Port 8081)**: Database management interface

### Local Development
- **App Connection**: Your Go app connects to pgpool on port 5432 (which routes to primary/replica)
- **Direct Database Access**: Use port 5433 for direct primary database access (migrations, admin)
- **pgAdmin**: Access database management at [http://localhost:8081](http://localhost:8081) (admin@admin.com / admin)
- Run migrations before starting the app.

## Configuration (Environment Variables)
- `SEED_ENABLED` (default: false): enable seeding on startup.
- `SEED_MODE` (default: `enforce`): seeding behavior.
  - `enforce`: upsert and restore soft-deleted users to match seeder data.
  - `exist-only`: insert only if missing; never update existing rows.
- `FREE_PLAN_MAX_LINKS` (default: 10): maximum number of links per user.
- `DATABASE_URL`: Postgres DSN (required in production).
- `PORT`: HTTP port (required in production).
- `GIN_MODE`: `debug` or `release` (required in production).
- `SEED_USERS_JSON`: optional JSON array to seed users (email and apikey only).

Example `SEED_USERS_JSON` (single line):
```env
SEED_USERS_JSON='[{"email":"test@example.com","apikey":"testkey12345678901234567890123456"},{"email":"demo@example.com","apikey":"demokey1234567890123456789012345"}]'
```

**Note**: API keys should be 32 characters long (hexadecimal format). The seeder will warn you if keys are not the correct length.

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
- **Clean Schema**: Single migration file with complete database setup
- **golang-migrate**: Professional migration management tool
- **Easy Reset**: Use `migrate down` to rollback, `migrate up` to apply

### Quick Database Setup
```bash
# 1. Create database (if it doesn't exist)
createdb -h localhost -p 5433 -U user urlshortener

# 2. Apply migrations
migrate -path ./backend/migrations -database "postgres://user:passhihihi@localhost:5433/urlshortener?sslmode=disable" up

# 3. To reset database (rollback all migrations)
migrate -path ./backend/migrations -database "postgres://user:passhihihi@localhost:5433/urlshortener?sslmode=disable" down
```

### Migration Files
- **`000001_init_schema.up.sql`**: Creates all tables and indexes
- **`000001_init_schema.down.sql`**: Removes all tables (clean rollback)

### Schema Overview
- **users**: Simple user accounts (email only)
- **apikeys**: API key authentication  
- **links**: Shortened URLs with click tracking
- **Indexes**: Optimized for common queries

## API Usage

### Authentication
- All `/api` endpoints require an `X-API-KEY` header.
- API keys are 32-character hexadecimal strings.
- Users can create API keys with just their email via `/create-api-key`.

### Endpoints
#### Create User (No Auth Required)





#### Create API Key (No Auth Required)
```
POST /create-api-key
Body: { "email": "user@example.com" }
Response: { 
  "message": "API key created successfully",
  "api_key": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "user_id": 123
}
```

#### Create Short Link
```
POST /api/links
Headers: X-API-KEY: <your-api-key>
Body: { 
  "long_url": "https://example.com",
  "short_code": "custom"  // Optional - auto-generated if not provided
}
Response: { 
  "shortened_url": "http://localhost:8080/custom",
  "short_code": "custom",
  "long_url": "https://example.com"
}
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

#### Delete Link
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

## Business Rules

### User Management
- **Simple**: Only email required to create account
- **Auto-Creation**: Users created automatically when requesting API key
- **No Plans/Roles**: All users are equal with same link limits
- **One API Key**: Each user can have one API key

### Link Management
- **Link Limits**: Configurable via `FREE_PLAN_MAX_LINKS` (default: 10)
- **Short Codes**: 6-character auto-generated or custom user-defined
- **Global Uniqueness**: Short codes must be unique across all users
- **Conflict Handling**: Returns error if custom code already exists

### Short Code Rules
- **Auto-Generated**: Random 6-character string if not specified
- **Custom Codes**: Users can specify their own short codes
- **Validation**: Custom codes checked for global uniqueness
- **Format**: Alphanumeric characters (0-9, a-z, A-Z)

## Development Notes
- Uses Uber Fx for dependency injection and lifecycle.
- Bun ORM models use soft delete and timestamps.
- All schema changes are managed by SQL migrations.
- API key authentication middleware is in `backend/internal/transport/middleware/apikey.go`.
- Routes are registered centrally in `backend/internal/transport/http/router/router.go`.
- Struct mapping uses [jinzhu/copier](https://github.com/jinzhu/copier).
- See `backend/internal/transport/http/handler/link_http_handler.go` for main API logic.
- See `backend/internal/transport/http/handler/user_http_handler.go` for user management.

## License
MIT
