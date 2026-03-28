# Order API Filter Parameters

## 📋 GET /api/v1/external/orders - Filter Parameters

Endpoint untuk mendapatkan list orders sekarang support filter parameters untuk pencarian dan filtering yang lebih spesifik.

---

## 🔍 Available Filters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `page` | integer | Page number (default: 1) | `1`, `2`, `3` |
| `limit` | integer | Items per page (default: 10) | `10`, `50`, `100` |
| `status` | string | Filter by order status | `PENDING`, `CONFIRMED`, `PREPARING`, `READY`, `COMPLETED`, `CANCELLED` |
| `method` | string | Filter by order method | `DINE_IN`, `TAKE_AWAY`, `DELIVERY` |
| `customer` | string | Search by customer name (partial, case-insensitive) | `john`, `doe` |
| `order_id` | string | Search by order ID (partial) | `550e8400`, `abc123` |

---

## 📝 Examples

### 1. Get All Orders (No Filter)
```
GET /api/v1/external/orders?page=1&limit=10
Authorization: Bearer YOUR_TOKEN
```

**Response:**
```json
{
  "success": true,
  "message": "Orders retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "customer_name": "John Doe",
      "table_number": "A1",
      "status": "PENDING",
      "order_method": "DINE_IN",
      "total_amount": 150000
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 50,
    "total_pages": 5
  }
}
```

### 2. Filter by Status
```
GET /api/v1/external/orders?status=PENDING
Authorization: Bearer YOUR_TOKEN
```

Hanya menampilkan order dengan status PENDING.

### 3. Filter by Order Method
```
GET /api/v1/external/orders?method=DINE_IN
Authorization: Bearer YOUR_TOKEN
```

Hanya menampilkan order dengan method DINE_IN.

### 4. Search by Customer Name
```
GET /api/v1/external/orders?customer=john
Authorization: Bearer YOUR_TOKEN
```

Menampilkan order yang customer name mengandung "john" (case-insensitive).

### 5. Search by Order ID
```
GET /api/v1/external/orders?order_id=550e8400
Authorization: Bearer YOUR_TOKEN
```

Menampilkan order yang ID mengandung "550e8400".

### 6. Multiple Filters (Kombinasi)
```
GET /api/v1/external/orders?status=PENDING&method=DINE_IN&page=1&limit=20
Authorization: Bearer YOUR_TOKEN
```

Filter order PENDING dengan method DINE_IN, page 1, 20 items per page.

### 7. Search Customer + Filter Status
```
GET /api/v1/external/orders?customer=john&status=READY
Authorization: Bearer YOUR_TOKEN
```

Order untuk customer "john" yang statusnya READY.

---

## 🎯 Use Cases

### Kitchen Display - Get Orders to Cook
```
GET /api/v1/external/orders?status=CONFIRMED&limit=50
Authorization: Bearer YOUR_TOKEN
```

Ambil semua order yang sudah confirmed dan perlu dimasak.

### Cashier - Get Dine-In Orders
```
GET /api/v1/external/orders?method=DINE_IN&status=PENDING
Authorization: Bearer YOUR_TOKEN
```

Ambil semua order dine-in yang pending.

### Waiter - Get Ready Orders
```
GET /api/v1/external/orders?status=READY
Authorization: Bearer YOUR_TOKEN
```

Ambil semua order yang ready untuk diambil.

### Customer Service - Search by Customer
```
GET /api/v1/external/orders?customer=john%20doe
Authorization: Bearer YOUR_TOKEN
```

Cari order berdasarkan nama customer.

### Order Tracking - Find Specific Order
```
GET /api/v1/external/orders?order_id=550e8400
Authorization: Bearer YOUR_TOKEN
```

Track order berdasarkan ID.

---

## 💻 Implementation Examples

### JavaScript/Fetch
```javascript
async function getOrders(filters = {}) {
  const token = localStorage.getItem('jwt_token');
  
  // Build query params
  const params = new URLSearchParams({
    page: filters.page || 1,
    limit: filters.limit || 10
  });
  
  if (filters.status) params.append('status', filters.status);
  if (filters.method) params.append('method', filters.method);
  if (filters.customer) params.append('customer', filters.customer);
  if (filters.order_id) params.append('order_id', filters.order_id);
  
  const response = await fetch(
    `http://localhost:8080/api/v1/external/orders?${params.toString()}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  
  return await response.json();
}

// Usage
const pendingOrders = await getOrders({ status: 'PENDING' });
const dineInOrders = await getOrders({ method: 'DINE_IN', page: 1, limit: 20 });
const customerOrders = await getOrders({ customer: 'john' });
```

### Axios
```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('jwt_token')}`
  }
});

// Get pending orders
const getPendingOrders = async () => {
  const response = await api.get('/external/orders', {
    params: { status: 'PENDING' }
  });
  return response.data;
};

// Get dine-in orders
const getDineInOrders = async () => {
  const response = await api.get('/external/orders', {
    params: { method: 'DINE_IN', limit: 50 }
  });
  return response.data;
};

// Search by customer
const searchByCustomer = async (customerName) => {
  const response = await api.get('/external/orders', {
    params: { customer: customerName }
  });
  return response.data;
};
```

