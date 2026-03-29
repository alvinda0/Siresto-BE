package repository

import (
	"project-name/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PromoRepository interface {
	Create(promo *entity.Promo) error
	Update(promo *entity.Promo) error
	Delete(id, companyID uuid.UUID, branchID *uuid.UUID) error
	FindByID(id, companyID uuid.UUID, branchID *uuid.UUID) (*entity.Promo, error)
	FindByCode(code string, companyID uuid.UUID, branchID *uuid.UUID) (*entity.Promo, error)
	FindByCompany(companyID uuid.UUID, page, limit int) ([]entity.Promo, int64, error)
	FindByBranch(companyID, branchID uuid.UUID, page, limit int) ([]entity.Promo, int64, error)
	CreatePromoProducts(promoID uuid.UUID, productIDs []uuid.UUID) error
	CreatePromoBundles(promoID uuid.UUID, bundles []entity.PromoBundle) error
	DeletePromoProducts(promoID uuid.UUID) error
	DeletePromoBundles(promoID uuid.UUID) error
}

type promoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) PromoRepository {
	return &promoRepository{db: db}
}

func (r *promoRepository) Create(promo *entity.Promo) error {
	return r.db.Create(promo).Error
}

func (r *promoRepository) Update(promo *entity.Promo) error {
	return r.db.Save(promo).Error
}

func (r *promoRepository) Delete(id, companyID uuid.UUID, branchID *uuid.UUID) error {
	query := r.db.Where("id = ? AND company_id = ?", id, companyID)
	
	if branchID != nil {
		// User with branch can only delete their own branch promos
		query = query.Where("branch_id = ?", *branchID)
	} else {
		// User without branch (OWNER) can only delete company-level promos
		query = query.Where("branch_id IS NULL")
	}
	
	return query.Delete(&entity.Promo{}).Error
}

func (r *promoRepository) FindByID(id, companyID uuid.UUID, branchID *uuid.UUID) (*entity.Promo, error) {
	var promo entity.Promo
	query := r.db.Preload("Company").Preload("Branch").
		Where("id = ? AND company_id = ?", id, companyID)
	
	// Users can view:
	// - Company-level promos (branch_id IS NULL) - visible to all
	// - Their own branch promos (branch_id = user's branch_id)
	if branchID != nil {
		query = query.Where("(branch_id IS NULL OR branch_id = ?)", *branchID)
	} else {
		// OWNER can view company-level promos only
		query = query.Where("branch_id IS NULL")
	}
	
	err := query.First(&promo).Error
	if err != nil {
		return nil, err
	}
	
	// Manually load PromoProducts and PromoBundles
	r.db.Preload("Product").Where("promo_id = ?", promo.ID).Find(&promo.PromoProducts)
	r.db.Preload("Product").Where("promo_id = ?", promo.ID).Find(&promo.PromoBundles)
	
	return &promo, nil
}

func (r *promoRepository) FindByCode(code string, companyID uuid.UUID, branchID *uuid.UUID) (*entity.Promo, error) {
	var promo entity.Promo
	query := r.db.Preload("Company").Preload("Branch").
		Where("code = ? AND company_id = ?", code, companyID)
	
	if branchID != nil {
		query = query.Where("(branch_id IS NULL OR branch_id = ?)", *branchID)
	} else {
		query = query.Where("branch_id IS NULL")
	}
	
	err := query.First(&promo).Error
	if err != nil {
		return nil, err
	}
	
	// Manually load PromoProducts and PromoBundles
	r.db.Preload("Product").Where("promo_id = ?", promo.ID).Find(&promo.PromoProducts)
	r.db.Preload("Product").Where("promo_id = ?", promo.ID).Find(&promo.PromoBundles)
	
	return &promo, nil
}

