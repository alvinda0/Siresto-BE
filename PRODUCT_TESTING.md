# Product API Testing Guide

## Prerequisites
1. Server running di `http://localhost:8080`
2. Database sudah di-migrate
3. Sudah ada user dengan role OWNER/ADMIN/CASHIER
4. Sudah ada category yang aktif

## Test Accounts
Gunakan akun dari `TEST_ACCOUNTS.md` atau buat akun baru.

---

## Test Case 1: Create Product

### Request
```bash
POST http://localhost:8080/api/v1/external/products
Authorization: Bearer <your_token>
Content-Type: application/json

{
  "branch_id": 1,
  "category_id": 1,
  "image": "https://images.unsplash.com/photo-1603133872878-684f208fb84b",
  "name": "Nasi Goreng Spesial",
  "description": "Nasi goreng dengan telur, ayam, dan sayuran segar",
  "stock": 50,
  "price": 25000,
  "position": "A1",
  "is_available": true
}
```

### Expected Response (201)
```json
{
  "status": "success",
  "message": "Product created successfully",
  "data": {
    "id": 1,
    "company_id": 1,
    "branch_id": 1,
    "category_id": 1,
    "image": "https://images.unsplash.com/photo-1603133872878-684f208fb84b",
    "name": "Nasi Goreng Spesial",
    "description": "Nasi goreng dengan telur, ayam, dan sayuran segar",
    "stock": 50,
    "price": 25000,
    "position": "A1",
    "is_available": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### cURL Command
```bash
curl -X POST http://localhost:8080/api/v1/external/products \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "branch_id": 1,
    "category_id": 1,
    "image": "https://images.unsplash.com/photo-1603133872878-684f208fb84b",
    "name": "Nasi Goreng Spesial",
    "description": "Nasi goreng dengan telur, ayam, dan sayuran segar",
    "stock": 50,
    "price": 25000,
    "position": "A1",
    "is_available": true
  }'
