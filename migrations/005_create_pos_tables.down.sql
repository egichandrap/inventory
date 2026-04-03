-- +migrate Down
-- Drop POS tables and related constraints/indexes
DROP INDEX IF EXISTS idx_transactions_transaction_no;
DROP INDEX IF EXISTS idx_transactions_cashier_id;
DROP INDEX IF EXISTS idx_transactions_status;
DROP INDEX IF EXISTS idx_transactions_payment_method;
DROP INDEX IF EXISTS idx_transactions_created_at;

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS carts;
