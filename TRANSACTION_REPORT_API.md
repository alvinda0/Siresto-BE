# Transaction Report API

API untuk mendapatkan laporan transaksi dengan filter tanggal, jam, search, dan pagination. Report otomatis terfilter berdasarkan company dan branch user yang login.

## Endpoint

```
GET /api/v1/external/reports/transactions
```

## Authentication

Memerlukan Bearer Token di header:
```
Authorization: Bearer <your_token>
```

## Query Parameters

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `start_date` | string | No | Tanggal mulai (Format: YYYY-MM-DD) | `2024-01-01` |
| `end_date` | string | No | Tanggal akhir (Format: YYYY-MM-DD) | `2024-01-31` |
| `start_time` | string | No | Jam mulai (Format: HH:MM) | `08:00` |
| `end_time` | string | No | Jam akhir (Format: HH:MM) | `22:00` |
| `search` | string | No | Search by customer name, phone, atau order ID | `John` |
| `status` | string | No | Filter by order status | `completed` |
| `payment_status` | string | No | Filter by payment status | `paid` |
| `payment_method` | string | No | Filter by payment method | `cash` |
| `order_method` | string | No | Filter by order method | `dine_in` |
| `page` | int | No | Halaman (default: 1) | `1` |
| `limit` | int | No | Jumlah data per halaman (default: 10) | `20` |

### Filter Values

**Order Status:**
- `pending` - Order baru dibuat
- `preparing` - Sedang diproses
- `ready` - Siap disajikan
- `completed` - Selesai
- `cancelled` - Dibatalkan

**Payment Status:**
- `unpaid` - Belum dibayar
- `paid` - Sudah dibayar
- `refunded` - Dikembalikan

**Order Method:**
- `dine_in` - Makan di tempat
- `takeaway` - Bawa pulang
- `delivery` - Delivery

## Response

### Success Response (200 OK)

```json
{
  "success": true,
  "message": "Transaction report retrieved successfully",
  "status": 200,
  "timestamp": "2024-01-15T10:30:00Z",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "order_number": "550e8400",
      "customer_name": "John Doe",
      "customer_phone": "081234567890",
      "table_number": "A1",
      "order_method": "dine_in",
      "status": "completed",
      "payment_status": "paid",
      "payment_method": "cash",
      "subtotal_amount": 100000,
      "tax_amount": 10000,
      "discount_amount": 5000,
      "total_amount": 105000,
      "paid_amount": 110000,
      "change_amount": 5000,
      "promo_code": "PROMO10",
      "company_name": "Restaurant ABC",
      "branch_name": "Cabang Jakarta",
      "created_at": "2024-01-15T10:30:00Z",
      "paid_at": "2024-01-15T11:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 50,
    "total_pages": 5
  }
}
```

### Error Response (400 Bad Request)

```json
{
  "success": false,
  "message": "Bad Request",
  "status": 400,
  "timestamp": "2024-01-15T10:30:00Z",
  "error": "Invalid query parameters"
}
```

### Error Response (401 Unauthorized)

```json
{
  "success": false,
  "message": "Unauthorized",
  "status": 401,
  "timestamp": "2024-01-15T10:30:00Z",
  "error": "Invalid or expired token"
}
```

### Error Response (500 Internal Server Error)

```json
{
  "success": false,
  "message": "Failed to get transaction report",
  "status": 500,
  "timestamp": "2024-01-15T10:30:00Z",
  "error": "database connection error"
}
```

## Examples

### 1. Get All Transactions (Default)

```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 2. Filter by Date Range

```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. Filter by Date and Time Range

```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?start_date=2024-01-15&end_date=2024-01-15&start_time=08:00&end_time=22:00" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. Search by Customer Name

```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?search=John" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 5. Filter by Status and Payment Method

```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?status=completed&payment_method=cash" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 6. With Pagination

```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?page=2&limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 7. Combined Filters

```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?start_date=2024-01-01&end_date=2024-01-31&status=completed&payment_status=paid&search=John&page=1&limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Notes

1. Report otomatis terfilter berdasarkan `company_id` dan `branch_id` dari user yang login
2. Jika `start_date` dan `end_date` tidak diisi, akan menampilkan semua transaksi
3. Filter `start_time` dan `end_time` hanya berlaku jika `start_date` atau `end_date` diisi
4. Search akan mencari di customer name, phone, dan order ID
5. Semua filter bersifat opsional dan bisa dikombinasikan
6. Default pagination: page=1, limit=10
7. Data diurutkan berdasarkan `created_at` DESC (terbaru dulu)

## Access Control

- Endpoint ini hanya bisa diakses oleh user dengan role `external` (OWNER, ADMIN, CASHIER)
- User hanya bisa melihat transaksi dari company dan branch mereka sendiri
- Filtering berdasarkan company dan branch dilakukan otomatis di backend
