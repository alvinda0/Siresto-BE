# Test Payment API
# PowerShell script untuk testing payment endpoint

$BASE_URL = "http://localhost:8080/api/v1"
$TOKEN = "your-token-here"

Write-Host "=== PAYMENT API TEST ===" -ForegroundColor Cyan
Write-Host ""

# Function untuk login dan get token
function Get-AuthToken {
    Write-Host "1. Login untuk mendapatkan token..." -ForegroundColor Yellow
    
    $loginBody = @{
        email = "owner@example.com"
        password = "password123"
    } | ConvertTo-Json

    $response = Invoke-RestMethod -Uri "$BASE_URL/login" -Method POST -Body $loginBody -ContentType "application/json"
    
    if ($response.status -eq "success") {
        Write-Host "✓ Login berhasil" -ForegroundColor Green
        return $response.data.token
    } else {
        Write-Host "✗ Login gagal" -ForegroundColor Red
        exit 1
    }
}

# Get token
$TOKEN = Get-AuthToken
$headers = @{
    "Authorization" = "Bearer $TOKEN"
    "Content-Type" = "application/json"
}

Write-Host ""
Write-Host "2. Create Order untuk testing payment..." -ForegroundColor Yellow

# Ganti dengan product_id yang valid dari database Anda
$createOrderBody = @{
    table_number = "A1"
    order_method = "DINE_IN"
    customer_name = "Test Customer"
    order_items = @(
        @{
            product_id = "your-product-id-here"
            quantity = 2
            note = "Extra pedas"
        }
    )
} | ConvertTo-Json

