# Payment Method di Response Order

## Overview
Field `payment_method` sudah tersedia di semua response order API.

## Field Payment Method di OrderResponse

```json
{
  "id": "uuid",
  "payment_method": "TUNAI|QRIS|DEBIT|CREDIT|GOPAY|OVO|COMPLIMENTARY",
  "payment_status": "UNPAID|PAID|PARTIAL",
  "paid_amount": 0,
  "change_amount": 0,
  "payment_note": "",
  "paid_at": "2024-01-01 12:00:00"
}
```

## Payment Method Options

| Method | Deskripsi | Behavior |
|--------|-----------|----------|
| `TUNAI` | Pembayaran tunai/cash | Bisa ada kembalian |
| `QRIS` | Pembayaran via QRIS | Paid amount = total |
| `DEBIT` | Kartu debit | Paid amount = total |
| `CREDIT` | Kartu kredit | Paid amount = total |
| `GOPAY` | GoPay | Paid amount = total |
| `OVO` | OVO | Paid amount = total |
| `COMPLIMENTARY` | Gratis/complimentary | **Total = 0, Paid = 0** |

### COMPLIMENTARY Payment Behavior

Ketika menggunakan payment method `COMPLIMENTARY`:
- `discount_amount` akan di-set = `subtotal_amount` (full discount)
- `tax_amount` akan di-set = 0
- `total_amount` akan di-set = 0
- `paid_amount` akan di-set = 0
- `change_amount` akan di-set = 0

Ini membuat order menjadi gratis sepenuhnya.

## Endpoint yang Mengembalikan Payment Method

### 1. Get All Orders
```bash
GET /api/v1/external/orders
```

Response:
```json
{
  "status": "success",
  "data": [
    {
      "id": "order-uuid",
      "payment_method": "TUNAI",
      "payment_status": "PAID",
      "paid_amount": 100000,
      "change_amount": 10000,
      ...
    }
  ]
}
```

### 2. Get Order by ID
```bash
GET /api/v1/external/orders/{id}
```

Response:
```json
{
  "status": "success",
  "data": {
    "id": "order-uuid",
    "payment_method": "QRIS",
    "payment_status": "PAID",
    "paid_amount": 90000,
    ...
  }
}
```

### 3. Create Order
```bash
POST /api/v1/external/orders
```

Response (payment_method akan kosong sebelum payment):
```json
{
  "status": "success",
  "data": {
    "id": "order-uuid",
    "payment_method": "",
    "payment_status": "UNPAID",
    "paid_amount": 0,
    ...
  }
}
```

### 4. Process Payment

```bash
POST /api/v1/external/orders/{id}/payment
```

Request:
```json
{
  "payment_method": "TUNAI",
  "paid_amount": 100000,
  "payment_note": "Pembayaran tunai"
}
```

Response:
```json
{
  "status": "success",
  "data": {
    "order_id": "order-uuid",
    "payment_method": "TUNAI",
    "payment_status": "PAID",
    "paid_amount": 100000,
    "change_amount": 10000,
    ...
  }
}
```

### 5. Process COMPLIMENTARY Payment

```bash
POST /api/v1/external/orders/{id}/payment
```

Request:
```json
{
  "payment_method": "COMPLIMENTARY",
  "paid_amount": 0,
  "payment_note": "Complimentary for VIP guest"
}
```

Response (semua amount jadi 0):
```json
{
  "status": "success",
  "data": {
    "order_id": "order-uuid",
    "payment_method": "COMPLIMENTARY",
    "payment_status": "PAID",
    "subtotal_amount": 50000,
    "discount_amount": 50000,
    "tax_amount": 0,
    "total_amount": 0,
    "paid_amount": 0,
    "change_amount": 0,
    "payment_note": "Complimentary for VIP guest"
  }
}
```

## Testing

### Test Payment Method Response
Jalankan test script untuk memverifikasi payment_method di response:

```powershell
./test_payment_method_response.ps1
```

### Test COMPLIMENTARY Payment
Jalankan test script untuk memverifikasi COMPLIMENTARY payment (gratis):

```powershell
./test_complimentary_payment.ps1
```

Hasil yang diharapkan untuk COMPLIMENTARY:
- Total Amount = 0
- Paid Amount = 0
- Discount Amount = Subtotal Amount (full discount)
- Tax Amount = 0

## Catatan

- Field `payment_method` akan kosong (`""`) untuk order yang belum dibayar
- Setelah proses payment, field ini akan terisi sesuai metode pembayaran yang dipilih
- Field ini selalu ada di response, baik order sudah dibayar atau belum
- **COMPLIMENTARY payment**: Membuat order gratis dengan set discount = subtotal, sehingga total = 0
- Untuk payment method selain TUNAI, paid_amount harus sama dengan total_amount
- Hanya TUNAI yang bisa memiliki change_amount (kembalian)
