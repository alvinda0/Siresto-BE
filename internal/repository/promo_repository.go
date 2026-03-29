package repository

import (
	"project-name/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PromoRepository interface {
	Create(promo *entity.Promo) error
	Update(promo *entity.Promo) error
	Delete(id, companyID uuid.UUID, branchID *uuid.UUID) error
	FindByID(id, companyID uuid.UUID, branchID *uuid.UUID) (*entity.Promo, error)
	FindByCode(code string, companyID uuid.UUID, branchID *uuid.UUID) (*entity.Promo, error)
	FindByCompany(companyID uuid.UUID, page, limit int) ([]entity.Promo, int64, error)
	FindByBranch(companyID, branchID uuid.UUID, page, limit int) ([]entity.Promo, int64, error)
}

type promoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) PromoRepository {
	return &promoRepository{db: db}
}

func (r *promoRepository) Create(promo *entity.Promo) error {
	return r.db.Create(promo).Error
}

func (r *promoRepository) Update(promo *entity.Promo) error {
	return r.db.Save(promo).Error
}

func (r *promoRepository) Delete(id, companyID uuid.UUID, branchID *uuid.UUID) error {
	query := r.db.Where("id = ? AND company_id = ?", id, companyID)
	
	if branchID != nil {
		// User with branch can only delete their own branch promos
		query = query.Where("branch_id = ?", *branchID)
	} else {
		// User without branch (OWNER) can only delete company-level promos
		query = query.Where("branch_id IS NULL")
	}
	
	return query.Delete(&entity.Promo{}).Error
}

func (r *promoRepository) FindByID(id, companyID uuid.UUID, branchID *uuid.UUID) (*entity.Promo, error) {
	var promo entity.Promo
	query := r.db.Preload("Company").Preload("Branch").
		Where("id = ? AND company_id = ?", id, companyID)
	
	// Users can view:
	// - Company-level promos (branch_id IS NULL) - visible to all
	// - Their own branch promos (branch_id = user's branch_id)
	if branchID != nil {
		query = query.Where("(branch_id IS NULL OR branch_id = ?)", *branchID)
	} else {
		// OWNER can view company-level promos only
		query = query.Where("branch_id IS NULL")
	}
	
	err := query.First(&promo).Error
	if err != nil {
		return nil, err
	}
	return &promo, nil
}

func (r *promoRepository) FindByCode(code string, companyID uuid.UUID, branchID *uuid.UUID) (*entity.Promo, error) {
	var promo entity.Promo
	query := r.db.Preload("Company").Preload("Branch").
		Where("code = ? AND company_id = ?", code, companyID)
	
	if branchID != nil {
		query = query.Where("(branch_id IS NULL OR branch_id = ?)", *branchID)
	} else {
		query = query.Where("branch_id IS NULL")
	}
	
	err := query.First(&promo).Error
	if err != nil {
		return nil, err
	}
	return &promo, nil
}

func (r *promoRepository) FindByCompany(companyID uuid.UUID, page, limit int) ([]entity.Promo, int64, error) {
	var promos []entity.Promo
	var total int64
	
	// Count total
	r.db.Model(&entity.Promo{}).
		Where("company_id = ? AND branch_id IS NULL", companyID).
		Count(&total)
	
	// Get paginated data
	offset := (page - 1) * limit
	err := r.db.Preload("Company").Preload("Branch").
		Where("company_id = ? AND branch_id IS NULL", companyID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&promos).Error
	
	return promos, total, err
}

func (r *promoRepository) FindByBranch(companyID, branchID uuid.UUID, page, limit int) ([]entity.Promo, int64, error) {
	var promos []entity.Promo
	var total int64
	
	// Count total (both company-level and branch-specific)
	r.db.Model(&entity.Promo{}).
		Where("company_id = ? AND (branch_id IS NULL OR branch_id = ?)", companyID, branchID).
		Count(&total)
	
	// Get paginated data
	offset := (page - 1) * limit
	err := r.db.Preload("Company").Preload("Branch").
		Where("company_id = ? AND (branch_id IS NULL OR branch_id = ?)", companyID, branchID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&promos).Error
	
	return promos, total, err
}
