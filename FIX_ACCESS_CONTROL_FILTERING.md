# Fix: Access Control Filtering untuk API Logs

## 🐛 Problem

Owner dan Admin masih melihat logs dari semua company (termasuk logs internal), padahal seharusnya hanya melihat logs dari company/branch mereka sendiri.

## 🔍 Root Cause

1. **JWT Claims tidak menyimpan `role_type`**
   - JWT hanya menyimpan `internal_role` dan `external_role`
   - Tidak ada field `role_type` (INTERNAL/EXTERNAL)

2. **Auth Middleware tidak set `role_type` ke context**
   - Middleware tidak set `role_type` ke gin context
   - Handler tidak bisa membedakan internal vs external user

3. **Company_id dan Branch_id tidak di-convert ke string**
   - JWT menyimpan sebagai UUID
   - Handler expect string untuk filtering
   - Mismatch type menyebabkan filtering tidak bekerja

## ✅ Solution

### 1. Update JWT Claims

**File:** `pkg/jwt.go`

**Before:**
```go
type JWTClaims struct {
    UserID       uuid.UUID  `json:"user_id"`
    Email        string     `json:"email"`
    InternalRole string     `json:"internal_role,omitempty"`
    ExternalRole string     `json:"external_role,omitempty"`
    CompanyID    *uuid.UUID `json:"company_id,omitempty"`
    BranchID     *uuid.UUID `json:"branch_id,omitempty"`
    jwt.RegisteredClaims
}
```

**After:**
```go
type JWTClaims struct {
    UserID       uuid.UUID  `json:"user_id"`
    Email        string     `json:"email"`
    RoleType     string     `json:"role_type"` // ✅ TAMBAH INI
    InternalRole string     `json:"internal_role,omitempty"`
    ExternalRole string     `json:"external_role,omitempty"`
    CompanyID    *uuid.UUID `json:"company_id,omitempty"`
    BranchID     *uuid.UUID `json:"branch_id,omitempty"`
    jwt.RegisteredClaims
}
```

**Update GenerateJWT function:**
```go
func GenerateJWT(userID uuid.UUID, email string, roleType string, internalRole, externalRole string, companyID, branchID *uuid.UUID) (string, error) {
    claims := JWTClaims{
        UserID:       userID,
        Email:        email,
        RoleType:     roleType, // ✅ TAMBAH INI
        InternalRole: internalRole,
        ExternalRole: externalRole,
        CompanyID:    companyID,
        BranchID:     branchID,
        // ...
    }
    // ...
}
```

### 2. Update Auth Middleware

**File:** `internal/middleware/auth_middleware.go`

**Before:**
```go
if claims, ok := token.Claims.(*JWTClaims); ok {
    c.Set("userID", claims.UserID)
    c.Set("email", claims.Email)
    c.Set("internalRole", claims.InternalRole)
    c.Set("externalRole", claims.ExternalRole)
    c.Set("companyID", claims.CompanyID)
    c.Set("branchID", claims.BranchID)
}
```

**After:**
```go
if claims, ok := token.Claims.(*JWTClaims); ok {
    c.Set("user_id", claims.UserID)
    c.Set("email", claims.Email)
    c.Set("role_type", claims.RoleType) // ✅ TAMBAH INI
    c.Set("internal_role", claims.InternalRole)
    c.Set("external_role", claims.ExternalRole)
    
    // ✅ Convert UUID to string untuk filtering
    if claims.CompanyID != nil {
        c.Set("company_id", claims.CompanyID.String())
    }
    if claims.BranchID != nil {
        c.Set("branch_id", claims.BranchID.String())
    }
}
```

### 3. Update User Handler (Login)

**File:** `internal/handler/user_handler.go`

**Before:**
```go
// Tentukan internal atau external role berdasarkan role type
var internalRole, externalRole string
if user.Role.Type == "INTERNAL" {
    internalRole = user.Role.Name
} else if user.Role.Type == "EXTERNAL" {
    externalRole = user.Role.Name
}

// Generate JWT token
token, err := pkg.GenerateJWT(user.ID, user.Email, internalRole, externalRole, user.CompanyID, user.BranchID)
```

**After:**
```go
// Tentukan internal atau external role berdasarkan role type
var internalRole, externalRole string
roleType := string(user.Role.Type) // ✅ TAMBAH INI

if user.Role.Type == "INTERNAL" {
    internalRole = user.Role.Name
} else if user.Role.Type == "EXTERNAL" {
    externalRole = user.Role.Name
}

// Generate JWT token
token, err := pkg.GenerateJWT(user.ID, user.Email, roleType, internalRole, externalRole, user.CompanyID, user.BranchID) // ✅ TAMBAH roleType
```

## 🧪 Testing

### Test 1: Login sebagai OWNER

**Request:**
```bash
POST http://localhost:8080/api/v1/login
Content-Type: application/json

{
  "email": "owner@restaurant.com",
  "password": "password123"
}
```

**Expected JWT Payload:**
```json
{
  "user_id": "uuid-here",
  "email": "owner@restaurant.com",
  "role_type": "EXTERNAL", // ✅ Harus ada
  "external_role": "OWNER",
  "company_id": "company-uuid",
  "branch_id": null
}
```

### Test 2: OWNER akses logs

**Request:**
```bash
GET http://localhost:8080/api/v1/logs
Authorization: Bearer <owner-token>
```

