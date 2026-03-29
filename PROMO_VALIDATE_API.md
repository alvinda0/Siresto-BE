# Promo Validate API

API untuk validasi kode promo sebelum digunakan dalam order.

## Endpoint

```
GET /api/v1/external/promos/validate/:code
```

## Description

Endpoint ini digunakan untuk mengecek apakah kode promo valid dan bisa digunakan. Validasi meliputi:
- Promo code exists
- Promo is active
- Promo has started
- Promo has not expired
- Promo quota is available

## Request

### Headers
```
Authorization: Bearer YOUR_TOKEN
Content-Type: application/json
```

### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `code` | String | ✅ | Kode promo yang akan divalidasi |

### Example Request
```bash
GET /api/v1/external/promos/validate/WEEKEND50
```

## Response

### Success Response (Valid Promo)

**Status Code:** `200 OK`

```json
{
  "success": true,
  "message": "Promo is valid and can be used",
  "data": {
    "valid": true,
    "message": "Promo is valid and can be used",
    "promo": {
      "id": "dfba1930-f5ee-48fd-b657-f740e5624ac4",
      "company_id": "2fa830c2-daec-4ddb-b061-2f7c50b7562b",
      "company_name": "Company Name",
      "branch_id": null,
      "branch_name": null,
      "name": "Weekend Special",
      "code": "WEEKEND50",
      "promo_category": "normal",
      "type": "percentage",
      "value": 50,
      "max_discount": 100000,
      "min_transaction": 200000,
      "quota": 100,
      "used_count": 25,
      "remaining_quota": 75,
      "start_date": "2024-12-01",
      "end_date": "2024-12-31",
      "is_active": true,
      "is_expired": false,
      "is_available": true,
      "created_at": "2024-12-01 10:00:00",
      "updated_at": "2024-12-01 10:00:00"
    }
  },
  "timestamp": "2024-12-15T10:30:00Z"
}
```

### Error Responses

#### 1. Promo Not Found

**Status Code:** `200 OK`

```json
{
  "success": true,
  "message": "Promo code not found",
  "data": {
    "valid": false,
    "message": "Promo code not found"
  },
  "timestamp": "2024-12-15T10:30:00Z"
}
```

#### 2. Promo Not Active

**Status Code:** `200 OK`

```json
{
  "success": true,
  "message": "Promo is not active",
  "data": {
    "valid": false,
    "message": "Promo is not active",
    "promo": {
      "id": "...",
      "name": "Weekend Special",
      "code": "WEEKEND50",
      "is_active": false,
      ...
    }
  },
  "timestamp": "2024-12-15T10:30:00Z"
}
```

#### 3. Promo Not Started Yet

**Status Code:** `200 OK`

```json
{
  "success": true,
  "message": "Promo will start on 2024-12-20",
  "data": {
    "valid": false,
    "message": "Promo will start on 2024-12-20",
    "promo": {
      "id": "...",
      "name": "Christmas Sale",
      "code": "XMAS2024",
      "start_date": "2024-12-20",
      ...
    }
  },
  "timestamp": "2024-12-15T10:30:00Z"
}
```

#### 4. Promo Expired

**Status Code:** `200 OK`

```json
{
  "success": true,
  "message": "Promo expired on 2024-11-30",
  "data": {
    "valid": false,
    "message": "Promo expired on 2024-11-30",
    "promo": {
      "id": "...",
      "name": "November Sale",
      "code": "NOV2024",
      "end_date": "2024-11-30",
      "is_expired": true,
      ...
    }
  },
  "timestamp": "2024-12-15T10:30:00Z"
}
```

#### 5. Quota Exhausted

**Status Code:** `200 OK`

```json
{
  "success": true,
  "message": "Promo quota has been exhausted",
  "data": {
    "valid": false,
    "message": "Promo quota has been exhausted",
    "promo": {
      "id": "...",
      "name": "Limited Offer",
      "code": "LIMITED100",
      "quota": 100,
      "used_count": 100,
      "remaining_quota": 0,
      ...
    }
  },
  "timestamp": "2024-12-15T10:30:00Z"
}
```

#### 6. Invalid Request

**Status Code:** `400 Bad Request`

```json
{
  "success": false,
  "message": "Invalid request",
  "error": "Promo code is required",
  "timestamp": "2024-12-15T10:30:00Z"
}
```

#### 7. Unauthorized

**Status Code:** `401 Unauthorized`

