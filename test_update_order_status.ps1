# Test Update Order Status API
# Script untuk testing endpoint PATCH /api/v1/external/orders/:id/status

$baseUrl = "http://localhost:8080/api/v1"

Write-Host "=== TEST UPDATE ORDER STATUS API ===" -ForegroundColor Cyan
Write-Host ""

# Step 1: Login
Write-Host "Step 1: Login..." -ForegroundColor Yellow
$loginBody = @{
    email = "owner@example.com"
    password = "password123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/login" -Method POST -Body $loginBody -ContentType "application/json"
    $token = $loginResponse.data.token
    Write-Host "✓ Login berhasil" -ForegroundColor Green
    Write-Host "Token: $token" -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host "✗ Login gagal: $_" -ForegroundColor Red
    exit
}

# Step 2: Create Quick Order
Write-Host "Step 2: Membuat order baru..." -ForegroundColor Yellow
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

# Ambil product ID dari environment atau gunakan default
$productId = $env:PRODUCT_ID
if (-not $productId) {
    Write-Host "PRODUCT_ID tidak ditemukan, gunakan ID dari database Anda" -ForegroundColor Red
    Write-Host "Contoh: `$env:PRODUCT_ID = 'your-product-uuid'" -ForegroundColor Yellow
    exit
}

$orderBody = @{
    table_number = "A1"
    order_method = "DINE_IN"
    order_items = @(
        @{
            product_id = $productId
            quantity = 2
            note = "Test order"
        }
    )
} | ConvertTo-Json -Depth 10

try {
    $orderResponse = Invoke-RestMethod -Uri "$baseUrl/external/orders/quick" -Method POST -Headers $headers -Body $orderBody
    $orderId = $orderResponse.data.id
    $currentStatus = $orderResponse.data.status
    Write-Host "✓ Order berhasil dibuat" -ForegroundColor Green
    Write-Host "Order ID: $orderId" -ForegroundColor Gray
    Write-Host "Status awal: $currentStatus" -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host "✗ Gagal membuat order: $_" -ForegroundColor Red
    exit
}

# Step 3: Update Status ke PROCESSING
Write-Host "Step 3: Update status dari PENDING ke PROCESSING..." -ForegroundColor Yellow
$statusBody = @{
    status = "PROCESSING"
} | ConvertTo-Json

try {
    $updateResponse = Invoke-RestMethod -Uri "$baseUrl/external/orders/$orderId/status" -Method PATCH -Headers $headers -Body $statusBody
    Write-Host "✓ Status berhasil diupdate" -ForegroundColor Green
    Write-Host "Status baru: $($updateResponse.data.status)" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Response:" -ForegroundColor Cyan
    $updateResponse.data | ConvertTo-Json -Depth 10
    Write-Host ""
} catch {
    Write-Host "✗ Gagal update status: $_" -ForegroundColor Red
    exit
}

# Step 4: Update Status ke READY
Write-Host "Step 4: Update status dari PROCESSING ke READY..." -ForegroundColor Yellow
$statusBody = @{
    status = "READY"
} | ConvertTo-Json

try {
    $updateResponse = Invoke-RestMethod -Uri "$baseUrl/external/orders/$orderId/status" -Method PATCH -Headers $headers -Body $statusBody
    Write-Host "✓ Status berhasil diupdate" -ForegroundColor Green
    Write-Host "Status baru: $($updateResponse.data.status)" -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host "✗ Gagal update status: $_" -ForegroundColor Red
    exit
}

# Step 5: Update Status ke COMPLETED
Write-Host "Step 5: Update status dari READY ke COMPLETED..." -ForegroundColor Yellow
$statusBody = @{
    status = "COMPLETED"
} | ConvertTo-Json

try {
    $updateResponse = Invoke-RestMethod -Uri "$baseUrl/external/orders/$orderId/status" -Method PATCH -Headers $headers -Body $statusBody
    Write-Host "✓ Status berhasil diupdate" -ForegroundColor Green
    Write-Host "Status baru: $($updateResponse.data.status)" -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host "✗ Gagal update status: $_" -ForegroundColor Red
    exit
}

# Step 6: Verify dengan Get Order by ID
Write-Host "Step 6: Verifikasi order..." -ForegroundColor Yellow
try {
    $getResponse = Invoke-RestMethod -Uri "$baseUrl/external/orders/$orderId" -Method GET -Headers $headers
    Write-Host "✓ Order berhasil diambil" -ForegroundColor Green
    Write-Host "Status final: $($getResponse.data.status)" -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host "✗ Gagal mengambil order: $_" -ForegroundColor Red
}

Write-Host "=== TEST SELESAI ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Catatan:" -ForegroundColor Yellow
Write-Host "- Pastikan server berjalan di http://localhost:8080" -ForegroundColor Gray
Write-Host "- Set PRODUCT_ID sebelum menjalankan script:" -ForegroundColor Gray
Write-Host "  `$env:PRODUCT_ID = 'your-product-uuid'" -ForegroundColor Gray
