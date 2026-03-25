package pkg

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Page  int
	Limit int
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// GetPaginationParams mengambil page dan limit dari query params
func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validasi page minimal 1
	if page < 1 {
		page = 1
	}

	// Validasi limit hanya boleh nilai tertentu
	validLimits := []int{10, 50, 100, 200, 500, 1000}
	isValidLimit := false
	for _, validLimit := range validLimits {
		if limit == validLimit {
			isValidLimit = true
			break
		}
	}

	// Jika limit tidak valid, set ke default 10
	if !isValidLimit {
		limit = 10
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// CalculateOffset menghitung offset untuk query database
func (p *PaginationParams) CalculateOffset() int {
	return (p.Page - 1) * p.Limit
}

// CreateMeta membuat metadata pagination
func (p *PaginationParams) CreateMeta(totalItems int64) PaginationMeta {
	totalPages := int(totalItems) / p.Limit
	if int(totalItems)%p.Limit > 0 {
		totalPages++
	}

	return PaginationMeta{
		Page:       p.Page,
		Limit:      p.Limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
