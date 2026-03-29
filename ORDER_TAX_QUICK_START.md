# Quick Start: Order Tax Calculation

## Setup

### 1. Jalankan Migration

```bash
# Windows PowerShell
.\run_tax_migration.ps1

# Linux/Mac
chmod +x run_tax_migration.sh
./run_tax_migration.sh
```

### 2. Buat Pajak

Login terlebih dahulu dan buat pajak:

```bash
# Login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@example.com",
    "password": "password123"
  }'

# Simpan token dari response
TOKEN="your_token_here"

# Buat Pajak PB1 (prioritas 1 - dihitung pertama)
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10,
    "prioritas": 1,
    "status": "active",
    "deskripsi": "Pajak Barang dan Jasa 10%"
  }'

# Buat Service Charge (prioritas 2 - dihitung kedua)
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "Service Charge",
    "tipe_pajak": "sc",
    "presentase": 5,
    "prioritas": 2,
    "status": "active",
    "deskripsi": "Biaya Layanan 5%"
  }'
```

### 3. Test Order dengan Pajak

#### Test Create Order
```bash
# Windows PowerShell
.\test_order_tax.ps1

# Linux/Mac
chmod +x test_order_tax.sh
./test_order_tax.sh
```

#### Test Get Order By ID (dengan tax breakdown)
```bash
# Windows PowerShell
.\test_get_order_by_id.ps1

# Linux/Mac
chmod +x test_get_order_by_id.sh
./test_get_order_by_id.sh
```

## Contoh Response

### Create Order & Get Order By ID Response

Kedua endpoint mengembalikan format yang sama dengan breakdown pajak lengkap.

#### Order dengan 2 item @ Rp 50.000 = Rp 100.000

```json
{
  "status": "success",
  "message": "Order created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "subtotal_amount": 100000,
    "tax_amount": 15500,
    "total_amount": 115500,
    "tax_details": [
      {
        "tax_id": "abc-123",
        "tax_name": "PB1",
        "percentage": 10,
        "priority": 1,
        "base_amount": 100000,
        "tax_amount": 10000
      },
      {
        "tax_id": "def-456",
        "tax_name": "Service Charge",
        "percentage": 5,
        "priority": 2,
        "base_amount": 110000,
        "tax_amount": 5500
      }
    ],
    "order_items": [...]
  }
}
```

### Breakdown Perhitungan:

1. **Subtotal**: Rp 100.000 (2 × Rp 50.000)
2. **PB1 (10%, prioritas 1)**: 
   - Base: Rp 100.000
   - Tax: Rp 10.000
   - Running total: Rp 110.000
3. **Service Charge (5%, prioritas 2)**:
   - Base: Rp 110.000 (sudah termasuk PB1)
   - Tax: Rp 5.500
   - Running total: Rp 115.500
4. **Total Akhir**: Rp 115.500

## Tips

- Pajak dengan **prioritas lebih KECIL** dihitung **PERTAMA**
  - Priority 1 = dihitung pertama
  - Priority 2 = dihitung kedua
  - Priority 3 = dihitung ketiga, dst.
- Setiap pajak dihitung dari base yang sudah termasuk pajak sebelumnya (kumulatif)
- Hanya pajak dengan status `active` yang dihitung
- Company-level tax (branch_id = NULL) berlaku untuk semua branch
- Branch-level tax hanya berlaku untuk branch tertentu

## Troubleshooting

### Order tidak kena pajak?

Cek:
1. Apakah ada pajak dengan status `active`?
2. Apakah pajak sudah di-assign ke company/branch yang benar?

```bash
# Cek pajak aktif
curl -X GET http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN"
```

### Perhitungan pajak salah?

Pastikan prioritas pajak sudah benar:
- Priority 1 = dihitung pertama
- Priority 2 = dihitung kedua, dst.
- Contoh: Service Charge (prioritas 1) dihitung sebelum PB1 (prioritas 2)
