# Error Dictionary

Dokumentasi terpusat untuk semua error codes yang digunakan dalam aplikasi JWT DDD Clean Architecture.

## Error Codes

### Validation Errors (HTTP 400)

| Error Code | Constant | Message | Description |
|------------|----------|---------|-------------|
| `ERR_VALIDATION` | `ErrValidationErr` | Validation failed | General validation error |
| `ERR_MISSING_FIELD` | `ErrMissingFieldErr` | Required field is missing | Required field not provided |
| `ERR_INVALID_FIELD` | `ErrInvalidFieldErr` | Invalid field value | Field value is invalid |
| `ERR_INVALID_CREDENTIALS` | `ErrInvalidCredentialsErr` | Invalid username or password | Authentication credentials invalid |

### Authentication Errors (HTTP 401)

| Error Code | Constant | Message | Description |
|------------|----------|---------|-------------|
| `ERR_UNAUTHENTICATED` | `ErrUnauthenticatedErr` | Authentication required | No authentication provided |
| `ERR_INVALID_TOKEN` | `ErrInvalidTokenErr` | Invalid token | Token is invalid or malformed |
| `ERR_EXPIRED_TOKEN` | `ErrExpiredTokenErr` | Token has expired | Token has passed expiration |
| `ERR_REVOKED_TOKEN` | `ErrRevokedTokenErr` | Token has been revoked | Token was manually revoked |

### Authorization Errors (HTTP 403)

| Error Code | Constant | Message | Description |
|------------|----------|---------|-------------|
| `ERR_UNAUTHORIZED` | `ErrUnauthorizedErr` | Not authorized to perform this action | Insufficient permissions |
| `ERR_FORBIDDEN` | `ErrForbiddenErr` | Access forbidden | Access denied |

### Not Found Errors (HTTP 404)

| Error Code | Constant | Message | Description |
|------------|----------|---------|-------------|
| `ERR_NOT_FOUND` | `ErrNotFoundErr` | Resource not found | Generic not found |
| `ERR_USER_NOT_FOUND` | `ErrUserNotFoundErr` | User not found | User does not exist |
| `ERR_TOKEN_NOT_FOUND` | `ErrTokenNotFoundErr` | Token not found | Token does not exist |

### Conflict Errors (HTTP 409)

| Error Code | Constant | Message | Description |
|------------|----------|---------|-------------|
| `ERR_CONFLICT` | `ErrConflictErr` | Resource conflict | Generic conflict |
| `ERR_TOKEN_EXISTS` | `ErrTokenExistsErr` | Token already exists | Token already exists |
| `ERR_USER_EXISTS` | `ErrUserExistsErr` | User already exists | User already exists |

### Internal Errors (HTTP 500)

| Error Code | Constant | Message | Description |
|------------|----------|---------|-------------|
| `ERR_INTERNAL` | `ErrInternalErr` | Internal server error | Generic internal error |
| `ERR_TOKEN_GENERATION` | `ErrTokenGenerationErr` | Failed to generate token | Token generation failed |
| `ERR_TOKEN_STORAGE` | `ErrTokenStorageErr` | Failed to store token | Token storage failed |
| `ERR_DATABASE` | `ErrDatabaseErr` | Database error | Database operation failed |

## Usage Examples

### Creating Basic Errors

```go
// Return predefined error
return apperrors.ErrInvalidTokenErr

// Add details to error
return apperrors.ErrValidationErr.WithDetails("Username must be at least 3 characters")
```

### Wrapping Errors

```go
// Wrap underlying error with context
err := database.Save(user)
if err != nil {
    return apperrors.Wrap(err, apperrors.ErrDatabase, "Failed to save user", 500)
}
```

### Helper Functions

```go
// Create validation error with field details
err := apperrors.NewValidationError("username", "must be at least 3 characters")

// Create not found error
err := apperrors.NewNotFoundError("User", userID)

// Create internal error with context
err := apperrors.NewInternalError("process payment", underlyingErr)
```

### Error Handling in Handlers

```go
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
    result, err := h.service.DoSomething()
    if err != nil {
        var appErr *apperrors.AppError
        if errors.As(err, &appErr) {
            w.WriteHeader(appErr.GetHTTPStatus())
            json.NewEncoder(w).Encode(appErr.ToResponse())
            return
        }
        // Fallback for unknown errors
        w.WriteHeader(http.StatusInternalServerError)
    }
}
```

### Checking Error Types

```go
// Check error code
var appErr *apperrors.AppError
if errors.As(err, &appErr) {
    if appErr.Code == apperrors.ErrInvalidToken {
        // Handle invalid token
    }
    if appErr.Code == apperrors.ErrExpiredToken {
        // Handle expired token
    }
}

// Check with ErrorIs
if errors.Is(err, apperrors.ErrInvalidTokenErr) {
    // Handle invalid token
}
```

## Error Response Format

All errors return a consistent JSON response:

```json
{
  "success": false,
  "error": {
    "code": "ERR_INVALID_TOKEN",
    "message": "Invalid token",
    "details": "Token signature mismatch"
  }
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Always `false` for errors |
| `error.code` | string | Error code (e.g., `ERR_INVALID_TOKEN`) |
| `error.message` | string | Human-readable error message |
| `error.details` | string | Optional additional context |

## Best Practices

1. **Always use predefined errors** from the error dictionary
2. **Add details** when additional context is helpful
3. **Wrap external errors** with appropriate error codes
4. **Return specific errors** from the service layer
5. **Let the HTTP layer handle** status code mapping
6. **Log full error details** internally, return sanitized versions to clients

## Error Flow

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  Infrastructure │ ──> │    Domain Layer  │ ──> │  HTTP Handler   │
│  (JWT, Repo)    │     │   (Service)      │     │  (Controller)   │
└─────────────────┘     └──────────────────┘     └─────────────────┘
        │                       │                        │
        │  apperrors.Err        │  apperrors.Err         │  Map to HTTP
        │  InvalidTokenErr      │  InvalidTokenErr       │  Response
        │                       │                        │
        ▼                       ▼                        ▼
   Return Error           Return Error              JSON Response
```
