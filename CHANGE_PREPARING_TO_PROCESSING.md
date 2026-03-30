# Perubahan: PREPARING → PROCESSING

## Summary

Status order `PREPARING` telah diubah menjadi `PROCESSING` sesuai permintaan.

## File yang Diubah

### 1. Entity (Code)
- ✅ `internal/entity/order.go` - Constant `OrderStatusPreparing` → `OrderStatusProcessing`

### 2. Dokumentasi
- ✅ `UPDATE_ORDER_STATUS_API.md` - Update semua referensi PREPARING → PROCESSING
- ✅ `UPDATE_ORDER_STATUS_QUICK_START.md` - Update contoh dan flow
- ✅ `UPDATE_ORDER_STATUS_SUMMARY.md` - Update summary

### 3. Test Scripts
- ✅ `test_update_order_status.ps1` - Update test case
- ✅ `test_update_order_status.sh` - Update test case

### 4. Migration
- ✅ `update_preparing_to_processing.go` - Script untuk update data existing
- ✅ `UPDATE_STATUS_MIGRATION.md` - Panduan migration

## Status Order Terbaru

```
PENDING → PROCESSING → READY → COMPLETED
         ↓
      CANCELLED
```

1. `PENDING` - Order baru dibuat
2. `CONFIRMED` - Order dikonfirmasi
3. `PROCESSING` - Order sedang diproses ⭐ (BARU)
4. `READY` - Order siap
5. `COMPLETED` - Order selesai
6. `CANCELLED` - Order dibatalkan

## Cara Menggunakan

### Request Body
```json
{
  "status": "PROCESSING"
}
```

### Contoh cURL
```bash
curl -X PATCH http://localhost:8080/api/v1/external/orders/{order_id}/status \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "PROCESSING"}'
```

## Migration (Jika Ada Data Lama)

Jika database sudah ada order dengan status `PREPARING`:

```bash
go run update_preparing_to_processing.go
```

## Testing

```powershell
# Windows
$env:PRODUCT_ID = "your-product-uuid"
.\test_update_order_status.ps1
```

```bash
# Linux/Mac
export PRODUCT_ID="your-product-uuid"
./test_update_order_status.sh
```

## Selesai! ✅

Semua file sudah diupdate dan siap digunakan dengan status `PROCESSING`.
