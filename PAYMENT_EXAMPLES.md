# Payment API - Contoh Penggunaan Praktis

## Skenario 1: Customer Makan di Restoran (Dine-In)

### Step 1: Buat Order
```json
POST /api/v1/external/orders
{
  "table_number": "A5",
  "order_method": "DINE_IN",
  "customer_name": "Budi Santoso",
  "customer_phone": "081234567890",
  "order_items": [
    {
      "product_id": "uuid-nasi-goreng",
      "quantity": 2,
      "note": "Pedas level 3"
    },
    {
      "product_id": "uuid-es-teh",
      "quantity": 2
    }
  ]
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "order-uuid-123",
    "table_number": "A5",
    "subtotal_amount": 100000,
    "tax_amount": 11000,
    "total_amount": 111000,
    "payment_status": "UNPAID"
  }
}
```

### Step 2: Bayar dengan Tunai
```json
POST /api/v1/external/orders/order-uuid-123/payment
{
  "payment_method": "TUNAI",
  "paid_amount": 150000,
  "payment_note": "Pembayaran tunai"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "order_id": "order-uuid-123",
    "payment_method": "TUNAI",
    "payment_status": "PAID",
    "total_amount": 111000,
    "paid_amount": 150000,
    "change_amount": 39000,
    "paid_at": "2024-03-20 15:30:45"
  }
}
```

## Skenario 2: Customer Pakai Promo

### Step 1: Buat Order dengan Promo
```json
POST /api/v1/external/orders
{
  "table_number": "B3",
  "order_method": "DINE_IN",
  "customer_name": "Siti Nurhaliza",
  "promo_code": "DISKON20",
  "order_items": [
    {
      "product_id": "uuid-paket-hemat",
      "quantity": 1
    }
  ]
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "order-uuid-456",
    "subtotal_amount": 150000,
    "discount_amount": 30000,
    "tax_amount": 13200,
    "total_amount": 133200,
    "promo_details": {
      "promo_name": "Diskon 20%",
      "promo_code": "DISKON20",
      "discount_amount": 30000
    }
  }
}
```

### Step 2: Bayar dengan QRIS
```json
POST /api/v1/external/orders/order-uuid-456/payment
{
  "payment_method": "QRIS",
  "paid_amount": 133200,
  "payment_note": "Scan QRIS berhasil"
}
```

## Skenario 3: Lupa Pakai Promo, Apply saat Payment

### Step 1: Buat Order (tanpa promo)
```json
POST /api/v1/external/orders
{
  "table_number": "C1",
  "order_method": "DINE_IN",
  "customer_name": "Ahmad Yani",
  "order_items": [
    {
      "product_id": "uuid-ayam-bakar",
      "quantity": 3
    }
  ]
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "order-uuid-789",
    "subtotal_amount": 200000,
    "tax_amount": 22000,
    "total_amount": 222000
  }
}
```

### Step 2: Bayar dengan Promo
```json
POST /api/v1/external/orders/order-uuid-789/payment
{
  "payment_method": "GOPAY",
  "paid_amount": 199800,
  "promo_code": "DISKON10",
  "payment_note": "Pembayaran GoPay dengan promo"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "order_id": "order-uuid-789",
    "payment_method": "GOPAY",
    "subtotal_amount": 200000,
    "discount_amount": 20000,
    "tax_amount": 19800,
    "total_amount": 199800,
    "paid_amount": 199800,
    "promo_details": {
      "promo_name": "Diskon 10%",
      "promo_code": "DISKON10",
      "discount_amount": 20000
    }
  }
}
```

## Skenario 4: VIP Customer (Complimentary)

### Step 1: Buat Order
```json
POST /api/v1/external/orders
{
  "table_number": "VIP-1",
  "order_method": "DINE_IN",
  "customer_name": "Direktur Utama",
  "notes": "VIP Guest",
  "order_items": [
    {
      "product_id": "uuid-steak",
      "quantity": 1
    }
  ]
}
```

