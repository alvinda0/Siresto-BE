# Flow Diagram: Perhitungan Pajak Order

## Flow Create Order

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Client Request                                           │
│    POST /api/v1/external/orders                             │
│    {                                                         │
│      "order_items": [{"product_id": "...", "quantity": 2}]  │
│    }                                                         │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Order Handler                                            │
│    - Validate request                                       │
│    - Get company_id & branch_id from context                │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Order Service - CreateOrder()                            │
│    - Validate products                                      │
│    - Calculate subtotal from items                          │
│      subtotal = Σ(price × quantity)                         │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Calculate Taxes - calculateTaxes()                       │
│    - Get active taxes (status = 'active')                   │
│    - Order by priority DESC                                 │
│                                                              │
│    current_amount = subtotal                                │
│    total_tax = 0                                            │
│                                                              │
│    For each tax (by priority):                              │
│      tax_amount = current_amount × (percentage / 100)       │
│      total_tax += tax_amount                                │
│      current_amount += tax_amount                           │
│                                                              │
│    Return: total_tax, tax_details[]                         │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. Save to Database                                         │
│    Order {                                                  │
│      subtotal_amount: 100000                                │
│      tax_amount: 15500                                      │
│      total_amount: 115500                                   │
│    }                                                         │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 6. Format Response - toOrderResponse()                      │
│    - Recalculate tax_details for breakdown                  │
│    - Return complete order with tax breakdown               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 7. Response to Client                                       │
│    {                                                         │
│      "subtotal_amount": 100000,                             │
│      "tax_amount": 15500,                                   │
│      "total_amount": 115500,                                │
│      "tax_details": [...]                                   │
│    }                                                         │
└─────────────────────────────────────────────────────────────┘
```

## Flow Get Order By ID

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Client Request                                           │
│    GET /api/v1/external/orders/{id}                         │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Order Handler                                            │
│    - Parse order ID                                         │
│    - Get company_id & branch_id from context                │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Order Service - GetOrderByID()                           │
│    - Find order from database                               │
│    - Check access control                                   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Format Response - toOrderResponse()                      │
│    - Get order data from DB:                                │
│      • subtotal_amount: 100000 (from DB)                    │
│      • tax_amount: 15500 (from DB)                          │
│      • total_amount: 115500 (from DB)                       │
│                                                              │
│    - Recalculate tax_details for breakdown:                 │
│      calculateTaxes(subtotal_amount, company_id, branch_id) │
│                                                              │
│    - Return complete order with tax breakdown               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. Response to Client                                       │
│    {                                                         │
│      "subtotal_amount": 100000,  ← from DB                  │
│      "tax_amount": 15500,        ← from DB                  │
│      "total_amount": 115500,     ← from DB                  │
│      "tax_details": [...]        ← recalculated             │
│    }                                                         │
└─────────────────────────────────────────────────────────────┘
```

## Contoh Perhitungan Bertingkat

```
Input:
  Subtotal: Rp 100.000
  Pajak 1: PB1 10% (prioritas 1)
  Pajak 2: Service Charge 5% (prioritas 2)

Step-by-step:

┌──────────────────────────────────────────────────────────┐
│ Initial State                                            │
│   current_amount = 100.000                               │
│   total_tax = 0                                          │
└──────────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────────┐
│ Tax 1: PB1 (10%, priority 1)                             │
│   base_amount = 100.000                                  │
│   tax_amount = 100.000 × 10% = 10.000                    │
│   current_amount = 100.000 + 10.000 = 110.000           │
│   total_tax = 0 + 10.000 = 10.000                        │
└──────────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────────┐
│ Tax 2: Service Charge (5%, priority 2)                   │
│   base_amount = 110.000  ← includes previous tax         │
│   tax_amount = 110.000 × 5% = 5.500                      │
│   current_amount = 110.000 + 5.500 = 115.500            │
│   total_tax = 10.000 + 5.500 = 15.500                    │
└──────────────────────────────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────────┐
│ Final Result                                             │
│   subtotal_amount: 100.000                               │
│   tax_amount: 15.500                                     │
│   total_amount: 115.500                                  │
│                                                           │
│   tax_details: [                                         │
│     {                                                     │
│       tax_name: "PB1",                                   │
│       base_amount: 100.000,                              │
│       tax_amount: 10.000                                 │
│     },                                                    │
│     {                                                     │
│       tax_name: "Service Charge",                        │
│       base_amount: 110.000,                              │
│       tax_amount: 5.500                                  │
│     }                                                     │
│   ]                                                       │
└──────────────────────────────────────────────────────────┘
```

## Database Schema

```
┌─────────────────────────────────────────────────────────┐
│ orders table                                            │
├─────────────────────────────────────────────────────────┤
│ id                UUID PRIMARY KEY                      │
│ company_id        UUID NOT NULL                         │
│ branch_id         UUID NOT NULL                         │
│ customer_name     VARCHAR(255)                          │
│ table_number      VARCHAR(50)                           │
│ order_method      VARCHAR(50)                           │
│ status            VARCHAR(50)                           │
│ subtotal_amount   DECIMAL(15,2) ← Total items           │
│ tax_amount        DECIMAL(15,2) ← Total taxes           │
│ total_amount      DECIMAL(15,2) ← Subtotal + Tax        │
│ created_at        TIMESTAMP                             │
│ updated_at        TIMESTAMP                             │
└─────────────────────────────────────────────────────────┘
                         │
                         │ 1:N
                         ▼
┌─────────────────────────────────────────────────────────┐
│ order_items table                                       │
├─────────────────────────────────────────────────────────┤
│ id                UUID PRIMARY KEY                      │
│ order_id          UUID FOREIGN KEY → orders.id          │
│ product_id        UUID FOREIGN KEY → products.id        │
│ quantity          INT                                   │
│ price             DECIMAL(15,2)                         │
│ note              TEXT                                  │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│ taxes table                                             │
├─────────────────────────────────────────────────────────┤
│ id                UUID PRIMARY KEY                      │
│ company_id        UUID NOT NULL                         │
│ branch_id         UUID (nullable)                       │
│ nama_pajak        VARCHAR(100)                          │
│ tipe_pajak        VARCHAR(10)                           │
│ presentase        DECIMAL(5,2)                          │
│ prioritas         INT ← Determines calculation order    │
│ status            VARCHAR(20) ← 'active' or 'inactive'  │
│ deskripsi         TEXT                                  │
└─────────────────────────────────────────────────────────┘
```

## Key Points

1. **Subtotal, Tax, Total disimpan di database** - Nilai tidak berubah meskipun konfigurasi pajak berubah
2. **Tax details dihitung ulang** - Setiap kali order ditampilkan, breakdown dihitung ulang dari pajak aktif
3. **Prioritas menentukan urutan** - Priority 1 dihitung pertama, Priority 2 kedua, dst.
4. **Perhitungan kumulatif** - Setiap pajak dihitung dari base yang sudah termasuk pajak sebelumnya
5. **Hanya pajak aktif** - Hanya pajak dengan status 'active' yang dihitung
