# Simple Test Order dengan Promo

$BASE_URL = "http://localhost:8080/api/v1"

Write-Host "=== LOGIN CASHIER ===" -ForegroundColor Cyan
$loginResponse = Invoke-RestMethod -Uri "$BASE_URL/auth/login" -Method POST -ContentType "application/json" -Body (@{
    email = "cashier@branch1.com"
    password = "password123"
} | ConvertTo-Json)

$TOKEN = $loginResponse.data.token
$HEADERS = @{
    "Authorization" = "Bearer $TOKEN"
    "Content-Type" = "application/json"
}

Write-Host "Token received" -ForegroundColor Green

# Get products
Write-Host "`nGetting products..." -ForegroundColor Cyan
$productsResponse = Invoke-RestMethod -Uri "$BASE_URL/products?limit=5" -Method GET -Headers $HEADERS
$product1 = $productsResponse.data[0]
Write-Host "Using product: $($product1.name) - Price: $($product1.price)" -ForegroundColor Yellow

# Get promos
Write-Host "`nGetting promos..." -ForegroundColor Cyan
$promosResponse = Invoke-RestMethod -Uri "$BASE_URL/promos?limit=10" -Method GET -Headers $HEADERS
$activePromos = $promosResponse.data | Where-Object { $_.is_available -eq $true }

if ($activePromos.Count -eq 0) {
    Write-Host "Creating test promo..." -ForegroundColor Yellow
    $today = Get-Date -Format "yyyy-MM-dd"
    $nextMonth = (Get-Date).AddMonths(1).ToString("yyyy-MM-dd")
    
    $createPromoResponse = Invoke-RestMethod -Uri "$BASE_URL/promos" -Method POST -Headers $HEADERS -Body (@{
        name = "Diskon 10%"
        code = "DISKON10"
        type = "percentage"
        value = 10
        max_discount = 50000
        min_transaction = 50000
        start_date = $today
        end_date = $nextMonth
        is_active = $true
    } | ConvertTo-Json)
    
    $promo = $createPromoResponse.data
} else {
    $promo = $activePromos[0]
}

Write-Host "Using promo: $($promo.code)" -ForegroundColor Green

# Create order WITHOUT promo
Write-Host "`n=== ORDER WITHOUT PROMO ===" -ForegroundColor Cyan
$orderNoPromo = @{
    table_number = "Table-10"
    customer_name = "Test Customer"
    order_method = "DINE_IN"
    order_items = @(
        @{
            product_id = $product1.id
            quantity = 5
        }
    )
}

$responseNoPromo = Invoke-RestMethod -Uri "$BASE_URL/orders" -Method POST -Headers $HEADERS -Body ($orderNoPromo | ConvertTo-Json -Depth 10)
$orderDataNoPromo = $responseNoPromo.data

Write-Host "Subtotal: $($orderDataNoPromo.subtotal_amount)" -ForegroundColor White
Write-Host "Discount: $($orderDataNoPromo.discount_amount)" -ForegroundColor White
Write-Host "Tax: $($orderDataNoPromo.tax_amount)" -ForegroundColor White
Write-Host "Total: $($orderDataNoPromo.total_amount)" -ForegroundColor White

# Create order WITH promo
Write-Host "`n=== ORDER WITH PROMO ===" -ForegroundColor Cyan
$orderWithPromo = @{
    table_number = "Table-11"
    customer_name = "Test Customer with Promo"
    order_method = "DINE_IN"
    promo_code = $promo.code
    order_items = @(
        @{
            product_id = $product1.id
            quantity = 5
        }
    )
}

try {
    $responsePromo = Invoke-RestMethod -Uri "$BASE_URL/orders" -Method POST -Headers $HEADERS -Body ($orderWithPromo | ConvertTo-Json -Depth 10)
    $orderDataPromo = $responsePromo.data

    Write-Host "Promo Code: $($orderDataPromo.promo_code)" -ForegroundColor Green
    Write-Host "Subtotal: $($orderDataPromo.subtotal_amount)" -ForegroundColor White
    Write-Host "Discount: -$($orderDataPromo.discount_amount)" -ForegroundColor Red
    Write-Host "After Discount: $($orderDataPromo.subtotal_amount - $orderDataPromo.discount_amount)" -ForegroundColor White
    Write-Host "Tax: $($orderDataPromo.tax_amount)" -ForegroundColor White
    Write-Host "Total: $($orderDataPromo.total_amount)" -ForegroundColor White

    if ($orderDataPromo.promo_details) {
        Write-Host "`nPromo Details:" -ForegroundColor Cyan
        Write-Host "  Name: $($orderDataPromo.promo_details.promo_name)" -ForegroundColor Gray
        Write-Host "  Type: $($orderDataPromo.promo_details.promo_type)" -ForegroundColor Gray
        Write-Host "  Value: $($orderDataPromo.promo_details.promo_value)" -ForegroundColor Gray
    }

    if ($orderDataPromo.tax_details) {
        Write-Host "`nTax Breakdown:" -ForegroundColor Cyan
        foreach ($tax in $orderDataPromo.tax_details) {
            Write-Host "  $($tax.tax_name) ($($tax.percentage)%) Priority $($tax.priority)" -ForegroundColor Gray
            Write-Host "    Base: $($tax.base_amount) -> Tax: $($tax.tax_amount)" -ForegroundColor Gray
        }
    }

    # Get order by ID
    Write-Host "`n=== GET ORDER BY ID ===" -ForegroundColor Cyan
    $getOrderResponse = Invoke-RestMethod -Uri "$BASE_URL/orders/$($orderDataPromo.id)" -Method GET -Headers $HEADERS
    $retrievedOrder = $getOrderResponse.data
    
    Write-Host "Retrieved successfully!" -ForegroundColor Green
    Write-Host "  Subtotal: $($retrievedOrder.subtotal_amount)" -ForegroundColor White
    Write-Host "  Discount: $($retrievedOrder.discount_amount)" -ForegroundColor White
    Write-Host "  Tax: $($retrievedOrder.tax_amount)" -ForegroundColor White
    Write-Host "  Total: $($retrievedOrder.total_amount)" -ForegroundColor White

} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nTest completed!" -ForegroundColor Cyan
