# Promo Categories - Visual Diagrams

## Database Schema

```
┌─────────────────────────────────────────────────────────────┐
│                         PROMOS                              │
├─────────────────────────────────────────────────────────────┤
│ id                UUID (PK)                                 │
│ company_id        UUID (FK → companies)                     │
│ branch_id         UUID (FK → branches) [nullable]           │
│ name              VARCHAR(100)                              │
│ code              VARCHAR(50)                               │
│ promo_category    VARCHAR(20) [normal/product/bundle] ◄──┐  │
│ type              VARCHAR(20) [percentage/fixed]          │  │
│ value             DECIMAL(15,2)                           │  │
│ max_discount      DECIMAL(15,2) [nullable]                │  │
│ min_transaction   DECIMAL(15,2) [nullable]                │  │
│ quota             INT [nullable]                           │  │
│ used_count        INT                                      │  │
│ start_date        DATE                                     │  │
│ end_date          DATE                                     │  │
│ is_active         BOOLEAN                                  │  │
│ created_at        TIMESTAMP                                │  │
│ updated_at        TIMESTAMP                                │  │
└─────────────────────────────────────────────────────────────┘  │
                                                                 │
                    ┌────────────────────────────────────────────┤
                    │                                            │
                    ▼                                            │
┌─────────────────────────────────────┐                         │
│       PROMO_PRODUCTS                │                         │
│   (untuk promo category: product)   │                         │
├─────────────────────────────────────┤                         │
│ id          UUID (PK)               │                         │
│ promo_id    UUID (FK → promos) ─────┼─────────────────────────┘
│ product_id  UUID (FK → products)    │
│ created_at  TIMESTAMP               │
│ UNIQUE(promo_id, product_id)        │
└─────────────────────────────────────┘
                    │
                    │
                    ▼
┌─────────────────────────────────────┐
│       PROMO_BUNDLES                 │
│   (untuk promo category: bundle)    │
├─────────────────────────────────────┤
│ id          UUID (PK)               │
│ promo_id    UUID (FK → promos) ─────┼─────────────────────────┐
│ product_id  UUID (FK → products)    │                         │
│ quantity    INT                     │                         │
│ created_at  TIMESTAMP               │                         │
│ UNIQUE(promo_id, product_id)        │                         │
└─────────────────────────────────────┘                         │
                                                                 │
                                                                 │
                    ┌────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│                         PRODUCTS                            │
├─────────────────────────────────────────────────────────────┤
│ id                UUID (PK)                                 │
│ name              VARCHAR(100)                              │
│ sku               VARCHAR(50)                               │
│ price             DECIMAL(15,2)                             │
│ ...                                                         │
└─────────────────────────────────────────────────────────────┘
```

## Promo Category Flow

```
                    ┌─────────────────┐
                    │  Create Promo   │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Pilih Category  │
                    └────────┬────────┘
                             │
            ┌────────────────┼────────────────┐
            │                │                │
            ▼                ▼                ▼
    ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
    │   NORMAL     │  │   PRODUCT    │  │   BUNDLE     │
    └──────┬───────┘  └──────┬───────┘  └──────┬───────┘
           │                 │                  │
           │                 │                  │
           ▼                 ▼                  ▼
    ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
    │ Diskon untuk │  │ Pilih produk │  │ Pilih bundle │
    │ semua produk │  │   tertentu   │  │   items      │
    └──────┬───────┘  └──────┬───────┘  └──────┬───────┘
           │                 │                  │
           │                 ▼                  ▼
           │          ┌──────────────┐  ┌──────────────┐
           │          │ Simpan ke    │  │ Simpan ke    │
           │          │promo_products│  │promo_bundles │
           │          └──────┬───────┘  └──────┬───────┘
           │                 │                  │
           └─────────────────┴──────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │  Promo Created  │
                    └─────────────────┘
```

## Promo Application Flow

