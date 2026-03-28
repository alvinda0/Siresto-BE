# WebSocket di Postman - Quick Guide

## ✅ Endpoint WebSocket yang Benar

```
ws://localhost:8080/api/v1/ws/orders?token={{jwt_token}}
```

**PENTING:** Endpoint sudah dipindahkan dari `/api/v1/external/orders/ws` ke `/api/v1/ws/orders`

---

## 🚀 Cara Connect di Postman

### Step 1: Buka WebSocket Request
1. Click **New** → **WebSocket Request**
2. Atau click **+** tab → pilih **WebSocket**

### Step 2: Masukkan URL
```
ws://localhost:8080/api/v1/ws/orders?token=YOUR_JWT_TOKEN
```

**Ganti `YOUR_JWT_TOKEN` dengan token dari login**

### Step 3: Click "Connect"
- Status akan berubah menjadi **Connected** (hijau)
- Jika error 401, berarti token invalid/expired

### Step 4: Test Realtime Updates
Buka terminal/Postman lain dan buat order baru:

```bash
POST http://localhost:8080/api/v1/orders
Authorization: Bearer YOUR_TOKEN
Content-Type: application/json

{
  "table_number": "A1",
  "order_method": "DINE_IN",
  "order_items": [
    {
      "product_id": "YOUR_PRODUCT_ID",
      "quantity": 2
    }
  ]
}
```

**Hasil:** WebSocket akan menerima message realtime!

---

## 📨 Format Message yang Diterima

### Order Created
```json
{
  "type": "order",
  "action": "created",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "table_number": "A1",
    "status": "PENDING",
    "total_amount": 150000,
    "order_items": [...]
  }
}
```

### Order Updated
```json
{
  "type": "order",
  "action": "updated",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "CONFIRMED",
    ...
  }
}
```

### Order Deleted
```json
{
  "type": "order",
  "action": "deleted",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

---

## 🔧 Troubleshooting

### Error 401: Unauthorized
**Penyebab:**
- Token invalid atau expired
- Token tidak dikirim di query parameter
- Format URL salah

**Solusi:**
1. Login ulang untuk mendapatkan token baru
2. Pastikan URL: `ws://localhost:8080/api/v1/ws/orders?token=YOUR_TOKEN`
3. Copy-paste token langsung dari response login

### Connection Timeout
**Penyebab:**
- Server tidak running
- Port salah
- Firewall blocking

**Solusi:**
1. Pastikan server running: `./server.exe`
2. Check port: default 8080
3. Test REST API dulu: `GET http://localhost:8080/api/v1/roles`

### Not Receiving Messages
**Penyebab:**
- Tidak ada order activity
- Company/Branch ID tidak match
- WebSocket disconnected

**Solusi:**
1. Check status connection (harus hijau "Connected")
2. Buat order baru untuk trigger message
3. Pastikan token dari user yang sama company/branch

---

## 💡 Tips

### 1. Gunakan Environment Variable
Di Postman, buat variable:
- `base_url`: `http://localhost:8080`
- `ws_url`: `ws://localhost:8080`
- `jwt_token`: (auto-set dari login response)

URL WebSocket jadi:
```
{{ws_url}}/api/v1/ws/orders?token={{jwt_token}}
```

### 2. Auto-Save Token dari Login
Di login request, tambahkan Test script:
```javascript
pm.test("Save token", function () {
    var jsonData = pm.response.json();
    pm.environment.set("jwt_token", jsonData.data.token);
});
```

### 3. Multiple WebSocket Connections
Postman support multiple WebSocket tabs. Buka beberapa tab untuk simulate multiple clients.

---

## 📋 Checklist Testing

- [ ] Server running
- [ ] Login berhasil dan dapat token
- [ ] WebSocket connected (status hijau)
- [ ] Buat order baru → terima message "created"
- [ ] Update order → terima message "updated"
- [ ] Delete order → terima message "deleted"

---

## 🎯 Use Cases

### Kitchen Display
```
1. Connect WebSocket
2. Filter messages dengan action: "created" atau "updated"
3. Tampilkan hanya order dengan status: CONFIRMED, PREPARING
```

### Cashier/POS
```
1. Connect WebSocket
2. Tampilkan semua order baru (action: "created")
3. Update list ketika ada perubahan status
```

### Waiter App
```
1. Connect WebSocket
2. Filter order dengan status: "READY"
3. Notifikasi waiter untuk ambil order
```

---

## 📖 Full Documentation

- **REST API:** `ORDER_API.md`
- **WebSocket Detail:** `ORDER_WEBSOCKET.md`
- **Quick Start:** `WEBSOCKET_QUICK_START.md`
- **HTML Test Client:** `test_websocket.html`
