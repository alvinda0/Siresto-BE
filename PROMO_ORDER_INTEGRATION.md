# Promo Categories - Order Integration Guide

## Overview

Panduan ini menjelaskan bagaimana mengintegrasikan 3 kategori promo (normal, product, bundle) dengan sistem order.

## Current Order Promo Implementation

Saat ini, order sudah support promo dengan field:
- `promo_code` - Kode promo yang digunakan
- `promo_discount` - Nilai diskon dari promo
- `promo_id` - ID promo yang digunakan

## Integration Requirements

### 1. Promo Normal
**Behavior**: Diskon diterapkan ke total order

**Validation**:
- ✅ Promo aktif dan belum expired
- ✅ Quota masih tersedia
- ✅ Minimum transaction terpenuhi
- ✅ Tidak perlu validasi produk

**Calculation**:
```go
if promo.Type == "percentage" {
    discount = subtotal * (promo.Value / 100)
    if promo.MaxDiscount != nil && discount > *promo.MaxDiscount {
        discount = *promo.MaxDiscount
    }
} else {
    discount = promo.Value
}
```

### 2. Promo Product
**Behavior**: Diskon hanya untuk produk yang terdaftar di promo

**Validation**:
- ✅ Promo aktif dan belum expired
- ✅ Quota masih tersedia
- ✅ Order mengandung minimal 1 produk yang terdaftar di promo

**Calculation**:
```go
// Hitung subtotal hanya untuk produk yang eligible
eligibleSubtotal := 0.0
for _, item := range orderItems {
    if isProductInPromo(item.ProductID, promo.PromoProducts) {
        eligibleSubtotal += item.Price * item.Quantity
    }
}

if promo.Type == "percentage" {
    discount = eligibleSubtotal * (promo.Value / 100)
    if promo.MaxDiscount != nil && discount > *promo.MaxDiscount {
        discount = *promo.MaxDiscount
    }
} else {
    discount = promo.Value
}
```

### 3. Promo Bundle
**Behavior**: Diskon diterapkan jika order mengandung semua produk bundle dengan quantity yang sesuai

**Validation**:
- ✅ Promo aktif dan belum expired
- ✅ Quota masih tersedia
- ✅ Order mengandung SEMUA produk bundle
- ✅ Quantity setiap produk >= quantity yang diminta di bundle

**Calculation**:
```go
// Validasi bundle terpenuhi
if !isBundleComplete(orderItems, promo.PromoBundles) {
    return error("Bundle requirements not met")
}

// Hitung berapa kali bundle bisa diterapkan
bundleCount := calculateBundleCount(orderItems, promo.PromoBundles)

if promo.Type == "percentage" {
    // Untuk percentage, hitung dari total bundle items
    bundleSubtotal := calculateBundleSubtotal(orderItems, promo.PromoBundles)
    discount = bundleSubtotal * (promo.Value / 100) * bundleCount
    if promo.MaxDiscount != nil && discount > *promo.MaxDiscount {
        discount = *promo.MaxDiscount
    }
} else {
    // Untuk fixed, kalikan dengan jumlah bundle
    discount = promo.Value * bundleCount
}
```

## Implementation Steps

### Step 1: Update Order Service

**File**: `internal/service/order_service.go`

Add helper functions:

```go
// Check if product is in promo
func isProductInPromo(productID uuid.UUID, promoProducts []entity.PromoProduct) bool {
    for _, pp := range promoProducts {
        if pp.ProductID == productID {
            return true
        }
    }
    return false
}

// Check if bundle is complete
func isBundleComplete(orderItems []entity.OrderItem, promoBundles []entity.PromoBundle) bool {
    for _, bundle := range promoBundles {
        found := false
        for _, item := range orderItems {
            if item.ProductID == bundle.ProductID && item.Quantity >= bundle.Quantity {
                found = true
                break
            }
        }
        if !found {
            return false
        }
    }
    return true
}

// Calculate how many times bundle can be applied
func calculateBundleCount(orderItems []entity.OrderItem, promoBundles []entity.PromoBundle) int {
    minCount := 999999
    for _, bundle := range promoBundles {
        for _, item := range orderItems {
            if item.ProductID == bundle.ProductID {
                count := item.Quantity / bundle.Quantity
                if count < minCount {
                    minCount = count
                }
                break
            }
        }
    }
    return minCount
}
```

### Step 2: Update ValidateAndApplyPromo Function