func (r *promoRepository) FindByCompany(companyID uuid.UUID, page, limit int) ([]entity.Promo, int64, error) {
	var promos []entity.Promo
	var total int64
	
	// Count total
	r.db.Model(&entity.Promo{}).
		Where("company_id = ? AND branch_id IS NULL", companyID).
		Count(&total)
	
	// Get paginated data
	offset := (page - 1) * limit
	err := r.db.Preload("Company").Preload("Branch").
		Where("company_id = ? AND branch_id IS NULL", companyID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&promos).Error
	
	// Manually load PromoProducts and PromoBundles
	if err == nil && len(promos) > 0 {
		promoIDs := make([]uuid.UUID, len(promos))
		for i, p := range promos {
			promoIDs[i] = p.ID
		}
		
		var promoProducts []entity.PromoProduct
		r.db.Preload("Product").Where("promo_id IN ?", promoIDs).Find(&promoProducts)
		
		promoProductsMap := make(map[uuid.UUID][]entity.PromoProduct)
		for _, pp := range promoProducts {
			promoProductsMap[pp.PromoID] = append(promoProductsMap[pp.PromoID], pp)
		}
		
		var promoBundles []entity.PromoBundle
		r.db.Preload("Product").Where("promo_id IN ?", promoIDs).Find(&promoBundles)
		
		promoBundlesMap := make(map[uuid.UUID][]entity.PromoBundle)
		for _, pb := range promoBundles {
			promoBundlesMap[pb.PromoID] = append(promoBundlesMap[pb.PromoID], pb)
		}
		
		for i := range promos {
			promos[i].PromoProducts = promoProductsMap[promos[i].ID]
			promos[i].PromoBundles = promoBundlesMap[promos[i].ID]
		}
	}
	
	return promos, total, err
}

func (r *promoRepository) FindByBranch(companyID, branchID uuid.UUID, page, limit int) ([]entity.Promo, int64, error) {
	var promos []entity.Promo
	var total int64
	
	// Count total (both company-level and branch-specific)
	r.db.Model(&entity.Promo{}).
		Where("company_id = ? AND (branch_id IS NULL OR branch_id = ?)", companyID, branchID).
		Count(&total)
	
	// Get paginated data
	offset := (page - 1) * limit
	err := r.db.Preload("Company").Preload("Branch").
		Where("company_id = ? AND (branch_id IS NULL OR branch_id = ?)", companyID, branchID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&promos).Error
	
	// Manually load PromoProducts and PromoBundles to avoid GORM preload issues
	if err == nil && len(promos) > 0 {
		promoIDs := make([]uuid.UUID, len(promos))
		for i, p := range promos {
			promoIDs[i] = p.ID
		}
		
		// Load PromoProducts
		var promoProducts []entity.PromoProduct
		r.db.Preload("Product").Where("promo_id IN ?", promoIDs).Find(&promoProducts)
		
		// Map to promos
		promoProductsMap := make(map[uuid.UUID][]entity.PromoProduct)
		for _, pp := range promoProducts {
			promoProductsMap[pp.PromoID] = append(promoProductsMap[pp.PromoID], pp)
		}
		
		// Load PromoBundles
		var promoBundles []entity.PromoBundle
		r.db.Preload("Product").Where("promo_id IN ?", promoIDs).Find(&promoBundles)
		
		// Map to promos
		promoBundlesMap := make(map[uuid.UUID][]entity.PromoBundle)
		for _, pb := range promoBundles {
			promoBundlesMap[pb.PromoID] = append(promoBundlesMap[pb.PromoID], pb)
		}
		
		// Assign to promos
		for i := range promos {
			promos[i].PromoProducts = promoProductsMap[promos[i].ID]
			promos[i].PromoBundles = promoBundlesMap[promos[i].ID]
		}
	}
	
	return promos, total, err
}


func (r *promoRepository) CreatePromoProducts(promoID uuid.UUID, productIDs []uuid.UUID) error {
	if len(productIDs) == 0 {
		return nil
	}

	promoProducts := make([]entity.PromoProduct, len(productIDs))
	for i, productID := range productIDs {
		promoProducts[i] = entity.PromoProduct{
			PromoID:   promoID,
			ProductID: productID,
		}
	}

	return r.db.Create(&promoProducts).Error
}

func (r *promoRepository) CreatePromoBundles(promoID uuid.UUID, bundles []entity.PromoBundle) error {
	if len(bundles) == 0 {
		return nil
	}

	for i := range bundles {
		bundles[i].PromoID = promoID
	}

	return r.db.Create(&bundles).Error
}

func (r *promoRepository) DeletePromoProducts(promoID uuid.UUID) error {
	return r.db.Where("promo_id = ?", promoID).Delete(&entity.PromoProduct{}).Error
}

func (r *promoRepository) DeletePromoBundles(promoID uuid.UUID) error {
	return r.db.Where("promo_id = ?", promoID).Delete(&entity.PromoBundle{}).Error
}
