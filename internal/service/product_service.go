package service

import (
	"project-name/internal/entity"
	"project-name/internal/repository"
	"errors"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(product *entity.Product) error
	GetAllProducts(companyID, branchID uuid.UUID, search string, page, limit int) ([]entity.Product, int64, error)
	GetProductByID(id, companyID, branchID uuid.UUID) (*entity.Product, error)
	UpdateProduct(product *entity.Product) error
	DeleteProduct(id, companyID, branchID uuid.UUID) error
}

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	branchRepo   repository.BranchRepository
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, branchRepo repository.BranchRepository) ProductService {
	return &productService{productRepo, categoryRepo, branchRepo}
}

func (s *productService) CreateProduct(product *entity.Product) error {
	// Validate category exists and belongs to company
	category, err := s.categoryRepo.FindByID(product.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}
	if category.CompanyID != product.CompanyID {
		return errors.New("category doesn't belong to your company")
	}

	// Validate branch exists and belongs to company
	branch, err := s.branchRepo.FindByID(product.BranchID)
	if err != nil {
		return errors.New("branch not found")
	}
	if branch.CompanyID != product.CompanyID {
		return errors.New("branch doesn't belong to your company")
	}

	return s.productRepo.Create(product)
}

func (s *productService) GetAllProducts(companyID, branchID uuid.UUID, search string, page, limit int) ([]entity.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Validate that branch belongs to company
	branch, err := s.branchRepo.FindByID(branchID)
	if err != nil {
		return nil, 0, errors.New("branch not found")
	}
	if branch.CompanyID != companyID {
		return nil, 0, errors.New("branch doesn't belong to your company")
	}

	return s.productRepo.FindAll(companyID, branchID, search, page, limit)
}

func (s *productService) GetProductByID(id, companyID, branchID uuid.UUID) (*entity.Product, error) {
	return s.productRepo.FindByID(id, companyID, branchID)
}

func (s *productService) UpdateProduct(product *entity.Product) error {
	// Check if product exists
	existing, err := s.productRepo.FindByID(product.ID, product.CompanyID, product.BranchID)
	if err != nil {
		return errors.New("product not found")
	}

	// Validate category if changed
	if product.CategoryID != existing.CategoryID {
		category, err := s.categoryRepo.FindByID(product.CategoryID)
		if err != nil {
			return errors.New("category not found")
		}
		if category.CompanyID != product.CompanyID {
			return errors.New("category doesn't belong to your company")
		}
	}

	return s.productRepo.Update(product)
}

func (s *productService) DeleteProduct(id, companyID, branchID uuid.UUID) error {
	// Check if product exists
	_, err := s.productRepo.FindByID(id, companyID, branchID)
	if err != nil {
		return errors.New("product not found")
	}

	return s.productRepo.Delete(id, companyID, branchID)
}
