-- +migrate Down
-- Remove additional indexes and constraints

DROP INDEX IF EXISTS idx_inventories_name;
DROP INDEX IF EXISTS idx_inventories_location;

ALTER TABLE inventories DROP CONSTRAINT IF EXISTS chk_inventories_quantity;
ALTER TABLE inventories DROP CONSTRAINT IF EXISTS chk_inventories_price;
