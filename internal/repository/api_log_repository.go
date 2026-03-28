package repository

import (
	"project-name/internal/entity"

	"gorm.io/gorm"
)

type APILogRepository interface {
	Create(log *entity.APILog) error
	FindAll(page, limit int, method, companyID, branchID string) ([]entity.APILog, int64, error)
	FindByID(id string, companyID, branchID string) (*entity.APILog, error)
}

type apiLogRepository struct {
	db *gorm.DB
}

func NewAPILogRepository(db *gorm.DB) APILogRepository {
	return &apiLogRepository{db: db}
}

func (r *apiLogRepository) Create(log *entity.APILog) error {
	return r.db.Create(log).Error
}

func (r *apiLogRepository) FindAll(page, limit int, method, companyID, branchID string) ([]entity.APILog, int64, error) {
	var logs []entity.APILog
	var total int64

	query := r.db.Model(&entity.APILog{})

	// Filter by method if provided
	if method != "" {
		query = query.Where("method = ?", method)
	}

	// Filter by company_id if provided (for external users)
	if companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	// Filter by branch_id if provided (for external users)
	if branchID != "" {
		query = query.Where("branch_id = ?", branchID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated results
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, total, err
}

func (r *apiLogRepository) FindByID(id string, companyID, branchID string) (*entity.APILog, error) {
	var log entity.APILog
	query := r.db.Model(&entity.APILog{})

	// Filter by company_id if provided (for external users)
	if companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	// Filter by branch_id if provided (for external users)
	if branchID != "" {
		query = query.Where("branch_id = ?", branchID)
	}

	err := query.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}
