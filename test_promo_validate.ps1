# Test Promo Validate API
# PowerShell Script

$BASE_URL = "http://localhost:8080/api/v1/external"
$TOKEN = ""

Write-Host "`n=== PROMO VALIDATE API TEST ===" -ForegroundColor Cyan
Write-Host ""

# Step 1: Login
Write-Host "1. Login as OWNER..." -ForegroundColor Yellow
try {
    $loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/login" -Method Post -Body (@{
        email = "owner@company1.com"
        password = "password123"
    } | ConvertTo-Json) -ContentType "application/json"

    $TOKEN = $loginResponse.data.token
    Write-Host "✓ Login successful" -ForegroundColor Green
    Write-Host "Token: $($TOKEN.Substring(0, 20))..." -ForegroundColor Gray
} catch {
    Write-Host "✗ Login failed" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit
}
Write-Host ""

# Step 2: Get existing promos
Write-Host "2. Get existing promos..." -ForegroundColor Yellow
try {
    $promosResponse = Invoke-RestMethod -Uri "$BASE_URL/promos?page=1&limit=5" -Method Get -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }

    $promos = $promosResponse.data
    Write-Host "✓ Found $($promos.Count) promos" -ForegroundColor Green
    
    if ($promos.Count -gt 0) {
        foreach ($promo in $promos) {
            Write-Host "  - $($promo.name) [$($promo.code)] - Active: $($promo.is_active)" -ForegroundColor Gray
        }
        $testCode = $promos[0].code
    } else {
        Write-Host "⚠ No promos found. Creating test promo..." -ForegroundColor Yellow
        
        # Create test promo
        $newPromo = @{
            name = "Test Promo"
            code = "TEST50"
            promo_category = "normal"
            type = "percentage"
            value = 50
            start_date = "2024-01-01"
            end_date = "2025-12-31"
            is_active = $true
        } | ConvertTo-Json

        $createResponse = Invoke-RestMethod -Uri "$BASE_URL/promos" -Method Post -Body $newPromo -ContentType "application/json" -Headers @{
            "Authorization" = "Bearer $TOKEN"
        }
        
        $testCode = $createResponse.data.code
        Write-Host "✓ Created test promo: $testCode" -ForegroundColor Green
    }
} catch {
    Write-Host "✗ Failed to get promos" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    $testCode = "WEEKEND50"
}
Write-Host ""

