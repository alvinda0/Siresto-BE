# Tax API - Multi-Tenant Update

Tax API telah diupdate untuk mendukung multi-tenant berdasarkan company_id dan branch_id dari user yang login.

## 🔄 Perubahan

### 1. Database Schema
Tabel `taxes` sekarang memiliki kolom tambahan:
- `company_id` (uuid, NOT NULL) - Foreign key ke companies
- `branch_id` (uuid, nullable) - Foreign key ke branches

```sql
CREATE TABLE taxes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id uuid NOT NULL,
    branch_id uuid,  -- NULL = company level, filled = branch specific
    nama_pajak varchar(100) NOT NULL,
    tipe_pajak varchar(10) NOT NULL,
    presentase decimal(5,2) NOT NULL,
    deskripsi text,
    status varchar(20) DEFAULT 'active',
    prioritas integer DEFAULT 0,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT fk_taxes_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT fk_taxes_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE
);
```

### 2. Access Control
Semua operasi CRUD sekarang:
- ✅ Otomatis menggunakan `company_id` dari user yang login
- ✅ Otomatis menggunakan `branch_id` dari user yang login (jika ada)
- ✅ User hanya bisa akses tax milik company/branch mereka sendiri
- ✅ Tidak bisa akses tax milik company/branch lain

### 3. Tax Levels

#### Company-Level Tax
Tax yang berlaku untuk seluruh company (semua branch):
- `branch_id` = NULL
- Dibuat oleh OWNER atau ADMIN tanpa specify branch_id
- Akan muncul di semua branch dalam company tersebut

#### Branch-Level Tax
Tax yang hanya berlaku untuk branch tertentu:
- `branch_id` = UUID branch tertentu
- Dibuat dengan specify `branch_id` di request body
- Hanya muncul di branch tersebut

---

## 📋 API Changes

### CREATE Tax
**POST** `/api/v1/external/tax`

**Body (Updated):**
```json
{
  "branch_id": "uuid-optional",  // NEW: optional, null = company level
  "nama_pajak": "PB1",
  "tipe_pajak": "pb1",
  "presentase": 10.00,
  "deskripsi": "Pajak Barang dan Jasa 1",
  "status": "active",
  "prioritas": 1
}
```

**Behavior:**
- Jika `branch_id` tidak dikirim atau null → Tax level company
- Jika `branch_id` dikirim → Tax level branch tersebut
- `company_id` otomatis dari user yang login

**Response:**
```json
{
  "status": "success",
  "message": "Tax created successfully",
  "data": {
    "id": "uuid",
    "company_id": "uuid",      // NEW
    "branch_id": "uuid-or-null", // NEW
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10.00,
    "deskripsi": "Pajak Barang dan Jasa 1",
    "status": "active",
    "prioritas": 1,
    "created_at": "2024-01-01 10:00:00",
    "updated_at": "2024-01-01 10:00:00"
  }
}
```

### GET All Taxes
**GET** `/api/v1/external/tax`

**Behavior:**
- Jika user punya `branch_id` → Return company-level + branch-specific taxes
- Jika user tidak punya `branch_id` (OWNER) → Return hanya company-level taxes

**Example:**
User ADMIN di Branch A akan melihat:
1. Tax company-level (branch_id = null)
2. Tax khusus Branch A (branch_id = Branch A UUID)

User OWNER akan melihat:
1. Hanya tax company-level (branch_id = null)

### UPDATE, DELETE, GET by ID
Semua operasi ini sekarang:
- ✅ Memvalidasi `company_id` sesuai user
- ✅ Memvalidasi `branch_id` sesuai user (jika ada)
- ✅ Return 404 jika tax bukan milik company/branch user

---

## 🧪 Testing Examples

### 1. Create Company-Level Tax (OWNER)
```bash
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $OWNER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "PB1 Company Wide",
    "tipe_pajak": "pb1",
    "presentase": 10.00,
    "status": "active",
    "prioritas": 1
  }'
```
Result: Tax berlaku untuk semua branch

### 2. Create Branch-Level Tax (ADMIN)
```bash
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "branch_id": "branch-uuid-here",
    "nama_pajak": "Service Charge Branch A",
    "tipe_pajak": "sc",
    "presentase": 5.00,
    "status": "active",
    "prioritas": 2
  }'
```
Result: Tax hanya berlaku untuk branch tersebut

### 3. Get All Taxes (ADMIN di Branch A)
```bash
curl -X GET http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```
Result: Akan return:
- Tax company-level (berlaku untuk semua branch)
- Tax khusus Branch A

### 4. Get All Taxes (OWNER)
```bash
curl -X GET http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $OWNER_TOKEN"
```
Result: Hanya return tax company-level

---

## 🔒 Security

### Access Control Rules:
1. ✅ User hanya bisa CRUD tax milik company mereka
2. ✅ User dengan branch_id hanya bisa akses:
   - Tax company-level (branch_id = null)
   - Tax branch mereka sendiri
3. ✅ User tanpa branch_id (OWNER) hanya bisa akses tax company-level
4. ✅ Tidak bisa akses tax company/branch lain (return 404)

### Foreign Key Constraints:
- `company_id` → CASCADE DELETE (hapus company = hapus semua tax-nya)
- `branch_id` → CASCADE DELETE (hapus branch = hapus tax khusus branch tersebut)

---

## 📊 Use Cases

### Use Case 1: Tax Nasional (Company-Level)
PB1 10% berlaku untuk semua branch:
```json
{
  "branch_id": null,
  "nama_pajak": "PB1",
  "tipe_pajak": "pb1",
  "presentase": 10.00
}
```

### Use Case 2: Service Charge per Branch
Branch A: 5%, Branch B: 7%:

**Branch A:**
```json
{
  "branch_id": "branch-a-uuid",
  "nama_pajak": "Service Charge",
  "tipe_pajak": "sc",
  "presentase": 5.00
}
```

**Branch B:**
```json
{
  "branch_id": "branch-b-uuid",
  "nama_pajak": "Service Charge",
  "tipe_pajak": "sc",
  "presentase": 7.00
}
```

### Use Case 3: Tax Khusus Daerah
Branch di Jakarta punya tax tambahan:
```json
{
  "branch_id": "jakarta-branch-uuid",
  "nama_pajak": "Pajak Daerah Jakarta",
  "tipe_pajak": "pb1",
  "presentase": 2.00
}
```

---

## 🔄 Migration

Jika tabel `taxes` sudah ada, migration akan otomatis:
1. Tambah kolom `company_id` dan `branch_id`
2. Tambah foreign key constraints
3. Tambah indexes

**Note:** Data lama yang tidak punya `company_id` perlu diupdate manual atau dihapus.

---

## ✅ Benefits

1. **Multi-Tenant Isolation**: Setiap company punya tax sendiri
2. **Flexibility**: Tax bisa company-level atau branch-level
3. **Security**: User tidak bisa akses tax company lain
4. **Scalability**: Mendukung banyak company dan branch
5. **Consistency**: Sama seperti Category dan Product API

---

## 📝 Notes

- Tax company-level (branch_id = null) akan muncul di semua branch
- Tax branch-level hanya muncul di branch tersebut
- User bisa create tax company-level atau branch-level sesuai kebutuhan
- Sorting tetap: prioritas DESC, nama_pajak ASC