```json
{
  "success": false,
  "message": "Unauthorized",
  "error": "Company ID not found in context",
  "timestamp": "2024-12-15T10:30:00Z"
}
```

## Validation Rules

### 1. Promo Code Exists
- Promo code harus ada di database
- Promo harus milik company yang sama dengan user

### 2. Promo Active
- `is_active` harus `true`

### 3. Promo Started
- Current date >= `start_date`

### 4. Promo Not Expired
- Current date <= `end_date`

### 5. Quota Available
- Jika `quota` tidak null: `used_count` < `quota`
- Jika `quota` null: unlimited

## Use Cases

### 1. Validate Before Order
```javascript
// Frontend: Check promo before creating order
const validatePromo = async (promoCode) => {
  const response = await fetch(`/api/v1/external/promos/validate/${promoCode}`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  const data = await response.json();
  
  if (data.data.valid) {
    // Promo valid, bisa digunakan
    console.log('Promo valid:', data.data.promo);
    return data.data.promo;
  } else {
    // Promo tidak valid
    alert(data.data.message);
    return null;
  }
};
```

### 2. Real-time Validation
```javascript
// Validate saat user mengetik promo code
const promoInput = document.getElementById('promo-code');
const promoStatus = document.getElementById('promo-status');

promoInput.addEventListener('blur', async () => {
  const code = promoInput.value.trim();
  if (!code) return;
  
  const response = await fetch(`/api/v1/external/promos/validate/${code}`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  
  const data = await response.json();
  
  if (data.data.valid) {
    promoStatus.textContent = `✓ ${data.data.message}`;
    promoStatus.className = 'success';
  } else {
    promoStatus.textContent = `✗ ${data.data.message}`;
    promoStatus.className = 'error';
  }
});
```

### 3. Display Promo Info
```javascript
// Show promo details after validation
const showPromoInfo = (promo) => {
  const info = `
    <div class="promo-info">
      <h3>${promo.name}</h3>
      <p>Code: ${promo.code}</p>
      <p>Discount: ${promo.type === 'percentage' ? promo.value + '%' : 'Rp ' + promo.value}</p>
      ${promo.min_transaction ? `<p>Min. Transaction: Rp ${promo.min_transaction}</p>` : ''}
      ${promo.max_discount ? `<p>Max. Discount: Rp ${promo.max_discount}</p>` : ''}
      ${promo.remaining_quota ? `<p>Remaining: ${promo.remaining_quota} uses</p>` : ''}
      <p>Valid until: ${promo.end_date}</p>
    </div>
  `;
  document.getElementById('promo-details').innerHTML = info;
};
```

## Testing

### Test 1: Valid Promo
```bash
curl -X GET "http://localhost:8080/api/v1/external/promos/validate/WEEKEND50" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected:** `valid: true`

### Test 2: Invalid Code
```bash
curl -X GET "http://localhost:8080/api/v1/external/promos/validate/INVALID123" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected:** `valid: false`, message: "Promo code not found"

### Test 3: Expired Promo
```bash
curl -X GET "http://localhost:8080/api/v1/external/promos/validate/EXPIRED2023" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected:** `valid: false`, message: "Promo expired on ..."

### Test 4: Inactive Promo
```bash
curl -X GET "http://localhost:8080/api/v1/external/promos/validate/INACTIVE" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected:** `valid: false`, message: "Promo is not active"

## Notes

1. **Case Sensitive**: Promo code adalah case-sensitive
2. **Multi-tenant**: Hanya bisa validate promo dari company yang sama
3. **Branch Access**: 
   - User dengan branch bisa validate company-level dan branch-level promos
   - OWNER bisa validate company-level promos saja
4. **No Side Effects**: Endpoint ini hanya validasi, tidak mengubah data
5. **Real-time**: Validasi dilakukan real-time berdasarkan waktu server

## Integration with Order

Setelah validasi berhasil, gunakan promo code di order:

```bash
POST /api/v1/external/orders
{
  "items": [...],
  "promo_code": "WEEKEND50",
  ...
}
```

Order service akan:
1. Validate promo lagi
2. Check promo category (normal/product/bundle)
3. Calculate discount
4. Apply to order
5. Increment used_count

---

**Version**: 1.0
**Endpoint**: `GET /api/v1/external/promos/validate/:code`
**Authentication**: Required (Bearer Token)
**Rate Limit**: None
