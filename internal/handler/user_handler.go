package handler

import (
	"project-name/internal/service"
	"project-name/pkg"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register untuk external user (owner restoran) dengan company
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=6"`
		CompanyName string `json:"company_name" binding:"required"`
		CompanyType string `json:"company_type" binding:"required"` // PT atau PERORANGAN
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, 400, "Invalid request", err.Error())
		return
	}
	
	// Validasi company type
	if req.CompanyType != "PT" && req.CompanyType != "CV" && req.CompanyType != "PERORANGAN" {
		pkg.ErrorResponse(c, 400, "Invalid company type", "Company type must be PT, CV, or PERORANGAN")
		return
	}
	
	user, company, err := h.service.RegisterWithCompany(req.Name, req.Email, req.Password, req.CompanyName, req.CompanyType)
	if err != nil {
		pkg.ErrorResponse(c, 500, "Registration failed", err.Error())
		return
	}
	
	pkg.SuccessResponse(c, 201, "User and company registered successfully", gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
		"company": gin.H{
			"id":   company.ID,
			"name": company.Name,
			"type": company.Type,
		},
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, 400, "Invalid request", err.Error())
		return
	}
	user, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		pkg.ErrorResponse(c, 401, "Login failed", err.Error())
		return
	}
	
	// Tentukan internal atau external role berdasarkan role type
	var internalRole, externalRole string
	roleType := string(user.Role.Type) // INTERNAL or EXTERNAL
	
	if user.Role.Type == "INTERNAL" {
		internalRole = user.Role.Name
	} else if user.Role.Type == "EXTERNAL" {
		externalRole = user.Role.Name
	}
	
	// Generate JWT token
	token, err := pkg.GenerateJWT(user.ID, user.Email, roleType, internalRole, externalRole, user.CompanyID, user.BranchID)
	if err != nil {
		pkg.ErrorResponse(c, 500, "Failed to generate token", err.Error())
		return
	}
	
	pkg.SuccessResponse(c, 200, "Login successful", gin.H{
		"token": token,
	})
}

// GetMe untuk mendapatkan informasi user yang sedang login
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		pkg.ErrorResponse(c, 401, "Unauthorized", "User ID not found in token")
		return
	}

	user, err := h.service.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		pkg.ErrorResponse(c, 404, "User not found", err.Error())
		return
	}

	pkg.SuccessResponse(c, 200, "User profile retrieved successfully", user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.ErrorResponse(c, 400, "Invalid user ID", err.Error())
		return
	}
	user, err := h.service.GetUserByID(id)
	if err != nil {
		pkg.ErrorResponse(c, 404, "User not found", err.Error())
		return
	}
	user.Password = ""
	pkg.SuccessResponse(c, 200, "User retrieved successfully", user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		pkg.ErrorResponse(c, 500, "Failed to retrieve users", err.Error())
		return
	}
	pkg.SuccessResponse(c, 200, "Users retrieved successfully", users)
}

// ===== EXTERNAL USER ENDPOINTS (untuk client restoran) =====

type CreateExternalUserRequest struct {
	Name      string     `json:"name" binding:"required"`
	Email     string     `json:"email" binding:"required,email"`
	Password  string     `json:"password" binding:"required,min=6"`
	RoleID    uuid.UUID  `json:"role_id" binding:"required"`
	CompanyID *uuid.UUID `json:"company_id"`
	BranchID  *uuid.UUID `json:"branch_id"`
}

// CreateExternalUser untuk owner membuat admin, cashier, kitchen, waiter
func (h *UserHandler) CreateExternalUser(c *gin.Context) {
	var req CreateExternalUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, 400, "Invalid request", err.Error())
		return
	}

	user, err := h.service.CreateExternalUser(
		req.Name,
		req.Email,
		req.Password,
		req.RoleID,
		req.CompanyID,
		req.BranchID,
	)
	if err != nil {
		pkg.ErrorResponse(c, 400, "Failed to create user", err.Error())
		return
	}

	user.Password = ""
	pkg.SuccessResponse(c, 201, "User created successfully", user)
}

