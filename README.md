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
- ✅ REST API endpoints
- ✅ Unit tests (94.9% coverage on domain layer)

## Installation

```bash
go mod tidy
```

## Usage

### Run as HTTP Server

```bash
go run cmd/main.go -server
# Or with custom host/port
go run cmd/main.go -server -host 0.0.0.0 -port 3000
```

### Run Demo Mode

```bash
go run cmd/main.go
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

- [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)
- [stretchr/testify](https://github.com/stretchr/testify)

## License

MIT
