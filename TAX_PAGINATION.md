# Tax API - Pagination

GET all taxes endpoint sekarang mendukung pagination dengan query parameters `page` dan `limit`.

## Endpoint

**GET** `/api/v1/external/tax?page=1&limit=10`

## Query Parameters

| Parameter | Type | Required | Default | Max | Description |
|-----------|------|----------|---------|-----|-------------|
| `page` | integer | No | 1 | - | Page number (starts from 1) |
| `limit` | integer | No | 10 | 100 | Items per page |

## Request Examples

### 1. Default (No Params)
```bash
GET /api/v1/external/tax
```
Returns: Page 1, 10 items per page

### 2. Specific Page
```bash
GET /api/v1/external/tax?page=2
```
Returns: Page 2, 10 items per page

### 3. Custom Limit
```bash
GET /api/v1/external/tax?limit=20
```
Returns: Page 1, 20 items per page

### 4. Page + Limit
```bash
GET /api/v1/external/tax?page=3&limit=5
```
Returns: Page 3, 5 items per page

---

## Response Format

### Success Response (200)
```json
{
  "success": true,
  "message": "Taxes retrieved successfully",
  "status": 200,
  "timestamp": "2026-03-29T05:10:00Z",
  "data": [
    {
      "id": "uuid",
      "company_id": "uuid",
      "company_name": "PT Restoran Sejahtera",
      "branch_id": "uuid",
      "branch_name": "Cabang Jakarta Pusat",
      "nama_pajak": "PB1",
      "tipe_pajak": "pb1",
      "presentase": 10.00,
      "deskripsi": "Pajak Barang dan Jasa 1",
      "status": "active",
      "prioritas": 1,
      "created_at": "2026-03-29 11:54:28",
      "updated_at": "2026-03-29 11:54:28"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

### Meta Fields

| Field | Type | Description |
|-------|------|-------------|
| `page` | integer | Current page number |
| `limit` | integer | Items per page |
| `total` | integer | Total number of items |
| `total_pages` | integer | Total number of pages |

---

## Validation Rules

1. **Page**: Must be positive integer (>= 1)
   - Invalid values default to 1
   - Example: `page=0` → defaults to `page=1`

2. **Limit**: Must be positive integer (1-100)
   - Invalid values default to 10
   - Max limit is 100
   - Example: `limit=200` → defaults to `limit=10`

---

## Examples with cURL

### Example 1: Get First Page
```bash
curl -X GET "http://localhost:8080/api/v1/external/tax?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Taxes retrieved successfully",
  "status": 200,
  "timestamp": "2026-03-29T05:10:00Z",
  "data": [ /* 10 items */ ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

### Example 2: Get Second Page
```bash
curl -X GET "http://localhost:8080/api/v1/external/tax?page=2&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Taxes retrieved successfully",
  "status": 200,
  "timestamp": "2026-03-29T05:10:05Z",
  "data": [ /* 10 items */ ],
  "meta": {
    "page": 2,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

### Example 3: Get Last Page
```bash
curl -X GET "http://localhost:8080/api/v1/external/tax?page=3&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Taxes retrieved successfully",
  "status": 200,
  "timestamp": "2026-03-29T05:10:10Z",
  "data": [ /* 5 items */ ],
  "meta": {
    "page": 3,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

### Example 4: Empty Page
```bash
curl -X GET "http://localhost:8080/api/v1/external/tax?page=10&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Taxes retrieved successfully",
  "status": 200,
  "timestamp": "2026-03-29T05:10:15Z",
  "data": [],
  "meta": {
    "page": 10,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

---

## Frontend Implementation

### JavaScript/TypeScript Example
```typescript
async function getTaxes(page = 1, limit = 10) {
  const response = await fetch(
    `http://localhost:8080/api/v1/external/tax?page=${page}&limit=${limit}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  
  const result = await response.json();
  
  console.log('Data:', result.data);
  console.log('Current Page:', result.meta.page);
  console.log('Total Items:', result.meta.total);
  console.log('Total Pages:', result.meta.total_pages);
  
  return result;
}

// Usage
getTaxes(1, 10); // Get first page
getTaxes(2, 20); // Get second page with 20 items
```

### React Example
```jsx
import { useState, useEffect } from 'react';

function TaxList() {
  const [taxes, setTaxes] = useState([]);
  const [meta, setMeta] = useState({});
  const [page, setPage] = useState(1);
  const limit = 10;

  useEffect(() => {
    fetchTaxes();
  }, [page]);

  const fetchTaxes = async () => {
    const response = await fetch(
      `/api/v1/external/tax?page=${page}&limit=${limit}`,
      {
        headers: { 'Authorization': `Bearer ${token}` }
      }
    );
    const result = await response.json();
    setTaxes(result.data);
    setMeta(result.meta);
  };

  return (
    <div>
      {/* Tax List */}
      {taxes.map(tax => (
        <div key={tax.id}>{tax.nama_pajak}</div>
      ))}

      {/* Pagination */}
      <div>
        <button 
          disabled={page === 1}
          onClick={() => setPage(page - 1)}
        >
          Previous
        </button>
        
        <span>Page {meta.page} of {meta.total_pages}</span>
        
        <button 
          disabled={page === meta.total_pages}
          onClick={() => setPage(page + 1)}
        >
          Next
        </button>
      </div>
    </div>
  );
}
```

---

## Benefits

1. ✅ **Performance**: Load only needed data
2. ✅ **UX**: Faster page loads
3. ✅ **Scalability**: Handle large datasets
4. ✅ **Bandwidth**: Reduce data transfer
5. ✅ **Flexibility**: Customizable page size

---

## Notes

- Default: `page=1`, `limit=10`
- Max limit: 100 items per page
- Empty pages return empty array with correct meta
- Sorting: Always by `prioritas DESC, nama_pajak ASC`
- Company-level taxes (branch_id = null) included for all branches
- Branch-level taxes only shown to users in that branch
