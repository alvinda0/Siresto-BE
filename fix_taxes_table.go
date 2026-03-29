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

	// Drop existing taxes table
	log.Println("Dropping existing taxes table...")
	if err := db.Exec("DROP TABLE IF EXISTS taxes CASCADE").Error; err != nil {
		log.Fatal("Failed to drop taxes table:", err)
	}
	log.Println("✓ Taxes table dropped")

	// Create new taxes table with correct schema
	log.Println("Creating new taxes table...")
	if err := db.Exec(`
		CREATE TABLE taxes (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			company_id uuid NOT NULL,
			branch_id uuid,
			nama_pajak varchar(100) NOT NULL,
			tipe_pajak varchar(10) NOT NULL CHECK (tipe_pajak IN ('sc', 'pb1')),
			presentase decimal(5,2) NOT NULL CHECK (presentase >= 0 AND presentase <= 100),
			deskripsi text,
			status varchar(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
			prioritas integer DEFAULT 0,
			created_at timestamptz DEFAULT NOW(),
			updated_at timestamptz DEFAULT NOW(),
			CONSTRAINT fk_taxes_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
			CONSTRAINT fk_taxes_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		log.Fatal("Failed to create taxes table:", err)
	}
	log.Println("✓ Taxes table created")

	// Create indexes
	log.Println("Creating indexes...")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_company_id ON taxes(company_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_branch_id ON taxes(branch_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_status ON taxes(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_prioritas ON taxes(prioritas)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_tipe_pajak ON taxes(tipe_pajak)")
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
	
	// Company-level tax (berlaku untuk semua branch)
	result := db.Exec(`
		INSERT INTO taxes (company_id, branch_id, nama_pajak, tipe_pajak, presentase, deskripsi, status, prioritas)
		VALUES ($1, NULL, 'PB1', 'pb1', 10.00, 'Pajak Barang dan Jasa 1 (Company Level)', 'active', 1)
	`, company.ID)
	
	if result.Error != nil {
		log.Printf("⚠ Failed to insert PB1: %v", result.Error)
	} else {
		log.Printf("✓ Inserted: PB1 (Company Level) for company %s", company.ID)
	}

	// Branch-level tax (hanya untuk branch tertentu)
	result = db.Exec(`
		INSERT INTO taxes (company_id, branch_id, nama_pajak, tipe_pajak, presentase, deskripsi, status, prioritas)
		VALUES ($1, $2, 'Service Charge', 'sc', 5.00, 'Biaya layanan (Branch Level)', 'active', 2)
	`, company.ID, branch.ID)
	
	if result.Error != nil {
		log.Printf("⚠ Failed to insert Service Charge: %v", result.Error)
	} else {
		log.Printf("✓ Inserted: Service Charge (Branch Level) for branch %s", branch.ID)
	}

	// Verify
	log.Println("\nVerifying taxes table...")
	var count int64
	db.Table("taxes").Count(&count)
	log.Printf("Total taxes in database: %d", count)

	// Show all taxes
	type Tax struct {
		ID         string  `gorm:"column:id"`
		CompanyID  string  `gorm:"column:company_id"`
		BranchID   *string `gorm:"column:branch_id"`
		NamaPajak  string  `gorm:"column:nama_pajak"`
		TipePajak  string  `gorm:"column:tipe_pajak"`
		Presentase float64 `gorm:"column:presentase"`
		Status     string  `gorm:"column:status"`
		Prioritas  int     `gorm:"column:prioritas"`
	}

	var taxes []Tax
	db.Table("taxes").Order("prioritas DESC, nama_pajak ASC").Find(&taxes)

	log.Println("\nCurrent taxes:")
	for _, tax := range taxes {
		branchInfo := "Company Level"
		if tax.BranchID != nil {
			branchInfo = "Branch: " + *tax.BranchID
		}
		log.Printf("- %s (%s): %.2f%% [%s] Priority: %d (%s)",
			tax.NamaPajak, tax.TipePajak, tax.Presentase, tax.Status, tax.Prioritas, branchInfo)
	}

	log.Println("\n✅ Migration completed successfully!")
	log.Println("\nNotes:")
	log.Println("- Company-level tax (branch_id = NULL) berlaku untuk semua branch")
	log.Println("- Branch-level tax hanya berlaku untuk branch tertentu")
	log.Println("- User akan melihat company-level + branch-level tax mereka")
}