```
                    ┌─────────────────┐
                    │  Customer Order │
                    │  with Promo     │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Validate Promo  │
                    │ - Active?       │
                    │ - Not expired?  │
                    │ - Quota OK?     │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Check Category  │
                    └────────┬────────┘
                             │
            ┌────────────────┼────────────────┐
            │                │                │
            ▼                ▼                ▼
    ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
    │   NORMAL     │  │   PRODUCT    │  │   BUNDLE     │
    └──────┬───────┘  └──────┬───────┘  └──────┬───────┘
           │                 │                  │
           ▼                 ▼                  ▼
    ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
    │ Apply to     │  │ Check if     │  │ Check if all │
    │ total order  │  │ order has    │  │ bundle items │
    │              │  │ eligible     │  │ are in order │
    │              │  │ products     │  │              │
    └──────┬───────┘  └──────┬───────┘  └──────┬───────┘
           │                 │                  │
           │                 ▼                  ▼
           │          ┌──────────────┐  ┌──────────────┐
           │          │ Calculate    │  │ Calculate    │
           │          │ discount for │  │ bundle count │
           │          │ eligible     │  │ & discount   │
           │          │ products     │  │              │
           │          └──────┬───────┘  └──────┬───────┘
           │                 │                  │
           └─────────────────┴──────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Apply Discount  │
                    │ to Order        │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Update          │
                    │ promo.used_count│
                    └─────────────────┘
```

## Request/Response Flow

### Promo Normal

```
REQUEST                         RESPONSE
┌─────────────────────┐        ┌─────────────────────┐
│ {                   │        │ {                   │
│   "promo_category": │        │   "id": "...",      │
│     "normal",       │   ──►  │   "promo_category": │
│   "name": "...",    │        │     "normal",       │
│   "code": "...",    │        │   "name": "...",    │
│   "type": "...",    │        │   "code": "...",    │
│   "value": 20       │        │   "value": 20       │
│ }                   │        │ }                   │
└─────────────────────┘        └─────────────────────┘
```

### Promo Product

```
REQUEST                         RESPONSE
┌─────────────────────┐        ┌─────────────────────┐
│ {                   │        │ {                   │
│   "promo_category": │        │   "id": "...",      │
│     "product",      │        │   "promo_category": │
│   "name": "...",    │   ──►  │     "product",      │
│   "code": "...",    │        │   "products": [     │
│   "product_ids": [  │        │     {               │
│     "uuid1",        │        │       "product_id": │
│     "uuid2"         │        │         "uuid1",    │
│   ]                 │        │       "name": "..." │
│ }                   │        │     }               │
└─────────────────────┘        │   ]                 │
                               │ }                   │
                               └─────────────────────┘
```

### Promo Bundle

```
REQUEST                         RESPONSE
┌─────────────────────┐        ┌─────────────────────┐
│ {                   │        │ {                   │
│   "promo_category": │        │   "id": "...",      │
│     "bundle",       │        │   "promo_category": │
│   "name": "...",    │   ──►  │     "bundle",       │
│   "bundle_items": [ │        │   "bundle_items": [ │
│     {               │        │     {               │
│       "product_id": │        │       "product_id": │
│         "uuid1",    │        │         "uuid1",    │
│       "quantity": 1 │        │       "name": "...",│
│     },              │        │       "quantity": 1 │
│     {               │        │     },              │
│       "product_id": │        │     {               │
│         "uuid2",    │        │       "product_id": │
│       "quantity": 2 │        │         "uuid2",    │
│     }               │        │       "name": "...",│
│   ]                 │        │       "quantity": 2 │
│ }                   │        │     }               │
└─────────────────────┘        │   ]                 │
                               │ }                   │
                               └─────────────────────┘
```

