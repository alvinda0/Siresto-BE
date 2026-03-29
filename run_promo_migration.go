package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to database")

	// Drop existing promos table if exists
	log.Println("Dropping existing promos table...")
	if err := db.Exec("DROP TABLE IF EXISTS promos CASCADE").Error; err != nil {
		log.Fatal("Failed to drop promos table:", err)
	}
	log.Println("✓ Promos table dropped")

	// Create new promos table
	log.Println("Creating new promos table...")
	if err := db.Exec(`
		CREATE TABLE promos (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			company_id uuid NOT NULL,
			branch_id uuid,
			name varchar(100) NOT NULL,
			code varchar(50) NOT NULL,
			type varchar(20) NOT NULL CHECK (type IN ('percentage', 'fixed')),
			value decimal(15,2) NOT NULL,
			max_discount decimal(15,2),
			min_transaction decimal(15,2),
			quota integer,
			used_count integer DEFAULT 0,
			start_date date NOT NULL,
			end_date date NOT NULL,
			is_active boolean DEFAULT true,
			created_at timestamptz DEFAULT NOW(),
			updated_at timestamptz DEFAULT NOW(),
			CONSTRAINT fk_promos_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
			CONSTRAINT fk_promos_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		log.Fatal("Failed to create promos table:", err)
	}
	log.Println("✓ Promos table created")

	// Create indexes
	log.Println("Creating indexes...")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_promos_company_id ON promos(company_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_promos_branch_id ON promos(branch_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_promos_code ON promos(code)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_promos_is_active ON promos(is_active)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_promos_start_date ON promos(start_date)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_promos_end_date ON promos(end_date)")
	log.Println("✓ Indexes created")

	// Get first company and branch for sample data
	type Company struct {
		ID string `gorm:"column:id"`
	}
	type Branch struct {
		ID        string `gorm:"column:id"`
		CompanyID string `gorm:"column:company_id"`
	}

	var company Company
	if err := db.Table("companies").First(&company).Error; err != nil {
		log.Println("⚠ No company found, skipping sample data")
		log.Println("\n✅ Migration completed successfully!")
		return
	}

	var branch Branch
	if err := db.Table("branches").Where("company_id = ?", company.ID).First(&branch).Error; err != nil {
		log.Println("⚠ No branch found, skipping sample data")
		log.Println("\n✅ Migration completed successfully!")
		return
	}

	// Insert sample data
	log.Println("\nInserting sample data...")

	// Company-level promo
	result := db.Exec(`
		INSERT INTO promos (company_id, branch_id, name, code, type, value, max_discount, min_transaction, quota, start_date, end_date, is_active)
		VALUES ($1, NULL, 'Diskon Lebaran', 'LEBARAN2026', 'percentage', 10.00, 50000, 100000, 100, '2026-03-20', '2026-04-10', true)
	`, company.ID)

	if result.Error != nil {
		log.Printf("⚠ Failed to insert company-level promo: %v", result.Error)
	} else {
		log.Printf("✓ Inserted: Diskon Lebaran (Company Level) for company %s", company.ID)
	}

	// Branch-level promo
	result = db.Exec(`
		INSERT INTO promos (company_id, branch_id, name, code, type, value, min_transaction, quota, start_date, end_date, is_active)
		VALUES ($1, $2, 'Promo Branch', 'BRANCH50', 'fixed', 50000, 200000, 50, '2026-03-25', '2026-04-05', true)
	`, company.ID, branch.ID)

	if result.Error != nil {
		log.Printf("⚠ Failed to insert branch-level promo: %v", result.Error)
	} else {
		log.Printf("✓ Inserted: Promo Branch (Branch Level) for branch %s", branch.ID)
	}

	// Verify
	log.Println("\nVerifying promos table...")
	var count int64
	db.Table("promos").Count(&count)
	log.Printf("Total promos in database: %d", count)

	// Show all promos
	type Promo struct {
		ID             string  `gorm:"column:id"`
		CompanyID      string  `gorm:"column:company_id"`
		BranchID       *string `gorm:"column:branch_id"`
		Name           string  `gorm:"column:name"`
		Code           string  `gorm:"column:code"`
		Type           string  `gorm:"column:type"`
		Value          float64 `gorm:"column:value"`
		MaxDiscount    *float64 `gorm:"column:max_discount"`
		MinTransaction *float64 `gorm:"column:min_transaction"`
		Quota          *int    `gorm:"column:quota"`
		UsedCount      int     `gorm:"column:used_count"`
		IsActive       bool    `gorm:"column:is_active"`
	}

	var promos []Promo
	db.Table("promos").Order("created_at DESC").Find(&promos)

	log.Println("\nCurrent promos:")
	for _, promo := range promos {
		branchInfo := "Company Level"
		if promo.BranchID != nil {
			branchInfo = "Branch: " + *promo.BranchID
		}
		log.Printf("- %s (%s): %s %.2f [%s] Quota: %d/%d (%s)",
			promo.Name, promo.Code, promo.Type, promo.Value,
			map[bool]string{true: "active", false: "inactive"}[promo.IsActive],
			promo.UsedCount, *promo.Quota, branchInfo)
	}

	log.Println("\n✅ Migration completed successfully!")
	log.Println("\nNotes:")
	log.Println("- Company-level promo (branch_id = NULL) berlaku untuk semua branch")
	log.Println("- Branch-level promo hanya berlaku untuk branch tertentu")
	log.Println("- Type: 'percentage' atau 'fixed'")
	log.Println("- Promo akan otomatis expired jika melewati end_date")
}
