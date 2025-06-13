package repository

import (
	"Edupay/model"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderRepository(order *model.Order) (*model.Order, error)
	GetAllOrdersRepository(page, limit int, userID, status string) ([]*model.Order, error)
	GetOrderByIDRepository(ID uuid.UUID) (*model.Order, error)
	UpdateOrderByIDRepository(ID uuid.UUID, order *model.Order) (*model.Order, error)
	DeleteOrderByIDRepository(ID uuid.UUID) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrderRepository(order *model.Order) (*model.Order, error) {
	// Mulai transaksi
	tx := r.db.Begin()

	// Ambil cart items untuk menghitung total amount
	var cartItems []*model.CartItem
	if err := tx.Where("cart_id = ? AND deleted_at IS NULL", order.CartID).Preload("Book").Preload("Shirt").Find(&cartItems).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get cart items: %v", err)
	}

	// Hitung TotalAmount
	totalAmount := 0.0
	for _, item := range cartItems {
		totalAmount += item.Price
	}
	order.TotalAmount = totalAmount
	order.Status = "pending"
	order.DeletedAt = gorm.DeletedAt{}

	// Buat order
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	// Muat ulang order dengan relasi User dan Cart
	if err := tx.Preload("User").Preload("Cart.User").Where("id = ?", order.ID).First(order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to load order relations: %v", err)
	}

	tx.Commit()
	return order, nil
}

func (r *orderRepository) GetAllOrdersRepository(page, limit int, userID, status string) ([]*model.Order, error) {
	var orders []*model.Order
	query := r.db.Model(&model.Order{})

	query = query.Where("deleted_at IS NULL")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Offset((page - 1) * limit).
		Limit(limit).
		Order("created_at DESC").
		Preload("User").
		Preload("Cart.User").
		Preload("Cart.CartItems").
		Preload("Cart.CartItems.Cart.User").
		Preload("Cart.CartItems.Book").
		Preload("Cart.CartItems.Shirt").
		Find(&orders).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

func (r *orderRepository) GetOrderByIDRepository(ID uuid.UUID) (*model.Order, error) {
	var order model.Order
	if err := r.db.Preload("User").
		Preload("Cart.User").
		Preload("Cart.CartItems").
		Preload("Cart.CartItems.Cart.User").
		Preload("Cart.CartItems.Book").
		Preload("Cart.CartItems.Shirt").
		First(&order, "id = ? AND deleted_at IS NULL", ID).Error; err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}
	return &order, nil
}

func (r *orderRepository) UpdateOrderByIDRepository(ID uuid.UUID, order *model.Order) (*model.Order, error) {
	if err := r.db.Model(&model.Order{}).Where("id = ? AND deleted_at IS NULL", ID).Updates(order).Error; err != nil {
		return nil, fmt.Errorf("failed to update order: %v", err)
	}
	// Muat ulang dengan relasi setelah update
	if err := r.db.Preload("User").Preload("Cart.User").Preload("Cart.CartItems.Cart.User").Preload("Cart.CartItems.Book").Preload("Cart.CartItems.Shirt").Where("id = ?", ID).First(order).Error; err != nil {
		return nil, fmt.Errorf("failed to load updated order: %v", err)
	}
	return order, nil
}

func (r *orderRepository) DeleteOrderByIDRepository(ID uuid.UUID) error {
	if err := r.db.Model(&model.Order{}).Where("id = ?", ID).Update("deleted_at", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to delete order: %v", err)
	}
	return nil
}