try {
    $orderResponse = Invoke-RestMethod -Uri "$BASE_URL/external/orders" -Method POST -Headers $headers -Body $createOrderBody
    
    if ($orderResponse.status -eq "success") {
        $orderId = $orderResponse.data.id
        $totalAmount = $orderResponse.data.total_amount
        
        Write-Host "✓ Order created successfully" -ForegroundColor Green
        Write-Host "  Order ID: $orderId" -ForegroundColor Gray
        Write-Host "  Total Amount: $totalAmount" -ForegroundColor Gray
        Write-Host ""
        
        # Test 1: Payment dengan TUNAI (Cash)
        Write-Host "3. Test Payment dengan TUNAI..." -ForegroundColor Yellow
        
        $paidAmount = [math]::Ceiling($totalAmount / 1000) * 1000 + 5000
        
        $paymentBody = @{
            payment_method = "TUNAI"
            paid_amount = $paidAmount
            payment_note = "Pembayaran tunai"
        } | ConvertTo-Json
        
        $paymentResponse = Invoke-RestMethod -Uri "$BASE_URL/external/orders/$orderId/payment" -Method POST -Headers $headers -Body $paymentBody
        
        if ($paymentResponse.status -eq "success") {
            Write-Host "✓ Payment TUNAI berhasil" -ForegroundColor Green
            Write-Host "  Payment Method: $($paymentResponse.data.payment_method)" -ForegroundColor Gray
            Write-Host "  Payment Status: $($paymentResponse.data.payment_status)" -ForegroundColor Gray
            Write-Host "  Total Amount: $($paymentResponse.data.total_amount)" -ForegroundColor Gray
            Write-Host "  Paid Amount: $($paymentResponse.data.paid_amount)" -ForegroundColor Gray
            Write-Host "  Change Amount: $($paymentResponse.data.change_amount)" -ForegroundColor Gray
            Write-Host "  Paid At: $($paymentResponse.data.paid_at)" -ForegroundColor Gray
        }
        
    } else {
        Write-Host "✗ Failed to create order" -ForegroundColor Red
        Write-Host $orderResponse | ConvertTo-Json -Depth 10
    }
} catch {
    Write-Host "✗ Error: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host $_.ErrorDetails.Message -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "4. Test Payment dengan QRIS..." -ForegroundColor Yellow

# Create another order for QRIS test
$createOrderBody2 = @{
    table_number = "B2"
    order_method = "DINE_IN"
    customer_name = "Test Customer 2"
    order_items = @(
        @{
            product_id = "your-product-id-here"
            quantity = 1
        }
    )
} | ConvertTo-Json

try {
    $orderResponse2 = Invoke-RestMethod -Uri "$BASE_URL/external/orders" -Method POST -Headers $headers -Body $createOrderBody2
    
    if ($orderResponse2.status -eq "success") {
        $orderId2 = $orderResponse2.data.id
        $totalAmount2 = $orderResponse2.data.total_amount
        
        Write-Host "✓ Order created for QRIS test" -ForegroundColor Green
        Write-Host "  Order ID: $orderId2" -ForegroundColor Gray
        Write-Host "  Total Amount: $totalAmount2" -ForegroundColor Gray
        
        # Payment dengan QRIS
        $paymentBody2 = @{
            payment_method = "QRIS"
            paid_amount = $totalAmount2
            payment_note = "Pembayaran via QRIS"
        } | ConvertTo-Json
        
        $paymentResponse2 = Invoke-RestMethod -Uri "$BASE_URL/external/orders/$orderId2/payment" -Method POST -Headers $headers -Body $paymentBody2
        
        if ($paymentResponse2.status -eq "success") {
            Write-Host "✓ Payment QRIS berhasil" -ForegroundColor Green
            Write-Host "  Payment Method: $($paymentResponse2.data.payment_method)" -ForegroundColor Gray
            Write-Host "  Payment Status: $($paymentResponse2.data.payment_status)" -ForegroundColor Gray
            Write-Host "  Total Amount: $($paymentResponse2.data.total_amount)" -ForegroundColor Gray
            Write-Host "  Paid Amount: $($paymentResponse2.data.paid_amount)" -ForegroundColor Gray
        }
    }
} catch {
    Write-Host "✗ Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "5. Test Payment dengan Promo..." -ForegroundColor Yellow

# Create order with promo
$createOrderBody3 = @{
    table_number = "C3"
    order_method = "DINE_IN"
    customer_name = "Test Customer 3"
    promo_code = "DISKON10"
    order_items = @(
        @{
            product_id = "your-product-id-here"
            quantity = 3
        }
    )
} | ConvertTo-Json

try {
    $orderResponse3 = Invoke-RestMethod -Uri "$BASE_URL/external/orders" -Method POST -Headers $headers -Body $createOrderBody3
    
    if ($orderResponse3.status -eq "success") {
        $orderId3 = $orderResponse3.data.id
        $totalAmount3 = $orderResponse3.data.total_amount
        $discountAmount = $orderResponse3.data.discount_amount
        
        Write-Host "✓ Order with promo created" -ForegroundColor Green
        Write-Host "  Order ID: $orderId3" -ForegroundColor Gray
        Write-Host "  Discount Amount: $discountAmount" -ForegroundColor Gray
        Write-Host "  Total Amount: $totalAmount3" -ForegroundColor Gray
        
        # Payment dengan GOPAY
        $paymentBody3 = @{
            payment_method = "GOPAY"
            paid_amount = $totalAmount3
            payment_note = "Pembayaran via GoPay dengan promo"
        } | ConvertTo-Json
        
        $paymentResponse3 = Invoke-RestMethod -Uri "$BASE_URL/external/orders/$orderId3/payment" -Method POST -Headers $headers -Body $paymentBody3
        
        if ($paymentResponse3.status -eq "success") {
            Write-Host "✓ Payment GOPAY dengan promo berhasil" -ForegroundColor Green
            Write-Host "  Payment Method: $($paymentResponse3.data.payment_method)" -ForegroundColor Gray
            Write-Host "  Discount Amount: $($paymentResponse3.data.discount_amount)" -ForegroundColor Gray
            Write-Host "  Total Amount: $($paymentResponse3.data.total_amount)" -ForegroundColor Gray
            
            if ($paymentResponse3.data.promo_details) {
                Write-Host "  Promo: $($paymentResponse3.data.promo_details.promo_name)" -ForegroundColor Gray
                Write-Host "  Promo Code: $($paymentResponse3.data.promo_details.promo_code)" -ForegroundColor Gray
            }
        }
    }
} catch {
    Write-Host "✗ Error: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host $_.ErrorDetails.Message -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "6. Test Payment COMPLIMENTARY..." -ForegroundColor Yellow

# Create order for complimentary
$createOrderBody4 = @{
    table_number = "VIP1"
    order_method = "DINE_IN"
    customer_name = "VIP Customer"
    order_items = @(
        @{
            product_id = "your-product-id-here"
            quantity = 1
        }
    )
} | ConvertTo-Json

try {
    $orderResponse4 = Invoke-RestMethod -Uri "$BASE_URL/external/orders" -Method POST -Headers $headers -Body $createOrderBody4
    
    if ($orderResponse4.status -eq "success") {
        $orderId4 = $orderResponse4.data.id
        
        Write-Host "✓ Order for complimentary created" -ForegroundColor Green
        Write-Host "  Order ID: $orderId4" -ForegroundColor Gray
        
        # Payment COMPLIMENTARY
        $paymentBody4 = @{
            payment_method = "COMPLIMENTARY"
            paid_amount = 0
            payment_note = "Complimentary untuk VIP"
        } | ConvertTo-Json
        
        $paymentResponse4 = Invoke-RestMethod -Uri "$BASE_URL/external/orders/$orderId4/payment" -Method POST -Headers $headers -Body $paymentBody4
        
        if ($paymentResponse4.status -eq "success") {
            Write-Host "✓ Payment COMPLIMENTARY berhasil" -ForegroundColor Green
            Write-Host "  Payment Method: $($paymentResponse4.data.payment_method)" -ForegroundColor Gray
            Write-Host "  Paid Amount: $($paymentResponse4.data.paid_amount)" -ForegroundColor Gray
        }
    }
} catch {
    Write-Host "✗ Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "=== TEST SELESAI ===" -ForegroundColor Cyan
