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

	// Create taxes table
	log.Println("Creating taxes table...")
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS taxes (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			nama_pajak varchar(100) NOT NULL,
			tipe_pajak varchar(10) NOT NULL CHECK (tipe_pajak IN ('sc', 'pb1')),
			presentase decimal(5,2) NOT NULL CHECK (presentase >= 0 AND presentase <= 100),
			deskripsi text,
			status varchar(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
			prioritas integer DEFAULT 0,
			created_at timestamptz DEFAULT NOW(),
			updated_at timestamptz DEFAULT NOW()
		)
	`).Error; err != nil {
		log.Fatal("Failed to create taxes table:", err)
	}

	log.Println("✓ Taxes table created")

	// Create indexes
	log.Println("Creating indexes...")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_status ON taxes(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_prioritas ON taxes(prioritas)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_tipe_pajak ON taxes(tipe_pajak)")
	log.Println("✓ Indexes created")

	// Insert sample data
	log.Println("Inserting sample data...")
	sampleData := []struct {
		NamaPajak  string
		TipePajak  string
		Presentase float64
		Deskripsi  string
		Status     string
		Prioritas  int
	}{
		{"PB1", "pb1", 10.00, "Pajak Barang dan Jasa 1", "active", 1},
		{"Service Charge", "sc", 5.00, "Biaya layanan", "active", 2},
	}

	for _, data := range sampleData {
		result := db.Exec(`
			INSERT INTO taxes (nama_pajak, tipe_pajak, presentase, deskripsi, status, prioritas)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT DO NOTHING
		`, data.NamaPajak, data.TipePajak, data.Presentase, data.Deskripsi, data.Status, data.Prioritas)

		if result.Error != nil {
			log.Printf("Warning: Failed to insert %s: %v", data.NamaPajak, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("✓ Inserted: %s", data.NamaPajak)
		} else {
			log.Printf("- Skipped (already exists): %s", data.NamaPajak)
		}
	}

	// Verify
	log.Println("\nVerifying taxes table...")
	var count int64
	db.Table("taxes").Count(&count)
	log.Printf("Total taxes in database: %d", count)

	// Show all taxes
	type Tax struct {
		ID         string  `gorm:"column:id"`
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
		log.Printf("- %s (%s): %.2f%% [%s] Priority: %d",
			tax.NamaPajak, tax.TipePajak, tax.Presentase, tax.Status, tax.Prioritas)
	}

	log.Println("\n✅ Migration completed successfully!")
}
