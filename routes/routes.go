package routes

import (
	"project-name/config"
	"project-name/internal/handler"
	"project-name/internal/middleware"
	"project-name/internal/repository"
	"project-name/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Serve static files (uploaded images)
	r.Static("/uploads", "./uploads")

	// Initialize dependencies
	userRepo := repository.NewUserRepository(config.DB)
	companyRepo := repository.NewCompanyRepository(config.DB)
	userService := service.NewUserServiceWithCompany(userRepo, companyRepo)
	userHandler := handler.NewUserHandler(userService)

	branchRepo := repository.NewBranchRepository(config.DB)
	companyService := service.NewCompanyService(companyRepo, userRepo)
	companyHandler := handler.NewCompanyHandler(companyService)

	branchService := service.NewBranchService(branchRepo, companyRepo, userRepo)
	branchHandler := handler.NewBranchHandler(branchService)

	roleRepo := repository.NewRoleRepository(config.DB)
	roleService := service.NewRoleService(roleRepo)
	roleHandler := handler.NewRoleHandler(roleService)

	categoryRepo := repository.NewCategoryRepository(config.DB)
	categoryService := service.NewCategoryService(categoryRepo, companyRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	productRepo := repository.NewProductRepository(config.DB)
	productService := service.NewProductService(productRepo, categoryRepo, branchRepo)
	productHandler := handler.NewProductHandler(productService)

	// API v1
	v1 := r.Group("/api/v1")

	// Public routes
	public := v1.Group("")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
	}

	// Auth routes (protected)
	auth := v1.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/me", userHandler.GetMe)
	}

	// Role routes (protected)
	roles := v1.Group("/roles")
	roles.Use(middleware.AuthMiddleware())
	{
		roles.GET("", roleHandler.GetAllRoles)
	}

	// ===== EXTERNAL API (untuk client restoran) =====
	external := v1.Group("/external")
	external.Use(middleware.AuthMiddleware())
	external.Use(middleware.RequireExternalRole()) // Hanya external users
	{
		// User management
		external.POST("/users", userHandler.CreateExternalUser)
		external.GET("/users", userHandler.GetExternalUsers)
		external.GET("/users/company/:company_id", userHandler.GetCompanyUsers)
		external.GET("/users/branch/:branch_id", userHandler.GetBranchUsers)

		// Company routes
		external.POST("/companies", companyHandler.CreateCompany)
		external.GET("/companies/detail/:id", companyHandler.GetCompany)
		external.GET("/companies/my", companyHandler.GetMyCompanies)

		// Branch routes
		external.POST("/branches", branchHandler.CreateBranch)
		external.GET("/branches/detail/:id", branchHandler.GetBranch)
		external.GET("/branches/company/:company_id", branchHandler.GetBranchesByCompany)

		// Category routes
		external.POST("/categories", categoryHandler.CreateCategory)
		external.PUT("/categories/:id", categoryHandler.UpdateCategory)
		external.DELETE("/categories/:id", categoryHandler.DeleteCategory)
		external.GET("/categories/:id", categoryHandler.GetCategory)
		external.GET("/categories", categoryHandler.GetCategories)

		// Product routes
		external.POST("/products", productHandler.CreateProduct)
		external.PUT("/products/:id", productHandler.UpdateProduct)
		external.DELETE("/products/:id", productHandler.DeleteProduct)
		external.GET("/products/:id", productHandler.GetProductByID)
		external.GET("/products", productHandler.GetAllProducts)
	}

	// ===== DASHBOARD API (untuk platform SIRESTO) =====
	dashboard := v1.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware())
	dashboard.Use(middleware.RequireInternalRole()) // Hanya internal users
	{
		// User management
		dashboard.POST("/users", userHandler.CreateInternalUser)
		dashboard.GET("/users", userHandler.GetInternalUsers)
		dashboard.GET("/users/:id", userHandler.GetUser)
		
		// Lihat semua companies (untuk monitoring)
		dashboard.GET("/companies", companyHandler.GetMyCompanies) // TODO: buat endpoint khusus untuk list all
		
		// Lihat semua external users (client restoran)
		dashboard.GET("/external-users", userHandler.GetExternalUsers)
	}
}
