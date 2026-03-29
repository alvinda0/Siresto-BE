Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Verifying Tax Priority Order" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

$BASE_URL = "http://localhost:8080/api/v1"

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

# Get taxes
Write-Host "2. Checking tax configuration..." -ForegroundColor Yellow
$headers = @{
    Authorization = "Bearer $TOKEN"
}

$taxes = Invoke-RestMethod -Uri "$BASE_URL/external/tax" -Method Get -Headers $headers

Write-Host "Active Taxes:" -ForegroundColor Cyan
$taxes.data | Where-Object { $_.status -eq "active" } | Sort-Object prioritas | ForEach-Object {
    Write-Host "  Priority $($_.prioritas): $($_.nama_pajak) ($($_.presentase)%)" -ForegroundColor White
}
Write-Host ""

# Verify priority order
$activeTaxes = $taxes.data | Where-Object { $_.status -eq "active" } | Sort-Object prioritas
if ($activeTaxes.Count -ge 2) {
    $firstTax = $activeTaxes[0]
    $secondTax = $activeTaxes[1]
    
    if ($firstTax.prioritas -lt $secondTax.prioritas) {
        Write-Host "✓ Priority order is correct!" -ForegroundColor Green
        Write-Host "  $($firstTax.nama_pajak) (priority $($firstTax.prioritas)) will be calculated FIRST" -ForegroundColor Green
        Write-Host "  $($secondTax.nama_pajak) (priority $($secondTax.prioritas)) will be calculated SECOND" -ForegroundColor Green
    } else {
        Write-Host "⚠ Warning: Priority order might be incorrect" -ForegroundColor Yellow
    }
} else {
    Write-Host "⚠ Less than 2 active taxes found" -ForegroundColor Yellow
}
Write-Host ""

# Get an order to verify
Write-Host "3. Checking order tax calculation..." -ForegroundColor Yellow
$orders = Invoke-RestMethod -Uri "$BASE_URL/external/orders?limit=1" -Method Get -Headers $headers

if ($orders.data.Count -gt 0) {
    $order = $orders.data[0]
    
    Write-Host "Order ID: $($order.id)" -ForegroundColor Cyan
    Write-Host "Subtotal: Rp $($order.subtotal_amount)" -ForegroundColor White
    Write-Host "Tax: Rp $($order.tax_amount)" -ForegroundColor Yellow
    Write-Host "Total: Rp $($order.total_amount)" -ForegroundColor Green
    Write-Host ""
    
    if ($order.tax_details -and $order.tax_details.Count -gt 0) {
        Write-Host "Tax Calculation Order:" -ForegroundColor Cyan
        $order.tax_details | ForEach-Object {
            Write-Host "  $($_.priority). $($_.tax_name) ($($_.percentage)%)" -ForegroundColor White
            Write-Host "     Base: Rp $($_.base_amount) → Tax: Rp $($_.tax_amount)" -ForegroundColor Gray
        }
        Write-Host ""
        
        # Verify order
        $sortedTaxDetails = $order.tax_details | Sort-Object priority
        $isCorrectOrder = $true
        for ($i = 0; $i -lt $sortedTaxDetails.Count - 1; $i++) {
            if ($sortedTaxDetails[$i].priority -gt $sortedTaxDetails[$i + 1].priority) {
                $isCorrectOrder = $false
                break
            }
        }
        
        if ($isCorrectOrder) {
            Write-Host "✓ Tax calculation order is correct!" -ForegroundColor Green
        } else {
            Write-Host "❌ Tax calculation order is incorrect!" -ForegroundColor Red
        }
    } else {
        Write-Host "⚠ No tax details found in order" -ForegroundColor Yellow
    }
} else {
    Write-Host "⚠ No orders found" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Verification Complete!" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
