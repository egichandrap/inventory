# QR Table Ordering System - Backend Implementation Guide

## 📋 Overview

This document describes the complete backend implementation for the QR Table Ordering System, enabling customers to scan QR codes at tables and place orders directly from their mobile devices.

---

## 🏗️ Architecture

### System Flow

```
Customer scans QR → Opens Web App → Browses Menu → Adds to Cart → Checkout
                                                    ↓
                                    Order created in database
                                                    ↓
                                    Kitchen/Staff Dashboard receives order
                                                    ↓
                                    Staff updates order status
```

### Components Created

```
Domain Layer:
├── model/table.go                    # Table entity
├── model/guestorder.go               # GuestOrder entity  
├── repository/table_repository.go    # Table repository interface
├── repository/guestorder_repository.go # GuestOrder repository interface
└── service/
    ├── qrcode_service.go             # QR code generation with logo
    ├── table_service.go              # Table management logic
    └── guestorder_service.go         # Guest order & checkout logic

Infrastructure Layer:
└── repository/
    ├── postgres_table_repository.go      # Table PostgreSQL implementation
    └── postgres_guestorder_repository.go # GuestOrder PostgreSQL implementation

Application Layer:
└── dto/
    ├── table_dto.go                  # Table request/response DTOs
    └── guestorder_dto.go             # GuestOrder request/response DTOs

Database Migrations:
├── 007_create_tables_table.up.sql
├── 007_create_tables_table.down.sql
├── 008_create_guest_orders_table.up.sql
└── 008_create_guest_orders_table.down.sql
```

---

## 📊 Domain Models

### Table Entity

**File**: `internal/domain/model/table.go`

**Key Features:**
- Table number (unique)
- Location (INDOOR, OUTDOOR, VIP, PATIO)
- Capacity (1-50 people)
- Status (AVAILABLE, OCCUPIED, RESERVED, MAINTENANCE)
- QR code storage (base64 encoded image)
- Status transitions with validation

**Example Usage:**
```go
table, _ := model.NewTable(5, model.LocationIndoor, 4, "Near window")
table.MarkOccupied()
table.GenerateQRCode("data:image/png;base64,...")
```

### GuestOrder Entity

**File**: `internal/domain/model/guestorder.go`

**Key Features:**
- Order number (auto-generated: ORD-YYYYMMDD-0001)
- Table reference
- Customer info (name, phone optional)
- Items with notes
- Financial calculations (subtotal, tax, discount, total)
- Payment tracking (method, status, amount, change)
- Status workflow (PENDING → CONFIRMED → PREPARING → READY → SERVED)
- Session tracking

**Example Usage:**
```go
order, _ := model.NewGuestOrder(tableID, 5, "John", "08123456789", "session123")
order.AddItem("prod1", "Nasi Goreng", 2, 25000, "Tidak pedas")
order.ApplyTax(11)
order.ProcessPayment(model.PaymentCash, 100000)
```

---

## 🔧 Services

### QRCodeService

**File**: `internal/domain/service/qrcode_service.go`

**Features:**
- Custom QR code generation with logo overlay
- Configurable colors, size, error correction
- Output as image, bytes, base64 string, or file
- Bulk generation for multiple tables

**Configuration:**
```go
config := QRCodeConfig{
    BaseURL:         "https://pos.restaurant.com",
    MerchantName:    "My Restaurant",
    MerchantLogo:    "/path/to/logo.png",
    Size:            256,
    ErrorCorrection: qrcode.Medium,
    ForegroundColor: color.Black,
    BackgroundColor: color.White,
}
```

**Usage:**
```go
qrService := NewQRCodeService(config)
qrBase64, _ := qrService.GenerateQRString(tableNumber, tableID)
// Returns: "data:image/png;base64,iVBOR..."
```

### TableService

**File**: `internal/domain/service/table_service.go`

**Methods:**
- `Create()` - Create new table with validation
- `GetByID()` - Get table by ID
- `Update()` - Update table details
- `Delete()` - Delete table (if not occupied)
- `List()` - List all tables with filtering
- `UpdateStatus()` - Change table status
- `GenerateQR()` - Generate QR code for table
- `GetAvailableTables()` - Get available tables by location
- `Count()` - Count total tables

### GuestOrderService

**File**: `internal/domain/service/guestorder_service.go`

**Methods:**
- `CreateOrder()` - Create new guest order
- `AddItem()` - Add item to order (with stock validation)
- `RemoveItem()` - Remove item from order
- `UpdateItemQuantity()` - Update item quantity
- `ProcessCheckout()` - Process payment & complete order
- `CancelOrder()` - Cancel order & restore inventory
- `UpdateOrderStatus()` - Update order status (staff action)
- `GetOrderByID()` - Get order details
- `GetPendingOrders()` - Get all pending orders
- `GetActiveOrders()` - Get all active orders
- `GetOrdersByTable()` - Get orders for specific table
- `GetTodaySales()` - Get today's sales summary

---

## 📡 API Endpoints (Planned)

### Public Endpoints (No Authentication)

