package repository

import (
	"project-name/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BranchRepository struct {
	DB *gorm.DB
}

func NewBranchRepository(db *gorm.DB) *BranchRepository {
	return &BranchRepository{DB: db}
}

func (r *BranchRepository) Create(branch *entity.Branch) error {
	return r.DB.Create(branch).Error
}

func (r *BranchRepository) FindByID(id uuid.UUID) (*entity.Branch, error) {
	var branch entity.Branch
	err := r.DB.Where("id = ?", id).
		Preload("Company.Owner.Role").
		Preload("Company.Owner").
		Preload("Company").
		First(&branch).Error
	return &branch, err
}

func (r *BranchRepository) FindByCompanyID(companyID uuid.UUID, limit, offset int) ([]entity.Branch, int64, error) {
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

func (r *BranchRepository) Update(branch *entity.Branch) error {
	return r.DB.Save(branch).Error
}

func (r *BranchRepository) Delete(id uuid.UUID) error {
	return r.DB.Where("id = ?", id).Delete(&entity.Branch{}).Error
}
