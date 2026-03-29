# Test Order dengan Promo
# Formula: ((Subtotal - Discount) + Tax Priority 1) + Tax Priority 2 = Total

$BASE_URL = "http://localhost:8080/api/v1"

# Login sebagai CASHIER
Write-Host "`n=== LOGIN CASHIER ===" -ForegroundColor Cyan
$loginResponse = Invoke-RestMethod -Uri "$BASE_URL/auth/login" -Method POST -ContentType "application/json" -Body (@{
    email = "cashier@branch1.com"
    password = "password123"
} | ConvertTo-Json)

$TOKEN = $loginResponse.data.token
$HEADERS = @{
    "Authorization" = "Bearer $TOKEN"
    "Content-Type" = "application/json"
}

Write-Host "Token: $TOKEN" -ForegroundColor Green

# Get products untuk order
Write-Host "`n=== GET PRODUCTS ===" -ForegroundColor Cyan
$productsResponse = Invoke-RestMethod -Uri "$BASE_URL/products?limit=5" -Method GET -Headers $HEADERS
$products = $productsResponse.data
Write-Host "Available products: $($products.Count)" -ForegroundColor Green

if ($products.Count -eq 0) {
    Write-Host "No products available!" -ForegroundColor Red
    exit
}

$product1 = $products[0]
Write-Host "Using product: $($product1.name) - Price: $($product1.price)" -ForegroundColor Yellow

# Get active promos
Write-Host "`n=== GET ACTIVE PROMOS ===" -ForegroundColor Cyan
$promosResponse = Invoke-RestMethod -Uri "$BASE_URL/promos?limit=10" -Method GET -Headers $HEADERS
$promos = $promosResponse.data | Where-Object { $_.is_available -eq $true }

if ($promos.Count -eq 0) {
    Write-Host "No active promos available!" -ForegroundColor Yellow
    Write-Host "Creating a test promo..." -ForegroundColor Yellow
    
    # Create test promo
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
    Write-Host "Created promo: $($promo.name) - Code: $($promo.code)" -ForegroundColor Green
} else {
    $promo = $promos[0]
    Write-Host "Using promo: $($promo.name) - Code: $($promo.code)" -ForegroundColor Green
    Write-Host "  Type: $($promo.type)" -ForegroundColor Gray
    Write-Host "  Value: $($promo.value)" -ForegroundColor Gray
    Write-Host "  Min Transaction: $($promo.min_transaction)" -ForegroundColor Gray
}

