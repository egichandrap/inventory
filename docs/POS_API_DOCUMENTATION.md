# POS (Point of Sale) API Documentation

## Overview
Ultimate POS System with user authentication, role-based access control, and complete POS functionality.

## Base URL
```
http://localhost:8080/api
```

## Authentication
All protected endpoints require JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Roles
- **SUPER_ADMIN**: Full access to everything, including user management
- **ADMIN**: Can manage inventory and access POS features
- **CASHIER**: Can only access POS features (cart, checkout, transactions)
- **VIEWER**: Read-only access

---

## 1. Authentication Endpoints

### 1.1 Login
**POST** `/auth/login`

**Request Body:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "id": "uuid",
      "username": "admin",
      "email": "admin@pos.local",
      "full_name": "Administrator",
      "role": "ADMIN",
      "status": "ACTIVE",
      "created_at": "2026-04-03T10:00:00Z",
      "updated_at": "2026-04-03T10:00:00Z",
      "last_login_at": "2026-04-03T10:00:00Z"
    }
  }
}
```

### 1.2 Register
**POST** `/auth/register`

**Request Body:**
```json
{
  "username": "cashier1",
  "email": "cashier1@pos.local",
  "password": "cashier123",
  "full_name": "Cashier One",
  "role": "CASHIER"
}
```

### 1.3 Logout
**POST** `/auth/logout`

**Headers:** Authorization required

### 1.4 Get Current User
**GET** `/auth/me`

**Headers:** Authorization required

### 1.5 Change Password
**POST** `/auth/change-password`

**Headers:** Authorization required

**Request Body:**
```json
{
  "old_password": "oldpass",
  "new_password": "newpass123"
}
```

### 1.6 Refresh Token
**POST** `/auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGci..."
}
```

---

## 2. User Management (Admin Only)

### 2.1 List Users
**GET** `/admin/users`

**Headers:** Authorization required (SUPER_ADMIN or ADMIN)

**Query Parameters:**
- `limit` (default: 20)
- `offset` (default: 0)
- `role` (SUPER_ADMIN, ADMIN, CASHIER, VIEWER)
- `status` (ACTIVE, INACTIVE, SUSPENDED)
- `search` (search by username, email, or full_name)

### 2.2 Get User by ID
**GET** `/admin/users/{id}`

**Headers:** Authorization required (SUPER_ADMIN or ADMIN)

### 2.3 Update User
**PUT** `/admin/users/{id}`

**Headers:** Authorization required (SUPER_ADMIN or ADMIN)

**Request Body:**
```json
{
  "email": "newemail@pos.local",
  "full_name": "New Name",
  "role": "ADMIN",
  "status": "ACTIVE"
}
```

### 2.4 Delete User
**DELETE** `/admin/users/{id}`

**Headers:** Authorization required (SUPER_ADMIN or ADMIN)

---

## 3. Inventory Endpoints

### 3.1 List Inventory
**GET** `/inventory`

**Headers:** Authorization required

**Query Parameters:**
- `sku`, `name`, `location`
- `min_qty`, `max_qty`
- `limit`, `offset`

### 3.2 Create Inventory (Admin Only)
**POST** `/inventory`

**Headers:** Authorization required (SUPER_ADMIN or ADMIN)

**Request Body:**
```json
{
  "sku": "PROD-001",
  "name": "Product Name",
  "description": "Product description",
  "quantity": 100,
  "unit": "pcs",
  "location": "Warehouse A",
  "min_stock": 10,
  "max_stock": 1000,
  "price": 50000
}
```

### 3.3 Get Inventory by ID
**GET** `/inventory/{id}`

**Headers:** Authorization required

### 3.4 Update Inventory (Admin Only)
**PUT** `/inventory/{id}`

**Headers:** Authorization required (SUPER_ADMIN or ADMIN)

### 3.5 Delete Inventory (Admin Only)
**DELETE** `/inventory/{id}`

**Headers:** Authorization required (SUPER_ADMIN or ADMIN)

### 3.6 Update Stock
**PUT** `/inventory/{id}/stock`

**Headers:** Authorization required

**Request Body:**
```json
{
  "quantity": 150
}
```

### 3.7 Adjust Stock
**POST** `/inventory/{id}/stock/adjust`

**Headers:** Authorization required

**Request Body:**
```json
{
  "amount": 10
}
```

---

## 4. POS Endpoints

### 4.1 Create Cart
**POST** `/pos/cart`

**Headers:** Authorization required

**Request Body:**
```json
{
  "customer_name": "John Doe"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Cart berhasil dibuat",
  "data": {
    "id": "cart-1234567890",
    "user_id": "user-uuid",
    "customer_name": "John Doe",
    "items": [],
    "total_amount": 0,
    "created_at": "2026-04-03T10:00:00Z",
    "updated_at": "2026-04-03T10:00:00Z"
  }
}
```

### 4.2 Get My Cart
**GET** `/pos/cart/my`

**Headers:** Authorization required

### 4.3 Get Cart by ID
**GET** `/pos/cart/{id}`

**Headers:** Authorization required

### 4.4 Add Item to Cart
**POST** `/pos/cart/{id}/items`

**Headers:** Authorization required

**Request Body:**
```json
{
  "product_id": "product-uuid",
  "quantity": 2
}
```

### 4.5 Update Cart Item Quantity
**PUT** `/pos/cart/{id}/items`

**Headers:** Authorization required

**Request Body:**
```json
{
  "product_id": "product-uuid",
  "quantity": 5
}
```

### 4.6 Remove Item from Cart
**DELETE** `/pos/cart/{id}/items`

**Headers:** Authorization required

**Request Body:**
```json
{
  "product_id": "product-uuid"
}
```

### 4.7 Clear Cart
**POST** `/pos/cart/{id}/clear`

**Headers:** Authorization required

### 4.8 Delete Cart
**DELETE** `/pos/cart/{id}`

**Headers:** Authorization required

### 4.9 Checkout
**POST** `/pos/checkout/{cart_id}`

**Headers:** Authorization required

**Request Body:**
```json
{
  "payment_method": "CASH",
  "payment_amount": 150000,
  "customer_name": "John Doe",
  "notes": "Please wrap the gift"
}
```

**Payment Methods:**
- `CASH` - Cash payment
- `CARD` - Credit/Debit card (TODO)
- `QRIS` - QRIS payment (TODO)
- `E_WALLET` - E-wallet (GoPay, OVO, etc.) (TODO)
- `TRANSFER` - Bank transfer (TODO)

**Response:**
```json
{
  "success": true,
  "message": "Checkout berhasil",
  "data": {
    "id": "transaction-uuid",
    "transaction_no": "TRX-20260403-0001",
    "cashier_id": "cashier-uuid",
    "cashier_name": "Cashier User",
    "customer_name": "John Doe",
    "items": [
      {
        "product_id": "product-uuid",
        "product_name": "Product Name",
        "sku": "PROD-001",
        "quantity": 2,
        "unit_price": 50000,
        "subtotal": 100000
      }
    ],
    "subtotal": 100000,
    "discount_amount": 0,
    "discount_percent": 0,
    "tax_amount": 11000,
    "tax_percent": 11,
    "total_amount": 111000,
    "payment_method": "CASH",
    "payment_amount": 150000,
    "change_amount": 39000,
    "status": "COMPLETED",
    "notes": "Please wrap the gift",
    "created_at": "2026-04-03T10:00:00Z"
  }
}
```

### 4.10 List Transactions
**GET** `/pos/transactions`

**Headers:** Authorization required

**Query Parameters:**
- `limit`, `offset`
- `status` (PENDING, COMPLETED, CANCELLED, REFUNDED)
- `payment_method` (CASH, CARD, QRIS, E_WALLET, TRANSFER)
- `search` (search by transaction_no or customer_name)

### 4.11 Get Transaction by ID
**GET** `/pos/transactions/{id}`

**Headers:** Authorization required

### 4.12 Cancel Transaction
**POST** `/pos/transactions/{id}/cancel`

**Headers:** Authorization required

**Note:** This will restore the inventory

### 4.13 Get Today's Sales Summary
**GET** `/pos/sales/today`

**Headers:** Authorization required

**Response:**
```json
{
  "success": true,
  "message": "Berhasil mengambil sales summary",
  "data": {
    "total_sales": 1500000,
    "total_transactions": 15,
    "total_items": 45,
    "date": "2026-04-03"
  }
}
```

---

## 5. Default Users

| Username | Password | Role | Email |
|----------|----------|------|-------|
| superadmin | admin123 | SUPER_ADMIN | superadmin@pos.local |
| admin | admin123 | ADMIN | admin@pos.local |
| cashier | cashier123 | CASHIER | cashier@pos.local |

---

## 6. Error Responses

All errors follow this format:
```json
{
  "success": false,
  "error": {
    "code": "ERR_VALIDATION",
    "message": "Error message",
    "details": "Optional details"
  }
}
```

**Common Error Codes:**
- `ERR_VALIDATION` - Validation error (400)
- `ERR_UNAUTHENTICATED` - Not authenticated (401)
- `ERR_FORBIDDEN` - Not authorized (403)
- `ERR_NOT_FOUND` - Resource not found (404)
- `ERR_CONFLICT` - Resource conflict (409)
- `ERR_INTERNAL` - Internal server error (500)

---

## 7. Quick Start

### 7.1 Start the Server
```bash
go run cmd/main.go -server
```

### 7.2 Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 7.3 Create Cart
```bash
curl -X POST http://localhost:8080/api/pos/cart \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"customer_name":"John Doe"}'
```

### 7.4 Add Item to Cart
```bash
curl -X POST http://localhost:8080/api/pos/cart/{cart_id}/items \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"product_id":"product-uuid","quantity":2}'
```

### 7.5 Checkout
```bash
curl -X POST http://localhost:8080/api/pos/checkout/{cart_id} \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "CASH",
    "payment_amount": 150000,
    "customer_name": "John Doe"
  }'
