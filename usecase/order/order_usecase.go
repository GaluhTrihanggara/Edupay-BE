package order

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type OrderUseCase interface {
	CreateOrderUseCase(payload *model.Order) (*model.Order, error)
	GetAllOrdersUseCase(page, limit int, userID, status string) ([]*model.Order, error)
	GetOrderByIdUseCase(orderId uuid.UUID) (*model.Order, error)
	UpdateOrderByIdUseCase(orderId uuid.UUID, payload *model.Order) (*model.Order, error)
	DeleteOrderByIdUseCase(orderId uuid.UUID) error
}

type orderUseCase struct {
	orderRepository    repository.OrderRepository
	cartItemRepository repository.CartItemRepository
}

func NewOrderUseCase(orderRepository repository.OrderRepository, cartItemRepository repository.CartItemRepository) *orderUseCase {
	return &orderUseCase{
		orderRepository:    orderRepository,
		cartItemRepository: cartItemRepository,
	}
}

// CreateOrderUseCase creates a new order
func (uc *orderUseCase) CreateOrderUseCase(payload *model.Order) (*model.Order, error) {
	if payload == nil {
		return nil, errors.New("order payload cannot be nil")
	}
	if payload.UserID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}
	if payload.CartID == uuid.Nil {
		return nil, errors.New("cart ID cannot be empty")
	}

	// Validasi bahwa cart memiliki item
	cartItems, err := uc.cartItemRepository.GetCartItemsByCartIDRepository(payload.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate cart items: %v", err)
	}
	if len(cartItems) == 0 {
		return nil, errors.New("cannot create order: cart is empty")
	}

	// Set status awal
	payload.Status = "pending"

	return uc.orderRepository.CreateOrderRepository(payload)
}

func (u *orderUseCase) GetAllOrdersUseCase(page, limit int, userID, status string) ([]*model.Order, error) {
	if page < 1 || limit < 1 {
		return nil, fmt.Errorf("page and limit must be greater than zero")
	}

	orders, err := u.orderRepository.GetAllOrdersRepository(page, limit, userID, status)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderByIdUseCase retrieves an order by ID
func (uc *orderUseCase) GetOrderByIdUseCase(orderId uuid.UUID) (*model.Order, error) {
	if orderId == uuid.Nil {
		return nil, errors.New("orderID cannot be empty")
	}
	order, err := uc.orderRepository.GetOrderByIDRepository(orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}
	return order, nil
}

// UpdateOrderByIdUseCase updates an order by ID
func (uc *orderUseCase) UpdateOrderByIdUseCase(orderId uuid.UUID, payload *model.Order) (*model.Order, error) {
	if orderId == uuid.Nil {
		return nil, errors.New("order ID cannot be empty")
	}
	if payload == nil {
		return nil, errors.New("order payload cannot be nil")
	}
	if payload.Status != "" && !isValidOrderStatus(payload.Status) {
		return nil, fmt.Errorf("invalid order status: %s", payload.Status)
	}

	return uc.orderRepository.UpdateOrderByIDRepository(orderId, payload)
}

// DeleteOrderByIdUseCase deletes an order by ID
func (uc *orderUseCase) DeleteOrderByIdUseCase(orderId uuid.UUID) error {
	if orderId == uuid.Nil {
		return errors.New("orderID cannot be empty")
	}
	return uc.orderRepository.DeleteOrderByIDRepository(orderId)
}

// isValidOrderStatus validates the order status
func isValidOrderStatus(status string) bool {
	validStatues := []string{"pending", "completed", "cancelled"}
	for _, s := range validStatues {
		if s == status {
			return true
		}
	}
	return false
}
