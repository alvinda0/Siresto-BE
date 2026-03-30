$BASE_URL = "http://localhost:8080/api/v1"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Testing Payment Method in Order Response" -ForegroundColor Cyan
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

# Create order
Write-Host "2. Creating order..." -ForegroundColor Yellow
$headers = @{
    Authorization = "Bearer $TOKEN"
}

$orderBody = @{
    table_number = "TEST-01"
    order_method = "DINE_IN"
    customer_name = "Test Customer"
    order_items = @(
        @{
            product_id = "PRODUCT_ID_HERE"
            quantity = 2
        }
    )
} | ConvertTo-Json

try {
    $createResponse = Invoke-RestMethod -Uri "$BASE_URL/external/orders" -Method Post -Body $orderBody -Headers $headers -ContentType "application/json"
    $ORDER_ID = $createResponse.data.id
    Write-Host "✓ Order created: $ORDER_ID" -ForegroundColor Green
    Write-Host ""
    
    # Display initial payment info
    Write-Host "Initial Payment Info:" -ForegroundColor Cyan
    Write-Host "  Payment Method: $($createResponse.data.payment_method)" -ForegroundColor $(if ($createResponse.data.payment_method) { "Green" } else { "Red" })
    Write-Host "  Payment Status: $($createResponse.data.payment_status)"
    Write-Host ""
    
    # Process payment
    Write-Host "3. Processing payment with TUNAI..." -ForegroundColor Yellow
    $paymentBody = @{
        payment_method = "TUNAI"
        paid_amount = 100000
        payment_note = "Test payment"
    } | ConvertTo-Json
    
    $paymentResponse = Invoke-RestMethod -Uri "$BASE_URL/external/orders/$ORDER_ID/payment" -Method Post -Body $paymentBody -Headers $headers -ContentType "application/json"
    Write-Host "✓ Payment processed" -ForegroundColor Green
    Write-Host ""
    
    # Display payment response
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Payment Response" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Payment Method: $($paymentResponse.data.payment_method)" -ForegroundColor Green
    Write-Host "Payment Status: $($paymentResponse.data.payment_status)"
    Write-Host "Total Amount: Rp $($paymentResponse.data.total_amount)"
    Write-Host "Paid Amount: Rp $($paymentResponse.data.paid_amount)"
    Write-Host "Change: Rp $($paymentResponse.data.change_amount)"
    Write-Host ""
    
    # Get order by ID to verify
    Write-Host "4. Verifying order details..." -ForegroundColor Yellow
    $orderDetail = Invoke-RestMethod -Uri "$BASE_URL/external/orders/$ORDER_ID" -Method Get -Headers $headers
    
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Order Details After Payment" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Order ID: $($orderDetail.data.id)"
    Write-Host "Payment Method: $($orderDetail.data.payment_method)" -ForegroundColor Green
    Write-Host "Payment Status: $($orderDetail.data.payment_status)"
    Write-Host "Paid Amount: Rp $($orderDetail.data.paid_amount)"
    Write-Host "Change Amount: Rp $($orderDetail.data.change_amount)"
    Write-Host "Payment Note: $($orderDetail.data.payment_note)"
    Write-Host "Paid At: $($orderDetail.data.paid_at)"
    Write-Host ""
    
    Write-Host "✓ Test completed successfully!" -ForegroundColor Green
    
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Response: $($_.ErrorDetails.Message)" -ForegroundColor Red
}
