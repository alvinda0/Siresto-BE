# Promo Categories Implementation Summary

## Overview
Sistem promo telah diupdate untuk mendukung 3 kategori promo:
1. **Normal** - Promo umum untuk semua produk
2. **Product** - Promo untuk produk tertentu
3. **Bundle** - Promo bundle (beli kombinasi produk)

## Files Modified

### 1. Entity Layer
**File: `internal/entity/promo.go`**
- ✅ Added `PromoCategory` field to `Promo` struct
- ✅ Added `PromoProduct` entity for product-specific promos
- ✅ Added `PromoBundle` entity for bundle promos
- ✅ Added relations: `PromoProducts` and `PromoBundles`

**File: `internal/entity/promo_dto.go`**
- ✅ Added `PromoCategory` to `CreatePromoRequest`
- ✅ Added `ProductIDs` field for product promos
- ✅ Added `BundleItems` field for bundle promos
- ✅ Added `BundleItem` struct
- ✅ Updated `UpdatePromoRequest` with same fields
- ✅ Added `PromoProductResponse` struct
- ✅ Added `PromoBundleResponse` struct
- ✅ Updated `PromoResponse` to include products and bundle_items

### 2. Repository Layer
**File: `internal/repository/promo_repository.go`**
- ✅ Added `CreatePromoProducts()` method
- ✅ Added `CreatePromoBundles()` method
- ✅ Added `DeletePromoProducts()` method
- ✅ Added `DeletePromoBundles()` method
- ✅ Updated `FindByID()` to preload PromoProducts and PromoBundles
- ✅ Updated `FindByCode()` to preload PromoProducts and PromoBundles
- ✅ Updated `FindByCompany()` to preload PromoProducts and PromoBundles
- ✅ Updated `FindByBranch()` to preload PromoProducts and PromoBundles

### 3. Service Layer
**File: `internal/service/promo_service.go`**
- ✅ Updated `CreatePromo()` to handle 3 promo categories
- ✅ Added validation for promo_category
- ✅ Added validation for product_ids (product promo)
- ✅ Added validation for bundle_items (bundle promo)
- ✅ Added logic to create promo_products
- ✅ Added logic to create promo_bundles
- ✅ Updated `UpdatePromo()` to handle category changes
- ✅ Added logic to update promo_products
- ✅ Added logic to update promo_bundles
- ✅ Updated `toResponse()` to include products and bundle_items

### 4. Handler Layer
**File: `internal/handler/promo_handler.go`**
- ℹ️ No changes needed - handler uses service layer

## Database Changes

### New Column
```sql
ALTER TABLE promos 
ADD COLUMN promo_category VARCHAR(20) NOT NULL DEFAULT 'normal';
```

### New Tables
```sql
-- Promo Products
CREATE TABLE promo_products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  promo_id UUID NOT NULL REFERENCES promos(id) ON DELETE CASCADE,
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(promo_id, product_id)
);

-- Promo Bundles
CREATE TABLE promo_bundles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  promo_id UUID NOT NULL REFERENCES promos(id) ON DELETE CASCADE,
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  quantity INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(promo_id, product_id)
);
```

### Indexes
```sql
CREATE INDEX idx_promo_products_promo_id ON promo_products(promo_id);
CREATE INDEX idx_promo_products_product_id ON promo_products(product_id);
CREATE INDEX idx_promo_bundles_promo_id ON promo_bundles(promo_id);
CREATE INDEX idx_promo_bundles_product_id ON promo_bundles(product_id);
CREATE INDEX idx_promos_promo_category ON promos(promo_category);
```

## Migration Files

### Created Files:
1. ✅ `add_promo_category_and_tables.go` - Migration script
2. ✅ `run_promo_category_migration.ps1` - PowerShell runner

## Documentation Files

### Created Files:
1. ✅ `PROMO_CATEGORIES.md` - Complete documentation
2. ✅ `PROMO_CATEGORIES_QUICK_START.md` - Quick start guide
3. ✅ `PROMO_CATEGORIES_IMPLEMENTATION.md` - This file
4. ✅ `test_promo_categories.ps1` - Testing script

## API Changes

### Request Format

