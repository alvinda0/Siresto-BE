package service

import (
	"errors"
	"project-name/internal/entity"
	"project-name/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PromoService interface {
	CreatePromo(companyID uuid.UUID, branchID *uuid.UUID, req *entity.CreatePromoRequest) (*entity.PromoResponse, error)
	UpdatePromo(id, companyID uuid.UUID, branchID *uuid.UUID, req *entity.UpdatePromoRequest) (*entity.PromoResponse, error)
	DeletePromo(id, companyID uuid.UUID, branchID *uuid.UUID) error
	GetPromoByID(id, companyID uuid.UUID, branchID *uuid.UUID) (*entity.PromoResponse, error)
	GetPromoByCode(code string, companyID uuid.UUID, branchID *uuid.UUID) (*entity.PromoResponse, error)
	GetAllPromos(companyID uuid.UUID, branchID *uuid.UUID, page, limit int) ([]entity.PromoResponse, map[string]interface{}, error)
}

type promoService struct {
	promoRepo repository.PromoRepository
}

func NewPromoService(promoRepo repository.PromoRepository) PromoService {
	return &promoService{promoRepo: promoRepo}
}

func (s *promoService) CreatePromo(companyID uuid.UUID, branchID *uuid.UUID, req *entity.CreatePromoRequest) (*entity.PromoResponse, error) {
	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errors.New("invalid end_date format, use YYYY-MM-DD")
	}

	// Validate dates
	if endDate.Before(startDate) {
		return nil, errors.New("end_date must be after start_date")
	}

	// Set default is_active
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Determine branch_id
	var finalBranchID *uuid.UUID
	if req.BranchID != nil {
		finalBranchID = req.BranchID
	} else {
		finalBranchID = branchID
	}

	promo := &entity.Promo{
		CompanyID:      companyID,
		BranchID:       finalBranchID,
		Name:           req.Name,
		Code:           req.Code,
		Type:           req.Type,
		Value:          req.Value,
		MaxDiscount:    req.MaxDiscount,
		MinTransaction: req.MinTransaction,
		Quota:          req.Quota,
		UsedCount:      0,
		StartDate:      startDate,
		EndDate:        endDate,
		IsActive:       isActive,
	}

	if err := s.promoRepo.Create(promo); err != nil {
		return nil, err
	}

	// Reload with relations
	reloadedPromo, err := s.promoRepo.FindByID(promo.ID, companyID, finalBranchID)
	if err != nil {
		return s.toResponse(promo), nil
	}

	return s.toResponse(reloadedPromo), nil
}

func (s *promoService) UpdatePromo(id, companyID uuid.UUID, branchID *uuid.UUID, req *entity.UpdatePromoRequest) (*entity.PromoResponse, error) {
	promo, err := s.promoRepo.FindByID(id, companyID, branchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("promo not found")
		}
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		promo.Name = req.Name
	}
	if req.Code != "" {
		promo.Code = req.Code
	}
	if req.Type != "" {
		promo.Type = req.Type
	}
	if req.Value > 0 {
		promo.Value = req.Value
	}
	if req.MaxDiscount != nil {
		promo.MaxDiscount = req.MaxDiscount
	}
	if req.MinTransaction != nil {
		promo.MinTransaction = req.MinTransaction
	}
	if req.Quota != nil {
		promo.Quota = req.Quota
	}
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, errors.New("invalid start_date format, use YYYY-MM-DD")
		}
		promo.StartDate = startDate
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end_date format, use YYYY-MM-DD")
		}
		promo.EndDate = endDate
	}
	if req.IsActive != nil {
		promo.IsActive = *req.IsActive
	}

	// Validate dates
	if promo.EndDate.Before(promo.StartDate) {
		return nil, errors.New("end_date must be after start_date")
	}

	if err := s.promoRepo.Update(promo); err != nil {
		return nil, err
	}

	return s.toResponse(promo), nil
}

func (s *promoService) DeletePromo(id, companyID uuid.UUID, branchID *uuid.UUID) error {
	_, err := s.promoRepo.FindByID(id, companyID, branchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("promo not found")
		}
		return err
	}

	return s.promoRepo.Delete(id, companyID, branchID)
}

func (s *promoService) GetPromoByID(id, companyID uuid.UUID, branchID *uuid.UUID) (*entity.PromoResponse, error) {
	promo, err := s.promoRepo.FindByID(id, companyID, branchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("promo not found")
		}
		return nil, err
	}

	return s.toResponse(promo), nil
}

func (s *promoService) GetPromoByCode(code string, companyID uuid.UUID, branchID *uuid.UUID) (*entity.PromoResponse, error) {
	promo, err := s.promoRepo.FindByCode(code, companyID, branchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("promo not found")
		}
		return nil, err
	}

	return s.toResponse(promo), nil
}

func (s *promoService) GetAllPromos(companyID uuid.UUID, branchID *uuid.UUID, page, limit int) ([]entity.PromoResponse, map[string]interface{}, error) {
	var promos []entity.Promo
	var total int64
	var err error

	if branchID != nil {
		promos, total, err = s.promoRepo.FindByBranch(companyID, *branchID, page, limit)
	} else {
		promos, total, err = s.promoRepo.FindByCompany(companyID, page, limit)
	}

	if err != nil {
		return nil, nil, err
	}

	responses := make([]entity.PromoResponse, len(promos))
	for i, promo := range promos {
		responses[i] = *s.toResponse(&promo)
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

func (s *promoService) toResponse(promo *entity.Promo) *entity.PromoResponse {
	now := time.Now()
	isExpired := now.After(promo.EndDate)
	
	// Calculate remaining quota
	var remainingQuota *int
	if promo.Quota != nil {
		remaining := *promo.Quota - promo.UsedCount
		if remaining < 0 {
			remaining = 0
		}
		remainingQuota = &remaining
	}
	
	// Check if promo is available
	isAvailable := promo.IsActive && 
		!isExpired && 
		now.After(promo.StartDate) &&
		(promo.Quota == nil || promo.UsedCount < *promo.Quota)

	response := &entity.PromoResponse{
		ID:             promo.ID,
		CompanyID:      promo.CompanyID,
		CompanyName:    promo.Company.Name,
		BranchID:       promo.BranchID,
		Name:           promo.Name,
		Code:           promo.Code,
		Type:           promo.Type,
		Value:          promo.Value,
		MaxDiscount:    promo.MaxDiscount,
		MinTransaction: promo.MinTransaction,
		Quota:          promo.Quota,
		UsedCount:      promo.UsedCount,
		RemainingQuota: remainingQuota,
		StartDate:      promo.StartDate.Format("2006-01-02"),
		EndDate:        promo.EndDate.Format("2006-01-02"),
		IsActive:       promo.IsActive,
		IsExpired:      isExpired,
		IsAvailable:    isAvailable,
		CreatedAt:      promo.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      promo.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Add branch name if exists
	if promo.Branch != nil {
		response.BranchName = &promo.Branch.Name
	}

	return response
}
