package repository

import (
	"project-name/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	DB *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{DB: db}
}

func (r *CompanyRepository) Create(company *entity.Company) error {
	return r.DB.Create(company).Error
}

func (r *CompanyRepository) FindByID(id uuid.UUID) (*entity.Company, error) {
	var company entity.Company
	err := r.DB.Where("id = ?", id).
		Preload("Owner.Role").
		Preload("Owner").
		Preload("Branches").
		First(&company).Error
	return &company, err
}

func (r *CompanyRepository) FindByOwnerID(ownerID uuid.UUID) ([]entity.Company, error) {
	var companies []entity.Company
	err := r.DB.Where("owner_id = ?", ownerID).
		Preload("Owner.Role").
		Preload("Owner").
		Preload("Branches").
		Find(&companies).Error
	return companies, err
}

func (r *CompanyRepository) Update(company *entity.Company) error {
	return r.DB.Save(company).Error
}

func (r *CompanyRepository) Delete(id uuid.UUID) error {
	return r.DB.Where("id = ?", id).Delete(&entity.Company{}).Error
}