```go
func (s *orderService) ValidateAndApplyPromo(
    promo *entity.Promo, 
    orderItems []entity.OrderItem, 
    subtotal float64,
) (float64, error) {
    // Basic validation
    now := time.Now()
    if !promo.IsActive {
        return 0, errors.New("promo is not active")
    }
    if now.Before(promo.StartDate) || now.After(promo.EndDate) {
        return 0, errors.New("promo is expired or not yet started")
    }
    if promo.Quota != nil && promo.UsedCount >= *promo.Quota {
        return 0, errors.New("promo quota exceeded")
    }
    if promo.MinTransaction != nil && subtotal < *promo.MinTransaction {
        return 0, errors.New("minimum transaction not met")
    }

    var discount float64

    switch promo.PromoCategory {
    case "normal":
        // Apply to total order
        discount = s.calculateNormalPromoDiscount(promo, subtotal)

    case "product":
        // Apply only to eligible products
        if len(promo.PromoProducts) == 0 {
            return 0, errors.New("promo has no products defined")
        }
        discount = s.calculateProductPromoDiscount(promo, orderItems)

    case "bundle":
        // Apply if bundle requirements met
        if len(promo.PromoBundles) < 2 {
            return 0, errors.New("promo bundle requires at least 2 products")
        }
        if !isBundleComplete(orderItems, promo.PromoBundles) {
            return 0, errors.New("bundle requirements not met")
        }
        discount = s.calculateBundlePromoDiscount(promo, orderItems)

    default:
        return 0, errors.New("invalid promo category")
    }

    return discount, nil
}

func (s *orderService) calculateNormalPromoDiscount(promo *entity.Promo, subtotal float64) float64 {
    var discount float64
    if promo.Type == "percentage" {
        discount = subtotal * (promo.Value / 100)
        if promo.MaxDiscount != nil && discount > *promo.MaxDiscount {
            discount = *promo.MaxDiscount
        }
    } else {
        discount = promo.Value
    }
    return discount
}

func (s *orderService) calculateProductPromoDiscount(promo *entity.Promo, orderItems []entity.OrderItem) float64 {
    eligibleSubtotal := 0.0
    for _, item := range orderItems {
        if isProductInPromo(item.ProductID, promo.PromoProducts) {
            eligibleSubtotal += item.Price * float64(item.Quantity)
        }
    }

    var discount float64
    if promo.Type == "percentage" {
        discount = eligibleSubtotal * (promo.Value / 100)
        if promo.MaxDiscount != nil && discount > *promo.MaxDiscount {
            discount = *promo.MaxDiscount
        }
    } else {
        discount = promo.Value
    }
    return discount
}

func (s *orderService) calculateBundlePromoDiscount(promo *entity.Promo, orderItems []entity.OrderItem) float64 {
    bundleCount := calculateBundleCount(orderItems, promo.PromoBundles)
    
    var discount float64
    if promo.Type == "percentage" {
        bundleSubtotal := 0.0
        for _, bundle := range promo.PromoBundles {
            for _, item := range orderItems {
                if item.ProductID == bundle.ProductID {
                    bundleSubtotal += item.Price * float64(bundle.Quantity)
                    break
                }
            }
        }
        discount = bundleSubtotal * (promo.Value / 100) * float64(bundleCount)
        if promo.MaxDiscount != nil && discount > *promo.MaxDiscount {
            discount = *promo.MaxDiscount
        }
    } else {
        discount = promo.Value * float64(bundleCount)
    }
    return discount
}
```

### Step 3: Update Repository to Preload Relations

Make sure when fetching promo by code, you preload products and bundles:

```go
func (r *promoRepository) FindByCode(code string, companyID uuid.UUID, branchID *uuid.UUID) (*entity.Promo, error) {
    var promo entity.Promo
    query := r.db.Preload("Company").Preload("Branch").
        Preload("PromoProducts.Product").
        Preload("PromoBundles.Product").
        Where("code = ? AND company_id = ?", code, companyID)
    
    // ... rest of the code
}
```

## Testing Scenarios

### Test 1: Promo Normal
```json
{
  "items": [
    {"product_id": "uuid1", "quantity": 2, "price": 100000},
    {"product_id": "uuid2", "quantity": 1, "price": 200000}
  ],
  "promo_code": "NEWYEAR2024"
}
```

Expected:
- Subtotal: 400,000
- Discount: 60,000 (15%)
- Total: 340,000

### Test 2: Promo Product
```json
{
  "items": [
    {"product_id": "uuid1", "quantity": 1, "price": 1000000},  // eligible
    {"product_id": "uuid3", "quantity": 1, "price": 500000}    // not eligible
  ],
  "promo_code": "LAPTOP50"
}
```

Expected:
- Subtotal: 1,500,000
- Eligible Subtotal: 1,000,000
- Discount: 500,000 (50% of eligible)
- Total: 1,000,000

### Test 3: Promo Bundle
```json
{
  "items": [
    {"product_id": "uuid1", "quantity": 1, "price": 5000000},  // bundle item 1 (qty: 1)
    {"product_id": "uuid2", "quantity": 2, "price": 200000},   // bundle item 2 (qty: 2)
    {"product_id": "uuid3", "quantity": 1, "price": 300000}    // bundle item 3 (qty: 1)
  ],
  "promo_code": "GAMING999"
}
```

Expected:
- Subtotal: 5,700,000
- Bundle Complete: Yes
- Discount: 1,000,000 (fixed)
- Total: 4,700,000

## Error Messages

### Promo Product
```json
{
  "error": "No eligible products in order for this promo"
}
```

### Promo Bundle
```json
{
  "error": "Bundle requirements not met. Required: Laptop x1, Mouse x2, Keyboard x1"
}
```

## Response Format

Add promo details to order response:

```json
{
  "id": "...",
  "promo_code": "GAMING999",
  "promo_discount": 1000000,
  "promo_details": {
    "id": "...",
    "name": "Paket Gaming",
    "category": "bundle",
    "type": "fixed",
    "value": 1000000,
    "bundle_items": [
      {
        "product_name": "Laptop",
        "quantity": 1
      },
      {
        "product_name": "Mouse",
        "quantity": 2
      }
    ]
  }
}
```

## Best Practices

1. **Always preload relations** when fetching promo
2. **Validate bundle completely** before applying discount
3. **Calculate eligible subtotal** for product promos
4. **Handle multiple bundle applications** correctly
5. **Show clear error messages** when promo can't be applied
6. **Log promo usage** for analytics
7. **Update used_count** after successful order

## Next Steps

1. Implement helper functions in order service
2. Update ValidateAndApplyPromo function
3. Add comprehensive error messages
4. Test all 3 promo categories
5. Update order response to include promo details
6. Add promo usage analytics

## Notes

- Promo normal paling mudah diimplementasikan
- Promo product perlu filter produk yang eligible
- Promo bundle paling kompleks, perlu validasi lengkap
- Pertimbangkan edge case: multiple bundles, partial bundles, dll
- Pastikan discount tidak melebihi subtotal
