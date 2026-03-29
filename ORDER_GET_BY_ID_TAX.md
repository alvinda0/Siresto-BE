# Get Order By ID - Tax Breakdown

## Endpoint

```
GET /api/v1/external/orders/{id}
```

## Headers

```
Authorization: Bearer YOUR_TOKEN
```

## Response

Response dari endpoint Get Order By ID akan menampilkan breakdown pajak lengkap yang sama dengan response Create Order.

### Response Structure

```json
{
  "status": "success",
  "message": "Order retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "company_id": "123e4567-e89b-12d3-a456-426614174000",
    "branch_id": "789e0123-e89b-12d3-a456-426614174000",
    "customer_name": "John Doe",
    "customer_phone": "081234567890",
    "table_number": "A1",
    "order_method": "DINE_IN",
    "status": "PENDING",
    
    "subtotal_amount": 100000,
    "tax_amount": 15500,
    "total_amount": 115500,
    
    "tax_details": [
      {
        "tax_id": "abc12345-e89b-12d3-a456-426614174000",
        "tax_name": "PB1",
        "percentage": 10,
        "priority": 1,
        "base_amount": 100000,
        "tax_amount": 10000
      },
      {
        "tax_id": "def67890-e89b-12d3-a456-426614174000",
        "tax_name": "Service Charge",
        "percentage": 5,
        "priority": 2,
        "base_amount": 110000,
        "tax_amount": 5500
      }
    ],
    
    "order_items": [
      {
        "id": "item-uuid",
        "product_id": "product-uuid",
        "product_name": "Nasi Goreng",
        "quantity": 2,
        "price": 50000,
        "subtotal": 100000,
        "note": ""
      }
    ],
    
    "created_at": "2024-01-15 10:30:00",
    "updated_at": "2024-01-15 10:30:00"
  }
}
```

## Field Descriptions

### Financial Fields

| Field | Type | Description |
|-------|------|-------------|
| `subtotal_amount` | float64 | Total harga semua item sebelum pajak |
| `tax_amount` | float64 | Total semua pajak yang dikenakan |
| `total_amount` | float64 | Total akhir (subtotal + tax) |

### Tax Details Array

| Field | Type | Description |
|-------|------|-------------|
| `tax_id` | uuid | ID pajak |
| `tax_name` | string | Nama pajak (contoh: "PB1", "Service Charge") |
| `percentage` | float64 | Persentase pajak |
| `priority` | int | Prioritas perhitungan (lebih tinggi = dihitung lebih dulu) |
| `base_amount` | float64 | Jumlah yang dikenakan pajak (sudah termasuk pajak sebelumnya) |
| `tax_amount` | float64 | Hasil perhitungan pajak ini |

## Contoh Penggunaan

### cURL

```bash
# Get order by ID
curl -X GET http://localhost:8080/api/v1/external/orders/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### PowerShell

```powershell
$headers = @{
    Authorization = "Bearer YOUR_TOKEN"
}

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/external/orders/550e8400-e29b-41d4-a716-446655440000" `
    -Method Get `
    -Headers $headers

# Display tax breakdown
Write-Host "Subtotal: Rp $($response.data.subtotal_amount)"
Write-Host "Tax: Rp $($response.data.tax_amount)"
Write-Host "Total: Rp $($response.data.total_amount)"
Write-Host ""
Write-Host "Tax Breakdown:"
$response.data.tax_details | ForEach-Object {
    Write-Host "  - $($_.tax_name) ($($_.percentage)%): Rp $($_.tax_amount)"
}
```

### JavaScript/Fetch

```javascript
const orderId = '550e8400-e29b-41d4-a716-446655440000';
const token = 'YOUR_TOKEN';

fetch(`http://localhost:8080/api/v1/external/orders/${orderId}`, {
  headers: {
    'Authorization': `Bearer ${token}`
  }
})
.then(response => response.json())
.then(data => {
  const order = data.data;
  
  console.log('Subtotal:', order.subtotal_amount);
  console.log('Tax:', order.tax_amount);
  console.log('Total:', order.total_amount);
  console.log('\nTax Breakdown:');
  
  order.tax_details.forEach(tax => {
    console.log(`  - ${tax.tax_name} (${tax.percentage}%): ${tax.tax_amount}`);
    console.log(`    Base: ${tax.base_amount}, Priority: ${tax.priority}`);
  });
});
```

## Cara Kerja

Ketika endpoint Get Order By ID dipanggil:

1. Order diambil dari database dengan semua field termasuk `subtotal_amount` dan `tax_amount`
2. Fungsi `toOrderResponse()` dipanggil untuk format response
3. Di dalam `toOrderResponse()`, fungsi `calculateTaxes()` dipanggil untuk menghitung ulang breakdown pajak
4. Breakdown pajak ditampilkan di field `tax_details`

### Mengapa Recalculate?

Tax details tidak disimpan di database, melainkan dihitung ulang setiap kali order ditampilkan. Ini memastikan:

- Jika ada perubahan pada pajak (nama, persentase), breakdown akan update otomatis
- Data tetap konsisten dengan konfigurasi pajak terkini
- Database lebih efisien (tidak perlu menyimpan data redundan)

**Note**: `subtotal_amount`, `tax_amount`, dan `total_amount` tetap disimpan di database dan tidak berubah meskipun konfigurasi pajak berubah. Hanya breakdown detail yang dihitung ulang.

## Testing

### Automated Test

```bash
# Windows
.\test_get_order_by_id.ps1

# Linux/Mac
./test_get_order_by_id.sh
```

### Manual Test

```bash
# 1. Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@example.com",
    "password": "password123"
  }' | jq -r '.data.token')

# 2. Get orders list
ORDER_ID=$(curl -s -X GET "http://localhost:8080/api/v1/external/orders?limit=1" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.data[0].id')

# 3. Get order by ID with tax breakdown
curl -X GET "http://localhost:8080/api/v1/external/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
```

## Contoh Output

```
==========================================
Order Details
==========================================
ID: 550e8400-e29b-41d4-a716-446655440000
Customer: John Doe
Table: A1
Status: PENDING

==========================================
Financial Breakdown
==========================================
Subtotal (before tax): Rp 100000
Tax Amount: Rp 15500
Total Amount: Rp 115500

==========================================
Tax Breakdown
==========================================

Tax: PB1
  Percentage: 10%
  Priority: 1
  Base Amount: Rp 100000
  Tax Amount: Rp 10000

Tax: Service Charge
  Percentage: 5%
  Priority: 2
  Base Amount: Rp 110000
  Tax Amount: Rp 5500

==========================================
Order Items
==========================================
2x Nasi Goreng @ Rp 50000 = Rp 100000
```

## Notes

- Tax breakdown selalu ditampilkan, bahkan jika order dibuat sebelum fitur pajak diimplementasikan
- Jika tidak ada pajak aktif saat order dibuat, `tax_details` akan array kosong
- Breakdown dihitung berdasarkan `subtotal_amount` yang tersimpan di database
- Prioritas pajak menentukan urutan perhitungan (DESC = tertinggi dulu)
