package service

import (
	"errors"
	"fmt"
	"project-name/internal/entity"
	"project-name/internal/repository"
	"project-name/pkg"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(req entity.CreateOrderRequest, companyID, branchID uuid.UUID) (*entity.OrderResponse, error)
	CreatePublicOrder(req entity.CreatePublicOrderRequest) (*entity.OrderResponse, error)
	UpdateOrder(id uuid.UUID, req entity.UpdateOrderRequest, companyID, branchID uuid.UUID) (*entity.OrderResponse, error)
	DeleteOrder(id uuid.UUID, companyID, branchID uuid.UUID) error
	GetOrderByID(id uuid.UUID, companyID, branchID *uuid.UUID) (*entity.OrderResponse, error)
	GetAllOrders(companyID, branchID *uuid.UUID, status, method, customer, orderID string, pagination pkg.PaginationParams) ([]entity.OrderResponse, *pkg.PaginationMeta, error)
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	branchRepo  repository.BranchRepository
	taxRepo     repository.TaxRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, branchRepo repository.BranchRepository, taxRepo repository.TaxRepository) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		branchRepo:  branchRepo,
		taxRepo:     taxRepo,
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

	// Calculate taxes
	taxAmount, _, err := s.calculateTaxes(subtotalAmount, companyID, branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate taxes: %v", err)
	}

	totalAmount := subtotalAmount + taxAmount

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
		PromoCode:      req.PromoCode,
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

	// Calculate taxes
	taxAmount, _, err := s.calculateTaxes(subtotalAmount, req.CompanyID, req.BranchID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate taxes: %v", err)
	}

	totalAmount := subtotalAmount + taxAmount

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
		PromoCode:      req.PromoCode,
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

		// Calculate taxes
		taxAmount, _, err := s.calculateTaxes(subtotalAmount, companyID, branchID)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate taxes: %v", err)
		}

		order.SubtotalAmount = subtotalAmount
		order.TaxAmount = taxAmount
		order.TotalAmount = subtotalAmount + taxAmount

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
	_, taxDetails, _ := s.calculateTaxes(order.SubtotalAmount, order.CompanyID, order.BranchID)

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
		Status:         order.Status,
		SubtotalAmount: order.SubtotalAmount,
		TaxAmount:      order.TaxAmount,
		TotalAmount:    order.TotalAmount,
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
