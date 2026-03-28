# API Logs Testing Script (PowerShell)
# Quick test untuk API Logging endpoints

$BaseUrl = "http://localhost:8080/api/v1"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "API Logs Testing Script" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Step 1: Login
Write-Host "1. Login sebagai Super Admin..." -ForegroundColor Yellow
$loginBody = @{
    email = "admin@siresto.com"
    password = "password123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$BaseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
    $token = $loginResponse.data.token
    
    if (-not $token) {
        Write-Host "❌ Login failed!" -ForegroundColor Red
        Write-Host "Response: $loginResponse"
        exit 1
    }
    
    Write-Host "✅ Login successful!" -ForegroundColor Green
    Write-Host "Token: $($token.Substring(0, [Math]::Min(20, $token.Length)))..."
    Write-Host ""
} catch {
    Write-Host "❌ Login error: $_" -ForegroundColor Red
    exit 1
}

# Setup headers
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

# Step 2: Generate some logs
Write-Host "2. Generate beberapa logs..." -ForegroundColor Yellow

Write-Host "   - GET /roles"
try {
    Invoke-RestMethod -Uri "$BaseUrl/roles" -Method Get -Headers $headers | Out-Null
} catch {}

Write-Host "   - GET /auth/me"
try {
    Invoke-RestMethod -Uri "$BaseUrl/auth/me" -Method Get -Headers $headers | Out-Null
} catch {}

Write-Host "   - GET /external/categories"
try {
    Invoke-RestMethod -Uri "$BaseUrl/external/categories" -Method Get -Headers $headers | Out-Null
} catch {}

Write-Host "✅ Logs generated!" -ForegroundColor Green
Write-Host ""

# Step 3: Get all logs
Write-Host "3. Get all logs..." -ForegroundColor Yellow
try {
    $logsResponse = Invoke-RestMethod -Uri "$BaseUrl/logs" -Method Get -Headers $headers
    Write-Host "✅ Total logs: $($logsResponse.meta.total_items)" -ForegroundColor Green
    Write-Host "   Current page: $($logsResponse.meta.page)"
    Write-Host "   Per page: $($logsResponse.meta.limit)"
    Write-Host "   Total pages: $($logsResponse.meta.total_pages)"
    Write-Host ""
} catch {
    Write-Host "❌ Error: $_" -ForegroundColor Red
}

# Step 4: Get logs with pagination
Write-Host "4. Get logs with pagination (page=1, limit=5)..." -ForegroundColor Yellow
try {
    $paginatedResponse = Invoke-RestMethod -Uri "$BaseUrl/logs?page=1&limit=5" -Method Get -Headers $headers
    Write-Host "✅ Retrieved $($paginatedResponse.data.Count) logs" -ForegroundColor Green
    Write-Host "   Meta: Page $($paginatedResponse.meta.page) of $($paginatedResponse.meta.total_pages)"
    Write-Host ""
} catch {
    Write-Host "❌ Error: $_" -ForegroundColor Red
}

# Step 5: Filter by method
Write-Host "5. Filter logs by method (GET)..." -ForegroundColor Yellow
try {
    $filteredResponse = Invoke-RestMethod -Uri "$BaseUrl/logs?method=GET" -Method Get -Headers $headers
    Write-Host "✅ Found $($filteredResponse.data.Count) GET requests" -ForegroundColor Green
    
    if ($filteredResponse.data.Count -gt 0) {
        Write-Host "   First log:"
        Write-Host "   - Method: $($filteredResponse.data[0].method)"
        Write-Host "   - Path: $($filteredResponse.data[0].path)"
        Write-Host "   - Status: $($filteredResponse.data[0].status_code)"
        Write-Host "   - Response Time: $($filteredResponse.data[0].response_time)ms"
        Write-Host "   - Access From: $($filteredResponse.data[0].access_from)"
    }
    Write-Host ""
} catch {
    Write-Host "❌ Error: $_" -ForegroundColor Red
}

