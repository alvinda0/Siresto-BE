# Add Item to Order API

Endpoint untuk menambahkan item ke pesanan yang sudah ada.

## Endpoint

```
POST /api/v1/external/orders/quick/:id
```

## Authentication

Memerlukan Bearer Token (login sebagai OWNER, ADMIN, atau CASHIER)

## Path Parameter

| Parameter | Type | Deskripsi |
|-----------|------|-----------|
| `id` | uuid | ID order yang mau ditambah item |

## Request Body

```json
{
  "product_id": "uuid-product",
  "quantity": 2,
  "note": "Extra pedas"
}
```

### Field Wajib

| Field | Type | Deskripsi |
|-------|------|-----------|
| `product_id` | uuid | ID produk yang mau ditambah |
| `quantity` | integer | Jumlah (minimal 1) |

### Field Opsional

| Field | Type | Deskripsi |
|-------|------|-----------|
| `note` | string | Catatan khusus untuk item |

## Response Success (200)

```json
{
  "status": "success",
  "message": "Item added to order successfully",
  "data": {
    "id": "uuid-order",
    "company_id": "uuid-company",
    "branch_id": "uuid-branch",
    "customer_name": "",
    "customer_phone": "",
    "table_number": "A5",
    "notes": "",
    "order_method": "DINE_IN",
    "status": "PENDING",
    "subtotal_amount": 75000,
    "tax_amount": 8250,
    "total_amount": 83250,
    "tax_details": [
      {
        "tax_id": "uuid-tax",
        "tax_name": "PB1",
        "percentage": 10,
        "priority": 1,
        "base_amount": 75000,
        "tax_amount": 7500
      }
    ],
    "order_items": [
      {
        "id": "uuid-item-1",
        "product_id": "uuid-product-1",
        "product_name": "Nasi Goreng",
        "quantity": 2,
        "price": 25000,
        "subtotal": 50000,
        "note": ""
      },
      {
        "id": "uuid-item-2",
        "product_id": "uuid-product-2",
        "product_name": "Es Teh",
        "quantity": 2,
        "price": 5000,
        "subtotal": 10000,
        "note": "Extra pedas"
      }
    ],
    "created_at": "2024-01-15 10:30:00",
    "updated_at": "2024-01-15 10:35:00"
  }
}
```

## Response Error

### 400 Bad Request
```json
{
  "status": "error",
  "message": "Invalid request body",
  "error": "Key: 'AddOrderItemRequest.ProductID' Error:Field validation for 'ProductID' failed on the 'required' tag"
}
```

### 404 Not Found
```json
{
  "status": "error",
  "message": "Failed to add item to order",
  "error": "order not found"
}
```

### 401 Unauthorized
```json
{
  "status": "error",
  "message": "Failed to add item to order",
  "error": "unauthorized to modify this order"
}
```

## Fitur

✅ Tambah item ke order existing
✅ Otomatis merge jika product_id sama (quantity dijumlahkan)
✅ Otomatis recalculate subtotal
✅ Otomatis recalculate diskon (jika ada promo)
✅ Otomatis recalculate pajak bertingkat
✅ Otomatis recalculate total
✅ Validasi produk tersedia
✅ Multi-tenant (company & branch)
✅ WebSocket broadcast real-time

## Contoh Penggunaan

### 1. Tambah 1 Item Sederhana
```bash
curl -X POST http://localhost:8080/api/v1/external/orders/quick/ORDER_UUID \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "uuid-es-teh",
    "quantity": 2
  }'
```

**Hasil:** Jika Es Teh belum ada, buat item baru. Jika sudah ada, quantity +2.

### 2. Tambah Item yang Sudah Ada
```bash
# Order sudah punya: Nasi Goreng x3
# Add lagi: Nasi Goreng x2
curl -X POST http://localhost:8080/api/v1/external/orders/quick/ORDER_UUID \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "uuid-nasi-goreng",
    "quantity": 2
  }'
```

**Hasil:** Nasi Goreng quantity jadi 5 (3+2), bukan bikin item baru.

### 3. Tambah Item Multiple Quantity
```bash
curl -X POST http://localhost:8080/api/v1/external/orders/quick/ORDER_UUID \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "uuid-kopi",
    "quantity": 5,
    "note": "Less sugar"
  }'
```

## Flow Proses

1. Validasi order ID exists
2. Validasi user punya akses ke order (company & branch)
3. Validasi product exists dan available
4. **Cek apakah product_id sudah ada di order items**
   - Jika sudah ada: Update quantity (quantity lama + quantity baru)
   - Jika belum ada: Buat item baru
5. Recalculate subtotal dari semua items
6. Recalculate diskon jika ada promo
7. Recalculate pajak bertingkat
8. Recalculate total
9. Update order
10. Broadcast via WebSocket
11. Return order lengkap dengan semua items

## Testing

Gunakan script PowerShell:
```powershell
.\test_quick_order.ps1
```

## Notes

- Endpoint ini untuk menambah item ke order yang sudah ada
- **Jika product_id sama, quantity akan dijumlahkan (tidak bikin item baru)**
- Jika order punya promo, diskon akan dihitung ulang otomatis
- Pajak dihitung ulang berdasarkan subtotal baru
- WebSocket akan broadcast update ke semua client yang subscribe
- Order status tidak berubah (tetap sesuai status sebelumnya)
- Note akan di-update jika item sudah ada dan note baru tidak kosong
