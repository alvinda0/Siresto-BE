# Dashboard Home - Quick Start Guide

Panduan cepat untuk menggunakan endpoint Dashboard Home.

## Apa yang Ditampilkan?

Endpoint `/api/v1/external/home` menampilkan statistik penjualan:

1. **Total Items by Date** - Total item terjual per tanggal (7 hari terakhir)
2. **Revenue by Date** - Pendapatan per tanggal (7 hari terakhir)
3. **Best Selling Daily** - Top 10 item terlaris hari ini
4. **Best Selling Weekly** - Top 10 item terlaris minggu ini
5. **Best Selling Monthly** - Top 10 item terlaris bulan ini
6. **Complimentary Items** - Item yang diberikan gratis (COMPLIMENTARY)

## Cara Menggunakan

### 1. Login terlebih dahulu

```bash
# Login untuk mendapatkan token
curl -X POST "http://localhost:8080/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@restaurant.com",
    "password": "password123"
  }'
```

Simpan `token` dari response.

### 2. Panggil endpoint home

```bash
curl -X GET "http://localhost:8080/api/v1/external/home" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. Gunakan PowerShell Script (Recommended)

```powershell
# Edit file test_dashboard_home.ps1
# Ganti YOUR_JWT_TOKEN_HERE dengan token Anda

# Jalankan script
.\test_dashboard_home.ps1
```

## Response Example

```json
{
  "success": true,
  "data": {
    "total_items_by_date": [
      {
        "date": "2024-03-24",
        "value": 145
      },
      {
        "date": "2024-03-25",
        "value": 167
      }
    ],
    "revenue_by_date": [
      {
        "date": "2024-03-24",
        "value": 2150000
      },
      {
        "date": "2024-03-25",
        "value": 2450000
      }
    ],
    "best_selling_daily": [
      {
        "product_id": "abc-123",
        "product_name": "Nasi Goreng Spesial",
        "total_qty": 45,
        "total_amount": 1350000
      }
    ],
    "best_selling_weekly": [...],
    "best_selling_monthly": [...],
    "complimentary_items": [
      {
        "product_id": "xyz-789",
        "product_name": "Air Mineral",
        "total_qty": 125
      }
    ]
  }
}
```

## Catatan Penting

1. **Hanya order PAID** - Statistik hanya menghitung order dengan status `payment_status = PAID`
2. **Data per Tanggal** - Menampilkan 7 hari terakhir (termasuk hari ini)
3. **Periode Waktu**:
   - Daily: Hari ini (00:00 - 23:59)
   - Weekly: Minggu ini (Minggu - Sabtu)
   - Monthly: Bulan ini (tanggal 1 - akhir bulan)
4. **Multi-Branch**: Jika user punya akses ke beberapa branch, data akan diagregasi dari semua branch
5. **Best Selling**: Diurutkan berdasarkan quantity terjual (terbanyak ke tersedikit)
6. **Tanggal Kosong**: Jika tidak ada transaksi di tanggal tertentu, tanggal tersebut tidak akan muncul di array

## Troubleshooting

### Error: "Company ID not found"
- Pastikan token valid dan belum expired
- Login ulang untuk mendapatkan token baru

### Data kosong
- Pastikan sudah ada order dengan status PAID
- Cek apakah order sudah dibayar (payment_status = PAID)
- Cek periode waktu (mungkin belum ada transaksi hari/minggu/bulan ini)

### Complimentary items kosong
- Normal jika tidak ada order dengan payment_method = COMPLIMENTARY
- Complimentary biasanya untuk item gratis seperti air mineral, kerupuk, dll

## File Terkait

- **API Documentation**: `DASHBOARD_HOME_API.md`
- **Testing Script**: `test_dashboard_home.ps1`
- **Handler**: `internal/handler/dashboard_handler.go`
- **Service**: `internal/service/dashboard_service.go`
- **Repository**: `internal/repository/dashboard_repository.go`
- **DTO**: `internal/entity/dashboard_dto.go`
