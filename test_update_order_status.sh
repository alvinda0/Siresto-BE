#!/bin/bash

# Test Update Order Status API
# Script untuk testing endpoint PATCH /api/v1/external/orders/:id/status

BASE_URL="http://localhost:8080/api/v1"

echo "=== TEST UPDATE ORDER STATUS API ==="
echo ""

# Step 1: Login
echo "Step 1: Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@example.com",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "✗ Login gagal"
  echo $LOGIN_RESPONSE | jq '.'
  exit 1
fi

echo "✓ Login berhasil"
echo "Token: $TOKEN"
echo ""

# Step 2: Create Quick Order
echo "Step 2: Membuat order baru..."

# Ambil product ID dari environment atau gunakan default
if [ -z "$PRODUCT_ID" ]; then
  echo "✗ PRODUCT_ID tidak ditemukan"
  echo "Gunakan: export PRODUCT_ID='your-product-uuid'"
  exit 1
fi

ORDER_RESPONSE=$(curl -s -X POST "$BASE_URL/external/orders/quick" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"table_number\": \"A1\",
    \"order_method\": \"DINE_IN\",
    \"order_items\": [
      {
        \"product_id\": \"$PRODUCT_ID\",
        \"quantity\": 2,
        \"note\": \"Test order\"
      }
    ]
  }")

ORDER_ID=$(echo $ORDER_RESPONSE | jq -r '.data.id')
CURRENT_STATUS=$(echo $ORDER_RESPONSE | jq -r '.data.status')

if [ "$ORDER_ID" == "null" ] || [ -z "$ORDER_ID" ]; then
  echo "✗ Gagal membuat order"
  echo $ORDER_RESPONSE | jq '.'
  exit 1
fi

echo "✓ Order berhasil dibuat"
echo "Order ID: $ORDER_ID"
echo "Status awal: $CURRENT_STATUS"
echo ""

# Step 3: Update Status ke PROCESSING
echo "Step 3: Update status dari PENDING ke PROCESSING..."
UPDATE_RESPONSE=$(curl -s -X PATCH "$BASE_URL/external/orders/$ORDER_ID/status" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "PROCESSING"
  }')

NEW_STATUS=$(echo $UPDATE_RESPONSE | jq -r '.data.status')

if [ "$NEW_STATUS" == "PROCESSING" ]; then
  echo "✓ Status berhasil diupdate"
  echo "Status baru: $NEW_STATUS"
  echo ""
  echo "Response:"
  echo $UPDATE_RESPONSE | jq '.data'
  echo ""
else
  echo "✗ Gagal update status"
  echo $UPDATE_RESPONSE | jq '.'
  exit 1
fi

# Step 4: Update Status ke READY
echo "Step 4: Update status dari PROCESSING ke READY..."
UPDATE_RESPONSE=$(curl -s -X PATCH "$BASE_URL/external/orders/$ORDER_ID/status" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "READY"
  }')

NEW_STATUS=$(echo $UPDATE_RESPONSE | jq -r '.data.status')

if [ "$NEW_STATUS" == "READY" ]; then
  echo "✓ Status berhasil diupdate"
  echo "Status baru: $NEW_STATUS"
  echo ""
else
  echo "✗ Gagal update status"
  echo $UPDATE_RESPONSE | jq '.'
  exit 1
fi

# Step 5: Update Status ke COMPLETED
echo "Step 5: Update status dari READY ke COMPLETED..."
UPDATE_RESPONSE=$(curl -s -X PATCH "$BASE_URL/external/orders/$ORDER_ID/status" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "COMPLETED"
  }')

NEW_STATUS=$(echo $UPDATE_RESPONSE | jq -r '.data.status')

if [ "$NEW_STATUS" == "COMPLETED" ]; then
  echo "✓ Status berhasil diupdate"
  echo "Status baru: $NEW_STATUS"
  echo ""
else
  echo "✗ Gagal update status"
  echo $UPDATE_RESPONSE | jq '.'
  exit 1
fi

# Step 6: Verify dengan Get Order by ID
echo "Step 6: Verifikasi order..."
GET_RESPONSE=$(curl -s -X GET "$BASE_URL/external/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

FINAL_STATUS=$(echo $GET_RESPONSE | jq -r '.data.status')

echo "✓ Order berhasil diambil"
echo "Status final: $FINAL_STATUS"
echo ""

echo "=== TEST SELESAI ==="
echo ""
echo "Catatan:"
echo "- Pastikan server berjalan di http://localhost:8080"
echo "- Set PRODUCT_ID sebelum menjalankan script:"
echo "  export PRODUCT_ID='your-product-uuid'"
