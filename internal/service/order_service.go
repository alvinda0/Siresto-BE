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
	GetAllOrders(companyID, branchID *uuid.UUID, pagination pkg.PaginationParams) ([]entity.OrderResponse, *pkg.PaginationMeta, error)
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	branchRepo  repository.BranchRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, branchRepo repository.BranchRepository) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		branchRepo:  branchRepo,
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

	// Calculate total and validate products
	var totalAmount float64
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
		totalAmount += subtotal

		orderItems = append(orderItems, entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Note:      item.Note,
		})
	}

	// Create order
	order := &entity.Order{
		CompanyID:     companyID,
		BranchID:      branchID,
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		TableNumber:   req.TableNumber,
		Notes:         req.Notes,
		ReferralCode:  req.ReferralCode,
		OrderMethod:   req.OrderMethod,
		PromoCode:     req.PromoCode,
		Status:        entity.OrderStatusPending,
		TotalAmount:   totalAmount,
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

	// Calculate total and validate products
	var totalAmount float64
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
		totalAmount += subtotal

		orderItems = append(orderItems, entity.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Note:      item.Note,
		})
	}

	// Create order
	order := &entity.Order{
		CompanyID:     req.CompanyID,
		BranchID:      req.BranchID,
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		TableNumber:   req.TableNumber,
		Notes:         req.Notes,
		ReferralCode:  req.ReferralCode,
		OrderMethod:   req.OrderMethod,
		PromoCode:     req.PromoCode,
		Status:        entity.OrderStatusPending,
		TotalAmount:   totalAmount,
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

		// Calculate new total
		var totalAmount float64
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
			totalAmount += subtotal

			orderItems = append(orderItems, entity.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Price,
				Note:      item.Note,
			})
		}

		order.TotalAmount = totalAmount

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

func (s *orderService) GetAllOrders(companyID, branchID *uuid.UUID, pagination pkg.PaginationParams) ([]entity.OrderResponse, *pkg.PaginationMeta, error) {
	orders, total, err := s.orderRepo.FindAll(companyID, branchID, pagination)
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

	return &entity.OrderResponse{
		ID:            order.ID,
		CompanyID:     order.CompanyID,
		BranchID:      order.BranchID,
		CustomerName:  order.CustomerName,
		CustomerPhone: order.CustomerPhone,
		TableNumber:   order.TableNumber,
		Notes:         order.Notes,
		ReferralCode:  order.ReferralCode,
		OrderMethod:   order.OrderMethod,
		PromoCode:     order.PromoCode,
		Status:        order.Status,
		TotalAmount:   order.TotalAmount,
		OrderItems:    items,
		CreatedAt:     order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
