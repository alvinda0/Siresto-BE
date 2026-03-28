# Quick Start - API Logging

Panduan cepat untuk menggunakan fitur API Logging.

## ✅ Sudah Siap Digunakan!

Sistem logging sudah otomatis aktif. Tidak perlu konfigurasi tambahan.

## 🚀 Cara Menggunakan

### 1. Jalankan Server
```bash
go run cmd/server/main.go
```

### 2. Login sebagai Admin
```bash
POST http://localhost:8080/api/v1/login
Content-Type: application/json

{
  "email": "admin@siresto.com",
  "password": "password123"
}
```

### 3. Lihat Semua Logs
```bash
GET http://localhost:8080/api/v1/logs
Authorization: Bearer <your-token>
```

### 4. Filter by Method
```bash
# Hanya POST requests
GET http://localhost:8080/api/v1/logs?method=POST
Authorization: Bearer <your-token>

# Hanya GET requests
GET http://localhost:8080/api/v1/logs?method=GET
Authorization: Bearer <your-token>

# Hanya DELETE requests
GET http://localhost:8080/api/v1/logs?method=DELETE
Authorization: Bearer <your-token>
```

### 5. Pagination
```bash
GET http://localhost:8080/api/v1/logs?page=1&limit=20
Authorization: Bearer <your-token>
```

### 6. Combine Filter & Pagination
```bash
GET http://localhost:8080/api/v1/logs?method=POST&page=1&limit=10
Authorization: Bearer <your-token>
```

### 7. Get Log Detail
```bash
GET http://localhost:8080/api/v1/logs/1
Authorization: Bearer <your-token>
```

## 📊 Apa yang Di-log?

Setiap request akan mencatat:
- ✅ HTTP Method (GET, POST, PUT, DELETE)
- ✅ Path/Endpoint
- ✅ Status Code
- ✅ Response Time (ms)
- ✅ IP Address
- ✅ User Agent
- ✅ Access From (postman, website, mobile, dll)
- ✅ User ID (jika login)
- ✅ Request Body
- ✅ Response Body
- ✅ Error Message (jika ada)

## 🔍 Access From Detection

Sistem otomatis mendeteksi sumber akses:
- **postman** - Dari Postman
- **website** - Dari browser
- **mobile** - Dari aplikasi mobile
- **curl** - Dari command line
- **insomnia** - Dari Insomnia
- **unknown** - Tidak terdeteksi

## 🔒 Access Control

### Internal Users (Lihat Semua Logs)
- ✅ SUPER_ADMIN - Bisa lihat semua logs dari semua company
- ✅ SUPPORT - Bisa lihat semua logs dari semua company
- ✅ FINANCE - Bisa lihat semua logs dari semua company

### External Users (Lihat Logs Company/Branch Sendiri)
- ✅ OWNER - Bisa lihat logs dari company sendiri
- ✅ ADMIN - Bisa lihat logs dari company dan branch sendiri
- ❌ CASHIER - Tidak bisa akses logs
- ❌ KITCHEN - Tidak bisa akses logs
- ❌ WAITER - Tidak bisa akses logs

## 📚 Dokumentasi Lengkap

- `API_LOGS_README.md` - Overview lengkap
- `API_LOGS_DOCUMENTATION.md` - API documentation
- `API_LOGS_TESTING.md` - Testing guide

## 💡 Tips

1. **Generate logs dulu** dengan hit beberapa endpoint
2. **Filter by method** untuk fokus ke operasi tertentu
3. **Check response_time** untuk identifikasi endpoint lambat
4. **Monitor error_message** untuk debugging
5. **Track user_id** untuk audit user activity

## 🎯 Use Cases

### Monitoring Traffic
```bash
GET /api/v1/logs?page=1&limit=100
```

### Debug Error
```bash
# Filter di aplikasi untuk status_code >= 400
GET /api/v1/logs
```

### Audit Changes
```bash
# Lihat semua POST/PUT/DELETE
GET /api/v1/logs?method=POST
GET /api/v1/logs?method=PUT
GET /api/v1/logs?method=DELETE
```

### Performance Analysis
```bash
# Ambil semua logs, analisis response_time
GET /api/v1/logs?limit=1000
```

## ⚡ Performance

- Logging dilakukan **async** (tidak memperlambat response)
- Request/response body **di-truncate** jika > 5000 chars
- Endpoint `/api/v1/logs` **tidak di-log** (avoid infinite loop)

## 🎉 Selesai!

Sistem logging sudah siap digunakan. Setiap request akan otomatis tercatat!
