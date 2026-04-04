# 🔒 Security Hardening - Ultimate POS System

## 📋 Overview

Dokumentasi ini menjelaskan semua security enhancements yang telah diimplementasikan untuk melindungi backend POS system dari berbagai serangan.

---

## ✅ Security Features Implemented

### 1. **CORS Middleware (Strict Policies)**
**File:** `internal/http/middleware/cors.go`

#### Fitur:
- ✅ Whitelist-based origin validation
- ✅ Credential support dengan strict origin checking
- ✅ Configurable allowed methods & headers
- ✅ Automatic preflight handling
- ✅ Default deny (empty origin list = reject all)

#### Cara Kerja:
```go
// Konfigurasi CORS
corsConfig := httpmiddleware.CORSConfig{
    AllowedOrigins:   []string{"https://pos.yourdomain.com"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-ID"},
    AllowCredentials: true,
}

// Apply middleware
r.Use(httpmiddleware.CORSMiddleware(corsConfig))
```

#### Protection:
- ❌ Mencegah cross-origin attacks dari domain tidak dikenal
- ❌ Mencegah data theft via CORS misconfiguration
- ✅ Hanya origin yang diizinkan yang bisa akses API

---

### 2. **Security Headers Middleware**
**File:** `internal/http/middleware/security_headers.go`

#### Headers Implemented:
| Header | Value | Protection |
|--------|-------|------------|
| `X-Frame-Options` | `DENY` | Clickjacking prevention |
| `X-Content-Type-Options` | `nosniff` | MIME sniffing prevention |
| `X-XSS-Protection` | `1; mode=block` | XSS filter activation |
| `Referrer-Policy` | `strict-origin-when-cross-origin` | Referrer info control |
| `Content-Security-Policy` | `default-src 'none'` | XSS & injection prevention |
| `Permissions-Policy` | `camera=(), microphone=()` | Feature restriction |
| `Strict-Transport-Security` | `max-age=63072000` | HTTPS enforcement |

#### Automatic:
- Semua response mendapat security headers
- Server header dihapus (information disclosure prevention)

---

### 3. **CSRF Protection**
**File:** `internal/http/middleware/csrf.go`

#### Features:
- ✅ Cryptographically secure token generation
- ✅ Per-request token validation
- ✅ One-time use tokens (replay attack prevention)
- ✅ Automatic token expiry
- ✅ Memory cleanup untuk prevent memory leak

#### How It Works:
```
1. Client requests CSRF token (GET)
2. Server generates token & set cookie
3. Client sends token in header (POST/PUT/DELETE)
4. Server validates token matches cookie
5. Token marked as used (one-time)
```

#### Implementation:
```go
csrf := httpmiddleware.NewCSRFMiddleware(httpmiddleware.DefaultCSRFConfig())

// Get token for forms
token, _ := csrf.GetToken(w, r)

// Validate in middleware
r.Use(csrf.Middleware())
```

#### Protection:
- ❌ Mencegah Cross-Site Request Forgery
- ❌ Mencegah unauthorized state changes
- ✅ Safe methods (GET/HEAD/OPTIONS) skip validation

---

### 4. **Request Body Size Limits**
**File:** `internal/http/middleware/body_limit.go`

#### Limits:
| Endpoint Type | Size Limit | Reason |
|---------------|------------|--------|
| General API | 1 MB | Prevent resource exhaustion |
| Login/Register | 10 KB | Strict limit for auth |
| File Upload | 100 KB | Strict for sensitive ops |

#### Usage:
```go
// Global limit
r.Use(httpmiddleware.MaxBodySizeMiddleware(1 << 20)) // 1MB

// Login endpoint strict limit
loginRouter.Use(httpmiddleware.LoginMaxBodyMiddleware())
```

#### Protection:
- ❌ Mencegah DoS via large payload
- ❌ Mencegah memory exhaustion
- ✅ Fast rejection of oversized requests

---

### 5. **Input Sanitization & Validation**
**File:** `internal/pkg/sanitizer/sanitizer.go`

#### Sanitization Functions:
| Function | Purpose | Protection |
|----------|---------|------------|
| `SanitizeString()` | General string sanitization | XSS, null bytes |
| `SanitizeHTML()` | HTML escaping | XSS prevention |
| `ValidateSQL()` | SQL injection detection | SQL injection |
| `ValidateXSS()` | XSS pattern detection | XSS attacks |
| `ValidatePathTraversal()` | Path traversal detection | Directory traversal |
| `SanitizeFilename()` | File name sanitization | Path traversal |
| `ValidateEmail()` | Email format validation | Injection prevention |
| `ValidateUsername()` | Username format validation | Injection prevention |
| `ValidatePassword()` | Password strength check | Weak passwords |
| `SanitizeSearchQuery()` | Search query cleaning | SQL/XSS injection |

