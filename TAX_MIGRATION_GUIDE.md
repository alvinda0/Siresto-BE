# Tax Migration Guide

Panduan untuk membuat tabel `taxes` di database.

## Pilihan Migration

Ada 3 cara untuk membuat tabel taxes:

### 1. Auto Migration (Recommended)

Tabel akan otomatis dibuat saat server start karena sudah ditambahkan di `config/config.go`.

```bash
# Jalankan server
go run cmd/server/main.go
```

Server akan otomatis:
- Create tabel `taxes` jika belum ada
- Create indexes
- Siap digunakan

---

### 2. Manual Migration dengan Go Script

Jalankan script migration standalone:

```bash
# Run migration
go run run_tax_migration.go
```

Script ini akan:
- Create tabel `taxes`
- Create indexes
- Insert sample data (PB1 dan Service Charge)
- Verify hasil migration

**Output yang diharapkan:**
```
Connected to database
Creating taxes table...
✓ Taxes table created
Creating indexes...
✓ Indexes created
Inserting sample data...
✓ Inserted: PB1
✓ Inserted: Service Charge
Total taxes in database: 2

Current taxes:
- Service Charge (sc): 5.00% [active] Priority: 2
- PB1 (pb1): 10.00% [active] Priority: 1

✅ Migration completed successfully!
```

---

### 3. Manual Migration dengan SQL

Jalankan SQL script langsung ke database:

```bash
# Menggunakan psql
psql $DATABASE_URL -f migrate_taxes.sql

# Atau copy-paste isi file migrate_taxes.sql ke database client Anda
```

---

## Verifikasi Migration

Setelah migration, verifikasi dengan query:

```sql
-- Check table exists
SELECT table_name 
FROM information_schema.tables 
WHERE table_name = 'taxes';

-- Check columns
SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns
WHERE table_name = 'taxes'
ORDER BY ordinal_position;

-- Check indexes
SELECT indexname, indexdef
FROM pg_indexes
WHERE tablename = 'taxes';

-- Check data
SELECT * FROM taxes ORDER BY prioritas DESC, nama_pajak ASC;
```

---

## Schema Details

```sql
CREATE TABLE taxes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_pajak varchar(100) NOT NULL,
    tipe_pajak varchar(10) NOT NULL CHECK (tipe_pajak IN ('sc', 'pb1')),
    presentase decimal(5,2) NOT NULL CHECK (presentase >= 0 AND presentase <= 100),
    deskripsi text,
    status varchar(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    prioritas integer DEFAULT 0,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);
```

**Indexes:**
- `idx_taxes_status` on `status`
- `idx_taxes_prioritas` on `prioritas`
- `idx_taxes_tipe_pajak` on `tipe_pajak`

**Constraints:**
- `tipe_pajak` must be 'sc' or 'pb1'
- `presentase` must be between 0 and 100
- `status` must be 'active' or 'inactive'

---

## Sample Data

Migration akan insert 2 sample data:

| nama_pajak | tipe_pajak | presentase | status | prioritas |
|------------|------------|------------|--------|-----------|
| PB1 | pb1 | 10.00 | active | 1 |
| Service Charge | sc | 5.00 | active | 2 |

---

## Rollback

Jika perlu rollback (hapus tabel):

```sql
-- Drop table
DROP TABLE IF EXISTS taxes CASCADE;
```

**Warning:** Ini akan menghapus semua data taxes!

---

## Troubleshooting

### Error: relation "taxes" already exists
Tabel sudah ada, tidak perlu migration lagi.

### Error: permission denied
Pastikan user database memiliki permission untuk CREATE TABLE.

### Error: database connection failed
Check `DATABASE_URL` di file `.env`:
```
DATABASE_URL=postgresql://user:password@localhost:5432/dbname
```

---

## Next Steps

Setelah migration berhasil:

1. ✅ Tabel `taxes` sudah siap
2. ✅ Indexes sudah dibuat
3. ✅ Sample data sudah ada (optional)
4. 🚀 Jalankan server: `go run cmd/server/main.go`
5. 🧪 Test API: `./test_tax.sh "your_token"`

Lihat `TAX_README.md` untuk quick start guide.
