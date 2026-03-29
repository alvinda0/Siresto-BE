# Test Promo Categories - Normal, Product, Bundle
# PowerShell Script

$BASE_URL = "http://localhost:8080/api"
$TOKEN = ""

Write-Host "=== PROMO CATEGORIES TEST ===" -ForegroundColor Cyan
Write-Host ""

# Step 1: Login
Write-Host "1. Login as OWNER..." -ForegroundColor Yellow
$loginResponse = Invoke-RestMethod -Uri "$BASE_URL/login" -Method Post -Body (@{
    email = "owner@company1.com"
    password = "password123"
} | ConvertTo-Json) -ContentType "application/json"

$TOKEN = $loginResponse.data.token
Write-Host "✓ Login successful" -ForegroundColor Green
Write-Host "Token: $TOKEN" -ForegroundColor Gray
Write-Host ""

# Step 2: Get Products
Write-Host "2. Get Products..." -ForegroundColor Yellow
$productsResponse = Invoke-RestMethod -Uri "$BASE_URL/products?page=1&limit=10" -Method Get -Headers @{
    "Authorization" = "Bearer $TOKEN"
}

$products = $productsResponse.data
Write-Host "✓ Found $($products.Count) products" -ForegroundColor Green

if ($products.Count -lt 3) {
    Write-Host "⚠ Need at least 3 products for testing. Please create more products first." -ForegroundColor Red
    exit
}

$product1 = $products[0].id
$product2 = $products[1].id
$product3 = $products[2].id

Write-Host "Product 1: $($products[0].name) ($product1)" -ForegroundColor Gray
Write-Host "Product 2: $($products[1].name) ($product2)" -ForegroundColor Gray
Write-Host "Product 3: $($products[2].name) ($product3)" -ForegroundColor Gray
Write-Host ""

# Step 3: Create Promo Normal
Write-Host "3. Create Promo Normal..." -ForegroundColor Yellow
$promoNormal = @{
    name = "Diskon Akhir Tahun"
    code = "NEWYEAR2024"
    promo_category = "normal"
    type = "percentage"
    value = 15
    max_discount = 100000
    min_transaction = 200000
    start_date = "2024-12-01"
    end_date = "2024-12-31"
    is_active = $true
} | ConvertTo-Json

