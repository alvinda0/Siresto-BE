package config

import (
	"log"
	"project-name/internal/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SeedDatabase() {
	log.Println("Starting database seeding...")

	// Hash password default
	hashPassword := func(password string) string {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hashed)
	}

	// Get role IDs
	getRoleID := func(roleName string) uuid.UUID {
		var role entity.Role
		DB.Where("name = ?", roleName).First(&role)
		return role.ID
	}

	// ===== SEED INTERNAL USERS (Platform SIRESTO) =====
	
	// 1. SUPER_ADMIN
	superAdmin := entity.User{
		Name:     "Super Admin",
		Email:    "superadmin@siresto.com",
		Password: hashPassword("admin123"),
		RoleID:   getRoleID("SUPER_ADMIN"),
		IsActive: true,
	}
	if err := DB.FirstOrCreate(&superAdmin, entity.User{Email: "superadmin@siresto.com"}).Error; err != nil {
		log.Printf("Error seeding SUPER_ADMIN: %v", err)
	} else {
		log.Println("✓ SUPER_ADMIN created: superadmin@siresto.com / admin123")
	}

	// 2. SUPPORT
	support := entity.User{
		Name:     "CS Support",
		Email:    "support@siresto.com",
		Password: hashPassword("support123"),
		RoleID:   getRoleID("SUPPORT"),
		IsActive: true,
	}
	if err := DB.FirstOrCreate(&support, entity.User{Email: "support@siresto.com"}).Error; err != nil {
		log.Printf("Error seeding SUPPORT: %v", err)
	} else {
		log.Println("✓ SUPPORT created: support@siresto.com / support123")
	}

	// 3. FINANCE
	finance := entity.User{
		Name:     "Finance Team",
		Email:    "finance@siresto.com",
		Password: hashPassword("finance123"),
		RoleID:   getRoleID("FINANCE"),
		IsActive: true,
	}
	if err := DB.FirstOrCreate(&finance, entity.User{Email: "finance@siresto.com"}).Error; err != nil {
		log.Printf("Error seeding FINANCE: %v", err)
	} else {
		log.Println("✓ FINANCE created: finance@siresto.com / finance123")
	}

	// ===== SEED EXTERNAL USERS (Client Restoran) =====

	// 4. OWNER - Buat company dulu
	owner := entity.User{
		Name:     "John Doe",
		Email:    "owner@restaurant.com",
		Password: hashPassword("owner123"),
		RoleID:   getRoleID("OWNER"),
		IsActive: true,
	}
	result := DB.FirstOrCreate(&owner, entity.User{Email: "owner@restaurant.com"})
	if result.Error != nil {
		log.Printf("Error seeding OWNER: %v", result.Error)
	} else {
		log.Printf("✓ OWNER created: owner@restaurant.com / owner123 (ID: %s)", owner.ID)
	}

	// Buat company untuk owner
	company := entity.Company{
		Name:    "PT Restoran Sejahtera",
		Type:    entity.CompanyTypePT,
		OwnerID: owner.ID,
	}
	if err := DB.FirstOrCreate(&company, entity.Company{Name: "PT Restoran Sejahtera"}).Error; err != nil {
		log.Printf("Error seeding Company: %v", err)
	} else {
		log.Printf("✓ Company created: PT Restoran Sejahtera (ID: %s)", company.ID)
	}

	// Buat branch untuk company
	branch := entity.Branch{
		CompanyID:  company.ID,
		Name:       "Cabang Jakarta Pusat",
		Address:    "Jl. Sudirman No. 123",
		City:       "Jakarta",
		Province:   "DKI Jakarta",
		PostalCode: "10220",
		Phone:      "021-12345678",
		IsActive:   true,
	}
	if err := DB.FirstOrCreate(&branch, entity.Branch{Name: "Cabang Jakarta Pusat", CompanyID: company.ID}).Error; err != nil {
		log.Printf("Error seeding Branch: %v", err)
	} else {
		log.Printf("✓ Branch created: Cabang Jakarta Pusat (ID: %s)", branch.ID)
	}

	// 5. ADMIN
	admin := entity.User{
		Name:      "Manager Cabang",
		Email:     "admin@restaurant.com",
		Password:  hashPassword("admin123"),
		RoleID:    getRoleID("ADMIN"),
		CompanyID: &company.ID,
		BranchID:  &branch.ID,
		IsActive:  true,
	}
	if err := DB.FirstOrCreate(&admin, entity.User{Email: "admin@restaurant.com"}).Error; err != nil {
		log.Printf("Error seeding ADMIN: %v", err)
	} else {
		log.Println("✓ ADMIN created: admin@restaurant.com / admin123")
	}

	// 6. CASHIER
	cashier := entity.User{
		Name:      "Kasir Jakarta",
		Email:     "cashier@restaurant.com",
		Password:  hashPassword("cashier123"),
		RoleID:    getRoleID("CASHIER"),
		CompanyID: &company.ID,
		BranchID:  &branch.ID,
		IsActive:  true,
	}
	if err := DB.FirstOrCreate(&cashier, entity.User{Email: "cashier@restaurant.com"}).Error; err != nil {
		log.Printf("Error seeding CASHIER: %v", err)
	} else {
		log.Println("✓ CASHIER created: cashier@restaurant.com / cashier123")
	}

	// 7. KITCHEN
	kitchen := entity.User{
		Name:      "Chef Dapur",
		Email:     "kitchen@restaurant.com",
		Password:  hashPassword("kitchen123"),
		RoleID:    getRoleID("KITCHEN"),
		CompanyID: &company.ID,
		BranchID:  &branch.ID,
		IsActive:  true,
	}
	if err := DB.FirstOrCreate(&kitchen, entity.User{Email: "kitchen@restaurant.com"}).Error; err != nil {
		log.Printf("Error seeding KITCHEN: %v", err)
	} else {
		log.Println("✓ KITCHEN created: kitchen@restaurant.com / kitchen123")
	}

	// 8. WAITER
	waiter := entity.User{
		Name:      "Pelayan Restoran",
		Email:     "waiter@restaurant.com",
		Password:  hashPassword("waiter123"),
		RoleID:    getRoleID("WAITER"),
		CompanyID: &company.ID,
		BranchID:  &branch.ID,
		IsActive:  true,
	}
	if err := DB.FirstOrCreate(&waiter, entity.User{Email: "waiter@restaurant.com"}).Error; err != nil {
		log.Printf("Error seeding WAITER: %v", err)
	} else {
		log.Println("✓ WAITER created: waiter@restaurant.com / waiter123")
	}

	// Seed Categories
	log.Println("\nSeeding categories...")
	
	// Main Categories (Company Level - berlaku untuk semua cabang)
	makanan := entity.Category{
		CompanyID:   company.ID,
		BranchID:    nil, // nil = berlaku untuk semua cabang
		Name:        "Makanan",
		Description: "Kategori makanan",
		Position:    1,
		IsActive:    true,
	}
	if err := DB.FirstOrCreate(&makanan, entity.Category{Name: "Makanan", CompanyID: company.ID, BranchID: nil}).Error; err != nil {
		log.Printf("Error seeding category Makanan: %v", err)
	} else {
		log.Printf("✓ Category created: Makanan (Company Level, ID: %s)", makanan.ID)
	}

	minuman := entity.Category{
		CompanyID:   company.ID,
		BranchID:    nil,
		Name:        "Minuman",
		Description: "Kategori minuman",
		Position:    2,
		IsActive:    true,
	}
	if err := DB.FirstOrCreate(&minuman, entity.Category{Name: "Minuman", CompanyID: company.ID, BranchID: nil}).Error; err != nil {
		log.Printf("Error seeding category Minuman: %v", err)
	} else {
		log.Printf("✓ Category created: Minuman (Company Level, ID: %s)", minuman.ID)
	}

	snack := entity.Category{
		CompanyID:   company.ID,
		BranchID:    nil,
		Name:        "Snack",
		Description: "Kategori snack dan cemilan",
		Position:    3,
		IsActive:    true,
	}
	if err := DB.FirstOrCreate(&snack, entity.Category{Name: "Snack", CompanyID: company.ID, BranchID: nil}).Error; err != nil {
		log.Printf("Error seeding category Snack: %v", err)
	} else {
		log.Printf("✓ Category created: Snack (Company Level, ID: %s)", snack.ID)
	}

	// Branch Specific Categories (Khusus Cabang Jakarta Pusat)
	menuSpesial := entity.Category{
		CompanyID:   company.ID,
		BranchID:    &branch.ID,
		Name:        "Menu Spesial Jakarta",
		Description: "Menu khusus cabang Jakarta Pusat",
		Position:    1,
		IsActive:    true,
	}
	if err := DB.FirstOrCreate(&menuSpesial, entity.Category{Name: "Menu Spesial Jakarta", CompanyID: company.ID, BranchID: &branch.ID}).Error; err != nil {
		log.Printf("Error seeding branch category: %v", err)
	} else {
		log.Printf("✓ Branch Category created: Menu Spesial Jakarta (Branch: %s)", branch.Name)
	}

	log.Println("Database seeding completed!")
	log.Println("\n========== TEST ACCOUNTS ==========")
	log.Println("INTERNAL USERS (Platform SIRESTO):")
	log.Println("  SUPER_ADMIN: superadmin@siresto.com / admin123")
	log.Println("  SUPPORT:     support@siresto.com / support123")
	log.Println("  FINANCE:     finance@siresto.com / finance123")
	log.Println("\nEXTERNAL USERS (Client Restoran):")
	log.Println("  OWNER:       owner@restaurant.com / owner123")
	log.Println("  ADMIN:       admin@restaurant.com / admin123")
	log.Println("  CASHIER:     cashier@restaurant.com / cashier123")
	log.Println("  KITCHEN:     kitchen@restaurant.com / kitchen123")
	log.Println("  WAITER:      waiter@restaurant.com / waiter123")
	log.Printf("\nCompany ID: %s", company.ID)
	log.Printf("Branch ID: %s\n", branch.ID)
	log.Println("===================================\n")
}
