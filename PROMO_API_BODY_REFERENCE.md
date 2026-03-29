# Promo API - Body Reference

Dokumentasi lengkap untuk request body POST dan PUT promo.

---

## 📝 POST /api/promos - Create Promo

### 1️⃣ Promo Normal

**Body:**
```json
{
  "branch_id": "uuid-optional",
  "name": "Diskon Ramadan 2024",
  "code": "RAMADAN20",
  "promo_category": "normal",
  "type": "percentage",
  "value": 20,
  "max_discount": 150000,
  "min_transaction": 300000,
  "quota": 100,
  "start_date": "2024-03-01",
  "end_date": "2024-04-30",
  "is_active": true
}
```

**Field Descriptions:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `branch_id` | UUID | ❌ | Null = company level, UUID = branch level |
| `name` | String | ✅ | Nama promo (max 100 char) |
| `code` | String | ✅ | Kode promo unik (max 50 char) |
| `promo_category` | String | ✅ | Harus "normal" |
| `type` | String | ✅ | "percentage" atau "fixed" |
| `value` | Number | ✅ | Nilai diskon (% atau nominal) |
| `max_discount` | Number | ❌ | Max diskon untuk percentage type |
| `min_transaction` | Number | ❌ | Minimum transaksi |
| `quota` | Integer | ❌ | Jumlah maksimal penggunaan |
| `start_date` | String | ✅ | Format: YYYY-MM-DD |
| `end_date` | String | ✅ | Format: YYYY-MM-DD |
| `is_active` | Boolean | ❌ | Default: true |

**Minimal Body:**
```json
{
  "name": "Diskon 10%",
  "code": "DISC10",
  "promo_category": "normal",
  "type": "percentage",
  "value": 10,
  "start_date": "2024-12-01",
  "end_date": "2024-12-31"
}
```

---

### 2️⃣ Promo Product

**Body:**
```json
{
  "branch_id": "uuid-optional",
  "name": "Flash Sale Laptop 50%",
  "code": "LAPTOP50",
  "promo_category": "product",
  "type": "percentage",
  "value": 50,
  "max_discount": 5000000,
  "min_transaction": null,
  "quota": 50,
  "product_ids": [
    "550e8400-e29b-41d4-a716-446655440001",
    "550e8400-e29b-41d4-a716-446655440002"
  ],
  "start_date": "2024-12-12",
  "end_date": "2024-12-12",
  "is_active": true
}
```

**Field Descriptions:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `branch_id` | UUID | ❌ | Null = company level, UUID = branch level |
| `name` | String | ✅ | Nama promo (max 100 char) |
| `code` | String | ✅ | Kode promo unik (max 50 char) |
| `promo_category` | String | ✅ | Harus "product" |
| `type` | String | ✅ | "percentage" atau "fixed" |
| `value` | Number | ✅ | Nilai diskon (% atau nominal) |
| `max_discount` | Number | ❌ | Max diskon untuk percentage type |
| `min_transaction` | Number | ❌ | Minimum transaksi |
| `quota` | Integer | ❌ | Jumlah maksimal penggunaan |
| `product_ids` | Array[UUID] | ✅ | **WAJIB, minimal 1 product** |
| `start_date` | String | ✅ | Format: YYYY-MM-DD |
| `end_date` | String | ✅ | Format: YYYY-MM-DD |
| `is_active` | Boolean | ❌ | Default: true |

**Minimal Body:**
```json
{
  "name": "Diskon Produk",
  "code": "PROD50",
  "promo_category": "product",
  "type": "percentage",
  "value": 50,
  "product_ids": [
    "550e8400-e29b-41d4-a716-446655440001"
  ],
  "start_date": "2024-12-01",
  "end_date": "2024-12-31"
}
```

---

### 3️⃣ Promo Bundle

**Body:**
```json
{
  "branch_id": "uuid-optional",
  "name": "Paket Gaming Lengkap",
  "code": "GAMING999",
  "promo_category": "bundle",
  "type": "fixed",
  "value": 1000000,
  "max_discount": null,
  "min_transaction": null,
  "quota": 30,
  "bundle_items": [
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440001",
      "quantity": 1
    },
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440002",
      "quantity": 2
    },
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440003",
      "quantity": 1
    }
  ],
  "start_date": "2024-12-01",
  "end_date": "2024-12-31",
  "is_active": true
}
```

