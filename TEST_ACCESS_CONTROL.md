# Testing Access Control - API Logs

Panduan testing untuk memverifikasi access control bekerja dengan benar.

## 🎯 Tujuan Testing

Memastikan:
1. Internal users bisa lihat semua logs
2. OWNER hanya bisa lihat logs company sendiri
3. ADMIN hanya bisa lihat logs company dan branch sendiri
4. External users tidak bisa lihat logs company lain

---

## 📋 Prerequisites

1. Server running di `http://localhost:8080`
2. Database sudah di-migrate dengan kolom `company_id` dan `branch_id`
3. Ada test accounts untuk berbagai role

---

## 🧪 Test Scenarios

### Scenario 1: Internal User (SUPER_ADMIN) - Lihat Semua Logs

**Step 1: Login sebagai SUPER_ADMIN**
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
  "data": {
    "token": "eyJhbGc...",
    "user": {
      "role": {
        "name": "SUPER_ADMIN",
        "type": "INTERNAL"
      }
    }
  }
}
```

**Step 2: Generate logs dari berbagai company**
```bash
# Login sebagai OWNER company A
POST /api/v1/login
{"email": "owner.companyA@test.com", "password": "password123"}

# Buat request (akan generate log dengan company_id A)
GET /api/v1/external/products
Authorization: Bearer <owner-A-token>

# Login sebagai OWNER company B
POST /api/v1/login
{"email": "owner.companyB@test.com", "password": "password123"}

# Buat request (akan generate log dengan company_id B)
GET /api/v1/external/products
Authorization: Bearer <owner-B-token>
```

**Step 3: SUPER_ADMIN lihat semua logs**
```bash
GET http://localhost:8080/api/v1/logs
Authorization: Bearer <super-admin-token>
```

**Expected Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "path": "/api/v1/external/products",
      "company_id": "company-A-uuid",
      "branch_id": "branch-A-uuid"
    },
    {
      "id": 2,
      "path": "/api/v1/external/products",
      "company_id": "company-B-uuid",
      "branch_id": "branch-B-uuid"
    }
  ]
}
```

**Validation:**
- ✅ SUPER_ADMIN dapat melihat logs dari company A
- ✅ SUPER_ADMIN dapat melihat logs dari company B
- ✅ SUPER_ADMIN dapat melihat semua logs

---

### Scenario 2: External User (OWNER) - Hanya Lihat Logs Company Sendiri

**Step 1: Login sebagai OWNER Company A**
```bash
POST http://localhost:8080/api/v1/login
Content-Type: application/json

{
  "email": "owner.companyA@test.com",
  "password": "password123"
}
```

**Expected Response:**
```json
{
  "status": "success",
  "data": {
    "token": "eyJhbGc...",
    "user": {
      "company_id": "company-A-uuid",
      "role": {
        "name": "OWNER",
        "type": "EXTERNAL"
      }
    }
  }
}
```

**Step 2: OWNER A lihat logs**
```bash
GET http://localhost:8080/api/v1/logs
Authorization: Bearer <owner-A-token>
```

**Expected Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "path": "/api/v1/external/products",
      "company_id": "company-A-uuid",
      "branch_id": "branch-A-uuid"
    }
    // Hanya logs dari company A
  ]
}
```

**Validation:**
- ✅ OWNER A hanya melihat logs dari company A
- ✅ OWNER A tidak melihat logs dari company B
- ✅ Filtering otomatis berdasarkan company_id

**Step 3: OWNER A coba akses log milik company B**
```bash
# Ambil ID log dari company B (misal ID 2)
GET http://localhost:8080/api/v1/logs/2
Authorization: Bearer <owner-A-token>
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
- ✅ OWNER A tidak bisa akses log milik company B
- ✅ Return 404 Not Found

---

### Scenario 3: External User (ADMIN) - Hanya Lihat Logs Company & Branch Sendiri

**Step 1: Login sebagai ADMIN Branch A1**
```bash
POST http://localhost:8080/api/v1/login
Content-Type: application/json

{
  "email": "admin.branchA1@test.com",
  "password": "password123"
}
```

**Expected Response:**
```json
{
  "status": "success",
  "data": {
    "token": "eyJhbGc...",
    "user": {
      "company_id": "company-A-uuid",
      "branch_id": "branch-A1-uuid",
      "role": {
        "name": "ADMIN",
        "type": "EXTERNAL"
      }
    }
  }
}
```

**Step 2: Generate logs dari berbagai branch**
```bash
# ADMIN Branch A1 buat request
GET /api/v1/external/products
Authorization: Bearer <admin-A1-token>
# Log: company_id=A, branch_id=A1

# ADMIN Branch A2 buat request
GET /api/v1/external/products
Authorization: Bearer <admin-A2-token>
# Log: company_id=A, branch_id=A2
```

**Step 3: ADMIN A1 lihat logs**
```bash
GET http://localhost:8080/api/v1/logs
Authorization: Bearer <admin-A1-token>
```

