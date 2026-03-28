# PowerShell script untuk testing Order API
# Usage: .\test_orders.ps1

$BASE_URL = "http://localhost:8080/api/v1"
$TOKEN = ""
$COMPANY_ID = ""
$BRANCH_ID = ""
$PRODUCT_ID = ""
$ORDER_ID = ""

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "   ORDER API TESTING SCRIPT" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

# Function to make HTTP requests
function Invoke-ApiRequest {
    param(
        [string]$Method,
        [string]$Url,
        [string]$Token = "",
        [object]$Body = $null
    )
    
    $headers = @{
        "Content-Type" = "application/json"
    }
    
    if ($Token) {
        $headers["Authorization"] = "Bearer $Token"
    }
    
    try {
        if ($Body) {
            $jsonBody = $Body | ConvertTo-Json -Depth 10
            $response = Invoke-RestMethod -Uri $Url -Method $Method -Headers $headers -Body $jsonBody
        } else {
            $response = Invoke-RestMethod -Uri $Url -Method $Method -Headers $headers
        }
        return $response
    } catch {
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
        if ($_.Exception.Response) {
            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $reader.BaseStream.Position = 0
            $responseBody = $reader.ReadToEnd()
            Write-Host "Response: $responseBody" -ForegroundColor Red
        }
        return $null
    }
}

# 1. Login untuk mendapatkan token
Write-Host "1. Login..." -ForegroundColor Yellow
$loginData = @{
    email = "owner@example.com"
    password = "password123"
}

$loginResponse = Invoke-ApiRequest -Method "POST" -Url "$BASE_URL/login" -Body $loginData

if ($loginResponse -and $loginResponse.data.token) {
    $TOKEN = $loginResponse.data.token
    $COMPANY_ID = $loginResponse.data.user.company_id
    $BRANCH_ID = $loginResponse.data.user.branch_id
    Write-Host "✓ Login successful" -ForegroundColor Green
    Write-Host "  Token: $($TOKEN.Substring(0, 20))..." -ForegroundColor Gray
    Write-Host "  Company ID: $COMPANY_ID" -ForegroundColor Gray
    Write-Host "  Branch ID: $BRANCH_ID" -ForegroundColor Gray
} else {
    Write-Host "✗ Login failed" -ForegroundColor Red
    exit
}

Write-Host ""

# 2. Get Products untuk mendapatkan product_id
Write-Host "2. Get Products..." -ForegroundColor Yellow
$productsResponse = Invoke-ApiRequest -Method "GET" -Url "$BASE_URL/external/products?limit=1" -Token $TOKEN

if ($productsResponse -and $productsResponse.data.Count -gt 0) {
    $PRODUCT_ID = $productsResponse.data[0].id
    Write-Host "✓ Product found" -ForegroundColor Green
    Write-Host "  Product ID: $PRODUCT_ID" -ForegroundColor Gray
    Write-Host "  Product Name: $($productsResponse.data[0].name)" -ForegroundColor Gray
    Write-Host "  Price: $($productsResponse.data[0].price)" -ForegroundColor Gray
} else {
    Write-Host "✗ No products found. Please create a product first." -ForegroundColor Red
    exit
}

Write-Host ""

# 3. Create Order (Authenticated)
Write-Host "3. Create Order (Authenticated)..." -ForegroundColor Yellow
$createOrderData = @{
    customer_name = "Sasa"
    customer_phone = "08123123123"
    table_number = "A1"
    notes = "Test Order from PowerShell"
    referral_code = ""
    order_method = "DINE_IN"
    promo_code = ""
    order_items = @(
        @{
            product_id = $PRODUCT_ID
            quantity = 3
            note = "Extra pedas"
        }
    )
}

$createResponse = Invoke-ApiRequest -Method "POST" -Url "$BASE_URL/orders" -Token $TOKEN -Body $createOrderData

if ($createResponse -and $createResponse.success) {
    $ORDER_ID = $createResponse.data.id
    Write-Host "✓ Order created successfully" -ForegroundColor Green
    Write-Host "  Order ID: $ORDER_ID" -ForegroundColor Gray
    Write-Host "  Status: $($createResponse.data.status)" -ForegroundColor Gray
    Write-Host "  Total Amount: $($createResponse.data.total_amount)" -ForegroundColor Gray
    Write-Host "  Items: $($createResponse.data.order_items.Count)" -ForegroundColor Gray
} else {
    Write-Host "✗ Failed to create order" -ForegroundColor Red
}