#### SQL Injection Patterns Detected:
```sql
SELECT, INSERT, UPDATE, DELETE, DROP, UNION, ALTER
CREATE, EXEC, DECLARE, SET, TABLE, DATABASE
Comment markers: --, /*, */
Special chars: ;, @, @@
```

#### XSS Patterns Detected:
```html
<script>, javascript:, on*=
<iframe>, <object>, <embed>
<link>, <style>, <img on*=
```

#### Usage:
```go
// Validate input
if !sanitizer.ValidateSQL(input) {
    return errors.New("SQL injection attempt detected")
}

// Sanitize string
clean := sanitizer.SanitizeString(userInput)
```

---

### 6. **Brute Force Protection**
**File:** `internal/http/middleware/brute_force.go`

#### Configuration:
```go
config := httpmiddleware.BruteForceConfig{
    MaxAttempts:     5,              // Lock after 5 failed attempts
    LockoutDuration: 15 * time.Minute, // Lock for 15 minutes
    WindowDuration:  10 * time.Minute, // Count attempts in 10 min window
}
```

#### Features:
- ✅ IP-based attempt tracking
- ✅ Automatic lockout after max attempts
- ✅ Sliding window for attempt counting
- ✅ Automatic cleanup of expired entries
- ✅ Memory-efficient storage

#### How It Works:
```
Attempt 1-4: Track & allow
Attempt 5:   Lock IP for 15 minutes
Attempt 6+:  Reject immediately
After 10 min window: Reset counter
After lockout expires: Allow again
```

#### Protection:
- ❌ Mencegah password guessing
- ❌ Mencegah credential stuffing
- ❌ Mencegah automated attacks
- ✅ Automatic recovery after timeout

---

### 7. **JWT Security Enhancements**
**File:** `internal/infrastructure/jwt/token_config.go`

#### Security Features:
| Feature | Implementation | Benefit |
|---------|----------------|---------|
| Short-lived tokens | 15 min access token | Reduce theft impact |
| Refresh tokens | 7 day expiry | User convenience |
| Token rotation | Max 10 refreshes | Detect token reuse |
| Session binding | Session ID in claims | Session management |
| Unique JTI | Per-token ID | Token tracking |
| Audience validation | Strict audience check | Prevent misuse |
| Issuer validation | Verify token source | Prevent forgery |

#### Token Claims:
```go
type TokenClaims struct {
    UserID       string    // User identifier
    Username     string    // Username
    Role         string    // User role
    TokenType    string    // access/refresh
    RefreshCount int       // Times refreshed
    SessionID    string    // Session binding
    jwt.RegisteredClaims   // Standard claims
}
```

#### Validation:
```go
// Validate all claims
err := jwt.ValidateTokenClaims(claims, "access", config)

// Check if should refresh
if claims.ShouldRefresh(5 * time.Minute) {
    // Auto refresh before expiry
}
```

#### Protection:
- ❌ Mencegah token theft exploitation
- ❌ Mencegah token replay
- ✅ Automatic token expiry
- ✅ Session tracking & validation

---

### 8. **Password Security**
**File:** `internal/pkg/security/password.go`

#### Features:
- ✅ Bcrypt hashing (cost factor configurable)
- ✅ Minimum cost enforcement
- ✅ Secure comparison (timing attack prevention)
- ✅ Password strength validation
- ✅ Common password blacklist

#### Password Requirements:
```
✓ Minimum 8 characters
✓ Maximum 128 characters
✓ At least 1 uppercase letter (A-Z)
✓ At least 1 lowercase letter (a-z)
✓ At least 1 number (0-9)
✓ At least 1 special character (!@#$%^&*...)
✗ Not in common passwords list
```

#### Common Passwords Blocked:
```
password, 123456, 12345678, qwerty, abc123,
monkey, letmein, dragon, superman, football,
baseball, shadow, sunshine, master, etc.
```

#### Usage:
```go
hasher := security.NewPasswordHasher(bcrypt.DefaultCost)

// Hash password
hashed, err := hasher.HashPassword(password)

// Verify password
err = hasher.ComparePassword(hashed, password)
```

#### Protection:
- ❌ Mencegah weak password usage
- ❌ Mencegah rainbow table attacks
- ❌ Mencegah timing attacks
- ✅ Secure password storage

---

### 9. **Validation Middleware**
**File:** `internal/http/middleware/validation.go`

#### Automatic Validations:
| Check | Scope | Protection |
|-------|-------|------------|
| Path traversal | URL path | Directory traversal |
| SQL injection | Query params | SQL injection |
| XSS patterns | Query params | XSS attacks |
| SQL injection | Headers | Header injection |
| XSS patterns | Headers | Header XSS |
| Content-Type | POST/PUT/PATCH | Content type spoofing |
| JSON validity | Request body | Malformed data |
| Unknown fields | Request body | Schema violation |

