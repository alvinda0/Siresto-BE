# Transaction Report - Response Format

Endpoint transaction report menggunakan format response standar dari `pkg/response.go` untuk konsistensi dengan endpoint lainnya.

## Standard Response Structure

Semua response mengikuti struktur berikut:

```json
{
  "success": boolean,
  "message": "string",
  "status": number,
  "timestamp": "ISO 8601 timestamp",
  "data": object | array,
  "meta": object,
  "error": "string (only on error)"
}
```

## Success Response

### Structure

```json
{
  "success": true,
  "message": "Transaction report retrieved successfully",
  "status": 200,
  "timestamp": "2024-01-15T10:30:00Z",
  "data": [...],
  "meta": {...}
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Always `true` for successful requests |
| `message` | string | Human-readable success message |
| `status` | number | HTTP status code (200) |
| `timestamp` | string | ISO 8601 timestamp in UTC |
| `data` | array | Array of transaction report objects |
| `meta` | object | Pagination metadata |

### Data Object

Each item in the `data` array contains:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "order_number": "550e8400",
  "customer_name": "John Doe",
  "customer_phone": "081234567890",
  "table_number": "A1",
  "order_method": "DINE_IN",
  "status": "COMPLETED",
  "payment_status": "PAID",
  "payment_method": "TUNAI",
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
```

### Meta Object

Pagination information:

```json
{
  "page": 1,
  "limit": 10,
  "total_items": 50,
  "total_pages": 5
}
```

## Error Response

### Structure

```json
{
  "success": false,
  "message": "Error category",
  "status": 400,
  "timestamp": "2024-01-15T10:30:00Z",
  "error": "Detailed error message"
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Always `false` for error responses |
| `message` | string | Error category or type |
| `status` | number | HTTP status code |
| `timestamp` | string | ISO 8601 timestamp in UTC |
| `error` | string | Detailed error message |

### Common Error Responses

#### 400 Bad Request

```json
{
  "success": false,
  "message": "Bad Request",
  "status": 400,
  "timestamp": "2024-01-15T10:30:00Z",
  "error": "Company ID not found in context"
}
```

#### 401 Unauthorized

```json
{
  "success": false,
  "message": "Unauthorized",
  "status": 401,
  "timestamp": "2024-01-15T10:30:00Z",
  "error": "Invalid or expired token"
}
```

#### 500 Internal Server Error

```json
{
  "success": false,
  "message": "Failed to get transaction report",
  "status": 500,
  "timestamp": "2024-01-15T10:30:00Z",
  "error": "database connection error"
}
```

## Response Handling Examples

### JavaScript/TypeScript

```javascript
// Success handling
fetch('/api/v1/external/reports/transactions', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
})
.then(response => response.json())
.then(result => {
  if (result.success) {
    console.log('Data:', result.data);
    console.log('Pagination:', result.meta);
  } else {
    console.error('Error:', result.error);
  }
});
```

### Go

```go
type TransactionReportResponse struct {
    Success   bool                          `json:"success"`
    Message   string                        `json:"message"`
    Status    int                           `json:"status"`
    Timestamp string                        `json:"timestamp"`
    Data      []entity.TransactionReportDTO `json:"data"`
    Meta      pkg.PaginationMeta            `json:"meta"`
    Error     string                        `json:"error,omitempty"`
}

// Parse response
var response TransactionReportResponse
err := json.Unmarshal(body, &response)
if err != nil {
    return err
}

if response.Success {
    // Handle success
    for _, report := range response.Data {
        fmt.Printf("Order: %s, Total: %.2f\n", report.OrderNumber, report.TotalAmount)
    }
} else {
    // Handle error
    return fmt.Errorf("API error: %s", response.Error)
}
```

### Python

```python
import requests

response = requests.get(
    'http://localhost:8080/api/v1/external/reports/transactions',
    headers={'Authorization': f'Bearer {token}'}
)

result = response.json()

if result['success']:
    for report in result['data']:
        print(f"Order: {report['order_number']}, Total: {report['total_amount']}")
    
    print(f"Page {result['meta']['page']} of {result['meta']['total_pages']}")
else:
    print(f"Error: {result['error']}")
```

## Benefits of Standard Response Format

1. **Consistency**: All endpoints use the same response structure
2. **Easy Parsing**: Client applications can use a single parser
3. **Metadata**: Timestamp and status included in every response
4. **Error Handling**: Clear distinction between success and error
5. **Type Safety**: Strongly typed response structure

## Notes

- All timestamps are in UTC and follow ISO 8601 format
- The `success` field provides quick boolean check
- The `status` field mirrors the HTTP status code
- The `error` field only appears in error responses
- The `data` and `meta` fields only appear in success responses
