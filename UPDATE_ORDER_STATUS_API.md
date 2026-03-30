# Update Order Status API

Endpoint untuk mengubah status order dari PENDING ke PROCESSING atau status lainnya.

## Endpoint

```
PATCH /api/v1/external/orders/:id/status
```

## Authentication

Memerlukan token JWT di header:
```
Authorization: Bearer <token>
```

## Request Body

```json
{
  "status": "PROCESSING"
}
```

### Status yang Valid

- `PENDING` - Order baru dibuat
- `CONFIRMED` - Order dikonfirmasi
- `PROCESSING` - Order sedang diproses/disiapkan
- `READY` - Order siap
- `COMPLETED` - Order selesai
- `CANCELLED` - Order dibatalkan

## Response Success (200)

```json
{
  "status": "success",
  "message": "Order status updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "company_id": "123e4567-e89b-12d3-a456-426614174000",
    "branch_id": "123e4567-e89b-12d3-a456-426614174001",
    "customer_name": "John Doe",
    "customer_phone": "081234567890",
    "table_number": "A1",
    "notes": "Extra pedas",
    "order_method": "DINE_IN",
    "status": "PROCESSING",
    "subtotal_amount": 50000,
    "discount_amount": 5000,
    "tax_amount": 4950,
    "total_amount": 49950,
    "order_items": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440000",
        "product_id": "770e8400-e29b-41d4-a716-446655440000",
        "product_name": "Nasi Goreng",
        "quantity": 2,
        "price": 25000,
        "subtotal": 50000,
        "note": "Extra pedas"
      }
    ],
    "tax_details": [
      {
        "tax_id": "880e8400-e29b-41d4-a716-446655440000",
        "tax_name": "PB1",
        "percentage": 10,
        "priority": 1,
        "base_amount": 45000,
        "tax_amount": 4500
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
  "error": "Key: 'UpdateOrderStatusRequest.Status' Error:Field validation for 'Status' failed on the 'required' tag"
}
```

### 401 Unauthorized
```json
{
  "status": "error",
  "message": "Company ID not found",
  "error": ""
}
```

### 404 Not Found
```json
{
  "status": "error",
  "message": "Failed to update order status",
  "error": "order not found"
}
```

## Contoh Penggunaan

### cURL
```bash
curl -X PATCH http://localhost:8080/api/v1/external/orders/550e8400-e29b-41d4-a716-446655440000/status \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "PROCESSING"
  }'
```

### PowerShell
```powershell
$token = "YOUR_TOKEN"
$orderId = "550e8400-e29b-41d4-a716-446655440000"

$body = @{
    status = "PROCESSING"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/external/orders/$orderId/status" `
    -Method PATCH `
    -Headers @{
        "Authorization" = "Bearer $token"
        "Content-Type" = "application/json"
    } `
    -Body $body
```

## Catatan

1. Endpoint ini hanya mengubah status order, tidak mengubah field lainnya
2. Order hanya bisa diupdate oleh user yang memiliki akses ke company dan branch yang sama
3. Perubahan status akan di-broadcast ke WebSocket clients
4. Anda bisa menambahkan validasi transisi status di service layer jika diperlukan (misalnya: PENDING hanya bisa ke CONFIRMED atau CANCELLED)

## WebSocket Notification

Setelah status diupdate, akan dikirim notifikasi WebSocket dengan format:
```json
{
  "action": "status_updated",
  "data": {
    // OrderResponse object
  }
}
```
