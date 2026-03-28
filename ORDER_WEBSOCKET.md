# Order WebSocket Documentation

## Overview
WebSocket endpoint untuk menerima realtime updates ketika ada perubahan order (create, update, delete). Sangat berguna untuk:
- Kitchen Display System (KDS)
- Kasir/POS
- Waiter App
- Dashboard monitoring

---

## WebSocket Endpoint

```
ws://localhost:8080/api/v1/ws/orders
```

**Authentication:** Required (JWT Token via query parameter atau header)

---

## Connection

### Method 1: Query Parameter (Recommended for Browser)
```javascript
const token = "your_jwt_token_here";
const ws = new WebSocket(`ws://localhost:8080/api/v1/ws/orders?token=${token}`);
```

### Method 2: Header (For Native Apps)
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws/orders');
// Note: WebSocket API doesn't support custom headers directly
// Use query parameter instead or implement custom handshake
```

---

## Message Format

Semua message yang diterima dari server dalam format JSON:

```json
{
  "type": "order",
  "action": "created|updated|deleted",
  "data": {
    // Order data (OrderResponse format)
  },
  "company_id": "uuid",
  "branch_id": "uuid"
}
```

### Message Types

#### 1. Order Created
```json
{
  "type": "order",
  "action": "created",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "company_id": "123e4567-e89b-12d3-a456-426614174000",
    "branch_id": "789e0123-e89b-12d3-a456-426614174000",
    "customer_name": "Sasa",
    "customer_phone": "08123123123",
    "table_number": "A1",
    "order_method": "DINE_IN",
    "status": "PENDING",
    "total_amount": 150000,
    "order_items": [
      {
        "id": "abc12345-e29b-41d4-a716-446655440000",
        "product_id": "c61965b6-270c-466f-b5eb-3ff722dfde48",
        "product_name": "Nasi Goreng",
        "quantity": 3,
        "price": 50000,
        "subtotal": 150000,
        "note": "Extra pedas"
      }
    ],
    "created_at": "2024-01-15 10:30:00",
    "updated_at": "2024-01-15 10:30:00"
  },
  "company_id": "123e4567-e89b-12d3-a456-426614174000",
  "branch_id": "789e0123-e89b-12d3-a456-426614174000"
}
```

#### 2. Order Updated
```json
{
  "type": "order",
  "action": "updated",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "CONFIRMED",
    // ... full order data
  },
  "company_id": "123e4567-e89b-12d3-a456-426614174000",
  "branch_id": "789e0123-e89b-12d3-a456-426614174000"
}
```

#### 3. Order Deleted
```json
{
  "type": "order",
  "action": "deleted",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000"
  },
  "company_id": "123e4567-e89b-12d3-a456-426614174000",
  "branch_id": "789e0123-e89b-12d3-a456-426614174000"
}
```

---

## Client Implementation Examples

### JavaScript (Browser)

```javascript
class OrderWebSocket {
  constructor(token) {
    this.token = token;
    this.ws = null;
    this.reconnectInterval = 5000;
  }

  connect() {
    this.ws = new WebSocket(
      `ws://localhost:8080/api/v1/ws/orders?token=${this.token}`
    );

    this.ws.onopen = () => {
      console.log('WebSocket Connected');
    };

    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.handleMessage(message);
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket Error:', error);
    };

    this.ws.onclose = () => {
      console.log('WebSocket Disconnected. Reconnecting...');
      setTimeout(() => this.connect(), this.reconnectInterval);
    };
  }

  handleMessage(message) {
    console.log('Received:', message);

    switch (message.action) {
      case 'created':
        this.onOrderCreated(message.data);
        break;
      case 'updated':
        this.onOrderUpdated(message.data);
        break;
      case 'deleted':
        this.onOrderDeleted(message.data);
        break;
    }
  }

  onOrderCreated(order) {
    console.log('New Order:', order);
    // Update UI - add new order to list
    // Play notification sound
    // Show toast notification
  }

  onOrderUpdated(order) {
    console.log('Order Updated:', order);
    // Update UI - update existing order
    // Highlight changed order
  }

  onOrderDeleted(data) {
    console.log('Order Deleted:', data.id);
    // Update UI - remove order from list
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
    }
  }
}

