Write-Host "Running tax fields migration for orders table..." -ForegroundColor Cyan

go run add_tax_fields_to_orders.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Migration completed successfully!" -ForegroundColor Green
} else {
    Write-Host "Migration failed!" -ForegroundColor Red
    exit 1
}