#### JSON Strict Mode:
```go
decoder.DisallowUnknownFields()
```
Rejects requests with unexpected fields, preventing schema abuse.

#### Helper Functions:
```go
ValidateUUID(uuid)         // Validate UUID format
ValidatePositiveInt(n)     // Check positive integer
ValidateNonNegative(n)     // Check non-negative
ValidateStringLength(s)    // Check string length
ValidateAlphanumeric(s)    // Check alphanumeric
```

---

### 10. **Rate Limiting**
**File:** `internal/http/middleware/rate_limiter.go`

#### Features:
- ✅ IP-based rate limiting
- ✅ Configurable limits & windows
- ✅ Automatic cleanup
- ✅ Memory-efficient storage

#### Configuration:
```go
limiter := httpmiddleware.NewRateLimiter(
    100,                    // 100 requests
    1 * time.Minute,        // per minute
)

r.Use(limiter.RateLimitMiddleware())
```

#### Protection:
- ❌ Mencegah API abuse
- ❌ Mencegah DoS attacks
- ❌ Mencegah resource exhaustion

---

## 🛡️ Security Architecture

### Layered Defense (Defense in Depth)

```
┌─────────────────────────────────────────┐
│  Layer 1: Network Level                 │
│  - Rate Limiting                        │
│  - IP-based Brute Force Protection      │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│  Layer 2: HTTP Level                    │
│  - CORS Validation                      │
│  - Security Headers                     │
│  - Body Size Limits                     │
│  - CSRF Protection                      │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│  Layer 3: Request Level                 │
│  - Input Validation                     │
│  - SQL Injection Detection              │
│  - XSS Detection                        │
│  - Path Traversal Detection             │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│  Layer 4: Authentication Level          │
│  - JWT Validation                       │
│  - Token Expiry Check                   │
│  - Role-based Access                    │
│  - Session Management                   │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│  Layer 5: Application Level             │
│  - Password Strength Check              │
│  - Input Sanitization                   │
│  - Business Logic Validation            │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│  Layer 6: Database Level                │
│  - Parameterized Queries                │
│  - SQL Injection Prevention             │
│  - Access Control                       │
└─────────────────────────────────────────┘
```

---

## 📊 Attack Prevention Matrix

| Attack Type | Prevention Layer | Implementation | Status |
|-------------|------------------|----------------|--------|
| **SQL Injection** | Input Validation, DB Layer | Pattern detection, parameterized queries | ✅ |
| **XSS** | Headers, Input Validation | CSP, X-XSS-Protection, sanitization | ✅ |
| **CSRF** | Middleware | Token validation, one-time tokens | ✅ |
| **Clickjacking** | Headers | X-Frame-Options: DENY | ✅ |
| **MIME Sniffing** | Headers | X-Content-Type-Options | ✅ |
| **Brute Force** | Middleware | IP tracking, automatic lockout | ✅ |
| **DoS (Large Body)** | Middleware | Body size limits | ✅ |
| **CORS Misconfig** | Middleware | Strict origin validation | ✅ |
| **Token Theft** | JWT Config | Short expiry, session binding | ✅ |
| **Weak Passwords** | Validation | Strength check, blacklist | ✅ |
| **Path Traversal** | Validation | Pattern detection | ✅ |
| **Header Injection** | Validation | Header sanitization | ✅ |
| **Resource Exhaustion** | Rate Limiter | Request limiting | ✅ |
| **Information Disclosure** | Headers | Server header removal | ✅ |

---

## 🔧 Implementation Details

### Files Created/Modified:

#### New Files (8):
1. `internal/http/middleware/cors.go` - CORS protection
2. `internal/http/middleware/security_headers.go` - Security headers
3. `internal/http/middleware/csrf.go` - CSRF protection
4. `internal/http/middleware/body_limit.go` - Body size limits
5. `internal/http/middleware/brute_force.go` - Brute force protection
6. `internal/http/middleware/validation.go` - Input validation
7. `internal/pkg/sanitizer/sanitizer.go` - Input sanitization
8. `internal/pkg/security/password.go` - Password hashing

#### Modified Files (2):
1. `internal/infrastructure/http/server.go` - Middleware integration
2. `internal/infrastructure/jwt/token_config.go` - JWT security

---

## 📝 Usage Examples

### 1. Applying Security Middleware