**Field Descriptions:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `branch_id` | UUID | ❌ | Null = company level, UUID = branch level |
| `name` | String | ✅ | Nama promo (max 100 char) |
| `code` | String | ✅ | Kode promo unik (max 50 char) |
| `promo_category` | String | ✅ | Harus "bundle" |
| `type` | String | ✅ | "percentage" atau "fixed" |
| `value` | Number | ✅ | Nilai diskon (% atau nominal) |
| `max_discount` | Number | ❌ | Max diskon untuk percentage type |
| `min_transaction` | Number | ❌ | Minimum transaksi |
| `quota` | Integer | ❌ | Jumlah maksimal penggunaan |
| `bundle_items` | Array[Object] | ✅ | **WAJIB, minimal 2 items** |
| `bundle_items[].product_id` | UUID | ✅ | ID produk |
| `bundle_items[].quantity` | Integer | ✅ | Jumlah produk (min: 1) |
| `start_date` | String | ✅ | Format: YYYY-MM-DD |
| `end_date` | String | ✅ | Format: YYYY-MM-DD |
| `is_active` | Boolean | ❌ | Default: true |

**Minimal Body:**
```json
{
  "name": "Paket Bundle",
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
      "product_id": "550e8400-e29b-41d4-a716-446655440002",
      "quantity": 1
    }
  ],
  "start_date": "2024-12-01",
  "end_date": "2024-12-31"
}
```

---

## 🔄 PUT /api/promos/:id - Update Promo

### 1️⃣ Update Promo Normal

**Body (semua field optional):**
```json
{
  "name": "Diskon Ramadan 2024 Updated",
  "code": "RAMADAN25",
  "promo_category": "normal",
  "type": "percentage",
  "value": 25,
  "max_discount": 200000,
  "min_transaction": 400000,
  "quota": 150,
  "start_date": "2024-03-01",
  "end_date": "2024-05-01",
  "is_active": false
}
```

**Update Partial (hanya field yang ingin diubah):**
```json
{
  "value": 25,
  "max_discount": 200000
}
```

---

### 2️⃣ Update Promo Product

**Body (semua field optional):**
```json
{
  "name": "Flash Sale Laptop 70%",
  "code": "LAPTOP70",
  "promo_category": "product",
  "type": "percentage",
  "value": 70,
  "max_discount": 7000000,
  "product_ids": [
    "550e8400-e29b-41d4-a716-446655440003",
    "550e8400-e29b-41d4-a716-446655440004"
  ],
  "start_date": "2024-12-15",
  "end_date": "2024-12-15",
  "is_active": true
}
```

**Update Products Only:**
```json
{
  "product_ids": [
    "550e8400-e29b-41d4-a716-446655440005",
    "550e8400-e29b-41d4-a716-446655440006"
  ]
}
```

**⚠️ Note:** Saat update `product_ids`, produk lama akan dihapus dan diganti dengan yang baru.

---

### 3️⃣ Update Promo Bundle

**Body (semua field optional):**
```json
{
  "name": "Paket Gaming Super",
  "code": "GAMING1500",
  "promo_category": "bundle",
  "type": "fixed",
  "value": 1500000,
  "bundle_items": [
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440001",
      "quantity": 2
    },
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440002",
      "quantity": 3
    }
  ],
  "start_date": "2024-12-01",
  "end_date": "2024-12-31",
  "is_active": true
}
```

**Update Bundle Items Only:**
```json
{
  "bundle_items": [
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440007",
      "quantity": 1
    },
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440008",
      "quantity": 2
    }
  ]
}
```

**⚠️ Note:** Saat update `bundle_items`, bundle items lama akan dihapus dan diganti dengan yang baru.

---

## 📋 Field Rules Summary

### Common Fields (All Categories)
| Field | POST | PUT | Validation |
|-------|------|-----|------------|
| `name` | ✅ Required | ❌ Optional | Max 100 char |
| `code` | ✅ Required | ❌ Optional | Max 50 char, unique |
| `promo_category` | ✅ Required | ❌ Optional | normal/product/bundle |
| `type` | ✅ Required | ❌ Optional | percentage/fixed |
| `value` | ✅ Required | ❌ Optional | Min: 0 |
| `max_discount` | ❌ Optional | ❌ Optional | Min: 0 (for percentage) |
| `min_transaction` | ❌ Optional | ❌ Optional | Min: 0 |
| `quota` | ❌ Optional | ❌ Optional | Min: 1 |
| `start_date` | ✅ Required | ❌ Optional | YYYY-MM-DD |
| `end_date` | ✅ Required | ❌ Optional | YYYY-MM-DD, > start_date |
| `is_active` | ❌ Optional | ❌ Optional | Boolean, default: true |
| `branch_id` | ❌ Optional | ❌ - | UUID or null |

