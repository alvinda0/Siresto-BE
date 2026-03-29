# Promo Categories - 3 Jenis Promo

Sistem promo sekarang mendukung 3 kategori promo yang berbeda:

## 1. Promo Normal
Promo umum yang berlaku untuk semua produk dalam transaksi.

### Karakteristik:
- Berlaku untuk semua produk
- Tidak ada batasan produk tertentu
- Diskon diterapkan ke total transaksi

### Contoh Request:
```json
{
  "name": "Diskon Akhir Tahun",
  "code": "NEWYEAR2024",
  "promo_category": "normal",
  "type": "percentage",
  "value": 15,
  "max_discount": 100000,
  "min_transaction": 200000,
  "start_date": "2024-12-01",
  "end_date": "2024-12-31",
  "is_active": true
}
```

## 2. Promo Product
Promo yang hanya berlaku untuk produk-produk tertentu.

### Karakteristik:
- Hanya berlaku untuk produk yang ditentukan
- Bisa untuk 1 atau lebih produk
- Diskon hanya diterapkan ke produk yang terdaftar

### Contoh Request:
```json
{
  "name": "Diskon Laptop Gaming",
  "code": "LAPTOP50",
  "promo_category": "product",
  "type": "percentage",
  "value": 50,
  "max_discount": 5000000,
  "product_ids": [
    "550e8400-e29b-41d4-a716-446655440001",
    "550e8400-e29b-41d4-a716-446655440002"
  ],
  "start_date": "2024-12-01",
  "end_date": "2024-12-31",
  "is_active": true
}
```

### Response:
```json
{
  "id": "...",
  "name": "Diskon Laptop Gaming",
  "code": "LAPTOP50",
  "promo_category": "product",
  "type": "percentage",
  "value": 50,
  "products": [
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440001",
      "product_name": "Laptop ASUS ROG"
    },
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440002",
      "product_name": "Laptop MSI Gaming"
    }
  ]
}
```

## 3. Promo Bundle
Promo bundle dimana customer harus membeli kombinasi produk tertentu dengan jumlah tertentu.

### Karakteristik:
- Memerlukan minimal 2 produk
- Setiap produk memiliki quantity yang harus dipenuhi
- Diskon diterapkan jika semua kondisi bundle terpenuhi

### Contoh Request:
```json
{
  "name": "Paket Hemat Laptop + Mouse",
  "code": "BUNDLE100",
  "promo_category": "bundle",
  "type": "fixed",
  "value": 500000,
  "bundle_items": [
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440001",
      "quantity": 1
    },
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440003",
      "quantity": 2
    }
  ],
  "start_date": "2024-12-01",
  "end_date": "2024-12-31",
  "is_active": true
}
```

### Response:
```json
{
  "id": "...",
  "name": "Paket Hemat Laptop + Mouse",
  "code": "BUNDLE100",
  "promo_category": "bundle",
  "type": "fixed",
  "value": 500000,
  "bundle_items": [
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440001",
      "product_name": "Laptop ASUS ROG",
      "quantity": 1
    },
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440003",
      "product_name": "Mouse Gaming",
      "quantity": 2
    }
  ]
}
```

## Validasi

### Promo Normal:
- Tidak memerlukan `product_ids` atau `bundle_items`

### Promo Product:
- WAJIB menyertakan `product_ids` (minimal 1 produk)
- Tidak boleh ada `bundle_items`

### Promo Bundle:
- WAJIB menyertakan `bundle_items` (minimal 2 produk)
- Setiap item harus memiliki `product_id` dan `quantity`
- Tidak boleh ada `product_ids`

## Database Schema

### Table: promos
```sql
ALTER TABLE promos 
ADD COLUMN promo_category VARCHAR(20) NOT NULL DEFAULT 'normal';
```

### Table: promo_products
```sql
CREATE TABLE promo_products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  promo_id UUID NOT NULL REFERENCES promos(id) ON DELETE CASCADE,
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(promo_id, product_id)
);
```

### Table: promo_bundles
```sql
CREATE TABLE promo_bundles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  promo_id UUID NOT NULL REFERENCES promos(id) ON DELETE CASCADE,
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  quantity INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(promo_id, product_id)
);
```

## Migration

Jalankan migration untuk menambahkan field dan tabel baru:

```bash
go run add_promo_category_and_tables.go
```

## API Endpoints

Semua endpoint promo yang sudah ada tetap sama, hanya menambahkan field baru:

- `POST /api/promos` - Create promo (dengan promo_category)
- `PUT /api/promos/:id` - Update promo (dengan promo_category)
- `GET /api/promos` - List promos (menampilkan products/bundle_items)
- `GET /api/promos/:id` - Get promo detail (menampilkan products/bundle_items)
- `DELETE /api/promos/:id` - Delete promo

## Contoh Penggunaan

### 1. Create Promo Normal
```bash
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Diskon Ramadan",
    "code": "RAMADAN2024",
    "promo_category": "normal",
    "type": "percentage",
    "value": 20,
    "max_discount": 150000,
    "min_transaction": 300000,
    "start_date": "2024-03-01",
    "end_date": "2024-04-30"
  }'
```

### 2. Create Promo Product
```bash
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Flash Sale Elektronik",
    "code": "FLASH50",
    "promo_category": "product",
    "type": "percentage",
    "value": 50,
    "product_ids": [
      "550e8400-e29b-41d4-a716-446655440001",
      "550e8400-e29b-41d4-a716-446655440002"
    ],
    "start_date": "2024-12-12",
    "end_date": "2024-12-12"
  }'
```

### 3. Create Promo Bundle
```bash
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Paket Lengkap Gaming",
    "code": "GAMING999",
    "promo_category": "bundle",
    "type": "fixed",
    "value": 1000000,
    "bundle_items": [
      {
        "product_id": "550e8400-e29b-41d4-a716-446655440001",
        "quantity": 1
      },
      {
        "product_id": "550e8400-e29b-41d4-a716-446655440003",
        "quantity": 1
      },
      {
        "product_id": "550e8400-e29b-41d4-a716-446655440004",
        "quantity": 1
      }
    ],
    "start_date": "2024-12-01",
    "end_date": "2024-12-31"
  }'
```

## Notes

1. Promo category tidak bisa diubah setelah dibuat (untuk menjaga konsistensi data)
2. Saat update promo product/bundle, product_ids/bundle_items yang lama akan dihapus dan diganti dengan yang baru
3. Promo bundle memerlukan minimal 2 produk untuk valid
4. Promo product memerlukan minimal 1 produk untuk valid
