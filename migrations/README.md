# Migrations

Database migrations for the JWT DDD Clean Architecture project.

## Running Migrations

### Automatic (Recommended)

The application automatically runs migrations on startup when using the database flag:

```bash
# Start server with database (auto-runs migrations)
./jwt-app -server -db

# Start server with in-memory repositories (no database, no migrations)
./jwt-app -server -db=false
```

Migrations are executed in alphabetical order based on the filename. The migration files are located in `internal/infrastructure/database/migrations/`.

### Manual Using golang-migrate CLI

Install golang-migrate:
```bash
go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Run migrations:
```bash
# Run all up migrations
migrate -path migrations -database "sqlite://./data/inventory.db" up

# Run down migrations
migrate -path migrations -database "sqlite://./data/inventory.db" down

# Run specific version
migrate -path migrations -database "sqlite://./data/inventory.db" goto 2

# Check migration status
migrate -path migrations -database "sqlite://./data/inventory.db" version
```

### Using the Application

The application automatically runs migrations on startup when using the database.

## Migration Files

| Version | File | Description |
|---------|------|-------------|
| 001 | `001_create_inventories_table` | Creates the main inventories table |
| 002 | `002_create_tokens_table` | Creates tokens and token_blacklist tables |
| 003 | `003_seed_inventory_data` | Seeds sample inventory data |

## Schema

### inventories

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT | Primary key |
| sku | TEXT | Unique SKU code |
| name | TEXT | Item name |
| description | TEXT | Item description |
| quantity | INTEGER | Current stock quantity |
| unit | TEXT | Unit of measurement |
| location | TEXT | Warehouse location |
| min_stock | INTEGER | Minimum stock level |
| max_stock | INTEGER | Maximum stock level |
| price | REAL | Item price |
| created_at | DATETIME | Creation timestamp |
| updated_at | DATETIME | Last update timestamp |

### tokens

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT | Primary key |
| user_id | TEXT | User identifier |
| token | TEXT | Refresh token value |
| token_type | TEXT | Type (refresh) |
| expires_at | DATETIME | Expiration timestamp |
| created_at | DATETIME | Creation timestamp |

### token_blacklist

| Column | Type | Description |
|--------|------|-------------|
| id | TEXT | Primary key |
| token | TEXT | Revoked token value |
| expires_at | DATETIME | Original expiration |
| created_at | DATETIME | Blacklist timestamp |
