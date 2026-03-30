# Test Dashboard Home API
# Pastikan sudah login dan punya token

$baseUrl = "http://localhost:8080/api/v1"

# Ganti dengan token yang valid
$token = "YOUR_JWT_TOKEN_HERE"

$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Testing Dashboard Home API" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Test 1: Get Home Stats
Write-Host "Test 1: Get Home Stats" -ForegroundColor Yellow
Write-Host "GET $baseUrl/external/home" -ForegroundColor Gray
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/external/home" `
        -Method Get `
        -Headers $headers
    
    Write-Host "Response:" -ForegroundColor Green
    $response | ConvertTo-Json -Depth 10
    Write-Host ""
    
    Write-Host "Summary:" -ForegroundColor Cyan
    Write-Host "  Total Items (7 days): $($response.data.total_items_by_date.Count) days data" -ForegroundColor White
    Write-Host "  Revenue (7 days): $($response.data.revenue_by_date.Count) days data" -ForegroundColor White
    Write-Host "  Best Selling Today: $($response.data.best_selling_daily.Count) items" -ForegroundColor White
    Write-Host "  Best Selling Weekly: $($response.data.best_selling_weekly.Count) items" -ForegroundColor White
    Write-Host "  Best Selling Monthly: $($response.data.best_selling_monthly.Count) items" -ForegroundColor White
    Write-Host "  Complimentary Items: $($response.data.complimentary_items.Count) items" -ForegroundColor White
    Write-Host ""
    
    # Display total items by date
    if ($response.data.total_items_by_date.Count -gt 0) {
        Write-Host "Total Items by Date (Last 7 Days):" -ForegroundColor Cyan
        $response.data.total_items_by_date | ForEach-Object {
            Write-Host "  $($_.date): $($_.value) items" -ForegroundColor White
        }
        Write-Host ""
    }
    
    # Display revenue by date
    if ($response.data.revenue_by_date.Count -gt 0) {
        Write-Host "Revenue by Date (Last 7 Days):" -ForegroundColor Cyan
        $response.data.revenue_by_date | ForEach-Object {
            Write-Host "  $($_.date): Rp $($_.value)" -ForegroundColor White
        }
        Write-Host ""
    }
    
    # Display top 3 best selling today
    if ($response.data.best_selling_daily.Count -gt 0) {
        Write-Host "Top 3 Best Selling Today:" -ForegroundColor Cyan
        $response.data.best_selling_daily | Select-Object -First 3 | ForEach-Object {
            Write-Host "  - $($_.product_name): $($_.total_qty) pcs (Rp $($_.total_amount))" -ForegroundColor White
        }
        Write-Host ""
    }
    
    # Display top 3 best selling weekly
    if ($response.data.best_selling_weekly.Count -gt 0) {
        Write-Host "Top 3 Best Selling Weekly:" -ForegroundColor Cyan
        $response.data.best_selling_weekly | Select-Object -First 3 | ForEach-Object {
            Write-Host "  - $($_.product_name): $($_.total_qty) pcs (Rp $($_.total_amount))" -ForegroundColor White
        }
        Write-Host ""
    }
    
    # Display complimentary items
    if ($response.data.complimentary_items.Count -gt 0) {
        Write-Host "Complimentary Items:" -ForegroundColor Cyan
        $response.data.complimentary_items | ForEach-Object {
            Write-Host "  - $($_.product_name): $($_.total_qty) pcs" -ForegroundColor White
        }
        Write-Host ""
    }
    
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "Details: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Testing Complete" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
