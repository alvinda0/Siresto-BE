# WebSocket Quick Start Guide

## 🚀 Cara Cepat Test WebSocket

### Step 1: Start Server
```bash
./server.exe
```

### Step 2: Login untuk Mendapatkan Token
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@example.com",
    "password": "password123"
  }'
```

Copy `token` dari response.

### Step 3: Buka WebSocket Test Client
1. Buka file `test_websocket.html` di browser
2. Paste JWT token ke input field
3. Click tombol "Connect"
4. Status akan berubah menjadi "Connected" (hijau)

### Step 4: Test Realtime Updates

**Terminal 1 - Buat Order Baru:**
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "table_number": "A1",
    "order_method": "DINE_IN",
    "order_items": [
      {
        "product_id": "YOUR_PRODUCT_ID",
        "quantity": 2
      }
    ]
  }'
```

**Browser - Lihat Update Realtime:**
- Order baru akan muncul INSTANT di panel "Active Orders"
- Message log akan menampilkan "CREATED"
- Notification sound akan berbunyi
- Counter "Total Orders" akan bertambah

### Step 5: Update Order Status

**Terminal - Update Status:**
```bash
curl -X PUT http://localhost:8080/api/v1/orders/ORDER_ID \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "CONFIRMED"
  }'
```

**Browser - Lihat Update:**
- Order status akan berubah INSTANT
- Message log akan menampilkan "UPDATED"
- Order card akan ter-highlight

---

## 🎯 Use Cases

### Kitchen Display System
```javascript
// Filter hanya order yang perlu dimasak
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  if (msg.action === 'created' || msg.action === 'updated') {
    const order = msg.data;
    if (['CONFIRMED', 'PREPARING'].includes(order.status)) {
      displayInKitchen(order);
      playAlertSound();
    }
  }
};
```

### Cashier/POS
```javascript
// Tampilkan semua order baru
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  if (msg.action === 'created') {
    showNewOrderNotification(msg.data);
    updateOrderList(msg.data);
  }
};
```

### Waiter App
```javascript
// Notifikasi ketika order ready
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  if (msg.action === 'updated' && msg.data.status === 'READY') {
    notifyWaiter(`Table ${msg.data.table_number} ready!`);
  }
};
```

---

## 🔧 Troubleshooting

### Connection Failed
- ✅ Pastikan server running
- ✅ Check JWT token valid
- ✅ Verify WebSocket URL correct

### Not Receiving Messages
- ✅ Check company_id dan branch_id match
- ✅ Verify order operations happening
- ✅ Check browser console for errors

### Frequent Disconnections
- ✅ Check network stability
- ✅ Verify token not expired
- ✅ Check server logs

---

## 📱 Integration Examples

### React
```jsx
import { useEffect, useState } from 'react';

function OrderDashboard() {
  const [orders, setOrders] = useState([]);
  
  useEffect(() => {
    const token = localStorage.getItem('jwt_token');
    const ws = new WebSocket(
      `ws://localhost:8080/api/v1/ws/orders?token=${token}`
    );
    
    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      if (msg.action === 'created') {
        setOrders(prev => [msg.data, ...prev]);
      }
    };
    
    return () => ws.close();
  }, []);
  
  return (
    <div>
      {orders.map(order => (
        <OrderCard key={order.id} order={order} />
      ))}
    </div>
  );
}
```

### Vue.js
```vue
<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

const orders = ref([]);
let ws = null;

onMounted(() => {
  const token = localStorage.getItem('jwt_token');
  ws = new WebSocket(
    `ws://localhost:8080/api/v1/ws/orders?token=${token}`
  );
  
  ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    if (msg.action === 'created') {
      orders.value.unshift(msg.data);
    }
  };
});

onUnmounted(() => {
  if (ws) ws.close();
});
</script>
```

---

## 🎨 Features

✅ **Realtime Updates** - Instant notification untuk semua perubahan order
✅ **Room-based** - Hanya menerima update untuk company/branch Anda
✅ **Auto Reconnect** - Otomatis reconnect jika koneksi terputus
✅ **Multiple Clients** - Support banyak client sekaligus
✅ **Notification Sound** - Audio alert untuk order baru
✅ **Visual Feedback** - Highlight dan animation untuk updates

---

## 📖 Full Documentation

Lihat `ORDER_WEBSOCKET.md` untuk dokumentasi lengkap.