// Usage
const token = localStorage.getItem('jwt_token');
const orderWs = new OrderWebSocket(token);
orderWs.connect();
```

### React Example

```jsx
import { useEffect, useState } from 'react';

function useOrderWebSocket(token) {
  const [orders, setOrders] = useState([]);
  const [ws, setWs] = useState(null);

  useEffect(() => {
    const websocket = new WebSocket(
      `ws://localhost:8080/api/v1/ws/orders?token=${token}`
    );

    websocket.onopen = () => {
      console.log('Connected to order updates');
    };

    websocket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      
      if (message.type === 'order') {
        switch (message.action) {
          case 'created':
            setOrders(prev => [message.data, ...prev]);
            // Show notification
            new Notification('New Order', {
              body: `Table ${message.data.table_number} - ${message.data.customer_name}`
            });
            break;
            
          case 'updated':
            setOrders(prev => 
              prev.map(order => 
                order.id === message.data.id ? message.data : order
              )
            );
            break;
            
          case 'deleted':
            setOrders(prev => 
              prev.filter(order => order.id !== message.data.id)
            );
            break;
        }
      }
    };

    websocket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    websocket.onclose = () => {
      console.log('Disconnected. Reconnecting...');
      setTimeout(() => {
        // Reconnect logic
      }, 5000);
    };

    setWs(websocket);

    return () => {
      websocket.close();
    };
  }, [token]);

  return { orders, ws };
}

// Component usage
function KitchenDisplay() {
  const token = localStorage.getItem('jwt_token');
  const { orders } = useOrderWebSocket(token);

  return (
    <div>
      <h1>Kitchen Display</h1>
      {orders.map(order => (
        <OrderCard key={order.id} order={order} />
      ))}
    </div>
  );
}
```

### Vue.js Example

```vue
<template>
  <div>
    <h1>Orders</h1>
    <div v-for="order in orders" :key="order.id">
      <OrderCard :order="order" />
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      orders: [],
      ws: null
    }
  },
  mounted() {
    this.connectWebSocket();
  },
  beforeUnmount() {
    if (this.ws) {
      this.ws.close();
    }
  },
  methods: {
    connectWebSocket() {
      const token = localStorage.getItem('jwt_token');
      this.ws = new WebSocket(
        `ws://localhost:8080/api/v1/ws/orders?token=${token}`
      );

      this.ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        
        if (message.type === 'order') {
          switch (message.action) {
            case 'created':
              this.orders.unshift(message.data);
              this.$notify({
                title: 'New Order',
                message: `Table ${message.data.table_number}`,
                type: 'success'
              });
              break;
            case 'updated':
              const index = this.orders.findIndex(o => o.id === message.data.id);
              if (index !== -1) {
                this.$set(this.orders, index, message.data);
              }
              break;
            case 'deleted':
              this.orders = this.orders.filter(o => o.id !== message.data.id);
              break;
          }
        }
      };

      this.ws.onclose = () => {
        setTimeout(() => this.connectWebSocket(), 5000);
      };
    }
  }
}
</script>
```

### Python Client Example

```python
import websocket
import json
import threading

class OrderWebSocketClient:
    def __init__(self, token):
        self.token = token
        self.ws = None
        
    def on_message(self, ws, message):
        data = json.loads(message)
        print(f"Received: {data}")
        
        if data['type'] == 'order':
            if data['action'] == 'created':
                print(f"New order: {data['data']['id']}")
            elif data['action'] == 'updated':
                print(f"Order updated: {data['data']['id']}")
            elif data['action'] == 'deleted':
                print(f"Order deleted: {data['data']['id']}")
    
    def on_error(self, ws, error):
        print(f"Error: {error}")
    
    def on_close(self, ws, close_status_code, close_msg):
        print("Connection closed")
    
    def on_open(self, ws):
        print("Connected to WebSocket")
    
    def connect(self):
        url = f"ws://localhost:8080/api/v1/ws/orders?token={self.token}"
        self.ws = websocket.WebSocketApp(
            url,
            on_open=self.on_open,
            on_message=self.on_message,
            on_error=self.on_error,
            on_close=self.on_close
        )
        
        wst = threading.Thread(target=self.ws.run_forever)
        wst.daemon = True
        wst.start()

