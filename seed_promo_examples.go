package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Promo struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CompanyID      uuid.UUID  `gorm:"type:uuid;not null"`
	BranchID       *uuid.UUID `gorm:"type:uuid"`
	Name           string     `gorm:"type:varchar(100);not null"`
	Code           string     `gorm:"type:varchar(50);not null"`
	PromoCategory  string     `gorm:"type:varchar(20);not null;default:'normal'"`
	Type           string     `gorm:"type:varchar(20);not null"`
	Value          float64    `gorm:"type:decimal(15,2);not null"`
	MaxDiscount    *float64   `gorm:"type:decimal(15,2)"`
	MinTransaction *float64   `gorm:"type:decimal(15,2)"`
	Quota          *int       `gorm:"type:int"`
	UsedCount      int        `gorm:"type:int;default:0"`
	StartDate      time.Time  `gorm:"type:date;not null"`
	EndDate        time.Time  `gorm:"type:date;not null"`
	IsActive       bool       `gorm:"type:boolean;default:true"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type PromoProduct struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PromoID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time
}

type PromoBundle struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PromoID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  int       `gorm:"type:int;not null;default:1"`
	CreatedAt time.Time
}

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Starting promo examples seeder...")

	// Get first company
	var companyID uuid.UUID
	if err := db.Raw("SELECT id FROM companies LIMIT 1").Scan(&companyID).Error; err != nil {
		log.Fatal("Failed to get company:", err)
	}
	log.Printf("Using company ID: %s", companyID)

	// Get some products
	var productIDs []uuid.UUID
	if err := db.Raw("SELECT id FROM products LIMIT 5").Scan(&productIDs).Error; err != nil {
		log.Fatal("Failed to get products:", err)
	}

	if len(productIDs) < 3 {
		log.Fatal("Need at least 3 products in database. Please create products first.")
	}

	log.Printf("Found %d products", len(productIDs))

	// 1. Create Promo Normal
	log.Println("\n1. Creating Promo Normal...")
	maxDiscount := 100000.0
	minTransaction := 200000.0
	quota := 100

	promoNormal := Promo{
		CompanyID:      companyID,
		Name:           "Diskon Akhir Tahun 2024",
		Code:           "NEWYEAR2024",
		PromoCategory:  "normal",
		Type:           "percentage",
		Value:          15,
		MaxDiscount:    &maxDiscount,
		MinTransaction: &minTransaction,
		Quota:          &quota,
		UsedCount:      0,
		StartDate:      time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		IsActive:       true,
	}

	if err := db.Create(&promoNormal).Error; err != nil {
		log.Printf("✗ Failed to create promo normal: %v", err)
	} else {
		log.Printf("✓ Created Promo Normal: %s (%s)", promoNormal.Name, promoNormal.ID)
	}

	// 2. Create Promo Product
	log.Println("\n2. Creating Promo Product...")
	maxDiscount2 := 5000000.0
	quota2 := 50

	promoProduct := Promo{
		CompanyID:     companyID,
		Name:          "Flash Sale Elektronik",
		Code:          "FLASH50",
		PromoCategory: "product",
		Type:          "percentage",
		Value:         50,
		MaxDiscount:   &maxDiscount2,
		Quota:         &quota2,
		UsedCount:     0,
		StartDate:     time.Date(2024, 12, 12, 0, 0, 0, 0, time.UTC),
		EndDate:       time.Date(2024, 12, 12, 0, 0, 0, 0, time.UTC),
		IsActive:      true,
	}

	if err := db.Create(&promoProduct).Error; err != nil {
		log.Printf("✗ Failed to create promo product: %v", err)
	} else {
		log.Printf("✓ Created Promo Product: %s (%s)", promoProduct.Name, promoProduct.ID)

		// Add products to promo
		promoProducts := []PromoProduct{
			{PromoID: promoProduct.ID, ProductID: productIDs[0]},
			{PromoID: promoProduct.ID, ProductID: productIDs[1]},
		}

		if err := db.Create(&promoProducts).Error; err != nil {
			log.Printf("✗ Failed to add products to promo: %v", err)
		} else {
			log.Printf("✓ Added %d products to promo", len(promoProducts))
		}
	}

	// 3. Create Promo Bundle
	log.Println("\n3. Creating Promo Bundle...")
	quota3 := 30

	promoBundle := Promo{
		CompanyID:     companyID,
		Name:          "Paket Hemat Gaming",
		Code:          "GAMING999",
		PromoCategory: "bundle",
		Type:          "fixed",
		Value:         1000000,
		Quota:         &quota3,
		UsedCount:     0,
		StartDate:     time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		EndDate:       time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		IsActive:      true,
	}

	if err := db.Create(&promoBundle).Error; err != nil {
		log.Printf("✗ Failed to create promo bundle: %v", err)
	} else {
		log.Printf("✓ Created Promo Bundle: %s (%s)", promoBundle.Name, promoBundle.ID)

		// Add bundle items
		promoBundles := []PromoBundle{
			{PromoID: promoBundle.ID, ProductID: productIDs[0], Quantity: 1},
			{PromoID: promoBundle.ID, ProductID: productIDs[1], Quantity: 2},
			{PromoID: promoBundle.ID, ProductID: productIDs[2], Quantity: 1},
		}

		if err := db.Create(&promoBundles).Error; err != nil {
			log.Printf("✗ Failed to add bundle items: %v", err)
		} else {
			log.Printf("✓ Added %d bundle items to promo", len(promoBundles))
		}
	}

	// 4. Create more examples
	log.Println("\n4. Creating additional examples...")

	// Promo Normal - Ramadan
	maxDiscount4 := 150000.0
	minTransaction4 := 300000.0
	promoRamadan := Promo{
		CompanyID:      companyID,
		Name:           "Diskon Ramadan 2024",
		Code:           "RAMADAN2024",
		PromoCategory:  "normal",
		Type:           "percentage",
		Value:          20,
		MaxDiscount:    &maxDiscount4,
		MinTransaction: &minTransaction4,
		StartDate:      time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC),
		IsActive:       true,
	}
	db.Create(&promoRamadan)
	log.Printf("✓ Created: %s", promoRamadan.Name)

	// Promo Product - Clearance
	promoProduct2 := Promo{
		CompanyID:     companyID,
		Name:          "Clearance Sale",
		Code:          "CLEAR70",
		PromoCategory: "product",
		Type:          "percentage",
		Value:         70,
		StartDate:     time.Date(2024, 12, 20, 0, 0, 0, 0, time.UTC),
		EndDate:       time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		IsActive:      true,
	}
	if err := db.Create(&promoProduct2).Error; err == nil {
		db.Create(&PromoProduct{PromoID: promoProduct2.ID, ProductID: productIDs[2]})
		log.Printf("✓ Created: %s", promoProduct2.Name)
	}

	// Promo Bundle - Office Package
	promoBundle2 := Promo{
		CompanyID:     companyID,
		Name:          "Paket Lengkap Office",
		Code:          "OFFICE500",
		PromoCategory: "bundle",
		Type:          "fixed",
		Value:         500000,
		StartDate:     time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		EndDate:       time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		IsActive:      true,
	}
	if err := db.Create(&promoBundle2).Error; err == nil {
		db.Create(&[]PromoBundle{
			{PromoID: promoBundle2.ID, ProductID: productIDs[1], Quantity: 1},
			{PromoID: promoBundle2.ID, ProductID: productIDs[2], Quantity: 1},
		})
		log.Printf("✓ Created: %s", promoBundle2.Name)
	}

	log.Println("\n✅ Seeding completed successfully!")
	log.Println("\nCreated promos:")
	log.Println("1. Promo Normal: Diskon Akhir Tahun 2024 (15%)")
	log.Println("2. Promo Product: Flash Sale Elektronik (50%)")
	log.Println("3. Promo Bundle: Paket Hemat Gaming (Rp 1.000.000)")
	log.Println("4. Promo Normal: Diskon Ramadan 2024 (20%)")
	log.Println("5. Promo Product: Clearance Sale (70%)")
	log.Println("6. Promo Bundle: Paket Lengkap Office (Rp 500.000)")
	log.Println("\nYou can now test the promos using the API!")
}
