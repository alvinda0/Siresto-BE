# Category API Documentation

## Overview
API untuk mengelola kategori dan sub-kategori produk. Kategori bisa dibuat di level:
1. **Company Level** (branch_id = null): Berlaku untuk semua cabang
2. **Branch Level** (branch_id = uuid): Khusus untuk cabang tertentu

Setiap kategori memiliki posisi yang dimulai dari 1 (bukan 0).

## Endpoints

### 1. Create Category
Membuat kategori baru (main category atau sub-category)

**Endpoint:** `POST /api/v1/external/categories`

**Request Body:**
```json
{
  "company_id": "uuid",
  "branch_id": "uuid (optional, jika null = company level)",
  "parent_id": "uuid (optional, untuk sub-category)",
  "name": "Makanan",
  "description": "Kategori makanan",
  "position": 1
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Category created successfully",
  "data": {
    "id": "uuid",
    "company_id": "uuid",
    "branch_id": null,
    "parent_id": null,
    "name": "Makanan",
    "description": "Kategori makanan",
    "position": 1,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Notes:**
- Jika `position` tidak diisi atau 0, sistem akan otomatis set ke posisi terakhir + 1
- Posisi dimulai dari 1, bukan 0
- `parent_id` diisi jika ingin membuat sub-category
- `branch_id` null = kategori berlaku untuk semua cabang (company level)
- `branch_id` diisi = kategori khusus untuk cabang tersebut (branch level)

---

### 2. Update Category
Update kategori yang sudah ada

**Endpoint:** `PUT /api/v1/external/categories/:id`

**Request Body:**
```json
{
  "parent_id": "uuid (optional)",
  "name": "Makanan Updated",
  "description": "Deskripsi baru",
  "position": 2,
  "is_active": true
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Category updated successfully",
  "data": null
}
```

---

### 3. Delete Category
Hapus kategori

**Endpoint:** `DELETE /api/v1/external/categories/:id`

**Response:**
```json
{
  "status": "success",
  "message": "Category deleted successfully",
  "data": null
}
```

**Notes:**
- Tidak bisa hapus kategori yang masih punya sub-category
- Hapus sub-category terlebih dahulu

---

### 4. Get Category by ID
Ambil detail kategori beserta sub-categories

**Endpoint:** `GET /api/v1/external/categories/:id`

**Response:**
```json
{
  "status": "success",
  "message": "Category retrieved successfully",
  "data": {
    "id": "uuid",
    "company_id": "uuid",
    "parent_id": null,
    "name": "Makanan",
    "description": "Kategori makanan",
    "position": 1,
    "is_active": true,
    "sub_categories": [
      {
        "id": "uuid",
        "company_id": "uuid",
        "parent_id": "uuid",
        "name": "Makanan Utama",
        "description": "Nasi, mie, pasta",
        "position": 1,
        "is_active": true
      }
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### 5. Get Categories by Company
Ambil semua kategori berdasarkan company dan/atau branch

**Endpoint:** `GET /api/v1/external/categories/company/:company_id`

**Query Parameters:**
- `branch_id` (optional): Filter by branch
  - Jika tidak diisi: ambil kategori company level (branch_id = null)
  - Jika diisi: ambil kategori khusus branch tersebut
- `parent_id` (optional): Filter by parent category
  - Jika tidak diisi: ambil semua main categories (parent_id = null)
  - Jika diisi: ambil sub-categories dari parent tersebut

**Examples:**

Get company level main categories:
```
GET /api/v1/external/categories/company/uuid-company
```

Get branch specific main categories:
```
GET /api/v1/external/categories/company/uuid-company?branch_id=uuid-branch
```

Get sub-categories of a parent (company level):
```
GET /api/v1/external/categories/company/uuid-company?parent_id=uuid-parent
```

Get sub-categories of a parent (branch level):
```
GET /api/v1/external/categories/company/uuid-company?branch_id=uuid-branch&parent_id=uuid-parent
```

**Response:**
```json
{
  "status": "success",
  "message": "Categories retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "company_id": "uuid",
      "parent_id": null,
      "name": "Makanan",
      "description": "Kategori makanan",
      "position": 1,
      "is_active": true,
      "sub_categories": [
        {
          "id": "uuid",
          "name": "Makanan Utama",
          "position": 1
        }
      ]
    },
    {
      "id": "uuid",
      "company_id": "uuid",
      "parent_id": null,
      "name": "Minuman",
      "description": "Kategori minuman",
      "position": 2,
      "is_active": true,
      "sub_categories": []
    }
  ]
}
```

---

### 6. Reorder Categories
Ubah urutan posisi kategori

**Endpoint:** `POST /api/v1/external/categories/company/:company_id/reorder`

**Query Parameters:**
- `branch_id` (optional): Untuk reorder kategori branch level
- `parent_id` (optional): Untuk reorder sub-categories

**Request Body:**
```json
{
  "category_ids": [
    "uuid-category-2",
    "uuid-category-1",
    "uuid-category-3"
  ]
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Categories reordered successfully",
  "data": null
}
```

**Notes:**
- Array `category_ids` menentukan urutan baru
- Index 0 akan jadi position 1, index 1 jadi position 2, dst
- Posisi dimulai dari 1

---

## Category Levels

### Company Level (branch_id = null)
Kategori yang berlaku untuk semua cabang dalam company tersebut.

**Use Case:**
- Menu standar yang sama di semua cabang
- Kategori umum seperti "Makanan", "Minuman"

**Example:**
```json
{
  "company_id": "uuid-company",
  "branch_id": null,
  "name": "Makanan"
}
```

### Branch Level (branch_id = uuid)
Kategori khusus untuk cabang tertentu.

**Use Case:**
- Menu spesial cabang tertentu
- Kategori regional
- Menu seasonal per cabang

**Example:**
```json
{
  "company_id": "uuid-company",
  "branch_id": "uuid-branch-jakarta",
  "name": "Menu Spesial Jakarta"
}
```

---

## Example Usage Flow

### 1. Membuat Company Level Categories
```bash
# Buat kategori Makanan (berlaku untuk semua cabang)
curl -X POST http://localhost:8080/api/v1/external/categories \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "uuid-company",
    "branch_id": null,
    "name": "Makanan",
    "description": "Kategori makanan"
  }'
```

### 2. Membuat Branch Specific Categories
```bash
# Buat kategori khusus Cabang Jakarta
curl -X POST http://localhost:8080/api/v1/external/categories \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "uuid-company",
    "branch_id": "uuid-branch-jakarta",
    "name": "Menu Spesial Jakarta",
    "description": "Menu khusus cabang Jakarta"
  }'
```

### 3. Get Company Level Categories
```bash
curl -X GET http://localhost:8080/api/v1/external/categories/company/uuid-company \
  -H "Authorization: Bearer <token>"
```

### 4. Get Branch Specific Categories
```bash
curl -X GET "http://localhost:8080/api/v1/external/categories/company/uuid-company?branch_id=uuid-branch" \
  -H "Authorization: Bearer <token>"
```

### 5. Reorder Categories
```bash
curl -X POST http://localhost:8080/api/v1/external/categories/company/uuid-company/reorder \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "category_ids": [
      "uuid-minuman",
      "uuid-makanan"
    ]
  }'
```

---

## Seeded Data

Setelah menjalankan seeder, akan ada data kategori berikut:

**Company Level Categories (berlaku untuk semua cabang):**
1. Makanan (position: 1)
   - Makanan Utama (position: 1)
   - Makanan Pembuka (position: 2)
   - Makanan Penutup (position: 3)
2. Minuman (position: 2)
   - Minuman Panas (position: 1)
   - Minuman Dingin (position: 2)

**Branch Level Categories (khusus Cabang Jakarta Pusat):**
1. Menu Spesial Jakarta (position: 1)

---

## Important Notes

1. **Posisi dimulai dari 1**, bukan 0
2. **Company Level** (branch_id = null): Kategori berlaku untuk semua cabang
3. **Branch Level** (branch_id = uuid): Kategori khusus untuk cabang tertentu
4. Jika tidak set position saat create, sistem otomatis set ke posisi terakhir + 1
5. Tidak bisa hapus kategori yang masih punya sub-category
6. Parent category dan sub-category harus dalam company yang sama
7. Parent category dan sub-category harus dalam branch yang sama (jika branch level)
8. Category tidak bisa jadi parent dari dirinya sendiri

---

## Use Cases

### Scenario 1: Menu Standar Semua Cabang
Buat kategori company level (branch_id = null) untuk menu yang sama di semua cabang.

```json
{
  "company_id": "uuid",
  "branch_id": null,
  "name": "Makanan"
}
```

### Scenario 2: Menu Khusus Per Cabang
Buat kategori branch level untuk menu spesial cabang tertentu.

```json
{
  "company_id": "uuid",
  "branch_id": "uuid-branch-jakarta",
  "name": "Menu Spesial Jakarta"
}
```

### Scenario 3: Kombinasi
- Company level: Menu standar (Makanan, Minuman)
- Branch level: Menu spesial per cabang (Menu Spesial Jakarta, Menu Spesial Bandung)
