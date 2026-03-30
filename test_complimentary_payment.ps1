$BASE_URL = "http://localhost:8080/api/v1"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Testing COMPLIMENTARY Payment (Free Order)" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Login
Write-Host "1. Login..." -ForegroundColor Yellow
$loginBody = @{
    email = "owner@example.com"
    password = "password123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$BASE_URL/login" -Method Post -Body $loginBody -ContentType "application/json"
    $TOKEN = $loginResponse.data.token
    
    if (-not $TOKEN) {
        Write-Host "❌ Login failed!" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "✓ Login successful" -ForegroundColor Green
    Write-Host ""
    
    # Get products
    Write-Host "2. Getting products..." -ForegroundColor Yellow
    $headers = @{
        Authorization = "Bearer $TOKEN"
    }
    
    $products = Invoke-RestMethod -Uri "$BASE_URL/external/products?limit=1" -Method Get -Headers $headers
    
    if ($products.data.Count -eq 0) {
        Write-Host "❌ No products found!" -ForegroundColor Red
        exit 1
    }
    
    $PRODUCT_ID = $products.data[0].id
    $PRODUCT_NAME = $products.data[0].name
    $PRODUCT_PRICE = $products.data[0].price
    
    Write-Host "✓ Using product: $PRODUCT_NAME (Rp $PRODUCT_PRICE)" -ForegroundColor Green
    Write-Host ""
    
    # Create order
    Write-Host "3. Creating order..." -ForegroundColor Yellow
    $orderBody = @{
        table_number = "COMP-01"
        order_method = "DINE_IN"
        customer_name = "VIP Guest"
        notes = "Complimentary order for special guest"
        order_items = @(
            @{
                product_id = $PRODUCT_ID
                quantity = 2
            }
        )
    } | ConvertTo-Json
    
    $createResponse = Invoke-RestMethod -Uri "$BASE_URL/external/orders" -Method Post -Body $orderBody -Headers $headers -ContentType "application/json"
    $ORDER_ID = $createResponse.data.id
    
    Write-Host "✓ Order created: $ORDER_ID" -ForegroundColor Green
    Write-Host ""
    
    # Display order before payment
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Order Before Payment" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Subtotal: Rp $($createResponse.data.subtotal_amount)" -ForegroundColor White
    Write-Host "Discount: Rp $($createResponse.data.discount_amount)" -ForegroundColor Yellow
    Write-Host "Tax: Rp $($createResponse.data.tax_amount)" -ForegroundColor Yellow
    Write-Host "Total: Rp $($createResponse.data.total_amount)" -ForegroundColor Green
    Write-Host "Payment Status: $($createResponse.data.payment_status)" -ForegroundColor Gray
    Write-Host ""
    
    # Process COMPLIMENTARY payment
    Write-Host "4. Processing COMPLIMENTARY payment..." -ForegroundColor Yellow
    $paymentBody = @{
        payment_method = "COMPLIMENTARY"
        paid_amount = 0
        payment_note = "Complimentary for VIP guest"
    } | ConvertTo-Json
    
    $paymentResponse = Invoke-RestMethod -Uri "$BASE_URL/external/orders/$ORDER_ID/payment" -Method Post -Body $paymentBody -Headers $headers -ContentType "application/json"
    
    Write-Host "✓ Payment processed" -ForegroundColor Green
    Write-Host ""
    
    # Display payment response
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Payment Response (COMPLIMENTARY)" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Payment Method: $($paymentResponse.data.payment_method)" -ForegroundColor Magenta
    Write-Host "Payment Status: $($paymentResponse.data.payment_status)" -ForegroundColor Green
    Write-Host ""
    Write-Host "Subtotal: Rp $($paymentResponse.data.subtotal_amount)" -ForegroundColor White
    Write-Host "Discount: Rp $($paymentResponse.data.discount_amount)" -ForegroundColor Yellow
    Write-Host "Tax: Rp $($paymentResponse.data.tax_amount)" -ForegroundColor Yellow
    Write-Host "Total: Rp $($paymentResponse.data.total_amount)" -ForegroundColor $(if ($paymentResponse.data.total_amount -eq 0) { "Green" } else { "Red" })
    Write-Host "Paid: Rp $($paymentResponse.data.paid_amount)" -ForegroundColor $(if ($paymentResponse.data.paid_amount -eq 0) { "Green" } else { "Red" })
    Write-Host "Change: Rp $($paymentResponse.data.change_amount)" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Payment Note: $($paymentResponse.data.payment_note)" -ForegroundColor Cyan
    Write-Host "Paid At: $($paymentResponse.data.paid_at)" -ForegroundColor Gray
    Write-Host ""
    
    # Verify amounts
    if ($paymentResponse.data.total_amount -eq 0 -and $paymentResponse.data.paid_amount -eq 0) {
        Write-Host "✓ COMPLIMENTARY payment verified: Total = 0, Paid = 0" -ForegroundColor Green
    } else {
        Write-Host "❌ ERROR: COMPLIMENTARY should make total = 0 and paid = 0!" -ForegroundColor Red
        Write-Host "   Current Total: $($paymentResponse.data.total_amount)" -ForegroundColor Red
        Write-Host "   Current Paid: $($paymentResponse.data.paid_amount)" -ForegroundColor Red
    }
    Write-Host ""
    
    # Get order by ID to verify
    Write-Host "5. Verifying order details..." -ForegroundColor Yellow
    $orderDetail = Invoke-RestMethod -Uri "$BASE_URL/external/orders/$ORDER_ID" -Method Get -Headers $headers
    
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Final Order Details" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Order ID: $($orderDetail.data.id)"
    Write-Host "Customer: $($orderDetail.data.customer_name)"
    Write-Host "Table: $($orderDetail.data.table_number)"
    Write-Host "Status: $($orderDetail.data.status)" -ForegroundColor Green
    Write-Host ""
    Write-Host "Payment Method: $($orderDetail.data.payment_method)" -ForegroundColor Magenta
    Write-Host "Payment Status: $($orderDetail.data.payment_status)" -ForegroundColor Green
    Write-Host ""
    Write-Host "Subtotal: Rp $($orderDetail.data.subtotal_amount)"
    Write-Host "Discount: Rp $($orderDetail.data.discount_amount)"
    Write-Host "Tax: Rp $($orderDetail.data.tax_amount)"
    Write-Host "Total: Rp $($orderDetail.data.total_amount)" -ForegroundColor $(if ($orderDetail.data.total_amount -eq 0) { "Green" } else { "Red" })
    Write-Host "Paid: Rp $($orderDetail.data.paid_amount)" -ForegroundColor $(if ($orderDetail.data.paid_amount -eq 0) { "Green" } else { "Red" })
    Write-Host ""
    
    Write-Host "✓ Test completed successfully!" -ForegroundColor Green
    
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "Response: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}
