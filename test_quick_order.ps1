# Test Add Item to Order Endpoint
# Endpoint: POST /api/v1/external/orders/quick/:id

$baseUrl = "http://localhost:8080/api/v1"

# Login dulu untuk dapat token
Write-Host "=== LOGIN ===" -ForegroundColor Cyan
$loginBody = @{
    email = "owner@example.com"
    password = "password123"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
$token = $loginResponse.data.token
Write-Host "Token: $token" -ForegroundColor Green

# Headers dengan token
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

# Ganti dengan Order ID yang valid
$orderId = "YOUR_ORDER_UUID_HERE"

Write-Host "`n=== ADD ITEM TO ORDER ===" -ForegroundColor Cyan

# Body: hanya item yang mau ditambah
$addItemBody = @{
    product_id = "PRODUCT_UUID_HERE"
    quantity = 2
    note = "Extra pedas"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/external/orders/quick/$orderId" -Method Post -Headers $headers -Body $addItemBody
    Write-Host "Success!" -ForegroundColor Green
    Write-Host ($response | ConvertTo-Json -Depth 10)
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host $_.ErrorDetails.Message -ForegroundColor Red
    }
}

Write-Host "`n=== NOTES ===" -ForegroundColor Yellow
Write-Host "Endpoint ini untuk menambah item ke order yang sudah ada"
Write-Host "Body hanya berisi:"
Write-Host "- product_id: UUID produk yang mau ditambah"
Write-Host "- quantity: Jumlah"
Write-Host "- note: Catatan (opsional)"
Write-Host ""
Write-Host "Sistem akan otomatis:"
Write-Host "- Validasi produk tersedia"
Write-Host "- Hitung ulang subtotal"
Write-Host "- Hitung ulang diskon (jika ada promo)"
Write-Host "- Hitung ulang pajak"
Write-Host "- Hitung ulang total"
Write-Host "- Broadcast via WebSocket"
