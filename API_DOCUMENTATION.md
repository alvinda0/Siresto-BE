# API Documentation - Sistem POS SIRESTO

## Deskripsi
Sistem POS multi-cabang dengan pemisahan role menggunakan UUID.

### 🔐 Internal Role (Platform SIRESTO)
- **SUPER_ADMIN**: Owner sistem SIRESTO
- **SUPPORT**: CS / Admin internal
- **FINANCE**: Lihat pembayaran subscription

### 🍔 External Role (Client Restoran)
- **OWNER**: Pemilik usaha restoran
- **ADMIN**: Manager cabang
- **CASHIER**: Kasir
- **KITCHEN**: Dapur
- **WAITER**: Pelayan

---

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Gunakan JWT token di header untuk endpoint yang dilindungi:
```
Authorization: Bearer <token>
```

---

# PUBLIC ENDPOINTS

## 1. Register (External User - Owner Restoran)
```http
POST /api/v1/register
```

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "company_name": "PT Restoran Sejahtera",
  "company_type": "PT"
}
```

**Field Descriptions:**
- `name` (required): Nama lengkap owner
- `email` (required): Email owner (harus valid dan unique)
- `password` (required): Password minimal 6 karakter
- `company_name` (required): Nama perusahaan
- `company_type` (required): Tipe perusahaan (`PT`, `CV`, atau `PERORANGAN`)

**Response:**
```json
{
  "status": "success",
  "message": "User and company registered successfully",
  "data": {
    "user": {
      "id": "uuid-user",
      "name": "John Doe",
      "email": "john@example.com"
    },
    "company": {
      "id": "uuid-company",
      "name": "PT Restoran Sejahtera",
      "type": "PT"
    }
  }
}
```

---

## 2. Login
```http
POST /api/v1/login
```

**Request Body:**
```json
{
  "email": "owner@restaurant.com",
  "password": "owner123"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

# EXTERNAL API (Client Restoran)
Base: `/api/v1/external`

## User Management

### 3. Create External User
```http
POST /api/v1/external/users
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "Kasir Jakarta",
  "email": "cashier.new@example.com",
  "password": "password123",
  "role_id": "uuid-role-cashier",
  "company_id": "uuid-company",
  "branch_id": "uuid-branch"
}
```

**Field Descriptions:**
- `name` (required): Nama lengkap user
- `email` (required): Email user (harus valid dan unique)
- `password` (required): Password minimal 6 karakter
- `role_id` (required): UUID role dari tabel roles (OWNER, ADMIN, CASHIER, KITCHEN, WAITER)
- `company_id` (optional): UUID company (wajib untuk semua external role)
- `branch_id` (optional): UUID branch (wajib untuk CASHIER, KITCHEN, WAITER)

**Response:**
```json
{
  "status": "success",
  "message": "User created successfully",
  "data": {
    "id": "uuid",
    "name": "Kasir Jakarta",
    "email": "cashier.new@example.com",
    "role_id": "uuid-role",
    "company_id": "uuid-company",
    "branch_id": "uuid-branch"
  }
}
```

---

### 4. Get All External Users
```http
GET /api/v1/external/users
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "External users retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "name": "John Doe",
      "email": "owner@restaurant.com",
      "role": {
        "id": "uuid",
        "name": "OWNER",
        "display_name": "Owner"
      }
    }
  ]
}
```

---

### 5. Get Users by Company
```http
GET /api/v1/external/users/company/:company_id
Authorization: Bearer <token>
```

---

### 6. Get Users by Branch
```http
GET /api/v1/external/users/branch/:branch_id
Authorization: Bearer <token>
```

---

## Company Management

### 7. Create Company
```http
POST /api/v1/external/companies
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "PT Restoran Sejahtera",
  "type": "PT"
}
```

**Type Options:**
- `PT` - Perseroan Terbatas
- `CV` - Commanditaire Vennootschap
- `PERORANGAN` - Usaha Perorangan

**Response:**
```json
{
  "status": "success",
  "message": "Company created successfully",
  "data": {
    "id": "uuid",
    "name": "PT Restoran Sejahtera",
    "type": "PT",
    "owner_id": "uuid"
  }
}
```

---

### 8. Get Company Detail
```http
GET /api/v1/external/companies/detail/:id
Authorization: Bearer <token>
```

---

### 9. Get My Companies
```http
GET /api/v1/external/companies/my
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Companies retrieved successfully",
  "data": [
    {
      "id": "uuid-company-id",
      "name": "PT Restoran Sejahtera",
      "type": "PT",
      "owner_id": "uuid"
    }
  ]
}
```

---

## Branch Management

### 10. Create Branch
```http
POST /api/v1/external/branches
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "company_id": "uuid-company",
  "name": "Cabang Jakarta Pusat",
  "address": "Jl. Sudirman No. 123",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "10220",
  "phone": "021-12345678"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Branch created successfully",
  "data": {
    "id": "uuid-branch-id",
    "company_id": "uuid-company",
    "name": "Cabang Jakarta Pusat",
    "address": "Jl. Sudirman No. 123"
  }
}
```

---

### 11. Get Branch Detail
```http
GET /api/v1/external/branches/detail/:id
Authorization: Bearer <token>
```

---

### 12. Get Branches by Company
```http
GET /api/v1/external/branches/company/:company_id
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "Branches retrieved successfully",
  "data": [
    {
      "id": "uuid-branch-id",
      "company_id": "uuid-company",
      "name": "Cabang Jakarta Pusat",
      "address": "Jl. Sudirman No. 123",
      "city": "Jakarta"
    }
  ]
}
```

---

# INTERNAL API (Platform SIRESTO)
Base: `/api/v1/internal`

## User Management

### 13. Create Internal User
```http
POST /api/v1/internal/users
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "CS Support",
  "email": "support@siresto.com",
  "password": "password123",
  "role_id": "uuid-role-support"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Internal user created successfully",
  "data": {
    "id": "uuid",
    "name": "CS Support",
    "email": "support@siresto.com",
    "role_id": "uuid-role"
  }
}
```

---

### 14. Get All Internal Users
```http
GET /api/v1/internal/users
Authorization: Bearer <token>
```

---

### 15. Get User by ID
```http
GET /api/v1/internal/users/:id
Authorization: Bearer <token>
```

---

### 16. Get All External Users (Monitoring)
```http
GET /api/v1/internal/external-users
Authorization: Bearer <token>
```

---

## Flow Penggunaan

### External (Client Restoran):
1. Register sebagai OWNER
2. Login
3. Buat Company (PT/Perorangan)
4. Buat Branch (cabang restoran)
5. Buat Admin untuk manage cabang
6. Buat Cashier, Kitchen, Waiter untuk operasional

### Internal (Platform SIRESTO):
1. SUPER_ADMIN login
2. Buat user SUPPORT untuk CS
3. Buat user FINANCE untuk monitoring pembayaran
4. Monitor semua external users dan companies

---

## Database Schema

### Users
- id (PK)
- name
- email (unique)
- password (hashed)
- internal_role (SUPER_ADMIN/SUPPORT/FINANCE) - nullable
- external_role (OWNER/ADMIN/CASHIER/KITCHEN/WAITER) - nullable
- company_id (FK -> companies.id) - nullable
- branch_id (FK -> branches.id) - nullable
- is_active
- created_at
- updated_at

**Note:** User hanya bisa punya 1 role (internal ATAU external, tidak bisa keduanya)

### Companies
- id (PK)
- name
- type (PT/PERORANGAN)
- owner_id (FK -> users.id)
- created_at
- updated_at

### Branches
- id (PK)
- company_id (FK -> companies.id)
- name
- address
- city
- province
- postal_code
- phone
- is_active
- created_at
- updated_at

---

## Perbedaan Internal vs External API

| Aspek | Internal API | External API |
|-------|-------------|--------------|
| **Base Path** | `/api/internal` | `/api/external` |
| **User Type** | Platform SIRESTO | Client Restoran |
| **Roles** | SUPER_ADMIN, SUPPORT, FINANCE | OWNER, ADMIN, CASHIER, KITCHEN, WAITER |
| **Scope** | Global (semua companies) | Per company/branch |
| **Purpose** | Monitoring & Management | Operasional restoran |


---

# HELPER ENDPOINTS

### 17. Get All Roles
```http
GET /api/v1/roles
```

**Response:**
```json
{
  "status": "success",
  "message": "Roles retrieved successfully",
  "data": [
    {
      "id": "uuid-role-super-admin",
      "name": "SUPER_ADMIN",
      "display_name": "Super Admin",
      "type": "INTERNAL",
      "description": "Owner sistem SIRESTO"
    },
    {
      "id": "uuid-role-owner",
      "name": "OWNER",
      "display_name": "Owner",
      "type": "EXTERNAL",
      "description": "Pemilik usaha restoran"
    },
    {
      "id": "uuid-role-cashier",
      "name": "CASHIER",
      "display_name": "Cashier",
      "type": "EXTERNAL",
      "description": "Kasir"
    }
  ]
}
```

---

## Flow Penggunaan

### External (Client Restoran):
1. Register sebagai OWNER → `POST /api/v1/register`
2. Login → `POST /api/v1/login`
3. Get Roles untuk mendapatkan role_id → `GET /api/v1/roles`
4. Buat Company → `POST /api/v1/external/companies`
5. Get Company ID → `GET /api/v1/external/companies/my`
6. Buat Branch → `POST /api/v1/external/branches`
7. Get Branch ID → `GET /api/v1/external/branches/company/:company_id`
8. Buat User (Admin/Cashier/Kitchen/Waiter) → `POST /api/v1/external/users`

### Internal (Platform SIRESTO):
1. SUPER_ADMIN login → `POST /api/v1/login`
2. Get Roles → `GET /api/v1/roles`
3. Buat user SUPPORT/FINANCE → `POST /api/v1/internal/users`
4. Monitor companies → `GET /api/v1/internal/companies`
5. Monitor external users → `GET /api/v1/internal/external-users`

---

## Contoh Lengkap: Membuat User Cashier

**Step 1: Login sebagai OWNER**
```bash
POST /api/v1/login
{
  "email": "owner@restaurant.com",
  "password": "owner123"
}
```

**Step 2: Get Roles**
```bash
GET /api/v1/roles
# Cari role dengan name "CASHIER", copy UUID-nya
```

**Step 3: Get Company ID**
```bash
GET /api/v1/external/companies/my
# Copy company ID dari response
```

**Step 4: Get Branch ID**
```bash
GET /api/v1/external/branches/company/{company_id}
# Copy branch ID dari response
```

**Step 5: Create Cashier**
```bash
POST /api/v1/external/users
{
  "name": "Kasir Baru",
  "email": "kasir.baru@restaurant.com",
  "password": "password123",
  "role_id": "uuid-dari-step-2",
  "company_id": "uuid-dari-step-3",
  "branch_id": "uuid-dari-step-4"
}
```

---

## Database Schema

### Roles
- id (PK, UUID)
- name (unique) - SUPER_ADMIN, OWNER, CASHIER, etc
- display_name - Super Admin, Owner, Cashier, etc
- type - INTERNAL or EXTERNAL
- description
- is_active
- created_at
- updated_at

### Users
- id (PK, UUID)
- name
- email (unique)
- password (hashed)
- role_id (FK -> roles.id)
- company_id (FK -> companies.id) - nullable
- branch_id (FK -> branches.id) - nullable
- is_active
- created_at
- updated_at

### Companies
- id (PK, UUID)
- name
- type (PT/PERORANGAN)
- owner_id (FK -> users.id)
- created_at
- updated_at

### Branches
- id (PK, UUID)
- company_id (FK -> companies.id)
- name
- address
- city
- province
- postal_code
- phone
- is_active
- created_at
- updated_at

---

## Test Accounts

Setelah seeder berjalan, gunakan akun berikut:

**Internal Users:**
- superadmin@siresto.com / admin123
- support@siresto.com / support123
- finance@siresto.com / finance123

**External Users:**
- owner@restaurant.com / owner123
- admin@restaurant.com / admin123
- cashier@restaurant.com / cashier123
- kitchen@restaurant.com / kitchen123
- waiter@restaurant.com / waiter123

**Note:** Seeder akan menampilkan Company ID dan Branch ID di console saat pertama kali dijalankan.
