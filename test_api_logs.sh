#!/bin/bash

# API Logs Testing Script
# Quick test untuk API Logging endpoints

BASE_URL="http://localhost:8080/api/v1"

echo "=========================================="
echo "API Logs Testing Script"
echo "=========================================="
echo ""

# Step 1: Login
echo "1. Login sebagai Super Admin..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@siresto.com",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Login failed!"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "✅ Login successful!"
echo "Token: ${TOKEN:0:20}..."
echo ""

# Step 2: Generate some logs
echo "2. Generate beberapa logs..."

echo "   - GET /roles"
curl -s -X GET "$BASE_URL/roles" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

echo "   - GET /auth/me"
curl -s -X GET "$BASE_URL/auth/me" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

echo "   - GET /external/categories"
curl -s -X GET "$BASE_URL/external/categories" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

echo "✅ Logs generated!"
echo ""

# Step 3: Get all logs
echo "3. Get all logs..."
LOGS_RESPONSE=$(curl -s -X GET "$BASE_URL/logs" \
  -H "Authorization: Bearer $TOKEN")

echo "$LOGS_RESPONSE" | head -c 500
echo "..."
echo ""

# Step 4: Get logs with pagination
echo "4. Get logs with pagination (page=1, limit=5)..."
PAGINATED_RESPONSE=$(curl -s -X GET "$BASE_URL/logs?page=1&limit=5" \
  -H "Authorization: Bearer $TOKEN")

echo "$PAGINATED_RESPONSE" | head -c 500
echo "..."
echo ""

# Step 5: Filter by method
echo "5. Filter logs by method (GET)..."
FILTERED_RESPONSE=$(curl -s -X GET "$BASE_URL/logs?method=GET" \
  -H "Authorization: Bearer $TOKEN")

echo "$FILTERED_RESPONSE" | head -c 500
echo "..."
echo ""

# Step 6: Get log by ID
echo "6. Get log by ID (id=1)..."
LOG_DETAIL=$(curl -s -X GET "$BASE_URL/logs/1" \
  -H "Authorization: Bearer $TOKEN")

echo "$LOG_DETAIL" | head -c 500
echo "..."
echo ""

# Step 7: Test with different User-Agent
echo "7. Test with different User-Agent (Mobile)..."
curl -s -X GET "$BASE_URL/roles" \
  -H "Authorization: Bearer $TOKEN" \
  -H "User-Agent: MyApp/1.0 (Android 12; Mobile)" > /dev/null

echo "✅ Request sent with mobile User-Agent"
echo ""

# Step 8: Verify mobile log
echo "8. Verify logs contain mobile access..."
MOBILE_LOGS=$(curl -s -X GET "$BASE_URL/logs?page=1&limit=1" \
  -H "Authorization: Bearer $TOKEN")

echo "$MOBILE_LOGS" | head -c 500
echo "..."
echo ""

echo "=========================================="
echo "✅ Testing Complete!"
echo "=========================================="
echo ""
echo "Untuk melihat hasil lengkap, akses:"
echo "  GET $BASE_URL/logs"
echo ""
echo "Atau gunakan Postman/Insomnia untuk testing lebih detail."
