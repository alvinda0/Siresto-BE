package service

import (
	"errors"
	"project-name/internal/entity"
	"project-name/internal/repository"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateCategory(category *entity.Category) error
	UpdateCategory(id uuid.UUID, category *entity.Category) error
	DeleteCategory(id uuid.UUID) error
	GetCategoryByID(id uuid.UUID) (*entity.Category, error)
	GetCategoriesByCompany(companyID uuid.UUID, branchID *uuid.UUID) ([]entity.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	companyRepo  *repository.CompanyRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository, companyRepo *repository.CompanyRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
		companyRepo:  companyRepo,
	}
}

func (s *categoryService) CreateCategory(category *entity.Category) error {
	// Validasi company exists
	_, err := s.companyRepo.FindByID(category.CompanyID)
	if err != nil {
		return errors.New("company not found")
	}

	// Set position otomatis jika tidak diset atau 0
	if category.Position <= 0 {
		maxPos, err := s.categoryRepo.GetMaxPosition(category.CompanyID, category.BranchID)
		if err != nil {
			return err
		}
		category.Position = maxPos + 1
	}

	return s.categoryRepo.Create(category)
}

func (s *categoryService) UpdateCategory(id uuid.UUID, category *entity.Category) error {
	existing, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	// Update fields
	existing.Name = category.Name
	existing.Description = category.Description
	existing.IsActive = category.IsActive
	
	// Update position jika diubah dan valid
	if category.Position > 0 {
		existing.Position = category.Position
	}

	return s.categoryRepo.Update(existing)
}

func (s *categoryService) DeleteCategory(id uuid.UUID) error {
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	return s.categoryRepo.Delete(id)
}

func (s *categoryService) GetCategoryByID(id uuid.UUID) (*entity.Category, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *categoryService) GetCategoriesByCompany(companyID uuid.UUID, branchID *uuid.UUID) ([]entity.Category, error) {
	return s.categoryRepo.FindByCompanyID(companyID, branchID)
}
