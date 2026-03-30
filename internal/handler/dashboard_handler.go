package handler

import (
	"net/http"
	"project-name/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetHomeStats menampilkan statistik untuk halaman home/dashboard
func (h *DashboardHandler) GetHomeStats(c *gin.Context) {
	// Ambil company_id dan branch_id dari context (dari middleware auth)
	companyIDInterface, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found"})
		return
	}
	companyID := companyIDInterface.(uuid.UUID)
	
	// Branch ID optional, jika tidak ada maka tampilkan semua branch
	var branchID uuid.UUID
	branchIDInterface, exists := c.Get("branch_id")
	if exists {
		branchID = branchIDInterface.(uuid.UUID)
	}
	
	// Get stats
	stats, err := h.dashboardService.GetHomeStats(companyID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}
