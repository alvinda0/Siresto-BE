# Product Image Upload Guide

## Overview
Endpoint POST dan PUT product sekarang support 2 cara untuk mengirim data:
1. **JSON** - Kirim URL gambar sebagai string
2. **Multipart Form-Data** - Upload file gambar langsung

---

## 1. CREATE PRODUCT dengan JSON (URL Gambar)

**POST** `/api/v1/external/products`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "branch_id": "b4506ffd-bc07-4fd7-8e30-f3de4b112f80",
  "category_id": "efcf879f-6da5-4854-a9a7-c53a785aa30d",
  "image": "https://example.com/image.jpg",
  "name": "Nasi Goreng",
  "description": "Nasi goreng enak",
  "stock": 50,
  "price": 25000,
  "position": "A1",
  "is_available": true
}
```

---

## 2. CREATE PRODUCT dengan Upload File

**POST** `/api/v1/external/products`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Body (Form-Data):**
```
branch_id: b4506ffd-bc07-4fd7-8e30-f3de4b112f80
category_id: efcf879f-6da5-4854-a9a7-c53a785aa30d
image: [FILE] (pilih file gambar)
name: Nasi Goreng
description: Nasi goreng enak
stock: 50
price: 25000
position: A1
is_available: true
```

### cURL Example:
```bash
curl -X POST http://localhost:8080/api/v1/external/products \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "branch_id=b4506ffd-bc07-4fd7-8e30-f3de4b112f80" \
  -F "category_id=efcf879f-6da5-4854-a9a7-c53a785aa30d" \
  -F "image=@/path/to/image.jpg" \
  -F "name=Nasi Goreng" \
  -F "description=Nasi goreng enak" \
  -F "stock=50" \
  -F "price=25000" \
  -F "position=A1" \
  -F "is_available=true"
```

### Postman:
1. Method: POST
2. URL: `http://localhost:8080/api/v1/external/products`
3. Headers: 
   - Authorization: Bearer YOUR_TOKEN
4. Body: 
   - Select "form-data"
   - Add fields:
     - branch_id: text
     - category_id: text
     - image: File (pilih file)
     - name: text
     - description: text
     - stock: text
     - price: text
     - position: text
     - is_available: text

---

## 3. UPDATE PRODUCT dengan JSON (URL Gambar)

**PUT** `/api/v1/external/products/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
  "category_id": "efcf879f-6da5-4854-a9a7-c53a785aa30d",
  "image": "https://example.com/new-image.jpg",
  "name": "Nasi Goreng Premium",
  "description": "Nasi goreng premium enak",
  "stock": 30,
  "price": 35000,
  "position": "A2",
  "is_available": true
}
```

---

## 4. UPDATE PRODUCT dengan Upload File

**PUT** `/api/v1/external/products/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Body (Form-Data):**
```
category_id: efcf879f-6da5-4854-a9a7-c53a785aa30d
image: [FILE] (pilih file gambar baru)
name: Nasi Goreng Premium
description: Nasi goreng premium enak
stock: 30
price: 35000
position: A2
is_available: true
```

### cURL Example:
```bash
curl -X PUT http://localhost:8080/api/v1/external/products/PRODUCT_ID \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "category_id=efcf879f-6da5-4854-a9a7-c53a785aa30d" \
  -F "image=@/path/to/new-image.jpg" \
  -F "name=Nasi Goreng Premium" \
  -F "description=Nasi goreng premium enak" \
  -F "stock=30" \
  -F "price=35000" \
  -F "position=A2" \
  -F "is_available=true"
```

---

## File Upload Specifications

### Allowed File Types:
- image/jpeg
- image/jpg
- image/png
- image/gif
- image/webp

### Maximum File Size:
- 5 MB (5,242,880 bytes)

### Upload Directory:
- Files disimpan di: `uploads/products/`
- File URL: `http://localhost:8080/uploads/products/filename.jpg`

### File Naming:
- Format: `{uuid}_{timestamp}.{extension}`
- Contoh: `a1b2c3d4-e5f6-7890-abcd-ef1234567890_1234567890.jpg`

---

## Response Examples

### Success Response (201/200):
```json
{
  "success": true,
  "message": "Product created successfully",
  "status": 201,
  "timestamp": "2026-03-28T12:00:00Z",
  "data": {
    "id": "36051b5a-7589-42e2-a5bb-396d3273f9fd",
    "company_id": "2fa830c2-daec-4ddb-b061-2f7c50b7562b",
    "branch_id": "8fe6b4f9-3dc9-4773-82a8-9f92f0d86458",
    "category_id": "e9b7273e-8df7-4462-bbb8-eff8e2aacf52",
    "image": "http://localhost:8080/uploads/products/a1b2c3d4_1234567890.jpg",
    "name": "Nasi Goreng",
    "description": "Nasi goreng enak",
    "stock": 50,
    "price": 25000,
    "position": "A1",
    "is_available": true,
    "created_at": "2026-03-28T12:00:00+07:00",
    "updated_at": "2026-03-28T12:00:00+07:00",
    "company": {
      "id": "2fa830c2-daec-4ddb-b061-2f7c50b7562b",
      "name": "PT Maju Jaya"
    },
    "branch": {
      "id": "8fe6b4f9-3dc9-4773-82a8-9f92f0d86458",
      "name": "Patron"
    },
    "category": {
      "id": "e9b7273e-8df7-4462-bbb8-eff8e2aacf52",
      "name": "Makanan"
    }
  }
}
```

### Error Response - File Too Large:
```json
{
  "success": false,
  "message": "Failed to upload image",
  "status": 400,
  "timestamp": "2026-03-28T12:00:00Z",
  "error": "file size exceeds maximum allowed size of 5242880 bytes"
}
```

### Error Response - Invalid File Type:
```json
{
  "success": false,
  "message": "Failed to upload image",
  "status": 400,
  "timestamp": "2026-03-28T12:00:00Z",
  "error": "file type application/pdf is not allowed. Allowed types: [image/jpeg image/jpg image/png image/gif image/webp]"
}
```

---

## Environment Variables

Tambahkan di `.env`:
```env
BASE_URL=http://localhost:8080
```

Untuk production:
```env
BASE_URL=https://api.yourdomain.com
```

---

## Notes

1. **Flexible Input**: Endpoint bisa menerima JSON atau Form-Data, sistem otomatis detect dari Content-Type header
2. **Optional Image**: Field `image` optional, bisa tidak dikirim
3. **URL or File**: Bisa kirim URL string atau upload file, tidak bisa keduanya sekaligus
4. **File Storage**: File disimpan di folder `uploads/products/` di server
5. **Static Serving**: File bisa diakses via `http://localhost:8080/uploads/products/filename.jpg`
6. **Update Image**: Saat update, gambar lama tidak otomatis terhapus (bisa ditambahkan fitur ini nanti)

---

## Testing Checklist

- [ ] Create product dengan JSON (URL)
- [ ] Create product dengan upload file
- [ ] Create product tanpa image
- [ ] Update product dengan JSON (URL)
- [ ] Update product dengan upload file
- [ ] Upload file > 5MB (harus error)
- [ ] Upload file bukan gambar (harus error)
- [ ] Upload file format tidak didukung (harus error)
- [ ] Akses gambar via URL setelah upload
