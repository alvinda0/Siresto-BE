# Test Accounts - SIRESTO POS

Akun-akun ini otomatis dibuat saat aplikasi pertama kali dijalankan.

## 🔐 INTERNAL USERS (Platform SIRESTO)

### 1. SUPER_ADMIN
```
Email: superadmin@siresto.com
Password: admin123
Role: SUPER_ADMIN
```
**Akses:**
- Semua endpoint internal
- Monitoring semua companies dan users
- Buat user SUPPORT dan FINANCE

### 2. SUPPORT
```
Email: support@siresto.com
Password: support123
Role: SUPPORT
```
**Akses:**
- Customer service
- Support client restoran
- Lihat data companies dan users

### 3. FINANCE
```
Email: finance@siresto.com
Password: finance123
Role: FINANCE
```
**Akses:**
- Monitoring pembayaran subscription
- Lihat data financial
- Report keuangan

---

## 🍔 EXTERNAL USERS (Client Restoran)

### Test Company
```
Name: PT Restoran Sejahtera
Type: PT
Branch: Cabang Jakarta Pusat
Address: Jl. Sudirman No. 123, Jakarta
```

### 4. OWNER
```
Email: owner@restaurant.com
Password: owner123
Role: OWNER
```
**Akses:**
- Buat dan kelola company
- Buat dan kelola branch
- Buat user ADMIN, CASHIER, KITCHEN, WAITER
- Full access ke semua data company

### 5. ADMIN
```
Email: admin@restaurant.com
Password: admin123
Role: ADMIN
Company: PT Restoran Sejahtera
```
**Akses:**
- Kelola operasional company
- Kelola semua branch
- Monitoring staff
- Report dan analytics

### 6. CASHIER
```
Email: cashier@restaurant.com
Password: cashier123
Role: CASHIER
Company: PT Restoran Sejahtera
Branch: Cabang Jakarta Pusat
```
**Akses:**
- Proses transaksi penjualan
- Terima pembayaran
- Print receipt
- Lihat menu dan harga

### 7. KITCHEN
```
Email: kitchen@restaurant.com
Password: kitchen123
Role: KITCHEN
Company: PT Restoran Sejahtera
Branch: Cabang Jakarta Pusat
```
**Akses:**
- Lihat order masuk
- Update status masakan
- Kelola inventory bahan
- Komunikasi dengan waiter

### 8. WAITER
```
Email: waiter@restaurant.com
Password: waiter123
Role: WAITER
Company: PT Restoran Sejahtera
Branch: Cabang Jakarta Pusat
```
**Akses:**
- Input order customer
- Lihat menu
- Update status order
- Komunikasi dengan kitchen

---

## 🧪 Testing API

### Login sebagai SUPER_ADMIN
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "superadmin@siresto.com",
    "password": "admin123"
  }'
```

### Login sebagai OWNER
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@restaurant.com",
    "password": "owner123"
  }'
```

### Login sebagai CASHIER
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "cashier@restaurant.com",
    "password": "cashier123"
  }'
```

---

## 📝 Notes

- Semua password menggunakan format: `{role}123`
- Internal users menggunakan domain: `@siresto.com`
- External users menggunakan domain: `@restaurant.com`
- Seeder menggunakan `FirstOrCreate` jadi aman dijalankan berulang kali
- Data tidak akan duplikat jika server di-restart
