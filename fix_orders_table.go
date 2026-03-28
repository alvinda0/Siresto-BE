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

	// Drop existing tables if they exist
	fmt.Println("\nDropping existing orders tables (if any)...")
	db.Exec("DROP TABLE IF EXISTS order_items CASCADE")
	db.Exec("DROP TABLE IF EXISTS orders CASCADE")
	fmt.Println("✓ Old tables dropped")

	// Create orders table with correct structure
	fmt.Println("\nCreating orders table...")
	if err := db.Exec(`
		CREATE TABLE orders (
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
	fmt.Println("Creating orders indexes...")
	db.Exec("CREATE INDEX idx_orders_company_id ON orders(company_id)")
	db.Exec("CREATE INDEX idx_orders_branch_id ON orders(branch_id)")
	db.Exec("CREATE INDEX idx_orders_status ON orders(status)")
	db.Exec("CREATE INDEX idx_orders_deleted_at ON orders(deleted_at)")
	fmt.Println("✓ Orders indexes created")

	// Create order_items table
	fmt.Println("\nCreating order_items table...")
	if err := db.Exec(`
		CREATE TABLE order_items (
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
	fmt.Println("Creating order_items indexes...")
	db.Exec("CREATE INDEX idx_order_items_order_id ON order_items(order_id)")
	db.Exec("CREATE INDEX idx_order_items_product_id ON order_items(product_id)")
	db.Exec("CREATE INDEX idx_order_items_deleted_at ON order_items(deleted_at)")
	fmt.Println("✓ Order_items indexes created")

	// Verify tables
	fmt.Println("\nVerifying tables...")
	var ordersExists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'orders')").Scan(&ordersExists)
	if ordersExists {
		fmt.Println("✓ Orders table exists")
	}

	var orderItemsExists bool
	db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'order_items')").Scan(&orderItemsExists)
	if orderItemsExists {
		fmt.Println("✓ Order_items table exists")
	}

	fmt.Println("\n✅ Migration completed successfully!")
	fmt.Println("\nYou can now:")
	fmt.Println("1. Start the server: ./server.exe")
	fmt.Println("2. Test create order:")
	fmt.Println("   POST http://localhost:8080/api/v1/external/orders")
}
