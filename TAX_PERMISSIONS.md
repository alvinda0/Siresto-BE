# Tax API - Permission Rules

Dokumentasi lengkap tentang permission rules untuk Tax API.

## Permission Matrix

| Action | Company-Level Tax | Branch-Level Tax |
|--------|-------------------|------------------|
| **CREATE** | OWNER only | OWNER, ADMIN, CASHIER (their branch) |
| **READ** | All users in company | All users in that branch |
| **UPDATE** | OWNER only | OWNER, ADMIN, CASHIER (their branch) |
| **DELETE** | OWNER only | OWNER, ADMIN, CASHIER (their branch) |

---

## Detailed Rules

### 1. CREATE Tax

#### Company-Level Tax (branch_id = null)
- ✅ **OWNER**: Can create
- ❌ **ADMIN, CASHIER, etc**: Cannot create

**Example:**
```json
POST /api/v1/external/tax
{
  "nama_pajak": "PB1",
  "tipe_pajak": "pb1",
  "presentase": 10.00
  // No branch_id = company level
}
```

#### Branch-Level Tax (branch_id = UUID)
- ✅ **OWNER**: Can create for any branch
- ✅ **ADMIN, CASHIER**: Can create for their own branch
- ❌ **ADMIN, CASHIER**: Cannot create for other branches

**Example:**
```json
POST /api/v1/external/tax
{
  "branch_id": "branch-uuid",
  "nama_pajak": "Service Charge",
  "tipe_pajak": "sc",
  "presentase": 5.00
}
```

---

### 2. READ Tax

#### GET All Taxes
- **OWNER** (no branch_id): See only company-level taxes
- **ADMIN, CASHIER** (has branch_id): See company-level + their branch taxes

**Example Response for ADMIN:**
```json
{
  "data": [
    {
      "id": "uuid-1",
      "branch_id": null,  // Company-level (visible to all)
      "nama_pajak": "PB1"
    },
    {
      "id": "uuid-2",
      "branch_id": "admin-branch-uuid",  // Their branch tax
      "nama_pajak": "Service Charge"
    }
  ]
}
```

#### GET Tax by ID
- **All users**: Can view company-level taxes
- **Branch users**: Can view their own branch taxes
- **Branch users**: Cannot view other branch taxes

---

### 3. UPDATE Tax

#### Company-Level Tax
- ✅ **OWNER**: Can update
- ❌ **ADMIN, CASHIER, etc**: Cannot update (403 Forbidden)

**Error Response:**
```json
{
  "success": false,
  "message": "Forbidden",
  "status": 403,
  "error": "Only OWNER can update company-level tax"
}
```

#### Branch-Level Tax
- ✅ **OWNER**: Can update any branch tax
- ✅ **ADMIN, CASHIER**: Can update their own branch tax
- ❌ **ADMIN, CASHIER**: Cannot update other branch tax (404)

---

### 4. DELETE Tax

#### Company-Level Tax
- ✅ **OWNER**: Can delete
- ❌ **ADMIN, CASHIER, etc**: Cannot delete (403 Forbidden)

**Error Response:**
```json
{
  "success": false,
  "message": "Forbidden",
  "status": 403,
  "error": "Only OWNER can delete company-level tax"
}
```

#### Branch-Level Tax
- ✅ **OWNER**: Can delete any branch tax
- ✅ **ADMIN, CASHIER**: Can delete their own branch tax
- ❌ **ADMIN, CASHIER**: Cannot delete other branch tax (404)

---

## Use Cases

### Use Case 1: OWNER Creates Company-Wide Tax
```bash
# Login as OWNER
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"owner@restaurant.com","password":"owner123"}'

# Create company-level tax (no branch_id)
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $OWNER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10.00
  }'
```
✅ Success: Tax created for all branches

### Use Case 2: ADMIN Tries to Delete Company-Level Tax
```bash
# Login as ADMIN
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@restaurant.com","password":"admin123"}'

# Try to delete company-level tax
curl -X DELETE http://localhost:8080/api/v1/external/tax/{company-tax-id} \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```
❌ Error 403: "Only OWNER can delete company-level tax"

### Use Case 3: ADMIN Creates Branch-Specific Tax
```bash
# Login as ADMIN (has branch_id)
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@restaurant.com","password":"admin123"}'

# Create branch-level tax (auto use their branch_id)
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "Service Charge",
    "tipe_pajak": "sc",
    "presentase": 5.00
  }'
```
✅ Success: Tax created for their branch only

### Use Case 4: ADMIN Deletes Their Branch Tax
```bash
# Delete their own branch tax
curl -X DELETE http://localhost:8080/api/v1/external/tax/{branch-tax-id} \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```
✅ Success: Tax deleted

---

## Error Responses

### 403 Forbidden - Not OWNER
```json
{
  "success": false,
  "message": "Forbidden",
  "status": 403,
  "timestamp": "2026-03-29T05:30:00Z",
  "error": "Only OWNER can delete company-level tax"
}
```

### 404 Not Found - Wrong Branch
```json
{
  "success": false,
  "message": "Tax not found",
  "status": 404,
  "timestamp": "2026-03-29T05:30:05Z",
  "error": "tax not found"
}
```

---

## Summary

### Company-Level Tax (branch_id = null)
- **Purpose**: Tax yang berlaku untuk semua branch
- **Who can manage**: OWNER only
- **Who can view**: All users in company
- **Example**: PB1 10% (nasional)

### Branch-Level Tax (branch_id = UUID)
- **Purpose**: Tax khusus untuk branch tertentu
- **Who can manage**: OWNER + users in that branch
- **Who can view**: Users in that branch
- **Example**: Service Charge 5% (hanya Branch A)

---

## Best Practices

1. ✅ Use company-level tax for national/standard taxes (PB1, etc)
2. ✅ Use branch-level tax for location-specific charges
3. ✅ OWNER manages company-wide policies
4. ✅ ADMIN/CASHIER manages branch-specific settings
5. ✅ Always check user role before allowing operations
