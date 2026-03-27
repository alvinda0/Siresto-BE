package service

import (
	"errors"
	"project-name/internal/entity"
	"project-name/internal/repository"

	"github.com/google/uuid"
)

type CompanyService struct {
	companyRepo *repository.CompanyRepository
	userRepo    repository.UserRepository
}

func NewCompanyService(companyRepo *repository.CompanyRepository, userRepo repository.UserRepository) *CompanyService {
	return &CompanyService{
		companyRepo: companyRepo,
		userRepo:    userRepo,
	}
}

func (s *CompanyService) CreateCompany(name string, companyType entity.CompanyType, ownerID uuid.UUID) (*entity.Company, error) {
	if name == "" {
		return nil, errors.New("nama perusahaan tidak boleh kosong")
	}
	
	if companyType != entity.CompanyTypePT && companyType != entity.CompanyTypePerorangan {
		return nil, errors.New("tipe perusahaan harus PT atau PERORANGAN")
	}

	company := &entity.Company{
		Name:    name,
		Type:    companyType,
		OwnerID: ownerID,
	}

	if err := s.companyRepo.Create(company); err != nil {
		return nil, err
	}

	return company, nil
}

func (s *CompanyService) GetCompanyByID(id uuid.UUID) (*entity.Company, error) {
	return s.companyRepo.FindByID(id)
}

func (s *CompanyService) GetCompaniesByOwner(ownerID uuid.UUID) ([]entity.Company, error) {
	return s.companyRepo.FindByOwnerID(ownerID)
}

func (s *CompanyService) GetMyCompanies(userID uuid.UUID) ([]entity.Company, error) {
	// Get user info untuk cek role dan branch
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Jika user adalah owner, tampilkan semua company miliknya dengan semua cabang
	companies, err := s.companyRepo.FindByOwnerID(userID)
	if err != nil {
		return nil, err
	}
	
	// Jika ada company sebagai owner, return semua
	if len(companies) > 0 {
		return companies, nil
	}

	// Jika bukan owner, cek apakah user punya company_id (admin/staff)
	if user.CompanyID != nil {
		// Get company user tersebut
		company, err := s.companyRepo.FindByID(*user.CompanyID)
		if err != nil {
			return nil, err
		}

		// Jika user punya branch_id, filter hanya cabang yang dia urus
		if user.BranchID != nil {
			filteredBranches := []entity.Branch{}
			for _, branch := range company.Branches {
				if branch.ID == *user.BranchID {
					filteredBranches = append(filteredBranches, branch)
					break
				}
			}
			company.Branches = filteredBranches
		}

		return []entity.Company{*company}, nil
	}

	// Jika tidak ada company sama sekali
	return []entity.Company{}, nil
}

func (s *CompanyService) UpdateCompany(id uuid.UUID, name string, companyType entity.CompanyType) error {
	company, err := s.companyRepo.FindByID(id)
	if err != nil {
		return err
	}

	if name != "" {
		company.Name = name
	}
	if companyType != "" {
		company.Type = companyType
	}

	return s.companyRepo.Update(company)
}

func (s *CompanyService) DeleteCompany(id uuid.UUID) error {
	return s.companyRepo.Delete(id)
}
