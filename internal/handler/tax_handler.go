package handler

import (
	"net/http"
	"project-name/internal/entity"
	"project-name/internal/service"
	"project-name/pkg"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaxHandler struct {
	taxService service.TaxService
}

func NewTaxHandler(taxService service.TaxService) *TaxHandler {
	return &TaxHandler{taxService: taxService}
}

// CreateTax godoc
// @Summary Create new tax
// @Description Create a new tax record
// @Tags Tax
// @Accept json
// @Produce json
// @Param tax body entity.CreateTaxRequest true "Tax data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/external/tax [post]
func (h *TaxHandler) CreateTax(c *gin.Context) {
	// Get company_id and branch_id from context (set by auth middleware)
	companyID, exists := c.Get("company_id")
	if !exists {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Company ID not found in context")
		return
	}

	branchIDVal, _ := c.Get("branch_id")
	var branchID *uuid.UUID
	if branchIDVal != nil {
		if bid, ok := branchIDVal.(uuid.UUID); ok {
			branchID = &bid
		}
	}

	var req entity.CreateTaxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	tax, err := h.taxService.CreateTax(companyID.(uuid.UUID), branchID, &req)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to create tax", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusCreated, "Tax created successfully", tax)
}

// UpdateTax godoc
// @Summary Update tax
// @Description Update an existing tax record
// @Tags Tax
// @Accept json
// @Produce json
// @Param id path string true "Tax ID"
// @Param tax body entity.UpdateTaxRequest true "Tax data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/external/tax/{id} [put]
func (h *TaxHandler) UpdateTax(c *gin.Context) {
	companyID, exists := c.Get("company_id")
	if !exists {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Company ID not found in context")
		return
	}

	branchIDVal, _ := c.Get("branch_id")
	var branchID *uuid.UUID
	if branchIDVal != nil {
		if bid, ok := branchIDVal.(uuid.UUID); ok {
			branchID = &bid
		}
	}

	// Get user role
	externalRole, _ := c.Get("external_role")
	userRole := ""
	if externalRole != nil {
		userRole = externalRole.(string)
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request", "Invalid tax ID")
		return
	}

	// Get tax first to check if it's company-level or branch-level
	existingTax, err := h.taxService.GetTaxByID(id, companyID.(uuid.UUID), branchID)
	if err != nil {
		if err.Error() == "tax not found" {
			pkg.ErrorResponse(c, http.StatusNotFound, "Tax not found", err.Error())
			return
		}
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve tax", err.Error())
		return
	}

	// Check permission: Only OWNER can update company-level tax
	if existingTax.BranchID == nil && userRole != "OWNER" {
		pkg.ErrorResponse(c, http.StatusForbidden, "Forbidden", "Only OWNER can update company-level tax")
		return
	}

	var req entity.UpdateTaxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	tax, err := h.taxService.UpdateTax(id, companyID.(uuid.UUID), branchID, &req)
	if err != nil {
		if err.Error() == "tax not found" {
			pkg.ErrorResponse(c, http.StatusNotFound, "Tax not found", err.Error())
			return
		}
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to update tax", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Tax updated successfully", tax)
}

// DeleteTax godoc
// @Summary Delete tax
// @Description Delete a tax record
// @Tags Tax
// @Produce json
// @Param id path string true "Tax ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/external/tax/{id} [delete]
func (h *TaxHandler) DeleteTax(c *gin.Context) {
	companyID, exists := c.Get("company_id")
	if !exists {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Company ID not found in context")
		return
	}

	branchIDVal, _ := c.Get("branch_id")
	var branchID *uuid.UUID
	if branchIDVal != nil {
		if bid, ok := branchIDVal.(uuid.UUID); ok {
			branchID = &bid
		}
	}

	// Get user role
	externalRole, _ := c.Get("external_role")
	userRole := ""
	if externalRole != nil {
		userRole = externalRole.(string)
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request", "Invalid tax ID")
		return
	}

	// Get tax first to check if it's company-level or branch-level
	tax, err := h.taxService.GetTaxByID(id, companyID.(uuid.UUID), branchID)
	if err != nil {
		if err.Error() == "tax not found" {
			pkg.ErrorResponse(c, http.StatusNotFound, "Tax not found", err.Error())
			return
		}
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve tax", err.Error())
		return
	}

	// Check permission: Only OWNER can delete company-level tax
	if tax.BranchID == nil && userRole != "OWNER" {
		pkg.ErrorResponse(c, http.StatusForbidden, "Forbidden", "Only OWNER can delete company-level tax")
		return
	}

	if err := h.taxService.DeleteTax(id, companyID.(uuid.UUID), branchID); err != nil {
		if err.Error() == "tax not found" {
			pkg.ErrorResponse(c, http.StatusNotFound, "Tax not found", err.Error())
			return
		}
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete tax", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Tax deleted successfully", nil)
}

// GetTaxByID godoc
// @Summary Get tax by ID
// @Description Get a single tax record by ID
// @Tags Tax
// @Produce json
// @Param id path string true "Tax ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/external/tax/{id} [get]
func (h *TaxHandler) GetTaxByID(c *gin.Context) {
	companyID, exists := c.Get("company_id")
	if !exists {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Company ID not found in context")
		return
	}

	branchIDVal, _ := c.Get("branch_id")
	var branchID *uuid.UUID
	if branchIDVal != nil {
		if bid, ok := branchIDVal.(uuid.UUID); ok {
			branchID = &bid
		}
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request", "Invalid tax ID")
		return
	}

	tax, err := h.taxService.GetTaxByID(id, companyID.(uuid.UUID), branchID)
	if err != nil {
		if err.Error() == "tax not found" {
			pkg.ErrorResponse(c, http.StatusNotFound, "Tax not found", err.Error())
			return
		}
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve tax", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Tax retrieved successfully", tax)
}

// GetAllTaxes godoc
// @Summary Get all taxes
// @Description Get all tax records for the user's company/branch with pagination
// @Tags Tax
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/external/tax [get]
func (h *TaxHandler) GetAllTaxes(c *gin.Context) {
	companyID, exists := c.Get("company_id")
	if !exists {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Company ID not found in context")
		return
	}

	branchIDVal, _ := c.Get("branch_id")
	var branchID *uuid.UUID
	if branchIDVal != nil {
		if bid, ok := branchIDVal.(uuid.UUID); ok {
			branchID = &bid
		}
	}

	// Parse pagination params
	page := 1
	limit := 10
	
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	taxes, meta, err := h.taxService.GetAllTaxes(companyID.(uuid.UUID), branchID, page, limit)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve taxes", err.Error())
		return
	}

	pkg.SuccessResponseWithMeta(c, http.StatusOK, "Taxes retrieved successfully", taxes, meta)
}
