# Transaction Report - Quick Start Guide

Panduan cepat untuk menggunakan Transaction Report API.

## Prerequisites

1. Server sudah running
2. Sudah memiliki token authentication (login sebagai OWNER/ADMIN/CASHIER)
3. Sudah ada data transaksi di database

## Step 1: Login dan Dapatkan Token

```bash
curl -X POST "http://localhost:8080/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@restaurant.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": { ... }
}
```

Simpan token untuk digunakan di request selanjutnya.

## Step 2: Get Transaction Report

### Basic Request (Semua Transaksi)

```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Filter by Date Range

```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Filter Today's Transactions

```bash
# Ganti dengan tanggal hari ini
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?start_date=2024-01-15&end_date=2024-01-15" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Search by Customer

```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?search=John" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Filter by Status

```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?status=completed&payment_status=paid" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Step 3: Using PowerShell Test Script

1. Edit file `test_transaction_report.ps1`
2. Ganti `YOUR_TOKEN_HERE` dengan token yang valid
3. Run script:

```powershell
.\test_transaction_report.ps1
```

## Common Use Cases

### 1. Laporan Harian

```bash
# Transaksi hari ini
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?start_date=2024-01-15&end_date=2024-01-15" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 2. Laporan Bulanan

```bash
# Transaksi bulan Januari 2024
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. Laporan Shift (dengan filter jam)

```bash
# Shift pagi (08:00 - 15:00)
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?start_date=2024-01-15&end_date=2024-01-15&start_time=08:00&end_time=15:00" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Shift malam (15:00 - 22:00)
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?start_date=2024-01-15&end_date=2024-01-15&start_time=15:00&end_time=22:00" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. Laporan Transaksi Cash

```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?payment_method=cash&payment_status=paid" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 5. Laporan Dine-in vs Takeaway

```bash
# Dine-in
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?order_method=dine_in" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Takeaway
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?order_method=takeaway" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Response Format

```json
{
  "success": true,
  "message": "Transaction report retrieved successfully",
  "status": 200,
  "timestamp": "2024-01-15T10:30:00Z",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "order_number": "550e8400",
      "customer_name": "John Doe",
      "customer_phone": "081234567890",
      "table_number": "A1",
      "order_method": "dine_in",
      "status": "completed",
      "payment_status": "paid",
      "payment_method": "cash",
      "subtotal_amount": 100000,
      "tax_amount": 10000,
      "discount_amount": 5000,
      "total_amount": 105000,
      "paid_amount": 110000,
      "change_amount": 5000,
      "promo_code": "PROMO10",
      "company_name": "Restaurant ABC",
      "branch_name": "Cabang Jakarta",
      "created_at": "2024-01-15T10:30:00Z",
      "paid_at": "2024-01-15T11:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 50,
    "total_pages": 5
  }
}
```

## Tips

1. **Pagination**: Gunakan `page` dan `limit` untuk mengatur jumlah data yang ditampilkan
2. **Date Format**: Gunakan format `YYYY-MM-DD` untuk tanggal (contoh: `2024-01-15`)
3. **Time Format**: Gunakan format `HH:MM` untuk jam (contoh: `08:00`, `22:30`)
4. **Search**: Search akan mencari di customer name, phone, dan order ID
5. **Combine Filters**: Semua filter bisa dikombinasikan sesuai kebutuhan
6. **Auto Filtering**: Report otomatis terfilter berdasarkan company dan branch user yang login

## Troubleshooting

### Error: "Company ID not found in context"
- Pastikan token valid dan belum expired
- Login ulang untuk mendapatkan token baru

### Error: "Unauthorized"
- Pastikan menggunakan Bearer token di header Authorization
- Format: `Authorization: Bearer YOUR_TOKEN`

### Empty Data
- Pastikan sudah ada transaksi di database
- Cek filter yang digunakan, mungkin terlalu spesifik
- Pastikan date range mencakup periode yang ada transaksinya

## Next Steps

- Lihat [TRANSACTION_REPORT_API.md](TRANSACTION_REPORT_API.md) untuk dokumentasi lengkap
- Gunakan script [test_transaction_report.ps1](test_transaction_report.ps1) untuk testing
