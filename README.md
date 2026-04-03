# JWT Token Generator - DDD + Clean Architecture

A Go project demonstrating JWT token generation using Domain-Driven Design (DDD) and Clean Architecture principles.

## Project Structure

```
jwt-ddd-clean/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── domain/
│   │   ├── model/              # Domain entities
│   │   │   ├── token.go
│   │   │   └── user.go
│   │   ├── repository/         # Repository interfaces
│   │   │   ├── token_repository.go
│   │   │   └── user_repository.go
│   │   └── service/            # Domain services
│   │       └── token_service.go
│   ├── infrastructure/
│   │   ├── jwt/                # JWT implementation
│   │   │   └── jwt_provider.go
│   │   ├── repository/         # Repository implementations
│   │   │   └── memory_token_repository.go
│   │   └── http/               # HTTP server & handlers
│   │       ├── server.go
│   │       └── token_http_handler.go
│   ├── handler/                # Application handlers
│   │   └── token_handler.go
│   └── dto/                    # Data Transfer Objects
│       └── token_dto.go
├── pkg/
│   └── jwt/                    # Public API
│       └── jwt.go
├── tests/
│   └── postman_collection.json # Postman collection
├── go.mod
└── README.md
```

## Architecture Layers

### Domain Layer (`internal/domain/`)
- **Entities**: Core business objects (Token, User)
- **Repository Interfaces**: Contracts for data access
- **Services**: Business logic implementation

### Infrastructure Layer (`internal/infrastructure/`)
- **JWT Provider**: Concrete JWT implementation using `golang-jwt/jwt/v5`
- **Repository Implementation**: In-memory token storage
- **HTTP Server**: REST API server

### Handler Layer (`internal/handler/`)
- **Token Handler**: Application-level request handling

### DTO Layer (`internal/dto/`)
- **Data Transfer Objects**: Request/Response structures

### Package Layer (`pkg/`)
- **Public API**: Clean interface for external consumers

## Features

- ✅ JWT Access Token generation
- ✅ JWT Refresh Token generation
- ✅ Token validation
- ✅ Token refresh mechanism
- ✅ Token revocation (blacklisting)
- ✅ In-memory token storage
- ✅ Clean Architecture separation
- ✅ DDD principles
- ✅ REST API endpoints with **Gorilla Mux** routing
- ✅ PostgreSQL database integration
- ✅ Auto-migrations on startup
- ✅ Environment-based configuration (.env)
- ✅ Unit tests (94.9% coverage on domain layer)

## Installation

```bash
go mod tidy
```

## Usage

### 1. Setup Configuration

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` file with your PostgreSQL credentials:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=jwt_ddd
```

### 2. Run as HTTP Server with PostgreSQL

```bash
# Make sure PostgreSQL is running
go run cmd/main.go -server

# Or with custom host/port
go run cmd/main.go -server -host 0.0.0.0 -port 3000
```

The application automatically:
- Connects to PostgreSQL
- Runs database migrations on startup
- Starts the HTTP API server

### 3. Run Demo Mode (Show Configuration)

```bash
go run cmd/main.go
```

### Command Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-server` | `false` | Run as HTTP server |
| `-host` | (from .env) | Server host |
| `-port` | (from .env) | Server port |

### PostgreSQL Setup

#### Option 1: Using Docker (Easiest)

```bash
# Start PostgreSQL container
docker run -d \
  --name jwt-postgres \
  -e POSTGRES_DB=jwt_ddd \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# Verify it's running
docker ps | grep jwt-postgres
```

#### Option 2: Local PostgreSQL Installation

**Ubuntu/Debian:**
```bash
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql
```

**macOS:**
```bash
brew install postgresql
brew services start postgresql
```

**Create Database and User:**
```bash
sudo -u postgres psql

CREATE USER jwt_user WITH PASSWORD 'jwt_password';
CREATE DATABASE jwt_ddd OWNER jwt_user;
GRANT ALL PRIVILEGES ON DATABASE jwt_ddd TO jwt_user;
\c jwt_ddd
GRANT ALL ON SCHEMA public TO jwt_user;
\q
```

Then update your `.env` file:
```env
DB_USER=jwt_user
DB_PASSWORD=jwt_password
DB_NAME=jwt_ddd
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | API info |
| GET | `/api/health` | Health check |
| POST | `/api/token/generate` | Generate new tokens |
| POST | `/api/token/refresh` | Refresh access token |
| POST | `/api/token/validate` | Validate token |
| POST | `/api/token/revoke` | Revoke token |

### API Examples

#### Generate Token
```bash
curl -X POST http://localhost:8080/api/token/generate \
  -H "Content-Type: application/json" \
  -d '{"username":"john_doe","password":"password123"}'
```

#### Validate Token
```bash
curl -X POST http://localhost:8080/api/token/validate \
  -H "Authorization: Bearer <access_token>"
```

#### Refresh Token
```bash
curl -X POST http://localhost:8080/api/token/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'
```

#### Revoke Token
```bash
curl -X POST http://localhost:8080/api/token/revoke \
  -H "Content-Type: application/json" \
  -d '{"token":"<access_token>"}'
```

## Testing

### Run Unit Tests

```bash
# Run all tests
go test -v ./...

# Run domain layer tests
go test -v ./internal/domain/service/...

# Run handler tests
go test -v ./internal/handler/...

# With coverage
go test -cover ./...
```

### Postman Collection

Import `postman_collection.json` into Postman to test the API.

**Steps:**
1. Open Postman
2. Click **Import**
3. Select `postman_collection.json`
4. Start the server: `go run cmd/main.go -server`
5. Run the requests in the collection

The collection includes:
- Individual endpoint tests
- Automated test scripts
- Complete token lifecycle workflow
- Pre-configured environment variables

## Configuration

```go
config := jwt.Config{
    SecretKey:       "your-secret-key",
    Issuer:          "your-app",
    AccessTokenTTL:  15 * time.Minute,
    RefreshTokenTTL: 7 * 24 * time.Hour,
}
```

## Dependencies

- [gorilla/mux](https://github.com/gorilla/mux) - HTTP router
- [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) - JWT implementation
- [lib/pq](https://github.com/lib/pq) - PostgreSQL driver
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver
- [joho/godotenv](https://github.com/joho/godotenv) - .env file loader
- [stretchr/testify](https://github.com/stretchr/testify) - Testing toolkit
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate) - Database migrations

## License

MIT
