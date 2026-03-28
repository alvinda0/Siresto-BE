# API Logging System

Sistem logging otomatis untuk tracking dan monitoring semua request API yang masuk ke sistem SIRESTO.

## 📋 Overview

Sistem ini secara otomatis mencatat setiap request yang masuk ke API, termasuk:
- HTTP Method (GET, POST, PUT, DELETE)
- Path/Endpoint yang diakses
- Status code response
- Response time (dalam milliseconds)
- IP Address client
- User Agent
- Sumber akses (website, mobile, postman, dll)
- User ID (jika authenticated)
- Request body
- Response body
- Error message (jika ada)

## 🎯 Features

### 1. Automatic Logging
Semua endpoint secara otomatis di-log tanpa perlu konfigurasi tambahan di setiap handler.

### 2. Access Source Detection
Sistem dapat mendeteksi sumber akses berdasarkan User-Agent:
- **Postman**: Request dari Postman REST client
- **Website**: Request dari web browser (Chrome, Firefox, Safari, dll)
- **Mobile**: Request dari aplikasi mobile (Android/iOS)
- **cURL**: Request dari command line curl
- **Insomnia**: Request dari Insomnia REST client
- **HTTPie**: Request dari HTTPie
- **Unknown**: Sumber tidak terdeteksi

### 3. Pagination & Filtering
- Pagination dengan parameter `page` dan `limit`
- Filter berdasarkan HTTP method (GET, POST, PUT, DELETE)
- Sorting berdasarkan waktu (terbaru dulu)

### 4. Performance Optimized
- Logging dilakukan secara asynchronous
- Tidak mempengaruhi response time endpoint utama
- Request/response body di-truncate jika terlalu besar (max 5000 chars)

### 5. Access Control
Hanya user dengan role INTERNAL (SUPER_ADMIN, SUPPORT, FINANCE) yang bisa mengakses logs.

## 🏗️ Architecture

```
Request → Logging Middleware → Handler → Response
              ↓ (async)
          API Log Service
              ↓
          API Log Repository
              ↓
          Database (api_logs table)
```

## 📁 File Structure

```
internal/
├── entity/
│   └── api_log.go              # Entity/Model untuk API Log
├── repository/
│   └── api_log_repository.go   # Database operations
├── service/
│   └── api_log_service.go      # Business logic
├── handler/
│   └── api_log_handler.go      # HTTP handlers
└── middleware/
    └── logging_middleware.go   # Middleware untuk auto-logging

routes/
└── routes.go                   # Route configuration

config/
└── config.go                   # Database migration
```

## 🗄️ Database Schema

```sql
CREATE TABLE api_logs (
    id SERIAL PRIMARY KEY,
    method VARCHAR(10) NOT NULL,
    path VARCHAR(255) NOT NULL,
    status_code INTEGER NOT NULL,
    response_time BIGINT NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    access_from VARCHAR(50),
    user_id INTEGER,
    request_body TEXT,
    response_body TEXT,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Indexes for better query performance
CREATE INDEX idx_api_logs_method ON api_logs(method);
CREATE INDEX idx_api_logs_path ON api_logs(path);
CREATE INDEX idx_api_logs_user_id ON api_logs(user_id);
CREATE INDEX idx_api_logs_access_from ON api_logs(access_from);
CREATE INDEX idx_api_logs_created_at ON api_logs(created_at);
CREATE INDEX idx_api_logs_deleted_at ON api_logs(deleted_at);
```

## 🚀 Usage

### Setup

1. **Migration sudah otomatis** saat server start pertama kali
2. **Middleware sudah terpasang** di `routes.go`
3. **Tidak perlu konfigurasi tambahan**

### Accessing Logs

#### Get All Logs
```bash
GET /api/v1/logs?page=1&limit=10&method=POST
Authorization: Bearer <token>
```

#### Get Log by ID
```bash
GET /api/v1/logs/:id
Authorization: Bearer <token>
```

## 📊 Use Cases

### 1. Traffic Monitoring
Monitor berapa banyak request yang masuk per endpoint, per method, per user.

### 2. Performance Analysis
Analisis response time untuk identifikasi endpoint yang lambat.

### 3. Error Debugging
Track error yang terjadi dengan melihat request/response body dan error message.

