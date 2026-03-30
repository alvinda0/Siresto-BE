# Double Promo Feature - Dokumentasi

## 🎉 Fitur Double Promo

Order sekarang bisa menggunakan **MULTIPLE PROMO** sekaligus!

### ✅ Aturan Double Promo

1. **Promo yang Sama Tidak Bisa Dipakai 2x**
   - ❌ WEEKEND10 + WEEKEND10 = DITOLAK
   - Error: "promo code has already been applied to this order"

2. **Promo Berbeda Bisa Dikombinasikan**
   - ✅ WEEKEND10 + LEBARAN2026 = BERHASIL
   - ✅ DISKON20 + CASHBACK50 = BERHASIL
   - Discount akan dijumlahkan

3. **Promo Bisa Ditambahkan Kapan Saja**
   - Promo pertama saat create order
   - Promo kedua saat payment
   - Atau semua promo saat payment

## 📋 Contoh Penggunaan

### Skenario 1: Double Promo (Create Order + Payment)

#### Step 1: Create Order dengan Promo WEEKEND10
```json
POST /api/v1/external/orders
{
  "table_number": "A1",
  "order_method": "DINE_IN",
  "customer_name": "John Doe",
  "promo_code": "WEEKEND10",
  "order_items": [
    {
      "product_id": "uuid-product",
      "quantity": 2
    }
  ]
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "order-uuid",
    "subtotal_amount": 100000,
    "discount_amount": 10000,
    "tax_amount": 9900,
    "total_amount": 99900,
    "promo_code": "WEEKEND10",
    "promo_details": [
      {
        "promo_code": "WEEKEND10",
        "promo_name": "Weekend Discount 10%",
        "promo_type": "percentage",
        "promo_value": 10
      }
    ]
  }
}
```

#### Step 2: Payment dengan Promo Tambahan LEBARAN2026
```json
POST /api/v1/external/orders/order-uuid/payment
{
  "payment_method": "TUNAI",
  "paid_amount": 100000,
  "promo_code": "LEBARAN2026",
  "payment_note": "Pembayaran tunai dengan double promo"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Payment processed successfully",
  "data": {
    "order_id": "order-uuid",
    "payment_method": "TUNAI",
    "payment_status": "PAID",
    "subtotal_amount": 100000,
    "discount_amount": 25000,
    "tax_amount": 8250,
    "total_amount": 83250,
    "paid_amount": 100000,
    "change_amount": 16750,
    "promo_details": [
      {
        "promo_code": "WEEKEND10",
        "promo_name": "Weekend Discount 10%",
        "promo_type": "percentage",
        "promo_value": 10
      },
      {
        "promo_code": "LEBARAN2026",
        "promo_name": "Lebaran Special 15%",
        "promo_type": "percentage",
        "promo_value": 15
      }
    ]
  }
}
```

### Skenario 2: Triple Promo!

#### Step 1: Create Order dengan 2 Promo
```json
POST /api/v1/external/orders
{
  "table_number": "B2",
  "order_method": "DINE_IN",
  "promo_code": "WEEKEND10",
  "order_items": [...]
}
```

#### Step 2: Payment dengan Promo Ketiga
```json
POST /api/v1/external/orders/{id}/payment
{
  "payment_method": "QRIS",
  "paid_amount": 75000,
  "promo_code": "CASHBACK5"
}
```

**Result:** Order akan punya 2 promo (WEEKEND10 + CASHBACK5)

### Skenario 3: Error - Promo Duplikat

#### Step 1: Create Order dengan WEEKEND10
```json
POST /api/v1/external/orders
{
  "promo_code": "WEEKEND10",
  ...
}
```

#### Step 2: Payment dengan WEEKEND10 Lagi (DITOLAK!)
```json
POST /api/v1/external/orders/{id}/payment
{
  "payment_method": "TUNAI",
  "paid_amount": 100000,
  "promo_code": "WEEKEND10"
}
```

