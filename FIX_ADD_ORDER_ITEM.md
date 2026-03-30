# Fix: Add Order Item Not Updating Totals

## Problem

Ketika menambahkan item ke order existing menggunakan endpoint `POST /api/v1/external/orders/quick/:id`, item dan quantity tidak bertambah di response `GET /orders/:id`, dan total amount juga tidak berubah.

## Root Cause

Ada 3 bug di fungsi `AddOrderItem`:

1. **Promo Usage Count Bug**: Ketika recalculate discount, fungsi memanggil `applyPromo()` yang akan increment `UsedCount` promo lagi, padahal promo sudah dipakai sebelumnya saat create order.

2. **GORM Save Not Updating**: Method `Update` menggunakan `db.Save()` yang tidak selalu update semua field, terutama jika nilai field tidak berubah atau bernilai 0. Ini menyebabkan `SubtotalAmount`, `TaxAmount`, dan `TotalAmount` tidak ter-update ke database.

3. **Subtotal Increment Bug**: Subtotal dihitung dengan increment (`order.SubtotalAmount + itemSubtotal`) bukan recalculate dari semua items. Ini menyebabkan jika ada item yang dihapus atau diubah sebelumnya, subtotal jadi tidak akurat.

4. **Duplicate Items**: Setiap kali add item dengan product_id yang sama, sistem membuat item baru di database. Seharusnya jika product_id sudah ada, quantity dijumlahkan saja.

## Solution

### 1. Tambah Method `FindByIDSimple` di Promo Repository

File: `internal/repository/promo_repository.go`

```go
func (r *promoRepository) FindByIDSimple(id uuid.UUID) (*entity.Promo, error) {
	var promo entity.Promo
	err := r.db.First(&promo, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &promo, nil
}
```

Method ini untuk fetch promo tanpa access control check, hanya untuk internal recalculation.

### 2. Fix Logic di `AddOrderItem`

File: `internal/service/order_service.go`

**Check for Existing Item:**

```go
// Check if item with same product_id already exists
var existingItem *entity.OrderItem
for i := range order.OrderItems {
    if order.OrderItems[i].ProductID == req.ProductID {
        existingItem = &order.OrderItems[i]
        break
    }
}

if existingItem != nil {
    // Update existing item quantity
    existingItem.Quantity += req.Quantity
    if req.Note != "" {
        existingItem.Note = req.Note
    }
    if err := s.orderRepo.UpdateOrderItem(existingItem); err != nil {
        return nil, err
    }
} else {
    // Create new order item
    newItem := entity.OrderItem{
        OrderID:   orderID,
        ProductID: req.ProductID,
        Quantity:  req.Quantity,
        Price:     product.Price,
        Note:      req.Note,
    }
    if err := s.orderRepo.CreateOrderItems([]entity.OrderItem{newItem}); err != nil {
        return nil, err
    }
}
```

**Recalculate Subtotal from All Items:**

**Before:**
```go
// Recalculate order totals
itemSubtotal := product.Price * float64(req.Quantity)
newSubtotal := order.SubtotalAmount + itemSubtotal
```

**After:**
```go
// Recalculate subtotal from ALL items (not just increment)
// Fetch fresh order with all items
order, err = s.orderRepo.FindByID(orderID)
if err != nil {
    return nil, errors.New("failed to fetch updated order")
}

var newSubtotal float64
for _, item := range order.OrderItems {
    newSubtotal += item.Price * float64(item.Quantity)
}
```

Sekarang subtotal dihitung ulang dari semua items yang ada, bukan hanya increment. Ini memastikan subtotal selalu akurat.

**Promo Recalculation Fix:**

**Before:**
```go
// Recalculate discount if promo exists
var discountAmount float64
if order.PromoCode != "" {
    discount, _, err := s.applyPromo(order.PromoCode, order.SubtotalAmount, companyID, branchID)
    if err != nil {
        order.PromoCode = ""
        order.PromoID = nil
        discountAmount = 0
    } else {
        discountAmount = discount
    }
}
```