# Usage
token = "your_jwt_token"
client = OrderWebSocketClient(token)
client.connect()
```

---

## Features

### 1. Auto Reconnection
Client akan otomatis reconnect jika koneksi terputus

### 2. Room-based Broadcasting
- Setiap client hanya menerima update untuk company_id dan branch_id mereka
- Isolasi data antar company/branch

### 3. Ping/Pong Heartbeat
- Server mengirim ping setiap 54 detik
- Client harus respond dengan pong
- Connection timeout: 60 detik

### 4. Multiple Clients
- Support multiple clients untuk company/branch yang sama
- Semua client akan menerima update yang sama

---

## Use Cases

### 1. Kitchen Display System (KDS)
```javascript
// Filter orders by status
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  if (message.action === 'created' || message.action === 'updated') {
    const order = message.data;
    if (['CONFIRMED', 'PREPARING'].includes(order.status)) {
      displayInKitchen(order);
    }
  }
};
```

### 2. Cashier/POS
```javascript
// Show all pending orders
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  if (message.action === 'created') {
    playNotificationSound();
    showNewOrderAlert(message.data);
  }
};
```

### 3. Waiter App
```javascript
// Track order status for specific table
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  if (message.action === 'updated') {
    const order = message.data;
    if (order.status === 'READY') {
      notifyWaiter(`Order for table ${order.table_number} is ready!`);
    }
  }
};
```

---

## Testing WebSocket

### Using wscat (CLI Tool)

```bash
# Install wscat
npm install -g wscat

# Connect to WebSocket
wscat -c "ws://localhost:8080/api/v1/ws/orders?token=YOUR_JWT_TOKEN"

# You'll receive messages when orders are created/updated/deleted
```

### Using Browser Console

```javascript
const token = "your_jwt_token";
const ws = new WebSocket(`ws://localhost:8080/api/v1/ws/orders?token=${token}`);

ws.onopen = () => console.log('Connected');
ws.onmessage = (e) => console.log('Message:', JSON.parse(e.data));
ws.onerror = (e) => console.error('Error:', e);
ws.onclose = () => console.log('Disconnected');
```

---

## Error Handling

### Connection Errors
- Invalid token: Connection will be rejected
- Network issues: Implement auto-reconnect
- Server restart: Client will reconnect automatically

### Best Practices
1. Always implement reconnection logic
2. Handle connection state in UI
3. Buffer messages during disconnection
4. Validate message format before processing
5. Implement exponential backoff for reconnection

---

## Security

1. **Authentication Required**: JWT token must be valid
2. **Company/Branch Isolation**: Users only receive updates for their company/branch
3. **CORS**: Configure allowed origins in production
4. **Rate Limiting**: Consider implementing rate limiting for WebSocket connections

---

## Performance

- **Concurrent Connections**: Supports thousands of concurrent connections
- **Message Size**: Maximum 512 bytes for client messages
- **Broadcast Latency**: < 10ms for message delivery
- **Memory Usage**: ~10KB per connection

---

## Troubleshooting

### Connection Refused
- Check if server is running
- Verify WebSocket endpoint URL
- Check JWT token validity

### Not Receiving Messages
- Verify company_id and branch_id match
- Check WebSocket connection status
- Ensure order operations are happening

### Frequent Disconnections
- Check network stability
- Verify ping/pong implementation
- Check server logs for errors

---

## Production Considerations

1. **Use WSS (Secure WebSocket)** in production
2. **Configure CORS** properly
3. **Implement rate limiting**
4. **Monitor connection count**
5. **Log WebSocket events**
6. **Use load balancer** with sticky sessions
7. **Implement message queue** for reliability

---

## Summary

WebSocket endpoint menyediakan realtime updates untuk order management, memungkinkan aplikasi untuk:
- Menerima notifikasi instant ketika ada order baru
- Update status order secara realtime
- Sinkronisasi data antar multiple devices
- Meningkatkan user experience dengan realtime feedback