#### Promo Normal
```json
{
  "promo_category": "normal",
  "name": "Diskon Umum",
  "code": "DISC20",
  "type": "percentage",
  "value": 20
}
```

#### Promo Product
```json
{
  "promo_category": "product",
  "name": "Diskon Laptop",
  "code": "LAPTOP50",
  "type": "percentage",
  "value": 50,
  "product_ids": ["uuid1", "uuid2"]
}
```

#### Promo Bundle
```json
{
  "promo_category": "bundle",
  "name": "Paket Gaming",
  "code": "GAMING999",
  "type": "fixed",
  "value": 1000000,
  "bundle_items": [
    {
      "product_id": "uuid1",
      "quantity": 1
    },
    {
      "product_id": "uuid2",
      "quantity": 2
    }
  ]
}
```

### Response Format

#### Promo Normal Response
```json
{
  "id": "...",
  "promo_category": "normal",
  "name": "Diskon Umum",
  ...
}
```

#### Promo Product Response
```json
{
  "id": "...",
  "promo_category": "product",
  "name": "Diskon Laptop",
  "products": [
    {
      "product_id": "...",
      "product_name": "Laptop ASUS",
      
    }
  ],
  ...
}
```

#### Promo Bundle Response
```json
{
  "id": "...",
  "promo_category": "bundle",
  "name": "Paket Gaming",
  "bundle_items": [
    {
      "product_id": "...",
      "product_name": "Laptop",
      
      "quantity": 1
    },
    {
      "product_id": "...",
      "product_name": "Mouse",
      
      "quantity": 2
    }
  ],
  ...
}
```

## Validation Rules

### Promo Normal
- ✅ No additional validation
- ✅ Works like before

### Promo Product
- ✅ Must have `product_ids` array
- ✅ Minimum 1 product required
- ✅ Products must exist in database

### Promo Bundle
- ✅ Must have `bundle_items` array
- ✅ Minimum 2 products required
- ✅ Each item must have `product_id` and `quantity`
- ✅ Quantity must be >= 1

## Business Logic

### Create Promo
1. Validate promo_category (normal/product/bundle)
2. Validate required fields based on category
3. Create promo record
4. If product category: create promo_products records
5. If bundle category: create promo_bundles records
6. Return response with relations

### Update Promo
1. Find existing promo
2. Update basic fields
3. If category is product and product_ids provided:
   - Delete old promo_products
   - Create new promo_products
4. If category is bundle and bundle_items provided:
   - Delete old promo_bundles
   - Create new promo_bundles
5. Return updated response with relations

### Delete Promo
- Cascade delete automatically removes promo_products and promo_bundles

## Testing

### Run Migration
```powershell
.\run_promo_category_migration.ps1
```

### Run Tests
```powershell
.\test_promo_categories.ps1
```

### Manual Testing
See `PROMO_CATEGORIES_QUICK_START.md` for cURL examples

## Backward Compatibility

✅ **Fully backward compatible**
- Existing promos will have `promo_category = 'normal'` by default
- Old API requests without `promo_category` will default to 'normal'
- No breaking changes to existing functionality

## Next Steps

### For Order Integration:
1. Update order service to check promo category
2. For product promo: validate if ordered products match promo products
3. For bundle promo: validate if order contains all bundle items with correct quantities
4. Apply discount based on promo rules

### For Frontend:
1. Add promo category selector in create/edit form
2. Show product selector for product promo
3. Show bundle items builder for bundle promo
4. Display products/bundle items in promo list/detail

## Performance Considerations

- ✅ Indexes added for foreign keys
- ✅ Preload used to avoid N+1 queries
- ✅ Unique constraints prevent duplicate entries
- ✅ Cascade delete for cleanup

## Security

- ✅ Multi-tenant isolation maintained
- ✅ RBAC permissions still apply
- ✅ Branch-level access control preserved
- ✅ Validation prevents invalid data

## Summary

✅ Entity layer updated
✅ Repository layer updated
✅ Service layer updated
✅ Migration created
✅ Documentation created
✅ Testing script created
✅ Backward compatible
✅ Multi-tenant safe
✅ Performance optimized

**Status: Ready for testing and deployment**
