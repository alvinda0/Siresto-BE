# Checklist: Implementasi Pajak Order

## ✅ Code Changes

- [x] Update `internal/entity/order.go` - tambah field `SubtotalAmount` dan `TaxAmount`
- [x] Update `internal/entity/order_dto.go` - tambah `TaxDetailDTO` dan update `OrderResponse`
- [x] Update `internal/repository/tax_repository.go` - tambah method `FindActiveTaxesByBranch()`
- [x] Update `internal/service/order_service.go` - tambah `calculateTaxes()` dan update semua fungsi order
- [x] Update `routes/routes.go` - tambah `taxRepo` ke `orderService`
- [x] Compile berhasil tanpa error

## ✅ Migration

- [x] Buat file `add_tax_fields_to_orders.go`
- [x] Buat script `run_tax_migration.ps1` (Windows)
- [x] Buat script `run_tax_migration.sh` (Linux/Mac)

## ✅ Testing Scripts

- [x] Buat `test_order_tax.ps1` (Windows) - Test create order
- [x] Buat `test_order_tax.sh` (Linux/Mac) - Test create order
- [x] Buat `test_get_order_by_id.ps1` (Windows) - Test get order by ID
- [x] Buat `test_get_order_by_id.sh` (Linux/Mac) - Test get order by ID

## ✅ Documentation

- [x] `ORDER_TAX_CALCULATION.md` - Dokumentasi lengkap cara kerja
- [x] `ORDER_TAX_QUICK_START.md` - Quick start guide
- [x] `ORDER_TAX_EXAMPLES.md` - Contoh perhitungan berbagai skenario
- [x] `ORDER_TAX_IMPLEMENTATION_SUMMARY.md` - Summary perubahan
- [x] `ORDER_TAX_CHECKLIST.md` - Checklist ini

## 📋 Next Steps (Manual)

### 1. Run Migration
```bash
# Windows
.\run_tax_migration.ps1

# Linux/Mac
./run_tax_migration.sh
```

### 2. Restart Server
```bash
go run cmd/server/main.go
```

### 3. Create Taxes

Login dan buat pajak:

```bash
# Login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@example.com",
    "password": "password123"
  }'

# Buat PB1 (prioritas 1)
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10,
    "prioritas": 1,
    "status": "active"
  }'

# Buat Service Charge (prioritas 2)
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "Service Charge",
    "tipe_pajak": "sc",
    "presentase": 5,
    "prioritas": 2,
    "status": "active"
  }'
```

### 4. Test Order

#### Test Create Order
```bash
# Windows
.\test_order_tax.ps1

# Linux/Mac
./test_order_tax.sh
```

#### Test Get Order By ID
```bash
# Windows
.\test_get_order_by_id.ps1

# Linux/Mac
./test_get_order_by_id.sh
```

### 5. Verify Response

Pastikan response order memiliki:
- `subtotal_amount`: Total item sebelum pajak
- `tax_amount`: Total semua pajak
- `total_amount`: Subtotal + Tax
- `tax_details`: Array berisi detail setiap pajak

## 🔍 Verification Points

### Database
- [ ] Tabel `orders` memiliki kolom `subtotal_amount`
- [ ] Tabel `orders` memiliki kolom `tax_amount`
- [ ] Existing orders sudah di-update dengan nilai default

### API Response
- [ ] Create order mengembalikan `subtotal_amount`, `tax_amount`, `total_amount`
- [ ] Create order mengembalikan `tax_details` array
- [ ] Get order by ID mengembalikan breakdown pajak lengkap
- [ ] Get all orders mengembalikan breakdown pajak untuk setiap order
- [ ] Update order recalculate pajak dengan benar
- [ ] Get order menampilkan breakdown pajak

### Calculation
- [ ] Pajak dihitung berdasarkan prioritas (DESC)
- [ ] Setiap pajak dihitung dari base yang sudah termasuk pajak sebelumnya
- [ ] Total amount = subtotal + sum(all taxes)
- [ ] Tax details menampilkan base_amount dan tax_amount untuk setiap pajak

### Edge Cases
- [ ] Order tanpa pajak aktif (tax_amount = 0)
- [ ] Order dengan 1 pajak
- [ ] Order dengan multiple pajak
- [ ] Update order items recalculate pajak

## 🐛 Troubleshooting

### Migration Error
```
Error: column already exists
```
**Solution**: Kolom sudah ada, skip migration atau drop kolom dulu

### Order tidak kena pajak
**Check**:
1. Apakah ada pajak dengan status `active`?
2. Apakah pajak di company/branch yang benar?
3. Cek dengan: `GET /api/v1/external/tax`

### Perhitungan salah
**Check**:
1. Urutan prioritas pajak (DESC = tertinggi dulu)
2. Persentase pajak sudah benar?
3. Lihat `tax_details` di response untuk debug

### Compile Error
**Check**:
1. `go mod tidy` untuk update dependencies
2. Pastikan semua import sudah benar
3. Cek dengan: `go build cmd/server/main.go`

## 📚 References

- [ORDER_TAX_CALCULATION.md](ORDER_TAX_CALCULATION.md) - Dokumentasi lengkap
- [ORDER_TAX_QUICK_START.md](ORDER_TAX_QUICK_START.md) - Quick start
- [ORDER_TAX_EXAMPLES.md](ORDER_TAX_EXAMPLES.md) - Contoh perhitungan
- [TAX_API.md](TAX_API.md) - API dokumentasi pajak
