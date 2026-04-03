-- +migrate Up
-- Cart items are stored in JSON format in the cart table.
-- In production, you might want separate cart_items table.

CREATE TABLE IF NOT EXISTS carts (
    id VARCHAR(100) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    customer_name VARCHAR(255),
    items JSONB NOT NULL DEFAULT '[]',
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id);
CREATE INDEX IF NOT EXISTS idx_carts_created_at ON carts(created_at);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_no VARCHAR(50) NOT NULL UNIQUE,
    cashier_id UUID NOT NULL REFERENCES users(id),
    cashier_name VARCHAR(255) NOT NULL,
    customer_name VARCHAR(255),
    items JSONB NOT NULL,
    subtotal DECIMAL(15, 2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    discount_percent DECIMAL(5, 2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    tax_percent DECIMAL(5, 2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    payment_method VARCHAR(20) NOT NULL,
    payment_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    change_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_transactions_transaction_no ON transactions(transaction_no);
CREATE INDEX IF NOT EXISTS idx_transactions_cashier_id ON transactions(cashier_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_payment_method ON transactions(payment_method);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC);

-- Add constraint for payment method
ALTER TABLE transactions ADD CONSTRAINT chk_transactions_payment_method 
    CHECK (payment_method IN ('CASH', 'CARD', 'QRIS', 'E_WALLET', 'TRANSFER'));

-- Add constraint for status
ALTER TABLE transactions ADD CONSTRAINT chk_transactions_status 
    CHECK (status IN ('PENDING', 'COMPLETED', 'CANCELLED', 'REFUNDED'));
