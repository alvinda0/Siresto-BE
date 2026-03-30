# Fix COMPLIMENTARY Payment - Total Jadi 0

## Masalah
Payment method COMPLIMENTARY seharusnya membuat order gratis (total = 0), tapi sebelumnya masih menghitung total normal.

## Solusi
Update logika di `ProcessPayment` untuk COMPLIMENTARY:
1. Set `discount_amount` = `subtotal_amount` (full discount)
2. Set `tax_amount` = 0
3. Set `total_amount` = 0
4. Set `paid_amount` = 0
5. Set `change_amount` = 0

## File yang Diubah
- `internal/service/order_service.go` - Fungsi `ProcessPayment`

## Perubahan Kode

```go
// For COMPLIMENTARY, make everything free
if req.PaymentMethod == entity.PaymentMethodComplimentary {
    // Set discount to cover full subtotal (making it free)
    order.DiscountAmount = order.SubtotalAmount
    order.TaxAmount = 0
    order.TotalAmount = 0
    req.PaidAmount = 0
    changeAmount = 0
} else if req.PaymentMethod == entity.PaymentMethodCash {
    // ... existing cash logic
}
```

## Testing

```powershell
# Test COMPLIMENTARY payment
./test_complimentary_payment.ps1
```

## Expected Result

Request:
```json
{
  "payment_method": "COMPLIMENTARY",
  "paid_amount": 0,
  "payment_note": "Complimentary for VIP"
}
```

Response:
```json
{
  "payment_method": "COMPLIMENTARY",
  "payment_status": "PAID",
  "subtotal_amount": 50000,
  "discount_amount": 50000,  // = subtotal
  "tax_amount": 0,            // no tax
  "total_amount": 0,          // FREE!
  "paid_amount": 0,
  "change_amount": 0
}
```

## Behavior Summary

| Payment Method | Total Calculation | Paid Amount |
|----------------|-------------------|-------------|
| TUNAI | (Subtotal - Discount) + Tax | Can be > total (ada kembalian) |
| QRIS/DEBIT/CREDIT/GOPAY/OVO | (Subtotal - Discount) + Tax | Must = total |
| COMPLIMENTARY | **0** (full discount) | **0** (gratis) |
