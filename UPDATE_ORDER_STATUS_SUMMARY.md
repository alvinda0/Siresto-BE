# Summary: Update Order Status Implementation

## ✅ Implementasi Selesai

Endpoint untuk mengubah status order dari PENDING ke PROCESSING (atau status lainnya) telah berhasil dibuat.

## ⚠️ Perubahan Status

Status `PREPARING` telah diubah menjadi `PROCESSING` untuk konsistensi penamaan.

Jika Anda memiliki data order dengan status `PREPARING` di database, jalankan migration:
```bash
go run update_preparing_to_processing.go
```

Lihat `UPDATE_STATUS_MIGRATION.md` untuk detail lengkap.

## Perubahan File

### 1. Entity & DTO
- **File**: `internal/entity/order_dto.go`
- **Perubahan**: Menambahkan `UpdateOrderStatusRequest` struct
```go
type UpdateOrderStatusRequest struct {
    Status OrderStatus `json:"status" binding:"required"`
}
```

### 2. Repository
- **File**: `internal/repository/order_repository.go`
- **Perubahan**: 
  - Menambahkan method `UpdateStatus` di interface
  - Implementasi method `UpdateStatus` untuk update status di database

### 3. Service
- **File**: `internal/service/order_service.go`
- **Perubahan**:
  - Menambahkan method `UpdateOrderStatus` di interface
  - Implementasi business logic untuk update status dengan validasi akses

### 4. Handler
- **File**: `internal/handler/order_handler.go`
- **Perubahan**: Menambahkan handler `UpdateOrderStatus` dengan WebSocket broadcast

### 5. Routes
- **File**: `routes/routes.go`
- **Perubahan**: Menambahkan route `PATCH /api/v1/external/orders/:id/status`

## Endpoint Baru

```
PATCH /api/v1/external/orders/:id/status
```

### Request Body
```json
{
  "status": "PREPARING"
}
```

### Response
```json
{
  "status": "success",
  "message": "Order status updated successfully",
  "data": {
    "id": "...",
    "status": "PREPARING",
    ...
  }
}
```

## Status yang Tersedia

1. `PENDING` - Order baru
2. `CONFIRMED` - Order dikonfirmasi
3. `PROCESSING` - Order sedang diproses ⭐
4. `READY` - Order siap
5. `COMPLETED` - Order selesai
6. `CANCELLED` - Order dibatalkan

## Fitur

✅ Update status order dengan endpoint khusus  
✅ Validasi akses (company & branch)  
✅ WebSocket notification otomatis  
✅ Lebih cepat dari PUT /orders/:id  
✅ Request body minimal  
✅ Tidak mengubah field lain  

## Testing

### PowerShell
```powershell
$env:PRODUCT_ID = "your-product-uuid"
.\test_update_order_status.ps1
```

### Bash
```bash
export PRODUCT_ID="your-product-uuid"
./test_update_order_status.sh
```

## Dokumentasi

1. **UPDATE_ORDER_STATUS_API.md** - Dokumentasi lengkap API
2. **UPDATE_ORDER_STATUS_QUICK_START.md** - Panduan cepat
3. **test_update_order_status.ps1** - Test script Windows
4. **test_update_order_status.sh** - Test script Linux/Mac

## Flow Penggunaan

```
1. Customer pesan → Order dibuat (PENDING)
2. Kitchen terima → Update ke PROCESSING ⭐
3. Masakan selesai → Update ke READY
4. Customer terima → Update ke COMPLETED
```

## Keuntungan Endpoint Ini

1. **Lebih Cepat**: Hanya update 1 field, tidak perlu kirim semua data
2. **Lebih Aman**: Tidak bisa accidentally mengubah field lain
3. **Real-time**: WebSocket notification otomatis
4. **Simple**: Request body minimal

## Next Steps (Opsional)

Jika ingin menambahkan validasi transisi status, edit di `internal/service/order_service.go`:

```go
func (s *orderService) UpdateOrderStatus(...) {
    // Validasi transisi status
    if order.Status == entity.OrderStatusCompleted {
        return nil, errors.New("cannot update completed order")
    }
    
    if req.Status == entity.OrderStatusProcessing && 
       order.Status != entity.OrderStatusPending {
        return nil, errors.New("can only process pending orders")
    }
    
    // ... dst
}
```

## Cara Menjalankan

1. Pastikan server berjalan:
```bash
go run cmd/server/main.go
```

2. Test dengan script atau manual curl/Postman

3. Monitor WebSocket untuk melihat real-time updates

## Selesai! 🎉

Endpoint update order status sudah siap digunakan.
