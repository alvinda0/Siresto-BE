# Tax API Documentation

API untuk mengelola data pajak (Tax) dalam sistem SIRESTO.

## Base URL
```
/api/v1/external/tax
```

## Authentication
Semua endpoint memerlukan:
- Header: `Authorization: Bearer <token>`
- Role: EXTERNAL (OWNER, ADMIN, CASHIER, dll)

## Endpoints

### 1. Create Tax
Membuat data pajak baru.

**Endpoint:** `POST /api/v1/external/tax`

**Request Body:**
```json
{
  "nama_pajak": "PB1",
  "tipe_pajak": "pb1",
  "presentase": 10.00,
  "deskripsi": "Pajak Barang dan Jasa 1",
  "status": "active",
  "prioritas": 1
}
```

**Field Validations:**
- `nama_pajak`: required, string
- `tipe_pajak`: required, enum (sc, pb1)
- `presentase`: required, float, min: 0, max: 100
- `deskripsi`: optional, string
- `status`: optional, enum (active, inactive), default: "active"
- `prioritas`: optional, integer, default: 0

**Response Success (201):**
```json
{
  "status": "success",
  "message": "Tax created successfully",
  "data": {
    "id": "uuid",
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10.00,
    "deskripsi": "Pajak Barang dan Jasa 1",
    "status": "active",
    "prioritas": 1,
    "created_at": "2024-01-01 10:00:00",
    "updated_at": "2024-01-01 10:00:00"
  }
}
```

---

### 2. Update Tax
Mengupdate data pajak yang sudah ada.

**Endpoint:** `PUT /api/v1/external/tax/:id`

**Request Body:**
```json
{
  "nama_pajak": "Service Charge Updated",
  "tipe_pajak": "sc",
  "presentase": 7.50,
  "deskripsi": "Service charge updated",
  "status": "inactive",
  "prioritas": 2
}
```

**Field Validations:**
- Semua field optional
- `tipe_pajak`: jika diisi, harus enum (sc, pb1)
- `presentase`: jika diisi, min: 0, max: 100
- `status`: jika diisi, harus enum (active, inactive)

**Response Success (200):**
```json
{
  "status": "success",
  "message": "Tax updated successfully",
  "data": {
    "id": "uuid",
    "nama_pajak": "Service Charge Updated",
    "tipe_pajak": "sc",
    "presentase": 7.50,
    "deskripsi": "Service charge updated",
    "status": "inactive",
    "prioritas": 2,
    "created_at": "2024-01-01 10:00:00",
    "updated_at": "2024-01-01 11:00:00"
  }
}
```

---

### 3. Get Tax by ID
Mengambil detail satu data pajak berdasarkan ID.

**Endpoint:** `GET /api/v1/external/tax/:id`

**Response Success (200):**
```json
{
  "status": "success",
  "message": "Tax retrieved successfully",
  "data": {
    "id": "uuid",
    "nama_pajak": "PB1",
    "tipe_pajak": "pb1",
    "presentase": 10.00,
    "deskripsi": "Pajak Barang dan Jasa 1",
    "status": "active",
    "prioritas": 1,
    "created_at": "2024-01-01 10:00:00",
    "updated_at": "2024-01-01 10:00:00"
  }
}
```

---

### 4. Get All Taxes
Mengambil semua data pajak.

**Endpoint:** `GET /api/v1/external/tax`

**Response Success (200):**
```json
{
  "status": "success",
  "message": "Taxes retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "nama_pajak": "PB1",
      "tipe_pajak": "pb1",
      "presentase": 10.00,
      "deskripsi": "Pajak Barang dan Jasa 1",
      "status": "active",
      "prioritas": 1,
      "created_at": "2024-01-01 10:00:00",
      "updated_at": "2024-01-01 10:00:00"
    },
    {
      "id": "uuid",
      "nama_pajak": "Service Charge",
      "tipe_pajak": "sc",
      "presentase": 5.00,
      "deskripsi": "Service charge untuk pelayanan",
      "status": "active",
      "prioritas": 2,
      "created_at": "2024-01-01 10:00:00",
      "updated_at": "2024-01-01 10:00:00"
    }
  ]
}
```

**Note:** Data diurutkan berdasarkan `prioritas DESC, nama_pajak ASC`

---

### 5. Delete Tax
Menghapus data pajak.

**Endpoint:** `DELETE /api/v1/external/tax/:id`

**Response Success (200):**
```json
{
  "status": "success",
  "message": "Tax deleted successfully"
}
```

---

## Error Responses

### 400 Bad Request
```json
{
  "status": "error",
  "message": "Invalid tax ID"
}
```

### 404 Not Found
```json
{
  "status": "error",
  "message": "tax not found"
}
```

### 500 Internal Server Error
```json
{
  "status": "error",
  "message": "error message"
}
```

---

## Tipe Pajak

### sc (Service Charge)
- Biaya layanan yang dikenakan kepada pelanggan
- Biasanya berkisar 5-10%

### pb1 (Pajak Barang dan Jasa 1)
- Pajak yang dikenakan pada barang dan jasa
- Sesuai dengan peraturan perpajakan yang berlaku

---

## Status

- `active`: Pajak aktif dan akan diterapkan
- `inactive`: Pajak tidak aktif, tidak akan diterapkan

---

## Prioritas

- Angka yang menentukan urutan penerapan pajak
- Semakin tinggi angka, semakin tinggi prioritas
- Default: 0

---

## Database Schema

```sql
CREATE TABLE taxes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_pajak varchar(100) NOT NULL,
    tipe_pajak varchar(10) NOT NULL,
    presentase decimal(5,2) NOT NULL,
    deskripsi text,
    status varchar(20) DEFAULT 'active',
    prioritas integer DEFAULT 0,
    created_at timestamptz,
    updated_at timestamptz
);

CREATE INDEX idx_taxes_status ON taxes(status);
CREATE INDEX idx_taxes_prioritas ON taxes(prioritas);
```
