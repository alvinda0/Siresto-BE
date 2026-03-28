# Order API Documentation

## Overview
API untuk mengelola order (pesanan) di sistem SIRESTO. Mendukung berbagai metode order seperti Dine In, Take Away, dan Delivery.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Semua endpoint (kecuali public order) memerlukan JWT token di header:
```
Authorization: Bearer <token>
```

---

## Endpoints

### 1. Create Order (Authenticated)
Membuat order baru dengan autentikasi (company_id dan branch_id dari JWT token).

**Endpoint:** `POST /api/v1/orders`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "customer_name": "Sasa",
  "customer_phone": "08123123123",
  "table_number": "A1",
  "notes": "Test Order",
  "referral_code": "",
  "order_method": "DINE_IN",
  "promo_code": "",
  "order_items": [
    {
      "product_id": "c61965b6-270c-466f-b5eb-3ff722dfde48",
      "quantity": 3,
      "note": "-"
    }
  ]
}
```

**Field Descriptions:**
- `customer_name` (optional): Nama pelanggan
- `customer_phone` (optional): Nomor telepon pelanggan
- `table_number` (required): Nomor meja
- `notes` (optional): Catatan order
- `referral_code` (optional): Kode referral
- `order_method` (required): Metode order - `DINE_IN`, `TAKE_AWAY`, atau `DELIVERY`
- `promo_code` (optional): Kode promo
- `order_items` (required): Array item yang dipesan
  - `product_id` (required): UUID produk
  - `quantity` (required): Jumlah item (min: 1)
  - `note` (optional): Catatan untuk item

**Response Success (201):**
```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "company_id": "123e4567-e89b-12d3-a456-426614174000",
    "branch_id": "789e0123-e89b-12d3-a456-426614174000",
    "customer_name": "Sasa",
    "customer_phone": "08123123123",
    "table_number": "A1",
    "notes": "Test Order",
    "referral_code": "",
    "order_method": "DINE_IN",
    "promo_code": "",
    "status": "PENDING",
    "total_amount": 150000,
    "order_items": [
      {
        "id": "abc12345-e29b-41d4-a716-446655440000",
        "product_id": "c61965b6-270c-466f-b5eb-3ff722dfde48",
        "product_name": "Nasi Goreng",
        "quantity": 3,
        "price": 50000,
        "subtotal": 150000,
        "note": "-"
      }
    ],
    "created_at": "2024-01-15 10:30:00",
    "updated_at": "2024-01-15 10:30:00"
  }
}
```

---

### 2. Create Public Order (No Authentication)
Membuat order tanpa autentikasi. Memerlukan company_id dan branch_id di request body.

**Endpoint:** `POST /api/v1/public/orders`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "company_id": "123e4567-e89b-12d3-a456-426614174000",
  "branch_id": "789e0123-e89b-12d3-a456-426614174000",
  "customer_name": "Sasa",
  "customer_phone": "08123123123",
  "table_number": "A1",
  "notes": "Test Order",
  "referral_code": "",
  "order_method": "DINE_IN",
  "promo_code": "",
  "order_items": [
    {
      "product_id": "c61965b6-270c-466f-b5eb-3ff722dfde48",
      "quantity": 3,
      "note": "-"
    }
  ]
}
```

**Response:** Same as Create Order (Authenticated)

---

### 3. Update Order
Update order yang sudah ada.

**Endpoint:** `PUT /api/v1/orders/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "customer_name": "Sasa Updated",
  "customer_phone": "08123123124",
  "table_number": "A2",
  "notes": "Updated notes",
  "order_method": "TAKE_AWAY",
  "status": "CONFIRMED",
  "order_items": [
    {
      "product_id": "c61965b6-270c-466f-b5eb-3ff722dfde48",
      "quantity": 2,
      "note": "Less spicy"
    }
  ]
}
```

**Field Descriptions:**
- Semua field optional
- `status`: `PENDING`, `CONFIRMED`, `PREPARING`, `READY`, `COMPLETED`, `CANCELLED`
- Jika `order_items` dikirim, akan mengganti semua item yang ada

**Response Success (200):**
```json
{
  "success": true,
  "message": "Order updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "CONFIRMED",
    ...
  }
}
```

