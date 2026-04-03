-- +migrate Up
-- Create inventories table
CREATE TABLE IF NOT EXISTS inventories (
    id TEXT PRIMARY KEY,
    sku TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    quantity INTEGER NOT NULL DEFAULT 0,
    unit TEXT NOT NULL,
    location TEXT,
    min_stock INTEGER NOT NULL DEFAULT 0,
    max_stock INTEGER NOT NULL DEFAULT 0,
    price REAL NOT NULL DEFAULT 0.0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on SKU for faster lookups
CREATE INDEX IF NOT EXISTS idx_inventories_sku ON inventories(sku);

-- Create index on location for filtering
CREATE INDEX IF NOT EXISTS idx_inventories_location ON inventories(location);

-- Create index on quantity for stock filtering
CREATE INDEX IF NOT EXISTS idx_inventories_quantity ON inventories(quantity);

-- +migrate Down
DROP INDEX IF EXISTS idx_inventories_quantity;
DROP INDEX IF EXISTS idx_inventories_location;
DROP INDEX IF EXISTS idx_inventories_sku;
DROP TABLE IF EXISTS inventories;
