# 📝 Environment Variables Documentation

## 📋 Overview

File `.env` dan `.env.example` telah dilengkapi dengan semua konfigurasi yang diperlukan untuk menjalankan sistem POS dengan aman dan optimal.

---

## 🔐 File Structure

### `.env` - Actual Configuration
File ini **TIDAK** boleh di-commit ke Git (ada di `.gitignore`). Berisi konfigurasi aktual untuk environment Anda.

### `.env.example` - Template
File ini **BOLEH** di-commit ke Git. Berisi template dengan nilai default dan dokumentasi lengkap.

---

## 📊 Configuration Categories

### 1. **Application Settings**
```bash
SERVER_HOST=localhost              # Server bind address
SERVER_PORT=8080                   # Server port number
SERVER_READ_TIMEOUT=15s            # HTTP read timeout
SERVER_WRITE_TIMEOUT=15s           # HTTP write timeout
SERVER_IDLE_TIMEOUT=60s            # HTTP idle timeout

APP_ENV=development                # Environment: development/production
APP_VERSION=2.2.0                  # Application version
APP_NAME=Ultimate POS System       # Application name
```

**Security Notes:**
- `APP_ENV=production` untuk production
- Timeouts mencegah DoS attacks
- Version tracking untuk monitoring

---

### 2. **Database Configuration**
```bash
DB_HOST=localhost                  # PostgreSQL host
DB_PORT=5432                       # PostgreSQL port
DB_USER=postgres                   # Database user
DB_PASSWORD=postgres               # Database password (CHANGE THIS!)
DB_NAME=pos_system                 # Database name
DB_SSLMODE=disable                 # SSL mode: disable/require/verify-full
DB_MAX_OPEN_CONNS=25               # Max open connections
DB_MAX_IDLE_CONNS=5                # Max idle connections
DB_CONN_MAX_LIFETIME=5m            # Connection max lifetime
DB_CONN_MAX_IDLE_TIME=15m          # Idle connection timeout
```

**Security Checklist:**
- [ ] Gunakan password yang kuat (min 16 chars)
- [ ] Production: `DB_SSLMODE=require` atau `verify-full`
- [ ] Jangan gunakan user `postgres` di production
- [ ] Batasi akses database hanya dari app server

---

### 3. **JWT Configuration**
```bash
JWT_SECRET=your-super-secret-key-change-in-production-min-32-chars
JWT_ISSUER=pos-system              # Token issuer
JWT_AUDIENCE=pos-client            # Token audience
JWT_ACCESS_TOKEN_TTL=15m           # Access token expiry
JWT_REFRESH_TOKEN_TTL=168h         # Refresh token expiry (7 days)
JWT_MAX_REFRESH_COUNT=10           # Max refresh count before re-login
```

**Security Notes:**
- [ ] `JWT_SECRET` minimal 32 karakter (generate: `openssl rand -base64 64`)
- [ ] Access token TTL short (15-30 menit)
- [ ] Refresh token rotation enabled
- [ ] Max refresh count mencegah token reuse

**Generate Secure JWT Secret:**
```bash
# Linux/macOS
openssl rand -base64 64

# Windows (PowerShell)
[System.Convert]::ToBase64String((1..64 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))
```

---

### 4. **CORS Settings**
```bash
CORS_ENABLED=true
CORS_ALLOWED_ORIGINS=https://pos.yourdomain.com,https://admin.yourdomain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,PATCH,OPTIONS
CORS_ALLOWED_HEADERS=Authorization,Content-Type,X-Request-ID,X-CSRF-Token
CORS_EXPOSED_HEADERS=X-Request-ID,X-RateLimit-Limit,X-RateLimit-Remaining
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=3600
```

**Security Notes:**
- [ ] Jangan gunakan `*` di production
- [ ] Specifik domain saja
- [ ] Enable credentials untuk cookies
- [ ] Development: `http://localhost:3000,http://localhost:8080`

---

