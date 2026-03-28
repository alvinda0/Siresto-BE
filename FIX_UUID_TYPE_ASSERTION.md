# Fix UUID Type Assertion

## 🐛 Masalah

Error panic saat akses endpoint products:
```
interface conversion: interface {} is uuid.UUID, not *uuid.UUID
```

## 🔍 Penyebab

Inkonsistensi tipe data `company_id` dan `branch_id` di context:

**Sebelum:**
- `AuthMiddleware()` set sebagai **string**: `c.Set("company_id", claims.CompanyID.String())`
- Handler expect **pointer UUID**: `companyID := *(companyIDVal.(*uuid.UUID))`

**Konflik:**
- Order handler butuh `uuid.UUID` (value)
- Product handler butuh `*uuid.UUID` (pointer)

## ✅ Solusi

Standardisasi semua middleware dan handler menggunakan **uuid.UUID** (value, bukan pointer):

### 1. Fix AuthMiddleware
```go
// internal/middleware/auth_middleware.go
if claims.CompanyID != nil {
    c.Set("company_id", *claims.CompanyID)  // Set as uuid.UUID value
}
if claims.BranchID != nil {
    c.Set("branch_id", *claims.BranchID)    // Set as uuid.UUID value
}
```

### 2. Fix WebSocketAuthMiddleware
```go
// internal/middleware/websocket_auth_middleware.go
if claims.CompanyID != nil {
    c.Set("company_id", *claims.CompanyID)  // Consistent with AuthMiddleware
}
if claims.BranchID != nil {
    c.Set("branch_id", *claims.BranchID)
}
```

### 3. Fix Product Handler
```go
// internal/handler/product_handler.go
// BEFORE:
companyID := *(companyIDVal.(*uuid.UUID))
branchID := *(branchIDVal.(*uuid.UUID))

// AFTER:
companyID := companyIDVal.(uuid.UUID)
branchID := branchIDVal.(uuid.UUID)
```

## 📝 File yang Diubah

1. `internal/middleware/auth_middleware.go` - Set UUID as value
2. `internal/middleware/websocket_auth_middleware.go` - Set UUID as value
3. `internal/handler/product_handler.go` - Type assertion tanpa pointer

## ✅ Hasil

- ✅ Products endpoint works
- ✅ Orders endpoint works
- ✅ WebSocket connection works
- ✅ Semua handler konsisten menggunakan `uuid.UUID` value

## 🧪 Testing

```bash
# Compile
go build -o server.exe ./cmd/server

# Test products
GET http://localhost:8080/api/v1/external/products
Authorization: Bearer YOUR_TOKEN

# Test orders
GET http://localhost:8080/api/v1/orders
Authorization: Bearer YOUR_TOKEN

# Test WebSocket
ws://localhost:8080/api/v1/ws/orders?token=YOUR_TOKEN
```

## 📚 Best Practice

**Gunakan value, bukan pointer untuk UUID di context:**

✅ **GOOD:**
```go
c.Set("company_id", uuid.UUID{...})
companyID := c.Get("company_id").(uuid.UUID)
```

❌ **BAD:**
```go
c.Set("company_id", &uuid.UUID{...})
companyID := *(c.Get("company_id").(*uuid.UUID))
```

**Alasan:**
- Lebih simple
- Tidak perlu dereference
- Konsisten dengan order handler
- Menghindari nil pointer panic
