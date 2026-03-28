package handler

import (
	"net/http"
	"project-name/internal/service"
	"project-name/pkg"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BranchHandler struct {
	branchService *service.BranchService
}

func NewBranchHandler(branchService *service.BranchService) *BranchHandler {
	return &BranchHandler{branchService: branchService}
}

type CreateBranchRequest struct {
	CompanyID  uuid.UUID `json:"company_id" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	Address    string    `json:"address" binding:"required"`
	City       string    `json:"city"`
	Province   string    `json:"province"`
	PostalCode string    `json:"postal_code"`
	Phone      string    `json:"phone"`
}

func (h *BranchHandler) CreateBranch(c *gin.Context) {
	var req CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	branch, err := h.branchService.CreateBranch(
		req.CompanyID,
		req.Name,
		req.Address,
		req.City,
		req.Province,
		req.PostalCode,
		req.Phone,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": branch})
}

func (h *BranchHandler) GetBranch(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid branch ID"})
		return
	}

	branch, err := h.branchService.GetBranchByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "branch not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": branch})
}

func (h *BranchHandler) GetBranchesByCompany(c *gin.Context) {
	companyID, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
		return
	}

	// Get current user ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	pagination := pkg.GetPaginationParams(c)
	branches, total, err := h.branchService.GetBranchesByCompanyFiltered(companyID, userID.(uuid.UUID), pagination.Limit, pagination.CalculateOffset())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := pagination.CreateMeta(total)
	pkg.SuccessResponseWithMeta(c, http.StatusOK, "Branches retrieved successfully", branches, meta)
}
