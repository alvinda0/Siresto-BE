# Quick Start - Promo Categories

Panduan cepat untuk menggunakan 3 jenis promo: Normal, Product, dan Bundle.

## 1. Jalankan Migration

```bash
go run add_promo_category_and_tables.go
```

Output yang diharapkan:
```
✓ Added promo_category column to promos table
✓ Created promo_products table
✓ Created promo_bundles table
✓ Created indexes
✅ Migration completed successfully!
```

## 2. Restart Server

```bash
go run cmd/server/main.go
```

## 3. Test Promo Categories

```powershell
.\test_promo_categories.ps1
```

## Perbedaan 3 Jenis Promo

### 🎯 Promo Normal
**Untuk apa:** Diskon umum untuk semua produk

**Contoh kasus:**
- Diskon Ramadan 20%
- Flash Sale Akhir Tahun
- Diskon Member Baru

**Request:**
```json
{
  "promo_category": "normal",
  "name": "Diskon Ramadan",
  "code": "RAMADAN20",
  "type": "percentage",
  "value": 20
}
```

### 🎯 Promo Product
**Untuk apa:** Diskon hanya untuk produk tertentu

**Contoh kasus:**
- Diskon 50% untuk Laptop Gaming
- Clearance Sale Produk Tertentu
- Promo Produk Baru

**Request:**
```json
{
  "promo_category": "product",
  "name": "Flash Sale Laptop",
  "code": "LAPTOP50",
  "type": "percentage",
  "value": 50,
  "product_ids": [
    "uuid-product-1",
    "uuid-product-2"
  ]
}
```

**Response akan include:**
```json
{
  "products": [
    {
      "product_id": "...",
      "product_name": "Laptop ASUS ROG",
      
    }
  ]
}
```

### 🎯 Promo Bundle
**Untuk apa:** Diskon saat beli kombinasi produk dengan jumlah tertentu

**Contoh kasus:**
- Beli Laptop + Mouse + Keyboard dapat diskon Rp 1.000.000
- Paket Hemat: 2 Baju + 1 Celana diskon 30%
- Bundle Gaming: PC + Monitor + Keyboard

**Request:**
```json
{
  "promo_category": "bundle",
  "name": "Paket Gaming",
  "code": "GAMING999",
  "type": "fixed",
  "value": 1000000,
  "bundle_items": [
    {
      "product_id": "uuid-laptop",
      "quantity": 1
    },
    {
      "product_id": "uuid-mouse",
      "quantity": 2
    }
  ]
}
```

**Response akan include:**
```json
{
  "bundle_items": [
    {
      "product_id": "...",
      "product_name": "Laptop ASUS",
      
      "quantity": 1
    },
    {
      "product_id": "...",
      "product_name": "Mouse Gaming",
      
      "quantity": 2
    }
  ]
}
```

## Field Wajib per Kategori

### Normal
```json
{
  "promo_category": "normal",
  "name": "...",
  "code": "...",
  "type": "percentage|fixed",
  "value": 0,
  "start_date": "YYYY-MM-DD",
  "end_date": "YYYY-MM-DD"
}
```

### Product
```json
{
  "promo_category": "product",
  "name": "...",
  "code": "...",
  "type": "percentage|fixed",
  "value": 0,
  "product_ids": ["uuid1", "uuid2"],  // WAJIB, min 1
  "start_date": "YYYY-MM-DD",
  "end_date": "YYYY-MM-DD"
}
```

### Bundle
```json
{
  "promo_category": "bundle",
  "name": "...",
  "code": "...",
  "type": "percentage|fixed",
  "value": 0,
  "bundle_items": [  // WAJIB, min 2
    {
      "product_id": "uuid",
      "quantity": 1
    }
  ],
  "start_date": "YYYY-MM-DD",
  "end_date": "YYYY-MM-DD"
}
```

## Testing Manual dengan cURL

### 1. Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@company1.com",
    "password": "password123"
  }'
```

### 2. Create Promo Normal
```bash
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Diskon Akhir Tahun",
    "code": "NEWYEAR2024",
    "promo_category": "normal",
    "type": "percentage",
    "value": 15,
    "max_discount": 100000,
    "start_date": "2024-12-01",
    "end_date": "2024-12-31"
  }'
```

### 3. Create Promo Product
```bash
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Flash Sale Laptop",
    "code": "LAPTOP50",
    "promo_category": "product",
    "type": "percentage",
    "value": 50,
    "product_ids": ["PRODUCT_UUID_1", "PRODUCT_UUID_2"],
    "start_date": "2024-12-12",
    "end_date": "2024-12-12"
  }'
```

### 4. Create Promo Bundle
```bash
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Paket Gaming",
    "code": "GAMING999",
    "promo_category": "bundle",
    "type": "fixed",
    "value": 1000000,
    "bundle_items": [
      {
        "product_id": "PRODUCT_UUID_1",
        "quantity": 1
      },
      {
        "product_id": "PRODUCT_UUID_2",
        "quantity": 2
      }
    ],
    "start_date": "2024-12-01",
    "end_date": "2024-12-31"
  }'
```

### 5. Get All Promos
```bash
curl -X GET "http://localhost:8080/api/promos?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 6. Get Promo by ID
```bash
curl -X GET "http://localhost:8080/api/promos/PROMO_UUID" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 7. Update Promo Product (change products)
```bash
curl -X PUT "http://localhost:8080/api/promos/PROMO_UUID" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_ids": ["NEW_PRODUCT_UUID_1", "NEW_PRODUCT_UUID_2"]
  }'
```

## Error Messages

### Promo Product tanpa product_ids
```json
{
  "error": "product_ids required for product promo"
}
```

### Promo Bundle dengan < 2 produk
```json
{
  "error": "bundle promo requires at least 2 products"
}
```

### Invalid promo_category
```json
{
  "error": "promo_category must be 'normal', 'product', or 'bundle'"
}
```

## Tips

1. **Promo Normal** paling sederhana, cocok untuk diskon umum
2. **Promo Product** untuk targeting produk tertentu
3. **Promo Bundle** untuk mendorong penjualan kombinasi produk
4. Gunakan `type: "percentage"` untuk diskon persen, `type: "fixed"` untuk diskon nominal
5. Set `max_discount` untuk promo percentage agar tidak unlimited
6. Set `min_transaction` untuk minimum pembelian
7. Set `quota` untuk membatasi jumlah penggunaan promo

## Next Steps

Setelah promo dibuat, Anda bisa:
1. Menggunakan promo di order dengan field `promo_code`
2. Melihat statistik penggunaan promo di field `used_count`
3. Monitoring remaining quota di field `remaining_quota`
4. Filter promo berdasarkan kategori saat listing
