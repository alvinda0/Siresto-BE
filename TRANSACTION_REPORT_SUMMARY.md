# Transaction Report - Implementation Summary

## Overview

Endpoint untuk mendapatkan laporan transaksi dengan fitur:
- Filter tanggal dan jam
- Search (customer name, phone, order ID)
- Filter status, payment method, order method
- Pagination
- Auto-filtering berdasarkan company dan branch user yang login

## Files Created/Modified

### 1. New Files

#### Entity/DTO
- `internal/entity/transaction_report_dto.go`
  - `TransactionReportDTO` - Response structure untuk report
  - `TransactionReportFilter` - Filter parameters
  - `ToReportDTO()` - Converter method

#### Documentation
- `TRANSACTION_REPORT_API.md` - API documentation lengkap
- `TRANSACTION_REPORT_QUICK_START.md` - Quick start guide
- `TRANSACTION_REPORT_SUMMARY.md` - Implementation summary (this file)
- `TRANSACTION_REPORT_RESPONSE_FORMAT.md` - Response format documentation

#### Testing
- `test_transaction_report.ps1` - PowerShell testing script

### 2. Modified Files

#### Repository Layer
- `internal/repository/order_repository.go`
  - Added `GetTransactionReport()` method to interface
  - Implemented `GetTransactionReport()` with filters:
    - Date range (start_date, end_date)
    - Time range (start_time, end_time)
    - Search (customer name, phone, order ID)
    - Status filters (order status, payment status, payment method, order method)
    - Pagination
    - Auto-filtering by company_id and branch_id

#### Service Layer
- `internal/service/order_service.go`
  - Added `GetTransactionReport()` method to interface
  - Implemented business logic for report generation
  - Convert orders to report DTOs
  - Generate pagination metadata

#### Handler Layer
- `internal/handler/order_handler.go`
  - Added `GetTransactionReport()` handler
  - Parse query parameters
  - Extract company_id and branch_id from context
  - Uses `pkg.SuccessResponseWithMeta()` for success response
  - Uses `pkg.ErrorResponse()` for error response
  - Return standardized JSON response with data and pagination meta

#### Routes
- `routes/routes.go`
  - Added route: `GET /api/v1/external/reports/transactions`
  - Protected by `AuthMiddleware()` and `RequireExternalRole()`

## API Endpoint

```
GET /api/v1/external/reports/transactions
```

### Query Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| start_date | string | Start date (YYYY-MM-DD) | 2024-01-01 |
| end_date | string | End date (YYYY-MM-DD) | 2024-01-31 |
| start_time | string | Start time (HH:MM) | 08:00 |
| end_time | string | End time (HH:MM) | 22:00 |
| search | string | Search keyword | John |
| status | string | Order status | completed |
| payment_status | string | Payment status | paid |
| payment_method | string | Payment method | cash |
| order_method | string | Order method | dine_in |
| page | int | Page number | 1 |
| limit | int | Items per page | 10 |

### Response Structure

```json
{
  "success": true,
  "message": "Transaction report retrieved successfully",
  "status": 200,
  "timestamp": "2024-01-15T10:30:00Z",
  "data": [
    {
      "id": "uuid",
      "order_number": "string",
      "customer_name": "string",
      "customer_phone": "string",
      "table_number": "string",
      "order_method": "string",
      "status": "string",
      "payment_status": "string",
      "payment_method": "string",
      "subtotal_amount": 0,
      "tax_amount": 0,
      "discount_amount": 0,
      "total_amount": 0,
      "paid_amount": 0,
      "change_amount": 0,
      "promo_code": "string",
      "company_name": "string",
      "branch_name": "string",
      "created_at": "timestamp",
      "paid_at": "timestamp"
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

## Features

### 1. Date & Time Filtering
- Filter by date range (start_date, end_date)
- Filter by time range (start_time, end_time)
- Combine date and time for precise filtering
- Support partial date filtering (only start or only end)

### 2. Search Functionality
- Search by customer name (case-insensitive, partial match)
- Search by customer phone (partial match)
- Search by order ID (partial match)

### 3. Status Filtering
- Filter by order status (pending, preparing, ready, completed, cancelled)
- Filter by payment status (unpaid, paid, refunded)
- Filter by payment method (cash, card, qris, etc.)
- Filter by order method (dine_in, takeaway, delivery)

### 4. Pagination
- Configurable page and limit
- Default: page=1, limit=10
- Returns pagination metadata (total_items, total_pages)

### 5. Access Control
- Auto-filtering by company_id and branch_id from JWT token
- User hanya bisa melihat transaksi dari company dan branch mereka
- Protected by AuthMiddleware and RequireExternalRole

### 6. Data Ordering
- Ordered by created_at DESC (newest first)

## Use Cases

### 1. Daily Report
```
GET /api/v1/external/reports/transactions?start_date=2024-01-15&end_date=2024-01-15
```

### 2. Monthly Report
```
GET /api/v1/external/reports/transactions?start_date=2024-01-01&end_date=2024-01-31
```

### 3. Shift Report
```
GET /api/v1/external/reports/transactions?start_date=2024-01-15&end_date=2024-01-15&start_time=08:00&end_time=15:00
```

### 4. Cash Transactions
```
GET /api/v1/external/reports/transactions?payment_method=cash&payment_status=paid
```

### 5. Search Customer
```
GET /api/v1/external/reports/transactions?search=John
```

## Testing

### Using PowerShell Script
```powershell
# Edit token in test_transaction_report.ps1
.\test_transaction_report.ps1
```

### Using cURL
```bash
curl -X GET "http://localhost:8080/api/v1/external/reports/transactions?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Security

1. **Authentication Required**: Bearer token must be provided
2. **Role-Based Access**: Only external users (OWNER, ADMIN, CASHIER) can access
3. **Multi-Tenant Isolation**: Auto-filtering by company_id and branch_id
4. **No Cross-Company Access**: Users cannot see other company's data

## Performance Considerations

1. **Database Indexing**: Ensure indexes on:
   - company_id
   - branch_id
   - created_at
   - status
   - payment_status

2. **Pagination**: Always use pagination to limit result set

3. **Query Optimization**: 
   - Preload only necessary relations (Company, Branch)
   - No OrderItems preload for better performance

## Future Enhancements

1. **Export to Excel/PDF**: Add export functionality
2. **Summary Statistics**: Add total revenue, average order value, etc.
3. **Date Presets**: Add quick filters (today, yesterday, this week, this month)
4. **Caching**: Implement caching for frequently accessed reports
5. **Real-time Updates**: WebSocket support for live report updates

## Documentation

- **API Documentation**: [TRANSACTION_REPORT_API.md](TRANSACTION_REPORT_API.md)
- **Quick Start Guide**: [TRANSACTION_REPORT_QUICK_START.md](TRANSACTION_REPORT_QUICK_START.md)
- **Response Format**: [TRANSACTION_REPORT_RESPONSE_FORMAT.md](TRANSACTION_REPORT_RESPONSE_FORMAT.md)
- **Testing Script**: [test_transaction_report.ps1](test_transaction_report.ps1)

## Checklist

- [x] Create DTO for transaction report
- [x] Implement repository method with filters
- [x] Implement service layer
- [x] Implement handler
- [x] Add route
- [x] Create API documentation
- [x] Create quick start guide
- [x] Create testing script
- [x] Verify no compilation errors
- [x] Document implementation

## Notes

- Report otomatis terfilter berdasarkan company dan branch user yang login
- Semua filter bersifat opsional dan bisa dikombinasikan
- Default sorting: created_at DESC (terbaru dulu)
- Search case-insensitive untuk user experience yang lebih baik
