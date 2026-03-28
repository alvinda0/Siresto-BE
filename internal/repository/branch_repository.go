package repository

import (
	"project-name/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BranchRepository interface {
	Create(branch *entity.Branch) error
	FindByID(id uuid.UUID) (*entity.Branch, error)
	FindByCompanyID(companyID uuid.UUID, limit, offset int) ([]entity.Branch, int64, error)
	Update(branch *entity.Branch) error
	Delete(id uuid.UUID) error
}

type branchRepository struct {
	DB *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{DB: db}
}

func (r *branchRepository) Create(branch *entity.Branch) error {
	return r.DB.Create(branch).Error
}

func (r *branchRepository) FindByID(id uuid.UUID) (*entity.Branch, error) {
	var branch entity.Branch
	err := r.DB.Where("id = ?", id).
		Preload("Company.Owner.Role").
		Preload("Company.Owner").
		Preload("Company").
		First(&branch).Error
	return &branch, err
}

func (r *branchRepository) FindByCompanyID(companyID uuid.UUID, limit, offset int) ([]entity.Branch, int64, error) {
	var branches []entity.Branch
	var total int64
	
	query := r.DB.Where("company_id = ?", companyID)
	
	// Count total
	query.Model(&entity.Branch{}).Count(&total)
	
	// Get paginated data
	err := query.Preload("Company.Owner.Role").
		Preload("Company.Owner").
		Preload("Company").
		Limit(limit).Offset(offset).
		Find(&branches).Error
	
	return branches, total, err
}

func (r *branchRepository) Update(branch *entity.Branch) error {
	return r.DB.Save(branch).Error
}

func (r *branchRepository) Delete(id uuid.UUID) error {
	return r.DB.Where("id = ?", id).Delete(&entity.Branch{}).Error
}
