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
	
	// Create categories table
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			company_id uuid NOT NULL,
			branch_id uuid,
			name text NOT NULL,
			description text,
			position integer NOT NULL DEFAULT 1,
			is_active boolean DEFAULT true,
			created_at timestamptz,
			updated_at timestamptz,
			CONSTRAINT fk_categories_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
			CONSTRAINT fk_categories_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		log.Fatal("Failed to create categories table:", err)
	}
	
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_categories_company_id ON categories(company_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_categories_branch_id ON categories(branch_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_categories_position ON categories(position)")
	
	// Create products table
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			company_id uuid NOT NULL,
			branch_id uuid NOT NULL,
			category_id uuid NOT NULL,
			image text,
			name text NOT NULL,
			description text,
			stock integer DEFAULT 0,
			price numeric(15,2) NOT NULL,
			position text,
			is_available boolean DEFAULT true,
			created_at timestamptz,
			updated_at timestamptz,
			deleted_at timestamptz,
			CONSTRAINT fk_products_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
			CONSTRAINT fk_products_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE,
			CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		log.Fatal("Failed to create products table:", err)
	}
	
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_products_company_id ON products(company_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_products_branch_id ON products(branch_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at)")
	
	// Create api_logs table
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS api_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			method VARCHAR(10) NOT NULL,
			path VARCHAR(255) NOT NULL,
			status_code INTEGER NOT NULL,
			response_time BIGINT NOT NULL,
			ip_address VARCHAR(45),
			user_agent TEXT,
			access_from VARCHAR(50),
			user_id UUID,
			company_id UUID,
			branch_id UUID,
			request_body TEXT,
			response_body TEXT,
			error_message TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		log.Fatal("Failed to create api_logs table:", err)
	}
	
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_api_logs_method ON api_logs(method)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_api_logs_path ON api_logs(path)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_api_logs_user_id ON api_logs(user_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_api_logs_company_id ON api_logs(company_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_api_logs_branch_id ON api_logs(branch_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_api_logs_access_from ON api_logs(access_from)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_api_logs_created_at ON api_logs(created_at)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_api_logs_deleted_at ON api_logs(deleted_at)")
	
	// Create orders table
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			company_id uuid NOT NULL,
			branch_id uuid NOT NULL,
			customer_name text,
			customer_phone text,
			table_number text,
			notes text,
			referral_code text,
			order_method text NOT NULL,
			promo_code text,
			status text DEFAULT 'PENDING',
			total_amount numeric(15,2) DEFAULT 0,
			created_at timestamptz,
			updated_at timestamptz,
			deleted_at timestamptz,
			CONSTRAINT fk_orders_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
			CONSTRAINT fk_orders_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		log.Fatal("Failed to create orders table:", err)
	}
	
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_orders_company_id ON orders(company_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_orders_branch_id ON orders(branch_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders(deleted_at)")
	
	// Create order_items table
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			order_id uuid NOT NULL,
			product_id uuid NOT NULL,
			quantity integer NOT NULL,
			price numeric(15,2) NOT NULL,
			note text,
			created_at timestamptz,
			updated_at timestamptz,
			deleted_at timestamptz,
			CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
			CONSTRAINT fk_order_items_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		log.Fatal("Failed to create order_items table:", err)
	}
	
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_order_items_deleted_at ON order_items(deleted_at)")
	
	// Create taxes table
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS taxes (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			company_id uuid NOT NULL,
			branch_id uuid,
			nama_pajak varchar(100) NOT NULL,
			tipe_pajak varchar(10) NOT NULL,
			presentase decimal(5,2) NOT NULL,
			deskripsi text,
			status varchar(20) DEFAULT 'active',
			prioritas integer DEFAULT 0,
			created_at timestamptz DEFAULT NOW(),
			updated_at timestamptz DEFAULT NOW(),
			CONSTRAINT fk_taxes_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
			CONSTRAINT fk_taxes_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		log.Fatal("Failed to create taxes table:", err)
	}
	
	// Add missing columns if table already exists
	DB.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='taxes' AND column_name='company_id') THEN
				ALTER TABLE taxes ADD COLUMN company_id uuid;
				ALTER TABLE taxes ADD CONSTRAINT fk_taxes_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE;
			END IF;
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='taxes' AND column_name='branch_id') THEN
				ALTER TABLE taxes ADD COLUMN branch_id uuid;
				ALTER TABLE taxes ADD CONSTRAINT fk_taxes_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE;
			END IF;
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='taxes' AND column_name='status') THEN
				ALTER TABLE taxes ADD COLUMN status varchar(20) DEFAULT 'active';
			END IF;
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='taxes' AND column_name='prioritas') THEN
				ALTER TABLE taxes ADD COLUMN prioritas integer DEFAULT 0;
			END IF;
		END $$;
	`)
	
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_company_id ON taxes(company_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_branch_id ON taxes(branch_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_status ON taxes(status)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_prioritas ON taxes(prioritas)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_taxes_tipe_pajak ON taxes(tipe_pajak)")
	
	// Create promos table
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS promos (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			company_id uuid NOT NULL,
			branch_id uuid,
			name varchar(100) NOT NULL,
			code varchar(50) NOT NULL,
			type varchar(20) NOT NULL,
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
	
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_promos_company_id ON promos(company_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_promos_branch_id ON promos(branch_id)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_promos_code ON promos(code)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_promos_is_active ON promos(is_active)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_promos_start_date ON promos(start_date)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_promos_end_date ON promos(end_date)")
	
	log.Println("Database migrated successfully")
}