### 5. **CSRF Protection**
```bash
CSRF_ENABLED=true
CSRF_COOKIE_NAME=csrf_token
CSRF_HEADER_NAME=X-CSRF-Token
CSRF_TOKEN_LENGTH=32
CSRF_COOKIE_SECURE=true            # true di production (HTTPS)
CSRF_COOKIE_TTL=24h
```

**Security Notes:**
- [ ] `CSRF_COOKIE_SECURE=true` di production (HTTPS only)
- [ ] Token length minimal 32 bytes
- [ ] One-time use tokens

---

### 6. **Rate Limiting**
```bash
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100            # Max requests per window
RATE_LIMIT_WINDOW=1m               # Time window
RATE_LIMIT_LOGIN_REQUESTS=5        # Login attempts
RATE_LIMIT_LOGIN_WINDOW=10m        # Login attempt window
```

**Recommended Values:**
| Environment | Requests | Window | Login Attempts |
|-------------|----------|--------|----------------|
| Development | 1000 | 1m | 20 |
| Staging | 500 | 1m | 10 |
| Production | 100 | 1m | 5 |

---

### 7. **Brute Force Protection**
```bash
BRUTE_FORCE_ENABLED=true
BRUTE_FORCE_MAX_ATTEMPTS=5         # Lock after N attempts
BRUTE_FORCE_LOCKOUT_DURATION=15m   # Lock duration
BRUTE_FORCE_WINDOW_DURATION=10m    # Counting window
```

**Security Notes:**
- [ ] Enable di semua environments
- [ ] Production: Consider IP + username tracking
- [ ] Monitor lockout events

---

### 8. **Request Limits**
```bash
MAX_BODY_SIZE=1048576              # 1MB general requests
LOGIN_MAX_BODY_SIZE=10240          # 10KB for login
UPLOAD_MAX_BODY_SIZE=10485760      # 10MB for uploads
```

**Protection:**
- Mencegah DoS via large payloads
- Fast rejection of oversized requests
- Different limits per endpoint type

---

### 9. **Password Policy**
```bash
PASSWORD_MIN_LENGTH=8
PASSWORD_MAX_LENGTH=128
PASSWORD_REQUIRE_UPPERCASE=true
PASSWORD_REQUIRE_LOWERCASE=true
PASSWORD_REQUIRE_NUMBER=true
PASSWORD_REQUIRE_SPECIAL=true
PASSWORD_BLOCK_COMMON=true         # Block common passwords
```

**Requirements Enforced:**
```
✓ Min 8 characters
✓ Max 128 characters
✓ 1+ uppercase (A-Z)
✓ 1+ lowercase (a-z)
✓ 1+ number (0-9)
✓ 1+ special (!@#$%^&*...)
✗ Not in common password list
```

---

### 10. **Security Headers**
```bash
SECURITY_HEADERS_ENABLED=true
SECURITY_HSTS_MAX_AGE=63072000     # 2 years
SECURITY_HSTS_INCLUDE_SUBDOMAINS=true
SECURITY_HSTS_PRELOAD=true
SECURITY_FRAME_OPTIONS=DENY
SECURITY_CONTENT_TYPE_OPTIONS=nosniff
SECURITY_XSS_PROTECTION=1; mode=block
```

**Headers Added:**
- `Strict-Transport-Security` - Force HTTPS
- `X-Frame-Options: DENY` - Prevent clickjacking
- `X-Content-Type-Options: nosniff` - Prevent MIME sniffing
- `X-XSS-Protection: 1; mode=block` - XSS filter

---

### 11. **Logging Configuration**
```bash
LOG_LEVEL=info                     # debug/info/warn/error
LOG_FORMAT=json                    # json/text
LOG_OUTPUT=stdout                  # stdout/file
LOG_FILE_PATH=logs/app.log
LOG_MAX_SIZE=100MB                 # Max file size
LOG_MAX_AGE=30d                    # Retention
LOG_MAX_BACKUPS=10                 # Max backup files
LOG_COMPRESS=true                  # Compress backups
```

