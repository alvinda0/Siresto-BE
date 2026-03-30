package service

import (
	"errors"
	"fmt"
	"project-name/internal/entity"
	"project-name/internal/repository"
	"project-name/pkg"
	"strings"
	"time"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(req entity.CreateOrderRequest, companyID, branchID uuid.UUID) (*entity.OrderResponse, error)
	QuickCreateOrder(req entity.QuickOrderRequest, companyID, branchID uuid.UUID) (*entity.OrderResponse, error)
	AddOrderItem(orderID uuid.UUID, req entity.AddOrderItemRequest, companyID, branchID uuid.UUID) (*entity.OrderResponse, error)
	CreatePublicOrder(req entity.CreatePublicOrderRequest) (*entity.OrderResponse, error)
	UpdateOrder(id uuid.UUID, req entity.UpdateOrderRequest, companyID, branchID uuid.UUID) (*entity.OrderResponse, error)
	UpdateOrderStatus(id uuid.UUID, req entity.UpdateOrderStatusRequest, companyID, branchID uuid.UUID) (*entity.OrderResponse, error)
	ProcessPayment(orderID uuid.UUID, req entity.ProcessPaymentRequest, companyID, branchID uuid.UUID) (*entity.PaymentResponse, error)
	DeleteOrder(id uuid.UUID, companyID, branchID uuid.UUID) error
	GetOrderByID(id uuid.UUID, companyID, branchID *uuid.UUID) (*entity.OrderResponse, error)
	GetAllOrders(companyID, branchID *uuid.UUID, status, method, customer, orderID string, pagination pkg.PaginationParams) ([]entity.OrderResponse, *pkg.PaginationMeta, error)
	GetTransactionReport(companyID, branchID uuid.UUID, filter entity.TransactionReportFilter) ([]entity.TransactionReportDTO, *pkg.PaginationMeta, error)
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	branchRepo  repository.BranchRepository
	taxRepo     repository.TaxRepository
	promoRepo   repository.PromoRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, branchRepo repository.BranchRepository, taxRepo repository.TaxRepository, promoRepo repository.PromoRepository) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		branchRepo:  branchRepo,
		taxRepo:     taxRepo,
		promoRepo:   promoRepo,
	}
}

func (s *orderService) CreateOrder(req entity.CreateOrderRequest, companyID, branchID uuid.UUID) (*entity.OrderResponse, error) {
	// Validate branch
	branch, err := s.branchRepo.FindByID(branchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	if branch.CompanyID != companyID {
		return nil, errors.New("branch does not belong to your company")
	}

	// Calculate subtotal and validate products
	var subtotalAmount float64
	var orderItems []entity.OrderItem

	for _, item := range req.OrderItems {
		product, err := s.productRepo.FindByID(item.ProductID, companyID, branchID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found", item.ProductID)
		}

		if product.BranchID != branchID {
			return nil, fmt.Errorf("product %s does not belong to this branch", product.Name)
		}

		if !product.IsAvailable {
			return nil, fmt.Errorf("product %s is not available", product.Name)
		}

		subtotal := product.Price * float64(item.Quantity)
		subtotalAmount += subtotal

		orderItems = append(orderItems, entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Note:      item.Note,
		})
	}

	// Apply promo if provided
	var discountAmount float64
	var promoID *uuid.UUID
	if req.PromoCode != "" {
		discount, pID, err := s.applyPromo(req.PromoCode, subtotalAmount, companyID, branchID)
		if err != nil {
			return nil, fmt.Errorf("promo error: %v", err)
		}
		discountAmount = discount
		promoID = pID
	}

	// Calculate taxes based on (subtotal - discount)
	amountAfterDiscount := subtotalAmount - discountAmount
	taxAmount, _, err := s.calculateTaxes(amountAfterDiscount, companyID, branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate taxes: %v", err)
	}

	totalAmount := amountAfterDiscount + taxAmount

	// Create order
	order := &entity.Order{
		CompanyID:      companyID,
		BranchID:       branchID,
		CustomerName:   req.CustomerName,
		CustomerPhone:  req.CustomerPhone,
		TableNumber:    req.TableNumber,
		Notes:          req.Notes,
		ReferralCode:   req.ReferralCode,
		OrderMethod:    req.OrderMethod,
		PromoID:        promoID,
		PromoCode:      req.PromoCode,
		DiscountAmount: discountAmount,
		Status:         entity.OrderStatusPending,
		SubtotalAmount: subtotalAmount,
		TaxAmount:      taxAmount,
		TotalAmount:    totalAmount,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Create order items
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}

	if err := s.orderRepo.CreateOrderItems(orderItems); err != nil {
		return nil, err
	}

	// Fetch complete order
	return s.GetOrderByID(order.ID, &companyID, &branchID)
}

