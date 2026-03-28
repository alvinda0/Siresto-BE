# Update: UUID sebagai Primary Key untuk API Logs

## 🔄 Perubahan

Primary key tabel `api_logs` diubah dari integer (SERIAL) menjadi UUID.

### Sebelumnya
```sql
CREATE TABLE api_logs (
    id SERIAL PRIMARY KEY,  -- ❌ Integer auto-increment
    ...
);
```

### Sekarang
```sql
CREATE TABLE api_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- ✅ UUID
    ...
);
```

---

## 🎯 Alasan Perubahan

### 1. Konsistensi dengan Tabel Lain
Semua tabel lain di sistem menggunakan UUID sebagai primary key:
- `users` → UUID
- `companies` → UUID
- `branches` → UUID
- `categories` → UUID
- `products` → UUID
- `roles` → UUID

### 2. Keamanan
- UUID tidak predictable (tidak bisa ditebak)
- Integer sequential bisa ditebak: `/api/logs/1`, `/api/logs/2`, dll
- UUID lebih aman: `/api/logs/550e8400-e29b-41d4-a716-446655440000`

### 3. Distributed System Ready
- UUID bisa di-generate di client/server tanpa collision
- Cocok untuk distributed logging system
- Tidak perlu koordinasi antar server untuk generate ID

### 4. Better for Merging Data
- Jika ada multiple database yang perlu di-merge
- UUID tidak akan conflict
- Integer bisa conflict (ID 1 di DB A vs ID 1 di DB B)

---

## 📊 Perubahan Detail

### 1. Entity (`internal/entity/api_log.go`)

**Before:**
```go
type APILog struct {
    ID            uint           `gorm:"primaryKey" json:"id"`
    UserID        *uint          `gorm:"index" json:"user_id,omitempty"`
    CompanyID     *string        `gorm:"type:uuid;index" json:"company_id,omitempty"`
    BranchID      *string        `gorm:"type:uuid;index" json:"branch_id,omitempty"`
    // ...
}
```

**After:**
```go
import "github.com/google/uuid"

type APILog struct {
    ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    UserID        *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
    CompanyID     *uuid.UUID     `gorm:"type:uuid;index" json:"company_id,omitempty"`
    BranchID      *uuid.UUID     `gorm:"type:uuid;index" json:"branch_id,omitempty"`
    // ...
}
```

### 2. DTO (`internal/entity/api_log_dto.go`)

**Before:**
```go
type APILogListDTO struct {
    ID            uint       `json:"id"`
    UserID        *uint      `json:"user_id,omitempty"`
    CompanyID     *string    `json:"company_id,omitempty"`
    BranchID      *string    `json:"branch_id,omitempty"`
    // ...
}

type APILogDetailDTO struct {
    ID            uint       `json:"id"`
    UserID        *uint      `json:"user_id,omitempty"`
    CompanyID     *string    `json:"company_id,omitempty"`
    BranchID      *string    `json:"branch_id,omitempty"`
    // ...
}
```

**After:**
```go
import "github.com/google/uuid"

type APILogListDTO struct {
    ID            uuid.UUID  `json:"id"`
    UserID        *uuid.UUID `json:"user_id,omitempty"`
    CompanyID     *uuid.UUID `json:"company_id,omitempty"`
    BranchID      *uuid.UUID `json:"branch_id,omitempty"`
    // ...
}

type APILogDetailDTO struct {
    ID            uuid.UUID  `json:"id"`
    UserID        *uuid.UUID `json:"user_id,omitempty"`
    CompanyID     *uuid.UUID `json:"company_id,omitempty"`
    BranchID      *uuid.UUID `json:"branch_id,omitempty"`
    // ...
}
```

### 3. Repository (`internal/repository/api_log_repository.go`)

**Before:**
```go
type APILogRepository interface {
    FindByID(id uint, companyID, branchID string) (*entity.APILog, error)
}

func (r *apiLogRepository) FindByID(id uint, companyID, branchID string) (*entity.APILog, error) {
    var log entity.APILog
    // ...
    err := query.First(&log, id).Error
    // ...
}
```

**After:**
```go
type APILogRepository interface {
    FindByID(id string, companyID, branchID string) (*entity.APILog, error)
}

func (r *apiLogRepository) FindByID(id string, companyID, branchID string) (*entity.APILog, error) {
    var log entity.APILog
    // ...
    err := query.Where("id = ?", id).First(&log).Error
    // ...
}
```

