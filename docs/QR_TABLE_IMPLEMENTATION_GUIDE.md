# QR Table Ordering System - Implementation Guide

## 📋 Overview

Dokumen ini menjelaskan implementasi lengkap **QR Table Ordering System** yang memungkinkan customer untuk memesan makanan/minuman langsung dari meja mereka dengan scan QR code.

---

## 🏗️ Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│  Presentation Layer (HTTP Handlers)                     │
│  ├── handler/table_handler.go                          │
│  └── handler/guestorder_handler.go                     │
└───────────────────┬─────────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────────┐
│  Application Layer (Usecases + DTOs)                   │
│  ├── usecase/table_usecase.go                          │
│  ├── usecase/guestorder_usecase.go                     │
│  ├── dto/table_dto.go                                  │
│  └── dto/guestorder_dto.go                             │
└───────────────────┬─────────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────────┐
│  Domain Layer (Models + Services + Repository Interfaces)│
│  ├── model/table.go                                     │
│  ├── model/guestorder.go                               │
│  ├── service/table_service.go                          │
│  ├── service/guestorder_service.go                     │
│  ├── service/qrcode_service.go                         │
│  ├── repository/table_repository.go                    │
│  └── repository/guestorder_repository.go               │
└───────────────────┬─────────────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────────────┐
│  Infrastructure Layer (Repository Implementations)      │
│  ├── repository/postgres_table_repository.go           │
│  ├── repository/postgres_guestorder_repository.go      │
│  ├── repository/memory_table_repository.go             │
│  └── repository/memory_guestorder_repository.go        │
└─────────────────────────────────────────────────────────┘
```

---

## 🗄️ Database Schema

### Tables Table

```sql
CREATE TABLE tables (
    id VARCHAR(100) PRIMARY KEY,
    number INTEGER NOT NULL UNIQUE CHECK (number > 0),
    location VARCHAR(20) NOT NULL CHECK (location IN ('INDOOR', 'OUTDOOR', 'VIP', 'PATIO')),
    capacity INTEGER NOT NULL CHECK (capacity BETWEEN 1 AND 50),
    status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE' 
        CHECK (status IN ('AVAILABLE', 'OCCUPIED', 'RESERVED', 'MAINTENANCE')),
    qr_code TEXT,
    qr_generated BOOLEAN DEFAULT FALSE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tables_number ON tables(number);
CREATE INDEX idx_tables_location ON tables(location);
CREATE INDEX idx_tables_status ON tables(status);
```

**Purpose:**
- Menyimpan informasi meja restaurant
- Track status meja (available, occupied, reserved, maintenance)
- Store QR code dalam format base64
- Location untuk grouping (indoor, outdoor, VIP, patio)

### Guest Orders Table

```sql
CREATE TABLE guest_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) NOT NULL UNIQUE,
    table_id VARCHAR(100) NOT NULL REFERENCES tables(id),
    table_number INTEGER NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(50),
    items JSONB NOT NULL DEFAULT '[]',
    subtotal DECIMAL(15, 2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    tax_percent DECIMAL(5, 2) NOT NULL DEFAULT 11,
    discount_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    discount_percent DECIMAL(5, 2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    payment_method VARCHAR(20) NOT NULL DEFAULT 'CASH'
        CHECK (payment_method IN ('CASH', 'CARD', 'QRIS', 'E_WALLET', 'TRANSFER')),
    payment_status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
        CHECK (payment_status IN ('PENDING', 'PAID', 'REFUNDED')),
    payment_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    change_amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('PENDING', 'CONFIRMED', 'PREPARING', 'READY', 'SERVED', 'CANCELLED')),
    notes TEXT,
    session_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_guest_orders_order_number ON guest_orders(order_number);
CREATE INDEX idx_guest_orders_table_id ON guest_orders(table_id);
CREATE INDEX idx_guest_orders_status ON guest_orders(status);
CREATE INDEX idx_guest_orders_created_at ON guest_orders(created_at DESC);
CREATE INDEX idx_guest_orders_session_id ON guest_orders(session_id);
```

**Purpose:**
- Menyimpan order dari customer yang tidak login
- Items disimpan dalam format JSONB untuk fleksibilitas
- Track order status workflow
- Support multiple payment methods
- Session tracking untuk analytics

---

## 📡 API Documentation

### 1. Table Management

#### Create Table

**Endpoint:** `POST /api/tables`

**Auth Required:** Admin/Super Admin

**Request Body:**
```json
{
  "number": 5,
  "location": "INDOOR",
  "capacity": 4,
  "description": "Near window, good for families"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Meja berhasil dibuat",
  "data": {
    "id": "uuid-here",
    "number": 5,
    "location": "INDOOR",
    "capacity": 4,
    "status": "AVAILABLE",
    "qr_code": "",
    "qr_generated": false,
    "description": "Near window, good for families",
    "created_at": "2026-04-04T10:00:00Z",
    "updated_at": "2026-04-04T10:00:00Z"
  }
}
```

#### Generate QR Code

**Endpoint:** `POST /api/tables/{id}/qr`

**Auth Required:** Admin/Super Admin

**Response:**
```json
{
  "success": true,
  "message": "QR code berhasil digenerate",
  "data": {
    "qr_code": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUg..."
  }
}
```

**Notes:**
- QR code berisi URL: `{base_url}/order?table={number}&id={uuid}`
- QR code disimpan di database dalam format base64
- Bisa di-download sebagai PNG untuk di-print

#### List Tables

**Endpoint:** `GET /api/tables?location=INDOOR&status=AVAILABLE&limit=50&offset=0`

**Auth Required:** Admin/Super Admin

**Query Parameters:**
- `location` (optional): Filter by location (INDOOR, OUTDOOR, VIP, PATIO)
- `status` (optional): Filter by status (AVAILABLE, OCCUPIED, RESERVED, MAINTENANCE)
- `min_capacity` (optional): Filter by minimum capacity
- `max_capacity` (optional): Filter by maximum capacity
- `search` (optional): Search by description
- `limit` (optional, default: 50): Items per page
- `offset` (optional, default: 0): Offset for pagination

#### Update Table Status

**Endpoint:** `POST /api/tables/{id}/status`

**Request Body:**
```json
{
  "status": "OCCUPIED"
}
```

**Valid Status Transitions:**
- `AVAILABLE` → `OCCUPIED`, `RESERVED`, `MAINTENANCE`
- `OCCUPIED` → `AVAILABLE`, `MAINTENANCE`
- `RESERVED` → `OCCUPIED`, `AVAILABLE`
- `MAINTENANCE` → `AVAILABLE`

### 2. Guest Ordering (Public - No Auth Required)

#### Create Order

**Endpoint:** `POST /api/guest/orders`

**Request Body:**
```json
{
  "table_id": "table-uuid",
  "customer_name": "John Doe",
  "customer_phone": "081234567890",
  "session_id": "session-uuid-optional"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Order berhasil dibuat",
  "data": {
    "id": "order-uuid",
    "order_number": "ORD-20260404-0001",
    "table_id": "table-uuid",
    "table_number": 5,
    "customer_name": "John Doe",
    "customer_phone": "081234567890",
    "items": [],
    "subtotal": 0,
    "tax_amount": 0,
    "tax_percent": 11,
    "discount_amount": 0,
    "discount_percent": 0,
    "total_amount": 0,
    "payment_method": "CASH",
    "payment_status": "PENDING",
    "payment_amount": 0,
    "change_amount": 0,
    "status": "PENDING",
    "created_at": "2026-04-04T10:00:00Z",
    "updated_at": "2026-04-04T10:00:00Z"
  }
}
```

#### Add Item to Order

**Endpoint:** `POST /api/guest/orders/{id}/items`

**Request Body:**
```json
{
  "product_id": "product-uuid",
  "product_name": "Nasi Goreng Spesial",
  "quantity": 2,
  "unit_price": 35000,
  "notes": "Tidak pedas, extra nasi"
}
```

**Notes:**
- Product diambil dari inventory system yang sudah ada
- Stock validation dilakukan saat add item
- Notes optional untuk special requests

#### Checkout Order

**Endpoint:** `POST /api/guest/orders/{id}/checkout`

**Request Body:**
```json
{
  "payment_method": "CASH",
  "payment_amount": 100000
}
```

**Supported Payment Methods:**
- `CASH` - Pay at counter (implemented)
- `CARD` - Card payment (TODO: Payment gateway)
- `QRIS` - QRIS payment (TODO: Payment gateway)
- `E_WALLET` - E-wallet (TODO: Payment gateway)
- `TRANSFER` - Bank transfer (TODO: Payment gateway)

**What Happens on Checkout:**
1. Stock deducted untuk semua items
2. Order status berubah dari `PENDING` → `CONFIRMED`
3. Payment status menjadi `PAID` (untuk cash) atau `PENDING` (untuk non-cash)
4. Table tetap `OCCUPIED` sampai order `SERVED` atau `CANCELLED`
5. Order muncul di kitchen dashboard

#### Cancel Order

**Endpoint:** `POST /api/guest/orders/{id}/cancel`

**Notes:**
- Bisa dilakukan kapan saja sebelum order `SERVED`
- Stock akan dikembalikan otomatis
- Table status berubah ke `AVAILABLE`

### 3. Order Management (Staff)

#### Get Pending Orders

**Endpoint:** `GET /api/orders/pending`

**Purpose:** Kitchen dashboard untuk melihat order yang perlu diproses

#### Update Order Status

**Endpoint:** `POST /api/orders/{id}/status`

**Request Body:**
```json
{
  "status": "PREPARING"
}
```

**Order Status Workflow:**
```
PENDING → CONFIRMED → PREPARING → READY → SERVED
   ↓
CANCELLED (anytime before SERVED)
```

**Notes:**
- `CONFIRMED`: Order diterima kitchen
- `PREPARING`: Sedang dibuat
- `READY`: Siap diantar ke meja
- `SERVED`: Sudah sampai ke customer
- Ketika `SERVED` atau `CANCELLED`, table otomatis jadi `AVAILABLE`

#### Get Today's Sales

**Endpoint:** `GET /api/reports/sales/today`

**Auth Required:** Admin/Super Admin

**Response:**
```json
{
  "success": true,
  "message": "Penjualan hari ini",
  "data": {
    "total_sales": 1500000,
    "total_orders": 25,
    "total_items": 75,
    "date": "2026-04-04"
  }
}
```

---

## 🔧 Configuration

### Environment Variables

```bash
# QR Code Configuration
QR_BASE_URL=http://localhost:8080
QR_MERCHANT_NAME=POS Restaurant
QR_SIZE=256
QR_LOGO_PATH=/path/to/logo.png  # Optional
```

### QR Code Service Configuration

QR code bisa dikonfigurasi di `server.go`:

```go
qrConfig := service.QRCodeConfig{
    BaseURL:         "https://pos.restaurant.com",
    MerchantName:    "My Restaurant",
    MerchantLogo:    "/path/to/logo.png",  // Optional
    Size:            256,
    ErrorCorrection: qrcode.Medium,
    ForegroundColor: color.Black,
    BackgroundColor: color.White,
}
```

---

## 📝 Usage Examples

### Example 1: Admin Creates Tables

```bash
# 1. Login as admin
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 2. Create table (save token)
curl -X POST http://localhost:8080/api/tables \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "number": 1,
    "location": "INDOOR",
    "capacity": 4,
    "description": "Window seat"
  }'

