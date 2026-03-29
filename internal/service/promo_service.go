package service

import (
	"errors"
	"fmt"
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
	ValidatePromoCode(code string, companyID uuid.UUID, branchID *uuid.UUID) (*entity.PromoValidationResponse, error)
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

	// Validate promo category
	if req.PromoCategory != "normal" && req.PromoCategory != "product" && req.PromoCategory != "bundle" {
		return nil, errors.New("promo_category must be 'normal', 'product', or 'bundle'")
	}

	// Validate product promo
	if req.PromoCategory == "product" && len(req.ProductIDs) == 0 {
		return nil, errors.New("product_ids required for product promo")
	}

	// Validate bundle promo
	if req.PromoCategory == "bundle" {
		if len(req.BundleItems) < 2 {
			return nil, errors.New("bundle promo requires at least 2 products")
		}
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
		PromoCategory:  req.PromoCategory,
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

	// Create promo products if product category
	if req.PromoCategory == "product" && len(req.ProductIDs) > 0 {
		if err := s.promoRepo.CreatePromoProducts(promo.ID, req.ProductIDs); err != nil {
			return nil, err
		}
	}

	// Create promo bundles if bundle category
	if req.PromoCategory == "bundle" && len(req.BundleItems) > 0 {
		bundles := make([]entity.PromoBundle, len(req.BundleItems))
		for i, item := range req.BundleItems {
			bundles[i] = entity.PromoBundle{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			}
		}
		if err := s.promoRepo.CreatePromoBundles(promo.ID, bundles); err != nil {
			return nil, err
		}
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
	if req.PromoCategory != "" {
		// Validate promo category
		if req.PromoCategory != "normal" && req.PromoCategory != "product" && req.PromoCategory != "bundle" {
			return nil, errors.New("promo_category must be 'normal', 'product', or 'bundle'")
		}
		promo.PromoCategory = req.PromoCategory
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

	// Update promo products if provided
	if req.PromoCategory == "product" || (promo.PromoCategory == "product" && len(req.ProductIDs) > 0) {
		// Delete existing products
		if err := s.promoRepo.DeletePromoProducts(promo.ID); err != nil {
			return nil, err
		}
		// Create new products
		if len(req.ProductIDs) > 0 {
			if err := s.promoRepo.CreatePromoProducts(promo.ID, req.ProductIDs); err != nil {
				return nil, err
			}
		}
	}

	// Update promo bundles if provided
	if req.PromoCategory == "bundle" || (promo.PromoCategory == "bundle" && len(req.BundleItems) > 0) {
		// Delete existing bundles
		if err := s.promoRepo.DeletePromoBundles(promo.ID); err != nil {
			return nil, err
		}
		// Create new bundles
		if len(req.BundleItems) > 0 {
			bundles := make([]entity.PromoBundle, len(req.BundleItems))
			for i, item := range req.BundleItems {
				bundles[i] = entity.PromoBundle{
					ProductID: item.ProductID,
					Quantity:  item.Quantity,
				}
			}
			if err := s.promoRepo.CreatePromoBundles(promo.ID, bundles); err != nil {
				return nil, err
			}
		}
	}

	if err := s.promoRepo.Update(promo); err != nil {
		return nil, err
	}

	// Reload with relations
	reloadedPromo, err := s.promoRepo.FindByID(promo.ID, companyID, branchID)
	if err != nil {
		return s.toResponse(promo), nil
	}

	return s.toResponse(reloadedPromo), nil
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
		PromoCategory:  promo.PromoCategory,
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

	// Add products for product promo
	if promo.PromoCategory == "product" && len(promo.PromoProducts) > 0 {
		response.Products = make([]entity.PromoProductResponse, len(promo.PromoProducts))
		for i, pp := range promo.PromoProducts {
			response.Products[i] = entity.PromoProductResponse{
				ProductID:   pp.ProductID,
				ProductName: pp.Product.Name,
			}
		}
	}

	// Add bundle items for bundle promo
	if promo.PromoCategory == "bundle" && len(promo.PromoBundles) > 0 {
		response.BundleItems = make([]entity.PromoBundleResponse, len(promo.PromoBundles))
		for i, pb := range promo.PromoBundles {
			response.BundleItems[i] = entity.PromoBundleResponse{
				ProductID:   pb.ProductID,
				ProductName: pb.Product.Name,
				Quantity:    pb.Quantity,
			}
		}
	}

	return response
}


func (s *promoService) ValidatePromoCode(code string, companyID uuid.UUID, branchID *uuid.UUID) (*entity.PromoValidationResponse, error) {
	promo, err := s.promoRepo.FindByCode(code, companyID, branchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entity.PromoValidationResponse{
				Valid:   false,
				Message: "Promo code not found",
			}, nil
		}
		return nil, err
	}

	now := time.Now()
	
	// Check if promo is active
	if !promo.IsActive {
		return &entity.PromoValidationResponse{
			Valid:   false,
			Message: "Promo is not active",
			Promo:   s.toResponse(promo),
		}, nil
	}

	// Check if promo has started
	if now.Before(promo.StartDate) {
		return &entity.PromoValidationResponse{
			Valid:   false,
			Message: fmt.Sprintf("Promo will start on %s", promo.StartDate.Format("2006-01-02")),
			Promo:   s.toResponse(promo),
		}, nil
	}

	// Check if promo has expired
	if now.After(promo.EndDate) {
		return &entity.PromoValidationResponse{
			Valid:   false,
			Message: fmt.Sprintf("Promo expired on %s", promo.EndDate.Format("2006-01-02")),
			Promo:   s.toResponse(promo),
		}, nil
	}

	// Check quota
	if promo.Quota != nil && promo.UsedCount >= *promo.Quota {
		return &entity.PromoValidationResponse{
			Valid:   false,
			Message: "Promo quota has been exhausted",
			Promo:   s.toResponse(promo),
		}, nil
	}

	// Promo is valid
	return &entity.PromoValidationResponse{
		Valid:   true,
		Message: "Promo is valid and can be used",
		Promo:   s.toResponse(promo),
	}, nil
}