**After:**
```go
// Recalculate discount if promo exists (without incrementing usage count)
var discountAmount float64
if order.PromoCode != "" && order.PromoID != nil {
    promo, err := s.promoRepo.FindByIDSimple(*order.PromoID)
    if err == nil && promo.IsActive {
        // Recalculate discount based on new subtotal
        if promo.Type == "percentage" {
            discountAmount = newSubtotal * (promo.Value / 100)
            if promo.MaxDiscount != nil && discountAmount > *promo.MaxDiscount {
                discountAmount = *promo.MaxDiscount
            }
        } else if promo.Type == "fixed" {
            discountAmount = promo.Value
            if discountAmount > newSubtotal {
                discountAmount = newSubtotal
            }
        }
    } else {
        // Promo no longer valid, remove it
        order.PromoCode = ""
        order.PromoID = nil
    }
}
```

### 3. Fix Repository Update Method

File: `internal/repository/order_repository.go`

**Before:**
```go
func (r *orderRepository) Update(order *entity.Order) error {
	return r.db.Save(order).Error
}
```

**After:**
```go
func (r *orderRepository) Update(order *entity.Order) error {
	return r.db.Model(order).Updates(map[string]interface{}{
		"customer_name":   order.CustomerName,
		"customer_phone":  order.CustomerPhone,
		"table_number":    order.TableNumber,
		"notes":           order.Notes,
		"order_method":    order.OrderMethod,
		"promo_id":        order.PromoID,
		"promo_code":      order.PromoCode,
		"discount_amount": order.DiscountAmount,
		"status":          order.Status,
		"subtotal_amount": order.SubtotalAmount,
		"tax_amount":      order.TaxAmount,
		"total_amount":    order.TotalAmount,
	}).Error
}
```

Menggunakan `Updates` dengan map untuk memastikan semua field ter-update, termasuk yang bernilai 0.

## Changes Made

### Files Modified

1. **internal/repository/promo_repository.go**
   - Added `FindByIDSimple()` method to interface
   - Implemented `FindByIDSimple()` method

2. **internal/repository/order_repository.go**
   - Added `UpdateOrderItem()` method to interface
   - Implemented `UpdateOrderItem()` method
   - Changed `Update()` method from `Save()` to `Updates()` with explicit field map
   - Ensures all fields including zero values are updated

3. **internal/service/order_service.go**
   - Fixed `AddOrderItem()` logic
   - Added check for existing item with same product_id
   - If exists: update quantity, if not: create new item
   - Changed promo recalculation to not increment usage count
   - Fixed subtotal calculation to recalculate from all items

## Testing

### Manual Test

1. Create order dengan 1 item
2. Note subtotal, tax, total, dan items count
3. Add item baru ke order
4. Get order by ID
5. Verify:
   - Items count bertambah
   - Subtotal bertambah sesuai harga item baru
   - Tax recalculated
   - Total recalculated

### Automated Test

```powershell
.\test_add_item_debug.ps1
```

Script ini akan:
1. Create quick order
2. Get order (before)
3. Add item
4. Get order (after)
5. Compare dan show difference

## Expected Behavior

### Before Fix
- Item ditambahkan ke database
- Tapi response GET order tidak menunjukkan item baru
- Subtotal, tax, dan total tidak berubah (tetap nilai lama)
- GORM `Save()` tidak update field yang tidak berubah

### After Fix
- Item ditambahkan ke database
- Response GET order menunjukkan semua items termasuk yang baru
- Subtotal = sum of all items (dihitung ulang dan disimpan)
- Tax = recalculated based on new subtotal
- Total = (subtotal - discount) + tax
- Jika ada promo, discount recalculated tanpa increment usage count
- GORM `Updates()` dengan map memastikan semua field ter-update

## Notes

- Fix ini juga berlaku untuk order yang punya promo code
- Promo usage count tidak akan bertambah saat add item
- Discount amount akan recalculated berdasarkan subtotal baru
- Jika promo sudah tidak valid (expired/inactive), akan dihapus dari order
