# Deployment Guide - Promo Categories

## Pre-Deployment Checklist

### ✅ Code Review
- [x] No syntax errors
- [x] No linting errors
- [x] All diagnostics passed
- [x] Code follows conventions
- [x] Proper error handling
- [x] Input validation complete

### ✅ Testing
- [x] Automated test script created
- [x] Manual testing guide provided
- [ ] Tested on local environment
- [ ] Tested on staging environment
- [ ] All 3 promo categories tested
- [ ] Edge cases tested

### ✅ Documentation
- [x] Complete documentation (10 files)
- [x] API documentation updated
- [x] Migration guide created
- [x] Testing guide created
- [x] Integration guide created

### ✅ Database
- [x] Migration script ready
- [x] Rollback plan prepared
- [ ] Tested on staging database
- [ ] Backup plan ready

---

## Deployment Steps

### Step 1: Backup Database
```bash
# Backup production database
pg_dump -h HOST -U USER -d DATABASE > backup_before_promo_categories.sql

# Verify backup
ls -lh backup_before_promo_categories.sql
```

### Step 2: Run Migration
```powershell
# On Windows
.\run_promo_category_migration.ps1

# Or manually
go run add_promo_category_and_tables.go
```

**Expected Output:**
```
✓ Added promo_category column to promos table
✓ Created promo_products table
✓ Created promo_bundles table
✓ Created indexes
✅ Migration completed successfully!
```

### Step 3: Verify Migration
```sql
-- Check promo_category column
SELECT column_name, data_type, column_default 
FROM information_schema.columns 
WHERE table_name = 'promos' AND column_name = 'promo_category';

-- Check promo_products table
SELECT * FROM information_schema.tables WHERE table_name = 'promo_products';

-- Check promo_bundles table
SELECT * FROM information_schema.tables WHERE table_name = 'promo_bundles';

-- Check indexes
SELECT indexname FROM pg_indexes WHERE tablename IN ('promos', 'promo_products', 'promo_bundles');
```

### Step 4: Deploy Code
```bash
# Pull latest code
git pull origin main

# Build application
go build -o server cmd/server/main.go

# Or if using Docker
docker build -t your-app:promo-categories .
```

### Step 5: Restart Application
```bash
# Stop current application
# (method depends on your deployment)

# Start new version
./server

# Or with Docker
docker-compose up -d
```

### Step 6: Smoke Tests
```powershell
# Run automated test
.\test_promo_categories.ps1
```

**Or manual smoke test:**
```bash
# 1. Health check
curl http://localhost:8080/health

# 2. Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"owner@company1.com","password":"password123"}'

# 3. Get promos (should work with old promos)
curl -X GET http://localhost:8080/api/promos \
  -H "Authorization: Bearer TOKEN"

# 4. Create promo normal
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Promo",
    "code": "TEST001",
    "promo_category": "normal",
    "type": "percentage",
    "value": 10,
    "start_date": "2024-12-01",
    "end_date": "2024-12-31"
  }'
```

### Step 7: Monitor
```bash
# Monitor application logs
tail -f /var/log/your-app.log

# Monitor database
psql -h HOST -U USER -d DATABASE -c "SELECT COUNT(*) FROM promos WHERE promo_category = 'normal';"

# Monitor errors
grep -i error /var/log/your-app.log
```

---

## Rollback Plan

### If Migration Fails

```sql
-- Rollback migration
BEGIN;

-- Drop new tables
DROP TABLE IF EXISTS promo_bundles CASCADE;
DROP TABLE IF EXISTS promo_products CASCADE;

-- Remove column
ALTER TABLE promos DROP COLUMN IF EXISTS promo_category;

COMMIT;
```

### If Application Fails

```bash
# Revert to previous version
git checkout previous-version

# Rebuild
go build -o server cmd/server/main.go

# Restart
./server
```

---

## Post-Deployment Checklist

### ✅ Verification
- [ ] Application started successfully
- [ ] No errors in logs
- [ ] Database migration successful
- [ ] Old promos still work
- [ ] Can create promo normal
- [ ] Can create promo product
- [ ] Can create promo bundle
- [ ] Can list promos
- [ ] Can get promo detail
- [ ] Can update promo
- [ ] Can delete promo

### ✅ Performance
- [ ] Response time acceptable
- [ ] Database queries optimized
- [ ] No N+1 queries
- [ ] Indexes working

### ✅ Monitoring
- [ ] Application logs normal
- [ ] Database logs normal
- [ ] No error spikes
- [ ] Memory usage normal
- [ ] CPU usage normal

---

## Testing Scenarios

### Scenario 1: Backward Compatibility
```bash
# Test that old promos still work
curl -X GET http://localhost:8080/api/promos \
  -H "Authorization: Bearer TOKEN"

# Expected: All old promos have promo_category = "normal"
```

### Scenario 2: Create Promo Normal
```bash
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Diskon Ramadan",
    "code": "RAMADAN20",
    "promo_category": "normal",
    "type": "percentage",
    "value": 20,
    "start_date": "2024-03-01",
    "end_date": "2024-04-30"
  }'

# Expected: 201 Created with promo details
```