## Layer Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         HANDLER                             │
│                  (promo_handler.go)                         │
│  - CreatePromo()                                            │
│  - UpdatePromo()                                            │
│  - GetPromo()                                               │
│  - DeletePromo()                                            │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                         SERVICE                             │
│                  (promo_service.go)                         │
│  - CreatePromo()                                            │
│    ├─ Validate promo_category                              │
│    ├─ Validate product_ids (if product)                    │
│    ├─ Validate bundle_items (if bundle)                    │
│    ├─ Create promo                                          │
│    ├─ Create promo_products (if product)                   │
│    └─ Create promo_bundles (if bundle)                     │
│                                                             │
│  - UpdatePromo()                                            │
│    ├─ Update promo fields                                  │
│    ├─ Delete old promo_products                            │
│    ├─ Create new promo_products                            │
│    ├─ Delete old promo_bundles                             │
│    └─ Create new promo_bundles                             │
│                                                             │
│  - toResponse()                                             │
│    ├─ Build basic response                                 │
│    ├─ Add products (if product category)                   │
│    └─ Add bundle_items (if bundle category)                │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                       REPOSITORY                            │
│                (promo_repository.go)                        │
│  - Create()                                                 │
│  - Update()                                                 │
│  - Delete()                                                 │
│  - FindByID() [with Preload]                               │
│  - FindByCode() [with Preload]                             │
│  - CreatePromoProducts()                                    │
│  - CreatePromoBundles()                                     │
│  - DeletePromoProducts()                                    │
│  - DeletePromoBundles()                                     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                        DATABASE                             │
│  - promos                                                   │
│  - promo_products                                           │
│  - promo_bundles                                            │
│  - products                                                 │
└─────────────────────────────────────────────────────────────┘
```

## Use Case Examples

### Use Case 1: Flash Sale (Normal)
```
Scenario: Diskon 20% untuk semua produk
┌──────────────────────────────────────┐
│ Customer membeli:                    │
│ - Laptop: Rp 5.000.000               │
│ - Mouse: Rp 200.000                  │
│ - Keyboard: Rp 300.000               │
├──────────────────────────────────────┤
│ Subtotal: Rp 5.500.000               │
│ Promo (20%): -Rp 1.100.000           │
│ Total: Rp 4.400.000                  │
└──────────────────────────────────────┘
```

### Use Case 2: Product Specific (Product)
```
Scenario: Diskon 50% untuk Laptop saja
┌──────────────────────────────────────┐
│ Customer membeli:                    │
│ - Laptop: Rp 5.000.000 ✓ eligible   │
│ - Mouse: Rp 200.000 ✗ not eligible  │
│ - Keyboard: Rp 300.000 ✗ not eligible│
├──────────────────────────────────────┤
│ Subtotal: Rp 5.500.000               │
│ Eligible: Rp 5.000.000               │
│ Promo (50%): -Rp 2.500.000           │
│ Total: Rp 3.000.000                  │
└──────────────────────────────────────┘
```

### Use Case 3: Bundle Deal (Bundle)
```
Scenario: Beli Laptop + Mouse x2 dapat diskon Rp 1.000.000
┌──────────────────────────────────────┐
│ Customer membeli:                    │
│ - Laptop x1: Rp 5.000.000 ✓          │
│ - Mouse x2: Rp 400.000 ✓             │
│ - Keyboard x1: Rp 300.000            │
├──────────────────────────────────────┤
│ Subtotal: Rp 5.700.000               │
│ Bundle Complete: YES                 │
│ Promo (fixed): -Rp 1.000.000         │
│ Total: Rp 4.700.000                  │
└──────────────────────────────────────┘
```

## State Diagram

```
                    ┌─────────────┐
                    │   CREATED   │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   ACTIVE    │◄──────┐
                    └──────┬──────┘       │
                           │              │
                ┌──────────┼──────────┐   │
                │          │          │   │
                ▼          ▼          ▼   │
         ┌──────────┐ ┌────────┐ ┌────────┐
         │  USED    │ │EXPIRED │ │INACTIVE│
         └──────────┘ └────────┘ └───┬────┘
                                     │
                                     └──────┘
                                   (can be reactivated)
```

## Validation Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    VALIDATION CHECKS                        │
└─────────────────────────────────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│   NORMAL     │   │   PRODUCT    │   │   BUNDLE     │
└──────┬───────┘   └──────┬───────┘   └──────┬───────┘
       │                  │                  │
       ▼                  ▼                  ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│ ✓ Active     │   │ ✓ Active     │   │ ✓ Active     │
│ ✓ Not expired│   │ ✓ Not expired│   │ ✓ Not expired│
│ ✓ Quota OK   │   │ ✓ Quota OK   │   │ ✓ Quota OK   │
│ ✓ Min trans  │   │ ✓ Min trans  │   │ ✓ Min trans  │
│              │   │ ✓ Has product│   │ ✓ Has bundle │
│              │   │   in order   │   │   complete   │
└──────┬───────┘   └──────┬───────┘   └──────┬───────┘
       │                  │                  │
       └──────────────────┴──────────────────┘
                          │
                          ▼
                   ┌──────────────┐
                   │ APPLY PROMO  │
                   └──────────────┘
```

---

**Note**: Diagram ini membantu visualisasi struktur dan flow dari sistem promo categories.