func (s *orderService) QuickCreateOrder(req entity.QuickOrderRequest, companyID, branchID uuid.UUID) (*entity.OrderResponse, error) {
	// Validate branch
	branch, err := s.branchRepo.FindByID(branchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	if branch.CompanyID != companyID {
		return nil, errors.New("branch does not belong to your company")
	}

	// Calculate subtotal and validate products
	var subtotalAmount float64
	var orderItems []entity.OrderItem

	for _, item := range req.OrderItems {
		product, err := s.productRepo.FindByID(item.ProductID, companyID, branchID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found", item.ProductID)
		}

		if product.BranchID != branchID {
			return nil, fmt.Errorf("product %s does not belong to this branch", product.Name)
		}

		if !product.IsAvailable {
			return nil, fmt.Errorf("product %s is not available", product.Name)
		}

		subtotal := product.Price * float64(item.Quantity)
		subtotalAmount += subtotal

		orderItems = append(orderItems, entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Note:      item.Note,
		})
	}

	// Calculate taxes (no promo for quick order)
	taxAmount, _, err := s.calculateTaxes(subtotalAmount, companyID, branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate taxes: %v", err)
	}

	totalAmount := subtotalAmount + taxAmount

	// Create order
	order := &entity.Order{
		CompanyID:      companyID,
		BranchID:       branchID,
		TableNumber:    req.TableNumber,
		OrderMethod:    req.OrderMethod,
		Status:         entity.OrderStatusPending,
		SubtotalAmount: subtotalAmount,
		TaxAmount:      taxAmount,
		TotalAmount:    totalAmount,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Create order items
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}

	if err := s.orderRepo.CreateOrderItems(orderItems); err != nil {
		return nil, err
	}

	// Fetch complete order
	return s.GetOrderByID(order.ID, &companyID, &branchID)
}

func (s *orderService) AddOrderItem(orderID uuid.UUID, req entity.AddOrderItemRequest, companyID, branchID uuid.UUID) (*entity.OrderResponse, error) {
	// Find existing order
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.CompanyID != companyID || order.BranchID != branchID {
		return nil, errors.New("unauthorized to modify this order")
	}

	// Validate product
	product, err := s.productRepo.FindByID(req.ProductID, companyID, branchID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	if product.BranchID != branchID {
		return nil, fmt.Errorf("product does not belong to this branch")
	}

	if !product.IsAvailable {
		return nil, fmt.Errorf("product %s is not available", product.Name)
	}

	// Check if item with same product_id already exists
	var existingItem *entity.OrderItem
	for i := range order.OrderItems {
		if order.OrderItems[i].ProductID == req.ProductID {
			existingItem = &order.OrderItems[i]
			break
		}
	}

	if existingItem != nil {
		// Update existing item quantity
		existingItem.Quantity += req.Quantity
		if req.Note != "" {
			existingItem.Note = req.Note
		}
		if err := s.orderRepo.UpdateOrderItem(existingItem); err != nil {
			return nil, err
		}
	} else {
		// Create new order item
		newItem := entity.OrderItem{
			OrderID:   orderID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Price:     product.Price,
			Note:      req.Note,
		}
		if err := s.orderRepo.CreateOrderItems([]entity.OrderItem{newItem}); err != nil {
			return nil, err
		}
	}

	// Recalculate subtotal from ALL items (not just increment)
	// Fetch fresh order with all items
	order, err = s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("failed to fetch updated order")
	}

	var newSubtotal float64
	for _, item := range order.OrderItems {
		newSubtotal += item.Price * float64(item.Quantity)
	}

	// Recalculate discount if promo exists (without incrementing usage count)
	var discountAmount float64
	if order.PromoCode != "" && order.PromoID != nil {
		promo, err := s.promoRepo.FindByIDSimple(*order.PromoID)
		if err == nil && promo.IsActive {
			// Recalculate discount based on new subtotal
			if promo.Type == "percentage" {
				discountAmount = newSubtotal * (promo.Value / 100)
				if promo.MaxDiscount != nil && discountAmount > *promo.MaxDiscount {
					discountAmount = *promo.MaxDiscount
				}
			} else if promo.Type == "fixed" {
				discountAmount = promo.Value
				if discountAmount > newSubtotal {
					discountAmount = newSubtotal
				}
			}
		} else {
			// Promo no longer valid, remove it
			order.PromoCode = ""
			order.PromoID = nil
		}
	}

	// Recalculate taxes
	amountAfterDiscount := newSubtotal - discountAmount
	taxAmount, _, err := s.calculateTaxes(amountAfterDiscount, companyID, branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate taxes: %v", err)
	}

	order.SubtotalAmount = newSubtotal
	order.DiscountAmount = discountAmount
	order.TaxAmount = taxAmount
	order.TotalAmount = amountAfterDiscount + taxAmount

	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return s.GetOrderByID(order.ID, &companyID, &branchID)
}

