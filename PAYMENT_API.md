# Payment API Documentation

## Endpoint: Process Payment

Endpoint untuk memproses pembayaran order dengan berbagai metode pembayaran.

### Endpoint
```
POST /api/v1/external/orders/{order_id}/payment
```

### Headers
```
Authorization: Bearer {token}
Content-Type: application/json
```

### Payment Methods
- `QRIS` - Pembayaran via QRIS
- `TUNAI` - Pembayaran tunai (cash)
- `DEBIT` - Kartu debit
- `CREDIT` - Kartu kredit
- `GOPAY` - GoPay
- `OVO` - OVO
- `COMPLIMENTARY` - Gratis/complimentary

### Request Body

```json
{
  "payment_method": "QRIS",
  "paid_amount": 115500,
  "promo_code": "DISKON10",
  "payment_note": "Pembayaran via QRIS"
}
```

#### Field Descriptions:
- `payment_method` (required): Metode pembayaran yang digunakan
- `paid_amount` (required): Jumlah yang dibayarkan
- `promo_code` (optional): Kode promo yang akan diaplikasikan saat payment (jika belum diaplikasikan saat create order)
- `payment_note` (optional): Catatan pembayaran

### Response Success (200 OK)

```json
{
  "status": "success",
  "message": "Payment processed successfully",
  "data": {
    "order_id": "123e4567-e89b-12d3-a456-426614174000",
    "payment_method": "QRIS",
    "payment_status": "PAID",
    "subtotal_amount": 100000,
    "discount_amount": 10000,
    "tax_amount": 9900,
    "total_amount": 99900,
    "paid_amount": 99900,
    "change_amount": 0,
    "payment_note": "Pembayaran via QRIS",
    "paid_at": "2024-03-20 15:30:45",
    "promo_details": {
      "promo_id": "456e7890-e89b-12d3-a456-426614174001",
      "promo_name": "Diskon 10%",
      "promo_code": "DISKON10",
      "promo_type": "percentage",
      "promo_value": 10,
      "discount_amount": 10000,
      "max_discount": 50000,
      "min_transaction": 50000
    },
    "tax_details": [
      {
        "tax_id": "789e0123-e89b-12d3-a456-426614174002",
        "tax_name": "Service Charge",
        "percentage": 5,
        "priority": 1,
        "base_amount": 90000,
        "tax_amount": 4500
      },
      {
        "tax_id": "012e3456-e89b-12d3-a456-426614174003",
        "tax_name": "PB1",
        "percentage": 10,
        "priority": 2,
        "base_amount": 94500,
        "tax_amount": 9450
      }
    ]
  }
}
```

### Response Error

#### Order Not Found (404)
```json
{
  "status": "error",
  "message": "Order not found",
  "error": "order not found"
}
```

#### Already Paid (400)
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "order has already been paid"
}
```

#### Invalid Payment Method (400)
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "invalid payment method"
}
```

#### Insufficient Payment (400)
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "paid amount is less than total amount"
}
```

#### Promo Already Applied - Same Promo (400)
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "promo error: promo code has already been applied to this order"
}
```

## Payment Flow

### Flow 1: Apply Promo saat Create Order, Payment di Akhir

1. **Create Order dengan Promo**
```bash
POST /api/v1/external/orders
{
  "table_number": "A1",
  "order_method": "DINE_IN",
  "promo_code": "DISKON10",
  "order_items": [
    {
      "product_id": "uuid",
      "quantity": 2
    }
  ]
}
```

2. **Process Payment (tanpa promo lagi)**
```bash
POST /api/v1/external/orders/{order_id}/payment
{
  "payment_method": "QRIS",
  "paid_amount": 99900
}
```

**PENTING:** Jika promo sudah dipakai saat create order, tidak bisa apply promo lagi saat payment (baik promo yang sama maupun berbeda).

### Flow 2: Apply Promo saat Payment

1. **Create Order tanpa Promo**
```bash
POST /api/v1/external/orders
{
  "table_number": "A1",
  "order_method": "DINE_IN",
  "order_items": [
    {
      "product_id": "uuid",
      "quantity": 2
    }
  ]
}
```

2. **Process Payment dengan Promo**
```bash
POST /api/v1/external/orders/{order_id}/payment
{
  "payment_method": "TUNAI",
  "paid_amount": 120000,
  "promo_code": "DISKON10"
}
```

**PENTING:** Promo hanya bisa diapply sekali per order, baik di awal (create order) atau di akhir (payment).

## Payment Method Rules

### TUNAI (Cash)
- `paid_amount` boleh lebih besar dari `total_amount`
- Sistem akan menghitung `change_amount` (kembalian)
- Contoh: Total 99.900, Bayar 100.000, Kembalian 100

### Non-Cash (QRIS, DEBIT, CREDIT, GOPAY, OVO)
- `paid_amount` harus sama dengan `total_amount`
- Tidak ada kembalian

### COMPLIMENTARY
- `paid_amount` otomatis 0
- Tidak ada kembalian
- Untuk order gratis/complimentary

## Order Status Changes

Setelah payment berhasil:
- `payment_status` berubah dari `UNPAID` menjadi `PAID`
- `status` otomatis berubah menjadi `COMPLETED`
- `paid_at` diisi dengan timestamp pembayaran

