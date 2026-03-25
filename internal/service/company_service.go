package service

import (
	"errors"
	"project-name/internal/entity"
	"project-name/internal/repository"

	"github.com/google/uuid"
)

type CompanyService struct {
	companyRepo *repository.CompanyRepository
}

func NewCompanyService(companyRepo *repository.CompanyRepository) *CompanyService {
	return &CompanyService{companyRepo: companyRepo}
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
