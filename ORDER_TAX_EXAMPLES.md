# Contoh Perhitungan Pajak Order

## Contoh 1: Pajak Tunggal

### Setup
- 1 Pajak: PB1 10%

### Order
- 1 × Nasi Goreng @ Rp 25.000 = Rp 25.000

### Perhitungan
```
Subtotal: Rp 25.000
PB1 (10%): 25.000 × 10% = Rp 2.500
Total: Rp 27.500
```

### Response
```json
{
  "subtotal_amount": 25000,
  "tax_amount": 2500,
  "total_amount": 27500,
  "tax_details": [
    {
      "tax_name": "PB1",
      "percentage": 10,
      "priority": 1,
      "base_amount": 25000,
      "tax_amount": 2500
    }
  ]
}
```

---

## Contoh 2: Dua Pajak Bertingkat

### Setup
- Pajak 1: Service Charge 5% (prioritas 1 - dihitung pertama)
- Pajak 2: PB1 10% (prioritas 2 - dihitung kedua)

### Order
- 2 × Nasi Goreng @ Rp 50.000 = Rp 100.000

### Perhitungan
```
Subtotal: Rp 100.000

Step 1 - Service Charge (prioritas 1):
  Base: Rp 100.000
  Tax: 100.000 × 5% = Rp 5.000
  Running Total: Rp 105.000

Step 2 - PB1 (prioritas 2):
  Base: Rp 105.000 (sudah termasuk Service Charge)
  Tax: 105.000 × 10% = Rp 10.500
  Running Total: Rp 115.500

Total Tax: Rp 5.000 + Rp 10.500 = Rp 15.500
Total Amount: Rp 115.500
```

### Response
```json
{
  "subtotal_amount": 100000,
  "tax_amount": 15500,
  "total_amount": 115500,
  "tax_details": [
    {
      "tax_name": "Service Charge",
      "percentage": 5,
      "priority": 1,
      "base_amount": 100000,
      "tax_amount": 5000
    },
    {
      "tax_name": "PB1",
      "percentage": 10,
      "priority": 2,
      "base_amount": 105000,
      "tax_amount": 10500
    }
  ]
}
```

---

## Contoh 3: Tiga Pajak Bertingkat

### Setup
- Pajak 1: PB1 10% (prioritas 1)
- Pajak 2: Service Charge 5% (prioritas 2)
- Pajak 3: Government Tax 2% (prioritas 3)

### Order
- 1 × Steak @ Rp 200.000 = Rp 200.000

### Perhitungan
```
Subtotal: Rp 200.000

Step 1 - PB1 (prioritas 1):
  Base: Rp 200.000
  Tax: 200.000 × 10% = Rp 20.000
  Running Total: Rp 220.000

Step 2 - Service Charge (prioritas 2):
  Base: Rp 220.000
  Tax: 220.000 × 5% = Rp 11.000
  Running Total: Rp 231.000

Step 3 - Government Tax (prioritas 3):
  Base: Rp 231.000
  Tax: 231.000 × 2% = Rp 4.620
  Running Total: Rp 235.620

Total Tax: Rp 20.000 + Rp 11.000 + Rp 4.620 = Rp 35.620
Total Amount: Rp 235.620
```

### Response
```json
{
  "subtotal_amount": 200000,
  "tax_amount": 35620,
  "total_amount": 235620,
  "tax_details": [
    {
      "tax_name": "PB1",
      "percentage": 10,
      "priority": 1,
      "base_amount": 200000,
      "tax_amount": 20000
    },
    {
      "tax_name": "Service Charge",
      "percentage": 5,
      "priority": 2,
      "base_amount": 220000,
      "tax_amount": 11000
    },
    {
      "tax_name": "Government Tax",
      "percentage": 2,
      "priority": 3,
      "base_amount": 231000,
      "tax_amount": 4620
    }
  ]
}
```

---

## Contoh 4: Multiple Items dengan Pajak