try {
    $normalResponse = Invoke-RestMethod -Uri "$BASE_URL/promos" -Method Post -Body $promoNormal -ContentType "application/json" -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }
    Write-Host "✓ Promo Normal created successfully" -ForegroundColor Green
    Write-Host "ID: $($normalResponse.data.id)" -ForegroundColor Gray
    Write-Host "Category: $($normalResponse.data.promo_category)" -ForegroundColor Gray
    $normalPromoId = $normalResponse.data.id
} catch {
    Write-Host "✗ Failed to create Promo Normal" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Step 4: Create Promo Product
Write-Host "4. Create Promo Product..." -ForegroundColor Yellow
$promoProduct = @{
    name = "Diskon Produk Pilihan"
    code = "PRODUCT50"
    promo_category = "product"
    type = "percentage"
    value = 50
    max_discount = 500000
    product_ids = @($product1, $product2)
    start_date = "2024-12-01"
    end_date = "2024-12-31"
    is_active = $true
} | ConvertTo-Json

try {
    $productResponse = Invoke-RestMethod -Uri "$BASE_URL/promos" -Method Post -Body $promoProduct -ContentType "application/json" -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }
    Write-Host "✓ Promo Product created successfully" -ForegroundColor Green
    Write-Host "ID: $($productResponse.data.id)" -ForegroundColor Gray
    Write-Host "Category: $($productResponse.data.promo_category)" -ForegroundColor Gray
    Write-Host "Products: $($productResponse.data.products.Count) items" -ForegroundColor Gray
    foreach ($p in $productResponse.data.products) {
        Write-Host "  - $($p.product_name) ($($p.product_sku))" -ForegroundColor Gray
    }
    $productPromoId = $productResponse.data.id
} catch {
    Write-Host "✗ Failed to create Promo Product" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Step 5: Create Promo Bundle
Write-Host "5. Create Promo Bundle..." -ForegroundColor Yellow
$promoBundle = @{
    name = "Paket Hemat Bundle"
    code = "BUNDLE100"
    promo_category = "bundle"
    type = "fixed"
    value = 500000
    bundle_items = @(
        @{
            product_id = $product1
            quantity = 1
        },
        @{
            product_id = $product2
            quantity = 2
        },
        @{
            product_id = $product3
            quantity = 1
        }
    )
    start_date = "2024-12-01"
    end_date = "2024-12-31"
    is_active = $true
} | ConvertTo-Json -Depth 10

try {
    $bundleResponse = Invoke-RestMethod -Uri "$BASE_URL/promos" -Method Post -Body $promoBundle -ContentType "application/json" -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }
    Write-Host "✓ Promo Bundle created successfully" -ForegroundColor Green
    Write-Host "ID: $($bundleResponse.data.id)" -ForegroundColor Gray
    Write-Host "Category: $($bundleResponse.data.promo_category)" -ForegroundColor Gray
    Write-Host "Bundle Items: $($bundleResponse.data.bundle_items.Count) items" -ForegroundColor Gray
    foreach ($b in $bundleResponse.data.bundle_items) {
        Write-Host "  - $($b.product_name) x$($b.quantity)" -ForegroundColor Gray
    }
    $bundlePromoId = $bundleResponse.data.id
} catch {
    Write-Host "✗ Failed to create Promo Bundle" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Step 6: Get All Promos
Write-Host "6. Get All Promos..." -ForegroundColor Yellow
try {
    $allPromosResponse = Invoke-RestMethod -Uri "$BASE_URL/promos?page=1&limit=10" -Method Get -Headers @{
        "Authorization" = "Bearer $TOKEN"
    }
    Write-Host "✓ Retrieved $($allPromosResponse.data.Count) promos" -ForegroundColor Green
    foreach ($promo in $allPromosResponse.data) {
        Write-Host "  - $($promo.name) [$($promo.promo_category)]" -ForegroundColor Gray
    }
} catch {
    Write-Host "✗ Failed to get promos" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# Step 7: Get Promo Normal Detail
if ($normalPromoId) {
    Write-Host "7. Get Promo Normal Detail..." -ForegroundColor Yellow
    try {
        $normalDetailResponse = Invoke-RestMethod -Uri "$BASE_URL/promos/$normalPromoId" -Method Get -Headers @{
            "Authorization" = "Bearer $TOKEN"
        }
        Write-Host "✓ Promo Normal Detail:" -ForegroundColor Green
        Write-Host ($normalDetailResponse.data | ConvertTo-Json -Depth 10) -ForegroundColor Gray
    } catch {
        Write-Host "✗ Failed to get promo detail" -ForegroundColor Red
    }
    Write-Host ""
}

# Step 8: Get Promo Product Detail
if ($productPromoId) {
    Write-Host "8. Get Promo Product Detail..." -ForegroundColor Yellow
    try {
        $productDetailResponse = Invoke-RestMethod -Uri "$BASE_URL/promos/$productPromoId" -Method Get -Headers @{
            "Authorization" = "Bearer $TOKEN"
        }
        Write-Host "✓ Promo Product Detail:" -ForegroundColor Green
        Write-Host ($productDetailResponse.data | ConvertTo-Json -Depth 10) -ForegroundColor Gray
    } catch {
        Write-Host "✗ Failed to get promo detail" -ForegroundColor Red
    }
    Write-Host ""
}

# Step 9: Get Promo Bundle Detail
if ($bundlePromoId) {
    Write-Host "9. Get Promo Bundle Detail..." -ForegroundColor Yellow
    try {
        $bundleDetailResponse = Invoke-RestMethod -Uri "$BASE_URL/promos/$bundlePromoId" -Method Get -Headers @{
            "Authorization" = "Bearer $TOKEN"
        }
        Write-Host "✓ Promo Bundle Detail:" -ForegroundColor Green
        Write-Host ($bundleDetailResponse.data | ConvertTo-Json -Depth 10) -ForegroundColor Gray
    } catch {
        Write-Host "✗ Failed to get promo detail" -ForegroundColor Red
    }
    Write-Host ""
}

# Step 10: Update Promo Product (change products)
if ($productPromoId) {
    Write-Host "10. Update Promo Product (change products)..." -ForegroundColor Yellow
    $updateProduct = @{
        product_ids = @($product2, $product3)
    } | ConvertTo-Json

    try {
        $updateResponse = Invoke-RestMethod -Uri "$BASE_URL/promos/$productPromoId" -Method Put -Body $updateProduct -ContentType "application/json" -Headers @{
            "Authorization" = "Bearer $TOKEN"
        }
        Write-Host "✓ Promo Product updated successfully" -ForegroundColor Green
        Write-Host "New Products: $($updateResponse.data.products.Count) items" -ForegroundColor Gray
        foreach ($p in $updateResponse.data.products) {
            Write-Host "  - $($p.product_name)" -ForegroundColor Gray
        }
    } catch {
        Write-Host "✗ Failed to update promo" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
    }
    Write-Host ""
}

Write-Host "=== TEST COMPLETED ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Summary:" -ForegroundColor Yellow
Write-Host "- Promo Normal: General discount for all products" -ForegroundColor Gray
Write-Host "- Promo Product: Discount for specific products only" -ForegroundColor Gray
Write-Host "- Promo Bundle: Discount when buying product combinations" -ForegroundColor Gray
