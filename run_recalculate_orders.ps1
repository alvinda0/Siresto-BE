Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Recalculating Existing Orders" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "This will recalculate subtotal, tax, and total for all existing orders" -ForegroundColor Yellow
Write-Host ""

$confirmation = Read-Host "Continue? (y/n)"
if ($confirmation -ne 'y') {
    Write-Host "Cancelled." -ForegroundColor Gray
    exit 0
}

Write-Host ""
Write-Host "Running recalculation..." -ForegroundColor Yellow
go run recalculate_existing_orders.go

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "✓ Recalculation completed successfully!" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "❌ Recalculation failed!" -ForegroundColor Red
    exit 1
}
