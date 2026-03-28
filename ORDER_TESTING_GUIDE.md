# Order Testing Guide

## ✅ Migration Berhasil!

Tables `orders` dan `order_items` sudah dibuat dengan benar.

---

## 🧪 Cara Test Create Order

### Step 1: Get Product ID

Pertama, ambil product_id yang valid:

```bash
GET http://localhost:8080/api/v1/external/products?limit=1
Authorization: Bearer YOUR_TOKEN
```

**Response:**
```json
{
  "data": [
    {
      "id": "c61965b6-270c-466f-b5eb-3ff722dfde48",  // <-- Copy this
      "name": "Nasi Goreng",
      "price": 50000
    }
  ]
}
```

### Step 2: Create Order

```bash
POST http://localhost:8080/api/v1/external/orders
Authorization: Bearer YOUR_TOKEN
Content-Type: application/json
```

**Body:**
```json
{
  "table_number": "A1",
  "customer_name": "John Doe",
  "customer_phone": "08123456789",
  "order_method": "DINE_IN",
  "notes": "Extra pedas",
  "order_items": [
    {
      "product_id": "c61965b6-270c-466f-b5eb-3ff722dfde48",
      "quantity": 2,
      "note": "Tanpa sayur"
    }
  ]
}
```

**Expected Response (201):**
```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "company_id": "...",
    "branch_id": "...",
    "customer_name": "John Doe",
    "table_number": "A1",
    "status": "PENDING",
    "total_amount": 100000,
    "order_items": [
      {
        "product_id": "c61965b6-270c-466f-b5eb-3ff722dfde48",
        "product_name": "Nasi Goreng",
        "quantity": 2,
        "price": 50000,
        "subtotal": 100000
      }
    ]
  }
}
```

### Step 3: Test WebSocket (Realtime)

1. **Buka `test_websocket.html` di browser**
2. **Paste JWT token**
3. **Click "Connect"**
4. **Buat order baru dari Postman**
5. **Lihat realtime update di browser!** 🎉

---

## 📋 All Order Endpoints

### 1. Create Order
```
POST /api/v1/external/orders
Authorization: Bearer TOKEN
```

### 2. Get All Orders
```
GET /api/v1/external/orders?page=1&limit=10
Authorization: Bearer TOKEN
```

### 3. Get Order by ID
```
GET /api/v1/external/orders/:id
Authorization: Bearer TOKEN
```

### 4. Update Order
```
PUT /api/v1/external/orders/:id
Authorization: Bearer TOKEN

Body:
{
  "status": "CONFIRMED",
  "notes": "Updated notes"
}
```

### 5. Delete Order
```
DELETE /api/v1/external/orders/:id
Authorization: Bearer TOKEN
```

### 6. Create Public Order (No Auth)
```
POST /api/v1/public/orders
Content-Type: application/json

Body:
{
  "company_id": "YOUR_COMPANY_ID",
  "branch_id": "YOUR_BRANCH_ID",
  "table_number": "A1",
  "order_method": "DINE_IN",
  "order_items": [...]
}
```

### 7. WebSocket (Realtime)
```
ws://localhost:8080/api/v1/ws/orders?token=YOUR_TOKEN
```

---

## 🔧 Troubleshooting

### Error: "product not found"
- Pastikan product_id valid
- Get product list dulu: `GET /api/v1/external/products`

### Error: "product does not belong to this branch"
- Product harus dari branch yang sama dengan user
- Check branch_id di JWT token

### Error: "product is not available"
- Product.is_available harus true
- Update product: `PUT /api/v1/external/products/:id`

### WebSocket tidak connect
- Gunakan endpoint baru: `ws://localhost:8080/api/v1/ws/orders`
- Token di query parameter, bukan header

---

## ✅ Checklist

- [ ] Migration berhasil (tables created)
- [ ] Server running
- [ ] Login berhasil (dapat token)
- [ ] Get products berhasil (dapat product_id)
- [ ] Create order berhasil (201 response)
- [ ] WebSocket connected
- [ ] Realtime update works

---

## 🎯 Next Steps

1. Test semua CRUD operations
2. Test WebSocket realtime updates
3. Test order status flow (PENDING → CONFIRMED → PREPARING → READY → COMPLETED)
4. Test multiple clients di WebSocket
5. Test public order endpoint

Happy testing! 🚀
