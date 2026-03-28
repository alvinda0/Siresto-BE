package handler

import (
	"project-name/internal/entity"
	"project-name/internal/service"
	"project-name/pkg"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService}
}

type ProductResponse struct {
	ID          string  `json:"id"`
	CompanyID   string  `json:"company_id"`
	BranchID    string  `json:"branch_id"`
	CategoryID  string  `json:"category_id"`
	Image       string  `json:"image"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
	Position    string  `json:"position"`
	IsAvailable bool    `json:"is_available"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Company     struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"company"`
	Branch struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"branch"`
	Category struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
}

func toProductResponse(product *entity.Product) ProductResponse {
	resp := ProductResponse{
		ID:          product.ID.String(),
		CompanyID:   product.CompanyID.String(),
		BranchID:    product.BranchID.String(),
		CategoryID:  product.CategoryID.String(),
		Image:       product.Image,
		Name:        product.Name,
		Description: product.Description,
		Stock:       product.Stock,
		Price:       product.Price,
		Position:    product.Position,
		IsAvailable: product.IsAvailable,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	
	resp.Company.ID = product.Company.ID.String()
	resp.Company.Name = product.Company.Name
	
	resp.Branch.ID = product.Branch.ID.String()
	resp.Branch.Name = product.Branch.Name
	
	resp.Category.ID = product.Category.ID.String()
	resp.Category.Name = product.Category.Name
	
	return resp
}

func toProductResponseList(products []entity.Product) []ProductResponse {
	responses := make([]ProductResponse, len(products))
	for i, product := range products {
		responses[i] = toProductResponse(&product)
	}
	return responses
}

type CreateProductRequest struct {
	BranchID    string  `json:"branch_id" binding:"required"`
	CategoryID  string  `json:"category_id" binding:"required"`
	Image       string  `json:"image"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price" binding:"required"`
	Position    string  `json:"position"`
	IsAvailable bool    `json:"is_available"`
}

type UpdateProductRequest struct {
	CategoryID  string  `json:"category_id" binding:"required"`
	Image       string  `json:"image"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price" binding:"required"`
	Position    string  `json:"position"`
	IsAvailable bool    `json:"is_available"`
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest

	// Check content type
	contentType := c.GetHeader("Content-Type")
	
	if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form data (with file upload)
		req.BranchID = strings.TrimSpace(c.PostForm("branch_id"))
		req.CategoryID = strings.TrimSpace(c.PostForm("category_id"))
		req.Name = strings.TrimSpace(c.PostForm("name"))
		req.Description = strings.TrimSpace(c.PostForm("description"))
		req.Position = strings.TrimSpace(c.PostForm("position"))
		
		// Parse numeric fields
		stockStr := strings.TrimSpace(c.PostForm("stock"))
		if stockStr != "" {
			stock, err := strconv.Atoi(stockStr)
			if err == nil {
				req.Stock = stock
			}
		}
		
		priceStr := strings.TrimSpace(c.PostForm("price"))
		if priceStr != "" {
			price, err := strconv.ParseFloat(priceStr, 64)
			if err == nil {
				req.Price = price
			}
		}
		
		availableStr := strings.TrimSpace(c.PostForm("is_available"))
		if availableStr != "" {
			req.IsAvailable = availableStr == "true" || availableStr == "1"
		} else {
			req.IsAvailable = true // default true
		}

		// Handle file upload
		if file, err := c.FormFile("image"); err == nil {
			config := pkg.DefaultImageUploadConfig()
			filePath, err := config.SaveFile(file)
			if err != nil {
				pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to upload image", err.Error())
				return
			}
			
			// Get base URL
			baseURL := os.Getenv("BASE_URL")
			if baseURL == "" {
				baseURL = "http://localhost:8080"
			}
			req.Image = pkg.GetFileURL(filePath, baseURL)
		} else if imageURL := strings.TrimSpace(c.PostForm("image")); imageURL != "" {
			// Use URL if provided
			req.Image = imageURL
		}

		// Validate required fields with detailed error messages
		if req.BranchID == "" {
			pkg.ErrorResponse(c, http.StatusBadRequest, "branch_id is required", "")
			return
		}
		if req.CategoryID == "" {
			pkg.ErrorResponse(c, http.StatusBadRequest, "category_id is required", "")
			return
		}
		if req.Name == "" {
			pkg.ErrorResponse(c, http.StatusBadRequest, "name is required", "")
			return
		}
		if req.Price <= 0 {
			pkg.ErrorResponse(c, http.StatusBadRequest, "price is required and must be greater than 0", "")
			return
		}
	} else {
		// Handle JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}
	}

	// Get company_id and branch_id from context (set by auth middleware)
	companyIDVal, exists := c.Get("companyID")
	if !exists || companyIDVal == nil {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Company ID not found", "")
		return
	}

	branchIDVal, exists := c.Get("branchID")
	if !exists || branchIDVal == nil {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Branch ID not found", "")
		return
	}

	// Extract UUID from pointer
	companyID := *(companyIDVal.(*uuid.UUID))
	branchID := *(branchIDVal.(*uuid.UUID))

	// Parse UUIDs
	reqBranchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid branch ID format", err.Error())
		return
	}

	reqCategoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID format", err.Error())
		return
	}

	// Validate that the requested branch_id matches the user's branch_id
	if reqBranchID != branchID {
		pkg.ErrorResponse(c, http.StatusForbidden, "You can only create products for your own branch", "")
		return
	}

	product := &entity.Product{
		CompanyID:   companyID,
		BranchID:    reqBranchID,
		CategoryID:  reqCategoryID,
		Image:       req.Image,
		Name:        req.Name,
		Description: req.Description,
		Stock:       req.Stock,
		Price:       req.Price,
		Position:    req.Position,
		IsAvailable: req.IsAvailable,
	}

	if err := h.productService.CreateProduct(product); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to create product", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusCreated, "Product created successfully", toProductResponse(product))
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	// Get company_id and branch_id from context
	companyIDVal, exists := c.Get("companyID")
	if !exists || companyIDVal == nil {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Company ID not found", "")
		return
	}

	// Get user role
	externalRole, _ := c.Get("externalRole")
	userBranchIDVal, _ := c.Get("branchID")

	// Extract company UUID from pointer
	companyID := *(companyIDVal.(*uuid.UUID))

	var branchID uuid.UUID

	// Check if user is OWNER
	if externalRole == "OWNER" {
		// OWNER must provide branch_id in query params
		branchIDParam := c.Query("branch_id")
		if branchIDParam == "" {
			pkg.ErrorResponse(c, http.StatusBadRequest, "OWNER must provide branch_id parameter", "")
			return
		}

		// Parse branch_id from query
		parsedBranchID, err := uuid.Parse(branchIDParam)
		if err != nil {
			pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid branch_id format", err.Error())
			return
		}

		// Validate that branch belongs to owner's company
		// This will be checked in service layer
		branchID = parsedBranchID
	} else {
		// Non-OWNER users use their branch_id from token
		if userBranchIDVal == nil {
			pkg.ErrorResponse(c, http.StatusUnauthorized, "Branch ID not found", "")
			return
		}
		branchID = *(userBranchIDVal.(*uuid.UUID))
	}

	// Get query parameters
	search := c.DefaultQuery("search", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := h.productService.GetAllProducts(
		companyID,
		branchID,
		search,
		page,
		limit,
	)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch products", err.Error())
		return
	}

	// Create pagination meta
	paginationParams := pkg.PaginationParams{Page: page, Limit: limit}
	meta := paginationParams.CreateMeta(total)

	pkg.SuccessResponseWithMeta(c, http.StatusOK, "Products retrieved successfully", toProductResponseList(products), meta)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	// Get company_id and branch_id from context
	companyIDVal, exists := c.Get("companyID")
	if !exists || companyIDVal == nil {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Company ID not found", "")
		return
	}

	branchIDVal, exists := c.Get("branchID")
	if !exists || branchIDVal == nil {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Branch ID not found", "")
		return
	}

	// Extract UUID from pointer
	companyID := *(companyIDVal.(*uuid.UUID))
	branchID := *(branchIDVal.(*uuid.UUID))

	product, err := h.productService.GetProductByID(id, companyID, branchID)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusNotFound, "Product not found", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", toProductResponse(product))
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	var req UpdateProductRequest

	// Check content type
	contentType := c.GetHeader("Content-Type")
	
	if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form data (with file upload)
		req.CategoryID = c.PostForm("category_id")
		req.Name = c.PostForm("name")
		req.Description = c.PostForm("description")
		req.Position = c.PostForm("position")
		
		// Parse numeric fields
		if stockStr := c.PostForm("stock"); stockStr != "" {
			stock, _ := strconv.Atoi(stockStr)
			req.Stock = stock
		}
		if priceStr := c.PostForm("price"); priceStr != "" {
			price, _ := strconv.ParseFloat(priceStr, 64)
			req.Price = price
		}
		if availableStr := c.PostForm("is_available"); availableStr != "" {
			req.IsAvailable = availableStr == "true" || availableStr == "1"
		}

		// Handle file upload
		if file, err := c.FormFile("image"); err == nil {
			config := pkg.DefaultImageUploadConfig()
			filePath, err := config.SaveFile(file)
			if err != nil {
				pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to upload image", err.Error())
				return
			}
			
			// Get base URL
			baseURL := os.Getenv("BASE_URL")
			if baseURL == "" {
				baseURL = "http://localhost:8080"
			}
			req.Image = pkg.GetFileURL(filePath, baseURL)
		} else if imageURL := c.PostForm("image"); imageURL != "" {
			// Use URL if provided
			req.Image = imageURL
		}

		// Validate required fields
		if req.CategoryID == "" || req.Name == "" || req.Price == 0 {
			pkg.ErrorResponse(c, http.StatusBadRequest, "Missing required fields: category_id, name, price", "")
			return
		}
	} else {
		// Handle JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}
	}

	// Get company_id and branch_id from context
	companyIDVal, exists := c.Get("companyID")
	if !exists || companyIDVal == nil {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Company ID not found", "")
		return
	}

	branchIDVal, exists := c.Get("branchID")
	if !exists || branchIDVal == nil {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Branch ID not found", "")
		return
	}

	// Extract UUID from pointer
	companyID := *(companyIDVal.(*uuid.UUID))
	branchID := *(branchIDVal.(*uuid.UUID))

	// Parse category UUID
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID format", err.Error())
		return
	}

	product := &entity.Product{
		ID:          id,
		CompanyID:   companyID,
		BranchID:    branchID,
		CategoryID:  categoryID,
		Image:       req.Image,
		Name:        req.Name,
		Description: req.Description,
		Stock:       req.Stock,
		Price:       req.Price,
		Position:    req.Position,
		IsAvailable: req.IsAvailable,
	}

	if err := h.productService.UpdateProduct(product); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to update product", err.Error())
		return
	}

	// Fetch updated product with relations
	updatedProduct, _ := h.productService.GetProductByID(id, companyID, branchID)
	pkg.SuccessResponse(c, http.StatusOK, "Product updated successfully", toProductResponse(updatedProduct))
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	// Get company_id and branch_id from context
	companyIDVal, exists := c.Get("companyID")
	if !exists || companyIDVal == nil {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Company ID not found", "")
		return
	}

	branchIDVal, exists := c.Get("branchID")
	if !exists || branchIDVal == nil {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Branch ID not found", "")
		return
	}

	// Extract UUID from pointer
	companyID := *(companyIDVal.(*uuid.UUID))
	branchID := *(branchIDVal.(*uuid.UUID))

	if err := h.productService.DeleteProduct(id, companyID, branchID); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to delete product", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}