**Expected Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 3,
      "path": "/api/v1/external/products",
      "company_id": "company-A-uuid",
      "branch_id": "branch-A1-uuid"
    }
    // Hanya logs dari branch A1
  ]
}
```

**Validation:**
- ✅ ADMIN A1 hanya melihat logs dari branch A1
- ✅ ADMIN A1 tidak melihat logs dari branch A2 (meskipun sama company)
- ✅ Filtering otomatis berdasarkan company_id DAN branch_id

---

### Scenario 4: Filter by Method dengan Access Control

**Test 1: SUPER_ADMIN filter by POST**
```bash
GET http://localhost:8080/api/v1/logs?method=POST
Authorization: Bearer <super-admin-token>
```

**Expected:**
- ✅ Semua POST logs dari semua company

**Test 2: OWNER filter by POST**
```bash
GET http://localhost:8080/api/v1/logs?method=POST
Authorization: Bearer <owner-A-token>
```

**Expected:**
- ✅ Hanya POST logs dari company A

**Test 3: ADMIN filter by DELETE**
```bash
GET http://localhost:8080/api/v1/logs?method=DELETE
Authorization: Bearer <admin-A1-token>
```

**Expected:**
- ✅ Hanya DELETE logs dari company A, branch A1

---

### Scenario 5: Pagination dengan Access Control

**Test 1: OWNER dengan pagination**
```bash
GET http://localhost:8080/api/v1/logs?page=1&limit=5
Authorization: Bearer <owner-A-token>
```

**Expected Response:**
```json
{
  "status": "success",
  "data": [...], // Max 5 items dari company A
  "meta": {
    "page": 1,
    "limit": 5,
    "total_items": 20, // Total logs dari company A saja
    "total_pages": 4
  }
}
```

**Validation:**
- ✅ Pagination hanya menghitung logs dari company sendiri
- ✅ Total items hanya dari company sendiri

---

### Scenario 6: Combine Filter & Pagination

**Test: ADMIN filter POST dengan pagination**
```bash
GET http://localhost:8080/api/v1/logs?method=POST&page=1&limit=10
Authorization: Bearer <admin-A1-token>
```

**Expected:**
- ✅ Hanya POST logs dari company A, branch A1
- ✅ Maksimal 10 items per page
- ✅ Total items hanya dari branch A1

---

## ✅ Validation Checklist

### Internal Users
- [ ] SUPER_ADMIN dapat melihat semua logs
- [ ] SUPPORT dapat melihat semua logs
- [ ] FINANCE dapat melihat semua logs
- [ ] Internal users dapat akses log dari company manapun

### External Users - OWNER
- [ ] OWNER hanya melihat logs dari company sendiri
- [ ] OWNER tidak melihat logs dari company lain
- [ ] OWNER tidak bisa akses log by ID dari company lain
- [ ] Filter by method tetap bekerja dengan access control
- [ ] Pagination menghitung hanya logs company sendiri

### External Users - ADMIN
- [ ] ADMIN hanya melihat logs dari company dan branch sendiri
- [ ] ADMIN tidak melihat logs dari branch lain (meskipun sama company)
- [ ] ADMIN tidak bisa akses log by ID dari branch lain
- [ ] Filter by method tetap bekerja dengan access control
- [ ] Pagination menghitung hanya logs branch sendiri

### Security
- [ ] Company_id dan branch_id tidak bisa di-manipulasi dari request
- [ ] Filtering dilakukan di database level
- [ ] External users tidak bisa bypass access control

---

## 🐛 Troubleshooting

### Issue: OWNER melihat logs dari company lain
**Penyebab:** company_id tidak tersimpan di log atau filtering tidak bekerja
**Solusi:** 
1. Check middleware menyimpan company_id dengan benar
2. Check repository filtering by company_id
3. Check handler mengambil company_id dari context

### Issue: ADMIN melihat logs dari semua branch
**Penyebab:** branch_id tidak difilter
**Solusi:**
1. Check middleware menyimpan branch_id dengan benar
2. Check repository filtering by branch_id
3. Check handler mengambil branch_id dari context

### Issue: Internal users tidak melihat semua logs
**Penyebab:** Filtering diterapkan ke internal users
**Solusi:**
1. Check handler hanya apply filter jika role_type = "EXTERNAL"
2. Check companyID dan branchID kosong untuk internal users

---

## 📊 Test Data Setup

### Companies
- Company A (UUID: company-A-uuid)
  - Branch A1 (UUID: branch-A1-uuid)
  - Branch A2 (UUID: branch-A2-uuid)
- Company B (UUID: company-B-uuid)
  - Branch B1 (UUID: branch-B1-uuid)

### Users
- SUPER_ADMIN: admin@siresto.com
- OWNER Company A: owner.companyA@test.com
- OWNER Company B: owner.companyB@test.com
- ADMIN Branch A1: admin.branchA1@test.com
- ADMIN Branch A2: admin.branchA2@test.com

---

## 🎉 Success Criteria

Testing dianggap berhasil jika:
1. ✅ Internal users dapat melihat semua logs
2. ✅ OWNER hanya melihat logs company sendiri
3. ✅ ADMIN hanya melihat logs company & branch sendiri
4. ✅ External users tidak bisa akses logs company/branch lain
5. ✅ Filter dan pagination bekerja dengan access control
6. ✅ Tidak ada cara untuk bypass access control

---

**Date:** March 28, 2026
**Status:** Ready for Testing