**Error Response:**
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "promo error: promo code has already been applied to this order"
}
```

## 📊 Response Format

### PromoCode Field
Promo codes disimpan sebagai comma-separated string:
```json
{
  "promo_code": "WEEKEND10,LEBARAN2026,CASHBACK5"
}
```

### PromoDetails Field
Array of promo details:
```json
{
  "promo_details": [
    {
      "promo_id": "uuid-1",
      "promo_name": "Weekend Discount 10%",
      "promo_code": "WEEKEND10",
      "promo_type": "percentage",
      "promo_value": 10,
      "max_discount": 50000,
      "min_transaction": 50000
    },
    {
      "promo_id": "uuid-2",
      "promo_name": "Lebaran Special 15%",
      "promo_code": "LEBARAN2026",
      "promo_type": "percentage",
      "promo_value": 15,
      "max_discount": 100000,
      "min_transaction": 100000
    }
  ]
}
```

### DiscountAmount Field
Total discount dari semua promo:
```json
{
  "discount_amount": 25000
}
```

## 💡 Perhitungan Discount

### Contoh: 2 Promo Percentage

**Order:**
- Subtotal: 100.000
- Promo 1: WEEKEND10 (10%)
- Promo 2: LEBARAN2026 (15%)

**Perhitungan:**
1. Discount dari WEEKEND10: 100.000 × 10% = 10.000
2. Discount dari LEBARAN2026: 100.000 × 15% = 15.000
3. Total Discount: 10.000 + 15.000 = 25.000
4. Amount after discount: 100.000 - 25.000 = 75.000
5. Tax (11%): 75.000 × 11% = 8.250
6. Total: 75.000 + 8.250 = 83.250

### Contoh: Promo Percentage + Fixed

**Order:**
- Subtotal: 200.000
- Promo 1: DISKON20 (20%)
- Promo 2: CASHBACK10K (fixed 10.000)

**Perhitungan:**
1. Discount dari DISKON20: 200.000 × 20% = 40.000
2. Discount dari CASHBACK10K: 10.000
3. Total Discount: 40.000 + 10.000 = 50.000
4. Amount after discount: 200.000 - 50.000 = 150.000
5. Tax (11%): 150.000 × 11% = 16.500
6. Total: 150.000 + 16.500 = 166.500

## 🔍 Cara Cek Promo yang Dipakai

### Get Order Detail
```bash
GET /api/v1/external/orders/{order_id}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "uuid",
    "promo_code": "WEEKEND10,LEBARAN2026",
    "discount_amount": 25000,
    "promo_details": [
      {
        "promo_code": "WEEKEND10",
        "promo_name": "Weekend Discount 10%"
      },
      {
        "promo_code": "LEBARAN2026",
        "promo_name": "Lebaran Special 15%"
      }
    ]
  }
}
```

## ⚠️ Validasi

1. **Promo yang sama tidak bisa dipakai 2x**
   - System akan cek apakah promo code sudah ada di list
   - Jika sudah ada, return error

2. **Setiap promo tetap harus valid**
   - Active status
   - Dalam periode valid
   - Quota masih ada
   - Minimum transaction terpenuhi

3. **Usage count tetap di-increment**
   - Setiap promo yang dipakai, usage count akan bertambah
   - Tidak ada decrement saat ganti promo (karena tidak ada ganti, hanya tambah)

## 💻 Implementation Details

### Database
- `promo_code` field menyimpan comma-separated promo codes
- `discount_amount` menyimpan total discount dari semua promo
- `promo_id` menyimpan ID promo terakhir (untuk backward compatibility)

### Logic
1. Split promo_code by comma
2. Check if new promo already in list
3. If not, apply new promo
4. Append to promo_code string
5. Add discount to total discount
6. Recalculate tax and total

## 🎯 Best Practice

1. **Validasi promo sebelum apply**
   ```bash
   GET /api/v1/external/promos/validate/{promo_code}
   ```

2. **Tampilkan semua promo yang aktif**
   - Show promo_details array di UI
   - Biarkan user tahu promo apa saja yang sudah dipakai

3. **Konfirmasi sebelum tambah promo**
   - "Tambah promo LEBARAN2026 ke order ini?"
   - "Total discount akan menjadi Rp 25.000"

4. **Handle error dengan baik**
   - Jika promo sama: "Promo WEEKEND10 sudah digunakan"
   - Jika promo tidak valid: "Promo tidak valid atau sudah expired"

## 🎉 Keuntungan Double Promo

1. **Flexibility** - Customer bisa pakai multiple promo
2. **Better UX** - Tidak perlu pilih satu promo saja
3. **Higher Discount** - Kombinasi promo = discount lebih besar
4. **Marketing** - Bisa buat campaign "Stack your promos!"

## 📝 Summary

- ✅ Multiple promo bisa dikombinasikan
- ✅ Promo yang sama tidak bisa dipakai 2x
- ✅ Discount dijumlahkan dari semua promo
- ✅ Promo bisa ditambahkan di create order atau payment
- ✅ Response menampilkan semua promo details
