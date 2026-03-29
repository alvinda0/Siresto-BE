# Tax API - Quick Start

CRUD API untuk mengelola data pajak (Tax) dalam sistem SIRESTO.

## Struktur Data

```go
type Tax struct {
    ID          uuid.UUID  // Primary key
    NamaPajak   string     // Nama pajak (contoh: "PB1", "Service Charge")
    TipePajak   string     // Tipe: "sc" atau "pb1"
    Presentase  float64    // Persentase pajak (0-100)
    Deskripsi   string     // Deskripsi pajak
    Status      string     // "active" atau "inactive"
    Prioritas   int        // Urutan prioritas (default: 0)
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

## Endpoints

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/api/v1/external/tax` | Create tax |
| PUT | `/api/v1/external/tax/:id` | Update tax |
| GET | `/api/v1/external/tax/:id` | Get tax by ID |
| GET | `/api/v1/external/tax` | Get all taxes |
| DELETE | `/api/v1/external/tax/:id` | Delete tax |

## Quick Test

```bash
# Set token
export TOKEN="your_jwt_token"

# Create tax
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10.00,
    "status": "active",
    "prioritas": 1
  }'

# Get all taxes
curl -X GET http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN"
```

## Testing Scripts

### Bash
```bash
chmod +x test_tax.sh
./test_tax.sh "your_token_here"
```

### PowerShell
```powershell
.\test_tax.ps1 -Token "your_token_here"
```

## Files Created

```
internal/
├── entity/
│   ├── tax.go           # Tax entity model
│   └── tax_dto.go       # Request/Response DTOs
├── repository/
│   └── tax_repository.go # Database operations
├── service/
│   └── tax_service.go    # Business logic
└── handler/
    └── tax_handler.go    # HTTP handlers

routes/routes.go          # Updated with tax routes
config/config.go          # Updated with tax migration

TAX_API.md               # Full API documentation
TAX_TESTING.md           # Testing guide
test_tax.sh              # Bash testing script
test_tax.ps1             # PowerShell testing script
```

## Tipe Pajak

- `sc` - Service Charge (biaya layanan)
- `pb1` - Pajak Barang dan Jasa 1

## Status

- `active` - Pajak aktif, akan diterapkan
- `inactive` - Pajak tidak aktif

## Validasi

- `nama_pajak`: required
- `tipe_pajak`: required, enum (sc, pb1)
- `presentase`: required, 0-100
- `status`: optional, enum (active, inactive), default: "active"
- `prioritas`: optional, integer, default: 0

## Database Migration

Ada 3 cara untuk membuat tabel `taxes`:

### 1. Auto Migration (Recommended)
```bash
go run cmd/server/main.go
```
Tabel otomatis dibuat saat server start.

### 2. Manual dengan Go Script
```bash
go run run_tax_migration.go
```
Akan create tabel + insert sample data.

### 3. Manual dengan SQL
```bash
psql $DATABASE_URL -f migrate_taxes.sql
```

Lihat `TAX_MIGRATION_GUIDE.md` untuk detail lengkap.

## Next Steps

1. Start server: `go run cmd/server/main.go`
2. Login untuk mendapatkan token
3. Test endpoints menggunakan script atau manual curl
4. Lihat `TAX_API.md` untuk dokumentasi lengkap
5. Lihat `TAX_TESTING.md` untuk panduan testing detail
