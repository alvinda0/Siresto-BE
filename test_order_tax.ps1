$BASE_URL = "http://localhost:8080/api/v1"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Testing Order Tax Calculation" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Login
Write-Host "1. Login..." -ForegroundColor Yellow
$loginBody = @{
    email = "owner@example.com"
    password = "password123"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod -Uri "$BASE_URL/login" -Method Post -Body $loginBody -ContentType "application/json"

$TOKEN = $loginResponse.data.token
$COMPANY_ID = $loginResponse.data.company_id
$BRANCH_ID = $loginResponse.data.branch_id

if (-not $TOKEN) {
    Write-Host "❌ Login failed!" -ForegroundColor Red
    $loginResponse | ConvertTo-Json
    exit 1
}

Write-Host "✓ Login successful" -ForegroundColor Green
Write-Host "  Company ID: $COMPANY_ID"
Write-Host "  Branch ID: $BRANCH_ID"
Write-Host ""

# Get products
Write-Host "2. Getting products..." -ForegroundColor Yellow
$headers = @{
    Authorization = "Bearer $TOKEN"
}

$products = Invoke-RestMethod -Uri "$BASE_URL/external/products?limit=1" -Method Get -Headers $headers

$PRODUCT_ID = $products.data[0].id
$PRODUCT_NAME = $products.data[0].name
$PRODUCT_PRICE = $products.data[0].price

if (-not $PRODUCT_ID) {
    Write-Host "❌ No products found!" -ForegroundColor Red
    exit 1
}

Write-Host "✓ Product found: $PRODUCT_NAME (Rp $PRODUCT_PRICE)" -ForegroundColor Green
Write-Host ""

# Get active taxes
Write-Host "3. Getting active taxes..." -ForegroundColor Yellow
$taxes = Invoke-RestMethod -Uri "$BASE_URL/external/tax" -Method Get -Headers $headers

$taxes.data | ForEach-Object {
    Write-Host "  - $($_.nama_pajak): $($_.presentase)% (Priority: $($_.prioritas), Status: $($_.status))"
}
Write-Host ""

# Create order
Write-Host "4. Creating order with tax calculation..." -ForegroundColor Yellow
$orderBody = @{
    table_number = "A1"
    customer_name = "Test Customer"
    customer_phone = "081234567890"
    order_method = "DINE_IN"
    order_items = @(
        @{
            product_id = $PRODUCT_ID
            quantity = 2
        }
    )
} | ConvertTo-Json

$orderResponse = Invoke-RestMethod -Uri "$BASE_URL/external/orders" -Method Post -Body $orderBody -ContentType "application/json" -Headers $headers

$ORDER_ID = $orderResponse.data.id

if (-not $ORDER_ID) {
    Write-Host "❌ Order creation failed!" -ForegroundColor Red
    $orderResponse | ConvertTo-Json
    exit 1
}

Write-Host "✓ Order created successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Order Details:" -ForegroundColor Cyan
Write-Host "  ID: $($orderResponse.data.id)"
Write-Host "  Subtotal: Rp $($orderResponse.data.subtotal_amount)"
Write-Host "  Tax Amount: Rp $($orderResponse.data.tax_amount)"
Write-Host "  Total: Rp $($orderResponse.data.total_amount)"
Write-Host ""

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Calculation Summary:" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Subtotal: Rp $($orderResponse.data.subtotal_amount)" -ForegroundColor White
Write-Host "Tax Amount: Rp $($orderResponse.data.tax_amount)" -ForegroundColor White
Write-Host "Total: Rp $($orderResponse.data.total_amount)" -ForegroundColor White
Write-Host ""
Write-Host "Tax Breakdown:" -ForegroundColor Cyan
$orderResponse.data.tax_details | ForEach-Object {
    Write-Host "  - $($_.tax_name) ($($_.percentage)%, Priority $($_.priority)): Base Rp $($_.base_amount) → Tax Rp $($_.tax_amount)"
}
Write-Host ""
Write-Host "✓ Test completed!" -ForegroundColor Green
