# Payment Promo Rules - Aturan Penggunaan Promo

## 📋 Aturan Utama

### ⚠️ SATU ORDER = SATU PROMO

Setiap order hanya bisa menggunakan **SATU promo** saja. Promo bisa diapply di:
1. Saat create order, ATAU
2. Saat payment

Tidak bisa keduanya!

## ✅ Skenario VALID

### 1. Promo di Create Order
```json
// Step 1: Create order dengan promo
POST /api/v1/external/orders
{
  "table_number": "A1",
  "promo_code": "WEEKEND10",
  "order_items": [...]
}

// Step 2: Payment tanpa promo
POST /api/v1/external/orders/{id}/payment
{
  "payment_method": "TUNAI",
  "paid_amount": 100000
  // Tidak ada promo_code
}
```
**Status:** ✅ BERHASIL

### 2. Promo di Payment
```json
// Step 1: Create order tanpa promo
POST /api/v1/external/orders
{
  "table_number": "A1",
  "order_items": [...]
  // Tidak ada promo_code
}

// Step 2: Payment dengan promo
POST /api/v1/external/orders/{id}/payment
{
  "payment_method": "QRIS",
  "paid_amount": 90000,
  "promo_code": "WEEKEND10"
}
```
**Status:** ✅ BERHASIL

### 3. Tanpa Promo Sama Sekali
```json
// Step 1: Create order tanpa promo
POST /api/v1/external/orders
{
  "table_number": "A1",
  "order_items": [...]
}

// Step 2: Payment tanpa promo
POST /api/v1/external/orders/{id}/payment
{
  "payment_method": "TUNAI",
  "paid_amount": 100000
}
```
**Status:** ✅ BERHASIL

## ❌ Skenario INVALID

### 1. Promo Duplikat (Promo yang Sama)
```json
// Step 1: Create order dengan WEEKEND10
POST /api/v1/external/orders
{
  "promo_code": "WEEKEND10",
  ...
}

// Step 2: Payment dengan WEEKEND10 lagi
POST /api/v1/external/orders/{id}/payment
{
  "payment_method": "TUNAI",
  "paid_amount": 100000,
  "promo_code": "WEEKEND10"  // ❌ DITOLAK!
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

### 2. Promo Berbeda
```json
// Step 1: Create order dengan WEEKEND10
POST /api/v1/external/orders
{
  "promo_code": "WEEKEND10",
  ...
}

// Step 2: Payment dengan DISKON20
POST /api/v1/external/orders/{id}/payment
{
  "payment_method": "QRIS",
  "paid_amount": 80000,
  "promo_code": "DISKON20"  // ❌ DITOLAK!
}
```

**Error Response:**
```json
{
  "status": "error",
  "message": "Failed to process payment",
  "error": "promo error: order already has a promo applied, cannot apply another promo"
}
```

## 📊 Tabel Ringkasan

| Create Order | Payment | Hasil |
|--------------|---------|-------|
| WEEKEND10 | - | ✅ Valid |
| - | WEEKEND10 | ✅ Valid |
| WEEKEND10 | WEEKEND10 | ❌ Error: promo sudah dipakai |
| WEEKEND10 | DISKON20 | ❌ Error: order sudah punya promo |
| - | - | ✅ Valid (tanpa promo) |

## 🎯 Alasan Aturan Ini

1. **Mencegah Double Discount**
   - Hindari customer mendapat diskon ganda
   - Satu promo per transaksi sudah cukup

2. **Konsistensi Perhitungan**
   - Total amount sudah dihitung dengan promo pertama
   - Menghindari recalculation yang kompleks

3. **Business Logic**
   - Promo biasanya tidak bisa dikombinasikan
   - Sesuai dengan praktik bisnis umum

## 💡 Best Practice untuk Frontend

### 1. Cek Promo Sebelum Payment
```javascript
// Sebelum tampilkan form payment, cek apakah order sudah punya promo
const order = await getOrderById(orderId);

if (order.promo_code) {
  // Disable input promo code
  // Tampilkan info: "Order ini sudah menggunakan promo {promo_code}"
  disablePromoInput();
} else {
  // Enable input promo code
  enablePromoInput();
}
```

### 2. Validasi Promo Sebelum Submit
```javascript
// Validasi promo sebelum create order
const validatePromo = async (promoCode) => {
  try {
    const response = await fetch(`/api/v1/external/promos/validate/${promoCode}`);
    if (response.ok) {
      return true;
    }
  } catch (error) {
    showError("Promo tidak valid");
    return false;
  }
};
```

### 3. Handle Error dengan Baik
```javascript
try {
  await processPayment(orderId, paymentData);
} catch (error) {
  if (error.message.includes("promo code has already been applied")) {
    showError("Promo sudah digunakan di order ini");
  } else if (error.message.includes("order already has a promo applied")) {
    showError("Order ini sudah menggunakan promo lain");
  }
}
```

## 🔍 Cara Cek Promo di Order

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
    "promo_code": "WEEKEND10",  // Jika ada promo
    "promo_id": "uuid",
    "discount_amount": 20000,
    "promo_details": {
      "promo_name": "Weekend Discount 10%",
      "promo_code": "WEEKEND10"
    }
  }
}
```

Jika `promo_code` tidak kosong, berarti order sudah menggunakan promo.

## 📝 Checklist untuk Developer

- [ ] Cek `promo_code` di order sebelum tampilkan form payment
- [ ] Disable input promo jika order sudah punya promo
- [ ] Tampilkan info promo yang sudah dipakai
- [ ] Handle error promo duplikat dengan baik
- [ ] Validasi promo sebelum submit
- [ ] Test skenario promo di create order
- [ ] Test skenario promo di payment
- [ ] Test error promo duplikat
- [ ] Test error promo berbeda

## 🎓 FAQ

**Q: Kenapa tidak bisa pakai 2 promo?**
A: Untuk mencegah double discount dan menjaga konsistensi perhitungan.

**Q: Bagaimana jika customer salah pilih promo?**
A: Customer harus cancel order dan buat order baru dengan promo yang benar.

**Q: Bisa ganti promo setelah create order?**
A: Tidak bisa. Promo hanya bisa diapply sekali dan tidak bisa diganti.

**Q: Bagaimana cara tahu order sudah pakai promo?**
A: Cek field `promo_code` di response GET order. Jika tidak kosong, berarti sudah pakai promo.

**Q: Bisa apply promo setelah payment?**
A: Tidak bisa. Promo hanya bisa diapply sebelum payment (saat create order atau saat payment).