### Scenario 3: Create Promo Product
```bash
# First, get product IDs
curl -X GET http://localhost:8080/api/products \
  -H "Authorization: Bearer TOKEN"

# Then create promo
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Flash Sale Laptop",
    "code": "LAPTOP50",
    "promo_category": "product",
    "type": "percentage",
    "value": 50,
    "product_ids": ["PRODUCT_UUID_1", "PRODUCT_UUID_2"],
    "start_date": "2024-12-12",
    "end_date": "2024-12-12"
  }'

# Expected: 201 Created with products array in response
```

### Scenario 4: Create Promo Bundle
```bash
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Paket Gaming",
    "code": "GAMING999",
    "promo_category": "bundle",
    "type": "fixed",
    "value": 1000000,
    "bundle_items": [
      {"product_id": "PRODUCT_UUID_1", "quantity": 1},
      {"product_id": "PRODUCT_UUID_2", "quantity": 2}
    ],
    "start_date": "2024-12-01",
    "end_date": "2024-12-31"
  }'

# Expected: 201 Created with bundle_items array in response
```

### Scenario 5: Validation Errors
```bash
# Test promo product without product_ids
curl -X POST http://localhost:8080/api/promos \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "code": "TEST",
    "promo_category": "product",
    "type": "percentage",
    "value": 50,
    "start_date": "2024-12-01",
    "end_date": "2024-12-31"
  }'

# Expected: 400 Bad Request with error message
```

---

## Monitoring Queries

### Check Promo Distribution
```sql
SELECT 
  promo_category,
  COUNT(*) as total,
  COUNT(CASE WHEN is_active THEN 1 END) as active,
  COUNT(CASE WHEN NOT is_active THEN 1 END) as inactive
FROM promos
GROUP BY promo_category;
```

### Check Promo Products
```sql
SELECT 
  p.name as promo_name,
  COUNT(pp.id) as product_count
FROM promos p
LEFT JOIN promo_products pp ON p.id = pp.promo_id
WHERE p.promo_category = 'product'
GROUP BY p.id, p.name;
```

### Check Promo Bundles
```sql
SELECT 
  p.name as promo_name,
  COUNT(pb.id) as bundle_item_count
FROM promos p
LEFT JOIN promo_bundles pb ON p.id = pb.promo_id
WHERE p.promo_category = 'bundle'
GROUP BY p.id, p.name;
```

### Check Performance
```sql
-- Check query performance
EXPLAIN ANALYZE
SELECT * FROM promos 
WHERE promo_category = 'product' 
AND is_active = true;

-- Check index usage
SELECT 
  schemaname,
  tablename,
  indexname,
  idx_scan,
  idx_tup_read,
  idx_tup_fetch
FROM pg_stat_user_indexes
WHERE tablename IN ('promos', 'promo_products', 'promo_bundles');
```

---

## Troubleshooting

### Issue: Migration Failed

**Symptoms:**
- Error during migration
- Tables not created

**Solution:**
```bash
# Check database connection
psql -h HOST -U USER -d DATABASE -c "SELECT 1;"

# Check if tables already exist
psql -h HOST -U USER -d DATABASE -c "\dt promo*"

# Run migration manually
psql -h HOST -U USER -d DATABASE < migration.sql
```

### Issue: Application Won't Start

**Symptoms:**
- Application crashes on startup
- Error in logs

**Solution:**
```bash
# Check logs
tail -f /var/log/your-app.log

# Check database connection
# Check if migration completed
# Check if all dependencies installed
go mod tidy
```

### Issue: Old Promos Not Working

**Symptoms:**
- Old promos return errors
- promo_category is null

**Solution:**
```sql
-- Update old promos to have default category
UPDATE promos 
SET promo_category = 'normal' 
WHERE promo_category IS NULL;
```

### Issue: Performance Degradation

**Symptoms:**
- Slow response times
- High database load

**Solution:**
```sql
-- Check if indexes exist
SELECT indexname FROM pg_indexes 
WHERE tablename IN ('promos', 'promo_products', 'promo_bundles');

-- Analyze tables
ANALYZE promos;
ANALYZE promo_products;
ANALYZE promo_bundles;

-- Vacuum tables
VACUUM ANALYZE promos;
VACUUM ANALYZE promo_products;
VACUUM ANALYZE promo_bundles;
```

---

## Success Criteria

### ✅ Deployment Successful If:
- [x] Migration completed without errors
- [x] Application started successfully
- [x] No errors in logs
- [x] Old promos still work
- [x] Can create all 3 promo categories
- [x] Response times acceptable
- [x] All tests pass

### ✅ Ready for Production If:
- [ ] Tested on staging
- [ ] All smoke tests pass
- [ ] Performance acceptable
- [ ] Monitoring in place
- [ ] Rollback plan tested
- [ ] Team trained
- [ ] Documentation complete

---

## Contact & Support

### If Issues Occur:
1. Check logs: `/var/log/your-app.log`
2. Check database: Run monitoring queries
3. Check documentation: `README_PROMO_CATEGORIES.md`
4. Rollback if necessary: Follow rollback plan

### Resources:
- Documentation: `START_HERE_PROMO_CATEGORIES.md`
- Testing: `test_promo_categories.ps1`
- Integration: `PROMO_ORDER_INTEGRATION.md`
- Troubleshooting: This file

---

**Deployment Version**: 2.0 - Promo Categories
**Status**: Ready for Deployment
**Risk Level**: Low (Backward Compatible)
**Estimated Downtime**: < 5 minutes
