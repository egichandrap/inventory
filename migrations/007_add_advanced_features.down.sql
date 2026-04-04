-- +migrate Down

-- Remove indexes
DROP INDEX IF EXISTS idx_transactions_customer_id;
DROP INDEX IF EXISTS idx_transactions_store_id;
DROP INDEX IF EXISTS idx_carts_store_id;
DROP INDEX IF EXISTS idx_inventories_store_id;
DROP INDEX IF EXISTS idx_inventories_category_id;
DROP INDEX IF EXISTS idx_inventories_barcode;
DROP INDEX IF EXISTS idx_inventory_alerts_acknowledged;
DROP INDEX IF EXISTS idx_inventory_alerts_severity;
DROP INDEX IF EXISTS idx_inventory_alerts_type;
DROP INDEX IF EXISTS idx_audit_logs_timestamp;
DROP INDEX IF EXISTS idx_audit_logs_entity_id;
DROP INDEX IF EXISTS idx_audit_logs_entity_type;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_stores_is_active;
DROP INDEX IF EXISTS idx_stores_code;
DROP INDEX IF EXISTS idx_customers_phone;
DROP INDEX IF EXISTS idx_customers_email;
DROP INDEX IF EXISTS idx_categories_is_active;
DROP INDEX IF EXISTS idx_categories_parent_id;
DROP INDEX IF EXISTS idx_categories_slug;

-- Remove columns
ALTER TABLE transactions DROP COLUMN IF EXISTS customer_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS store_id;
ALTER TABLE carts DROP COLUMN IF EXISTS store_id;
ALTER TABLE inventories DROP COLUMN IF EXISTS store_id;
ALTER TABLE inventories DROP COLUMN IF EXISTS category_id;
ALTER TABLE inventories DROP COLUMN IF EXISTS barcode_type;
ALTER TABLE inventories DROP COLUMN IF EXISTS barcode;

-- Drop tables
DROP TABLE IF EXISTS inventory_alerts;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS stores;
DROP TABLE IF EXISTS categories;