# Create order WITHOUT promo first
Write-Host "`n=== CREATE ORDER WITHOUT PROMO ===" -ForegroundColor Cyan
$orderWithoutPromo = @{
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

$responseWithoutPromo = Invoke-RestMethod -Uri "$BASE_URL/orders" -Method POST -Headers $HEADERS -Body ($orderWithoutPromo | ConvertTo-Json -Depth 10)
$orderNoPromo = $responseWithoutPromo.data

Write-Host "`nOrder WITHOUT Promo:" -ForegroundColor Yellow
Write-Host "  Order ID: $($orderNoPromo.id)" -ForegroundColor White
Write-Host "  Subtotal: $($orderNoPromo.subtotal_amount)" -ForegroundColor White
Write-Host "  Discount: $($orderNoPromo.discount_amount)" -ForegroundColor White
Write-Host "  Tax Amount: $($orderNoPromo.tax_amount)" -ForegroundColor White
Write-Host "  Total: $($orderNoPromo.total_amount)" -ForegroundColor White

if ($orderNoPromo.tax_details) {
    Write-Host "`n  Tax Breakdown:" -ForegroundColor Cyan
    foreach ($tax in $orderNoPromo.tax_details) {
        Write-Host "    - $($tax.tax_name) ($($tax.percentage)%) Priority $($tax.priority)" -ForegroundColor Gray
        Write-Host "      Base: $($tax.base_amount) -> Tax: $($tax.tax_amount)" -ForegroundColor Gray
    }
}

# Create order WITH promo
Write-Host "`n=== CREATE ORDER WITH PROMO ===" -ForegroundColor Cyan
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
    $responseWithPromo = Invoke-RestMethod -Uri "$BASE_URL/orders" -Method POST -Headers $HEADERS -Body ($orderWithPromo | ConvertTo-Json -Depth 10)
    $orderPromo = $responseWithPromo.data

    Write-Host "`nOrder WITH Promo:" -ForegroundColor Yellow
    Write-Host "  Order ID: $($orderPromo.id)" -ForegroundColor White
    Write-Host "  Promo Code: $($orderPromo.promo_code)" -ForegroundColor Green
    Write-Host "  Subtotal: $($orderPromo.subtotal_amount)" -ForegroundColor White
    Write-Host "  Discount: -$($orderPromo.discount_amount)" -ForegroundColor Red
    Write-Host "  After Discount: $($orderPromo.subtotal_amount - $orderPromo.discount_amount)" -ForegroundColor White
    Write-Host "  Tax Amount: $($orderPromo.tax_amount)" -ForegroundColor White
    Write-Host "  Total: $($orderPromo.total_amount)" -ForegroundColor White

    if ($orderPromo.promo_details) {
        Write-Host "`n  Promo Details:" -ForegroundColor Cyan
        Write-Host "    Name: $($orderPromo.promo_details.promo_name)" -ForegroundColor Gray
        Write-Host "    Type: $($orderPromo.promo_details.promo_type)" -ForegroundColor Gray
        Write-Host "    Value: $($orderPromo.promo_details.promo_value)" -ForegroundColor Gray
        Write-Host "    Discount: $($orderPromo.promo_details.discount_amount)" -ForegroundColor Gray
    }

    if ($orderPromo.tax_details) {
        Write-Host "`n  Tax Breakdown:" -ForegroundColor Cyan
        foreach ($tax in $orderPromo.tax_details) {
            Write-Host "    - $($tax.tax_name) ($($tax.percentage)%) Priority $($tax.priority)" -ForegroundColor Gray
            Write-Host "      Base: $($tax.base_amount) -> Tax: $($tax.tax_amount)" -ForegroundColor Gray
        }
    }

    # Verify calculation
    Write-Host "`n=== CALCULATION VERIFICATION ===" -ForegroundColor Cyan
    $subtotal = $orderPromo.subtotal_amount
    $discount = $orderPromo.discount_amount
    $afterDiscount = $subtotal - $discount
    $calculatedTotal = $afterDiscount
    
    Write-Host "Formula: ((Subtotal - Discount) + Tax1) + Tax2" -ForegroundColor Yellow
    Write-Host "Step 1: Subtotal = $subtotal" -ForegroundColor White
    Write-Host "Step 2: Discount = $discount" -ForegroundColor White
    Write-Host "Step 3: After Discount = $afterDiscount" -ForegroundColor White
    
    foreach ($tax in $orderPromo.tax_details) {
        Write-Host "Step $($tax.priority + 3): + Tax $($tax.priority) ($($tax.tax_name)) = $($tax.tax_amount)" -ForegroundColor White
        $calculatedTotal += $tax.tax_amount
    }
    
    Write-Host "Calculated Total: $calculatedTotal" -ForegroundColor Green
    Write-Host "Actual Total: $($orderPromo.total_amount)" -ForegroundColor Green
    
    if ([Math]::Abs($calculatedTotal - $orderPromo.total_amount) -lt 0.01) {
        Write-Host "✓ CALCULATION CORRECT!" -ForegroundColor Green
    } else {
        Write-Host "✗ CALCULATION MISMATCH!" -ForegroundColor Red
    }

    # Get order by ID to verify breakdown
    Write-Host "`n=== GET ORDER BY ID ===" -ForegroundColor Cyan
    $getOrderResponse = Invoke-RestMethod -Uri "$BASE_URL/orders/$($orderPromo.id)" -Method GET -Headers $HEADERS
    $retrievedOrder = $getOrderResponse.data
    
    Write-Host "Retrieved Order:" -ForegroundColor Yellow
    Write-Host "  Subtotal: $($retrievedOrder.subtotal_amount)" -ForegroundColor White
    Write-Host "  Discount: $($retrievedOrder.discount_amount)" -ForegroundColor White
    Write-Host "  Tax: $($retrievedOrder.tax_amount)" -ForegroundColor White
    Write-Host "  Total: $($retrievedOrder.total_amount)" -ForegroundColor White
    
    if ($retrievedOrder.promo_details) {
        Write-Host "  Promo: $($retrievedOrder.promo_details.promo_name) ($($retrievedOrder.promo_details.promo_code))" -ForegroundColor Green
    }

} catch {
    Write-Host "Error creating order with promo:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        $errorDetails = $_.ErrorDetails.Message | ConvertFrom-Json
        Write-Host "Details: $($errorDetails.message)" -ForegroundColor Red
    }
}

Write-Host "`n=== TEST COMPLETED ===" -ForegroundColor Cyan
