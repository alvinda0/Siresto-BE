package handler

import (
	"net/http"
	"project-name/internal/entity"
	"project-name/internal/service"
	"project-name/internal/websocket"
	"project-name/pkg"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder godoc
// @Summary Create a new order
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body entity.CreateOrderRequest true "Order data"
// @Success 201 {object} pkg.Response{data=entity.OrderResponse}
// @Router /api/v1/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req entity.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Get company_id and branch_id from context (set by auth middleware)
	companyID, exists := c.Get("company_id")
	if !exists {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Company ID not found", "")
		return
	}

	branchID, exists := c.Get("branch_id")
	if !exists {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Branch ID not found", "")
		return
	}

	order, err := h.orderService.CreateOrder(req, companyID.(uuid.UUID), branchID.(uuid.UUID))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to create order", err.Error())
		return
	}

	// Broadcast to WebSocket clients
	hub := websocket.GetHub()
	hub.BroadcastOrderUpdate("created", order, order.CompanyID, order.BranchID)

	pkg.SuccessResponse(c, http.StatusCreated, "Order created successfully", order)
}

// QuickCreateOrder godoc
// @Summary Quick create order (minimal fields)
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body entity.QuickOrderRequest true "Quick order data"
// @Success 201 {object} pkg.Response{data=entity.OrderResponse}
// @Router /api/v1/orders/quick [post]
func (h *OrderHandler) QuickCreateOrder(c *gin.Context) {
	var req entity.QuickOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Get company_id and branch_id from context (set by auth middleware)
	companyID, exists := c.Get("company_id")
	if !exists {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Company ID not found", "")
		return
	}

	branchID, exists := c.Get("branch_id")
	if !exists {
		pkg.ErrorResponse(c, http.StatusUnauthorized, "Branch ID not found", "")
		return
	}

	order, err := h.orderService.QuickCreateOrder(req, companyID.(uuid.UUID), branchID.(uuid.UUID))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to create order", err.Error())
		return
	}

	// Broadcast to WebSocket clients
	hub := websocket.GetHub()
	hub.BroadcastOrderUpdate("created", order, order.CompanyID, order.BranchID)

	pkg.SuccessResponse(c, http.StatusCreated, "Quick order created successfully", order)
}

// AddOrderItem godoc
// @Summary Add item to existing order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param item body entity.AddOrderItemRequest true "Item to add"
// @Success 200 {object} pkg.Response{data=entity.OrderResponse}
// @Router /api/v1/orders/quick/{id} [post]
func (h *OrderHandler) AddOrderItem(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID", err.Error())
		return
	}

	var req entity.AddOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	companyID, _ := c.Get("company_id")
	branchID, _ := c.Get("branch_id")

	order, err := h.orderService.AddOrderItem(orderID, req, companyID.(uuid.UUID), branchID.(uuid.UUID))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to add item to order", err.Error())
		return
	}

	// Broadcast to WebSocket clients
	hub := websocket.GetHub()
	hub.BroadcastOrderUpdate("updated", order, order.CompanyID, order.BranchID)

	pkg.SuccessResponse(c, http.StatusOK, "Item added to order successfully", order)
}

// CreatePublicOrder godoc
// @Summary Create a new order without authentication
// @Tags Public Orders
// @Accept json
// @Produce json
// @Param order body entity.CreatePublicOrderRequest true "Order data"
// @Success 201 {object} pkg.Response{data=entity.OrderResponse}
// @Router /api/v1/public/orders [post]
func (h *OrderHandler) CreatePublicOrder(c *gin.Context) {
	var req entity.CreatePublicOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	order, err := h.orderService.CreatePublicOrder(req)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to create order", err.Error())
		return
	}

	// Broadcast to WebSocket clients
	hub := websocket.GetHub()
	hub.BroadcastOrderUpdate("created", order, order.CompanyID, order.BranchID)

	pkg.SuccessResponse(c, http.StatusCreated, "Order created successfully", order)
}

// UpdateOrder godoc
// @Summary Update an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param order body entity.UpdateOrderRequest true "Order data"
// @Success 200 {object} pkg.Response{data=entity.OrderResponse}
// @Router /api/v1/orders/{id} [put]
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID", err.Error())
		return
	}

	var req entity.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	companyID, _ := c.Get("company_id")
	branchID, _ := c.Get("branch_id")

	order, err := h.orderService.UpdateOrder(id, req, companyID.(uuid.UUID), branchID.(uuid.UUID))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to update order", err.Error())
		return
	}

	// Broadcast to WebSocket clients
	hub := websocket.GetHub()
	hub.BroadcastOrderUpdate("updated", order, order.CompanyID, order.BranchID)

	pkg.SuccessResponse(c, http.StatusOK, "Order updated successfully", order)
}

// DeleteOrder godoc
// @Summary Delete an order
// @Tags Orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} pkg.Response
// @Router /api/v1/orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID", err.Error())
		return
	}

	companyID, _ := c.Get("company_id")
	branchID, _ := c.Get("branch_id")

	if err := h.orderService.DeleteOrder(id, companyID.(uuid.UUID), branchID.(uuid.UUID)); err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to delete order", err.Error())
		return
	}

	// Broadcast to WebSocket clients
	hub := websocket.GetHub()
	hub.BroadcastOrderUpdate("deleted", map[string]interface{}{"id": id}, companyID.(uuid.UUID), branchID.(uuid.UUID))

	pkg.SuccessResponse(c, http.StatusOK, "Order deleted successfully", nil)
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Tags Orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} pkg.Response{data=entity.OrderResponse}
// @Router /api/v1/orders/{id} [get]
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.ErrorResponse(c, http.StatusBadRequest, "Invalid order ID", err.Error())
		return
	}

	companyID, _ := c.Get("company_id")
	branchID, _ := c.Get("branch_id")

	companyUUID := companyID.(uuid.UUID)
	branchUUID := branchID.(uuid.UUID)

	order, err := h.orderService.GetOrderByID(id, &companyUUID, &branchUUID)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusNotFound, "Order not found", err.Error())
		return
	}

	pkg.SuccessResponse(c, http.StatusOK, "Order retrieved successfully", order)
}

// GetAllOrders godoc
// @Summary Get all orders
// @Tags Orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status (PENDING, CONFIRMED, PREPARING, READY, COMPLETED, CANCELLED)"
// @Param method query string false "Filter by order method (DINE_IN, TAKE_AWAY, DELIVERY)"
// @Param customer query string false "Search by customer name (partial match)"
// @Param order_id query string false "Search by order ID (partial match)"
// @Success 200 {object} pkg.Response{data=[]entity.OrderResponse}
// @Router /api/v1/orders [get]
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	pagination := pkg.GetPaginationParams(c)

	// Get filter parameters from query
	status := c.Query("status")
	method := c.Query("method")
	customer := c.Query("customer")
	orderID := c.Query("order_id")

	companyID, _ := c.Get("company_id")
	branchID, _ := c.Get("branch_id")

	companyUUID := companyID.(uuid.UUID)
	branchUUID := branchID.(uuid.UUID)

	orders, meta, err := h.orderService.GetAllOrders(&companyUUID, &branchUUID, status, method, customer, orderID, pagination)
	if err != nil {
		pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve orders", err.Error())
		return
	}

	pkg.SuccessResponseWithMeta(c, http.StatusOK, "Orders retrieved successfully", orders, meta)
}
