# 💳 Payment API - Complete Guide

## 📖 Daftar Isi

1. [Quick Start](#quick-start)
2. [Payment Methods](#payment-methods)
3. [API Endpoint](#api-endpoint)
4. [Payment Flow](#payment-flow)
5. [Promo Rules](#promo-rules)
6. [Examples](#examples)
7. [Testing](#testing)

## 🚀 Quick Start

### 1. Migration sudah dijalankan ✅
```bash
go run add_payment_fields_to_orders.go
```

### 2. Restart server
```bash
go run cmd/server/main.go
```

### 3. Endpoint siap digunakan
```
POST /api/v1/external/orders/{order_id}/payment
```

## 💳 Payment Methods

| Method | Kode | Kembalian | Validasi |
|--------|------|-----------|----------|
| QRIS | `QRIS` | ❌ | paid_amount = total_amount |
| Tunai | `TUNAI` | ✅ | paid_amount ≥ total_amount |
| Debit | `DEBIT` | ❌ | paid_amount = total_amount |
| Credit | `CREDIT` | ❌ | paid_amount = total_amount |
| GoPay | `GOPAY` | ❌ | paid_amount = total_amount |
| OVO | `OVO` | ❌ | paid_amount = total_amount |
| Complimentary | `COMPLIMENTARY` | ❌ | paid_amount = 0 |

## 🔌 API Endpoint

### Request
```http
POST /api/v1/external/orders/{order_id}/payment
Authorization: Bearer {token}
Content-Type: application/json

{
  "payment_method": "TUNAI",
  "paid_amount": 120000,
  "promo_code": "DISKON10",
  "payment_note": "Pembayaran tunai"
}
```

### Response Success
```json
{
  "status": "success",
  "message": "Payment processed successfully",
  "data": {
    "order_id": "uuid",
    "payment_method": "TUNAI",
    "payment_status": "PAID",
    "subtotal_amount": 100000,
    "discount_amount": 10000,
    "tax_amount": 9900,
    "total_amount": 99900,
    "paid_amount": 120000,
    "change_amount": 20100,
    "payment_note": "Pembayaran tunai",
    "paid_at": "2024-03-20 15:30:45"
  }
}
```

## 🔄 Payment Flow

### Flow 1: Promo di Create Order
```
1. POST /api/v1/external/orders
   Body: { promo_code: "DISKON10", ... }
   
2. POST /api/v1/external/orders/{id}/payment
   Body: { payment_method: "TUNAI", paid_amount: 120000 }
```

### Flow 2: Promo di Payment
```
1. POST /api/v1/external/orders
   Body: { ... } (tanpa promo)
   
2. POST /api/v1/external/orders/{id}/payment
   Body: { 
     payment_method: "QRIS", 
     paid_amount: 99900,
     promo_code: "DISKON10"
   }
```

### ⚠️ Aturan Promo

**PENTING:** Satu order hanya bisa menggunakan SATU promo!

- ✅ Promo di create order → Payment tanpa promo
- ✅ Tanpa promo di create order → Promo di payment
- ❌ Promo di create order → Promo yang sama di payment (DITOLAK)
- ❌ Promo di create order → Promo berbeda di payment (DITOLAK)

**Contoh Error:**
```json
// Jika promo sudah dipakai (promo yang sama)
{
  "error": "promo code has already been applied to this order"
}

// Jika coba pakai promo berbeda
{
  "error": "order already has a promo applied, cannot apply another promo"
}
```

📚 **Lihat dokumentasi lengkap:** [PAYMENT_PROMO_RULES.md](PAYMENT_PROMO_RULES.md)

## 📝 Examples

### Example 1: TUNAI (Cash)
```bash
curl -X POST http://localhost:8080/api/v1/external/orders/{order_id}/payment \
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
  "data": {
    "payment_method": "TUNAI",
    "total_amount": 99900,
    "paid_amount": 120000,
    "change_amount": 20100
  }
}
```

### Example 2: QRIS
```bash
curl -X POST http://localhost:8080/api/v1/external/orders/{order_id}/payment \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "QRIS",
    "paid_amount": 99900,
    "payment_note": "Pembayaran via QRIS"
  }'
```

### Example 3: GOPAY + Promo
```bash
curl -X POST http://localhost:8080/api/v1/external/orders/{order_id}/payment \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "GOPAY",
    "paid_amount": 89910,
    "promo_code": "DISKON10",
    "payment_note": "Pembayaran GoPay dengan promo"
  }'
```

### Example 4: COMPLIMENTARY
```bash
curl -X POST http://localhost:8080/api/v1/external/orders/{order_id}/payment \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "COMPLIMENTARY",
    "paid_amount": 0,
    "payment_note": "Complimentary untuk VIP"
  }'
```

## 🧪 Testing

### Menggunakan PowerShell Script

1. Edit `test_payment.ps1`:
   - Ganti `your-product-id-here` dengan product ID yang valid
   - Sesuaikan email/password login

2. Jalankan test:
```powershell
.\test_payment.ps1
```

### Manual Testing dengan Postman

1. **Login** untuk mendapatkan token
2. **Create Order** untuk mendapatkan order_id
3. **Process Payment** dengan berbagai metode

## ⚠️ Error Handling

### Order Not Found
```json
{
  "status": "error",
  "message": "Order not found",
  "error": "order not found"
}
```

### Already Paid
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "order has already been paid"
}
```

### Invalid Payment Method
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "invalid payment method"
}
```

### Insufficient Payment
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "paid amount is less than total amount"
}
```

## 📊 Database Changes

Migration menambahkan kolom berikut ke tabel `orders`:

```sql
payment_method   VARCHAR(50)
payment_status   VARCHAR(50) DEFAULT 'UNPAID'
paid_amount      DECIMAL(15,2) DEFAULT 0
change_amount    DECIMAL(15,2) DEFAULT 0
payment_note     TEXT
paid_at          TIMESTAMP
```

## 🎯 Features

✅ 7 metode pembayaran (QRIS, TUNAI, DEBIT, CREDIT, GOPAY, OVO, COMPLIMENTARY)
✅ Perhitungan kembalian otomatis untuk TUNAI
✅ Apply promo saat create order atau saat payment
✅ Validasi payment amount sesuai metode
✅ Auto complete order setelah payment
✅ WebSocket broadcast untuk real-time update
✅ Tax calculation dengan promo
✅ Payment history dengan timestamp

## 📚 Dokumentasi Lengkap

- `PAYMENT_API.md` - Dokumentasi API lengkap
- `PAYMENT_QUICK_START.md` - Quick start guide
- `PAYMENT_IMPLEMENTATION_SUMMARY.md` - Summary implementasi
- `test_payment.ps1` - Test script

## 🎉 Selesai!

Endpoint payment sudah siap digunakan dengan semua fitur yang diminta!