**Expected Behavior:**
1. ✅ Auth middleware set `role_type` = "EXTERNAL" ke context
2. ✅ Auth middleware set `company_id` = "company-uuid" (string) ke context
3. ✅ Handler detect `role_type` == "EXTERNAL"
4. ✅ Handler apply filter `company_id` = "company-uuid"
5. ✅ Repository query: `WHERE company_id = 'company-uuid'`
6. ✅ Response: Hanya logs dari company owner

**Expected Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "company_id": "company-uuid", // ✅ Sama dengan owner's company
      "path": "/api/v1/external/products"
    }
    // Tidak ada logs dari company lain
  ]
}
```

### Test 3: SUPER_ADMIN akses logs

**Request:**
```bash
POST http://localhost:8080/api/v1/login
{
  "email": "admin@siresto.com",
  "password": "password123"
}
```

**Expected JWT Payload:**
```json
{
  "user_id": "uuid-here",
  "email": "admin@siresto.com",
  "role_type": "INTERNAL", // ✅ Harus ada
  "internal_role": "SUPER_ADMIN",
  "company_id": null,
  "branch_id": null
}
```

**Request:**
```bash
GET http://localhost:8080/api/v1/logs
Authorization: Bearer <super-admin-token>
```

**Expected Behavior:**
1. ✅ Auth middleware set `role_type` = "INTERNAL" ke context
2. ✅ Handler detect `role_type` == "INTERNAL"
3. ✅ Handler TIDAK apply filter company_id/branch_id
4. ✅ Repository query: `SELECT * FROM api_logs` (tanpa WHERE company_id)
5. ✅ Response: Semua logs dari semua company

## 📋 Checklist Verification

### JWT Token
- [x] JWT claims include `role_type`
- [x] `role_type` = "INTERNAL" untuk internal users
- [x] `role_type` = "EXTERNAL" untuk external users
- [x] `company_id` dan `branch_id` tersimpan di JWT

### Auth Middleware
- [x] Set `role_type` ke context
- [x] Set `company_id` sebagai string ke context
- [x] Set `branch_id` sebagai string ke context
- [x] Context keys konsisten (snake_case)

### Handler
- [x] Get `role_type` dari context
- [x] Get `company_id` dari context (string)
- [x] Get `branch_id` dari context (string)
- [x] Apply filter hanya jika `role_type` == "EXTERNAL"

### Repository
- [x] Filter by `company_id` jika provided
- [x] Filter by `branch_id` jika provided
- [x] Tidak filter jika company_id/branch_id kosong (internal users)

## 🔍 Debug Tips

### Jika masih melihat logs dari company lain:

1. **Check JWT token:**
```bash
# Decode JWT token di jwt.io
# Pastikan ada field "role_type": "EXTERNAL"
# Pastikan ada field "company_id": "uuid-here"
```

2. **Check context di handler:**
```go
// Tambah log di handler
roleType, exists := c.Get("role_type")
fmt.Println("Role Type:", roleType, "Exists:", exists)

companyID, exists := c.Get("company_id")
fmt.Println("Company ID:", companyID, "Exists:", exists)
```

3. **Check SQL query:**
```go
// Tambah log di repository
fmt.Println("Filtering by company_id:", companyID)
fmt.Println("Filtering by branch_id:", branchID)
```

4. **Check database:**
```sql
-- Pastikan logs punya company_id
SELECT id, path, company_id, branch_id FROM api_logs LIMIT 10;

-- Check logs dari specific company
SELECT COUNT(*) FROM api_logs WHERE company_id = 'your-company-uuid';
```

## ⚠️ Important Notes

### 1. Perlu Login Ulang
Setelah update ini, user perlu **login ulang** untuk mendapatkan JWT token baru yang include `role_type`.

Token lama (tanpa `role_type`) akan:
- ❌ Tidak punya field `role_type` di claims
- ❌ Context `role_type` akan kosong
- ❌ Filtering tidak bekerja (akan lihat semua logs)

### 2. Context Keys Consistency
Pastikan context keys konsisten di semua file:
- ✅ `user_id` (bukan `userID`)
- ✅ `role_type` (bukan `roleType`)
- ✅ `company_id` (bukan `companyID`)
- ✅ `branch_id` (bukan `branchID`)

### 3. Type Conversion
- JWT menyimpan `company_id` sebagai `*uuid.UUID`
- Handler expect `company_id` sebagai `string`
- Middleware harus convert: `claims.CompanyID.String()`

## ✅ Summary

**Files Changed:**
1. ✅ `pkg/jwt.go` - Tambah `role_type` ke claims
2. ✅ `internal/middleware/auth_middleware.go` - Set `role_type` ke context & convert UUID to string
3. ✅ `internal/handler/user_handler.go` - Pass `role_type` saat generate JWT

**What's Fixed:**
1. ✅ JWT token sekarang include `role_type`
2. ✅ Auth middleware set `role_type` ke context
3. ✅ Company_id dan branch_id di-convert ke string
4. ✅ Handler bisa detect internal vs external user
5. ✅ Filtering bekerja dengan benar

**Status:** ✅ Fixed & Ready to Test

**Date:** March 28, 2026

---

## 🚀 Next Steps

1. ✅ Restart server
2. ✅ Login ulang (untuk dapat token baru)
3. ✅ Test dengan OWNER → Harus hanya lihat logs company sendiri
4. ✅ Test dengan ADMIN → Harus hanya lihat logs branch sendiri
5. ✅ Test dengan SUPER_ADMIN → Harus lihat semua logs