### Step 2: Bayar Complimentary
```json
POST /api/v1/external/orders/order-uuid-vip/payment
{
  "payment_method": "COMPLIMENTARY",
  "paid_amount": 0,
  "payment_note": "Complimentary untuk VIP"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "payment_method": "COMPLIMENTARY",
    "payment_status": "PAID",
    "total_amount": 500000,
    "paid_amount": 0,
    "change_amount": 0
  }
}
```

## Skenario 5: Take Away dengan Kartu Debit

### Step 1: Buat Order
```json
POST /api/v1/external/orders
{
  "table_number": "TA-01",
  "order_method": "TAKE_AWAY",
  "customer_name": "Rina Wijaya",
  "customer_phone": "082345678901",
  "order_items": [
    {
      "product_id": "uuid-paket-1",
      "quantity": 2
    },
    {
      "product_id": "uuid-paket-2",
      "quantity": 1
    }
  ]
}
```

### Step 2: Bayar dengan Debit
```json
POST /api/v1/external/orders/order-uuid-ta/payment
{
  "payment_method": "DEBIT",
  "paid_amount": 175000,
  "payment_note": "Kartu Debit BCA"
}
```

## Skenario 6: Delivery dengan OVO

### Step 1: Buat Order
```json
POST /api/v1/external/orders
{
  "table_number": "DEL-123",
  "order_method": "DELIVERY",
  "customer_name": "Dewi Lestari",
  "customer_phone": "083456789012",
  "notes": "Alamat: Jl. Sudirman No. 123",
  "order_items": [
    {
      "product_id": "uuid-pizza",
      "quantity": 2
    }
  ]
}
```

### Step 2: Bayar dengan OVO
```json
POST /api/v1/external/orders/order-uuid-del/payment
{
  "payment_method": "OVO",
  "paid_amount": 250000,
  "payment_note": "Pembayaran via OVO"
}
```

## Skenario 7: Split Bill (Multiple Payments)

Untuk split bill, buat order terpisah untuk setiap customer:

### Customer 1
```json
POST /api/v1/external/orders
{
  "table_number": "D5",
  "order_method": "DINE_IN",
  "customer_name": "Customer 1",
  "notes": "Split bill 1/3",
  "order_items": [
    { "product_id": "uuid-item-1", "quantity": 1 }
  ]
}
```

### Customer 2
```json
POST /api/v1/external/orders
{
  "table_number": "D5",
  "order_method": "DINE_IN",
  "customer_name": "Customer 2",
  "notes": "Split bill 2/3",
  "order_items": [
    { "product_id": "uuid-item-2", "quantity": 1 }
  ]
}
```

Kemudian masing-masing bayar dengan metode yang berbeda.

## Tips Penggunaan

### 1. Validasi Promo Sebelum Payment
```json
GET /api/v1/external/promos/validate/DISKON10
```

### 2. Cek Order Detail Sebelum Payment
```json
GET /api/v1/external/orders/{order_id}
```

### 3. Filter Order by Payment Status
```json
GET /api/v1/external/orders?payment_status=UNPAID
```

## Error Handling

### Promo Tidak Valid
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "promo error: promo code not found"
}
```

### Promo Sudah Dipakai (Promo yang Sama)
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "promo error: promo code has already been applied to this order"
}
```

### Promo Sudah Dipakai (Promo Berbeda)
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "promo error: order already has a promo applied, cannot apply another promo"
}
```

### Pembayaran Kurang
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "paid amount is less than total amount"
}
```

