# Fix: UUID Migration Error

## 🐛 Error

```json
{
  "error": "sql: Scan error on column index 0, name \"id\": Scan: unable to scan type int64 into UUID"
}
```

## 🔍 Root Cause

Tabel `api_logs` di database masih menggunakan tipe `SERIAL` (int64), tapi kode Go sudah expect `UUID`.

**Mismatch:**
- Database: `id SERIAL` (int64)
- Go Code: `id uuid.UUID`

## ✅ Solution

Perlu drop dan recreate tabel `api_logs` dengan schema UUID yang baru.

---

## 🔧 Fix Steps

### Option 1: Manual SQL (Recommended)

**Step 1: Connect to database**
```bash
psql -U your_username -d your_database
```

**Step 2: Run migration script**
```sql
-- Drop old table
DROP TABLE IF EXISTS api_logs;

-- Create new table with UUID
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

-- Create indexes
CREATE INDEX idx_api_logs_method ON api_logs(method);
CREATE INDEX idx_api_logs_path ON api_logs(path);
CREATE INDEX idx_api_logs_user_id ON api_logs(user_id);
CREATE INDEX idx_api_logs_company_id ON api_logs(company_id);
CREATE INDEX idx_api_logs_branch_id ON api_logs(branch_id);
CREATE INDEX idx_api_logs_access_from ON api_logs(access_from);
CREATE INDEX idx_api_logs_created_at ON api_logs(created_at);
CREATE INDEX idx_api_logs_deleted_at ON api_logs(deleted_at);
```

**Step 3: Verify**
```sql
\d api_logs
```

Expected output:
```
Column        | Type      | Nullable | Default
--------------+-----------+----------+-------------------------
id            | uuid      | not null | gen_random_uuid()
method        | varchar   | not null |
...
```

### Option 2: Using Migration File

**Step 1: Run migration file**
```bash
psql -U your_username -d your_database -f migrate_api_logs_to_uuid.sql
```

### Option 3: Restart Server (Auto Migration)

Jika `config.go` sudah diupdate dengan schema UUID:

**Step 1: Drop table manually**
```sql
DROP TABLE IF EXISTS api_logs;
```

**Step 2: Restart server**
```bash
go run cmd/server/main.go
```

Server akan otomatis create table dengan schema UUID yang baru.

---

## 🧪 Verification

### Test 1: Check Table Schema

**SQL:**
```sql
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'api_logs' 
ORDER BY ordinal_position;
```

**Expected:**
```
column_name    | data_type
---------------+-----------
id             | uuid       ✅
method         | character varying
path           | character varying
status_code    | integer
response_time  | bigint
ip_address     | character varying
user_agent     | text
access_from    | character varying
user_id        | uuid       ✅
company_id     | uuid       ✅
branch_id      | uuid       ✅
...
```

### Test 2: Test API

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
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",  ✅ UUID format
      "method": "GET",
      "path": "/api/v1/products"
    }
  ]
}
```

### Test 3: Create New Log

**Request:**
```bash
GET http://localhost:8080/api/v1/products
Authorization: Bearer <token>
```

**Then check logs:**
```bash
GET http://localhost:8080/api/v1/logs
Authorization: Bearer <token>
```

**Validation:**
- ✅ New log created with UUID
- ✅ No scan error
- ✅ Response successful

---

## ⚠️ Important Notes

### 1. Data Loss
**WARNING:** Dropping table akan menghapus semua existing logs!

Jika perlu backup:
```sql
-- Backup old data
CREATE TABLE api_logs_backup AS SELECT * FROM api_logs;

-- After migration, if needed, you can reference old data
SELECT * FROM api_logs_backup;
```

### 2. Production Consideration

Untuk production, pertimbangkan:
- Backup data terlebih dahulu
- Lakukan migration saat traffic rendah
- Monitor error setelah migration

### 3. Alternative: Keep Old Data

Jika ingin keep old data (not recommended):
```sql
-- Rename old table
ALTER TABLE api_logs RENAME TO api_logs_old;

-- Create new table
CREATE TABLE api_logs (...);

-- Old data tetap ada di api_logs_old
-- Tapi tidak bisa di-query via API (karena schema berbeda)
```

---

## 🔍 Troubleshooting

### Issue: Still getting scan error after migration

**Check:**
1. Pastikan table sudah di-drop dan recreate
2. Restart server setelah migration
3. Clear any connection pool

**Verify:**
```sql
-- Check table exists
SELECT EXISTS (
    SELECT FROM information_schema.tables 
    WHERE table_name = 'api_logs'
);

-- Check id column type
SELECT data_type 
FROM information_schema.columns 
WHERE table_name = 'api_logs' AND column_name = 'id';
-- Should return: uuid
```

### Issue: Table not found after drop

**Solution:**
Restart server, migration akan auto-create table dengan schema baru.

### Issue: Permission denied to drop table

**Solution:**
```sql
-- Grant permission
GRANT ALL PRIVILEGES ON TABLE api_logs TO your_username;

-- Then drop
DROP TABLE api_logs;
```

---

## 📋 Quick Fix Checklist

- [ ] Backup existing logs (if needed)
- [ ] Connect to database
- [ ] Run: `DROP TABLE IF EXISTS api_logs;`
- [ ] Run migration script or restart server
- [ ] Verify table schema (id should be UUID)
- [ ] Test API endpoint
- [ ] Verify new logs are created with UUID
- [ ] Monitor for errors

---

## 🎉 Summary

**Problem:** Database table masih integer, code expect UUID

**Solution:** Drop dan recreate table dengan UUID schema

**Steps:**
1. ✅ Drop table: `DROP TABLE IF EXISTS api_logs;`
2. ✅ Restart server (auto create dengan UUID)
3. ✅ Test API

**Status:** Ready to Fix

**Date:** March 28, 2026

---

## 🚀 Quick Command

```bash
# One-liner to fix (adjust connection details)
psql -U postgres -d siresto_db -c "DROP TABLE IF EXISTS api_logs;"

# Then restart server
go run cmd/server/main.go
```
