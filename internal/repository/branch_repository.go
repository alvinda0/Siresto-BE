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

func (r *BranchRepository) FindByCompanyID(companyID uuid.UUID) ([]entity.Branch, error) {
	var branches []entity.Branch
	err := r.DB.Where("company_id = ?", companyID).
		Preload("Company.Owner.Role").
		Preload("Company.Owner").
		Preload("Company").
		Find(&branches).Error
	return branches, err
}

func (r *BranchRepository) Update(branch *entity.Branch) error {
	return r.DB.Save(branch).Error
}

func (r *BranchRepository) Delete(id uuid.UUID) error {
	return r.DB.Where("id = ?", id).Delete(&entity.Branch{}).Error
}
