package handler

import (
	"net/http"
	"project-name/internal/entity"
	"project-name/internal/service"
	"project-name/pkg"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

type CreateCategoryRequest struct {
	CompanyID   uuid.UUID  `json:"company_id" binding:"required"`
	BranchID    *uuid.UUID `json:"branch_id"`
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	Position    int        `json:"position"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Position    int    `json:"position"`
	IsActive    bool   `json:"is_active"`
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	category := &entity.Category{
		CompanyID:   req.CompanyID,
		BranchID:    req.BranchID,
		Name:        req.Name,
		Description: req.Description,
		Position:    req.Position,
		IsActive:    true,
	}

	if err := h.categoryService.CreateCategory(category); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to create category", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusCreated, "Category created successfully", category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", err.Error())
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	category := &entity.Category{
		Name:        req.Name,
		Description: req.Description,
		Position:    req.Position,
		IsActive:    req.IsActive,
	}

	if err := h.categoryService.UpdateCategory(id, category); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to update category", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Category updated successfully", nil)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", err.Error())
		return
	}

	if err := h.categoryService.DeleteCategory(id); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to delete category", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", err.Error())
		return
	}

	category, err := h.categoryService.GetCategoryByID(id)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusNotFound, "Category not found", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	// Query parameters
	companyIDStr := c.Query("company_id")
	branchIDStr := c.Query("branch_id")

	// Validasi company_id required
	if companyIDStr == "" {
		pkg.ErrorResponse(c, http.StatusBadRequest, "company_id is required", "company_id query parameter is missing")
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid company ID", err.Error())
		return
	}

	// Optional branch_id
	var branchID *uuid.UUID
	if branchIDStr != "" {
		bid, err := uuid.Parse(branchIDStr)
		if err != nil {
			pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid branch ID", err.Error())
			return
		}
		branchID = &bid
	}

	categories, err := h.categoryService.GetCategoriesByCompany(companyID, branchID)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve categories", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", categories)
}
