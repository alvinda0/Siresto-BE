# API Logs Documentation

API untuk monitoring dan tracking semua request yang masuk ke sistem.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Semua endpoint memerlukan authentication token (Bearer Token) dan hanya bisa diakses oleh user dengan role INTERNAL (SUPER_ADMIN, SUPPORT, FINANCE).

## Endpoints

### 1. Get All API Logs

Mendapatkan semua log API dengan pagination dan filter.

**Endpoint:** `GET /logs`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| page | integer | No | 1 | Nomor halaman |
| limit | integer | No | 10 | Jumlah data per halaman |
| method | string | No | - | Filter berdasarkan HTTP method (GET, POST, PUT, DELETE) |

**Example Request:**
```bash
# Get all logs (page 1, 10 items)
curl -X GET "http://localhost:8080/api/v1/logs" \
  -H "Authorization: Bearer <token>"

# Get logs with pagination
curl -X GET "http://localhost:8080/api/v1/logs?page=2&limit=20" \
  -H "Authorization: Bearer <token>"

# Filter by method
curl -X GET "http://localhost:8080/api/v1/logs?method=POST" \
  -H "Authorization: Bearer <token>"

# Combine filters
curl -X GET "http://localhost:8080/api/v1/logs?page=1&limit=50&method=DELETE" \
  -H "Authorization: Bearer <token>"
```

**Success Response (200 OK):**
```json
{
  "status": "success",
  "message": "Logs retrieved successfully",
  "data": [
    {
      "id": 1,
      "method": "POST",
      "path": "/api/v1/external/products",
      "status_code": 201,
      "response_time": 145,
      "ip_address": "192.168.1.100",
      "user_agent": "PostmanRuntime/7.32.3",
      "access_from": "postman",
      "user_id": 5,
      "request_body": "{\"name\":\"Nasi Goreng\",\"price\":25000}",
      "response_body": "{\"status\":\"success\",\"message\":\"Product created\"}",
      "error_message": "",
      "created_at": "2026-03-28T10:30:00Z",
      "updated_at": "2026-03-28T10:30:00Z"
    },
    {
      "id": 2,
      "method": "GET",
      "path": "/api/v1/external/products",
      "status_code": 200,
      "response_time": 52,
      "ip_address": "192.168.1.101",
      "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0",
      "access_from": "website",
      "user_id": 3,
      "request_body": "",
      "response_body": "{\"status\":\"success\",\"data\":[...]}",
      "error_message": "",
      "created_at": "2026-03-28T10:25:00Z",
      "updated_at": "2026-03-28T10:25:00Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 150,
    "total_pages": 15
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "status": "error",
  "message": "Unauthorized",
  "error": "Invalid or missing token"
}
```

**Error Response (403 Forbidden):**
```json
{
  "status": "error",
  "message": "Access denied",
  "error": "Only internal users can access this resource"
}
```

---

### 2. Get API Log by ID

Mendapatkan detail log API berdasarkan ID.

**Endpoint:** `GET /logs/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | integer | Yes | ID log yang ingin diambil |

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/logs/1" \
  -H "Authorization: Bearer <token>"
```

**Success Response (200 OK):**
```json
{
  "status": "success",
  "message": "Log retrieved successfully",
  "data": {
    "id": 1,
    "method": "POST",
    "path": "/api/v1/external/products",
    "status_code": 201,
    "response_time": 145,
    "ip_address": "192.168.1.100",
    "user_agent": "PostmanRuntime/7.32.3",
    "access_from": "postman",
    "user_id": 5,
    "request_body": "{\"name\":\"Nasi Goreng\",\"price\":25000,\"category_id\":\"uuid-here\"}",
    "response_body": "{\"status\":\"success\",\"message\":\"Product created successfully\",\"data\":{\"id\":\"product-uuid\"}}",
    "error_message": "",
    "created_at": "2026-03-28T10:30:00Z",
    "updated_at": "2026-03-28T10:30:00Z"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "status": "error",
  "message": "Log not found",
  "error": "log not found"
}
```

---

## Access From Values

Field `access_from` mendeteksi sumber akses berdasarkan User-Agent:

