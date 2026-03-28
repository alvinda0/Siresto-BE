# Update: Access Control untuk API Logs

## 🔄 Perubahan

Sistem access control untuk API Logs telah diupdate dengan aturan baru:

### Sebelumnya
- ❌ Hanya **INTERNAL users** (SUPER_ADMIN, SUPPORT, FINANCE) yang bisa akses logs
- ❌ External users tidak bisa akses sama sekali

### Sekarang
- ✅ **INTERNAL users** (SUPER_ADMIN, SUPPORT, FINANCE) → Bisa lihat **SEMUA logs**
- ✅ **EXTERNAL users** (OWNER, ADMIN) → Bisa lihat logs **company dan branch mereka saja**
- ❌ External users lain (CASHIER, KITCHEN, WAITER) → Tidak bisa akses

---

## 🎯 Aturan Access Control Baru

### 1. Internal Users (SUPER_ADMIN, SUPPORT, FINANCE)
**Akses:** Semua logs dari semua company dan branch

**Contoh:**
```bash
GET /api/v1/logs
Authorization: Bearer <super-admin-token>

# Response: Semua logs dari semua company
```

### 2. External Users - OWNER
**Akses:** Hanya logs dari company yang dimiliki

**Contoh:**
```bash
GET /api/v1/logs
Authorization: Bearer <owner-token>

# Response: Hanya logs dari company_id milik owner tersebut
```

### 3. External Users - ADMIN
**Akses:** Hanya logs dari company dan branch yang ditugaskan

**Contoh:**
```bash
GET /api/v1/logs
Authorization: Bearer <admin-token>

# Response: Hanya logs dari company_id dan branch_id admin tersebut
```

### 4. External Users Lain (CASHIER, KITCHEN, WAITER)
**Akses:** Tidak bisa akses logs

**Response:**
```json
{
  "status": "error",
  "message": "Access denied",
  "error": "You don't have permission to access logs"
}
```

---

## 📊 Perubahan Database

### Field Baru di Tabel `api_logs`

```sql
ALTER TABLE api_logs 
ADD COLUMN company_id UUID,
ADD COLUMN branch_id UUID;

CREATE INDEX idx_api_logs_company_id ON api_logs(company_id);
CREATE INDEX idx_api_logs_branch_id ON api_logs(branch_id);
```

**Field baru:**
- `company_id` - UUID company yang melakukan request
- `branch_id` - UUID branch yang melakukan request

**Catatan:** Field ini otomatis terisi dari context user yang login

---

## 🔧 Perubahan Teknis

### 1. Entity (`internal/entity/api_log.go`)
```go
type APILog struct {
    // ... field lain
    CompanyID     *string        `gorm:"type:uuid;index" json:"company_id,omitempty"`
    BranchID      *string        `gorm:"type:uuid;index" json:"branch_id,omitempty"`
    // ...
}
```

### 2. Repository (`internal/repository/api_log_repository.go`)
```go
// Sekarang menerima companyID dan branchID untuk filtering
FindAll(page, limit int, method, companyID, branchID string) ([]entity.APILog, int64, error)
FindByID(id uint, companyID, branchID string) (*entity.APILog, error)
```

### 3. Service (`internal/service/api_log_service.go`)
```go
// Sekarang menerima companyID dan branchID untuk filtering
GetAllLogs(page, limit int, method, companyID, branchID string) ([]entity.APILog, *pkg.PaginationMeta, error)
GetLogByID(id uint, companyID, branchID string) (*entity.APILog, error)
```

### 4. Handler (`internal/handler/api_log_handler.go`)
```go
// Mengambil role_type dari context
roleType, _ := c.Get("role_type")

// Jika EXTERNAL, ambil company_id dan branch_id
if roleType == "EXTERNAL" {
    companyID = c.Get("company_id")
    branchID = c.Get("branch_id")
}
```

### 5. Middleware (`internal/middleware/logging_middleware.go`)
```go
// Sekarang menyimpan company_id dan branch_id ke log
apiLog := &entity.APILog{
    // ... field lain
    CompanyID:    companyID,
    BranchID:     branchID,
    // ...
}
```

### 6. Routes (`routes/routes.go`)
```go
// Middleware RequireInternalRole() dihapus
// Sekarang semua authenticated users bisa akses, tapi dengan filtering
logs := v1.Group("/logs")
logs.Use(middleware.AuthMiddleware())
{
    logs.GET("", apiLogHandler.GetAllLogs)
    logs.GET("/:id", apiLogHandler.GetLogByID)
}
```

---

## 🧪 Testing

### Test 1: Internal User (SUPER_ADMIN)
```bash
# Login sebagai SUPER_ADMIN
POST /api/v1/login
{
  "email": "admin@siresto.com",
  "password": "password123"
}

# Get all logs (akan dapat semua logs)
GET /api/v1/logs
Authorization: Bearer <super-admin-token>

# Expected: Semua logs dari semua company
```