---

### 4. Get All Orders
Mendapatkan daftar semua order dengan pagination.

**Endpoint:** `GET /api/v1/orders`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (optional): Nomor halaman (default: 1)
- `limit` (optional): Jumlah item per halaman (default: 10)

**Example:**
```
GET /api/v1/orders?page=1&limit=10
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Orders retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "company_id": "123e4567-e89b-12d3-a456-426614174000",
      "branch_id": "789e0123-e89b-12d3-a456-426614174000",
      "customer_name": "Sasa",
      "table_number": "A1",
      "order_method": "DINE_IN",
      "status": "PENDING",
      "total_amount": 150000,
      "order_items": [...],
      "created_at": "2024-01-15 10:30:00",
      "updated_at": "2024-01-15 10:30:00"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

---

### 5. Get Order by ID
Mendapatkan detail order berdasarkan ID.

**Endpoint:** `GET /api/v1/orders/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Example:**
```
GET /api/v1/orders/550e8400-e29b-41d4-a716-446655440000
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Order retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "company_id": "123e4567-e89b-12d3-a456-426614174000",
    "branch_id": "789e0123-e89b-12d3-a456-426614174000",
    "customer_name": "Sasa",
    "customer_phone": "08123123123",
    "table_number": "A1",
    "notes": "Test Order",
    "order_method": "DINE_IN",
    "status": "PENDING",
    "total_amount": 150000,
    "order_items": [
      {
        "id": "abc12345-e29b-41d4-a716-446655440000",
        "product_id": "c61965b6-270c-466f-b5eb-3ff722dfde48",
        "product_name": "Nasi Goreng",
        "quantity": 3,
        "price": 50000,
        "subtotal": 150000,
        "note": "-"
      }
    ],
    "created_at": "2024-01-15 10:30:00",
    "updated_at": "2024-01-15 10:30:00"
  }
}
```

---

## Order Status Flow

```
PENDING → CONFIRMED → PREPARING → READY → COMPLETED
                                        ↓
                                   CANCELLED
```

**Status Descriptions:**
- `PENDING`: Order baru dibuat, menunggu konfirmasi
- `CONFIRMED`: Order dikonfirmasi, siap diproses
- `PREPARING`: Order sedang diproses/dimasak
- `READY`: Order siap diambil/diantar
- `COMPLETED`: Order selesai
- `CANCELLED`: Order dibatalkan

---

## Order Methods

- `DINE_IN`: Makan di tempat
- `TAKE_AWAY`: Bungkus/dibawa pulang
- `DELIVERY`: Delivery/antar

---

## Error Responses

**400 Bad Request:**
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "validation error details"
}
```

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "invalid or missing token"
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "Order not found",
  "error": "order with specified ID does not exist"
}
```

---

## Business Rules

1. Order hanya bisa dibuat jika produk tersedia (`is_available = true`)
2. Produk harus berada di branch yang sama dengan order
3. Total amount dihitung otomatis dari (price × quantity) semua items
4. User hanya bisa melihat/update order dari company dan branch mereka
5. Saat update order items, semua item lama akan diganti dengan yang baru

---

## Testing Examples

### cURL Examples

**Create Order:**
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "table_number": "A1",
    "order_method": "DINE_IN",
    "order_items": [
      {
        "product_id": "c61965b6-270c-466f-b5eb-3ff722dfde48",
        "quantity": 2
      }
    ]
  }'
```

**Create Public Order:**
```bash
curl -X POST http://localhost:8080/api/v1/public/orders \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "123e4567-e89b-12d3-a456-426614174000",
    "branch_id": "789e0123-e89b-12d3-a456-426614174000",
    "table_number": "A1",
    "order_method": "DINE_IN",
    "order_items": [
      {
        "product_id": "c61965b6-270c-466f-b5eb-3ff722dfde48",
        "quantity": 2
      }
    ]
  }'
```

**Get All Orders:**
```bash
curl -X GET "http://localhost:8080/api/v1/orders?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Get Order by ID:**
```bash
curl -X GET http://localhost:8080/api/v1/orders/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Update Order:**
```bash
curl -X PUT http://localhost:8080/api/v1/orders/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "CONFIRMED"
  }'
```