Write-Host ""

# 4. Get Order by ID
Write-Host "4. Get Order by ID..." -ForegroundColor Yellow
$getOrderResponse = Invoke-ApiRequest -Method "GET" -Url "$BASE_URL/orders/$ORDER_ID" -Token $TOKEN

if ($getOrderResponse -and $getOrderResponse.success) {
    Write-Host "✓ Order retrieved successfully" -ForegroundColor Green
    Write-Host "  Customer: $($getOrderResponse.data.customer_name)" -ForegroundColor Gray
    Write-Host "  Table: $($getOrderResponse.data.table_number)" -ForegroundColor Gray
    Write-Host "  Method: $($getOrderResponse.data.order_method)" -ForegroundColor Gray
    Write-Host "  Status: $($getOrderResponse.data.status)" -ForegroundColor Gray
} else {
    Write-Host "✗ Failed to get order" -ForegroundColor Red
}

Write-Host ""

# 5. Get All Orders
Write-Host "5. Get All Orders..." -ForegroundColor Yellow
$getAllResponse = Invoke-ApiRequest -Method "GET" -Url "$BASE_URL/orders?page=1&limit=10" -Token $TOKEN

if ($getAllResponse -and $getAllResponse.success) {
    Write-Host "✓ Orders retrieved successfully" -ForegroundColor Green
    Write-Host "  Total Orders: $($getAllResponse.meta.total)" -ForegroundColor Gray
    Write-Host "  Current Page: $($getAllResponse.meta.page)" -ForegroundColor Gray
    Write-Host "  Total Pages: $($getAllResponse.meta.total_pages)" -ForegroundColor Gray
} else {
    Write-Host "✗ Failed to get orders" -ForegroundColor Red
}

Write-Host ""

# 6. Update Order
Write-Host "6. Update Order..." -ForegroundColor Yellow
$updateOrderData = @{
    customer_name = "Sasa Updated"
    table_number = "A2"
    status = "CONFIRMED"
    notes = "Updated from PowerShell"
}

$updateResponse = Invoke-ApiRequest -Method "PUT" -Url "$BASE_URL/orders/$ORDER_ID" -Token $TOKEN -Body $updateOrderData

if ($updateResponse -and $updateResponse.success) {
    Write-Host "✓ Order updated successfully" -ForegroundColor Green
    Write-Host "  New Status: $($updateResponse.data.status)" -ForegroundColor Gray
    Write-Host "  New Table: $($updateResponse.data.table_number)" -ForegroundColor Gray
} else {
    Write-Host "✗ Failed to update order" -ForegroundColor Red
}

Write-Host ""

# 7. Create Public Order (No Auth)
Write-Host "7. Create Public Order (No Authentication)..." -ForegroundColor Yellow
$publicOrderData = @{
    company_id = $COMPANY_ID
    branch_id = $BRANCH_ID
    customer_name = "Public Customer"
    customer_phone = "08199999999"
    table_number = "B5"
    notes = "Public order test"
    order_method = "TAKE_AWAY"
    order_items = @(
        @{
            product_id = $PRODUCT_ID
            quantity = 1
            note = "No note"
        }
    )
}

$publicResponse = Invoke-ApiRequest -Method "POST" -Url "$BASE_URL/public/orders" -Body $publicOrderData

if ($publicResponse -and $publicResponse.success) {
    Write-Host "✓ Public order created successfully" -ForegroundColor Green
    Write-Host "  Order ID: $($publicResponse.data.id)" -ForegroundColor Gray
    Write-Host "  Method: $($publicResponse.data.order_method)" -ForegroundColor Gray
    Write-Host "  Total: $($publicResponse.data.total_amount)" -ForegroundColor Gray
} else {
    Write-Host "✗ Failed to create public order" -ForegroundColor Red
}

Write-Host ""
Write-Host "==================================" -ForegroundColor Cyan
Write-Host "   TESTING COMPLETED" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