**Recommended:**
| Environment | Log Level | Output |
|-------------|-----------|--------|
| Development | debug | stdout |
| Staging | info | stdout + file |
| Production | warn | file + external |

---

### 12. **Audit Logging**
```bash
AUDIT_ENABLED=true
AUDIT_LOG_LEVEL=all                # all/security/none
AUDIT_LOG_SENSITIVE_DATA=false     # Don't log sensitive data
AUDIT_RETENTION_DAYS=365           # Keep for 1 year
```

**What's Logged:**
- Login/logout events
- CRUD operations
- Security events
- Transaction operations
- Password changes

---

### 13. **Session Management**
```bash
SESSION_TIMEOUT=30m                # Inactive timeout
SESSION_ABSOLUTE_TIMEOUT=8h        # Max session duration
SESSION_CONCURRENT_LIMIT=3         # Max concurrent sessions
SESSION_SECURE_COOKIE=true         # HTTPS only
SESSION_HTTP_ONLY=true             # No JS access
SESSION_SAME_SITE=strict           # CSRF protection
```

**Security Notes:**
- [ ] `SESSION_SECURE_COOKIE=true` di production
- [ ] Limit concurrent sessions
- [ ] Short timeout untuk sensitive apps

---

### 14. **Cache Configuration (Optional)**
```bash
CACHE_ENABLED=false
CACHE_TYPE=memory                  # memory/redis
CACHE_REDIS_HOST=localhost
CACHE_REDIS_PORT=6379
CACHE_REDIS_PASSWORD=
CACHE_REDIS_DB=0
CACHE_DEFAULT_TTL=1h
CACHE_MAX_SIZE=10000
```

**Production Setup:**
```bash
CACHE_ENABLED=true
CACHE_TYPE=redis
CACHE_REDIS_HOST=redis.yourdomain.com
CACHE_REDIS_PORT=6379
CACHE_REDIS_PASSWORD=your-redis-password
CACHE_REDIS_DB=0
```

---

### 15. **Email Configuration (Optional)**
```bash
SMTP_ENABLED=false
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@pos.local
SMTP_FROM_NAME=POS System
SMTP_TLS=true
SMTP_TIMEOUT=10s
```

**Setup Gmail:**
1. Enable 2FA
2. Generate App Password
3. Use app password (not account password)

---

### 16. **File Upload Configuration**
```bash
UPLOAD_ENABLED=true
UPLOAD_MAX_SIZE=10MB
UPLOAD_ALLOWED_TYPES=image/jpeg,image/png,image/gif,application/pdf
UPLOAD_PATH=./uploads
UPLOAD_SANITIZE_FILENAME=true
```

**Security:**
- [ ] Sanitize filenames
- [ ] Validate file types
- [ ] Set size limits
- [ ] Store outside webroot

---

### 17. **Backup Configuration**
```bash
BACKUP_ENABLED=false
BACKUP_PATH=./backups
BACKUP_SCHEDULE=0 2 * * *          # Daily at 2 AM
BACKUP_RETENTION_DAYS=30
BACKUP_DATABASE=true
BACKUP_FILES=false
```

**Production Setup:**
```bash
BACKUP_ENABLED=true
BACKUP_PATH=/var/backups/pos
BACKUP_SCHEDULE=0 2 * * *
BACKUP_RETENTION_DAYS=90
BACKUP_DATABASE=true
BACKUP_FILES=true
```

---

### 18. **API Configuration**
```bash
API_VERSION=v1
API_PREFIX=/api
API_DOCS_ENABLED=true
API_DOCS_PATH=/docs
API_PAGINATION_DEFAULT=20
API_PAGINATION_MAX=100
```

**Production:**
```bash
API_DOCS_ENABLED=false             # Disable in production
```

---