func (s *orderService) CreatePublicOrder(req entity.CreatePublicOrderRequest) (*entity.OrderResponse, error) {
	// Validate branch
	branch, err := s.branchRepo.FindByID(req.BranchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	if branch.CompanyID != req.CompanyID {
		return nil, errors.New("branch does not belong to the specified company")
	}

	// Calculate subtotal and validate products
	var subtotalAmount float64
	var orderItems []entity.OrderItem

	for _, item := range req.OrderItems {
		product, err := s.productRepo.FindByID(item.ProductID, req.CompanyID, req.BranchID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found", item.ProductID)
		}

		if product.BranchID != req.BranchID {
			return nil, fmt.Errorf("product %s does not belong to this branch", product.Name)
		}

		if !product.IsAvailable {
			return nil, fmt.Errorf("product %s is not available", product.Name)
		}

		subtotal := product.Price * float64(item.Quantity)
		subtotalAmount += subtotal

		orderItems = append(orderItems, entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Note:      item.Note,
		})
	}

	// Apply promo if provided
	var discountAmount float64
	var promoID *uuid.UUID
	if req.PromoCode != "" {
		discount, pID, err := s.applyPromo(req.PromoCode, subtotalAmount, req.CompanyID, req.BranchID)
		if err != nil {
			return nil, fmt.Errorf("promo error: %v", err)
		}
		discountAmount = discount
		promoID = pID
	}

	// Calculate taxes based on (subtotal - discount)
	amountAfterDiscount := subtotalAmount - discountAmount
	taxAmount, _, err := s.calculateTaxes(amountAfterDiscount, req.CompanyID, req.BranchID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate taxes: %v", err)
	}

	totalAmount := amountAfterDiscount + taxAmount

	// Create order
	order := &entity.Order{
		CompanyID:      req.CompanyID,
		BranchID:       req.BranchID,
		CustomerName:   req.CustomerName,
		CustomerPhone:  req.CustomerPhone,
		TableNumber:    req.TableNumber,
		Notes:          req.Notes,
		ReferralCode:   req.ReferralCode,
		OrderMethod:    req.OrderMethod,
		PromoID:        promoID,
		PromoCode:      req.PromoCode,
		DiscountAmount: discountAmount,
		Status:         entity.OrderStatusPending,
		SubtotalAmount: subtotalAmount,
		TaxAmount:      taxAmount,
		TotalAmount:    totalAmount,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Create order items
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}

	if err := s.orderRepo.CreateOrderItems(orderItems); err != nil {
		return nil, err
	}

	// Fetch complete order
	return s.GetOrderByID(order.ID, &req.CompanyID, &req.BranchID)
}

