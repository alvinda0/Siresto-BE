# Tax API - Standard Response Format

Tax API sekarang menggunakan standard response format yang konsisten dengan handler lain menggunakan `pkg/response.go`.

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation message",
  "status": 200,
  "timestamp": "2026-03-29T04:00:00Z",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "status": 400,
  "timestamp": "2026-03-29T04:00:00Z",
  "error": "Detailed error description"
}
```

---

## Examples

### 1. CREATE Tax - Success (201)
```json
{
  "success": true,
  "message": "Tax created successfully",
  "status": 201,
  "timestamp": "2026-03-29T04:57:23Z",
  "data": {
    "id": "367a2292-a422-4dad-9e1f-79732a7adce9",
    "company_id": "2fa830c2-daec-4ddb-b061-2f7c50b7562b",
    "branch_id": "8fe6b4f9-3dc9-4773-82a8-9f92f0d86458",
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10.00,
    "deskripsi": "Pajak Barang dan Jasa 1",
    "status": "active",
    "prioritas": 1,
    "created_at": "2026-03-29 11:57:23",
    "updated_at": "2026-03-29 11:57:23"
  }
}
```

### 2. GET All Taxes - Success (200)
```json
{
  "success": true,
  "message": "Taxes retrieved successfully",
  "status": 200,
  "timestamp": "2026-03-29T04:57:30Z",
  "data": [
    {
      "id": "uuid-1",
      "company_id": "uuid",
      "branch_id": null,
      "nama_pajak": "PB1",
      "tipe_pajak": "pb1",
      "presentase": 10.00,
      "deskripsi": "Pajak Barang dan Jasa 1 (Company Level)",
      "status": "active",
      "prioritas": 1,
      "created_at": "2026-03-29 11:54:28",
      "updated_at": "2026-03-29 11:54:28"
    },
    {
      "id": "uuid-2",
      "company_id": "uuid",
      "branch_id": "uuid",
      "nama_pajak": "Service Charge",
      "tipe_pajak": "sc",
      "presentase": 5.00,
      "deskripsi": "Biaya layanan (Branch Level)",
      "status": "active",
      "prioritas": 2,
      "created_at": "2026-03-29 11:54:28",
      "updated_at": "2026-03-29 11:54:28"
    }
  ]
}
```

### 3. GET Tax by ID - Success (200)
```json
{
  "success": true,
  "message": "Tax retrieved successfully",
  "status": 200,
  "timestamp": "2026-03-29T04:57:35Z",
  "data": {
    "id": "367a2292-a422-4dad-9e1f-79732a7adce9",
    "company_id": "2fa830c2-daec-4ddb-b061-2f7c50b7562b",
    "branch_id": "8fe6b4f9-3dc9-4773-82a8-9f92f0d86458",
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10.00,
    "deskripsi": "Pajak Barang dan Jasa 1",
    "status": "active",
    "prioritas": 1,
    "created_at": "2026-03-29 11:57:23",
    "updated_at": "2026-03-29 11:57:23"
  }
}
```

### 4. UPDATE Tax - Success (200)
```json
{
  "success": true,
  "message": "Tax updated successfully",
  "status": 200,
  "timestamp": "2026-03-29T04:57:40Z",
  "data": {
    "id": "367a2292-a422-4dad-9e1f-79732a7adce9",
    "company_id": "2fa830c2-daec-4ddb-b061-2f7c50b7562b",
    "branch_id": "8fe6b4f9-3dc9-4773-82a8-9f92f0d86458",
    "nama_pajak": "PB1 Updated",
    "tipe_pajak": "pb1",
    "presentase": 11.00,
    "deskripsi": "Updated description",
    "status": "inactive",
    "prioritas": 5,
    "created_at": "2026-03-29 11:57:23",
    "updated_at": "2026-03-29 11:58:15"
  }
}
```

### 5. DELETE Tax - Success (200)
```json
{
  "success": true,
  "message": "Tax deleted successfully",
  "status": 200,
  "timestamp": "2026-03-29T04:58:20Z"
}
```

---

## Error Responses

### 400 Bad Request - Invalid Input
```json
{
  "success": false,
  "message": "Invalid request",
  "status": 400,
  "timestamp": "2026-03-29T04:58:25Z",
  "error": "Key: 'CreateTaxRequest.TipePajak' Error:Field validation for 'TipePajak' failed on the 'oneof' tag"
}
```

### 401 Unauthorized - No Company ID
```json
{
  "success": false,
  "message": "Unauthorized",
  "status": 401,
  "timestamp": "2026-03-29T04:58:30Z",
  "error": "Company ID not found in context"
}
```

### 404 Not Found - Tax Not Found
```json
{
  "success": false,
  "message": "Tax not found",
  "status": 404,
  "timestamp": "2026-03-29T04:58:35Z",
  "error": "tax not found"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Failed to create tax",
  "status": 500,
  "timestamp": "2026-03-29T04:58:40Z",
  "error": "database connection error"
}
```

---

## Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | `true` for success, `false` for error |
| `message` | string | Human-readable message |
| `status` | integer | HTTP status code |
| `timestamp` | string | ISO 8601 timestamp (UTC) |
| `data` | object/array | Response data (success only) |
| `error` | string | Error details (error only) |

---

## Benefits

1. ✅ **Konsisten** dengan handler lain (Category, Product, Order, dll)
2. ✅ **Timestamp** untuk tracking
3. ✅ **Success flag** untuk easy checking
4. ✅ **Status code** included in body
5. ✅ **Structured error** messages
6. ✅ **Easy parsing** di frontend

---

## Migration from Old Format

### Old Format (Before)
```json
{
  "status": "success",
  "message": "Tax created successfully",
  "data": { ... }
}
```

### New Format (After)
```json
{
  "success": true,
  "message": "Tax created successfully",
  "status": 201,
  "timestamp": "2026-03-29T04:57:23Z",
  "data": { ... }
}
```

### Changes:
- ✅ `status` (string) → `success` (boolean)
- ✅ Added `status` (integer) for HTTP code
- ✅ Added `timestamp` (string) for tracking
- ✅ Error format: `error` field instead of `message` only

---

## Implementation

Tax handler sekarang menggunakan:
- `pkg.SuccessResponse(c, status, message, data)` untuk success
- `pkg.ErrorResponse(c, status, message, error)` untuk error

Sama seperti handler lain di project ini.
