# Fix: Urutan Prioritas Pajak

## Masalah yang Ditemukan

Dari response order:
```json
{
  "subtotal_amount": 0,  ← Seharusnya 75000
  "tax_amount": 0,
  "total_amount": 75000,
  "tax_details": [
    {
      "tax_name": "PB1",
      "priority": 2,  ← Dihitung duluan (salah!)
      "base_amount": 0
    },
    {
      "tax_name": "Service Charge",
      "priority": 1,  ← Seharusnya dihitung duluan
      "base_amount": 0
    }
  ]
}
```

### Masalah:
1. **subtotal_amount = 0** - Order lama belum di-migrate
2. **Urutan prioritas terbalik** - Priority 1 seharusnya dihitung pertama, bukan priority 2

## Solusi yang Diterapkan

### 1. Update Migration Script

File: `add_tax_fields_to_orders.go`

```sql
UPDATE orders 
SET subtotal_amount = total_amount, 
    tax_amount = 0 
WHERE subtotal_amount = 0 OR subtotal_amount IS NULL
```

Ini akan update semua existing orders dengan:
- `subtotal_amount` = `total_amount` (karena dulu belum ada pajak)
- `tax_amount` = 0

### 2. Ubah Urutan Sorting Prioritas

**Sebelum:**
```go
// Order by priority DESC (highest first)
Order("prioritas DESC")
```
- Priority 2 dihitung pertama
- Priority 1 dihitung kedua

**Sesudah:**
```go
// Order by priority ASC (lowest first)
Order("prioritas ASC")
```
- Priority 1 dihitung pertama ✓
- Priority 2 dihitung kedua ✓
- Priority 3 dihitung ketiga ✓

### 3. Update Dokumentasi

Semua dokumentasi diupdate untuk menjelaskan:
- **Priority 1 = dihitung PERTAMA**
- **Priority 2 = dihitung KEDUA**
- **Priority 3 = dihitung KETIGA**, dst.

## Cara Kerja Sekarang

### Contoh: Order Rp 100.000

Dengan pajak:
- Service Charge 5% (priority 1)
- PB1 10% (priority 2)

**Perhitungan:**
```
1. Subtotal: Rp 100.000

2. Service Charge (priority 1):
   Base: 100.000
   Tax: 100.000 × 5% = 5.000
   Total: 105.000

3. PB1 (priority 2):
   Base: 105.000
   Tax: 105.000 × 10% = 10.500
   Total: 115.500

Final: Rp 115.500
```

## Testing

### 1. Jalankan Migration

```bash
go run add_tax_fields_to_orders.go
```

Ini akan update existing orders dengan subtotal_amount yang benar.

### 2. Restart Server

```bash
go run cmd/server/main.go
```

### 3. Test Get Order By ID

```bash
# Windows
.\test_get_order_by_id.ps1

# Linux/Mac
./test_get_order_by_id.sh
```

### Expected Result

```json
{
  "subtotal_amount": 75000,  ✓ Sudah benar
  "tax_amount": 7875,        ✓ Dihitung dengan benar
  "total_amount": 82875,     ✓ Subtotal + Tax
  "tax_details": [
    {
      "tax_name": "Service Charge",
      "priority": 1,           ✓ Dihitung pertama
      "base_amount": 75000,
      "tax_amount": 3750
    },
    {
      "tax_name": "PB1",
      "priority": 2,           ✓ Dihitung kedua
      "base_amount": 78750,
      "tax_amount": 7875
    }
  ]
}
```

## Perubahan File

### Code Changes
1. `internal/repository/tax_repository.go` - Ubah `ORDER BY prioritas DESC` → `ASC`
2. `internal/service/order_service.go` - Update comment di `calculateTaxes()`
3. `add_tax_fields_to_orders.go` - Update migration untuk handle NULL values

### Documentation Updates
1. `ORDER_TAX_CALCULATION.md`
2. `ORDER_TAX_QUICK_START.md`
3. `ORDER_TAX_EXAMPLES.md`
4. `ORDER_TAX_FLOW.md`

## Catatan Penting

### Untuk Existing Orders
- Setelah migration, `subtotal_amount` akan sama dengan `total_amount`
- `tax_amount` akan 0 (karena dulu belum ada fitur pajak)
- Tax breakdown akan dihitung ulang setiap kali order ditampilkan

### Untuk New Orders
- `subtotal_amount` = total item sebelum pajak
- `tax_amount` = total semua pajak
- `total_amount` = subtotal + tax
- Tax breakdown dihitung dan disimpan di response

## Rekomendasi Prioritas

Untuk restoran pada umumnya:

```
Priority 1: Service Charge (5%)
  → Dihitung dari subtotal murni

Priority 2: PB1/VAT (10%)
  → Dihitung dari subtotal + service charge

Priority 3: Pajak tambahan (jika ada)
  → Dihitung dari total sebelumnya
```

Tapi ini bisa disesuaikan dengan kebutuhan bisnis masing-masing!