```

---

## Test Case 2: Create Multiple Products

### Product 2 - Minuman
```json
{
  "branch_id": 1,
  "category_id": 2,
  "image": "https://images.unsplash.com/photo-1546173159-315724a31696",
  "name": "Es Teh Manis",
  "description": "Teh manis dingin segar",
  "stock": 100,
  "price": 5000,
  "position": "B1",
  "is_available": true
}
```

### Product 3 - Makanan
```json
{
  "branch_id": 1,
  "category_id": 1,
  "image": "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38",
  "name": "Mie Goreng",
  "description": "Mie goreng pedas dengan telur",
  "stock": 30,
  "price": 20000,
  "position": "A2",
  "is_available": true
}
```

### Product 4 - Snack
```json
{
  "branch_id": 1,
  "category_id": 3,
  "image": "https://images.unsplash.com/photo-1599490659213-e2b9527bd087",
  "name": "Kentang Goreng",
  "description": "Kentang goreng crispy dengan saus",
  "stock": 40,
  "price": 15000,
  "position": "C1",
  "is_available": true
}
```

---

## Test Case 3: Get All Products

### Request
```bash
GET http://localhost:8080/api/v1/external/products
Authorization: Bearer <your_token>
```

### With Pagination
```bash
GET http://localhost:8080/api/v1/external/products?page=1&limit=10
Authorization: Bearer <your_token>
```

### With Search
```bash
GET http://localhost:8080/api/v1/external/products?search=goreng
Authorization: Bearer <your_token>
```

### Combined
```bash
GET http://localhost:8080/api/v1/external/products?search=nasi&page=1&limit=5
Authorization: Bearer <your_token>
```

### Expected Response (200)
```json
{
  "status": "success",
  "message": "Products retrieved successfully",
  "data": [
    {
      "id": 1,
      "company_id": 1,
      "branch_id": 1,
      "category_id": 1,
      "image": "https://images.unsplash.com/photo-1603133872878-684f208fb84b",
      "name": "Nasi Goreng Spesial",
      "description": "Nasi goreng dengan telur, ayam, dan sayuran segar",
      "stock": 50,
      "price": 25000,
      "position": "A1",
      "is_available": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "company": {
        "id": 1,
        "name": "Restoran ABC"
      },
      "branch": {
        "id": 1,
        "name": "Cabang Pusat"
      },
      "category": {
        "id": 1,
        "name": "Makanan Utama"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 4,
    "total_pages": 1
  }
}
```

### cURL Command
```bash
# Get all
curl -X GET http://localhost:8080/api/v1/external/products \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# With search
curl -X GET "http://localhost:8080/api/v1/external/products?search=goreng&page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## Test Case 4: Get Product by ID

### Request
```bash
GET http://localhost:8080/api/v1/external/products/1
Authorization: Bearer <your_token>
```

### Expected Response (200)
```json
{
  "status": "success",
  "message": "Product retrieved successfully",
  "data": {
    "id": 1,
    "company_id": 1,
    "branch_id": 1,
    "category_id": 1,
    "image": "https://images.unsplash.com/photo-1603133872878-684f208fb84b",
    "name": "Nasi Goreng Spesial",
    "description": "Nasi goreng dengan telur, ayam, dan sayuran segar",
    "stock": 50,
    "price": 25000,
    "position": "A1",
    "is_available": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "company": {
      "id": 1,
      "name": "Restoran ABC"
    },
    "branch": {
      "id": 1,
      "name": "Cabang Pusat"
    },
    "category": {
      "id": 1,
      "name": "Makanan Utama"
    }
  }
}
```

### cURL Command
```bash
curl -X GET http://localhost:8080/api/v1/external/products/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## Test Case 5: Update Product

### Request
```bash
PUT http://localhost:8080/api/v1/external/products/1
Authorization: Bearer <your_token>
Content-Type: application/json

{
  "category_id": 1,
  "image": "https://images.unsplash.com/photo-1603133872878-684f208fb84b",
  "name": "Nasi Goreng Spesial Premium",
  "description": "Nasi goreng dengan telur, ayam, udang, dan sayuran segar",
  "stock": 30,
  "price": 35000,
  "position": "A1",
  "is_available": true
}
```

### Expected Response (200)
```json
{
  "status": "success",
  "message": "Product updated successfully",
  "data": {
    "id": 1,
    "company_id": 1,
    "branch_id": 1,
    "category_id": 1,
    "image": "https://images.unsplash.com/photo-1603133872878-684f208fb84b",
    "name": "Nasi Goreng Spesial Premium",
    "description": "Nasi goreng dengan telur, ayam, udang, dan sayuran segar",
    "stock": 30,
    "price": 35000,
    "position": "A1",
    "is_available": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

### cURL Command
```bash
curl -X PUT http://localhost:8080/api/v1/external/products/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "image": "https://images.unsplash.com/photo-1603133872878-684f208fb84b",
    "name": "Nasi Goreng Spesial Premium",
    "description": "Nasi goreng dengan telur, ayam, udang, dan sayuran segar",
    "stock": 30,
    "price": 35000,
    "position": "A1",
    "is_available": true
  }'
```

---

## Test Case 6: Update Stock Only

### Request
```bash
PUT http://localhost:8080/api/v1/external/products/1
Authorization: Bearer <your_token>
Content-Type: application/json

{
  "category_id": 1,
  "image": "https://images.unsplash.com/photo-1603133872878-684f208fb84b",
  "name": "Nasi Goreng Spesial Premium",
  "description": "Nasi goreng dengan telur, ayam, udang, dan sayuran segar",
  "stock": 10,
  "price": 35000,
  "position": "A1",
  "is_available": true
}
```

---

## Test Case 7: Update Availability Status

### Request - Set to Unavailable
```bash
PUT http://localhost:8080/api/v1/external/products/1
Authorization: Bearer <your_token>
Content-Type: application/json

{
  "category_id": 1,
  "image": "https://images.unsplash.com/photo-1603133872878-684f208fb84b",
  "name": "Nasi Goreng Spesial Premium",
  "description": "Nasi goreng dengan telur, ayam, udang, dan sayuran segar",
  "stock": 0,
  "price": 35000,
  "position": "A1",
  "is_available": false
}
```

---

## Test Case 8: Delete Product

### Request
```bash
DELETE http://localhost:8080/api/v1/external/products/1
Authorization: Bearer <your_token>
```

### Expected Response (200)
```json
{
  "status": "success",
  "message": "Product deleted successfully",
  "data": null
}
```

### cURL Command
```bash
curl -X DELETE http://localhost:8080/api/v1/external/products/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## Error Test Cases

### Test Case 9: Create Product with Invalid Category
```json
{
  "branch_id": 1,
  "category_id": 999,
  "name": "Test Product",
  "price": 10000,
  "is_available": true
}
```

**Expected Response (400):**
```json
{
  "status": "error",
  "message": "Failed to create product",
  "error": "category not found or doesn't belong to your company"
}
```

---

### Test Case 10: Create Product with Wrong Branch
```json
{
  "branch_id": 999,
  "category_id": 1,
  "name": "Test Product",
  "price": 10000,
  "is_available": true
}
```

**Expected Response (403):**
```json
{
  "status": "error",
  "message": "You can only create products for your own branch",
  "error": null
}
```

---

### Test Case 11: Get Product from Different Branch
```bash
GET http://localhost:8080/api/v1/external/products/999
Authorization: Bearer <your_token>
```

**Expected Response (404):**
```json
{
  "status": "error",
  "message": "Product not found",
  "error": "record not found"
}
```

---

### Test Case 12: Create Product without Required Fields
```json
{
  "branch_id": 1,
  "category_id": 1
}
```

**Expected Response (400):**
```json
{
  "status": "error",
  "message": "Invalid request body",
  "error": "Key: 'CreateProductRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

---

## Postman Collection

### Setup Environment Variables
```
base_url: http://localhost:8080
token: <your_jwt_token>
branch_id: 1
category_id: 1
product_id: 1
```

### Collection Structure
```
Product API
├── Auth
│   └── Login
├── Products
│   ├── Create Product
│   ├── Get All Products
│   ├── Get All Products with Search
│   ├── Get Product by ID
│   ├── Update Product
│   └── Delete Product
└── Error Cases
    ├── Invalid Category
    ├── Wrong Branch
    └── Missing Required Fields
```

---

## Testing Checklist

- [ ] Create product dengan semua field
- [ ] Create product dengan field minimal (required only)
- [ ] Get all products tanpa parameter
- [ ] Get all products dengan pagination
- [ ] Get all products dengan search
- [ ] Get product by ID yang valid
- [ ] Get product by ID yang tidak ada
- [ ] Update product semua field
- [ ] Update product sebagian field
- [ ] Update stock produk
- [ ] Update status ketersediaan
- [ ] Delete product yang valid
- [ ] Delete product yang tidak ada
- [ ] Create product dengan category invalid
- [ ] Create product dengan branch yang salah
- [ ] Akses product dari company lain (harus gagal)
- [ ] Akses product dari branch lain (harus gagal)

---

## Notes

1. Ganti `YOUR_TOKEN_HERE` dengan token JWT yang valid
2. Pastikan category_id dan branch_id sudah ada di database
3. Product ID akan auto-increment, sesuaikan dengan data Anda
4. Semua endpoint otomatis filter berdasarkan company_id dan branch_id dari token
5. Soft delete digunakan, data tidak benar-benar terhapus