# 3. Generate QR code
curl -X POST http://localhost:8080/api/tables/{table_id}/qr \
  -H "Authorization: Bearer YOUR_TOKEN"

# 4. Download QR code (response contains base64 image)
# Save to file and print for table placement
```

### Example 2: Customer Orders via QR Scan

```bash
# 1. Customer scans QR code at table 5
#    URL: http://pos.restaurant.com/order?table=5&id=uuid

# 2. Create order (from mobile web app)
curl -X POST http://localhost:8080/api/guest/orders \
  -H "Content-Type: application/json" \
  -d '{
    "table_id": "table-uuid",
    "customer_name": "John Doe",
    "customer_phone": "081234567890"
  }'

# 3. Add items to order
curl -X POST http://localhost:8080/api/guest/orders/{order_id}/items \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "product-uuid-1",
    "product_name": "Nasi Goreng",
    "quantity": 2,
    "unit_price": 35000,
    "notes": "Tidak pedas"
  }'

# 4. Add another item
curl -X POST http://localhost:8080/api/guest/orders/{order_id}/items \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "product-uuid-2",
    "product_name": "Es Teh Manis",
    "quantity": 2,
    "unit_price": 10000
  }'

# 5. Checkout
curl -X POST http://localhost:8080/api/guest/orders/{order_id}/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "CASH",
    "payment_amount": 100000
  }'