### 4. Service (`internal/service/api_log_service.go`)

**Before:**
```go
type APILogService interface {
    GetLogByID(id uint, companyID, branchID string) (*entity.APILogDetailDTO, error)
}

func (s *apiLogService) GetLogByID(id uint, companyID, branchID string) (*entity.APILogDetailDTO, error) {
    // ...
}
```

**After:**
```go
type APILogService interface {
    GetLogByID(id string, companyID, branchID string) (*entity.APILogDetailDTO, error)
}

func (s *apiLogService) GetLogByID(id string, companyID, branchID string) (*entity.APILogDetailDTO, error) {
    // ...
}
```

### 5. Handler (`internal/handler/api_log_handler.go`)

**Before:**
```go
import "strconv"

func (h *APILogHandler) GetLogByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid log ID", err.Error())
        return
    }
    
    log, err := h.service.GetLogByID(uint(id), companyID, branchID)
    // ...
}
```

**After:**
```go
import "github.com/google/uuid"

func (h *APILogHandler) GetLogByID(c *gin.Context) {
    id := c.Param("id")
    
    // Validate UUID format
    if _, err := uuid.Parse(id); err != nil {
        pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid log ID format", "ID must be a valid UUID")
        return
    }
    
    log, err := h.service.GetLogByID(id, companyID, branchID)
    // ...
}
```

### 6. Middleware (`internal/middleware/logging_middleware.go`)

**Before:**
```go
var userID *uint
var companyID *string
var branchID *string

if id, exists := c.Get("user_id"); exists {
    if uid, ok := id.(uint); ok {
        userID = &uid
    }
}

if cid, exists := c.Get("company_id"); exists {
    if cidStr, ok := cid.(string); ok {
        companyID = &cidStr
    }
}
```

**After:**
```go
import "github.com/google/uuid"

var userID *uuid.UUID
var companyID *uuid.UUID
var branchID *uuid.UUID

if id, exists := c.Get("user_id"); exists {
    if uid, ok := id.(uuid.UUID); ok {
        userID = &uid
    }
}

if cid, exists := c.Get("company_id"); exists {
    if cidStr, ok := cid.(string); ok {
        if parsed, err := uuid.Parse(cidStr); err == nil {
            companyID = &parsed
        }
    }
}
```

### 7. Database Migration (`config/config.go`)

**Before:**
```sql
CREATE TABLE IF NOT EXISTS api_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    company_id UUID,
    branch_id UUID,
    ...
);
```

**After:**
```sql
CREATE TABLE IF NOT EXISTS api_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    company_id UUID,
    branch_id UUID,
    ...
);
```

---

## 📋 API Changes

### GET `/api/v1/logs/:id`

**Before:**
```bash
GET /api/v1/logs/1
GET /api/v1/logs/123
```

**After:**
```bash
GET /api/v1/logs/550e8400-e29b-41d4-a716-446655440000
GET /api/v1/logs/6ba7b810-9dad-11d1-80b4-00c04fd430c8
```

### Response Format

**Before:**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "user_id": 123,
    "company_id": "550e8400-e29b-41d4-a716-446655440000",
    "branch_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
  }
}
```

**After:**
```json
{
  "status": "success",
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "company_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "branch_id": "7c9e6679-7425-40de-944b-e07fc1f90ae7"
  }
}
```

---

## 🧪 Testing

### Test 1: Get All Logs (Verify UUID format)

**Request:**
```bash
GET http://localhost:8080/api/v1/logs
Authorization: Bearer <token>
```

**Expected Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",  // ✅ UUID format
      "method": "GET",
      "path": "/api/v1/products"
    }
  ]
}
```

**Validation:**
- ✅ Field `id` adalah UUID (format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)
- ✅ Field `user_id`, `company_id`, `branch_id` juga UUID

### Test 2: Get Log by UUID

**Request:**
```bash
GET http://localhost:8080/api/v1/logs/a1b2c3d4-e5f6-7890-abcd-ef1234567890
Authorization: Bearer <token>
```

**Expected Response:**
```json
{
  "status": "success",
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "method": "GET",
    "path": "/api/v1/products"
  }
}
```

