# Quick Start: Order dengan Promo

## 1. Jalankan Migration

```bash
go run add_promo_fields_to_orders.go
```

Output yang diharapkan:
```
Adding promo fields to orders table...
✓ Added promo_id column
✓ Added discount_amount column
✓ Updated existing orders
✓ Migration completed successfully!
```

## 2. Restart Server

```bash
go run cmd/server/main.go
```

## 3. Test Order dengan Promo

### PowerShell:
```powershell
.\test_order_with_promo.ps1
```

### Bash:
```bash
./test_order_with_promo.sh
```

## 4. Manual Test dengan Postman/cURL

### Step 1: Login
```bash
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json

{
  "email": "cashier@branch1.com",
  "password": "password123"
}
```

### Step 2: Create Promo (jika belum ada)
```bash
POST http://localhost:8080/api/v1/promos
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Diskon 10%",
  "code": "DISKON10",
  "type": "percentage",
  "value": 10,
  "max_discount": 50000,
  "min_transaction": 50000,
  "start_date": "2026-03-29",
  "end_date": "2026-04-29",
  "is_active": true
}
```

### Step 3: Get Products
```bash
GET http://localhost:8080/api/v1/products?limit=5
Authorization: Bearer {token}
```

### Step 4: Create Order dengan Promo
```bash
POST http://localhost:8080/api/v1/orders
Authorization: Bearer {token}
Content-Type: application/json

{
  "table_number": "Table-1",
  "customer_name": "John Doe",
  "order_method": "DINE_IN",
  "promo_code": "DISKON10",
  "order_items": [
    {
      "product_id": "{product_id}",
      "quantity": 5
    }
  ]
}
```

### Step 5: Get Order by ID
```bash
GET http://localhost:8080/api/v1/orders/{order_id}
Authorization: Bearer {token}
```

## Expected Response

```json
{
  "status": "success",
  "message": "Order retrieved successfully",
  "data": {
    "id": "uuid",
    "subtotal_amount": 100000,
    "discount_amount": 10000,
    "tax_amount": 13950,
    "total_amount": 103950,
    "promo_code": "DISKON10",
    "promo_id": "uuid",
    "promo_details": {
      "promo_id": "uuid",
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
        "tax_id": "uuid",
        "tax_name": "Service Charge",
        "percentage": 5,
        "priority": 1,
        "base_amount": 90000,
        "tax_amount": 4500
      },
      {
        "tax_id": "uuid",
        "tax_name": "PB1",
        "percentage": 10,
        "priority": 2,
        "base_amount": 94500,
        "tax_amount": 9450
      }
    ],
    "order_items": [...]
  }
}
```

## Calculation Breakdown

```
Subtotal:        100,000
Discount (10%):  -10,000
─────────────────────────
After Discount:   90,000

Tax Priority 1 (Service Charge 5%):
  Base: 90,000
  Tax:  4,500
─────────────────────────
After Tax 1:      94,500

Tax Priority 2 (PB1 10%):
  Base: 94,500
  Tax:  9,450
─────────────────────────
TOTAL:           103,950
```

## Troubleshooting

### Error: "promo code not found"
- Pastikan promo sudah dibuat
- Cek kode promo benar

### Error: "promo has expired"
- Cek `end_date` promo
- Update promo dengan tanggal yang valid

### Error: "minimum transaction is X"
- Subtotal order harus >= `min_transaction` promo
- Tambah quantity atau pilih product lain

### Error: "promo quota has been exhausted"
- Promo sudah mencapai quota maksimal
- Buat promo baru atau update quota

## Features

✅ Percentage discount dengan max_discount
✅ Fixed discount
✅ Minimum transaction validation
✅ Quota tracking
✅ Date range validation
✅ Tax calculation setelah discount
✅ Promo details di response
✅ Tax breakdown dengan priority
✅ Auto increment promo usage count

## Next Steps

- Test dengan berbagai tipe promo (percentage, fixed)
- Test dengan berbagai kombinasi tax priority
- Test validasi (expired, quota, min transaction)
- Integrate dengan frontend/mobile app
