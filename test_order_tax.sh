#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "=========================================="
echo "Testing Order Tax Calculation"
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
COMPANY_ID=$(echo $LOGIN_RESPONSE | jq -r '.data.company_id')
BRANCH_ID=$(echo $LOGIN_RESPONSE | jq -r '.data.branch_id')

if [ "$TOKEN" == "null" ]; then
    echo "❌ Login failed!"
    echo $LOGIN_RESPONSE | jq '.'
    exit 1
fi

echo "✓ Login successful"
echo "  Company ID: $COMPANY_ID"
echo "  Branch ID: $BRANCH_ID"
echo ""

# Get products
echo "2. Getting products..."
PRODUCTS=$(curl -s -X GET "$BASE_URL/external/products?limit=1" \
  -H "Authorization: Bearer $TOKEN")

PRODUCT_ID=$(echo $PRODUCTS | jq -r '.data[0].id')
PRODUCT_NAME=$(echo $PRODUCTS | jq -r '.data[0].name')
PRODUCT_PRICE=$(echo $PRODUCTS | jq -r '.data[0].price')

if [ "$PRODUCT_ID" == "null" ]; then
    echo "❌ No products found!"
    exit 1
fi

echo "✓ Product found: $PRODUCT_NAME (Rp $PRODUCT_PRICE)"
echo ""

# Get active taxes
echo "3. Getting active taxes..."
TAXES=$(curl -s -X GET "$BASE_URL/external/tax" \
  -H "Authorization: Bearer $TOKEN")

echo $TAXES | jq '.data[] | {nama_pajak, presentase, prioritas, status}'
echo ""

# Create order
echo "4. Creating order with tax calculation..."
ORDER_RESPONSE=$(curl -s -X POST "$BASE_URL/external/orders" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"table_number\": \"A1\",
    \"customer_name\": \"Test Customer\",
    \"customer_phone\": \"081234567890\",
    \"order_method\": \"DINE_IN\",
    \"order_items\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 2
      }
    ]
  }")

ORDER_ID=$(echo $ORDER_RESPONSE | jq -r '.data.id')

if [ "$ORDER_ID" == "null" ]; then
    echo "❌ Order creation failed!"
    echo $ORDER_RESPONSE | jq '.'
    exit 1
fi

echo "✓ Order created successfully!"
echo ""
echo "Order Details:"
echo $ORDER_RESPONSE | jq '{
  id: .data.id,
  subtotal_amount: .data.subtotal_amount,
  tax_amount: .data.tax_amount,
  total_amount: .data.total_amount,
  tax_details: .data.tax_details,
  order_items: .data.order_items
}'
echo ""

# Verify calculation
SUBTOTAL=$(echo $ORDER_RESPONSE | jq -r '.data.subtotal_amount')
TAX_AMOUNT=$(echo $ORDER_RESPONSE | jq -r '.data.tax_amount')
TOTAL=$(echo $ORDER_RESPONSE | jq -r '.data.total_amount')

echo "=========================================="
echo "Calculation Summary:"
echo "=========================================="
echo "Subtotal: Rp $SUBTOTAL"
echo "Tax Amount: Rp $TAX_AMOUNT"
echo "Total: Rp $TOTAL"
echo ""
echo "Tax Breakdown:"
echo $ORDER_RESPONSE | jq -r '.data.tax_details[] | "  - \(.tax_name) (\(.percentage)%, Priority \(.priority)): Base Rp \(.base_amount) → Tax Rp \(.tax_amount)"'
echo ""
echo "✓ Test completed!"
