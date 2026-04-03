-- +migrate Up
-- Add additional performance optimizations and constraints

-- Add additional index for inventory search
CREATE INDEX IF NOT EXISTS idx_inventories_name ON inventories(name);
CREATE INDEX IF NOT EXISTS idx_inventories_location ON inventories(location);

-- Add check constraint for inventory
ALTER TABLE inventories ADD CONSTRAINT chk_inventories_quantity
    CHECK (quantity >= 0);

ALTER TABLE inventories ADD CONSTRAINT chk_inventories_price
    CHECK (price >= 0);

-- Add comments to tables
COMMENT ON TABLE users IS 'User accounts with role-based access control';
COMMENT ON TABLE carts IS 'Shopping carts for POS transactions';
COMMENT ON TABLE transactions IS 'Sales transactions with payment details';
