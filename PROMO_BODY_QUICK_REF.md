# Promo Body - Quick Reference

## POST /api/promos

### 🟢 Normal
```json
{
  "name": "Diskon 20%",
  "code": "DISC20",
  "promo_category": "normal",
  "type": "percentage",
  "value": 20,
  "start_date": "2024-12-01",
  "end_date": "2024-12-31"
}
```

### 🟡 Product
```json
{
  "name": "Diskon Laptop 50%",
  "code": "LAPTOP50",
  "promo_category": "product",
  "type": "percentage",
  "value": 50,
  "product_ids": ["uuid1", "uuid2"],
  "start_date": "2024-12-01",
  "end_date": "2024-12-31"
}
```

### 🔴 Bundle
```json
{
  "name": "Paket Gaming",
  "code": "GAMING999",
  "promo_category": "bundle",
  "type": "fixed",
  "value": 1000000,
  "bundle_items": [
    {"product_id": "uuid1", "quantity": 1},
    {"product_id": "uuid2", "quantity": 2}
  ],
  "start_date": "2024-12-01",
  "end_date": "2024-12-31"
}
```

---

## PUT /api/promos/:id

### Update Any Field (Optional)
```json
{
  "name": "New Name",
  "value": 30,
  "is_active": false
}
```

### Update Products
```json
{
  "product_ids": ["uuid3", "uuid4"]
}
```

### Update Bundle Items
```json
{
  "bundle_items": [
    {"product_id": "uuid5", "quantity": 2},
    {"product_id": "uuid6", "quantity": 1}
  ]
}
```

---

## Optional Fields (All Categories)

```json
{
  "branch_id": "uuid-or-null",
  "max_discount": 100000,
  "min_transaction": 200000,
  "quota": 100,
  "is_active": true
}
```

---

## Field Types

| Field | Type | Example |
|-------|------|---------|
| `name` | String | "Diskon Ramadan" |
| `code` | String | "RAMADAN20" |
| `promo_category` | String | "normal" / "product" / "bundle" |
| `type` | String | "percentage" / "fixed" |
| `value` | Number | 20 atau 100000 |
| `max_discount` | Number | 150000 |
| `min_transaction` | Number | 300000 |
| `quota` | Integer | 100 |
| `product_ids` | Array[UUID] | ["uuid1", "uuid2"] |
| `bundle_items` | Array[Object] | [{"product_id": "uuid", "quantity": 1}] |
| `start_date` | String | "2024-12-01" |
| `end_date` | String | "2024-12-31" |
| `is_active` | Boolean | true / false |
| `branch_id` | UUID/null | "uuid" atau null |

---

## Validation Rules

✅ **Normal**: No product_ids, no bundle_items
✅ **Product**: product_ids required (min 1)
✅ **Bundle**: bundle_items required (min 2)
✅ **Dates**: end_date > start_date
✅ **Type**: percentage or fixed
✅ **Value**: >= 0

---

## Quick Test

```bash
# Login
TOKEN=$(curl -X POST http://localhost:8080/api/login \
  -d '{"email":"owner@company1.com","password":"password123"}' \
  | jq -r '.data.token')

# Create Normal
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Test","code":"TEST","promo_category":"normal","type":"percentage","value":10,"start_date":"2024-12-01","end_date":"2024-12-31"}'

# Create Product
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Test","code":"TEST2","promo_category":"product","type":"percentage","value":50,"product_ids":["UUID"],"start_date":"2024-12-01","end_date":"2024-12-31"}'

# Create Bundle
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Test","code":"TEST3","promo_category":"bundle","type":"fixed","value":100000,"bundle_items":[{"product_id":"UUID1","quantity":1},{"product_id":"UUID2","quantity":2}],"start_date":"2024-12-01","end_date":"2024-12-31"}'
```

---

**Full Documentation**: `PROMO_API_BODY_REFERENCE.md`