### 19. **Feature Flags**
```bash
FEATURE_MULTI_STORE=false
FEATURE_LOYALTY_POINTS=true
FEATURE_BARCODE=true
FEATURE_RECEIPT_PRINTING=true
FEATURE_EXPORT_CSV=true
FEATURE_AUDIT_LOG=true
FEATURE_INVENTORY_ALERTS=true
FEATURE_CATEGORY=true
FEATURE_CUSTOMER=true
```

**Usage:**
- Enable/disable features without code changes
- Gradual rollout
- A/B testing support
- Environment-specific features

---

### 20. **Development Settings**
```bash
DEBUG=true
DEBUG_SQL=false
DEBUG_HTTP=false
DEBUG_JWT=false
RELOAD_ON_CHANGE=false
SWAGGER_ENABLED=true
```

**Environment Settings:**
| Setting | Development | Production |
|---------|-------------|------------|
| DEBUG | true | false |
| DEBUG_SQL | true | false |
| DEBUG_HTTP | true | false |
| DEBUG_JWT | true | false |
| SWAGGER_ENABLED | true | false |

---

## 🔒 Production Checklist

### **Before Deploying to Production:**

```bash
# 1. Generate secure secrets
openssl rand -base64 64 > jwt_secret.txt
openssl rand -base64 32 > db_password.txt

# 2. Update .env file
# - Set APP_ENV=production
# - Set DEBUG=false
# - Update JWT_SECRET
# - Update DB_PASSWORD
# - Configure CORS_ALLOWED_ORIGINS
# - Enable CSRF_COOKIE_SECURE
# - Enable DB_SSLMODE=require

# 3. Security checks
# - All security features enabled
# - Rate limiting configured
# - Brute force protection enabled
# - Password policy enforced
# - Security headers enabled
# - Audit logging enabled

# 4. Test security
# - Test CORS policy
# - Test CSRF protection
# - Test rate limiting
# - Test brute force protection
# - Test password policy
# - Test security headers

# 5. Monitoring
# - Configure log aggregation
# - Set up alerts
# - Enable metrics
# - Configure backup schedule

# 6. Documentation
# - Update API docs
# - Document custom configs
# - Create runbook
# - Incident response plan
```

---

## 📈 Monitoring & Alerts

### **Metrics to Monitor:**
1. Failed login attempts (brute force detection)
2. Rate limit violations
3. CSRF token mismatches
4. SQL injection attempts
5. XSS attempts
6. JWT validation failures
7. Session timeouts
8. Database connection pool usage
9. Response times
10. Error rates

### **Alerts to Configure:**
- Failed login rate > 10/min
- Rate limit violations > 100/min
- Error rate > 5%
- Response time P95 > 1s
- Database connections > 80%
- Disk space < 20%
- Memory usage > 80%

---

## 🔧 Quick Start

### **Development:**
```bash
# 1. Copy example file
cp .env.example .env

# 2. Start PostgreSQL
docker run -d \
  --name pos-postgres \
  -e POSTGRES_DB=pos_system \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:16-alpine

# 3. Run application
go run cmd/main.go -server

# 4. Test API
curl http://localhost:8081/api/health
```

### **Production:**
```bash
# 1. Generate secrets
export JWT_SECRET=$(openssl rand -base64 64)
export DB_PASSWORD=$(openssl rand -base64 32)

# 2. Create .env with production values
# (See Production Checklist above)

# 3. Build binary
go build -o pos-app ./cmd/main.go

# 4. Run with systemd
# Create /etc/systemd/system/pos.service
# Enable and start service

# 5. Verify
curl https://pos.yourdomain.com/api/health
```

---

## 📚 References

- [12 Factor App - Config](https://12factor.net/config)
- [OWASP Security Cheat Sheet](https://cheatsheetseries.owasp.org/)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [PostgreSQL Security](https://www.postgresql.org/docs/current/security.html)

---

**Last Updated:** April 4, 2026  
**Version:** 2.2.0  
**Total Configurations:** 100+
