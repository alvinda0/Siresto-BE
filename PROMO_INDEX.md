# Promo System - Index Dokumentasi

## 📚 Dokumentasi Lengkap

### 1. Quick Start
**File: `PROMO_CATEGORIES_QUICK_START.md`**
- Panduan cepat untuk memulai
- Contoh request untuk 3 jenis promo
- Testing manual dengan cURL
- Tips dan best practices

👉 **Mulai dari sini jika baru pertama kali**

### 2. Complete Documentation
**File: `PROMO_CATEGORIES.md`**
- Penjelasan detail 3 kategori promo
- Database schema lengkap
- Validasi rules
- API endpoints
- Contoh penggunaan lengkap

👉 **Baca ini untuk pemahaman mendalam**

### 3. Implementation Summary
**File: `PROMO_CATEGORIES_IMPLEMENTATION.md`**
- File-file yang dimodifikasi
- Database changes
- API changes
- Business logic
- Testing guide
- Performance & security considerations

👉 **Untuk developer yang ingin tahu detail implementasi**

## 🚀 Getting Started

### Step 1: Run Migration
```powershell
.\run_promo_category_migration.ps1
```

Atau manual:
```bash
go run add_promo_category_and_tables.go
```

### Step 2: Restart Server
```bash
go run cmd/server/main.go
```

### Step 3: Test
```powershell
.\test_promo_categories.ps1
```

## 📋 3 Jenis Promo

### 1️⃣ Promo Normal
Diskon umum untuk semua produk
```json
{
  "promo_category": "normal",
  "type": "percentage",
  "value": 20
}
```

### 2️⃣ Promo Product
Diskon untuk produk tertentu
```json
{
  "promo_category": "product",
  "type": "percentage",
  "value": 50,
  "product_ids": ["uuid1", "uuid2"]
}
```

### 3️⃣ Promo Bundle
Diskon untuk kombinasi produk
```json
{
  "promo_category": "bundle",
  "type": "fixed",
  "value": 1000000,
  "bundle_items": [
    {"product_id": "uuid1", "quantity": 1},
    {"product_id": "uuid2", "quantity": 2}
  ]
}
```

## 🗂️ File Structure

### Migration Files
- `add_promo_category_and_tables.go` - Migration script
- `run_promo_category_migration.ps1` - Migration runner

### Code Files
- `internal/entity/promo.go` - Entity definitions
- `internal/entity/promo_dto.go` - Request/Response DTOs
- `internal/repository/promo_repository.go` - Database operations
- `internal/service/promo_service.go` - Business logic
- `internal/handler/promo_handler.go` - HTTP handlers (no changes)

### Documentation Files
- `PROMO_INDEX.md` - This file
- `PROMO_CATEGORIES.md` - Complete documentation
- `PROMO_CATEGORIES_QUICK_START.md` - Quick start guide
- `PROMO_CATEGORIES_IMPLEMENTATION.md` - Implementation details

### Testing Files
- `test_promo_categories.ps1` - Automated testing script

## 🔑 Key Features

✅ 3 kategori promo: Normal, Product, Bundle
✅ Multi-tenant support
✅ Branch-level access control
✅ Backward compatible
✅ Cascade delete
✅ Optimized with indexes
✅ Complete validation
✅ Detailed error messages

## 📊 Database Tables

### promos
- Added: `promo_category` column

### promo_products (NEW)
- Links promo to specific products
- For product category promos

### promo_bundles (NEW)
- Links promo to product bundles
- Includes quantity per product
- For bundle category promos

## 🎯 Use Cases

### Promo Normal
- Diskon Ramadan 20%
- Flash Sale Akhir Tahun
- Member Baru Diskon 10%
- Diskon Hari Kemerdekaan

### Promo Product
- Clearance Sale Laptop 50%
- Diskon Produk Baru
- Flash Sale Elektronik
- Promo Kategori Tertentu

### Promo Bundle
- Paket Gaming: PC + Monitor + Keyboard
- Paket Hemat: 2 Baju + 1 Celana
- Bundle Laptop + Accessories
- Paket Lengkap Office

## 🧪 Testing

### Automated Test
```powershell
.\test_promo_categories.ps1
```

Test akan:
1. Login sebagai OWNER
2. Get products untuk testing
3. Create promo normal
4. Create promo product
5. Create promo bundle
6. Get all promos
7. Get detail setiap promo
8. Update promo product

### Manual Test
Lihat `PROMO_CATEGORIES_QUICK_START.md` untuk contoh cURL

## ⚠️ Important Notes

1. **Promo Product** memerlukan minimal 1 product_id
2. **Promo Bundle** memerlukan minimal 2 bundle_items
3. Saat update promo product/bundle, data lama akan dihapus dan diganti
4. Cascade delete otomatis menghapus promo_products dan promo_bundles
5. Backward compatible - promo lama otomatis jadi category 'normal'

## 🔄 Migration Status

Run migration untuk:
- ✅ Add promo_category column
- ✅ Create promo_products table
- ✅ Create promo_bundles table
- ✅ Create indexes

## 📞 API Endpoints

Semua endpoint sama seperti sebelumnya:

- `POST /api/promos` - Create promo
- `GET /api/promos` - List promos
- `GET /api/promos/:id` - Get promo detail
- `PUT /api/promos/:id` - Update promo
- `DELETE /api/promos/:id` - Delete promo

Yang berubah hanya request/response format untuk mendukung 3 kategori.

## 🎓 Learning Path

1. **Pemula**: Baca `PROMO_CATEGORIES_QUICK_START.md`
2. **Intermediate**: Baca `PROMO_CATEGORIES.md`
3. **Advanced**: Baca `PROMO_CATEGORIES_IMPLEMENTATION.md`
4. **Testing**: Jalankan `test_promo_categories.ps1`

## 💡 Tips

- Gunakan promo normal untuk diskon umum
- Gunakan promo product untuk targeting spesifik
- Gunakan promo bundle untuk cross-selling
- Set max_discount untuk promo percentage
- Set min_transaction untuk minimum pembelian
- Set quota untuk membatasi penggunaan

## 🚦 Status

✅ **Ready for Production**

- All code implemented
- Migration ready
- Documentation complete
- Testing script ready
- No breaking changes
- Backward compatible

## 📝 Changelog

### Version 2.0 - Promo Categories
- Added promo_category field (normal/product/bundle)
- Added promo_products table
- Added promo_bundles table
- Updated entity, repository, service layers
- Added comprehensive documentation
- Added testing scripts
- Maintained backward compatibility

---

**Need Help?**
- Quick Start: `PROMO_CATEGORIES_QUICK_START.md`
- Full Docs: `PROMO_CATEGORIES.md`
- Implementation: `PROMO_CATEGORIES_IMPLEMENTATION.md`
