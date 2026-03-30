# Test Add Item to Order - Debug Version
# Endpoint: POST /api/v1/external/orders/quick/:id

$baseUrl = "http://localhost:8080/api/v1"

# Login
Write-Host "=== LOGIN ===" -ForegroundColor Cyan
$loginBody = @{
    email = "owner@example.com"
    password = "password123"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
$token = $loginResponse.data.token
Write-Host "Token: $token`n" -ForegroundColor Green

$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

# Step 1: Create Quick Order
Write-Host "=== STEP 1: CREATE QUICK ORDER ===" -ForegroundColor Cyan
$quickOrderBody = @{
    table_number = "A5"
    order_method = "DINE_IN"
    order_items = @(
        @{
            product_id = "PRODUCT_UUID_1"
            quantity = 2
        }
    )
} | ConvertTo-Json -Depth 10

try {
    $createResponse = Invoke-RestMethod -Uri "$baseUrl/external/orders/quick" -Method Post -Headers $headers -Body $quickOrderBody
    $orderId = $createResponse.data.id
    
    Write-Host "Order Created!" -ForegroundColor Green
    Write-Host "Order ID: $orderId"
    Write-Host "Subtotal: $($createResponse.data.subtotal_amount)"
    Write-Host "Tax: $($createResponse.data.tax_amount)"
    Write-Host "Total: $($createResponse.data.total_amount)"
    Write-Host "Items Count: $($createResponse.data.order_items.Count)`n"
} catch {
    Write-Host "Error creating order: $($_.Exception.Message)" -ForegroundColor Red
    exit
}

# Step 2: Get Order by ID (Before Add Item)
Write-Host "=== STEP 2: GET ORDER BY ID (BEFORE) ===" -ForegroundColor Cyan
try {
    $beforeResponse = Invoke-RestMethod -Uri "$baseUrl/external/orders/$orderId" -Method Get -Headers $headers
    Write-Host "Before Add Item:" -ForegroundColor Yellow
    Write-Host "Subtotal: $($beforeResponse.data.subtotal_amount)"
    Write-Host "Tax: $($beforeResponse.data.tax_amount)"
    Write-Host "Total: $($beforeResponse.data.total_amount)"
    Write-Host "Items Count: $($beforeResponse.data.order_items.Count)"
    Write-Host "Items:"
    foreach ($item in $beforeResponse.data.order_items) {
        Write-Host "  - $($item.product_name) x$($item.quantity) = $($item.subtotal)"
    }
    Write-Host ""
} catch {
    Write-Host "Error getting order: $($_.Exception.Message)" -ForegroundColor Red
}

# Step 3: Add Item to Order
Write-Host "=== STEP 3: ADD ITEM TO ORDER ===" -ForegroundColor Cyan
$addItemBody = @{
    product_id = "PRODUCT_UUID_2"
    quantity = 1
    note = "Extra pedas"
} | ConvertTo-Json

try {
    $addResponse = Invoke-RestMethod -Uri "$baseUrl/external/orders/quick/$orderId" -Method Post -Headers $headers -Body $addItemBody
    Write-Host "Item Added!" -ForegroundColor Green
    Write-Host "Subtotal: $($addResponse.data.subtotal_amount)"
    Write-Host "Tax: $($addResponse.data.tax_amount)"
    Write-Host "Total: $($addResponse.data.total_amount)"
    Write-Host "Items Count: $($addResponse.data.order_items.Count)`n"
} catch {
    Write-Host "Error adding item: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host $_.ErrorDetails.Message -ForegroundColor Red
    }
    exit
}

# Step 4: Get Order by ID (After Add Item)
Write-Host "=== STEP 4: GET ORDER BY ID (AFTER) ===" -ForegroundColor Cyan
try {
    $afterResponse = Invoke-RestMethod -Uri "$baseUrl/external/orders/$orderId" -Method Get -Headers $headers
    Write-Host "After Add Item:" -ForegroundColor Yellow
    Write-Host "Subtotal: $($afterResponse.data.subtotal_amount)"
    Write-Host "Tax: $($afterResponse.data.tax_amount)"
    Write-Host "Total: $($afterResponse.data.total_amount)"
    Write-Host "Items Count: $($afterResponse.data.order_items.Count)"
    Write-Host "Items:"
    foreach ($item in $afterResponse.data.order_items) {
        Write-Host "  - $($item.product_name) x$($item.quantity) = $($item.subtotal)"
    }
    Write-Host ""
} catch {
    Write-Host "Error getting order: $($_.Exception.Message)" -ForegroundColor Red
}

# Step 5: Comparison
Write-Host "=== COMPARISON ===" -ForegroundColor Cyan
$subtotalDiff = $afterResponse.data.subtotal_amount - $beforeResponse.data.subtotal_amount
$taxDiff = $afterResponse.data.tax_amount - $beforeResponse.data.tax_amount
$totalDiff = $afterResponse.data.total_amount - $beforeResponse.data.total_amount
$itemsDiff = $afterResponse.data.order_items.Count - $beforeResponse.data.order_items.Count

Write-Host "Subtotal Difference: +$subtotalDiff" -ForegroundColor $(if ($subtotalDiff -gt 0) { "Green" } else { "Red" })
Write-Host "Tax Difference: +$taxDiff" -ForegroundColor $(if ($taxDiff -gt 0) { "Green" } else { "Red" })
Write-Host "Total Difference: +$totalDiff" -ForegroundColor $(if ($totalDiff -gt 0) { "Green" } else { "Red" })
Write-Host "Items Difference: +$itemsDiff" -ForegroundColor $(if ($itemsDiff -gt 0) { "Green" } else { "Red" })

if ($itemsDiff -gt 0 -and $totalDiff -gt 0) {
    Write-Host "`n✅ TEST PASSED! Item added successfully and totals recalculated." -ForegroundColor Green
} else {
    Write-Host "`n❌ TEST FAILED! Item or totals not updated correctly." -ForegroundColor Red
}
