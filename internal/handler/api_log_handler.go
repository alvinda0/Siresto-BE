package handler

import (
	"project-name/internal/service"
	"project-name/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type APILogHandler struct {
	service service.APILogService
}

func NewAPILogHandler(service service.APILogService) *APILogHandler {
	return &APILogHandler{service: service}
}

// GetAllLogs godoc
// @Summary Get all API logs
// @Description Get all API logs with pagination and optional method filter. Internal users see all logs, external users only see their company/branch logs. Response body is excluded from list view.
// @Tags API Logs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param method query string false "Filter by HTTP method (GET, POST, PUT, DELETE)"
// @Security BearerAuth
// @Success 200 {object} pkg.Response
// @Failure 500 {object} pkg.Response
// @Router /api/logs [get]
func (h *APILogHandler) GetAllLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	method := c.Query("method")

	// Get user role type
	roleType, _ := c.Get("role_type")
	
	// Get company_id and branch_id for filtering
	var companyID, branchID string
	
	// If external user, filter by their company and branch
	if roleType == "EXTERNAL" {
		if cid, exists := c.Get("company_id"); exists {
			companyID = cid.(string)
		}
		if bid, exists := c.Get("branch_id"); exists {
			branchID = bid.(string)
		}
	}
	// If internal user, companyID and branchID remain empty (see all logs)

	logs, meta, err := h.service.GetAllLogs(page, limit, method, companyID, branchID)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch logs", err.Error())
		return
	}

	pkg.SuccessResponseWithMeta(c, http.StatusOK, "Logs retrieved successfully", logs, meta)
}

// GetLogByID godoc
// @Summary Get API log by ID
// @Description Get detailed information of a specific API log including response body. Internal users can see any log, external users only see logs from their company/branch
// @Tags API Logs
// @Accept json
// @Produce json
// @Param id path string true "Log UUID"
// @Security BearerAuth
// @Success 200 {object} pkg.Response
// @Failure 404 {object} pkg.Response
// @Failure 500 {object} pkg.Response
// @Router /api/logs/{id} [get]
func (h *APILogHandler) GetLogByID(c *gin.Context) {
	id := c.Param("id")
	
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid log ID format", "ID must be a valid UUID")
		return
	}

	// Get user role type
	roleType, _ := c.Get("role_type")
	
	// Get company_id and branch_id for filtering
	var companyID, branchID string
	
	// If external user, filter by their company and branch
	if roleType == "EXTERNAL" {
		if cid, exists := c.Get("company_id"); exists {
			companyID = cid.(string)
		}
		if bid, exists := c.Get("branch_id"); exists {
			branchID = bid.(string)
		}
	}
	// If internal user, companyID and branchID remain empty (can see any log)

	log, err := h.service.GetLogByID(id, companyID, branchID)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusNotFound, "Log not found", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Log retrieved successfully", log)
}
