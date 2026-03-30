# Migration: PREPARING → PROCESSING

## Perubahan

Status order `PREPARING` telah diubah menjadi `PROCESSING` untuk konsistensi penamaan.

## Status Sebelum dan Sesudah

| Sebelum | Sesudah |
|---------|---------|
| PREPARING | PROCESSING |

Status lainnya tetap sama:
- PENDING
- CONFIRMED
- READY
- COMPLETED
- CANCELLED

## Cara Menjalankan Migration

Jika Anda sudah memiliki data order dengan status `PREPARING` di database, jalankan script migration:

### 1. Jalankan Migration Script

```bash
go run update_preparing_to_processing.go
```

Script ini akan:
- Update semua order dengan status `PREPARING` menjadi `PROCESSING`
- Menampilkan jumlah record yang diupdate

### 2. Verifikasi

Cek database untuk memastikan tidak ada lagi status `PREPARING`:

```sql
SELECT status, COUNT(*) 
FROM orders 
GROUP BY status;
```

## Jika Database Masih Kosong

Jika database Anda masih kosong atau belum ada order dengan status `PREPARING`, Anda tidak perlu menjalankan migration. Cukup gunakan status `PROCESSING` untuk order baru.

## Update Aplikasi Client

Jika Anda memiliki aplikasi client (mobile app, web app, dll) yang menggunakan status `PREPARING`, update kode client untuk menggunakan `PROCESSING`:

### Sebelum:
```javascript
if (order.status === 'PREPARING') {
  // ...
}
```

### Sesudah:
```javascript
if (order.status === 'PROCESSING') {
  // ...
}
```

## Rollback (Jika Diperlukan)

Jika perlu rollback ke `PREPARING`:

```sql
UPDATE orders 
SET status = 'PREPARING' 
WHERE status = 'PROCESSING';
```

Dan kembalikan constant di `internal/entity/order.go`:

```go
OrderStatusPreparing OrderStatus = "PREPARING"
```

## Catatan

- Migration ini aman dan tidak mengubah data lain
- Hanya field `status` yang diupdate
- Tidak ada data yang hilang
- Proses reversible (bisa di-rollback)
