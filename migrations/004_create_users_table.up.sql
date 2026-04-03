-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'CASHIER',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- Add constraint for role
ALTER TABLE users ADD CONSTRAINT chk_users_role 
    CHECK (role IN ('SUPER_ADMIN', 'ADMIN', 'CASHIER', 'VIEWER'));

-- Add constraint for status
ALTER TABLE users ADD CONSTRAINT chk_users_status 
    CHECK (status IN ('ACTIVE', 'INACTIVE', 'SUSPENDED'));

-- Insert default super admin user (password: admin123)
INSERT INTO users (id, username, password_hash, email, full_name, role, status)
VALUES (
    gen_random_uuid(),
    'superadmin',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'superadmin@pos.local',
    'Super Administrator',
    'SUPER_ADMIN',
    'ACTIVE'
) ON CONFLICT (username) DO NOTHING;

-- Insert default admin user (password: admin123)
INSERT INTO users (id, username, password_hash, email, full_name, role, status)
VALUES (
    gen_random_uuid(),
    'admin',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'admin@pos.local',
    'Administrator',
    'ADMIN',
    'ACTIVE'
) ON CONFLICT (username) DO NOTHING;

-- Insert default cashier user (password: cashier123)
INSERT INTO users (id, username, password_hash, email, full_name, role, status)
VALUES (
    gen_random_uuid(),
    'cashier',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'cashier@pos.local',
    'Cashier User',
    'CASHIER',
    'ACTIVE'
) ON CONFLICT (username) DO NOTHING;