### Order Sudah Dibayar
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "order has already been paid"
}
```

## Best Practices

1. **Selalu validasi promo** sebelum apply
2. **Cek total amount** sebelum payment
3. **Gunakan TUNAI** jika customer bayar lebih
4. **Gunakan non-cash** untuk exact amount
5. **Simpan payment_note** untuk audit trail
6. **Handle error** dengan baik di frontend
7. **Show change_amount** untuk TUNAI payment


## Skenario 8: Error - Promo Duplikat (Promo yang Sama)

### Step 1: Buat Order dengan Promo WEEKEND10
```json
POST /api/v1/external/orders
{
  "table_number": "E1",
  "order_method": "DINE_IN",
  "customer_name": "Customer Error 1",
  "promo_code": "WEEKEND10",
  "order_items": [
    {
      "product_id": "uuid-item",
      "quantity": 2
    }
  ]
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "order-uuid-error1",
    "promo_code": "WEEKEND10",
    "discount_amount": 20000,
    "total_amount": 100000
  }
}
```

### Step 2: Coba Bayar dengan Promo WEEKEND10 Lagi (DITOLAK!)
```json
POST /api/v1/external/orders/order-uuid-error1/payment
{
  "payment_method": "TUNAI",
  "paid_amount": 100000,
  "promo_code": "WEEKEND10"
}
```

**Response Error:**
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "promo error: promo code has already been applied to this order"
}
```

### Step 3: Bayar Tanpa Promo (BERHASIL)
```json
POST /api/v1/external/orders/order-uuid-error1/payment
{
  "payment_method": "TUNAI",
  "paid_amount": 100000
}
```

**Response Success:**
```json
{
  "status": "success",
  "data": {
    "payment_status": "PAID",
    "total_amount": 100000,
    "paid_amount": 100000
  }
}
```

## Skenario 9: Error - Promo Berbeda

### Step 1: Buat Order dengan Promo WEEKEND10
```json
POST /api/v1/external/orders
{
  "table_number": "F2",
  "order_method": "DINE_IN",
  "customer_name": "Customer Error 2",
  "promo_code": "WEEKEND10",
  "order_items": [
    {
      "product_id": "uuid-item",
      "quantity": 3
    }
  ]
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "order-uuid-error2",
    "promo_code": "WEEKEND10",
    "discount_amount": 30000,
    "total_amount": 150000
  }
}
```

### Step 2: Coba Bayar dengan Promo DISKON20 (DITOLAK!)
```json
POST /api/v1/external/orders/order-uuid-error2/payment
{
  "payment_method": "QRIS",
  "paid_amount": 120000,
  "promo_code": "DISKON20"
}
```

**Response Error:**
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "promo error: order already has a promo applied, cannot apply another promo"
}
```

### Step 3: Bayar Tanpa Promo (BERHASIL)
```json
POST /api/v1/external/orders/order-uuid-error2/payment
{
  "payment_method": "QRIS",
  "paid_amount": 150000
}
```

## Skenario 10: Valid - Promo Hanya di Payment

### Step 1: Buat Order TANPA Promo
```json
POST /api/v1/external/orders
{
  "table_number": "G3",
  "order_method": "DINE_IN",
  "customer_name": "Customer Valid",
  "order_items": [
    {
      "product_id": "uuid-item",
      "quantity": 2
    }
  ]
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "order-uuid-valid",
    "promo_code": "",
    "discount_amount": 0,
    "total_amount": 120000
  }
}
```

### Step 2: Bayar dengan Promo WEEKEND10 (BERHASIL)
```json
POST /api/v1/external/orders/order-uuid-valid/payment
{
  "payment_method": "GOPAY",
  "paid_amount": 100000,
  "promo_code": "WEEKEND10"
}
```

**Response Success:**
```json
{
  "status": "success",
  "data": {
    "payment_status": "PAID",
    "discount_amount": 20000,
    "total_amount": 100000,
    "paid_amount": 100000,
    "promo_details": {
      "promo_code": "WEEKEND10",
      "discount_amount": 20000
    }
  }
}
```

## Ringkasan Aturan Promo

| Skenario | Create Order | Payment | Status |
|----------|--------------|---------|--------|
| Promo di awal | WEEKEND10 | - | ✅ Valid |
| Promo di akhir | - | WEEKEND10 | ✅ Valid |
| Promo duplikat | WEEKEND10 | WEEKEND10 | ❌ Error: promo sudah dipakai |
| Promo berbeda | WEEKEND10 | DISKON20 | ❌ Error: order sudah punya promo |
| Tanpa promo | - | - | ✅ Valid |

**Kesimpulan:** Satu order hanya bisa menggunakan SATU promo, baik di awal (create order) atau di akhir (payment).
