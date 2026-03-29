# Tax API Testing Script (PowerShell)
# Usage: .\test_tax.ps1 -Token "your_token_here"

param(
    [Parameter(Mandatory=$true)]
    [string]$Token
)

$BaseUrl = "http://localhost:8080/api/v1/external"
$Headers = @{
    "Authorization" = "Bearer $Token"
    "Content-Type" = "application/json"
}

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Tax API Testing" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# 1. Create PB1 Tax
Write-Host "`n1. Creating PB1 Tax..." -ForegroundColor Yellow
$createPB1Body = @{
    nama_pajak = "PB1"
    tipe_pajak = "pb1"
    presentase = 10.00
    deskripsi = "Pajak Barang dan Jasa 1"
    status = "active"
    prioritas = 1
} | ConvertTo-Json

$createPB1 = Invoke-RestMethod -Uri "$BaseUrl/tax" -Method Post -Headers $Headers -Body $createPB1Body
$createPB1 | ConvertTo-Json -Depth 10
$taxId1 = $createPB1.data.id
Write-Host "Created Tax ID: $taxId1" -ForegroundColor Green

# 2. Create Service Charge
Write-Host "`n2. Creating Service Charge..." -ForegroundColor Yellow
$createSCBody = @{
    nama_pajak = "Service Charge"
    tipe_pajak = "sc"
    presentase = 5.00
    deskripsi = "Biaya layanan"
    status = "active"
    prioritas = 2
} | ConvertTo-Json

$createSC = Invoke-RestMethod -Uri "$BaseUrl/tax" -Method Post -Headers $Headers -Body $createSCBody
$createSC | ConvertTo-Json -Depth 10
$taxId2 = $createSC.data.id
Write-Host "Created Tax ID: $taxId2" -ForegroundColor Green

# 3. Get All Taxes
Write-Host "`n3. Getting All Taxes..." -ForegroundColor Yellow
$allTaxes = Invoke-RestMethod -Uri "$BaseUrl/tax" -Method Get -Headers $Headers
$allTaxes | ConvertTo-Json -Depth 10

# 4. Get Tax by ID
Write-Host "`n4. Getting Tax by ID ($taxId1)..." -ForegroundColor Yellow
$taxById = Invoke-RestMethod -Uri "$BaseUrl/tax/$taxId1" -Method Get -Headers $Headers
$taxById | ConvertTo-Json -Depth 10

# 5. Update Tax
Write-Host "`n5. Updating Tax ($taxId1)..." -ForegroundColor Yellow
$updateBody = @{
    nama_pajak = "PB1 Updated"
    presentase = 11.00
    deskripsi = "Updated description"
} | ConvertTo-Json

$updated = Invoke-RestMethod -Uri "$BaseUrl/tax/$taxId1" -Method Put -Headers $Headers -Body $updateBody
$updated | ConvertTo-Json -Depth 10

# 6. Update Status to Inactive
Write-Host "`n6. Setting Tax to Inactive ($taxId2)..." -ForegroundColor Yellow
$statusBody = @{
    status = "inactive"
} | ConvertTo-Json

$statusUpdate = Invoke-RestMethod -Uri "$BaseUrl/tax/$taxId2" -Method Put -Headers $Headers -Body $statusBody
$statusUpdate | ConvertTo-Json -Depth 10

# 7. Get All Taxes Again
Write-Host "`n7. Getting All Taxes (after updates)..." -ForegroundColor Yellow
$allTaxesAfter = Invoke-RestMethod -Uri "$BaseUrl/tax" -Method Get -Headers $Headers
$allTaxesAfter | ConvertTo-Json -Depth 10

# 8. Test Validation - Invalid tipe_pajak
Write-Host "`n8. Testing Validation - Invalid tipe_pajak..." -ForegroundColor Yellow
try {
    $invalidTypeBody = @{
        nama_pajak = "Invalid Tax"
        tipe_pajak = "invalid"
        presentase = 10.00
    } | ConvertTo-Json
    
    $invalidType = Invoke-RestMethod -Uri "$BaseUrl/tax" -Method Post -Headers $Headers -Body $invalidTypeBody
    $invalidType | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Expected Error: $($_.Exception.Message)" -ForegroundColor Red
}

# 9. Test Validation - Presentase > 100
Write-Host "`n9. Testing Validation - Presentase > 100..." -ForegroundColor Yellow
try {
    $invalidPercentBody = @{
        nama_pajak = "Invalid Tax"
        tipe_pajak = "pb1"
        presentase = 150.00
    } | ConvertTo-Json
    
    $invalidPercent = Invoke-RestMethod -Uri "$BaseUrl/tax" -Method Post -Headers $Headers -Body $invalidPercentBody
    $invalidPercent | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Expected Error: $($_.Exception.Message)" -ForegroundColor Red
}

# 10. Delete Tax
Write-Host "`n10. Deleting Tax ($taxId1)..." -ForegroundColor Yellow
$delete1 = Invoke-RestMethod -Uri "$BaseUrl/tax/$taxId1" -Method Delete -Headers $Headers
$delete1 | ConvertTo-Json -Depth 10

# 11. Verify Deletion
Write-Host "`n11. Verifying Deletion (should return 404)..." -ForegroundColor Yellow
try {
    $verify = Invoke-RestMethod -Uri "$BaseUrl/tax/$taxId1" -Method Get -Headers $Headers
    $verify | ConvertTo-Json -Depth 10
} catch {
    Write-Host "Expected Error (404): $($_.Exception.Message)" -ForegroundColor Red
}

# 12. Delete Second Tax
Write-Host "`n12. Deleting Second Tax ($taxId2)..." -ForegroundColor Yellow
$delete2 = Invoke-RestMethod -Uri "$BaseUrl/tax/$taxId2" -Method Delete -Headers $Headers
$delete2 | ConvertTo-Json -Depth 10

# 13. Final Check - Get All Taxes
Write-Host "`n13. Final Check - Get All Taxes..." -ForegroundColor Yellow
$finalCheck = Invoke-RestMethod -Uri "$BaseUrl/tax" -Method Get -Headers $Headers
$finalCheck | ConvertTo-Json -Depth 10

Write-Host "`n==========================================" -ForegroundColor Cyan
Write-Host "Testing Complete!" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