### Category-Specific Fields

**Promo Product:**
| Field | POST | PUT | Validation |
|-------|------|-----|------------|
| `product_ids` | ✅ Required | ❌ Optional | Array[UUID], min: 1 |

**Promo Bundle:**
| Field | POST | PUT | Validation |
|-------|------|-----|------------|
| `bundle_items` | ✅ Required | ❌ Optional | Array[Object], min: 2 |
| `bundle_items[].product_id` | ✅ Required | ✅ Required | UUID |
| `bundle_items[].quantity` | ✅ Required | ✅ Required | Integer, min: 1 |

---

## ⚠️ Validation Errors

### Error: Missing promo_category
```json
{
  "error": "promo_category is required"
}
```

### Error: Invalid promo_category
```json
{
  "error": "promo_category must be 'normal', 'product', or 'bundle'"
}
```

### Error: Product promo without product_ids
```json
{
  "error": "product_ids required for product promo"
}
```

### Error: Bundle promo with < 2 items
```json
{
  "error": "bundle promo requires at least 2 products"
}
```

### Error: Invalid date format
```json
{
  "error": "invalid start_date format, use YYYY-MM-DD"
}
```

### Error: End date before start date
```json
{
  "error": "end_date must be after start_date"
}
```

---

## 📝 Complete Examples

### Example 1: Create Promo Normal (Percentage)
```bash
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Diskon Akhir Tahun 15%",
    "code": "NEWYEAR15",
    "promo_category": "normal",
    "type": "percentage",
    "value": 15,
    "max_discount": 100000,
    "min_transaction": 200000,
    "quota": 100,
    "start_date": "2024-12-01",
    "end_date": "2024-12-31",
    "is_active": true
  }'
```

### Example 2: Create Promo Normal (Fixed)
```bash
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Potongan Rp 50.000",
    "code": "POTONG50K",
    "promo_category": "normal",
    "type": "fixed",
    "value": 50000,
    "min_transaction": 300000,
    "start_date": "2024-12-01",
    "end_date": "2024-12-31"
  }'
```

### Example 3: Create Promo Product
```bash
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Flash Sale Laptop 50%",
    "code": "LAPTOP50",
    "promo_category": "product",
    "type": "percentage",
    "value": 50,
    "max_discount": 5000000,
    "product_ids": [
      "550e8400-e29b-41d4-a716-446655440001",
      "550e8400-e29b-41d4-a716-446655440002"
    ],
    "start_date": "2024-12-12",
    "end_date": "2024-12-12"
  }'
```

### Example 4: Create Promo Bundle
```bash
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Paket Gaming Lengkap",
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
        "product_id": "550e8400-e29b-41d4-a716-446655440002",
        "quantity": 2
      }
    ],
    "start_date": "2024-12-01",
    "end_date": "2024-12-31"
  }'
```

### Example 5: Update Promo (Partial)
```bash
curl -X PUT http://localhost:8080/api/promos/PROMO_UUID \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "value": 20,
    "is_active": false
  }'
```

### Example 6: Update Product IDs
```bash
curl -X PUT http://localhost:8080/api/promos/PROMO_UUID \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_ids": [
      "550e8400-e29b-41d4-a716-446655440005",
      "550e8400-e29b-41d4-a716-446655440006"
    ]
  }'
```

---

## 💡 Tips

1. **Promo Normal**: Paling sederhana, tidak perlu product_ids atau bundle_items
2. **Promo Product**: Wajib ada product_ids, minimal 1 produk
3. **Promo Bundle**: Wajib ada bundle_items, minimal 2 produk
4. **Type Percentage**: Gunakan max_discount untuk membatasi diskon maksimal
5. **Type Fixed**: Tidak perlu max_discount
6. **Update**: Hanya kirim field yang ingin diubah
7. **Product/Bundle Update**: Saat update product_ids atau bundle_items, data lama akan diganti

---

**Version**: 2.0 - Promo Categories
**Last Updated**: 2024
