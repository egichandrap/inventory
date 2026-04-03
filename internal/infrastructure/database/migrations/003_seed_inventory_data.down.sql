-- +migrate Down
DELETE FROM inventories WHERE id IN ('inv-001', 'inv-002', 'inv-003', 'inv-004', 'inv-005');
