package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func GetAllowedOrigins() []string {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		return []string{"http://localhost:3000"} // default
	}
	return strings.Split(origins, ",")
}

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")
}

func MigrateDB() {
	// Enable UUID extension for PostgreSQL
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Println("Warning: Failed to create uuid-ossp extension:", err)
	}
	
	// Create tables without dropping (safe for production)
	log.Println("Creating tables if not exists...")
	
	// Create roles table first (no dependencies)
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			name text UNIQUE NOT NULL,
			display_name text NOT NULL,
			type text NOT NULL,
			description text,
			is_active boolean DEFAULT true,
			created_at timestamptz,
			updated_at timestamptz
		)
	`).Error; err != nil {
		log.Fatal("Failed to create roles table:", err)
	}
	
	// Seed default roles
	roles := []struct {
		Name        string
		DisplayName string
		Type        string
		Description string
	}{
		// Internal roles
		{"SUPER_ADMIN", "Super Admin", "INTERNAL", "Owner sistem SIRESTO"},
		{"SUPPORT", "Support", "INTERNAL", "CS / Admin internal"},
		{"FINANCE", "Finance", "INTERNAL", "Lihat pembayaran subscription"},
		// External roles
		{"OWNER", "Owner", "EXTERNAL", "Pemilik usaha restoran"},
		{"ADMIN", "Admin", "EXTERNAL", "Manager cabang"},
		{"CASHIER", "Cashier", "EXTERNAL", "Kasir"},
		{"KITCHEN", "Kitchen", "EXTERNAL", "Dapur"},
		{"WAITER", "Waiter", "EXTERNAL", "Pelayan"},
	}
	
	for _, role := range roles {
		DB.Exec(`
			INSERT INTO roles (name, display_name, type, description, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, true, NOW(), NOW())
			ON CONFLICT (name) DO NOTHING
		`, role.Name, role.DisplayName, role.Type, role.Description)
	}
	log.Println("✓ Roles seeded")
	
	// Create users table first (no dependencies except roles)
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			name text NOT NULL,
			email text UNIQUE NOT NULL,
			password text NOT NULL,
			role_id uuid NOT NULL,
			company_id uuid,
			branch_id uuid,
			is_active boolean DEFAULT true,
			created_at timestamptz,
			updated_at timestamptz,
			CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
		)
	`).Error; err != nil {
		log.Fatal("Failed to create users table:", err)
	}
	
	// Create companies table with foreign key to users
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS companies (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			name text NOT NULL,
			type text NOT NULL,
			owner_id uuid NOT NULL,
			created_at timestamptz,
			updated_at timestamptz,
			CONSTRAINT fk_companies_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		log.Fatal("Failed to create companies table:", err)
	}
	
	// Create branches table with foreign key to companies
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS branches (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			company_id uuid NOT NULL,
			name text NOT NULL,
			address text NOT NULL,
			city text,
			province text,
			postal_code text,
			phone text,
			is_active boolean DEFAULT true,
			created_at timestamptz,
			updated_at timestamptz,
			CONSTRAINT fk_branches_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		log.Fatal("Failed to create branches table:", err)
	}
	
	// Add foreign keys back to users table (if not exists)
	DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_users_company'
			) THEN
				ALTER TABLE users 
				ADD CONSTRAINT fk_users_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;
			END IF;
		END $$;
	`)
	
	DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_users_branch'
			) THEN
				ALTER TABLE users 
				ADD CONSTRAINT fk_users_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE SET NULL;
			END IF;
		END $$;
	`)
	
	// Create indexes
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_branch_id ON users(branch_id)")
	
	log.Println("Database migrated successfully")
}
