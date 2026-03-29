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

type PromoHandler struct {
	promoService service.PromoService
}

func NewPromoHandler(promoService service.PromoService) *PromoHandler {
	return &PromoHandler{promoService: promoService}
}

// CreatePromo godoc
// @Summary Create new promo
// @Description Create a new promo
// @Tags Promo
// @Accept json
// @Produce json
// @Param promo body entity.CreatePromoRequest true "Promo data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/external/promos [post]
func (h *PromoHandler) CreatePromo(c *gin.Context) {
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

	var req entity.CreatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	promo, err := h.promoService.CreatePromo(companyID.(uuid.UUID), branchID, &req)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to create promo", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusCreated, "Promo created successfully", promo)
}

// UpdatePromo godoc
// @Summary Update promo
// @Description Update an existing promo
// @Tags Promo
// @Accept json
// @Produce json
// @Param id path string true "Promo ID"
// @Param promo body entity.UpdatePromoRequest true "Promo data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/external/promos/{id} [put]
func (h *PromoHandler) UpdatePromo(c *gin.Context) {
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
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request", "Invalid promo ID")
		return
	}

	// Get promo first to check if it's company-level or branch-level
	existingPromo, err := h.promoService.GetPromoByID(id, companyID.(uuid.UUID), branchID)
	if err != nil {
		if err.Error() == "promo not found" {
			pkg.ErrorResponse(c, http.StatusNotFound, "Promo not found", err.Error())
			return
		}
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve promo", err.Error())
		return
	}

	// Check permission: Only OWNER and ADMIN can update company-level promo
	if existingPromo.BranchID == nil && userRole != "OWNER" && userRole != "ADMIN" {
		pkg.ErrorResponse(c, http.StatusForbidden, "Forbidden", "Only OWNER and ADMIN can update company-level promo")
		return
	}

	var req entity.UpdatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	promo, err := h.promoService.UpdatePromo(id, companyID.(uuid.UUID), branchID, &req)
	if err != nil {
		if err.Error() == "promo not found" {
			pkg.ErrorResponse(c, http.StatusNotFound, "Promo not found", err.Error())
			return
		}
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to update promo", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Promo updated successfully", promo)
}

// DeletePromo godoc
// @Summary Delete promo
// @Description Delete a promo
// @Tags Promo
// @Produce json
// @Param id path string true "Promo ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/external/promos/{id} [delete]
func (h *PromoHandler) DeletePromo(c *gin.Context) {
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
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request", "Invalid promo ID")
		return
	}

	// Get promo first to check if it's company-level or branch-level
	promo, err := h.promoService.GetPromoByID(id, companyID.(uuid.UUID), branchID)
	if err != nil {
		if err.Error() == "promo not found" {
			pkg.ErrorResponse(c, http.StatusNotFound, "Promo not found", err.Error())
			return
		}
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve promo", err.Error())
		return
	}

	// Check permission: Only OWNER and ADMIN can delete company-level promo
	if promo.BranchID == nil && userRole != "OWNER" && userRole != "ADMIN" {
		pkg.ErrorResponse(c, http.StatusForbidden, "Forbidden", "Only OWNER and ADMIN can delete company-level promo")
		return
	}

	if err := h.promoService.DeletePromo(id, companyID.(uuid.UUID), branchID); err != nil {
		if err.Error() == "promo not found" {
			pkg.ErrorResponse(c, http.StatusNotFound, "Promo not found", err.Error())
			return
		}
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete promo", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Promo deleted successfully", nil)
}

// GetPromoByID godoc
// @Summary Get promo by ID
// @Description Get a single promo by ID
// @Tags Promo
// @Produce json
// @Param id path string true "Promo ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/external/promos/{id} [get]
func (h *PromoHandler) GetPromoByID(c *gin.Context) {
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
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request", "Invalid promo ID")
		return
	}

	promo, err := h.promoService.GetPromoByID(id, companyID.(uuid.UUID), branchID)
	if err != nil {
		if err.Error() == "promo not found" {
			pkg.ErrorResponse(c, http.StatusNotFound, "Promo not found", err.Error())
			return
		}
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve promo", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Promo retrieved successfully", promo)
}

// GetAllPromos godoc
// @Summary Get all promos
// @Description Get all promos with pagination
// @Tags Promo
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/external/promos [get]
func (h *PromoHandler) GetAllPromos(c *gin.Context) {
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

	promos, meta, err := h.promoService.GetAllPromos(companyID.(uuid.UUID), branchID, page, limit)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve promos", err.Error())
		return
	}

	pkg.SuccessResponseWithMeta(c, http.StatusOK, "Promos retrieved successfully", promos, meta)
}


// ValidatePromoCode godoc
// @Summary Validate promo code
// @Description Check if a promo code is valid and can be used
// @Tags Promo
// @Produce json
// @Param code path string true "Promo Code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/external/promos/validate/{code} [get]
func (h *PromoHandler) ValidatePromoCode(c *gin.Context) {
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

	code := c.Param("code")
	if code == "" {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request", "Promo code is required")
		return
	}

	validation, err := h.promoService.ValidatePromoCode(code, companyID.(uuid.UUID), branchID)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to validate promo", err.Error())
		return
	}

	if validation.Valid {
		pkg.SuccessResponse(c, http.StatusOK, validation.Message, validation)
	} else {
		pkg.SuccessResponse(c, http.StatusOK, validation.Message, validation)
	}
}
