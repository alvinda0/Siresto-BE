# Dashboard Home API

Endpoint untuk menampilkan statistik di halaman home/dashboard.

## Endpoint

```
GET /api/v1/external/home
```

## Authentication

Memerlukan token JWT di header:
```
Authorization: Bearer <token>
```

## Response

### Success Response (200 OK)

```json
{
  "success": true,
  "data": {
    "total_items": 1250,
    "weekly_revenue": 15750000,
    "best_selling_daily": [
      {
        "product_id": "uuid",
        "product_name": "Nasi Goreng Spesial",
        "total_qty": 45,
        "total_amount": 1350000
      },
      {
        "product_id": "uuid",
        "product_name": "Es Teh Manis",
        "total_qty": 38,
        "total_amount": 190000
      }
    ],
    "best_selling_weekly": [
      {
        "product_id": "uuid",
        "product_name": "Nasi Goreng Spesial",
        "total_qty": 285,
        "total_amount": 8550000
      },
      {
        "product_id": "uuid",
        "product_name": "Ayam Bakar",
        "total_qty": 210,
        "total_amount": 7350000
      }
    ],
    "best_selling_monthly": [
      {
        "product_id": "uuid",
        "product_name": "Nasi Goreng Spesial",
        "total_qty": 1150,
        "total_amount": 34500000
      },
      {
        "product_id": "uuid",
        "product_name": "Ayam Bakar",
        "total_qty": 890,
        "total_amount": 31150000
      }
    ],
    "complimentary_items": [
      {
        "product_id": "uuid",
        "product_name": "Air Mineral",
        "total_qty": 125
      },
      {
        "product_id": "uuid",
        "product_name": "Kerupuk",
        "total_qty": 98
      }
    ]
  }
}
```

## Field Descriptions

### Response Data

| Field | Type | Description |
|-------|------|-------------|
| `total_items_by_date` | array | Total item terjual per tanggal (7 hari terakhir) |
| `revenue_by_date` | array | Total pendapatan per tanggal (7 hari terakhir) |
| `best_selling_daily` | array | Top 10 item terlaris hari ini |
| `best_selling_weekly` | array | Top 10 item terlaris minggu ini |
| `best_selling_monthly` | array | Top 10 item terlaris bulan ini |
| `complimentary_items` | array | Item yang dibayar dengan metode COMPLIMENTARY bulan ini |

### Daily Stats

| Field | Type | Description |
|-------|------|-------------|
| `date` | string | Tanggal (format: YYYY-MM-DD) |
| `value` | float | Nilai (total items atau revenue) |

### Best Selling Item

| Field | Type | Description |
|-------|------|-------------|
| `product_id` | string (uuid) | ID produk |
| `product_name` | string | Nama produk |
| `total_qty` | integer | Total quantity terjual |
| `total_amount` | float | Total nilai penjualan |

### Complimentary Item

| Field | Type | Description |
|-------|------|-------------|
| `product_id` | string (uuid) | ID produk |
| `product_name` | string | Nama produk |
| `total_qty` | integer | Total quantity yang diberikan complimentary |

## Notes

1. Endpoint ini hanya menghitung order dengan `payment_status = PAID`
2. Data best selling diurutkan berdasarkan quantity terjual (DESC)
3. Data per tanggal menampilkan 7 hari terakhir (termasuk hari ini)
4. Periode waktu:
   - Daily: Hari ini (00:00 - 23:59)
   - Weekly: Minggu ini (Minggu - Sabtu)
   - Monthly: Bulan ini (tanggal 1 - akhir bulan)
5. Jika user memiliki akses ke multiple branches, data akan diagregasi dari semua branch
6. Jika user hanya memiliki akses ke 1 branch, data hanya dari branch tersebut
7. Data per tanggal hanya menampilkan tanggal yang ada transaksi (tidak menampilkan tanggal tanpa transaksi)

## Error Responses

### 401 Unauthorized
```json
{
  "error": "Company ID not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "error message"
}
```

## Example Usage

### cURL
```bash
curl -X GET "http://localhost:8080/api/v1/external/home" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### PowerShell
```powershell
$token = "YOUR_JWT_TOKEN"
$headers = @{
    "Authorization" = "Bearer $token"
}

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/external/home" `
    -Method Get `
    -Headers $headers
```