**Validation:**
- ✅ Bisa akses log dengan UUID
- ✅ Response berisi detail lengkap

### Test 3: Invalid UUID Format

**Request:**
```bash
GET http://localhost:8080/api/v1/logs/123
Authorization: Bearer <token>
```

**Expected Response:**
```json
{
  "status": "error",
  "message": "Invalid log ID format",
  "error": "ID must be a valid UUID"
}
```

**Validation:**
- ✅ Return 400 Bad Request
- ✅ Error message jelas

### Test 4: Non-existent UUID

**Request:**
```bash
GET http://localhost:8080/api/v1/logs/00000000-0000-0000-0000-000000000000
Authorization: Bearer <token>
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
- ✅ Return 404 Not Found
- ✅ UUID valid tapi data tidak ada

---

## ⚠️ Breaking Changes

### 1. API Endpoint Parameter
- ❌ Old: `/api/v1/logs/1` (integer)
- ✅ New: `/api/v1/logs/550e8400-e29b-41d4-a716-446655440000` (UUID)

### 2. Response Format
- Field `id` sekarang UUID string, bukan integer
- Field `user_id`, `company_id`, `branch_id` sekarang UUID, bukan integer/string

### 3. Database Schema
- Tabel `api_logs` perlu di-recreate atau migrate
- Data lama dengan integer ID tidak compatible

---

## 🔄 Migration Strategy

### Option 1: Drop and Recreate (Development)
```sql
-- Backup data jika perlu
-- DROP TABLE api_logs;

-- Restart server, table akan dibuat ulang dengan UUID
```

### Option 2: Migrate Existing Data (Production)
```sql
-- 1. Rename old table
ALTER TABLE api_logs RENAME TO api_logs_old;

-- 2. Create new table with UUID
CREATE TABLE api_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    method VARCHAR(10) NOT NULL,
    path VARCHAR(255) NOT NULL,
    status_code INTEGER NOT NULL,
    response_time BIGINT NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    access_from VARCHAR(50),
    user_id UUID,
    company_id UUID,
    branch_id UUID,
    request_body TEXT,
    response_body TEXT,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- 3. Migrate data (if needed)
-- Note: user_id perlu di-convert dari integer ke UUID
-- Jika tidak ada mapping, skip user_id

INSERT INTO api_logs (
    method, path, status_code, response_time,
    ip_address, user_agent, access_from,
    company_id, branch_id,
    request_body, response_body, error_message,
    created_at, updated_at, deleted_at
)
SELECT 
    method, path, status_code, response_time,
    ip_address, user_agent, access_from,
    company_id::uuid, branch_id::uuid,
    request_body, response_body, error_message,
    created_at, updated_at, deleted_at
FROM api_logs_old;

-- 4. Drop old table
DROP TABLE api_logs_old;
```

---

## ✅ Checklist

### Code Changes
- [x] Update entity APILog (id, user_id, company_id, branch_id → UUID)
- [x] Update DTO (APILogListDTO, APILogDetailDTO)
- [x] Update repository interface & implementation
- [x] Update service interface & implementation
- [x] Update handler (validate UUID, remove strconv)
- [x] Update middleware (parse UUID from string)
- [x] Update database migration
- [x] Add uuid import where needed

### Testing
- [ ] Test GET /logs returns UUID
- [ ] Test GET /logs/:id with valid UUID
- [ ] Test GET /logs/:id with invalid UUID format
- [ ] Test GET /logs/:id with non-existent UUID
- [ ] Test filtering still works with UUID

### Documentation
- [x] Create UPDATE_UUID_PRIMARY_KEY.md
- [ ] Update API_LOGS_DOCUMENTATION.md
- [ ] Update testing guides

---

## 🎉 Summary

**Perubahan:**
1. ✅ Primary key `id` dari integer → UUID
2. ✅ Field `user_id` dari integer → UUID
3. ✅ Field `company_id` dari string → UUID
4. ✅ Field `branch_id` dari string → UUID
5. ✅ API endpoint parameter dari integer → UUID
6. ✅ Validation UUID format di handler

**Benefits:**
- 🔒 Lebih aman (tidak predictable)
- 🔄 Konsisten dengan tabel lain
- 🌐 Distributed system ready
- 🔀 Merge-friendly

**Status:** ✅ Complete & Ready to Test

**Date:** March 28, 2026
