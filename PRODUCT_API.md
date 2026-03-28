# Product API Documentation

## Base URL
```
/api/v1/external
```

## Authentication
Semua endpoint membutuhkan:
- Header: `Authorization: Bearer <token>`
- Token didapat dari endpoint `/api/v1/login`

## Endpoints

### 1. Create Product
**POST** `/products`

Membuat produk baru. Data akan otomatis tersimpan sesuai company dan branch user yang login.

**Request Body:**
```json
{
  "branch_id": 1,
  "category_id": 1,
  "image": "https://example.com/image.jpg",
  "name": "Nasi Goreng Spesial",
  "description": "Nasi goreng dengan telur, ayam, dan sayuran",
  "stock": 50,
  "price": 25000,
  "position": "A1",
  "is_available": true
}
```

**Field Descriptions:**
- `branch_id` (required): ID cabang (harus sesuai dengan branch user yang login)
- `category_id` (required): ID kategori produk
- `image` (optional): URL gambar produk
- `name` (required): Nama produk
- `description` (optional): Deskripsi produk
- `stock` (optional): Jumlah stok, default: 0
- `price` (required): Harga jual produk
- `position` (optional): Posisi produk (misal: rak, nomor meja, dll)
- `is_available` (optional): Status ketersediaan, default: true

**Response Success (201):**
```json
{
  "status": "success",
  "message": "Product created successfully",
  "data": {
    "id": 1,
    "company_id": 1,
    "branch_id": 1,
    "category_id": 1,
    "image": "https://example.com/image.jpg",
    "name": "Nasi Goreng Spesial",
    "description": "Nasi goreng dengan telur, ayam, dan sayuran",
    "stock": 50,
    "price": 25000,
    "position": "A1",
    "is_available": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### 2. Get All Products
**GET** `/products`

Mengambil semua produk sesuai company dan branch user yang login.

**Query Parameters:**
- `search` (optional): Pencarian berdasarkan nama atau deskripsi produk
- `page` (optional): Nomor halaman, default: 1
- `limit` (optional): Jumlah data per halaman, default: 10

**Example Request:**
```
GET /api/v1/external/products?search=nasi&page=1&limit=10
```

**Response Success (200):**
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
      "image": "https://example.com/image.jpg",
      "name": "Nasi Goreng Spesial",
      "description": "Nasi goreng dengan telur, ayam, dan sayuran",
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
    "total": 1,
    "total_pages": 1
  }
}
```

---

### 3. Get Product by ID
**GET** `/products/:id`

Mengambil detail produk berdasarkan ID. Hanya bisa mengakses produk dari company dan branch sendiri.

**Example Request:**
```
GET /api/v1/external/products/1
```

**Response Success (200):**
```json
{
  "status": "success",
  "message": "Product retrieved successfully",
  "data": {
    "id": 1,
    "company_id": 1,
    "branch_id": 1,
    "category_id": 1,
    "image": "https://example.com/image.jpg",
    "name": "Nasi Goreng Spesial",
    "description": "Nasi goreng dengan telur, ayam, dan sayuran",
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

**Response Error (404):**
```json
{
  "status": "error",
  "message": "Product not found",
  "error": "record not found"
}
```

---

### 4. Update Product
**PUT** `/products/:id`

Mengupdate produk. Hanya bisa mengupdate produk dari company dan branch sendiri.

**Request Body:**
```json
{
  "category_id": 1,
  "image": "https://example.com/new-image.jpg",
  "name": "Nasi Goreng Spesial Premium",
  "description": "Nasi goreng dengan telur, ayam, udang, dan sayuran",
  "stock": 30,
  "price": 35000,
  "position": "A2",
  "is_available": true
}
```

**Response Success (200):**
```json
{
  "status": "success",
  "message": "Product updated successfully",
  "data": {
    "id": 1,
    "company_id": 1,
    "branch_id": 1,
    "category_id": 1,
    "image": "https://example.com/new-image.jpg",
    "name": "Nasi Goreng Spesial Premium",
    "description": "Nasi goreng dengan telur, ayam, udang, dan sayuran",
    "stock": 30,
    "price": 35000,
    "position": "A2",
    "is_available": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

---

### 5. Delete Product
**DELETE** `/products/:id`

Menghapus produk (soft delete). Hanya bisa menghapus produk dari company dan branch sendiri.

**Example Request:**
```
DELETE /api/v1/external/products/1
```

**Response Success (200):**
```json
{
  "status": "success",
  "message": "Product deleted successfully",
  "data": null
}
```

**Response Error (404):**
```json
{
  "status": "error",
  "message": "Failed to delete product",
  "error": "product not found"
}
```

---

## Error Responses

### 400 Bad Request
```json
{
  "status": "error",
  "message": "Invalid request body",
  "error": "Key: 'CreateProductRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

### 401 Unauthorized
```json
{
  "status": "error",
  "message": "Company ID not found",
  "error": null
}
```

### 403 Forbidden
```json
{
  "status": "error",
  "message": "You can only create products for your own branch",
  "error": null
}
```

### 404 Not Found
```json
{
  "status": "error",
  "message": "Product not found",
  "error": "record not found"
}
```

---

## Notes

1. **Company & Branch Filtering**: Semua data produk otomatis difilter berdasarkan `company_id` dan `branch_id` dari user yang login. Tidak perlu mengirim parameter company atau branch di query.

2. **Category Validation**: Kategori yang dipilih harus milik company yang sama dengan user yang login.

3. **Branch Validation**: Saat create product, `branch_id` yang dikirim harus sama dengan `branch_id` user yang login.

4. **Search**: Parameter search akan mencari di field `name` dan `description`.

5. **Soft Delete**: Delete menggunakan soft delete, data tidak benar-benar dihapus dari database.

6. **Image**: Field image menerima URL string. Untuk upload gambar, perlu endpoint terpisah untuk upload file.

7. **Stock Management**: Field stock bisa diupdate manual atau nanti bisa diintegrasikan dengan sistem transaksi.

8. **Price**: Harga dalam format decimal(15,2), mendukung hingga 13 digit sebelum koma dan 2 digit desimal.

---

## Testing Flow

1. Login sebagai user external (OWNER/ADMIN/CASHIER)
2. Pastikan sudah ada category (buat dulu jika belum ada)
3. Create product dengan category_id yang valid
4. Get all products untuk melihat list
5. Get product by ID untuk detail
6. Update product untuk mengubah data
7. Delete product untuk menghapus