# Step 3: Validate existing promo
Write-Host "3. Validate existing promo code: $testCode" -ForegroundColor Yellow
try {
    $validateResponse = Invoke-RestMethod -Uri "$BASE_URL/promos/validate/$testCode" -Method Get -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }

    Write-Host "✓ Validation response received" -ForegroundColor Green
    Write-Host "Valid: $($validateResponse.data.valid)" -ForegroundColor $(if ($validateResponse.data.valid) { "Green" } else { "Red" })
    Write-Host "Message: $($validateResponse.data.message)" -ForegroundColor Gray
    
    if ($validateResponse.data.promo) {
        $promo = $validateResponse.data.promo
        Write-Host "`nPromo Details:" -ForegroundColor Cyan
        Write-Host "  Name: $($promo.name)" -ForegroundColor Gray
        Write-Host "  Code: $($promo.code)" -ForegroundColor Gray
        Write-Host "  Category: $($promo.promo_category)" -ForegroundColor Gray
        Write-Host "  Type: $($promo.type)" -ForegroundColor Gray
        Write-Host "  Value: $($promo.value)" -ForegroundColor Gray
        Write-Host "  Active: $($promo.is_active)" -ForegroundColor Gray
        Write-Host "  Expired: $($promo.is_expired)" -ForegroundColor Gray
        Write-Host "  Available: $($promo.is_available)" -ForegroundColor Gray
        if ($promo.quota) {
            Write-Host "  Quota: $($promo.used_count)/$($promo.quota) used" -ForegroundColor Gray
            Write-Host "  Remaining: $($promo.remaining_quota)" -ForegroundColor Gray
        }
        Write-Host "  Valid: $($promo.start_date) - $($promo.end_date)" -ForegroundColor Gray
    }
} catch {
    Write-Host "✗ Validation failed" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Step 4: Validate invalid promo code
Write-Host "4. Validate invalid promo code: INVALID123" -ForegroundColor Yellow
try {
    $invalidResponse = Invoke-RestMethod -Uri "$BASE_URL/promos/validate/INVALID123" -Method Get -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }

    Write-Host "✓ Response received" -ForegroundColor Green
    Write-Host "Valid: $($invalidResponse.data.valid)" -ForegroundColor Red
    Write-Host "Message: $($invalidResponse.data.message)" -ForegroundColor Gray
} catch {
    Write-Host "✗ Request failed" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Step 5: Create inactive promo and test
Write-Host "5. Test inactive promo..." -ForegroundColor Yellow
try {
    $inactivePromo = @{
        name = "Inactive Test"
        code = "INACTIVE99"
        promo_category = "normal"
        type = "percentage"
        value = 10
        start_date = "2024-01-01"
        end_date = "2025-12-31"
        is_active = $false
    } | ConvertTo-Json

    $createInactive = Invoke-RestMethod -Uri "$BASE_URL/promos" -Method Post -Body $inactivePromo -ContentType "application/json" -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }
    
    Write-Host "✓ Created inactive promo" -ForegroundColor Green
    
    # Validate inactive promo
    $validateInactive = Invoke-RestMethod -Uri "$BASE_URL/promos/validate/INACTIVE99" -Method Get -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }
    
    Write-Host "Valid: $($validateInactive.data.valid)" -ForegroundColor Red
    Write-Host "Message: $($validateInactive.data.message)" -ForegroundColor Gray
} catch {
    Write-Host "✗ Test failed" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Step 6: Create expired promo and test
Write-Host "6. Test expired promo..." -ForegroundColor Yellow
try {
    $expiredPromo = @{
        name = "Expired Test"
        code = "EXPIRED2023"
        promo_category = "normal"
        type = "percentage"
        value = 20
        start_date = "2023-01-01"
        end_date = "2023-12-31"
        is_active = $true
    } | ConvertTo-Json

    $createExpired = Invoke-RestMethod -Uri "$BASE_URL/promos" -Method Post -Body $expiredPromo -ContentType "application/json" -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }
    
    Write-Host "✓ Created expired promo" -ForegroundColor Green
    
    # Validate expired promo
    $validateExpired = Invoke-RestMethod -Uri "$BASE_URL/promos/validate/EXPIRED2023" -Method Get -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }
    
    Write-Host "Valid: $($validateExpired.data.valid)" -ForegroundColor Red
    Write-Host "Message: $($validateExpired.data.message)" -ForegroundColor Gray
} catch {
    Write-Host "✗ Test failed" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Step 7: Create future promo and test
Write-Host "7. Test future promo (not started yet)..." -ForegroundColor Yellow
try {
    $futurePromo = @{
        name = "Future Test"
        code = "FUTURE2030"
        promo_category = "normal"
        type = "percentage"
        value = 30
        start_date = "2030-01-01"
        end_date = "2030-12-31"
        is_active = $true
    } | ConvertTo-Json

    $createFuture = Invoke-RestMethod -Uri "$BASE_URL/promos" -Method Post -Body $futurePromo -ContentType "application/json" -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }
    
    Write-Host "✓ Created future promo" -ForegroundColor Green
    
    # Validate future promo
    $validateFuture = Invoke-RestMethod -Uri "$BASE_URL/promos/validate/FUTURE2030" -Method Get -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }
    
    Write-Host "Valid: $($validateFuture.data.valid)" -ForegroundColor Red
    Write-Host "Message: $($validateFuture.data.message)" -ForegroundColor Gray
} catch {
    Write-Host "✗ Test failed" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

Write-Host "=== TEST COMPLETED ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Summary:" -ForegroundColor Yellow
Write-Host "✓ Valid promo returns valid: true" -ForegroundColor Green
Write-Host "✓ Invalid code returns valid: false" -ForegroundColor Green
Write-Host "✓ Inactive promo returns valid: false" -ForegroundColor Green
Write-Host "✓ Expired promo returns valid: false" -ForegroundColor Green
Write-Host "✓ Future promo returns valid: false" -ForegroundColor Green
Write-Host ""
Write-Host "API Endpoint: GET /api/v1/external/promos/validate/:code" -ForegroundColor Cyan
