# Tax API Testing Guide

Panduan untuk testing Tax API menggunakan curl atau Postman.

## Prerequisites

1. Server berjalan di `http://localhost:8080`
2. Sudah memiliki token authentication (login sebagai EXTERNAL user)
3. Set token ke environment variable:
```bash
export TOKEN="your_jwt_token_here"
```

---

## 1. Create Tax

### Test Case: Create PB1 Tax
```bash
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10.00,
    "deskripsi": "Pajak Barang dan Jasa 1",
    "status": "active",
    "prioritas": 1
  }'
```

### Test Case: Create Service Charge
```bash
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "Service Charge",
    "tipe_pajak": "sc",
    "presentase": 5.00,
    "deskripsi": "Biaya layanan",
    "status": "active",
    "prioritas": 2
  }'
```

### Test Case: Validation Error - Invalid tipe_pajak
```bash
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "Invalid Tax",
    "tipe_pajak": "invalid",
    "presentase": 10.00
  }'
```
Expected: 400 Bad Request

### Test Case: Validation Error - Presentase > 100
```bash
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "Invalid Tax",
    "tipe_pajak": "pb1",
    "presentase": 150.00
  }'
```
Expected: 400 Bad Request

---

## 2. Get All Taxes

```bash
curl -X GET http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN"
```

Expected: List of all taxes, sorted by prioritas DESC, nama_pajak ASC

---

## 3. Get Tax by ID

```bash
# Replace {tax_id} with actual UUID
curl -X GET http://localhost:8080/api/v1/external/tax/{tax_id} \
  -H "Authorization: Bearer $TOKEN"
```

### Test Case: Invalid UUID
```bash
curl -X GET http://localhost:8080/api/v1/external/tax/invalid-uuid \
  -H "Authorization: Bearer $TOKEN"
```
Expected: 400 Bad Request

### Test Case: Non-existent Tax
```bash
curl -X GET http://localhost:8080/api/v1/external/tax/00000000-0000-0000-0000-000000000000 \
  -H "Authorization: Bearer $TOKEN"
```
Expected: 404 Not Found

---

## 4. Update Tax

### Test Case: Update All Fields
```bash
# Replace {tax_id} with actual UUID
curl -X PUT http://localhost:8080/api/v1/external/tax/{tax_id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "PB1 Updated",
    "tipe_pajak": "pb1",
    "presentase": 11.00,
    "deskripsi": "Updated description",
    "status": "inactive",
    "prioritas": 5
  }'
```

### Test Case: Partial Update (Only nama_pajak)
```bash
curl -X PUT http://localhost:8080/api/v1/external/tax/{tax_id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "New Name Only"
  }'
```

### Test Case: Update Status to Inactive
```bash
curl -X PUT http://localhost:8080/api/v1/external/tax/{tax_id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "inactive"
  }'
```

---

## 5. Delete Tax

```bash
# Replace {tax_id} with actual UUID
curl -X DELETE http://localhost:8080/api/v1/external/tax/{tax_id} \
  -H "Authorization: Bearer $TOKEN"
```

### Test Case: Delete Non-existent Tax
```bash
curl -X DELETE http://localhost:8080/api/v1/external/tax/00000000-0000-0000-0000-000000000000 \
  -H "Authorization: Bearer $TOKEN"
```
Expected: 404 Not Found

---

## Complete Test Flow

```bash
#!/bin/bash

# Set your token
export TOKEN="your_jwt_token_here"

echo "=== 1. Create PB1 Tax ==="
TAX_ID=$(curl -s -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10.00,
    "deskripsi": "Pajak Barang dan Jasa 1",
    "status": "active",
    "prioritas": 1
  }' | jq -r '.data.id')

echo "Created Tax ID: $TAX_ID"

echo -e "\n=== 2. Create Service Charge ==="
curl -X POST http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "Service Charge",
    "tipe_pajak": "sc",
    "presentase": 5.00,
    "deskripsi": "Biaya layanan",
    "status": "active",
    "prioritas": 2
  }'

echo -e "\n=== 3. Get All Taxes ==="
curl -X GET http://localhost:8080/api/v1/external/tax \
  -H "Authorization: Bearer $TOKEN"

echo -e "\n=== 4. Get Tax by ID ==="
curl -X GET http://localhost:8080/api/v1/external/tax/$TAX_ID \
  -H "Authorization: Bearer $TOKEN"

echo -e "\n=== 5. Update Tax ==="
curl -X PUT http://localhost:8080/api/v1/external/tax/$TAX_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "PB1 Updated",
    "presentase": 11.00,
    "status": "inactive"
  }'

echo -e "\n=== 6. Delete Tax ==="
curl -X DELETE http://localhost:8080/api/v1/external/tax/$TAX_ID \
  -H "Authorization: Bearer $TOKEN"

echo -e "\n=== 7. Verify Deletion ==="
curl -X GET http://localhost:8080/api/v1/external/tax/$TAX_ID \
  -H "Authorization: Bearer $TOKEN"
```

---

## Postman Collection

### Environment Variables
```
base_url: http://localhost:8080
token: your_jwt_token_here
tax_id: (will be set automatically)
```

### Collection Structure
```
Tax API
├── Create Tax (PB1)
├── Create Tax (Service Charge)
├── Get All Taxes
├── Get Tax by ID
├── Update Tax
├── Delete Tax
└── Validation Tests
    ├── Invalid tipe_pajak
    ├── Invalid presentase
    ├── Invalid UUID
    └── Non-existent Tax
```

### Test Scripts (Postman)

**Create Tax - Tests tab:**
```javascript
if (pm.response.code === 201) {
    const response = pm.response.json();
    pm.environment.set("tax_id", response.data.id);
    pm.test("Tax created successfully", () => {
        pm.expect(response.status).to.eql("success");
    });
}
```

**Get All Taxes - Tests tab:**
```javascript
pm.test("Status is success", () => {
    const response = pm.response.json();
    pm.expect(response.status).to.eql("success");
});

pm.test("Data is array", () => {
    const response = pm.response.json();
    pm.expect(response.data).to.be.an("array");
});
```

---

## Expected Results

### Success Scenarios
- Create: 201 Created with tax data
- Get All: 200 OK with array of taxes
- Get by ID: 200 OK with single tax
- Update: 200 OK with updated tax data
- Delete: 200 OK with success message

### Error Scenarios
- Invalid UUID: 400 Bad Request
- Non-existent tax: 404 Not Found
- Invalid tipe_pajak: 400 Bad Request
- Invalid presentase: 400 Bad Request
- Missing required fields: 400 Bad Request
- No authentication: 401 Unauthorized