```go
// In server.go
func setupRoutes(r *mux.Router, ...) {
    // Global security middleware
    r.Use(httpmiddleware.SecurityHeadersMiddleware)
    r.Use(httpmiddleware.ValidationMiddleware)
    r.Use(httpmiddleware.MaxBodySizeMiddleware(1 << 20))
    
    // Protected routes
    protectedRouter := r.PathPrefix("/api").Subrouter()
    protectedRouter.Use(authMiddleware.Authenticate)
    protectedRouter.Use(csrfMiddleware.Middleware())
    
    // Setup routes...
}
```

### 2. Password Validation

```go
// Validate password strength
valid, errMsg := sanitizer.ValidatePassword(userInput)
if !valid {
    return errors.New(errMsg)
}

// Hash password
hasher := security.NewPasswordHasher(bcrypt.DefaultCost)
hashed, _ := hasher.HashPassword(password)
```

### 3. Input Sanitization

```go
// Sanitize user input
clean := sanitizer.SanitizeString(request.Name)

// Validate against injection
if !sanitizer.ValidateSQL(clean) {
    return errors.New("Invalid input")
}

// Validate email
if !sanitizer.ValidateEmail(request.Email) {
    return errors.New("Invalid email")
}
```

### 4. Brute Force Protection

```go
// Create middleware
bf := httpmiddleware.NewBruteForceMiddleware(
    httpmiddleware.DefaultBruteForceConfig(),
)

// Apply to login routes
loginRouter.Use(bf.Middleware())

// Periodic cleanup
go func() {
    ticker := time.NewTicker(1 * time.Hour)
    for range ticker.C {
        bf.Cleanup()
    }
}()
```

---

## 🎯 Security Best Practices

### ✅ DO:
- Use security middleware on all routes
- Validate all user inputs
- Use parameterized queries
- Implement rate limiting
- Enable security headers
- Use strong password policies
- Log security events
- Rotate JWT secrets regularly
- Update dependencies

### ❌ DON'T:
- Store plaintext passwords
- Trust user input
- Disable security headers
- Use wildcard CORS in production
- Allow unlimited request sizes
- Skip token validation
- Log sensitive data
- Use weak JWT secrets
- Ignore security warnings

---

## 📈 Security Checklist

| Item | Status | Notes |
|------|--------|-------|
| CORS properly configured | ✅ | Whitelist only |
| CSRF protection | ✅ | Token-based |
| Security headers | ✅ | All critical headers |
| Input validation | ✅ | All inputs validated |
| Password hashing | ✅ | Bcrypt with cost |
| Rate limiting | ✅ | Per IP limiting |
| Brute force protection | ✅ | Auto lockout |
| JWT security | ✅ | Short expiry, rotation |
| SQL injection prevention | ✅ | Validation + parameterized |
| XSS prevention | ✅ | Headers + sanitization |
| Body size limits | ✅ | Prevent DoS |
| Error handling | ✅ | No information leakage |
| Logging | ✅ | Audit trail |
| Dependencies updated | ⚠️ | Regular updates needed |

---

## 🚨 Security Monitoring

### Events to Monitor:
1. Failed login attempts (brute force detection)
2. CSRF token mismatches
3. SQL injection attempts
4. XSS attempts
5. Rate limit violations
6. JWT validation failures
7. Unauthorized access attempts
8. Large payload rejections

### Logging:
```go
// Log security events
auditSvc.LogWithSuccess(
    ctx,
    userID, userName,
    model.ActionCreate,
    "SECURITY",
    "LOGIN_FAILED",
    map[string]interface{}{
        "ip": ipAddress,
        "username": attemptedUsername,
        "reason": "invalid_password",
    },
    ipAddress, userAgent,
    false,
    "Failed login attempt",
)
```

---

## 🔮 Future Enhancements

1. **WAF Integration** - Web Application Firewall
2. **IP Reputation** - Block known malicious IPs
3. **Geolocation Blocking** - Restrict by country
4. **Device Fingerprinting** - Detect suspicious devices
5. **Behavioral Analysis** - Detect anomalous patterns
6. **2FA/MFA** - Two-factor authentication
7. **API Key Rotation** - Automatic key rotation
8. **Security Scanning** - Automated vulnerability scanning
9. **Penetration Testing** - Regular security audits
10. **Incident Response** - Automated incident handling

---

## 📚 References

- OWASP Top 10: https://owasp.org/www-project-top-ten/
- JWT Best Practices: https://tools.ietf.org/html/rfc8725
- CORS Specification: https://fetch.spec.whatwg.org/#http-cors-protocol
- Content Security Policy: https://content-security-policy.com/
- Bcrypt: https://en.wikipedia.org/wiki/Bcrypt

---

**Last Updated:** April 4, 2026  
**Version:** 2.2.0 - Security Hardened  
**Build Status:** ✅ Passing  
**Security Status:** 🔒 Production Ready