```
GET  /api/public/menu              - Browse menu (reuse existing inventory)
POST /api/guest/orders             - Create new order
GET  /api/guest/orders/{id}        - Get order details
POST /api/guest/orders/{id}/items  - Add item to order
PUT  /api/guest/orders/{id}/items  - Update item quantity
DELETE /api/guest/orders/{id}/items/{productId} - Remove item
POST /api/guest/orders/{id}/checkout - Process checkout
```

### Protected Endpoints (Staff/Admin)

```
# Table Management
GET    /api/tables                  - List all tables
POST   /api/tables                  - Create new table
GET    /api/tables/{id}             - Get table details
PUT    /api/tables/{id}             - Update table
DELETE /api/tables/{id}             - Delete table
POST   /api/tables/{id}/qr          - Generate QR code
POST   /api/tables/{id}/status      - Update table status

# Order Management
GET    /api/orders                  - List all orders (with filters)
GET    /api/orders/pending          - Get pending orders
GET    /api/orders/active           - Get active orders
GET    /api/orders/{id}             - Get order details
POST   /api/orders/{id}/status      - Update order status
POST   /api/orders/{id}/cancel      - Cancel order
GET    /api/orders/table/{tableId}  - Get orders by table

# Reports
GET    /api/reports/sales/today     - Today's sales summary
```

---

## 🗄️ Database Schema

### Tables Table

```sql
CREATE TABLE tables (
    id VARCHAR(100) PRIMARY KEY,
    number INTEGER NOT NULL UNIQUE,
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

### Guest Orders Table

```sql
CREATE TABLE guest_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) NOT NULL UNIQUE,
    table_id VARCHAR(100) NOT NULL REFERENCES tables(id),
    table_number INTEGER NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(50),
    items JSONB NOT NULL,
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

---

## 🚀 Next Steps

### Remaining Components to Create:

1. **Usecases** (`internal/application/usecase/`)
   - `table_usecase.go` - Table usecase interface & implementation
   - `guestorder_usecase.go` - Guest order usecase interface & implementation

2. **HTTP Handlers** (`internal/handler/`)
   - `table_handler.go` - Table CRUD handlers
   - `guestorder_handler.go` - Guest order handlers

3. **Configuration** (`internal/infrastructure/config/`)
   - Add QR code configuration fields

4. **Server Wiring** (`internal/infrastructure/http/server.go`)
   - Update `buildApp()` to create new services
   - Update `setupRoutes()` to add new routes

5. **Dependencies** (`go.mod`)
   - Add `github.com/skip2/go-qrcode`

6. **Testing & Verification**
   - Build verification
   - Test all endpoints
   - QR code generation test

---

## 📝 Implementation Notes

### QR Code Generation Flow

```
1. Admin creates table via API
2. System generates unique QR code with embedded URL
3. QR code stored in database as base64 string
4. Admin downloads QR code (PNG/PDF)
5. Print and place QR code on physical table
6. Customer scans → opens web app with table context
```

### Order Flow

```
1. Customer scans QR code
   → URL: https://pos.app/order?table=5&id=uuid
2. Web app extracts table ID from URL
3. Customer browses menu (public API)
4. Customer creates order (POST /api/guest/orders)
   → Order created with status: PENDING
   → Table marked as OCCUPIED
5. Customer adds items (POST /api/guest/orders/{id}/items)
   → Stock validation
6. Customer checks out (POST /api/guest/orders/{id}/checkout)
   → Inventory deducted
   → Table remains OCCUPIED
   → Order status: CONFIRMED
7. Kitchen receives order notification
8. Staff updates order status:
   CONFIRMED → PREPARING → READY → SERVED
9. When SERVED, table marked as AVAILABLE
```

### Security Considerations

- **Public endpoints** have no authentication (designed for guest access)
- **Rate limiting** should be added to prevent abuse
- **Session tracking** via `session_id` for analytics
- **Phone number** optional for privacy
- **Order validation** prevents unauthorized modifications after checkout

### Performance Optimizations

- **Batch queries** for order listing with pagination
- **JSONB storage** for order items (flexible schema)
- **Indexes** on frequently queried fields
- **Connection pooling** via PostgreSQL settings

---

## 🎯 QR Code Design Specifications

### QR Code Content

```
https://yourdomain.com/order?table={TABLE_NUMBER}&id={TABLE_UUID}
```

### QR Code Features

- **Size**: 256x256 pixels (scalable to 512x512 for print)
- **Error Correction**: Medium (15% data recovery)
- **Logo Overlay**: Center, 20% of QR size
- **Format**: PNG with transparent background
- **Colors**: Black on white (customizable)

### Print Specifications

- **Recommended Size**: 10cm x 10cm minimum
- **Material**: Waterproof sticker or acrylic
- **Placement**: Center of table or table stand
- **Design**: Include table number and restaurant logo

---

## 📚 Related Documentation

- [README.md](../README.md) - Main project documentation
- [OPTIMIZATION_SUMMARY.md](../OPTIMIZATION_SUMMARY.md) - Performance optimizations
- [CHANGELOG.md](../CHANGELOG.md) - Version history

---

**Version**: 1.0.0  
**Date**: April 4, 2026  
**Status**: Backend Implementation In Progress
