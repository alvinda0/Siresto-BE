# Testing Category API

## Setup
1. Jalankan server: `go run cmd/server/main.go`
2. Login sebagai OWNER untuk mendapatkan token
3. Gunakan Company ID dari seeder

## Test Scenarios

### 1. Login sebagai OWNER
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@restaurant.com",
    "password": "owner123"
  }'
```

Simpan `token` dan `company_id` dari response.

---

### 2. Get All Main Categories
```bash
curl -X GET "http://localhost:8080/api/v1/external/categories/company/{COMPANY_ID}" \
  -H "Authorization: Bearer {TOKEN}"
```

Expected: Dapat list kategori Makanan (position: 1) dan Minuman (position: 2)

---

### 3. Get Sub-Categories Makanan
```bash
# Ganti {MAKANAN_ID} dengan ID kategori Makanan dari step 2
curl -X GET "http://localhost:8080/api/v1/external/categories/company/{COMPANY_ID}?parent_id={MAKANAN_ID}" \
  -H "Authorization: Bearer {TOKEN}"
```

Expected: Dapat list:
- Makanan Utama (position: 1)
- Makanan Pembuka (position: 2)
- Makanan Penutup (position: 3)

---

### 4. Create New Main Category
```bash
curl -X POST http://localhost:8080/api/v1/external/categories \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "{COMPANY_ID}",
    "name": "Snack",
    "description": "Kategori snack dan cemilan"
  }'
```

Expected: Category baru dengan position: 3 (otomatis)

---

### 5. Create Sub-Category
```bash
curl -X POST http://localhost:8080/api/v1/external/categories \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "{COMPANY_ID}",
    "parent_id": "{SNACK_ID}",
    "name": "Snack Manis",
    "description": "Kue, cookies, dll"
  }'
```

Expected: Sub-category baru dengan position: 1

---

### 6. Get Category Detail
```bash
curl -X GET "http://localhost:8080/api/v1/external/categories/{CATEGORY_ID}" \
  -H "Authorization: Bearer {TOKEN}"
```

Expected: Detail kategori dengan sub_categories (jika ada)

---

### 7. Update Category
```bash
curl -X PUT "http://localhost:8080/api/v1/external/categories/{CATEGORY_ID}" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Snack Updated",
    "description": "Deskripsi baru",
    "position": 1,
    "is_active": true
  }'
```

Expected: Category berhasil diupdate

---

### 8. Reorder Main Categories
```bash
# Tukar posisi Minuman dan Makanan
curl -X POST "http://localhost:8080/api/v1/external/categories/company/{COMPANY_ID}/reorder" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "category_ids": [
      "{MINUMAN_ID}",
      "{MAKANAN_ID}",
      "{SNACK_ID}"
    ]
  }'
```

Expected: 
- Minuman jadi position: 1
- Makanan jadi position: 2
- Snack jadi position: 3

---

### 9. Reorder Sub-Categories
```bash
# Reorder sub-categories Makanan
curl -X POST "http://localhost:8080/api/v1/external/categories/company/{COMPANY_ID}/reorder?parent_id={MAKANAN_ID}" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "category_ids": [
      "{MAKANAN_PENUTUP_ID}",
      "{MAKANAN_UTAMA_ID}",
      "{MAKANAN_PEMBUKA_ID}"
    ]
  }'
```

Expected:
- Makanan Penutup jadi position: 1
- Makanan Utama jadi position: 2
- Makanan Pembuka jadi position: 3

---

### 10. Try Delete Category with Sub-Categories
```bash
curl -X DELETE "http://localhost:8080/api/v1/external/categories/{MAKANAN_ID}" \
  -H "Authorization: Bearer {TOKEN}"
```

Expected: Error "cannot delete category with sub-categories"

---

### 11. Delete Sub-Category
```bash
curl -X DELETE "http://localhost:8080/api/v1/external/categories/{MAKANAN_PEMBUKA_ID}" \
  -H "Authorization: Bearer {TOKEN}"
```

Expected: Sub-category berhasil dihapus

---

### 12. Now Delete Main Category
```bash
# Setelah semua sub-categories dihapus
curl -X DELETE "http://localhost:8080/api/v1/external/categories/{SNACK_ID}" \
  -H "Authorization: Bearer {TOKEN}"
```

Expected: Main category berhasil dihapus

---

## Validation Tests

### Test 1: Position dimulai dari 1
- Create category tanpa set position
- Check response, position harus >= 1, tidak boleh 0

### Test 2: Auto increment position
- Create 3 categories tanpa set position
- Check positions: 1, 2, 3 (atau lanjutan dari yang sudah ada)

### Test 3: Parent validation
- Try create sub-category dengan parent_id dari company lain
- Expected: Error "parent category must belong to the same company"

### Test 4: Self-parent validation
- Try update category dengan parent_id = category.id
- Expected: Error "category cannot be its own parent"

### Test 5: Delete protection
- Try delete category yang masih punya sub-categories
- Expected: Error "cannot delete category with sub-categories"

---

## Quick Test Script

Buat file `test_category.sh`:

```bash
#!/bin/bash

# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"owner@restaurant.com","password":"owner123"}' \
  | jq -r '.data.token')

COMPANY_ID=$(curl -s -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN" \
  | jq -r '.data.company_id')

echo "Token: $TOKEN"
echo "Company ID: $COMPANY_ID"

# Get all main categories
echo -e "\n=== Main Categories ==="
curl -s -X GET "http://localhost:8080/api/v1/external/categories/company/$COMPANY_ID" \
  -H "Authorization: Bearer $TOKEN" | jq

# Create new category
echo -e "\n=== Create Category ==="
curl -s -X POST http://localhost:8080/api/v1/external/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"company_id\": \"$COMPANY_ID\",
    \"name\": \"Test Category\",
    \"description\": \"Testing\"
  }" | jq
```

Jalankan: `bash test_category.sh`