| Value | Description |
|-------|-------------|
| `postman` | Request dari Postman |
| `mobile` | Request dari aplikasi mobile (Android/iOS) |
| `website` | Request dari web browser |
| `curl` | Request dari curl command |
| `insomnia` | Request dari Insomnia REST client |
| `httpie` | Request dari HTTPie |
| `unknown` | Sumber tidak terdeteksi |

---

## Log Fields Description

| Field | Type | Description |
|-------|------|-------------|
| id | integer | ID unik log |
| method | string | HTTP method (GET, POST, PUT, DELETE) |
| path | string | URL path yang diakses |
| status_code | integer | HTTP status code response |
| response_time | integer | Waktu response dalam milliseconds |
| ip_address | string | IP address client |
| user_agent | string | User agent string dari client |
| access_from | string | Sumber akses (postman, mobile, website, dll) |
| user_id | integer | ID user yang melakukan request (null jika tidak authenticated) |
| request_body | string | Body request (max 5000 chars) |
| response_body | string | Body response (max 5000 chars) |
| error_message | string | Error message jika ada |
| created_at | timestamp | Waktu log dibuat |
| updated_at | timestamp | Waktu log diupdate |

---

## Use Cases

### 1. Monitoring Traffic
```bash
# Lihat semua request hari ini
curl -X GET "http://localhost:8080/api/v1/logs?page=1&limit=100" \
  -H "Authorization: Bearer <token>"
```

### 2. Debug Error
```bash
# Filter hanya request yang error (status 4xx atau 5xx)
# Bisa dilakukan di aplikasi frontend dengan filter status_code
curl -X GET "http://localhost:8080/api/v1/logs" \
  -H "Authorization: Bearer <token>"
```

### 3. Analisis Performance
```bash
# Ambil semua logs untuk analisis response time
curl -X GET "http://localhost:8080/api/v1/logs?limit=1000" \
  -H "Authorization: Bearer <token>"
```

### 4. Security Audit
```bash
# Filter POST/PUT/DELETE untuk audit perubahan data
curl -X GET "http://localhost:8080/api/v1/logs?method=DELETE" \
  -H "Authorization: Bearer <token>"
```

### 5. User Activity Tracking
```bash
# Lihat detail aktivitas user tertentu
# Filter di aplikasi berdasarkan user_id
curl -X GET "http://localhost:8080/api/v1/logs" \
  -H "Authorization: Bearer <token>"
```

---

## Notes

1. **Automatic Logging**: Semua endpoint secara otomatis di-log kecuali endpoint `/api/v1/logs` itu sendiri untuk menghindari infinite loop.

2. **Async Processing**: Log disimpan secara asynchronous untuk tidak mempengaruhi response time endpoint utama.

3. **Body Truncation**: Request dan response body dibatasi maksimal 5000 karakter untuk menghemat storage.

4. **Access Control**: 
   - **Internal users** (SUPER_ADMIN, SUPPORT, FINANCE) dapat melihat **semua logs** dari semua company
   - **External users** (OWNER, ADMIN) hanya dapat melihat logs dari **company dan branch mereka sendiri**
   - External users lain (CASHIER, KITCHEN, WAITER) tidak dapat mengakses logs

5. **Performance**: Index sudah dibuat pada kolom yang sering diquery (method, path, user_id, access_from, created_at).

---

## Testing Examples

### Test dengan Postman
1. Login sebagai SUPER_ADMIN
2. Copy token dari response
3. GET `/api/v1/logs` dengan Bearer token
4. Coba filter dengan query parameters

### Test dengan cURL
```bash
# Login dulu
TOKEN=$(curl -X POST "http://localhost:8080/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@siresto.com","password":"password123"}' \
  | jq -r '.data.token')

# Get all logs
curl -X GET "http://localhost:8080/api/v1/logs" \
  -H "Authorization: Bearer $TOKEN"

# Get logs with filter
curl -X GET "http://localhost:8080/api/v1/logs?method=POST&page=1&limit=20" \
  -H "Authorization: Bearer $TOKEN"

# Get specific log
curl -X GET "http://localhost:8080/api/v1/logs/1" \
  -H "Authorization: Bearer $TOKEN"
```
