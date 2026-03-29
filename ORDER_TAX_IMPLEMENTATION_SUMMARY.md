# Summary: Implementasi Perhitungan Pajak Bertingkat pada Order

## Perubahan yang Dilakukan

### 1. Entity Changes

#### `internal/entity/order.go`
- Menambahkan field `SubtotalAmount` (total item sebelum pajak)
- Menambahkan field `TaxAmount` (total semua pajak)
- Field `TotalAmount` sekarang = SubtotalAmount + TaxAmount

#### `internal/entity/order_dto.go`
- Menambahkan `SubtotalAmount`, `TaxAmount` di `OrderResponse`
- Menambahkan `TaxDetails []TaxDetailDTO` untuk breakdown perhitungan
- Membuat struct baru `TaxDetailDTO` untuk detail setiap pajak

### 2. Repository Changes

#### `internal/repository/tax_repository.go`
- Menambahkan method `FindActiveTaxesByBranch()` untuk mengambil pajak aktif
- Pajak diurutkan berdasarkan prioritas DESC (tertinggi dulu)

### 3. Service Changes

#### `internal/service/order_service.go`
- Menambahkan `taxRepo` ke struct `orderService`
- Membuat fungsi `calculateTaxes()` untuk perhitungan pajak bertingkat
- Update `CreateOrder()` untuk menghitung pajak
- Update `CreatePublicOrder()` untuk menghitung pajak
- Update `UpdateOrder()` untuk recalculate pajak saat order items berubah
- Update `toOrderResponse()` untuk menampilkan tax details

### 4. Routes Changes

#### `routes/routes.go`
- Update inisialisasi `orderService` untuk include `taxRepo`
- Memindahkan inisialisasi `taxRepo` sebelum `orderService`

### 5. Migration

#### `add_tax_fields_to_orders.go`
- Script migration untuk menambahkan kolom baru ke tabel orders
- Menambahkan `subtotal_amount` dan `tax_amount`
- Update existing orders dengan nilai default

### 6. Testing Scripts

- `test_order_tax.ps1` (Windows)
- `test_order_tax.sh` (Linux/Mac)
- `run_tax_migration.ps1` (Windows)
- `run_tax_migration.sh` (Linux/Mac)

### 7. Documentation

- `ORDER_TAX_CALCULATION.md` - Dokumentasi lengkap
- `ORDER_TAX_QUICK_START.md` - Quick start guide
- `ORDER_TAX_IMPLEMENTATION_SUMMARY.md` - Summary ini

## Cara Kerja Perhitungan Pajak

### Formula

```
Subtotal = Σ(price × quantity) untuk semua items

Untuk setiap pajak (diurutkan berdasarkan prioritas DESC):
  tax_amount = current_base × (percentage / 100)
  current_base = current_base + tax_amount

Total Tax = Σ semua tax_amount
Total Amount = Subtotal + Total Tax
```

### Contoh

Input:
- Item: 2 × Rp 50.000 = Rp 100.000
- Pajak 1: PB1 10% (prioritas 1)
- Pajak 2: Service Charge 5% (prioritas 2)

Perhitungan:
1. Subtotal: Rp 100.000
2. PB1: 100.000 × 10% = Rp 10.000 → Base: Rp 110.000
3. SC: 110.000 × 5% = Rp 5.500 → Base: Rp 115.500
4. Total: Rp 115.500

## API Changes

### Request (tidak berubah)

```json
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

### Response (ditambahkan field baru)

```json
{
  "status": "success",
  "data": {
    "id": "uuid",
    "subtotal_amount": 100000,    // NEW
    "tax_amount": 15500,          // NEW
    "total_amount": 115500,       // UPDATED (dulu hanya subtotal)
    "tax_details": [              // NEW
      {
        "tax_id": "uuid",
        "tax_name": "PB1",
        "percentage": 10,
        "priority": 1,
        "base_amount": 100000,
        "tax_amount": 10000
      },
      {
        "tax_id": "uuid",
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

## Setup Instructions

### 1. Run Migration

```bash
# Windows
.\run_tax_migration.ps1

# Linux/Mac
./run_tax_migration.sh
```

### 2. Create Taxes

Buat pajak melalui API `/api/v1/external/tax` dengan:
- `prioritas`: angka untuk urutan perhitungan (lebih tinggi = dihitung lebih dulu)
- `status`: "active" agar dihitung
- `presentase`: persentase pajak

### 3. Test

```bash
# Windows
.\test_order_tax.ps1

# Linux/Mac
./test_order_tax.sh
```

## Backward Compatibility

- Existing orders akan di-migrate dengan `subtotal_amount = total_amount` dan `tax_amount = 0`
- API request tidak berubah, hanya response yang ditambahkan field baru
- Jika tidak ada pajak aktif, order tetap bisa dibuat dengan `tax_amount = 0`

## Notes

- Pajak dihitung otomatis saat create/update order
- Hanya pajak dengan status `active` yang dihitung
- Company-level tax (branch_id = NULL) berlaku untuk semua branch
- Branch-level tax hanya berlaku untuk branch tertentu
- Prioritas menentukan urutan perhitungan (DESC = tertinggi dulu)
