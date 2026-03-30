# Add Item to Order - Summary

## Endpoint Baru

```
POST /api/v1/external/orders/quick/:id
```

Endpoint untuk menambahkan item ke pesanan yang sudah ada.

## Request

```json
{
  "product_id": "uuid-product",
  "quantity": 2,
  "note": "Extra pedas"
}
```

## Keunggulan

- ✅ Tambah item ke order existing
- ✅ Body minimal: product_id, quantity, note (opsional)
- ✅ Otomatis recalculate subtotal, diskon, pajak, total
- ✅ Validasi produk dan ketersediaan
- ✅ WebSocket broadcast real-time
- ✅ Multi-tenant support

## Use Case

Kasir/waiter bisa menambah item ke order yang sudah dibuat tanpa perlu update seluruh order.

Contoh:
1. Customer pesan Nasi Goreng (order dibuat)
2. 5 menit kemudian customer mau tambah Es Teh
3. Kasir tinggal POST ke `/orders/quick/:id` dengan body Es Teh
4. Sistem otomatis hitung ulang total

## Testing

```powershell
.\test_quick_order.ps1
```

## Files Modified

1. `internal/entity/order_dto.go` - Added `AddOrderItemRequest`
2. `internal/service/order_service.go` - Added `AddOrderItem()` method
3. `internal/handler/order_handler.go` - Added `AddOrderItem()` handler
4. `routes/routes.go` - Added route `POST /external/orders/quick/:id`

## Files Updated

1. `QUICK_ORDER_API.md` - Updated with add item documentation
2. `test_quick_order.ps1` - Updated with add item test
3. `QUICK_ORDER_SUMMARY.md` - This file