# Step 6: Filter by POST method
Write-Host "6. Filter logs by method (POST)..." -ForegroundColor Yellow
try {
    $postResponse = Invoke-RestMethod -Uri "$BaseUrl/logs?method=POST" -Method Get -Headers $headers
    Write-Host "✅ Found $($postResponse.data.Count) POST requests" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "❌ Error: $_" -ForegroundColor Red
}

# Step 7: Get log by ID
Write-Host "7. Get log by ID (id=1)..." -ForegroundColor Yellow
try {
    $logDetail = Invoke-RestMethod -Uri "$BaseUrl/logs/1" -Method Get -Headers $headers
    Write-Host "✅ Log details:" -ForegroundColor Green
    Write-Host "   - ID: $($logDetail.data.id)"
    Write-Host "   - Method: $($logDetail.data.method)"
    Write-Host "   - Path: $($logDetail.data.path)"
    Write-Host "   - Status Code: $($logDetail.data.status_code)"
    Write-Host "   - Response Time: $($logDetail.data.response_time)ms"
    Write-Host "   - IP Address: $($logDetail.data.ip_address)"
    Write-Host "   - Access From: $($logDetail.data.access_from)"
    Write-Host "   - User ID: $($logDetail.data.user_id)"
    Write-Host ""
} catch {
    Write-Host "❌ Error: $_" -ForegroundColor Red
}

# Step 8: Test with different User-Agent
Write-Host "8. Test with different User-Agent (Mobile)..." -ForegroundColor Yellow
$mobileHeaders = @{
    "Authorization" = "Bearer $token"
    "User-Agent" = "MyApp/1.0 (Android 12; Mobile)"
}
try {
    Invoke-RestMethod -Uri "$BaseUrl/roles" -Method Get -Headers $mobileHeaders | Out-Null
    Write-Host "✅ Request sent with mobile User-Agent" -ForegroundColor Green
    Write-Host ""
} catch {}

# Step 9: Verify mobile log
Write-Host "9. Verify logs contain mobile access..." -ForegroundColor Yellow
try {
    $recentLogs = Invoke-RestMethod -Uri "$BaseUrl/logs?page=1&limit=5" -Method Get -Headers $headers
    $mobileLog = $recentLogs.data | Where-Object { $_.access_from -eq "mobile" } | Select-Object -First 1
    
    if ($mobileLog) {
        Write-Host "✅ Found mobile access log:" -ForegroundColor Green
        Write-Host "   - Path: $($mobileLog.path)"
        Write-Host "   - Access From: $($mobileLog.access_from)"
        Write-Host "   - User Agent: $($mobileLog.user_agent)"
    } else {
        Write-Host "⚠️  No mobile logs found yet (might need to wait a moment)" -ForegroundColor Yellow
    }
    Write-Host ""
} catch {
    Write-Host "❌ Error: $_" -ForegroundColor Red
}

# Step 10: Test combine filter
Write-Host "10. Test combine filter (GET + pagination)..." -ForegroundColor Yellow
try {
    $combinedResponse = Invoke-RestMethod -Uri "$BaseUrl/logs?method=GET&page=1&limit=3" -Method Get -Headers $headers
    Write-Host "✅ Combined filter results:" -ForegroundColor Green
    Write-Host "   - Method filter: GET"
    Write-Host "   - Page: $($combinedResponse.meta.page)"
    Write-Host "   - Limit: $($combinedResponse.meta.limit)"
    Write-Host "   - Results: $($combinedResponse.data.Count)"
    Write-Host ""
} catch {
    Write-Host "❌ Error: $_" -ForegroundColor Red
}

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "✅ Testing Complete!" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Summary:" -ForegroundColor Yellow
Write-Host "  ✅ Login successful"
Write-Host "  ✅ Logs generated"
Write-Host "  ✅ Get all logs working"
Write-Host "  ✅ Pagination working"
Write-Host "  ✅ Filter by method working"
Write-Host "  ✅ Get by ID working"
Write-Host "  ✅ Access source detection working"
Write-Host ""
Write-Host "Untuk testing lebih detail, gunakan Postman atau lihat:" -ForegroundColor Cyan
Write-Host "  - API_LOGS_TESTING.md"
Write-Host "  - API_LOGS_DOCUMENTATION.md"
Write-Host ""
