-- +migrate Up

-- Guest orders table (orders from tables via QR, no login required)
CREATE TABLE IF NOT EXISTS guest_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) NOT NULL UNIQUE,
    table_id UUID NOT NULL REFERENCES tables(id),
    table_number INT NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(50),
    items JSONB NOT NULL DEFAULT '[]',
    subtotal DECIMAL(15, 2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    tax_percent DECIMAL(5, 2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    discount_percent DECIMAL(5, 2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    payment_method VARCHAR(20) NOT NULL DEFAULT 'CASH',
    payment_status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    payment_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    change_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    notes TEXT,
    session_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_guest_orders_order_number ON guest_orders(order_number);
CREATE INDEX IF NOT EXISTS idx_guest_orders_table_id ON guest_orders(table_id);
CREATE INDEX IF NOT EXISTS idx_guest_orders_status ON guest_orders(status);
CREATE INDEX IF NOT EXISTS idx_guest_orders_created_at ON guest_orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_guest_orders_session_id ON guest_orders(session_id);

-- Add constraint for valid payment method
ALTER TABLE guest_orders ADD CONSTRAINT chk_guest_orders_payment_method
    CHECK (payment_method IN ('CASH', 'CARD', 'QRIS', 'E_WALLET', 'TRANSFER'));

-- Add constraint for valid payment status
ALTER TABLE guest_orders ADD CONSTRAINT chk_guest_orders_payment_status
    CHECK (payment_status IN ('PENDING', 'PAID', 'REFUNDED'));

-- Add constraint for valid order status
ALTER TABLE guest_orders ADD CONSTRAINT chk_guest_orders_status
    CHECK (status IN ('PENDING', 'CONFIRMED', 'PREPARING', 'READY', 'SERVED', 'CANCELLED'));