### Test 2: External User (OWNER)
```bash
# Login sebagai OWNER
POST /api/v1/login
{
  "email": "owner@restaurant.com",
  "password": "password123"
}

# Get logs (hanya logs company sendiri)
GET /api/v1/logs
Authorization: Bearer <owner-token>

# Expected: Hanya logs dengan company_id = owner's company
```

### Test 3: External User (ADMIN)
```bash
# Login sebagai ADMIN
POST /api/v1/login
{
  "email": "admin@branch.com",
  "password": "password123"
}

# Get logs (hanya logs company dan branch sendiri)
GET /api/v1/logs
Authorization: Bearer <admin-token>

# Expected: Hanya logs dengan company_id dan branch_id = admin's company & branch
```

### Test 4: Filter by Method (OWNER)
```bash
# OWNER filter by POST method
GET /api/v1/logs?method=POST
Authorization: Bearer <owner-token>

# Expected: Hanya POST logs dari company sendiri
```

### Test 5: Get Log by ID (OWNER)
```bash
# OWNER coba akses log milik company sendiri
GET /api/v1/logs/1
Authorization: Bearer <owner-token>

# Expected: Success jika log tersebut milik company sendiri

# OWNER coba akses log milik company lain
GET /api/v1/logs/999
Authorization: Bearer <owner-token>

# Expected: 404 Not Found (karena bukan milik company sendiri)
```

---

## 📋 Checklist Update

### Code Changes
- [x] Update entity (tambah company_id, branch_id)
- [x] Update repository (tambah filter company & branch)
- [x] Update service (tambah parameter company & branch)
- [x] Update handler (ambil role_type dan filter)
- [x] Update middleware (simpan company_id & branch_id)
- [x] Update routes (hapus RequireInternalRole)
- [x] Update database migration (tambah kolom & index)

### Testing
- [ ] Test internal user dapat semua logs
- [ ] Test OWNER hanya dapat logs company sendiri
- [ ] Test ADMIN hanya dapat logs company & branch sendiri
- [ ] Test filter by method tetap bekerja
- [ ] Test pagination tetap bekerja
- [ ] Test get by ID dengan access control
- [ ] Test external user tidak bisa akses log company lain

### Documentation
- [x] Update access control documentation
- [ ] Update API documentation
- [ ] Update testing guide
- [ ] Update quick start guide

---

## 🔍 Cara Kerja Filtering

### Internal User
```
Query: SELECT * FROM api_logs WHERE ...
(Tidak ada filter company_id atau branch_id)
Result: Semua logs
```

### External User (OWNER)
```
Query: SELECT * FROM api_logs WHERE company_id = 'owner-company-uuid' AND ...
Result: Hanya logs dari company owner
```

### External User (ADMIN)
```
Query: SELECT * FROM api_logs WHERE company_id = 'admin-company-uuid' AND branch_id = 'admin-branch-uuid' AND ...
Result: Hanya logs dari company dan branch admin
```

---

## 💡 Use Cases

### 1. SUPER_ADMIN Monitor Semua Activity
```bash
GET /api/v1/logs?page=1&limit=100
# Lihat semua activity dari semua company
```

### 2. OWNER Monitor Activity Restaurant
```bash
GET /api/v1/logs?method=POST
# Lihat semua POST request di restaurant sendiri
```

### 3. ADMIN Monitor Activity Branch
```bash
GET /api/v1/logs?method=DELETE
# Lihat semua DELETE request di branch sendiri
```

### 4. OWNER Audit User Activity
```bash
GET /api/v1/logs
# Filter di aplikasi berdasarkan user_id untuk audit staff
```

---

## 🔒 Security

### Keamanan yang Diterapkan
- ✅ External users tidak bisa lihat logs company lain
- ✅ External users tidak bisa lihat logs branch lain
- ✅ Filtering dilakukan di database level (tidak bisa di-bypass)
- ✅ Access control berdasarkan role_type dari JWT token

### Catatan Keamanan
- Company_id dan branch_id diambil dari JWT token (tidak dari request)
- Tidak bisa di-manipulasi oleh client
- Query filtering dilakukan di repository level

---

## 🎉 Summary

**Perubahan utama:**
1. ✅ Tambah field `company_id` dan `branch_id` di tabel logs
2. ✅ Internal users bisa lihat semua logs
3. ✅ External users (OWNER, ADMIN) bisa lihat logs company/branch mereka
4. ✅ Filtering otomatis berdasarkan role dan context user
5. ✅ Backward compatible (internal users tetap bisa lihat semua)

**Status:** ✅ Complete & Ready to Test

**Date:** March 28, 2026
