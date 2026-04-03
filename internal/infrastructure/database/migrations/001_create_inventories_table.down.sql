-- +migrate Down
DROP INDEX IF EXISTS idx_inventories_quantity;
DROP INDEX IF EXISTS idx_inventories_location;
DROP INDEX IF EXISTS idx_inventories_sku;
DROP TABLE IF EXISTS inventories;
