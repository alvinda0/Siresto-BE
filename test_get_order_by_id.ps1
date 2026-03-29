$BASE_URL = "http://localhost:8080/api/v1"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Testing Get Order By ID with Tax Breakdown" -ForegroundColor Cyan
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

if (-not $TOKEN) {
    Write-Host "❌ Login failed!" -ForegroundColor Red
    exit 1
}

Write-Host "✓ Login successful" -ForegroundColor Green
Write-Host ""

# Get all orders
Write-Host "2. Getting orders list..." -ForegroundColor Yellow
$headers = @{
    Authorization = "Bearer $TOKEN"
}

$orders = Invoke-RestMethod -Uri "$BASE_URL/external/orders?limit=1" -Method Get -Headers $headers

if ($orders.data.Count -eq 0) {
    Write-Host "❌ No orders found! Please create an order first." -ForegroundColor Red
    exit 1
}

$ORDER_ID = $orders.data[0].id

Write-Host "✓ Found order: $ORDER_ID" -ForegroundColor Green
Write-Host ""

# Get order by ID
Write-Host "3. Getting order by ID with tax breakdown..." -ForegroundColor Yellow
$orderDetail = Invoke-RestMethod -Uri "$BASE_URL/external/orders/$ORDER_ID" -Method Get -Headers $headers

Write-Host "✓ Order retrieved successfully!" -ForegroundColor Green
Write-Host ""

# Display order details
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Order Details" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "ID: $($orderDetail.data.id)"
Write-Host "Customer: $($orderDetail.data.customer_name)"
Write-Host "Table: $($orderDetail.data.table_number)"
Write-Host "Status: $($orderDetail.data.status)"
Write-Host ""

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Financial Breakdown" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Subtotal (before tax): Rp $($orderDetail.data.subtotal_amount)" -ForegroundColor White
Write-Host "Tax Amount: Rp $($orderDetail.data.tax_amount)" -ForegroundColor Yellow
Write-Host "Total Amount: Rp $($orderDetail.data.total_amount)" -ForegroundColor Green
Write-Host ""

# Display tax breakdown
if ($orderDetail.data.tax_details -and $orderDetail.data.tax_details.Count -gt 0) {
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Tax Breakdown" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    
    $orderDetail.data.tax_details | ForEach-Object {
        Write-Host ""
        Write-Host "Tax: $($_.tax_name)" -ForegroundColor Yellow
        Write-Host "  Percentage: $($_.percentage)%"
        Write-Host "  Priority: $($_.priority)"
        Write-Host "  Base Amount: Rp $($_.base_amount)"
        Write-Host "  Tax Amount: Rp $($_.tax_amount)" -ForegroundColor Green
    }
    Write-Host ""
} else {
    Write-Host "No taxes applied to this order" -ForegroundColor Gray
    Write-Host ""
}

# Display order items
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Order Items" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
$orderDetail.data.order_items | ForEach-Object {
    Write-Host "$($_.quantity)x $($_.product_name) @ Rp $($_.price) = Rp $($_.subtotal)"
}
Write-Host ""

Write-Host "✓ Test completed!" -ForegroundColor Green
