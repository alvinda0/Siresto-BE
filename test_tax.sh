#!/bin/bash

# Tax API Testing Script
# Usage: ./test_tax.sh <token>

if [ -z "$1" ]; then
    echo "Usage: ./test_tax.sh <token>"
    exit 1
fi

TOKEN=$1
BASE_URL="http://localhost:8080/api/v1/external"

echo "=========================================="
echo "Tax API Testing"
echo "=========================================="

# 1. Create PB1 Tax
echo -e "\n1. Creating PB1 Tax..."
CREATE_PB1=$(curl -s -X POST "$BASE_URL/tax" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10.00,
    "deskripsi": "Pajak Barang dan Jasa 1",
    "status": "active",
    "prioritas": 1
  }')

echo "$CREATE_PB1" | jq '.'
TAX_ID_1=$(echo "$CREATE_PB1" | jq -r '.data.id')
echo "Created Tax ID: $TAX_ID_1"

# 2. Create Service Charge
echo -e "\n2. Creating Service Charge..."
CREATE_SC=$(curl -s -X POST "$BASE_URL/tax" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "Service Charge",
    "tipe_pajak": "sc",
    "presentase": 5.00,
    "deskripsi": "Biaya layanan",
    "status": "active",
    "prioritas": 2
  }')

echo "$CREATE_SC" | jq '.'
TAX_ID_2=$(echo "$CREATE_SC" | jq -r '.data.id')
echo "Created Tax ID: $TAX_ID_2"

# 3. Get All Taxes
echo -e "\n3. Getting All Taxes..."
curl -s -X GET "$BASE_URL/tax" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 4. Get Tax by ID
echo -e "\n4. Getting Tax by ID ($TAX_ID_1)..."
curl -s -X GET "$BASE_URL/tax/$TAX_ID_1" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 5. Update Tax
echo -e "\n5. Updating Tax ($TAX_ID_1)..."
curl -s -X PUT "$BASE_URL/tax/$TAX_ID_1" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "PB1 Updated",
    "presentase": 11.00,
    "deskripsi": "Updated description"
  }' | jq '.'

# 6. Update Status to Inactive
echo -e "\n6. Setting Tax to Inactive ($TAX_ID_2)..."
curl -s -X PUT "$BASE_URL/tax/$TAX_ID_2" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "inactive"
  }' | jq '.'

# 7. Get All Taxes Again
echo -e "\n7. Getting All Taxes (after updates)..."
curl -s -X GET "$BASE_URL/tax" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 8. Test Validation - Invalid tipe_pajak
echo -e "\n8. Testing Validation - Invalid tipe_pajak..."
curl -s -X POST "$BASE_URL/tax" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "Invalid Tax",
    "tipe_pajak": "invalid",
    "presentase": 10.00
  }' | jq '.'

# 9. Test Validation - Presentase > 100
echo -e "\n9. Testing Validation - Presentase > 100..."
curl -s -X POST "$BASE_URL/tax" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pajak": "Invalid Tax",
    "tipe_pajak": "pb1",
    "presentase": 150.00
  }' | jq '.'

# 10. Delete Tax
echo -e "\n10. Deleting Tax ($TAX_ID_1)..."
curl -s -X DELETE "$BASE_URL/tax/$TAX_ID_1" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 11. Verify Deletion
echo -e "\n11. Verifying Deletion (should return 404)..."
curl -s -X GET "$BASE_URL/tax/$TAX_ID_1" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 12. Delete Second Tax
echo -e "\n12. Deleting Second Tax ($TAX_ID_2)..."
curl -s -X DELETE "$BASE_URL/tax/$TAX_ID_2" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 13. Final Check - Get All Taxes
echo -e "\n13. Final Check - Get All Taxes..."
curl -s -X GET "$BASE_URL/tax" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo -e "\n=========================================="
echo "Testing Complete!"
echo "=========================================="
