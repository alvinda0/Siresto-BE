#!/bin/bash

# Test Order dengan Promo
# Formula: ((Subtotal - Discount) + Tax Priority 1) + Tax Priority 2 = Total

BASE_URL="http://localhost:8080/api/v1"

echo -e "\n\033[1;36m=== LOGIN CASHIER ===\033[0m"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "cashier@branch1.com",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
echo -e "\033[1;32mToken: $TOKEN\033[0m"

# Get products
echo -e "\n\033[1;36m=== GET PRODUCTS ===\033[0m"
PRODUCTS_RESPONSE=$(curl -s -X GET "$BASE_URL/products?limit=5" \
  -H "Authorization: Bearer $TOKEN")

PRODUCT_ID=$(echo $PRODUCTS_RESPONSE | jq -r '.data[0].id')
PRODUCT_NAME=$(echo $PRODUCTS_RESPONSE | jq -r '.data[0].name')
PRODUCT_PRICE=$(echo $PRODUCTS_RESPONSE | jq -r '.data[0].price')

echo -e "\033[1;33mUsing product: $PRODUCT_NAME - Price: $PRODUCT_PRICE\033[0m"

# Get active promos
echo -e "\n\033[1;36m=== GET ACTIVE PROMOS ===\033[0m"
PROMOS_RESPONSE=$(curl -s -X GET "$BASE_URL/promos?limit=10" \
  -H "Authorization: Bearer $TOKEN")

PROMO_CODE=$(echo $PROMOS_RESPONSE | jq -r '.data[] | select(.is_available == true) | .code' | head -1)

if [ -z "$PROMO_CODE" ] || [ "$PROMO_CODE" == "null" ]; then
    echo -e "\033[1;33mNo active promos available! Creating test promo...\033[0m"
    
    TODAY=$(date +%Y-%m-%d)
    NEXT_MONTH=$(date -d "+1 month" +%Y-%m-%d)
    
    CREATE_PROMO_RESPONSE=$(curl -s -X POST "$BASE_URL/promos" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"name\": \"Diskon 10%\",
        \"code\": \"DISKON10\",
        \"type\": \"percentage\",
        \"value\": 10,
        \"max_discount\": 50000,
        \"min_transaction\": 50000,
        \"start_date\": \"$TODAY\",
        \"end_date\": \"$NEXT_MONTH\",
        \"is_active\": true
      }")
    
    PROMO_CODE=$(echo $CREATE_PROMO_RESPONSE | jq -r '.data.code')
    echo -e "\033[1;32mCreated promo: $PROMO_CODE\033[0m"
else
    echo -e "\033[1;32mUsing promo: $PROMO_CODE\033[0m"
fi

# Create order WITHOUT promo
echo -e "\n\033[1;36m=== CREATE ORDER WITHOUT PROMO ===\033[0m"
ORDER_NO_PROMO=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"table_number\": \"Table-10\",
    \"customer_name\": \"Test Customer\",
    \"order_method\": \"DINE_IN\",
    \"order_items\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 5
      }
    ]
  }")

echo -e "\n\033[1;33mOrder WITHOUT Promo:\033[0m"
echo $ORDER_NO_PROMO | jq '{
  order_id: .data.id,
  subtotal: .data.subtotal_amount,
  discount: .data.discount_amount,
  tax: .data.tax_amount,
  total: .data.total_amount,
  tax_details: .data.tax_details
}'

# Create order WITH promo
echo -e "\n\033[1;36m=== CREATE ORDER WITH PROMO ===\033[0m"
ORDER_WITH_PROMO=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"table_number\": \"Table-11\",
    \"customer_name\": \"Test Customer with Promo\",
    \"order_method\": \"DINE_IN\",
    \"promo_code\": \"$PROMO_CODE\",
    \"order_items\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 5
      }
    ]
  }")

echo -e "\n\033[1;33mOrder WITH Promo:\033[0m"
echo $ORDER_WITH_PROMO | jq '{
  order_id: .data.id,
  promo_code: .data.promo_code,
  subtotal: .data.subtotal_amount,
  discount: .data.discount_amount,
  after_discount: (.data.subtotal_amount - .data.discount_amount),
  tax: .data.tax_amount,
  total: .data.total_amount,
  promo_details: .data.promo_details,
  tax_details: .data.tax_details
}'

# Verify calculation
echo -e "\n\033[1;36m=== CALCULATION VERIFICATION ===\033[0m"
SUBTOTAL=$(echo $ORDER_WITH_PROMO | jq -r '.data.subtotal_amount')
DISCOUNT=$(echo $ORDER_WITH_PROMO | jq -r '.data.discount_amount')
TAX=$(echo $ORDER_WITH_PROMO | jq -r '.data.tax_amount')
TOTAL=$(echo $ORDER_WITH_PROMO | jq -r '.data.total_amount')

AFTER_DISCOUNT=$(echo "$SUBTOTAL - $DISCOUNT" | bc)
CALCULATED_TOTAL=$(echo "$AFTER_DISCOUNT + $TAX" | bc)

echo -e "\033[1;33mFormula: ((Subtotal - Discount) + Tax1) + Tax2\033[0m"
echo "Subtotal: $SUBTOTAL"
echo "Discount: $DISCOUNT"
echo "After Discount: $AFTER_DISCOUNT"
echo "Tax: $TAX"
echo "Calculated Total: $CALCULATED_TOTAL"
echo "Actual Total: $TOTAL"

if [ "$CALCULATED_TOTAL" == "$TOTAL" ]; then
    echo -e "\033[1;32m✓ CALCULATION CORRECT!\033[0m"
else
    echo -e "\033[1;31m✗ CALCULATION MISMATCH!\033[0m"
fi

# Get order by ID
ORDER_ID=$(echo $ORDER_WITH_PROMO | jq -r '.data.id')
echo -e "\n\033[1;36m=== GET ORDER BY ID ===\033[0m"
GET_ORDER=$(curl -s -X GET "$BASE_URL/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

echo -e "\n\033[1;33mRetrieved Order:\033[0m"
echo $GET_ORDER | jq '{
  subtotal: .data.subtotal_amount,
  discount: .data.discount_amount,
  tax: .data.tax_amount,
  total: .data.total_amount,
  promo: .data.promo_details
}'

echo -e "\n\033[1;36m=== TEST COMPLETED ===\033[0m"
