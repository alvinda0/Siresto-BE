# Implementasi Order dengan Promo

## Overview

Implementasi perhitungan order dengan promo menggunakan formula:
```
Total = ((Subtotal - Discount) + Tax Priority 1) + Tax Priority 2
```

## Formula Perhitungan

### Langkah-langkah:
1. **Subtotal**: Hitung total harga semua item (quantity × price)
2. **Discount**: Apply promo jika ada (percentage atau fixed)
3. **After Discount**: Subtotal - Discount
4. **Tax Priority 1**: Hitung pajak pertama dari (Subtotal - Discount)
5. **Tax Priority 2**: Hitung pajak kedua dari (Subtotal - Discount + Tax Priority 1)
6. **Total**: After Discount + Tax Priority 1 + Tax Priority 2

### Contoh Perhitungan:

#### Tanpa Promo:
```
Subtotal: 100,000
Discount: 0
After Discount: 100,000

Service Charge 5% (Priority 1): 100,000 × 5% = 5,000
After Tax 1: 105,000

PB1 10% (Priority 2): 105,000 × 10% = 10,500
Total: 115,500
```

#### Dengan Promo 10%:
```
Subtotal: 100,000
Discount: 10,000 (10%)
After Discount: 90,000

Service Charge 5% (Priority 1): 90,000 × 5% = 4,500
After Tax 1: 94,500

PB1 10% (Priority 2): 94,500 × 10% = 9,450
Total: 103,950
```

#### Dengan Promo Fixed 15,000:
```
Subtotal: 100,000
Discount: 15,000
After Discount: 85,000

Service Charge 5% (Priority 1): 85,000 × 5% = 4,250
After Tax 1: 89,250

PB1 10% (Priority 2): 89,250 × 10% = 8,925
Total: 98,175
```

## Database Schema

### Perubahan pada Tabel `orders`:

```sql
ALTER TABLE orders 
ADD COLUMN promo_id UUID REFERENCES promos(id) ON DELETE SET NULL;

ALTER TABLE orders 
ADD COLUMN discount_amount DECIMAL(15,2) DEFAULT 0;
```

### Field Baru:
- `promo_id`: UUID (nullable) - ID promo yang digunakan
- `discount_amount`: DECIMAL(15,2) - Jumlah diskon yang didapat

## API Changes

### Create Order Request

```json
{
  "table_number": "Table-1",
  "customer_name": "John Doe",
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

### Order Response

```json
{
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
```

## Validasi Promo

Sistem akan memvalidasi:

1. **Promo Exists**: Kode promo harus ada di database
2. **Active**: Promo harus aktif (`is_active = true`)
3. **Date Range**: Tanggal sekarang harus dalam range `start_date` - `end_date`
4. **Quota**: Jika ada quota, `used_count` harus < `quota`
5. **Minimum Transaction**: Subtotal harus >= `min_transaction` (jika ada)

### Error Messages:

- "promo code not found"
- "promo is not active"
- "promo has not started yet"
- "promo has expired"
- "promo quota has been exhausted"
- "minimum transaction is X"

## Tipe Promo

### 1. Percentage Discount
```json
{
  "type": "percentage",
  "value": 10,
  "max_discount": 50000
}
```
- Diskon = Subtotal × (value / 100)
- Jika diskon > max_discount, gunakan max_discount

### 2. Fixed Discount
```json
{
  "type": "fixed",
  "value": 15000
}
```
- Diskon = value
- Jika diskon > subtotal, gunakan subtotal

## Promo Usage Tracking

Setiap kali promo digunakan:
1. Increment `used_count` di tabel `promos`
2. Simpan `promo_id` di order
3. Simpan `discount_amount` di order

## Migration

Jalankan migration untuk menambahkan field baru:

```bash
go run add_promo_fields_to_orders.go
```

## Testing

Jalankan test script:

```bash
# PowerShell
.\test_order_with_promo.ps1

# Bash
./test_order_with_promo.sh
```

Test akan:
1. Login sebagai CASHIER
2. Get products dan promos
3. Create order tanpa promo
4. Create order dengan promo
5. Verify perhitungan
6. Get order by ID untuk verify breakdown

## Code Changes

### Files Modified:
1. `internal/entity/order.go` - Added promo_id, discount_amount fields
2. `internal/entity/order_dto.go` - Added PromoDetailDTO
3. `internal/service/order_service.go` - Added promo logic
4. `internal/repository/order_repository.go` - Added Promo preload
5. `routes/routes.go` - Added PromoRepository to OrderService

### New Files:
1. `add_promo_fields_to_orders.go` - Migration script
2. `test_order_with_promo.ps1` - Test script
3. `ORDER_PROMO_IMPLEMENTATION.md` - Documentation

## Get Order by ID - Breakdown

Ketika get order by ID, response akan include:

1. **Order Items**: List semua item dengan subtotal
2. **Promo Details**: Detail promo yang digunakan (jika ada)
3. **Tax Details**: Breakdown perhitungan setiap pajak dengan priority
4. **Amounts**:
   - `subtotal_amount`: Total sebelum diskon & pajak
   - `discount_amount`: Jumlah diskon
   - `tax_amount`: Total pajak
   - `total_amount`: Grand total

## Notes

- Promo hanya bisa digunakan sekali per order
- Diskon diterapkan ke subtotal sebelum pajak
- Tax dihitung dari amount setelah diskon
- Tax priority menentukan urutan perhitungan (1 first, 2 second, dst)
- Promo usage count otomatis increment saat order dibuat
- Jika promo tidak valid, order akan gagal dengan error message yang jelas