### Setup
- Pajak 1: PB1 10% (prioritas 1)
- Pajak 2: Service Charge 5% (prioritas 2)

### Order
- 2 × Nasi Goreng @ Rp 25.000 = Rp 50.000
- 1 × Es Teh @ Rp 5.000 = Rp 5.000
- 3 × Ayam Goreng @ Rp 30.000 = Rp 90.000

### Perhitungan
```
Subtotal: Rp 50.000 + Rp 5.000 + Rp 90.000 = Rp 145.000

Step 1 - PB1 (prioritas 1):
  Base: Rp 145.000
  Tax: 145.000 × 10% = Rp 14.500
  Running Total: Rp 159.500

Step 2 - Service Charge (prioritas 2):
  Base: Rp 159.500
  Tax: 159.500 × 5% = Rp 7.975
  Running Total: Rp 167.475

Total Tax: Rp 14.500 + Rp 7.975 = Rp 22.475
Total Amount: Rp 167.475
```

### Response
```json
{
  "subtotal_amount": 145000,
  "tax_amount": 22475,
  "total_amount": 167475,
  "tax_details": [
    {
      "tax_name": "PB1",
      "percentage": 10,
      "priority": 1,
      "base_amount": 145000,
      "tax_amount": 14500
    },
    {
      "tax_name": "Service Charge",
      "percentage": 5,
      "priority": 2,
      "base_amount": 159500,
      "tax_amount": 7975
    }
  ],
  "order_items": [
    {
      "product_name": "Nasi Goreng",
      "quantity": 2,
      "price": 25000,
      "subtotal": 50000
    },
    {
      "product_name": "Es Teh",
      "quantity": 1,
      "price": 5000,
      "subtotal": 5000
    },
    {
      "product_name": "Ayam Goreng",
      "quantity": 3,
      "price": 30000,
      "subtotal": 90000
    }
  ]
}
```

---

## Contoh 5: Tanpa Pajak

### Setup
- Tidak ada pajak aktif atau semua pajak inactive

### Order
- 1 × Nasi Goreng @ Rp 25.000 = Rp 25.000

### Perhitungan
```
Subtotal: Rp 25.000
Tax: Rp 0 (tidak ada pajak aktif)
Total: Rp 25.000
```

### Response
```json
{
  "subtotal_amount": 25000,
  "tax_amount": 0,
  "total_amount": 25000,
  "tax_details": []
}
```

---

## Perbandingan: Prioritas Berbeda

### Skenario A: PB1 dihitung dulu (prioritas 1), lalu SC (prioritas 2)

```
Subtotal: Rp 100.000
PB1 (10%): 100.000 × 10% = Rp 10.000 → Rp 110.000
SC (5%): 110.000 × 5% = Rp 5.500 → Rp 115.500
Total: Rp 115.500
```

### Skenario B: SC dihitung dulu (prioritas 1), lalu PB1 (prioritas 2)

```
Subtotal: Rp 100.000
SC (5%): 100.000 × 5% = Rp 5.000 → Rp 105.000
PB1 (10%): 105.000 × 10% = Rp 10.500 → Rp 115.500
Total: Rp 115.500
```

**Note**: Dalam contoh ini hasilnya sama, tapi dengan persentase yang berbeda, urutan prioritas akan mempengaruhi hasil akhir!

---

## Tips Konfigurasi Prioritas

1. **Service Charge**: Prioritas 1 (dihitung pertama dari subtotal)
2. **Pajak Pemerintah** (PB1, VAT): Prioritas 2 (dihitung dari subtotal + service charge)
3. **Pajak Tambahan**: Prioritas 3+ (dihitung dari total sebelumnya)

Ini memastikan service charge dihitung dari subtotal murni, sedangkan pajak pemerintah dihitung dari amount yang sudah termasuk service charge.

**Catatan**: Urutan ini bisa disesuaikan dengan kebutuhan bisnis Anda. Yang penting adalah memahami bahwa Priority 1 = dihitung pertama, Priority 2 = dihitung kedua, dst.
