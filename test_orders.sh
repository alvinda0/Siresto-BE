#!/bin/bash

# Bash script untuk testing Order API
# Usage: ./test_orders.sh

BASE_URL="http://localhost:8080/api/v1"
TOKEN=""
COMPANY_ID=""
BRANCH_ID=""
PRODUCT_ID=""
ORDER_ID=""

echo "=================================="
echo "   ORDER API TESTING SCRIPT"
echo "=================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

# 1. Login
echo -e "${YELLOW}1. Login...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@example.com",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')
COMPANY_ID=$(echo $LOGIN_RESPONSE | jq -r '.data.user.company_id')
BRANCH_ID=$(echo $LOGIN_RESPONSE | jq -r '.data.user.branch_id')

if [ "$TOKEN" != "null" ] && [ "$TOKEN" != "" ]; then
    echo -e "${GREEN}✓ Login successful${NC}"
    echo -e "${GRAY}  Token: ${TOKEN:0:20}...${NC}"
    echo -e "${GRAY}  Company ID: $COMPANY_ID${NC}"
    echo -e "${GRAY}  Branch ID: $BRANCH_ID${NC}"
else
    echo -e "${RED}✗ Login failed${NC}"
    exit 1
fi

echo ""

# 2. Get Products
echo -e "${YELLOW}2. Get Products...${NC}"
PRODUCTS_RESPONSE=$(curl -s -X GET "$BASE_URL/external/products?limit=1" \
  -H "Authorization: Bearer $TOKEN")

PRODUCT_ID=$(echo $PRODUCTS_RESPONSE | jq -r '.data[0].id')
PRODUCT_NAME=$(echo $PRODUCTS_RESPONSE | jq -r '.data[0].name')
PRODUCT_PRICE=$(echo $PRODUCTS_RESPONSE | jq -r '.data[0].price')

if [ "$PRODUCT_ID" != "null" ] && [ "$PRODUCT_ID" != "" ]; then
    echo -e "${GREEN}✓ Product found${NC}"
    echo -e "${GRAY}  Product ID: $PRODUCT_ID${NC}"
    echo -e "${GRAY}  Product Name: $PRODUCT_NAME${NC}"
    echo -e "${GRAY}  Price: $PRODUCT_PRICE${NC}"
else
    echo -e "${RED}✗ No products found. Please create a product first.${NC}"
    exit 1
fi

echo ""

# 3. Create Order (Authenticated)
echo -e "${YELLOW}3. Create Order (Authenticated)...${NC}"
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/orders" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"customer_name\": \"Sasa\",
    \"customer_phone\": \"08123123123\",
    \"table_number\": \"A1\",
    \"notes\": \"Test Order from Bash\",
    \"referral_code\": \"\",
    \"order_method\": \"DINE_IN\",
    \"promo_code\": \"\",
    \"order_items\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 3,
        \"note\": \"Extra pedas\"
      }
    ]
  }")

ORDER_ID=$(echo $CREATE_RESPONSE | jq -r '.data.id')
ORDER_STATUS=$(echo $CREATE_RESPONSE | jq -r '.data.status')
ORDER_TOTAL=$(echo $CREATE_RESPONSE | jq -r '.data.total_amount')
ORDER_ITEMS=$(echo $CREATE_RESPONSE | jq -r '.data.order_items | length')

if [ "$ORDER_ID" != "null" ] && [ "$ORDER_ID" != "" ]; then
    echo -e "${GREEN}✓ Order created successfully${NC}"
    echo -e "${GRAY}  Order ID: $ORDER_ID${NC}"
    echo -e "${GRAY}  Status: $ORDER_STATUS${NC}"
    echo -e "${GRAY}  Total Amount: $ORDER_TOTAL${NC}"
    echo -e "${GRAY}  Items: $ORDER_ITEMS${NC}"
else
    echo -e "${RED}✗ Failed to create order${NC}"
fi

echo ""

