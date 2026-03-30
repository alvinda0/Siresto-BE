# PowerShell script untuk testing Transaction Report API
# Pastikan server sudah running di http://localhost:8080

$BASE_URL = "http://localhost:8080/api/v1"
$TOKEN = "YOUR_TOKEN_HERE"  # Ganti dengan token yang valid

Write-Host "=== TRANSACTION REPORT API TESTING ===" -ForegroundColor Cyan
Write-Host ""

# Function untuk membuat request
function Invoke-ApiRequest {
    param (
        [string]$Method,
        [string]$Endpoint,
        [string]$Description
    )
    
    Write-Host "Testing: $Description" -ForegroundColor Yellow
    Write-Host "Endpoint: $Method $Endpoint" -ForegroundColor Gray
    
    try {
        $headers = @{
            "Authorization" = "Bearer $TOKEN"
            "Content-Type" = "application/json"
        }
        
        $response = Invoke-RestMethod -Uri "$BASE_URL$Endpoint" -Method $Method -Headers $headers
        Write-Host "Response:" -ForegroundColor Green
        $response | ConvertTo-Json -Depth 10
        Write-Host ""
    }
    catch {
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host ""
    }
}

# 1. Get All Transactions (Default)
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions" -Description "Get all transactions (default pagination)"

# 2. Get Transactions with Custom Pagination
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions?page=1&limit=5" -Description "Get transactions with limit 5"

# 3. Filter by Date Range
$startDate = "2024-01-01"
$endDate = "2024-12-31"
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions?start_date=$startDate&end_date=$endDate" -Description "Filter by date range"

# 4. Filter by Date and Time Range
$startDate = "2024-01-15"
$endDate = "2024-01-15"
$startTime = "08:00"
$endTime = "22:00"
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions?start_date=$startDate&end_date=$endDate&start_time=$startTime&end_time=$endTime" -Description "Filter by date and time range"

# 5. Search by Customer Name
$search = "John"
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions?search=$search" -Description "Search by customer name"

# 6. Filter by Status
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions?status=completed" -Description "Filter by completed status"

# 7. Filter by Payment Status
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions?payment_status=paid" -Description "Filter by paid status"

# 8. Filter by Payment Method
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions?payment_method=cash" -Description "Filter by cash payment"

# 9. Filter by Order Method
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions?order_method=dine_in" -Description "Filter by dine-in orders"

# 10. Combined Filters
$combinedParams = "start_date=2024-01-01&end_date=2024-12-31&status=completed&payment_status=paid&page=1&limit=10"
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions?$combinedParams" -Description "Combined filters"

# 11. Search with Pagination
$searchParams = "search=customer&page=1&limit=20"
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions?$searchParams" -Description "Search with pagination"

# 12. Today's Transactions
$today = Get-Date -Format "yyyy-MM-dd"
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions?start_date=$today&end_date=$today" -Description "Today's transactions"

# 13. This Month's Transactions
$firstDay = Get-Date -Day 1 -Format "yyyy-MM-dd"
$lastDay = Get-Date -Format "yyyy-MM-dd"
Invoke-ApiRequest -Method "GET" -Endpoint "/external/reports/transactions?start_date=$firstDay&end_date=$lastDay" -Description "This month's transactions"

Write-Host "=== TESTING COMPLETED ===" -ForegroundColor Cyan
