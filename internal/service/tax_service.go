package service

import (
	"project-name/internal/entity"
	"project-name/internal/repository"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaxService interface {
	CreateTax(companyID uuid.UUID, branchID *uuid.UUID, req *entity.CreateTaxRequest) (*entity.TaxResponse, error)
	UpdateTax(id, companyID uuid.UUID, branchID *uuid.UUID, req *entity.UpdateTaxRequest) (*entity.TaxResponse, error)
	DeleteTax(id, companyID uuid.UUID, branchID *uuid.UUID) error
	GetTaxByID(id, companyID uuid.UUID, branchID *uuid.UUID) (*entity.TaxResponse, error)
	GetAllTaxes(companyID uuid.UUID, branchID *uuid.UUID, page, limit int) ([]entity.TaxResponse, map[string]interface{}, error)
}

type taxService struct {
	taxRepo repository.TaxRepository
}

func NewTaxService(taxRepo repository.TaxRepository) TaxService {
	return &taxService{taxRepo: taxRepo}
}

func (s *taxService) CreateTax(companyID uuid.UUID, branchID *uuid.UUID, req *entity.CreateTaxRequest) (*entity.TaxResponse, error) {
	// Set default status if not provided
	if req.Status == "" {
		req.Status = "active"
	}

	// Determine branch_id: use from request if provided, otherwise use from context
	var finalBranchID *uuid.UUID
	if req.BranchID != nil {
		finalBranchID = req.BranchID
	} else {
		finalBranchID = branchID
	}

	tax := &entity.Tax{
		CompanyID:  companyID,
		BranchID:   finalBranchID,
		NamaPajak:  req.NamaPajak,
		TipePajak:  req.TipePajak,
		Presentase: req.Presentase,
		Deskripsi:  req.Deskripsi,
		Status:     req.Status,
		Prioritas:  req.Prioritas,
	}

	if err := s.taxRepo.Create(tax); err != nil {
		return nil, err
	}

	// Reload with relations
	reloadedTax, err := s.taxRepo.FindByID(tax.ID, companyID, finalBranchID)
	if err != nil {
		// If reload fails, return without relations
		return s.toResponse(tax), nil
	}

	return s.toResponse(reloadedTax), nil
}

func (s *taxService) UpdateTax(id, companyID uuid.UUID, branchID *uuid.UUID, req *entity.UpdateTaxRequest) (*entity.TaxResponse, error) {
	tax, err := s.taxRepo.FindByID(id, companyID, branchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tax not found")
		}
		return nil, err
	}

	// Update only provided fields
	if req.NamaPajak != "" {
		tax.NamaPajak = req.NamaPajak
	}
	if req.TipePajak != "" {
		tax.TipePajak = req.TipePajak
	}
	if req.Presentase > 0 {
		tax.Presentase = req.Presentase
	}
	if req.Deskripsi != "" {
		tax.Deskripsi = req.Deskripsi
	}
	if req.Status != "" {
		tax.Status = req.Status
	}
	if req.Prioritas != 0 {
		tax.Prioritas = req.Prioritas
	}

	if err := s.taxRepo.Update(tax); err != nil {
		return nil, err
	}

	return s.toResponse(tax), nil
}

func (s *taxService) DeleteTax(id, companyID uuid.UUID, branchID *uuid.UUID) error {
	_, err := s.taxRepo.FindByID(id, companyID, branchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tax not found")
		}
		return err
	}

	return s.taxRepo.Delete(id, companyID, branchID)
}

func (s *taxService) GetTaxByID(id, companyID uuid.UUID, branchID *uuid.UUID) (*entity.TaxResponse, error) {
	tax, err := s.taxRepo.FindByID(id, companyID, branchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tax not found")
		}
		return nil, err
	}

	return s.toResponse(tax), nil
}

func (s *taxService) GetAllTaxes(companyID uuid.UUID, branchID *uuid.UUID, page, limit int) ([]entity.TaxResponse, map[string]interface{}, error) {
	var taxes []entity.Tax
	var total int64
	var err error

	if branchID != nil {
		// Get both company-level and branch-specific taxes
		taxes, total, err = s.taxRepo.FindByBranch(companyID, *branchID, page, limit)
	} else {
		// Get only company-level taxes
		taxes, total, err = s.taxRepo.FindByCompany(companyID, page, limit)
	}

	if err != nil {
		return nil, nil, err
	}

	responses := make([]entity.TaxResponse, len(taxes))
	for i, tax := range taxes {
		responses[i] = *s.toResponse(&tax)
	}

	// Calculate pagination metadata
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
	}

	return responses, meta, nil
}

func (s *taxService) toResponse(tax *entity.Tax) *entity.TaxResponse {
	response := &entity.TaxResponse{
		ID:          tax.ID,
		CompanyID:   tax.CompanyID,
		CompanyName: tax.Company.Name,
		BranchID:    tax.BranchID,
		NamaPajak:   tax.NamaPajak,
		TipePajak:   tax.TipePajak,
		Presentase:  tax.Presentase,
		Deskripsi:   tax.Deskripsi,
		Status:      tax.Status,
		Prioritas:   tax.Prioritas,
		CreatedAt:   tax.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   tax.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	
	// Add branch name if exists
	if tax.Branch != nil {
		response.BranchName = &tax.Branch.Name
	}
	
	return response
}
