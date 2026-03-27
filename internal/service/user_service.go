package service

import (
	"errors"
	"log"
	"project-name/internal/entity"
	"project-name/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterWithCompany(name, email, password, companyName, companyType string) (*entity.User, *entity.Company, error)
	Login(email, password string) (*entity.User, error)
	GetUserByID(id uuid.UUID) (*entity.User, error)
	GetAllUsers() ([]entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id uuid.UUID) error
	CreateExternalUser(name, email, password string, roleID uuid.UUID, companyID, branchID *uuid.UUID) (*entity.User, error)
	CreateInternalUser(name, email, password string, roleID uuid.UUID) (*entity.User, error)
	GetUsersByCompany(companyID uuid.UUID, limit, offset int) ([]entity.User, int64, error)
	GetUsersByCompanyFiltered(companyID, currentUserID uuid.UUID, limit, offset int) ([]entity.User, int64, error)
	GetUsersByBranch(branchID uuid.UUID, limit, offset int) ([]entity.User, int64, error)
	GetInternalUsers(limit, offset int) ([]entity.User, int64, error)
	GetExternalUsers(limit, offset int) ([]entity.User, int64, error)
}

type userService struct {
	repo        repository.UserRepository
	companyRepo *repository.CompanyRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo:        repo,
		companyRepo: nil,
	}
}

func NewUserServiceWithCompany(repo repository.UserRepository, companyRepo *repository.CompanyRepository) UserService {
	return &userService{
		repo:        repo,
		companyRepo: companyRepo,
	}
}

// RegisterWithCompany untuk registrasi owner dengan company sekaligus
func (s *userService) RegisterWithCompany(name, email, password, companyName, companyType string) (*entity.User, *entity.Company, error) {
	if s.companyRepo == nil {
		return nil, nil, errors.New("company repository not initialized")
	}
	
	// Get OWNER role ID
	var ownerRole entity.Role
	if err := s.repo.DB().Where("name = ?", "OWNER").First(&ownerRole).Error; err != nil {
		return nil, nil, errors.New("OWNER role not found")
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}
	
	// Create user
	user := &entity.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		RoleID:   ownerRole.ID,
		IsActive: true,
	}
	
	if err := s.repo.Create(user); err != nil {
		return nil, nil, err
	}
	
	// Create company
	var compType entity.CompanyType
	if companyType == "PT" {
		compType = entity.CompanyTypePT
	} else if companyType == "CV" {
		compType = entity.CompanyTypeCV
	} else {
		compType = entity.CompanyTypePerorangan
	}
	
	company := &entity.Company{
		Name:    companyName,
		Type:    compType,
		OwnerID: user.ID,
	}
	
	if err := s.companyRepo.Create(company); err != nil {
		// Rollback user creation if company creation fails
		s.repo.Delete(user.ID)
		return nil, nil, err
	}
	
	// Update user with company_id
	user.CompanyID = &company.ID
	if err := s.repo.Update(user); err != nil {
		return nil, nil, err
	}
	
	return user, company, nil
}

func (s *userService) Login(email, password string) (*entity.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		log.Printf("Login failed - user not found: %s, error: %v", email, err)
		return nil, errors.New("invalid credentials")
	}
	
	log.Printf("Login attempt - email: %s, stored hash length: %d", email, len(user.Password))
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Login failed - password mismatch for: %s, error: %v", email, err)
		return nil, errors.New("invalid credentials")
	}
	
	log.Printf("Login successful for: %s", email)
	return user, nil
}

func (s *userService) GetUserByID(id uuid.UUID) (*entity.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetAllUsers() ([]entity.User, error) {
	return s.repo.FindAll()
}

func (s *userService) UpdateUser(user *entity.User) error {
	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// CreateExternalUser untuk membuat user restoran (admin, cashier, dll)
func (s *userService) CreateExternalUser(name, email, password string, roleID uuid.UUID, companyID, branchID *uuid.UUID) (*entity.User, error) {
	// Verify role exists and is external
	var role entity.Role
	if err := s.repo.DB().Where("id = ? AND type = ?", roleID, "EXTERNAL").First(&role).Error; err != nil {
		return nil, errors.New("invalid external role")
	}

	// Validasi berdasarkan role name
	if role.Name == "ADMIN" && (companyID == nil || *companyID == uuid.Nil) {
		return nil, errors.New("admin harus terhubung dengan perusahaan")
	}

	if (role.Name == "CASHIER" || role.Name == "KITCHEN" || role.Name == "WAITER") {
		if companyID == nil || *companyID == uuid.Nil || branchID == nil || *branchID == uuid.Nil {
			return nil, errors.New("user harus terhubung dengan perusahaan dan cabang")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		RoleID:    roleID,
		CompanyID: companyID,
		BranchID:  branchID,
		IsActive:  true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// CreateInternalUser untuk membuat user internal platform (super_admin, support, finance)
func (s *userService) CreateInternalUser(name, email, password string, roleID uuid.UUID) (*entity.User, error) {
	// Verify role exists and is internal
	var role entity.Role
	if err := s.repo.DB().Where("id = ? AND type = ?", roleID, "INTERNAL").First(&role).Error; err != nil {
		return nil, errors.New("invalid internal role")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		RoleID:   roleID,
		IsActive: true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsersByCompany untuk mendapatkan semua user dalam perusahaan
func (s *userService) GetUsersByCompany(companyID uuid.UUID, limit, offset int) ([]entity.User, int64, error) {
	return s.repo.FindByCompanyID(companyID, limit, offset)
}

// GetUsersByCompanyFiltered untuk mendapatkan user berdasarkan role yang login
func (s *userService) GetUsersByCompanyFiltered(companyID, currentUserID uuid.UUID, limit, offset int) ([]entity.User, int64, error) {
	// Get current user info
	currentUser, err := s.repo.FindByID(currentUserID)
	if err != nil {
		return nil, 0, err
	}

	// Jika OWNER, tampilkan semua user di company
	if currentUser.Role.Name == "OWNER" {
		return s.repo.FindByCompanyID(companyID, limit, offset)
	}

	// Jika ADMIN, hanya tampilkan user di cabang yang dia urus
	if currentUser.Role.Name == "ADMIN" && currentUser.BranchID != nil {
		return s.repo.FindByBranchID(*currentUser.BranchID, limit, offset)
	}

	// Jika role lain (CASHIER, KITCHEN, WAITER), tidak boleh akses endpoint ini
	return nil, 0, errors.New("unauthorized to view company users")
}

// GetUsersByBranch untuk mendapatkan semua user dalam cabang
func (s *userService) GetUsersByBranch(branchID uuid.UUID, limit, offset int) ([]entity.User, int64, error) {
	return s.repo.FindByBranchID(branchID, limit, offset)
}

// GetInternalUsers untuk mendapatkan semua internal user platform
func (s *userService) GetInternalUsers(limit, offset int) ([]entity.User, int64, error) {
	return s.repo.FindInternalUsers(limit, offset)
}

// GetExternalUsers untuk mendapatkan semua external user (client restoran)
func (s *userService) GetExternalUsers(limit, offset int) ([]entity.User, int64, error) {
	return s.repo.FindExternalUsers(limit, offset)
}
