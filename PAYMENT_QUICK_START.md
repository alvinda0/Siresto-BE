# Payment API - Quick Start Guide

## 🚀 Quick Start

### 1. Jalankan Migration

```bash
go run add_payment_fields_to_orders.go
```

### 2. Restart Server

```bash
go run cmd/server/main.go
```

### 3. Test Payment

Edit `test_payment.ps1` dan ganti:
- `your-product-id-here` dengan product ID yang valid
- Email/password login sesuai dengan data Anda

Jalankan test:
```powershell
.\test_payment.ps1
```

## 📋 Payment Methods

| Method | Kode | Deskripsi |
|--------|------|-----------|
| QRIS | `QRIS` | Pembayaran via QRIS |
| Tunai | `TUNAI` | Pembayaran cash |
| Debit | `DEBIT` | Kartu debit |
| Credit | `CREDIT` | Kartu kredit |
| GoPay | `GOPAY` | GoPay |
| OVO | `OVO` | OVO |
| Complimentary | `COMPLIMENTARY` | Gratis |

## 🔄 Payment Flow

### Skenario 1: Promo di Awal
```
1. Create Order + Promo Code
2. Process Payment
```

### Skenario 2: Promo di Akhir
```
1. Create Order (tanpa promo)
2. Process Payment + Promo Code
```

## 💡 Contoh Request

### Payment TUNAI (dengan kembalian)
```json
POST /api/v1/external/orders/{order_id}/payment
{
  "payment_method": "TUNAI",
  "paid_amount": 120000,
  "payment_note": "Pembayaran tunai"
}
```

### Payment QRIS
```json
POST /api/v1/external/orders/{order_id}/payment
{
  "payment_method": "QRIS",
  "paid_amount": 99900,
  "payment_note": "Pembayaran via QRIS"
}
```

### Payment dengan Promo
```json
POST /api/v1/external/orders/{order_id}/payment
{
  "payment_method": "GOPAY",
  "paid_amount": 89910,
  "promo_code": "DISKON10",
  "payment_note": "Pembayaran GoPay dengan promo"
}
```

### Payment COMPLIMENTARY
```json
POST /api/v1/external/orders/{order_id}/payment
{
  "payment_method": "COMPLIMENTARY",
  "paid_amount": 0,
  "payment_note": "Complimentary untuk VIP"
}
```

## ✅ Response Success

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

## ⚠️ Validasi

### TUNAI
- ✅ `paid_amount` boleh > `total_amount`
- ✅ Ada kembalian (`change_amount`)

### Non-Cash (QRIS, DEBIT, CREDIT, GOPAY, OVO)
- ❌ `paid_amount` harus = `total_amount`
- ❌ Tidak ada kembalian

### COMPLIMENTARY
- ✅ `paid_amount` = 0
- ✅ Gratis

## 📊 Order Status Changes

Setelah payment berhasil:
- `payment_status`: `UNPAID` → `PAID`
- `status`: `PENDING` → `COMPLETED`
- `paid_at`: Diisi timestamp

## 🔗 Related Endpoints

- `POST /api/v1/external/orders` - Create order
- `POST /api/v1/external/orders/quick` - Quick order
- `GET /api/v1/external/orders/{id}` - Get order detail
- `POST /api/v1/external/orders/{id}/payment` - Process payment

## 📖 Full Documentation

Lihat `PAYMENT_API.md` untuk dokumentasi lengkap.
