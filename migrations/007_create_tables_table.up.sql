-- +migrate Up

-- Tables table (restaurant tables for QR ordering)
CREATE TABLE IF NOT EXISTS tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number INT NOT NULL UNIQUE,
    location VARCHAR(20) NOT NULL,
    capacity INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE',
    qr_code TEXT NOT NULL DEFAULT '',
    qr_generated BOOLEAN NOT NULL DEFAULT false,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_tables_number ON tables(number);
CREATE INDEX IF NOT EXISTS idx_tables_location ON tables(location);
CREATE INDEX IF NOT EXISTS idx_tables_status ON tables(status);

-- Add constraint for valid location values
ALTER TABLE tables ADD CONSTRAINT chk_tables_location
    CHECK (location IN ('INDOOR', 'OUTDOOR', 'VIP', 'PATIO'));

-- Add constraint for valid status values
ALTER TABLE tables ADD CONSTRAINT chk_tables_status
    CHECK (status IN ('AVAILABLE', 'OCCUPIED', 'RESERVED', 'MAINTENANCE'));

-- Add constraint for valid capacity
ALTER TABLE tables ADD CONSTRAINT chk_tables_capacity
    CHECK (capacity > 0 AND capacity <= 50);