```

---

## 8. Architecture

This project follows **Clean Architecture** with **Domain-Driven Design (DDD)** principles:

```
Domain Layer (Entities, Repositories Interfaces, Services)
    ↓
Handler Layer (Use Cases)
    ↓
Infrastructure Layer (HTTP, Database, JWT, Repository Implementations)
```

### Key Components:
- **Domain Models**: User, Cart, Transaction, Inventory
- **Repository Interfaces**: Define data operations
- **Services**: Business logic
- **Handlers**: HTTP request handlers
- **Middleware**: Authentication and authorization

---

## 9. TODO / Future Enhancements

- [ ] Implement PostgreSQL repositories for Cart and Transaction
- [ ] Payment gateway integration (Midtrans, Xendit, Stripe)
- [ ] QRIS payment support
- [ ] E-wallet payment support
- [ ] Card payment support
- [ ] Refund functionality
- [ ] Payment reconciliation
- [ ] Advanced reporting and analytics
- [ ] Barcode/QR code scanning
- [ ] Multi-store support
- [ ] Customer loyalty program
- [ ] Inventory alerts (low stock, out of stock)
- [ ] Batch operations for inventory
- [ ] Export transactions to CSV/Excel
- [ ] Receipt generation and printing

---

## 10. Database Migrations

Migrations are located in `migrations/` directory:
- `001_create_inventories_table.up.sql`
- `002_create_tokens_table.up.sql`
- `003_seed_inventory_data.up.sql`
- `004_create_users_table.up.sql`
- `005_create_pos_tables.up.sql`

Run migrations automatically on server startup.

---

## License
MIT License