### React Hook
```jsx
import { useState, useEffect } from 'react';
import axios from 'axios';

function useOrders(filters = {}) {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(false);
  const [meta, setMeta] = useState(null);

  useEffect(() => {
    const fetchOrders = async () => {
      setLoading(true);
      try {
        const token = localStorage.getItem('jwt_token');
        const response = await axios.get(
          'http://localhost:8080/api/v1/external/orders',
          {
            params: filters,
            headers: { Authorization: `Bearer ${token}` }
          }
        );
        setOrders(response.data.data);
        setMeta(response.data.meta);
      } catch (error) {
        console.error('Error fetching orders:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchOrders();
  }, [JSON.stringify(filters)]);

  return { orders, loading, meta };
}

// Usage in component
function OrderList() {
  const [filters, setFilters] = useState({
    status: 'PENDING',
    page: 1,
    limit: 10
  });

  const { orders, loading, meta } = useOrders(filters);

  return (
    <div>
      <select 
        value={filters.status} 
        onChange={(e) => setFilters({...filters, status: e.target.value})}
      >
        <option value="">All Status</option>
        <option value="PENDING">Pending</option>
        <option value="CONFIRMED">Confirmed</option>
        <option value="PREPARING">Preparing</option>
        <option value="READY">Ready</option>
      </select>

      {loading ? (
        <p>Loading...</p>
      ) : (
        <ul>
          {orders.map(order => (
            <li key={order.id}>{order.customer_name} - {order.status}</li>
          ))}
        </ul>
      )}

      <p>Page {meta?.page} of {meta?.total_pages}</p>
    </div>
  );
}
```

### Vue.js
```vue
<template>
  <div>
    <div class="filters">
      <select v-model="filters.status">
        <option value="">All Status</option>
        <option value="PENDING">Pending</option>
        <option value="CONFIRMED">Confirmed</option>
        <option value="PREPARING">Preparing</option>
        <option value="READY">Ready</option>
      </select>

      <select v-model="filters.method">
        <option value="">All Methods</option>
        <option value="DINE_IN">Dine In</option>
        <option value="TAKE_AWAY">Take Away</option>
        <option value="DELIVERY">Delivery</option>
      </select>

      <input 
        v-model="filters.customer" 
        placeholder="Search customer..."
      />
    </div>

    <div v-if="loading">Loading...</div>
    <div v-else>
      <div v-for="order in orders" :key="order.id">
        {{ order.customer_name }} - {{ order.status }}
      </div>
    </div>

    <div class="pagination">
      Page {{ meta.page }} of {{ meta.total_pages }}
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      orders: [],
      loading: false,
      meta: {},
      filters: {
        status: '',
        method: '',
        customer: '',
        page: 1,
        limit: 10
      }
    }
  },
  watch: {
    filters: {
      deep: true,
      handler() {
        this.fetchOrders();
      }
    }
  },
  methods: {
    async fetchOrders() {
      this.loading = true;
      try {
        const token = localStorage.getItem('jwt_token');
        const response = await axios.get(
          'http://localhost:8080/api/v1/external/orders',
          {
            params: this.filters,
            headers: { Authorization: `Bearer ${token}` }
          }
        );
        this.orders = response.data.data;
        this.meta = response.data.meta;
      } catch (error) {
        console.error('Error:', error);
      } finally {
        this.loading = false;
      }
    }
  },
  mounted() {
    this.fetchOrders();
  }
}
</script>
```

---

## 🔍 Filter Behavior

### Status Filter
- **Exact match**: Harus sama persis
- **Case-sensitive**: `PENDING` ≠ `pending`
- **Valid values**: `PENDING`, `CONFIRMED`, `PREPARING`, `READY`, `COMPLETED`, `CANCELLED`

### Method Filter
- **Exact match**: Harus sama persis
- **Case-sensitive**: `DINE_IN` ≠ `dine_in`
- **Valid values**: `DINE_IN`, `TAKE_AWAY`, `DELIVERY`

### Customer Search
- **Partial match**: Mengandung substring
- **Case-insensitive**: `john` = `John` = `JOHN`
- **Example**: `john` akan match "John Doe", "johnny", "John Smith"

### Order ID Search
- **Partial match**: Mengandung substring
- **Case-insensitive**: `550e` = `550E`
- **Example**: `550e8400` akan match order dengan ID yang mengandung "550e8400"

---

## 📊 Performance Tips

1. **Use pagination**: Jangan ambil semua data sekaligus
2. **Combine filters**: Lebih spesifik = lebih cepat
3. **Index optimization**: Database sudah punya index untuk status, method, company_id, branch_id
4. **Limit results**: Gunakan limit yang reasonable (10-50 items)

---

## 🧪 Testing with cURL

```bash
# Get all orders
curl -X GET "http://localhost:8080/api/v1/external/orders" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Filter by status
curl -X GET "http://localhost:8080/api/v1/external/orders?status=PENDING" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Filter by method
curl -X GET "http://localhost:8080/api/v1/external/orders?method=DINE_IN" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Search by customer
curl -X GET "http://localhost:8080/api/v1/external/orders?customer=john" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Multiple filters
curl -X GET "http://localhost:8080/api/v1/external/orders?status=PENDING&method=DINE_IN&page=1&limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 📖 Related Documentation

- **Order API**: `ORDER_API.md`
- **WebSocket Filters**: `WEBSOCKET_FILTERS.md`
- **Testing Guide**: `ORDER_TESTING_GUIDE.md`
