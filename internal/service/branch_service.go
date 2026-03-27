package service

import (
	"errors"
	"project-name/internal/entity"
	"project-name/internal/repository"

	"github.com/google/uuid"
)

type BranchService struct {
	branchRepo  *repository.BranchRepository
	companyRepo *repository.CompanyRepository
	userRepo    repository.UserRepository
}

func NewBranchService(branchRepo *repository.BranchRepository, companyRepo *repository.CompanyRepository, userRepo repository.UserRepository) *BranchService {
	return &BranchService{
		branchRepo:  branchRepo,
		companyRepo: companyRepo,
		userRepo:    userRepo,
	}
}

func (s *BranchService) CreateBranch(companyID uuid.UUID, name, address, city, province, postalCode, phone string) (*entity.Branch, error) {
	// Validasi company exists
	_, err := s.companyRepo.FindByID(companyID)
	if err != nil {
		return nil, errors.New("perusahaan tidak ditemukan")
	}

	if name == "" {
		return nil, errors.New("nama cabang tidak boleh kosong")
	}
	if address == "" {
		return nil, errors.New("alamat tidak boleh kosong")
	}

	branch := &entity.Branch{
		CompanyID:  companyID,
		Name:       name,
		Address:    address,
		City:       city,
		Province:   province,
		PostalCode: postalCode,
		Phone:      phone,
		IsActive:   true,
	}

	if err := s.branchRepo.Create(branch); err != nil {
		return nil, err
	}

	// Load branch with company relation
	return s.branchRepo.FindByID(branch.ID)
}

func (s *BranchService) GetBranchByID(id uuid.UUID) (*entity.Branch, error) {
	return s.branchRepo.FindByID(id)
}

func (s *BranchService) GetBranchesByCompany(companyID uuid.UUID, limit, offset int) ([]entity.Branch, int64, error) {
	return s.branchRepo.FindByCompanyID(companyID, limit, offset)
}

func (s *BranchService) GetBranchesByCompanyFiltered(companyID, currentUserID uuid.UUID, limit, offset int) ([]entity.Branch, int64, error) {
	// Get current user info
	currentUser, err := s.userRepo.FindByID(currentUserID)
	if err != nil {
		return nil, 0, err
	}

	// Jika OWNER, tampilkan semua cabang di company
	if currentUser.Role.Name == "OWNER" {
		return s.branchRepo.FindByCompanyID(companyID, limit, offset)
	}

	// Jika ADMIN atau role lain yang punya branch_id, hanya tampilkan cabang yang dia urus
	if currentUser.BranchID != nil {
		branch, err := s.branchRepo.FindByID(*currentUser.BranchID)
		if err != nil {
			return nil, 0, err
		}
		
		// Pastikan branch tersebut memang milik company yang diminta
		if branch.CompanyID != companyID {
			return []entity.Branch{}, 0, nil
		}
		
		return []entity.Branch{*branch}, 1, nil
	}

	// Jika tidak punya branch_id, return empty
	return []entity.Branch{}, 0, nil
}

func (s *BranchService) UpdateBranch(id uuid.UUID, name, address, city, province, postalCode, phone string, isActive *bool) error {
	branch, err := s.branchRepo.FindByID(id)
	if err != nil {
		return err
	}

	if name != "" {
		branch.Name = name
	}
	if address != "" {
		branch.Address = address
	}
	if city != "" {
		branch.City = city
	}
	if province != "" {
		branch.Province = province
	}
	if postalCode != "" {
		branch.PostalCode = postalCode
	}
	if phone != "" {
		branch.Phone = phone
	}
	if isActive != nil {
		branch.IsActive = *isActive
	}

	return s.branchRepo.Update(branch)
}

func (s *BranchService) DeleteBranch(id uuid.UUID) error {
	return s.branchRepo.Delete(id)
}