func (s *orderService) UpdateOrder(id uuid.UUID, req entity.UpdateOrderRequest, companyID, branchID uuid.UUID) (*entity.OrderResponse, error) {
	// Find existing order
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.CompanyID != companyID || order.BranchID != branchID {
		return nil, errors.New("unauthorized to update this order")
	}

	// Update order fields
	if req.CustomerName != "" {
		order.CustomerName = req.CustomerName
	}
	if req.CustomerPhone != "" {
		order.CustomerPhone = req.CustomerPhone
	}
	if req.TableNumber != "" {
		order.TableNumber = req.TableNumber
	}
	if req.Notes != "" {
		order.Notes = req.Notes
	}
	if req.OrderMethod != "" {
		order.OrderMethod = req.OrderMethod
	}
	if req.Status != "" {
		order.Status = req.Status
	}

	// Update order items if provided
	if len(req.OrderItems) > 0 {
		// Delete existing items
		if err := s.orderRepo.DeleteOrderItems(order.ID); err != nil {
			return nil, err
		}

		// Calculate new subtotal
		var subtotalAmount float64
		var orderItems []entity.OrderItem

		for _, item := range req.OrderItems {
			product, err := s.productRepo.FindByID(item.ProductID, companyID, branchID)
			if err != nil {
				return nil, fmt.Errorf("product %s not found", item.ProductID)
			}

			if product.BranchID != branchID {
				return nil, fmt.Errorf("product %s does not belong to this branch", product.Name)
			}

			subtotal := product.Price * float64(item.Quantity)
			subtotalAmount += subtotal

			orderItems = append(orderItems, entity.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Price,
				Note:      item.Note,
			})
		}

		// Apply promo if exists
		var discountAmount float64
		if order.PromoCode != "" {
			discount, _, err := s.applyPromo(order.PromoCode, subtotalAmount, companyID, branchID)
			if err != nil {
				// If promo is no longer valid, remove it
				order.PromoCode = ""
				order.PromoID = nil
				discountAmount = 0
			} else {
				discountAmount = discount
			}
		}

		// Calculate taxes based on (subtotal - discount)
		amountAfterDiscount := subtotalAmount - discountAmount
		taxAmount, _, err := s.calculateTaxes(amountAfterDiscount, companyID, branchID)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate taxes: %v", err)
		}

		order.SubtotalAmount = subtotalAmount
		order.DiscountAmount = discountAmount
		order.TaxAmount = taxAmount
		order.TotalAmount = amountAfterDiscount + taxAmount

		// Create new items
		if err := s.orderRepo.CreateOrderItems(orderItems); err != nil {
			return nil, err
		}
	}

	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return s.GetOrderByID(order.ID, &companyID, &branchID)
}

func (s *orderService) UpdateOrderStatus(id uuid.UUID, req entity.UpdateOrderStatusRequest, companyID, branchID uuid.UUID) (*entity.OrderResponse, error) {
	// Find existing order
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.CompanyID != companyID || order.BranchID != branchID {
		return nil, errors.New("unauthorized to update this order")
	}

	// Validate status transition (optional - add business rules here)
	// For example: PENDING -> PREPARING -> READY -> COMPLETED
	// You can add validation logic here if needed

	// Update status
	if err := s.orderRepo.UpdateStatus(id, req.Status); err != nil {
		return nil, err
	}

	return s.GetOrderByID(id, &companyID, &branchID)
}

func (s *orderService) DeleteOrder(id uuid.UUID, companyID, branchID uuid.UUID) error {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	if order.CompanyID != companyID || order.BranchID != branchID {
		return errors.New("unauthorized to delete this order")
	}

	return s.orderRepo.Delete(id)
}

func (s *orderService) GetOrderByID(id uuid.UUID, companyID, branchID *uuid.UUID) (*entity.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Check access control
	if companyID != nil && order.CompanyID != *companyID {
		return nil, errors.New("unauthorized to view this order")
	}
	if branchID != nil && order.BranchID != *branchID {
		return nil, errors.New("unauthorized to view this order")
	}

	return s.toOrderResponse(order), nil
}

func (s *orderService) GetAllOrders(companyID, branchID *uuid.UUID, status, method, customer, orderID string, pagination pkg.PaginationParams) ([]entity.OrderResponse, *pkg.PaginationMeta, error) {
	orders, total, err := s.orderRepo.FindAll(companyID, branchID, status, method, customer, orderID, pagination)
	if err != nil {
		return nil, nil, err
	}

	var responses []entity.OrderResponse
	for _, order := range orders {
		responses = append(responses, *s.toOrderResponse(&order))
	}

	meta := pagination.CreateMeta(total)

	return responses, &meta, nil
}