# Response includes:
# - total_amount: 99000 (including 11% tax)
# - change_amount: 1000 (kembalian)
# - status: CONFIRMED
```

### Example 3: Kitchen Staff Manages Orders

```bash
# 1. Login as staff
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"cashier","password":"cashier123"}'

# 2. Get pending orders
curl -X GET http://localhost:8080/api/orders/pending \
  -H "Authorization: Bearer STAFF_TOKEN"

# 3. Update order status to PREPARING
curl -X POST http://localhost:8080/api/orders/{order_id}/status \
  -H "Authorization: Bearer STAFF_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "PREPARING"}'

# 4. Update to READY
curl -X POST http://localhost:8080/api/orders/{order_id}/status \
  -H "Authorization: Bearer STAFF_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "READY"}'

# 5. Update to SERVED (order complete, table becomes AVAILABLE)
curl -X POST http://localhost:8080/api/orders/{order_id}/status \
  -H "Authorization: Bearer STAFF_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "SERVED"}'
```

---

## 🚀 Deployment Checklist

### Pre-Deployment

- [ ] PostgreSQL database setup
- [ ] Run migrations (007 & 008)
- [ ] Configure environment variables
- [ ] Set up QR code logo (optional)
- [ ] Update base URL in config

### Post-Deployment

- [ ] Create initial tables via API
- [ ] Generate QR codes for all tables
- [ ] Print & place QR codes on tables
- [ ] Test order flow end-to-end
- [ ] Train staff on kitchen dashboard
- [ ] Monitor first few orders

---

## 🐛 Troubleshooting

### Issue: QR Code Not Generating

**Symptoms:** Error saat generate QR code

**Solutions:**
1. Check `go-qrcode` library installed: `go mod tidy`
2. Verify base URL configuration
3. Check logo path if using custom logo

### Issue: Order Creation Fails

**Symptoms:** Error saat create guest order

**Solutions:**
1. Verify table exists and is not in MAINTENANCE status
2. Check required fields (table_id, customer_name)
3. Review server logs for detailed error

### Issue: Stock Not Deducting

**Symptoms:** Checkout berhasil tapi stock tidak berkurang

**Solutions:**
1. Verify product_id exists in inventory
2. Check stock availability before checkout
3. Review inventory repository connection

---

## 📚 Related Documentation

- [README.md](../README.md) - Main project documentation
- [QR_TABLE_ORDERING_SYSTEM.md](QR_TABLE_ORDERING_SYSTEM.md) - Technical architecture
- [OPTIMIZATION_SUMMARY.md](../OPTIMIZATION_SUMMARY.md) - Performance optimizations
- [CHANGELOG.md](../CHANGELOG.md) - Version history

---

**Version**: 3.1.0  
**Last Updated**: April 4, 2026  
**Status**: Production Ready ✅
