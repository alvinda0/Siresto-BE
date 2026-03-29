package repository

import (
	"project-name/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaxRepository interface {
	Create(tax *entity.Tax) error
	Update(tax *entity.Tax) error
	Delete(id, companyID uuid.UUID, branchID *uuid.UUID) error
	FindByID(id, companyID uuid.UUID, branchID *uuid.UUID) (*entity.Tax, error)
	FindByCompany(companyID uuid.UUID, page, limit int) ([]entity.Tax, int64, error)
	FindByBranch(companyID, branchID uuid.UUID, page, limit int) ([]entity.Tax, int64, error)
}

type taxRepository struct {
	db *gorm.DB
}

func NewTaxRepository(db *gorm.DB) TaxRepository {
	return &taxRepository{db: db}
}

func (r *taxRepository) Create(tax *entity.Tax) error {
	return r.db.Create(tax).Error
}

func (r *taxRepository) Update(tax *entity.Tax) error {
	return r.db.Save(tax).Error
}

func (r *taxRepository) Delete(id, companyID uuid.UUID, branchID *uuid.UUID) error {
	// Delete logic:
	// - If user has no branch_id (OWNER): can delete company-level taxes only
	// - If user has branch_id (ADMIN, etc): can delete their branch taxes only
	// Note: Permission check for company-level tax is done in handler
	
	query := r.db.Where("id = ? AND company_id = ?", id, companyID)
	
	if branchID != nil {
		// User with branch can only delete their own branch taxes
		query = query.Where("branch_id = ?", *branchID)
	} else {
		// User without branch (OWNER) can only delete company-level taxes
		query = query.Where("branch_id IS NULL")
	}
	
	return query.Delete(&entity.Tax{}).Error
}

func (r *taxRepository) FindByID(id, companyID uuid.UUID, branchID *uuid.UUID) (*entity.Tax, error) {
	var tax entity.Tax
	query := r.db.Preload("Company").Preload("Branch").
		Where("id = ? AND company_id = ?", id, companyID)
	
	// Users can view:
	// - Company-level taxes (branch_id IS NULL) - visible to all
	// - Their own branch taxes (branch_id = user's branch_id)
	if branchID != nil {
		query = query.Where("(branch_id IS NULL OR branch_id = ?)", *branchID)
	} else {
		// OWNER can view company-level taxes only
		query = query.Where("branch_id IS NULL")
	}
	
	err := query.First(&tax).Error
	if err != nil {
		return nil, err
	}
	return &tax, nil
}

func (r *taxRepository) FindByCompany(companyID uuid.UUID, page, limit int) ([]entity.Tax, int64, error) {
	var taxes []entity.Tax
	var total int64
	
	// Count total
	r.db.Model(&entity.Tax{}).
		Where("company_id = ? AND branch_id IS NULL", companyID).
		Count(&total)
	
	// Get paginated data
	offset := (page - 1) * limit
	err := r.db.Preload("Company").Preload("Branch").
		Where("company_id = ? AND branch_id IS NULL", companyID).
		Order("prioritas DESC, nama_pajak ASC").
		Limit(limit).
		Offset(offset).
		Find(&taxes).Error
	
	return taxes, total, err
}

func (r *taxRepository) FindByBranch(companyID, branchID uuid.UUID, page, limit int) ([]entity.Tax, int64, error) {
	var taxes []entity.Tax
	var total int64
	
	// Count total (both company-level and branch-specific)
	r.db.Model(&entity.Tax{}).
		Where("company_id = ? AND (branch_id IS NULL OR branch_id = ?)", companyID, branchID).
		Count(&total)
	
	// Get paginated data
	offset := (page - 1) * limit
	err := r.db.Preload("Company").Preload("Branch").
		Where("company_id = ? AND (branch_id IS NULL OR branch_id = ?)", companyID, branchID).
		Order("prioritas DESC, nama_pajak ASC").
		Limit(limit).
		Offset(offset).
		Find(&taxes).Error
	
	return taxes, total, err
}
