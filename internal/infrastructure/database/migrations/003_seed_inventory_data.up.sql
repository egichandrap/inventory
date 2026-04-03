-- +migrate Up
-- Insert sample inventory data
INSERT INTO inventories (id, sku, name, description, quantity, unit, location, min_stock, max_stock, price, created_at, updated_at)
VALUES 
    ('inv-001', 'SKU-001', 'Laptop Dell XPS 15', 'High-performance laptop with Intel i7', 50, 'unit', 'Warehouse A', 10, 100, 1299.99, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('inv-002', 'SKU-002', 'Mouse Logitech MX Master 3', 'Wireless mouse with ergonomic design', 200, 'unit', 'Warehouse A', 20, 500, 99.99, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('inv-003', 'SKU-003', 'Keyboard Mechanical RGB', 'Mechanical keyboard with RGB lighting', 150, 'unit', 'Warehouse B', 15, 300, 149.99, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('inv-004', 'SKU-004', 'Monitor Samsung 27"', '4K UHD monitor with HDR support', 75, 'unit', 'Warehouse A', 5, 150, 399.99, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('inv-005', 'SKU-005', 'USB-C Hub 7-in-1', 'Multi-port USB-C hub with HDMI output', 300, 'unit', 'Warehouse B', 30, 600, 49.99, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- +migrate Down
DELETE FROM inventories WHERE id IN ('inv-001', 'inv-002', 'inv-003', 'inv-004', 'inv-005');
