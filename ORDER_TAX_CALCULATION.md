# Perhitungan Pajak Bertingkat pada Order

## Overview

Sistem order sekarang mendukung perhitungan pajak bertingkat berdasarkan prioritas. Pajak dihitung secara kumulatif, di mana pajak dengan prioritas lebih KECIL (1, 2, 3, ...) dihitung terlebih dahulu, dan hasilnya menjadi base untuk pajak berikutnya.

**Urutan Prioritas:**
- Priority 1 = dihitung PERTAMA
- Priority 2 = dihitung KEDUA  
- Priority 3 = dihitung KETIGA, dst.

## Cara Kerja

### Contoh Perhitungan

Misalkan:
- Total item (subtotal): Rp 100.000
- Pajak 1: Service Charge 5% (prioritas 1) ← dihitung pertama
- Pajak 2: PB1 10% (prioritas 2) ← dihitung kedua

Karena diurutkan berdasarkan prioritas ASC (1, 2, 3, ...), maka:

1. **Service Charge (prioritas 1)**: 
   - Base: Rp 100.000
   - Pajak: 100.000 × 5% = Rp 5.000
   - Total: Rp 105.000

2. **Pajak PB1 (prioritas 2)**:
   - Base: Rp 105.000 (sudah termasuk pajak sebelumnya)
   - Pajak: 105.000 × 10% = Rp 10.500
   - Total: Rp 115.500

**Total Akhir**: Rp 115.500

**Catatan Penting**: 
- Priority 1 = dihitung PERTAMA
- Priority 2 = dihitung KEDUA
- Priority 3 = dihitung KETIGA, dst.
- Angka prioritas yang LEBIH KECIL dihitung LEBIH DULU

## Struktur Data

### Order Entity

```go
type Order struct {
    // ... fields lainnya
    SubtotalAmount float64 // Total item sebelum pajak
    TaxAmount      float64 // Total semua pajak
    TotalAmount    float64 // Subtotal + Tax
}
```

### Order Response

```go
type OrderResponse struct {
    // ... fields lainnya
    SubtotalAmount float64        // Total item sebelum pajak
    TaxAmount      float64        // Total semua pajak
    TotalAmount    float64        // Subtotal + Tax
    TaxDetails     []TaxDetailDTO // Detail perhitungan setiap pajak
}

type TaxDetailDTO struct {
    TaxID      uuid.UUID // ID pajak
    TaxName    string    // Nama pajak
    Percentage float64   // Persentase pajak
    Priority   int       // Prioritas pajak
    BaseAmount float64   // Jumlah yang dikenakan pajak
    TaxAmount  float64   // Hasil perhitungan pajak
}
```

## API Response Example

### Create Order Response & Get Order By ID Response

Kedua endpoint ini mengembalikan format response yang sama dengan breakdown pajak lengkap:

```json
{
  "status": "success",
  "message": "Order created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "company_id": "123e4567-e89b-12d3-a456-426614174000",
    "branch_id": "789e0123-e89b-12d3-a456-426614174000",
    "customer_name": "John Doe",
    "table_number": "A1",
    "order_method": "DINE_IN",
    "status": "PENDING",
    "subtotal_amount": 100000,
    "tax_amount": 15500,
    "total_amount": 115500,
    "tax_details": [
      {
        "tax_id": "abc12345-e89b-12d3-a456-426614174000",
        "tax_name": "PB1",
        "percentage": 10,
        "priority": 1,
        "base_amount": 100000,
        "tax_amount": 10000
      },
      {
        "tax_id": "def67890-e89b-12d3-a456-426614174000",
        "tax_name": "Service Charge",
        "percentage": 5,
        "priority": 2,
        "base_amount": 110000,
        "tax_amount": 5500
      }
    ],
    "order_items": [
      {
        "id": "item-uuid",
        "product_id": "product-uuid",
        "product_name": "Nasi Goreng",
        "quantity": 2,
        "price": 50000,
        "subtotal": 100000
      }
    ],
    "created_at": "2024-01-15 10:30:00",
    "updated_at": "2024-01-15 10:30:00"
  }
}
```

## Konfigurasi Pajak

Untuk menggunakan fitur ini, pastikan pajak sudah dikonfigurasi dengan benar:

1. **Buat Pajak** melalui endpoint `/api/v1/external/tax`
2. **Set Prioritas**: Priority 1 = dihitung pertama, Priority 2 = dihitung kedua, dst.
3. **Set Status**: Hanya pajak dengan status `active` yang akan dihitung
4. **Scope**: 
   - Company-level tax (branch_id = NULL): berlaku untuk semua branch
   - Branch-level tax: hanya berlaku untuk branch tertentu

### Contoh Create Tax

```bash
# Service Charge (prioritas 1 - dihitung pertama)
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "Service Charge",
    "tipe_pajak": "sc",
    "presentase": 5,
    "prioritas": 1,
    "status": "active",
    "deskripsi": "Biaya Layanan"
  }'

# Pajak PB1 (prioritas 2 - dihitung kedua)
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10,
    "prioritas": 2,
    "status": "active",
    "deskripsi": "Pajak Barang dan Jasa"
  }'
```

## Migration

Untuk menambahkan field baru ke tabel orders, jalankan:

```bash
go run add_tax_fields_to_orders.go
```

Migration ini akan:
1. Menambahkan kolom `subtotal_amount` ke tabel orders
2. Menambahkan kolom `tax_amount` ke tabel orders
3. Update existing orders: set subtotal_amount = total_amount, tax_amount = 0

## Testing

### Test Create Order dengan Pajak

```bash
# Windows
.\test_order_tax.ps1

# Linux/Mac
./test_order_tax.sh
```

### Test Get Order By ID dengan Tax Breakdown

```bash
# Windows
.\test_get_order_by_id.ps1

# Linux/Mac
./test_get_order_by_id.sh
```

### Manual Test Create Order

```bash
# 1. Login dan dapatkan token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@example.com",
    "password": "password123"
  }' | jq -r '.data.token')

# 2. Create order
curl -X POST http://localhost:8080/api/v1/external/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "table_number": "A1",
    "customer_name": "John Doe",
    "order_method": "DINE_IN",
    "order_items": [
      {
        "product_id": "YOUR_PRODUCT_ID",
        "quantity": 2
      }
    ]
  }'

# 3. Get order by ID (akan menampilkan tax breakdown juga)
curl -X GET http://localhost:8080/api/v1/external/orders/ORDER_ID \
  -H "Authorization: Bearer $TOKEN"
```

## Notes

- Pajak dihitung otomatis saat create order dan update order
- Jika tidak ada pajak aktif, order akan dibuat tanpa pajak (tax_amount = 0)
- Perhitungan pajak menggunakan prioritas DESC (tertinggi dulu)
- Tax details ditampilkan di response untuk transparansi perhitungan