### 4. Security Audit
Monitor aktivitas user, terutama untuk operasi yang mengubah data (POST, PUT, DELETE).

### 5. User Activity Tracking
Lihat aktivitas spesifik user berdasarkan user_id.

### 6. Access Source Analytics
Analisis dari mana user mengakses API (mobile app, website, atau tools lain).

## 🔒 Security

### Access Control
- Hanya user dengan role INTERNAL yang bisa akses logs
- External users (OWNER, ADMIN, CASHIER, dll) tidak bisa akses

### Data Privacy
- Request/response body di-truncate untuk menghemat storage
- Sensitive data sebaiknya tidak di-log (password, token, dll)
- Bisa dikustomisasi untuk exclude certain paths

### Self-Exclusion
Endpoint `/api/v1/logs` tidak di-log untuk menghindari infinite loop.

## ⚙️ Configuration

### Exclude Certain Paths
Edit `logging_middleware.go` untuk exclude path tertentu:

```go
func LoggingMiddleware(logService service.APILogService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip logging for certain paths
        if strings.HasPrefix(c.Request.URL.Path, "/api/logs") ||
           strings.HasPrefix(c.Request.URL.Path, "/uploads") {
            c.Next()
            return
        }
        // ... rest of the code
    }
}
```

### Adjust Body Size Limit
Edit `logging_middleware.go`:

```go
// Change from 5000 to your desired limit
if len(responseBody) > 10000 {
    responseBody = responseBody[:10000] + "... (truncated)"
}
```

### Custom Access From Detection
Edit `determineAccessFrom()` function di `logging_middleware.go`:

```go
func determineAccessFrom(userAgent string) string {
    userAgent = strings.ToLower(userAgent)
    
    // Add your custom detection
    if strings.Contains(userAgent, "myapp") {
        return "my-custom-app"
    }
    
    // ... rest of the code
}
```

## 📈 Performance Considerations

### Async Logging
Logging dilakukan secara asynchronous menggunakan goroutine:

```go
go func() {
    _ = logService.CreateLog(apiLog)
}()
```

Ini memastikan logging tidak mempengaruhi response time endpoint utama.

### Database Indexes
Index sudah dibuat pada kolom yang sering di-query:
- `method` - untuk filter by method
- `path` - untuk filter by path
- `user_id` - untuk filter by user
- `access_from` - untuk analytics
- `created_at` - untuk sorting

### Storage Management
Pertimbangkan untuk:
- Implement log rotation (hapus log lama secara periodik)
- Archive old logs ke cold storage
- Implement retention policy (misal: simpan log 90 hari terakhir)

## 🧪 Testing

Lihat file `API_LOGS_TESTING.md` untuk panduan testing lengkap.

Quick test:
```bash
# 1. Login
curl -X POST "http://localhost:8080/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@siresto.com","password":"password123"}'

# 2. Get logs
curl -X GET "http://localhost:8080/api/v1/logs" \
  -H "Authorization: Bearer <token>"
```

## 📚 Documentation

- `API_LOGS_DOCUMENTATION.md` - API documentation lengkap
- `API_LOGS_TESTING.md` - Testing guide
- `API_LOGS_README.md` - Overview (file ini)

## 🔮 Future Enhancements

Beberapa enhancement yang bisa ditambahkan:

1. **Dashboard Analytics**
   - Grafik traffic per endpoint
   - Response time trends
   - Error rate monitoring

2. **Real-time Monitoring**
   - WebSocket untuk real-time log streaming
   - Alert system untuk error spike

3. **Advanced Filtering**
   - Filter by date range
   - Filter by status code
   - Filter by response time
   - Full-text search di request/response body

4. **Export Functionality**
   - Export logs ke CSV/Excel
   - Export untuk audit purposes

5. **Log Aggregation**
   - Summary statistics per endpoint
   - Daily/weekly/monthly reports

6. **Retention Policy**
   - Auto-delete old logs
   - Archive to S3/cloud storage

## 🤝 Contributing

Jika ingin menambahkan fitur atau fix bug:
1. Pastikan tidak break existing functionality
2. Update documentation
3. Add tests
4. Follow Go best practices

## 📝 License

Internal use only - SIRESTO Platform
