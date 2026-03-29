# Run Promo Category Migration
# PowerShell Script

Write-Host "=== PROMO CATEGORY MIGRATION ===" -ForegroundColor Cyan
Write-Host ""

Write-Host "This will:" -ForegroundColor Yellow
Write-Host "1. Add promo_category column to promos table" -ForegroundColor Gray
Write-Host "2. Create promo_products table" -ForegroundColor Gray
Write-Host "3. Create promo_bundles table" -ForegroundColor Gray
Write-Host "4. Create necessary indexes" -ForegroundColor Gray
Write-Host ""

$confirm = Read-Host "Continue? (y/n)"
if ($confirm -ne "y") {
    Write-Host "Migration cancelled" -ForegroundColor Red
    exit
}

Write-Host ""
Write-Host "Running migration..." -ForegroundColor Yellow
Write-Host ""

try {
    go run add_promo_category_and_tables.go
    
    Write-Host ""
    Write-Host "=== MIGRATION COMPLETED ===" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "1. Restart your server: go run cmd/server/main.go" -ForegroundColor Gray
    Write-Host "2. Test the new promo categories: .\test_promo_categories.ps1" -ForegroundColor Gray
    Write-Host "3. Read documentation: PROMO_CATEGORIES.md" -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host ""
    Write-Host "=== MIGRATION FAILED ===" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    Write-Host ""
    Write-Host "Please check:" -ForegroundColor Yellow
    Write-Host "1. Database connection in .env file" -ForegroundColor Gray
    Write-Host "2. Database is running" -ForegroundColor Gray
    Write-Host "3. Go dependencies are installed" -ForegroundColor Gray
}