## Examples

### Example 1: Payment dengan TUNAI (Cash)

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/external/orders/123e4567-e89b-12d3-a456-426614174000/payment \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "TUNAI",
    "paid_amount": 120000,
    "payment_note": "Pembayaran tunai"
  }'
```

**Response:**
```json
{
  "status": "success",
  "message": "Payment processed successfully",
  "data": {
    "order_id": "123e4567-e89b-12d3-a456-426614174000",
    "payment_method": "TUNAI",
    "payment_status": "PAID",
    "total_amount": 115500,
    "paid_amount": 120000,
    "change_amount": 4500
  }
}
```

### Example 2: Payment dengan QRIS + Promo

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/external/orders/123e4567-e89b-12d3-a456-426614174000/payment \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "QRIS",
    "paid_amount": 99900,
    "promo_code": "DISKON10",
    "payment_note": "Pembayaran via QRIS dengan promo"
  }'
```

### Example 3: Payment COMPLIMENTARY

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/external/orders/123e4567-e89b-12d3-a456-426614174000/payment \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "COMPLIMENTARY",
    "paid_amount": 0,
    "payment_note": "Complimentary untuk VIP"
  }'
```

## Migration

Jalankan migration untuk menambahkan payment fields:

```bash
go run add_payment_fields_to_orders.go
```

Migration akan menambahkan kolom:
- `payment_method` - Metode pembayaran
- `payment_status` - Status pembayaran (UNPAID/PAID/PARTIAL)
- `paid_amount` - Jumlah yang dibayar
- `change_amount` - Kembalian
- `payment_note` - Catatan pembayaran
- `paid_at` - Timestamp pembayaran


## Promo Rules

### ⚠️ PENTING: Aturan Penggunaan Promo

1. **Promo yang Sama Tidak Bisa Dipakai 2x**
   - Jika order sudah pakai promo WEEKEND10 di awal
   - Tidak bisa apply promo WEEKEND10 lagi saat payment
   - Error: "promo code has already been applied to this order"

2. **Promo Berbeda Boleh Menggantikan**
   - Jika order pakai promo WEEKEND10 di awal
   - Bisa diganti dengan promo LEBARAN2026 saat payment
   - Promo lama akan diganti dengan promo baru
   - Usage count promo lama akan dikurangi

3. **Promo Bisa di Awal atau di Akhir**
   - Apply promo saat create order, atau
   - Apply promo saat payment, atau
   - Ganti promo saat payment (jika promo berbeda)

### Contoh Skenario

#### ✅ VALID: Promo di Create Order
```json
// Step 1: Create order dengan promo
POST /api/v1/external/orders
{
  "promo_code": "WEEKEND10",
  ...
}

// Step 2: Payment tanpa promo
POST /api/v1/external/orders/{id}/payment
{
  "payment_method": "TUNAI",
  "paid_amount": 100000
}
```

#### ✅ VALID: Promo di Payment
```json
// Step 1: Create order tanpa promo
POST /api/v1/external/orders
{
  ...
}

// Step 2: Payment dengan promo
POST /api/v1/external/orders/{id}/payment
{
  "payment_method": "QRIS",
  "paid_amount": 90000,
  "promo_code": "WEEKEND10"
}
```

#### ✅ VALID: Ganti Promo Berbeda
```json
// Step 1: Create order dengan promo WEEKEND10
POST /api/v1/external/orders
{
  "promo_code": "WEEKEND10",
  ...
}

// Step 2: Payment dengan promo LEBARAN2026 (BERHASIL - promo diganti)
POST /api/v1/external/orders/{id}/payment
{
  "payment_method": "GOPAY",
  "paid_amount": 85000,
  "promo_code": "LEBARAN2026"  // ✅ Promo WEEKEND10 diganti dengan LEBARAN2026
}
```

#### ❌ INVALID: Promo Duplikat
```json
// Step 1: Create order dengan promo WEEKEND10
POST /api/v1/external/orders
{
  "promo_code": "WEEKEND10",
  ...
}

// Step 2: Payment dengan promo WEEKEND10 lagi (DITOLAK!)
POST /api/v1/external/orders/{id}/payment
{
  "payment_method": "TUNAI",
  "paid_amount": 100000,
  "promo_code": "WEEKEND10"  // ❌ Error: promo yang sama sudah dipakai
}
```

### Best Practice

1. **Tentukan kapan apply promo:**
   - Jika customer sudah tahu promo dari awal → apply saat create order
   - Jika customer baru ingat promo saat bayar → apply saat payment
   - Jika customer mau ganti promo → apply promo baru saat payment

2. **Validasi promo sebelum create order:**
   ```bash
   GET /api/v1/external/promos/validate/{promo_code}
   ```

3. **Cek order detail sebelum payment:**
   ```bash
   GET /api/v1/external/orders/{order_id}
   ```
   Lihat promo apa yang sedang dipakai

4. **Handle error dengan baik:**
   - Tampilkan pesan error yang jelas ke user
   - Jika promo sama, tampilkan: "Promo ini sudah digunakan"
   - Jika mau ganti promo, tampilkan konfirmasi: "Ganti promo WEEKEND10 dengan LEBARAN2026?"
