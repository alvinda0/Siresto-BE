# API Logs Testing Guide

Panduan testing untuk API Logs endpoint.

## Prerequisites

1. Server sudah running di `http://localhost:8080`
2. Database sudah di-migrate (tabel `api_logs` sudah dibuat)
3. Punya akun dengan role INTERNAL (SUPER_ADMIN, SUPPORT, atau FINANCE)

## Test Accounts

Gunakan salah satu akun internal dari `TEST_ACCOUNTS.md`:

```json
{
  "email": "admin@siresto.com",
  "password": "password123"
}
```

---

## Testing Steps

### Step 1: Login untuk Mendapatkan Token

**Request:**
```bash
POST http://localhost:8080/api/v1/login
Content-Type: application/json

{
  "email": "admin@siresto.com",
  "password": "password123"
}
```

**Expected Response:**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid-here",
      "name": "Super Admin",
      "email": "admin@siresto.com",
      "role": {
        "name": "SUPER_ADMIN",
        "display_name": "Super Admin"
      }
    }
  }
}
```

**Copy token untuk digunakan di request selanjutnya.**

---

### Step 2: Generate Some Logs

Sebelum test endpoint logs, buat beberapa request ke endpoint lain untuk generate logs:

**Request 1 - GET:**
```bash
GET http://localhost:8080/api/v1/roles
Authorization: Bearer <your-token>
```

**Request 2 - POST (akan error karena tidak ada body, tapi akan generate log):**
```bash
POST http://localhost:8080/api/v1/external/categories
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "name": "Test Category"
}
```

**Request 3 - GET dengan Postman:**
```bash
GET http://localhost:8080/api/v1/auth/me
Authorization: Bearer <your-token>
User-Agent: PostmanRuntime/7.32.3
```

---

### Step 3: Test Get All Logs (Basic)

**Request:**
```bash
GET http://localhost:8080/api/v1/logs
Authorization: Bearer <your-token>
```

**Expected Response:**
```json
{
  "status": "success",
  "message": "Logs retrieved successfully",
  "data": [
    {
      "id": 1,
      "method": "GET",
      "path": "/api/v1/roles",
      "status_code": 200,
      "response_time": 45,
      "ip_address": "::1",
      "user_agent": "PostmanRuntime/7.32.3",
      "access_from": "postman",
      "user_id": 1,
      "request_body": "",
      "response_body": "{...}",
      "error_message": "",
      "created_at": "2026-03-28T10:30:00Z",
      "updated_at": "2026-03-28T10:30:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 5,
    "total_pages": 1
  }
}
```

**Validation:**
- ✅ Status code: 200
- ✅ Response memiliki field `data` (array)
- ✅ Response memiliki field `meta` dengan pagination info
- ✅ Setiap log memiliki semua field yang diperlukan

---

### Step 4: Test Get All Logs with Pagination

**Request:**
```bash
GET http://localhost:8080/api/v1/logs?page=1&limit=5
Authorization: Bearer <your-token>
```

**Expected Response:**
```json
{
  "status": "success",
  "message": "Logs retrieved successfully",
  "data": [...],
  "meta": {
    "current_page": 1,
    "per_page": 5,
    "total": 20,
    "total_pages": 4
  }
}
```

**Validation:**
- ✅ `meta.per_page` = 5
- ✅ `data` array length ≤ 5
- ✅ `meta.total_pages` dihitung dengan benar

---

### Step 5: Test Filter by Method (GET)

**Request:**
```bash
GET http://localhost:8080/api/v1/logs?method=GET
Authorization: Bearer <your-token>
```

**Expected Response:**
```json
{
  "status": "success",
  "message": "Logs retrieved successfully",
  "data": [
    {
      "id": 1,
      "method": "GET",
      ...
    },
    {
      "id": 3,
      "method": "GET",
      ...
    }
  ],
  "meta": {...}
}
```

**Validation:**
- ✅ Semua log dalam `data` memiliki `method` = "GET"
- ✅ Tidak ada log dengan method lain

---

### Step 6: Test Filter by Method (POST)

**Request:**
```bash
GET http://localhost:8080/api/v1/logs?method=POST
Authorization: Bearer <your-token>
```

**Expected Response:**
```json
{
  "status": "success",
  "message": "Logs retrieved successfully",
  "data": [
    {
      "id": 2,
      "method": "POST",
      ...
    }
  ],
  "meta": {...}
}
```

**Validation:**
- ✅ Semua log dalam `data` memiliki `method` = "POST"

---

### Step 7: Test Filter by Method (PUT)

**Request:**
```bash
GET http://localhost:8080/api/v1/logs?method=PUT
Authorization: Bearer <your-token>
```

**Expected Response:**
```json
{
  "status": "success",
  "message": "Logs retrieved successfully",
  "data": [],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 0,
    "total_pages": 0
  }
}
```

**Validation:**
- ✅ `data` array kosong jika tidak ada log dengan method PUT
- ✅ `meta.total` = 0

---

### Step 8: Test Filter by Method (DELETE)

**Request:**
```bash
GET http://localhost:8080/api/v1/logs?method=DELETE
Authorization: Bearer <your-token>
```

**Expected Response:**
```json
{
  "status": "success",
  "message": "Logs retrieved successfully",
  "data": [],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 0,
    "total_pages": 0
  }
}
```

**Validation:**
- ✅ `data` array kosong jika tidak ada log dengan method DELETE

---

### Step 9: Test Combine Filter and Pagination

**Request:**
```bash
GET http://localhost:8080/api/v1/logs?method=GET&page=1&limit=3
Authorization: Bearer <your-token>
```

**Expected Response:**
```json
{
  "status": "success",
  "message": "Logs retrieved successfully",
  "data": [
    {
      "id": 5,
      "method": "GET",
      ...
    },
    {
      "id": 4,
      "method": "GET",
      ...
    },
    {
      "id": 3,
      "method": "GET",
      ...
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 3,
    "total": 10,
    "total_pages": 4
  }
}
```

**Validation:**
- ✅ Hanya menampilkan log dengan method GET
- ✅ Maksimal 3 items per page
- ✅ Pagination bekerja dengan benar

---

### Step 10: Test Get Log by ID

**Request:**
```bash
GET http://localhost:8080/api/v1/logs/1
Authorization: Bearer <your-token>
```

**Expected Response:**
```json
{
  "status": "success",
  "message": "Log retrieved successfully",
  "data": {
    "id": 1,
    "method": "GET",
    "path": "/api/v1/roles",
    "status_code": 200,
    "response_time": 45,
    "ip_address": "::1",
    "user_agent": "PostmanRuntime/7.32.3",
    "access_from": "postman",
    "user_id": 1,
    "request_body": "",
    "response_body": "{\"status\":\"success\",...}",
    "error_message": "",
    "created_at": "2026-03-28T10:30:00Z",
    "updated_at": "2026-03-28T10:30:00Z"
  }
}
```

**Validation:**
- ✅ Status code: 200
- ✅ Response memiliki field `data` (object, bukan array)
- ✅ Log memiliki semua field lengkap

---

### Step 11: Test Get Log by Invalid ID

**Request:**
```bash
GET http://localhost:8080/api/v1/logs/99999
Authorization: Bearer <your-token>
```

**Expected Response:**
```json
{
  "status": "error",
  "message": "Log not found",
  "error": "log not found"
}
```

**Validation:**
- ✅ Status code: 404
- ✅ Error message sesuai

---

### Step 12: Test Access from Different Sources

**Test 1 - Postman:**
```bash
GET http://localhost:8080/api/v1/roles
Authorization: Bearer <your-token>
User-Agent: PostmanRuntime/7.32.3
```

**Test 2 - Website (Browser):**
```bash
GET http://localhost:8080/api/v1/roles
Authorization: Bearer <your-token>
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0
```

**Test 3 - Mobile:**
```bash
GET http://localhost:8080/api/v1/roles
Authorization: Bearer <your-token>
User-Agent: MyApp/1.0 (Android 12; Mobile)
```

**Test 4 - cURL:**
```bash
curl -X GET "http://localhost:8080/api/v1/roles" \
  -H "Authorization: Bearer <your-token>"
```

**Kemudian check logs:**
```bash
GET http://localhost:8080/api/v1/logs
Authorization: Bearer <your-token>
```

**Validation:**
- ✅ Log dengan User-Agent Postman memiliki `access_from` = "postman"
- ✅ Log dengan User-Agent browser memiliki `access_from` = "website"
- ✅ Log dengan User-Agent mobile memiliki `access_from` = "mobile"
- ✅ Log dengan User-Agent curl memiliki `access_from` = "curl"

---

### Step 13: Test Unauthorized Access

**Request (tanpa token):**
```bash
GET http://localhost:8080/api/v1/logs
```

**Expected Response:**
```json
{
  "status": "error",
  "message": "Unauthorized",
  "error": "Authorization header required"
}
```

**Validation:**
- ✅ Status code: 401
- ✅ Error message sesuai

---

### Step 14: Test Access with External Role

Login dengan akun OWNER atau ADMIN (external role):

**Request:**
```bash
POST http://localhost:8080/api/v1/login
Content-Type: application/json

{
  "email": "owner@restaurant.com",
  "password": "password123"
}
```

**Kemudian coba akses logs:**
```bash
GET http://localhost:8080/api/v1/logs
Authorization: Bearer <external-user-token>
```

**Expected Response:**
```json
{
  "status": "error",
  "message": "Access denied",
  "error": "Only internal users can access this resource"
}
```

**Validation:**
- ✅ Status code: 403
- ✅ External users tidak bisa akses logs

---

## Checklist Testing

### Basic Functionality
- [ ] Get all logs berhasil
- [ ] Get log by ID berhasil
- [ ] Pagination bekerja dengan benar
- [ ] Filter by method GET bekerja
- [ ] Filter by method POST bekerja
- [ ] Filter by method PUT bekerja
- [ ] Filter by method DELETE bekerja
- [ ] Combine filter dan pagination bekerja

### Access Control
- [ ] Internal users (SUPER_ADMIN) bisa akses
- [ ] Internal users (SUPPORT) bisa akses
- [ ] Internal users (FINANCE) bisa akses
- [ ] External users (OWNER) tidak bisa akses
- [ ] External users (ADMIN) tidak bisa akses
- [ ] Unauthenticated users tidak bisa akses

### Access From Detection
- [ ] Postman terdeteksi sebagai "postman"
- [ ] Browser terdeteksi sebagai "website"
- [ ] Mobile app terdeteksi sebagai "mobile"
- [ ] cURL terdeteksi sebagai "curl"

### Error Handling
- [ ] Invalid log ID return 404
- [ ] Missing token return 401
- [ ] Wrong role return 403

### Performance
- [ ] Response time < 200ms untuk get all logs
- [ ] Response time < 100ms untuk get by ID
- [ ] Logging tidak memperlambat endpoint lain

---

## Notes

1. **Automatic Logging**: Setiap request ke endpoint lain akan otomatis tercatat di logs
2. **Async Processing**: Logging dilakukan secara async, tidak mempengaruhi response time
3. **Self-Exclusion**: Endpoint `/api/v1/logs` tidak di-log untuk menghindari infinite loop
4. **Body Truncation**: Request/response body > 5000 chars akan di-truncate

---

## Troubleshooting

### Logs tidak muncul
- Pastikan middleware sudah terpasang di routes
- Check database connection
- Pastikan tabel `api_logs` sudah dibuat

### Access from selalu "unknown"
- Check User-Agent header di request
- Pastikan User-Agent tidak kosong

### Pagination tidak bekerja
- Check query parameters (page, limit)
- Pastikan nilai page dan limit > 0
