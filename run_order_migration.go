package main

import (
	"fmt"
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

	fmt.Println("Connected to database successfully")

	// Create orders table
	fmt.Println("Creating orders table...")
	if err := db.Exec(`
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
			created_at timestamptz DEFAULT NOW(),
			updated_at timestamptz DEFAULT NOW(),
			deleted_at timestamptz,
			CONSTRAINT fk_orders_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
			CONSTRAINT fk_orders_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		log.Fatal("Failed to create orders table:", err)
	}
	fmt.Println("✓ Orders table created")

	// Create indexes for orders
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_company_id ON orders(company_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_branch_id ON orders(branch_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders(deleted_at)")
	fmt.Println("✓ Orders indexes created")

	// Create order_items table
	fmt.Println("Creating order_items table...")
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			order_id uuid NOT NULL,
			product_id uuid NOT NULL,
			quantity integer NOT NULL,
			price numeric(15,2) NOT NULL,
			note text,
			created_at timestamptz DEFAULT NOW(),
			updated_at timestamptz DEFAULT NOW(),
			deleted_at timestamptz,
			CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
			CONSTRAINT fk_order_items_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		log.Fatal("Failed to create order_items table:", err)
	}
	fmt.Println("✓ Order_items table created")

	// Create indexes for order_items
	db.Exec("CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_order_items_deleted_at ON order_items(deleted_at)")
	fmt.Println("✓ Order_items indexes created")

	fmt.Println("\n✅ Migration completed successfully!")
	fmt.Println("\nYou can now:")
	fmt.Println("1. Start the server: ./server.exe")
	fmt.Println("2. Create orders via API")
}
