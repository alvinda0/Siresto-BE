# Dashboard Home - Implementation Summary

## Overview

Endpoint dashboard home untuk menampilkan statistik penjualan restoran dengan data per tanggal (7 hari terakhir) dan item terlaris.

## Endpoint

```
GET /api/v1/external/home
```

## Response Structure

```json
{
  "success": true,
  "data": {
    "total_items_by_date": [
      {"date": "2024-03-24", "value": 145},
      {"date": "2024-03-25", "value": 167}
    ],
    "revenue_by_date": [
      {"date": "2024-03-24", "value": 2150000},
      {"date": "2024-03-25", "value": 2450000}
    ],
    "best_selling_daily": [...],
    "best_selling_weekly": [...],
    "best_selling_monthly": [...],
    "complimentary_items": [...]
  }
}
```

## Features

1. **Total Items by Date** - Total item terjual per tanggal (7 hari terakhir)
2. **Revenue by Date** - Pendapatan per tanggal (7 hari terakhir)
3. **Best Selling Daily** - Top 10 item terlaris hari ini
4. **Best Selling Weekly** - Top 10 item terlaris minggu ini
5. **Best Selling Monthly** - Top 10 item terlaris bulan ini
6. **Complimentary Items** - Item yang dibayar dengan metode COMPLIMENTARY

## Files Created

### Backend Files
- `internal/entity/dashboard_dto.go` - Response DTOs
- `internal/repository/dashboard_repository.go` - Database queries
- `internal/service/dashboard_service.go` - Business logic
- `internal/handler/dashboard_handler.go` - HTTP handler
- `routes/routes.go` - Route registration (updated)

### Documentation Files
- `DASHBOARD_HOME_API.md` - Complete API documentation
- `DASHBOARD_HOME_QUICK_START.md` - Quick start guide
- `DASHBOARD_HOME_SUMMARY.md` - This file
- `test_dashboard_home.ps1` - PowerShell testing script

## Key Implementation Details

### Repository Layer

**GetTotalItemsByDate()**
- Query: Aggregate SUM(quantity) per tanggal dari order_items
- Filter: payment_status = PAID, 7 hari terakhir
- Group by: TO_CHAR(created_at, 'YYYY-MM-DD')

**GetRevenueByDate()**
- Query: Aggregate SUM(total_amount) per tanggal dari orders
- Filter: payment_status = PAID, 7 hari terakhir
- Group by: TO_CHAR(created_at, 'YYYY-MM-DD')

**GetBestSellingItems()**
- Query: Aggregate SUM(quantity) dan SUM(quantity * price) per product
- Filter: payment_status = PAID, periode tertentu
- Order by: total_qty DESC
- Limit: 10 items

**GetComplimentaryItems()**
- Query: Aggregate SUM(quantity) per product
- Filter: payment_method = COMPLIMENTARY, payment_status = PAID
- Order by: total_qty DESC

### Service Layer

Menghitung periode waktu:
- **Daily**: Hari ini (00:00 - 23:59)
- **Weekly**: Minggu ini (Minggu - Sabtu)
- **Monthly**: Bulan ini (tanggal 1 - akhir bulan)

### Handler Layer

- Mengambil company_id dan branch_id dari context (middleware auth)
- Jika branch_id tidak ada, data diagregasi dari semua branch
- Return JSON response dengan format standar

## Database Queries

### Total Items by Date
```sql
SELECT 
  TO_CHAR(orders.created_at, 'YYYY-MM-DD') as date,
  COALESCE(SUM(order_items.quantity), 0) as value
FROM order_items
JOIN orders ON orders.id = order_items.order_id
WHERE orders.company_id = ?
  AND orders.payment_status = 'PAID'
  AND orders.created_at >= CURRENT_DATE - INTERVAL '6 day'
GROUP BY TO_CHAR(orders.created_at, 'YYYY-MM-DD')
ORDER BY date ASC
```

### Revenue by Date
```sql
SELECT 
  TO_CHAR(created_at, 'YYYY-MM-DD') as date,
  COALESCE(SUM(total_amount), 0) as value
FROM orders
WHERE company_id = ?
  AND payment_status = 'PAID'
  AND created_at >= CURRENT_DATE - INTERVAL '6 day'
GROUP BY TO_CHAR(created_at, 'YYYY-MM-DD')
ORDER BY date ASC
```

### Best Selling Items
```sql
SELECT 
  order_items.product_id,
  products.name as product_name,
  SUM(order_items.quantity) as total_qty,
  SUM(order_items.quantity * order_items.price) as total_amount
FROM order_items
JOIN orders ON orders.id = order_items.order_id
JOIN products ON products.id = order_items.product_id
WHERE orders.company_id = ?
  AND orders.payment_status = 'PAID'
  AND orders.created_at BETWEEN ? AND ?
GROUP BY order_items.product_id, products.name
ORDER BY total_qty DESC
LIMIT 10
```

### Complimentary Items
```sql
SELECT 
  order_items.product_id,
  products.name as product_name,
  SUM(order_items.quantity) as total_qty
FROM order_items
JOIN orders ON orders.id = order_items.order_id
JOIN products ON products.id = order_items.product_id
WHERE orders.company_id = ?
  AND orders.payment_method = 'COMPLIMENTARY'
  AND orders.payment_status = 'PAID'
  AND orders.created_at BETWEEN ? AND ?
GROUP BY order_items.product_id, products.name
ORDER BY total_qty DESC
```

## Testing

### Using PowerShell Script
```powershell
# Edit test_dashboard_home.ps1
# Ganti YOUR_JWT_TOKEN_HERE dengan token Anda

.\test_dashboard_home.ps1
```

### Using cURL
```bash
curl -X GET "http://localhost:8080/api/v1/external/home" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Notes

1. Hanya menghitung order dengan `payment_status = PAID`
2. Data per tanggal hanya menampilkan tanggal yang ada transaksi
3. Best selling diurutkan berdasarkan quantity (DESC)
4. Multi-tenant: Data otomatis difilter berdasarkan company_id dan branch_id dari token
5. Complimentary items biasanya untuk item gratis seperti air mineral, kerupuk, dll

## Next Steps

1. Restart server untuk apply changes
2. Test endpoint menggunakan script `test_dashboard_home.ps1`
3. Integrate dengan frontend untuk menampilkan chart/grafik
4. Pertimbangkan caching untuk performa (jika data besar)

## Related Documentation

- **API Documentation**: `DASHBOARD_HOME_API.md`
- **Quick Start Guide**: `DASHBOARD_HOME_QUICK_START.md`
- **Testing Script**: `test_dashboard_home.ps1`