# 4. Get Order by ID
echo -e "${YELLOW}4. Get Order by ID...${NC}"
GET_ORDER_RESPONSE=$(curl -s -X GET "$BASE_URL/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

GET_SUCCESS=$(echo $GET_ORDER_RESPONSE | jq -r '.success')
CUSTOMER_NAME=$(echo $GET_ORDER_RESPONSE | jq -r '.data.customer_name')
TABLE_NUMBER=$(echo $GET_ORDER_RESPONSE | jq -r '.data.table_number')
ORDER_METHOD=$(echo $GET_ORDER_RESPONSE | jq -r '.data.order_method')

if [ "$GET_SUCCESS" == "true" ]; then
    echo -e "${GREEN}✓ Order retrieved successfully${NC}"
    echo -e "${GRAY}  Customer: $CUSTOMER_NAME${NC}"
    echo -e "${GRAY}  Table: $TABLE_NUMBER${NC}"
    echo -e "${GRAY}  Method: $ORDER_METHOD${NC}"
    echo -e "${GRAY}  Status: $ORDER_STATUS${NC}"
else
    echo -e "${RED}✗ Failed to get order${NC}"
fi

echo ""

# 5. Get All Orders
echo -e "${YELLOW}5. Get All Orders...${NC}"
GET_ALL_RESPONSE=$(curl -s -X GET "$BASE_URL/orders?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN")

TOTAL_ORDERS=$(echo $GET_ALL_RESPONSE | jq -r '.meta.total')
CURRENT_PAGE=$(echo $GET_ALL_RESPONSE | jq -r '.meta.page')
TOTAL_PAGES=$(echo $GET_ALL_RESPONSE | jq -r '.meta.total_pages')

if [ "$TOTAL_ORDERS" != "null" ]; then
    echo -e "${GREEN}✓ Orders retrieved successfully${NC}"
    echo -e "${GRAY}  Total Orders: $TOTAL_ORDERS${NC}"
    echo -e "${GRAY}  Current Page: $CURRENT_PAGE${NC}"
    echo -e "${GRAY}  Total Pages: $TOTAL_PAGES${NC}"
else
    echo -e "${RED}✗ Failed to get orders${NC}"
fi

echo ""

# 6. Update Order
echo -e "${YELLOW}6. Update Order...${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Sasa Updated",
    "table_number": "A2",
    "status": "CONFIRMED",
    "notes": "Updated from Bash"
  }')

NEW_STATUS=$(echo $UPDATE_RESPONSE | jq -r '.data.status')
NEW_TABLE=$(echo $UPDATE_RESPONSE | jq -r '.data.table_number')

if [ "$NEW_STATUS" != "null" ]; then
    echo -e "${GREEN}✓ Order updated successfully${NC}"
    echo -e "${GRAY}  New Status: $NEW_STATUS${NC}"
    echo -e "${GRAY}  New Table: $NEW_TABLE${NC}"
else
    echo -e "${RED}✗ Failed to update order${NC}"
fi

echo ""

# 7. Create Public Order (No Auth)
echo -e "${YELLOW}7. Create Public Order (No Authentication)...${NC}"
PUBLIC_RESPONSE=$(curl -s -X POST "$BASE_URL/public/orders" \
  -H "Content-Type: application/json" \
  -d "{
    \"company_id\": \"$COMPANY_ID\",
    \"branch_id\": \"$BRANCH_ID\",
    \"customer_name\": \"Public Customer\",
    \"customer_phone\": \"08199999999\",
    \"table_number\": \"B5\",
    \"notes\": \"Public order test\",
    \"order_method\": \"TAKE_AWAY\",
    \"order_items\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 1,
        \"note\": \"No note\"
      }
    ]
  }")

PUBLIC_ORDER_ID=$(echo $PUBLIC_RESPONSE | jq -r '.data.id')
PUBLIC_METHOD=$(echo $PUBLIC_RESPONSE | jq -r '.data.order_method')
PUBLIC_TOTAL=$(echo $PUBLIC_RESPONSE | jq -r '.data.total_amount')

if [ "$PUBLIC_ORDER_ID" != "null" ] && [ "$PUBLIC_ORDER_ID" != "" ]; then
    echo -e "${GREEN}✓ Public order created successfully${NC}"
    echo -e "${GRAY}  Order ID: $PUBLIC_ORDER_ID${NC}"
    echo -e "${GRAY}  Method: $PUBLIC_METHOD${NC}"
    echo -e "${GRAY}  Total: $PUBLIC_TOTAL${NC}"
else
    echo -e "${RED}✗ Failed to create public order${NC}"
fi

echo ""
echo "=================================="
echo "   TESTING COMPLETED"
echo "=================================="
