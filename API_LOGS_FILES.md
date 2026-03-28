# API Logging - File Structure

Daftar file yang dibuat untuk fitur API Logging.

## 📁 Core Files

### 1. Entity/Model
**File:** `internal/entity/api_log.go`
- Definisi struktur data API Log
- Field: id, method, path, status_code, response_time, ip_address, user_agent, access_from, user_id, request_body, response_body, error_message, timestamps

### 2. Repository Layer
**File:** `internal/repository/api_log_repository.go`
- Interface dan implementasi database operations
- Methods:
  - `Create()` - Simpan log baru
  - `FindAll()` - Get semua logs dengan pagination dan filter
  - `FindByID()` - Get log by ID

### 3. Service Layer
**File:** `internal/service/api_log_service.go`
- Business logic untuk API logs
- Methods:
  - `CreateLog()` - Create log entry
  - `GetAllLogs()` - Get logs dengan pagination dan filter
  - `GetLogByID()` - Get single log

### 4. Handler Layer
**File:** `internal/handler/api_log_handler.go`
- HTTP handlers untuk API endpoints
- Endpoints:
  - `GET /api/v1/logs` - Get all logs
  - `GET /api/v1/logs/:id` - Get log by ID

### 5. Middleware
**File:** `internal/middleware/logging_middleware.go`
- Middleware untuk automatic logging
- Capture request/response
- Detect access source
- Save log asynchronously

## 🔧 Configuration Files

### 6. Routes
**File:** `routes/routes.go`
- Integrasi logging middleware
- Setup API log endpoints
- Access control (internal users only)

### 7. Database Migration
**File:** `config/config.go`
- Create `api_logs` table
- Create indexes untuk performance
- Auto-run saat server start

## 📚 Documentation Files

### 8. API Documentation
**File:** `API_LOGS_DOCUMENTATION.md`
- Complete API documentation
- Request/response examples
- Query parameters
- Error handling
- Use cases

### 9. Testing Guide
**File:** `API_LOGS_TESTING.md`
- Step-by-step testing guide
- Test cases untuk semua scenarios
- Expected responses
- Validation checklist

### 10. README
**File:** `API_LOGS_README.md`
- Overview sistem logging
- Architecture
- Features
- Configuration
- Performance considerations
- Future enhancements

### 11. Quick Start
**File:** `QUICK_START_LOGS.md`
- Panduan cepat untuk mulai menggunakan
- Basic commands
- Common use cases
- Tips & tricks

### 12. File Structure (This File)
**File:** `API_LOGS_FILES.md`
- Daftar semua file yang dibuat
- Penjelasan setiap file

## 🗂️ File Tree

```
backend-golang/
├── internal/
│   ├── entity/
│   │   └── api_log.go                    # ✅ NEW
│   ├── repository/
│   │   └── api_log_repository.go         # ✅ NEW
│   ├── service/
│   │   └── api_log_service.go            # ✅ NEW
│   ├── handler/
│   │   └── api_log_handler.go            # ✅ NEW
│   └── middleware/
│       └── logging_middleware.go         # ✅ NEW
├── routes/
│   └── routes.go                         # ✏️ MODIFIED
├── config/
│   └── config.go                         # ✏️ MODIFIED
├── API_LOGS_DOCUMENTATION.md             # ✅ NEW
├── API_LOGS_TESTING.md                   # ✅ NEW
├── API_LOGS_README.md                    # ✅ NEW
├── QUICK_START_LOGS.md                   # ✅ NEW
└── API_LOGS_FILES.md                     # ✅ NEW (this file)
```

## 📊 Summary

### New Files Created: 10
1. `internal/entity/api_log.go`
2. `internal/repository/api_log_repository.go`
3. `internal/service/api_log_service.go`
4. `internal/handler/api_log_handler.go`
5. `internal/middleware/logging_middleware.go`
6. `API_LOGS_DOCUMENTATION.md`
7. `API_LOGS_TESTING.md`
8. `API_LOGS_README.md`
9. `QUICK_START_LOGS.md`
10. `API_LOGS_FILES.md`

### Modified Files: 2
1. `routes/routes.go` - Added logging middleware & log endpoints
2. `config/config.go` - Added api_logs table migration

## 🎯 What Each Layer Does

### Entity Layer
Definisi struktur data yang akan disimpan di database.

### Repository Layer
Komunikasi langsung dengan database (CRUD operations).

### Service Layer
Business logic, validasi, dan orchestration.

### Handler Layer
HTTP request/response handling, parsing parameters.

### Middleware Layer
Intercept semua requests, capture data, save logs.

## 🔄 Data Flow

```
Request
  ↓
Logging Middleware (capture request)
  ↓
Handler (process request)
  ↓
Service (business logic)
  ↓
Repository (database)
  ↓
Response
  ↓
Logging Middleware (capture response)
  ↓
Save Log (async)
```

## 🚀 Next Steps

1. ✅ Jalankan server: `go run cmd/server/main.go`
2. ✅ Database migration otomatis jalan
3. ✅ Middleware otomatis aktif
4. ✅ Hit beberapa endpoint untuk generate logs
5. ✅ Access logs via `/api/v1/logs`

## 💡 Tips

- Semua file sudah siap pakai, tidak perlu modifikasi
- Middleware sudah terpasang global di routes
- Database migration otomatis saat server start
- Logs disimpan async, tidak affect performance

## 🔍 Where to Look

- **Mau lihat struktur data?** → `internal/entity/api_log.go`
- **Mau lihat query database?** → `internal/repository/api_log_repository.go`
- **Mau lihat business logic?** → `internal/service/api_log_service.go`
- **Mau lihat API endpoints?** → `internal/handler/api_log_handler.go`
- **Mau lihat cara logging?** → `internal/middleware/logging_middleware.go`
- **Mau test API?** → `API_LOGS_TESTING.md`
- **Mau quick start?** → `QUICK_START_LOGS.md`
