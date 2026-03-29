#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "=========================================="
echo "Testing Get Order By ID with Tax Breakdown"
echo "=========================================="
echo ""

# Login
echo "1. Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@example.com",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')

if [ "$TOKEN" == "null" ]; then
    echo "❌ Login failed!"
    exit 1
fi

echo "✓ Login successful"
echo ""

# Get all orders
echo "2. Getting orders list..."
ORDERS=$(curl -s -X GET "$BASE_URL/external/orders?limit=1" \
  -H "Authorization: Bearer $TOKEN")

ORDER_ID=$(echo $ORDERS | jq -r '.data[0].id')

if [ "$ORDER_ID" == "null" ]; then
    echo "❌ No orders found! Please create an order first."
    exit 1
fi

echo "✓ Found order: $ORDER_ID"
echo ""

# Get order by ID
echo "3. Getting order by ID with tax breakdown..."
ORDER_DETAIL=$(curl -s -X GET "$BASE_URL/external/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "✓ Order retrieved successfully!"
echo ""

# Display order details
echo "=========================================="
echo "Order Details"
echo "=========================================="
echo "ID: $(echo $ORDER_DETAIL | jq -r '.data.id')"
echo "Customer: $(echo $ORDER_DETAIL | jq -r '.data.customer_name')"
echo "Table: $(echo $ORDER_DETAIL | jq -r '.data.table_number')"
echo "Status: $(echo $ORDER_DETAIL | jq -r '.data.status')"
echo ""

echo "=========================================="
echo "Financial Breakdown"
echo "=========================================="
echo "Subtotal (before tax): Rp $(echo $ORDER_DETAIL | jq -r '.data.subtotal_amount')"
echo "Tax Amount: Rp $(echo $ORDER_DETAIL | jq -r '.data.tax_amount')"
echo "Total Amount: Rp $(echo $ORDER_DETAIL | jq -r '.data.total_amount')"
echo ""

# Display tax breakdown
TAX_COUNT=$(echo $ORDER_DETAIL | jq '.data.tax_details | length')

if [ "$TAX_COUNT" -gt 0 ]; then
    echo "=========================================="
    echo "Tax Breakdown"
    echo "=========================================="
    
    echo $ORDER_DETAIL | jq -r '.data.tax_details[] | 
        "\nTax: \(.tax_name)\n  Percentage: \(.percentage)%\n  Priority: \(.priority)\n  Base Amount: Rp \(.base_amount)\n  Tax Amount: Rp \(.tax_amount)"'
    echo ""
else
    echo "No taxes applied to this order"
    echo ""
fi

# Display order items
echo "=========================================="
echo "Order Items"
echo "=========================================="
echo $ORDER_DETAIL | jq -r '.data.order_items[] | 
    "\(.quantity)x \(.product_name) @ Rp \(.price) = Rp \(.subtotal)"'
echo ""

echo "✓ Test completed!"
