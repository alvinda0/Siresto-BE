# Quick Start: Update Order Status

Panduan cepat untuk menggunakan endpoint update status order.

## Endpoint

```
PATCH /api/v1/external/orders/:id/status
```

## Status Order yang Tersedia

1. `PENDING` - Order baru dibuat (default)
2. `CONFIRMED` - Order dikonfirmasi
3. `PROCESSING` - Order sedang diproses/disiapkan
4. `READY` - Order siap diambil/diantar
5. `COMPLETED` - Order selesai
6. `CANCELLED` - Order dibatalkan

## Flow Status Order yang Umum

```
PENDING → PROCESSING → READY → COMPLETED
         ↓
      CANCELLED (bisa dari status manapun)
```

## Contoh Request

### 1. Update dari PENDING ke PROCESSING

```bash
curl -X PATCH http://localhost:8080/api/v1/external/orders/{order_id}/status \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "PROCESSING"}'
```

### 2. Update ke READY

```bash
curl -X PATCH http://localhost:8080/api/v1/external/orders/{order_id}/status \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "READY"}'
```

### 3. Update ke COMPLETED

```bash
curl -X PATCH http://localhost:8080/api/v1/external/orders/{order_id}/status \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "COMPLETED"}'
```

### 4. Cancel Order

```bash
curl -X PATCH http://localhost:8080/api/v1/external/orders/{order_id}/status \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "CANCELLED"}'
```

## Testing

### PowerShell (Windows)

1. Set product ID:
```powershell
$env:PRODUCT_ID = "your-product-uuid"
```

2. Jalankan test script:
```powershell
.\test_update_order_status.ps1
```

### Bash (Linux/Mac)

1. Set product ID:
```bash
export PRODUCT_ID="your-product-uuid"
```

2. Jalankan test script:
```bash
chmod +x test_update_order_status.sh
./test_update_order_status.sh
```

## Response Success

```json
{
  "status": "success",
  "message": "Order status updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "PROCESSING",
    "table_number": "A1",
    "order_method": "DINE_IN",
    "subtotal_amount": 50000,
    "discount_amount": 0,
    "tax_amount": 5000,
    "total_amount": 55000,
    "order_items": [...],
    "created_at": "2024-01-15 10:30:00",
    "updated_at": "2024-01-15 10:35:00"
  }
}
```

## Fitur

- ✅ Update status order dengan satu endpoint
- ✅ Validasi akses (hanya bisa update order di company/branch sendiri)
- ✅ WebSocket notification otomatis ke semua client
- ✅ Tidak mengubah field lain (hanya status)
- ✅ Lebih cepat dari endpoint PUT /orders/:id (tidak perlu kirim semua data)

## Perbedaan dengan PUT /orders/:id

| Fitur | PATCH /orders/:id/status | PUT /orders/:id |
|-------|-------------------------|-----------------|
| Update status | ✅ | ✅ |
| Update customer info | ❌ | ✅ |
| Update order items | ❌ | ✅ |
| Recalculate total | ❌ | ✅ |
| Kecepatan | Lebih cepat | Lebih lambat |
| Request body | Minimal | Lengkap |

## Use Case

1. **Kitchen Display System**: Update status saat makanan mulai dimasak (PROCESSING)
2. **Waiter App**: Update status saat makanan siap diantar (READY)
3. **Cashier**: Update status saat order selesai dibayar (COMPLETED)
4. **Customer Service**: Cancel order jika ada masalah (CANCELLED)

## Catatan Penting

1. Endpoint ini hanya mengubah status, tidak mengubah field lainnya
2. Perubahan status akan di-broadcast ke WebSocket clients secara real-time
3. Anda bisa menambahkan validasi transisi status di service layer jika diperlukan
4. Order hanya bisa diupdate oleh user yang memiliki akses ke company dan branch yang sama

## File Terkait

- `internal/entity/order_dto.go` - DTO untuk request
- `internal/service/order_service.go` - Business logic
- `internal/handler/order_handler.go` - HTTP handler
- `internal/repository/order_repository.go` - Database operations
- `routes/routes.go` - Route definition
- `UPDATE_ORDER_STATUS_API.md` - Dokumentasi lengkap API
- `test_update_order_status.ps1` - Test script untuk Windows
- `test_update_order_status.sh` - Test script untuk Linux/Mac