func (s *orderService) toOrderResponse(order *entity.Order) *entity.OrderResponse {
	var items []entity.OrderItemDTO
	for _, item := range order.OrderItems {
		productName := ""
		if item.Product.ID != uuid.Nil {
			productName = item.Product.Name
		}

		items = append(items, entity.OrderItemDTO{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Price * float64(item.Quantity),
			Note:        item.Note,
		})
	}

	// Calculate tax details for response
	amountAfterDiscount := order.SubtotalAmount - order.DiscountAmount
	_, taxDetails, _ := s.calculateTaxes(amountAfterDiscount, order.CompanyID, order.BranchID)

	// Get promo details if promo was used
	var promoDetailsList []entity.PromoDetailDTO
	if order.PromoCode != "" {
		promoCodes := strings.Split(order.PromoCode, ",")
		for _, code := range promoCodes {
			code = strings.TrimSpace(code)
			if code != "" {
				// Fetch promo details by code
				promo, err := s.promoRepo.FindByCode(code, order.CompanyID, &order.BranchID)
				if err == nil {
					promoDetailsList = append(promoDetailsList, entity.PromoDetailDTO{
						PromoID:        promo.ID,
						PromoName:      promo.Name,
						PromoCode:      promo.Code,
						PromoType:      promo.Type,
						PromoValue:     promo.Value,
						DiscountAmount: 0, // We'll calculate individual discount later if needed
						MaxDiscount:    promo.MaxDiscount,
						MinTransaction: promo.MinTransaction,
					})
				}
			}
		}
	}

	paidAtStr := ""
	if order.PaidAt != nil {
		paidAtStr = order.PaidAt.Format("2006-01-02 15:04:05")
	}

	return &entity.OrderResponse{
		ID:             order.ID,
		CompanyID:      order.CompanyID,
		BranchID:       order.BranchID,
		CustomerName:   order.CustomerName,
		CustomerPhone:  order.CustomerPhone,
		TableNumber:    order.TableNumber,
		Notes:          order.Notes,
		ReferralCode:   order.ReferralCode,
		OrderMethod:    order.OrderMethod,
		PromoCode:      order.PromoCode,
		PromoID:        order.PromoID,
		DiscountAmount: order.DiscountAmount,
		PromoDetails:   promoDetailsList,
		Status:         order.Status,
		SubtotalAmount: order.SubtotalAmount,
		TaxAmount:      order.TaxAmount,
		TotalAmount:    order.TotalAmount,
		PaymentMethod:  order.PaymentMethod,
		PaymentStatus:  order.PaymentStatus,
		PaidAmount:     order.PaidAmount,
		ChangeAmount:   order.ChangeAmount,
		PaymentNote:    order.PaymentNote,
		PaidAt:         paidAtStr,
		TaxDetails:     taxDetails,
		OrderItems:     items,
		CreatedAt:      order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// calculateTaxes menghitung pajak bertingkat berdasarkan prioritas
// Priority 1 = dihitung pertama, Priority 2 = dihitung kedua, dst.
// Contoh: subtotal 100.000, Service Charge 5% (prioritas 1), PB1 10% (prioritas 2)
// - Base: 100.000
// - Service Charge (prioritas 1): 100.000 * 5% = 5.000 -> Total: 105.000
// - PB1 (prioritas 2): 105.000 * 10% = 10.500 -> Total: 115.500
func (s *orderService) calculateTaxes(subtotal float64, companyID, branchID uuid.UUID) (float64, []entity.TaxDetailDTO, error) {
	// Get active taxes ordered by priority ASC (1, 2, 3, ...)
	taxes, err := s.taxRepo.FindActiveTaxesByBranch(companyID, branchID)
	if err != nil {
		return 0, nil, err
	}

	var taxDetails []entity.TaxDetailDTO
	var totalTax float64
	currentAmount := subtotal

	// Calculate taxes in order of priority (1 first, then 2, then 3, ...)
	for _, tax := range taxes {
		taxAmount := currentAmount * (tax.Presentase / 100)
		
		taxDetails = append(taxDetails, entity.TaxDetailDTO{
			TaxID:      tax.ID,
			TaxName:    tax.NamaPajak,
			Percentage: tax.Presentase,
			Priority:   tax.Prioritas,
			BaseAmount: currentAmount,
			TaxAmount:  taxAmount,
		})

		totalTax += taxAmount
		currentAmount += taxAmount // Untuk pajak berikutnya, base-nya adalah amount + pajak sebelumnya
	}

	return totalTax, taxDetails, nil
}

// applyPromo validates and calculates discount from promo code
func (s *orderService) applyPromo(promoCode string, subtotal float64, companyID, branchID uuid.UUID) (float64, *uuid.UUID, error) {
	// Find promo by code
	promo, err := s.promoRepo.FindByCode(promoCode, companyID, &branchID)
	if err != nil {
		return 0, nil, errors.New("promo code not found")
	}

	// Validate promo is active
	if !promo.IsActive {
		return 0, nil, errors.New("promo is not active")
	}

	// Validate promo date range
	now := time.Now()
	if now.Before(promo.StartDate) {
		return 0, nil, errors.New("promo has not started yet")
	}
	if now.After(promo.EndDate) {
		return 0, nil, errors.New("promo has expired")
	}

	// Validate quota
	if promo.Quota != nil && promo.UsedCount >= *promo.Quota {
		return 0, nil, errors.New("promo quota has been exhausted")
	}

	// Validate minimum transaction
	if promo.MinTransaction != nil && subtotal < *promo.MinTransaction {
		return 0, nil, fmt.Errorf("minimum transaction is %.2f", *promo.MinTransaction)
	}

	// Calculate discount
	var discount float64
	if promo.Type == "percentage" {
		discount = subtotal * (promo.Value / 100)
		// Apply max discount if set
		if promo.MaxDiscount != nil && discount > *promo.MaxDiscount {
			discount = *promo.MaxDiscount
		}
	} else if promo.Type == "fixed" {
		discount = promo.Value
		// Discount cannot exceed subtotal
		if discount > subtotal {
			discount = subtotal
		}
	} else {
		return 0, nil, errors.New("invalid promo type")
	}

	// Increment used count
	promo.UsedCount++
	if err := s.promoRepo.Update(promo); err != nil {
		return 0, nil, fmt.Errorf("failed to update promo usage: %v", err)
	}

	return discount, &promo.ID, nil
}

func (s *orderService) ProcessPayment(orderID uuid.UUID, req entity.ProcessPaymentRequest, companyID, branchID uuid.UUID) (*entity.PaymentResponse, error) {
	// Find existing order
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.CompanyID != companyID || order.BranchID != branchID {
		return nil, errors.New("unauthorized to process payment for this order")
	}

	// Check if order is already paid
	if order.PaymentStatus == entity.PaymentStatusPaid {
		return nil, errors.New("order has already been paid")
	}

	// Apply promo if provided
	if req.PromoCode != "" {
		// Check if trying to apply the same promo code again
		existingPromoCodes := []string{}
		if order.PromoCode != "" {
			existingPromoCodes = strings.Split(order.PromoCode, ",")
		}
		
		// Check if promo already exists
		for _, code := range existingPromoCodes {
			if strings.TrimSpace(code) == req.PromoCode {
				return nil, errors.New("promo code has already been applied to this order")
			}
		}
		
		// Apply additional promo
		additionalDiscount, promoID, err := s.applyPromo(req.PromoCode, order.SubtotalAmount, companyID, branchID)
		if err != nil {
			return nil, fmt.Errorf("promo error: %v", err)
		}
		
		// Add to existing promo codes
		if order.PromoCode != "" {
			order.PromoCode = order.PromoCode + "," + req.PromoCode
		} else {
			order.PromoCode = req.PromoCode
		}
		
		// Add to existing discount
		order.DiscountAmount += additionalDiscount
		
		// Store the new promo ID (we'll keep the last one for backward compatibility)
		order.PromoID = promoID
		
		// Recalculate taxes based on (subtotal - total discount)
		amountAfterDiscount := order.SubtotalAmount - order.DiscountAmount
		taxAmount, _, err := s.calculateTaxes(amountAfterDiscount, companyID, branchID)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate taxes: %v", err)
		}
		
		order.TaxAmount = taxAmount
		order.TotalAmount = amountAfterDiscount + taxAmount
	}

	// Validate payment method
	validMethods := []entity.PaymentMethod{
		entity.PaymentMethodQRIS,
		entity.PaymentMethodCash,
		entity.PaymentMethodDebit,
		entity.PaymentMethodCredit,
		entity.PaymentMethodGoPay,
		entity.PaymentMethodOVO,
		entity.PaymentMethodComplimentary,
	}
	
	isValidMethod := false
	for _, method := range validMethods {
		if req.PaymentMethod == method {
			isValidMethod = true
			break
		}
	}
	
	if !isValidMethod {
		return nil, errors.New("invalid payment method")
	}

	// Calculate change (only for TUNAI)
	var changeAmount float64
	
	// For COMPLIMENTARY, make everything free
	if req.PaymentMethod == entity.PaymentMethodComplimentary {
		// Set discount to cover full subtotal (making it free)
		order.DiscountAmount = order.SubtotalAmount
		order.TaxAmount = 0
		order.TotalAmount = 0
		req.PaidAmount = 0
		changeAmount = 0
	} else if req.PaymentMethod == entity.PaymentMethodCash {
		if req.PaidAmount < order.TotalAmount {
			return nil, errors.New("paid amount is less than total amount")
		}
		changeAmount = req.PaidAmount - order.TotalAmount
	} else {
		// For non-cash payments, paid amount must match total amount
		if req.PaidAmount != order.TotalAmount {
			return nil, errors.New("paid amount must match total amount for non-cash payments")
		}
	}

	// Update order with payment info
	now := time.Now()
	order.PaymentMethod = req.PaymentMethod
	order.PaymentStatus = entity.PaymentStatusPaid
	order.PaidAmount = req.PaidAmount
	order.ChangeAmount = changeAmount
	order.PaymentNote = req.PaymentNote
	order.PaidAt = &now
	order.Status = entity.OrderStatusCompleted // Auto complete when paid

	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	// Prepare response
	_, taxDetails, _ := s.calculateTaxes(order.SubtotalAmount-order.DiscountAmount, companyID, branchID)
	
	// Get all promo details
	var promoDetailsList []entity.PromoDetailDTO
	if order.PromoCode != "" {
		promoCodes := strings.Split(order.PromoCode, ",")
		for _, code := range promoCodes {
			code = strings.TrimSpace(code)
			if code != "" {
				promo, err := s.promoRepo.FindByCode(code, companyID, &branchID)
				if err == nil {
					promoDetailsList = append(promoDetailsList, entity.PromoDetailDTO{
						PromoID:        promo.ID,
						PromoName:      promo.Name,
						PromoCode:      promo.Code,
						PromoType:      promo.Type,
						PromoValue:     promo.Value,
						DiscountAmount: 0, // Total discount shown in main field
						MaxDiscount:    promo.MaxDiscount,
						MinTransaction: promo.MinTransaction,
					})
				}
			}
		}
	}

	paidAtStr := ""
	if order.PaidAt != nil {
		paidAtStr = order.PaidAt.Format("2006-01-02 15:04:05")
	}

	return &entity.PaymentResponse{
		OrderID:        order.ID,
		PaymentMethod:  order.PaymentMethod,
		PaymentStatus:  order.PaymentStatus,
		SubtotalAmount: order.SubtotalAmount,
		DiscountAmount: order.DiscountAmount,
		TaxAmount:      order.TaxAmount,
		TotalAmount:    order.TotalAmount,
		PaidAmount:     order.PaidAmount,
		ChangeAmount:   order.ChangeAmount,
		PaymentNote:    order.PaymentNote,
		PaidAt:         paidAtStr,
		PromoDetails:   promoDetailsList,
		TaxDetails:     taxDetails,
	}, nil
}

func (s *orderService) GetTransactionReport(companyID, branchID uuid.UUID, filter entity.TransactionReportFilter) ([]entity.TransactionReportDTO, *pkg.PaginationMeta, error) {
	// Get orders from repository
	orders, total, err := s.orderRepo.GetTransactionReport(companyID, branchID, filter)
	if err != nil {
		return nil, nil, err
	}

	// Convert to DTOs
	reportDTOs := make([]entity.TransactionReportDTO, len(orders))
	for i, order := range orders {
		reportDTOs[i] = order.ToReportDTO()
	}

	// Create pagination meta
	meta := &pkg.PaginationMeta{
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalItems: total,
		TotalPages: int((total + int64(filter.Limit) - 1) / int64(filter.Limit)),
	}

	return reportDTOs, meta, nil
}