// GetCompanyUsers untuk mendapatkan semua user dalam perusahaan
func (h *UserHandler) GetCompanyUsers(c *gin.Context) {
	companyID, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		pkg.ErrorResponse(c, 400, "Invalid company ID", err.Error())
		return
	}

	// Get current user ID
	userID, exists := c.Get("user_id")
	if !exists {
		pkg.ErrorResponse(c, 401, "Unauthorized", "User ID not found")
		return
	}

	pagination := pkg.GetPaginationParams(c)
	users, total, err := h.service.GetUsersByCompanyFiltered(companyID, userID.(uuid.UUID), pagination.Limit, pagination.CalculateOffset())
	if err != nil {
		pkg.ErrorResponse(c, 500, "Failed to retrieve users", err.Error())
		return
	}

	meta := pagination.CreateMeta(total)
	pkg.SuccessResponseWithMeta(c, 200, "Users retrieved successfully", users, meta)
}

// GetBranchUsers untuk mendapatkan semua user dalam cabang
func (h *UserHandler) GetBranchUsers(c *gin.Context) {
	branchID, err := uuid.Parse(c.Param("branch_id"))
	if err != nil {
		pkg.ErrorResponse(c, 400, "Invalid branch ID", err.Error())
		return
	}

	pagination := pkg.GetPaginationParams(c)
	users, total, err := h.service.GetUsersByBranch(branchID, pagination.Limit, pagination.CalculateOffset())
	if err != nil {
		pkg.ErrorResponse(c, 500, "Failed to retrieve users", err.Error())
		return
	}

	meta := pagination.CreateMeta(total)
	pkg.SuccessResponseWithMeta(c, 200, "Users retrieved successfully", users, meta)
}

// GetExternalUsers untuk mendapatkan semua external user (client restoran)
func (h *UserHandler) GetExternalUsers(c *gin.Context) {
	pagination := pkg.GetPaginationParams(c)
	users, total, err := h.service.GetExternalUsers(pagination.Limit, pagination.CalculateOffset())
	if err != nil {
		pkg.ErrorResponse(c, 500, "Failed to retrieve users", err.Error())
		return
	}

	meta := pagination.CreateMeta(total)
	pkg.SuccessResponseWithMeta(c, 200, "External users retrieved successfully", users, meta)
}

// ===== INTERNAL USER ENDPOINTS (untuk platform SIRESTO) =====

type CreateInternalUserRequest struct {
	Name     string    `json:"name" binding:"required"`
	Email    string    `json:"email" binding:"required,email"`
	Password string    `json:"password" binding:"required,min=6"`
	RoleID   uuid.UUID `json:"role_id" binding:"required"`
}

// CreateInternalUser untuk super_admin membuat support, finance
func (h *UserHandler) CreateInternalUser(c *gin.Context) {
	var req CreateInternalUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, 400, "Invalid request", err.Error())
		return
	}

	user, err := h.service.CreateInternalUser(
		req.Name,
		req.Email,
		req.Password,
		req.RoleID,
	)
	if err != nil {
		pkg.ErrorResponse(c, 400, "Failed to create user", err.Error())
		return
	}

	user.Password = ""
	pkg.SuccessResponse(c, 201, "Internal user created successfully", user)
}

// GetInternalUsers untuk mendapatkan semua internal user platform
func (h *UserHandler) GetInternalUsers(c *gin.Context) {
	pagination := pkg.GetPaginationParams(c)
	users, total, err := h.service.GetInternalUsers(pagination.Limit, pagination.CalculateOffset())
	if err != nil {
		pkg.ErrorResponse(c, 500, "Failed to retrieve users", err.Error())
		return
	}

	meta := pagination.CreateMeta(total)
	pkg.SuccessResponseWithMeta(c, 200, "Internal users retrieved successfully", users, meta)
}
