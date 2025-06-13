package repository

import (
	"Edupay/model"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	AddCartItemRepository(cartItem *model.CartItem) (*model.CartItem, error)
	GetAllCartItemsRepository(page, limit int, cartID string) ([]*model.CartItem, error)
	GetCartItemByItemIDRepository(cartID, itemID uuid.UUID) (*model.CartItem, error)
	GetCartItemsByCartIDRepository(cartID uuid.UUID) ([]*model.CartItem, error)
	UpdateCartItemByIDRepository(ID uuid.UUID, cartItem *model.CartItem) (*model.CartItem, error)
	DeleteCartItemByIDRepository(ID uuid.UUID) error
}

type cartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) CartItemRepository {
	return &cartItemRepository{db: db}
}

func (r *cartItemRepository) AddCartItemRepository(cartItem *model.CartItem) (*model.CartItem, error) {
	if cartItem.ItemType != "book" && cartItem.ItemType != "shirt" {
		return nil, fmt.Errorf("invalid item type: %s", cartItem.ItemType)
	}

	// Validasi item_id berdasarkan item_type
	if cartItem.ItemType == "book" {
		var book model.Book
		if err := r.db.Where("id = ? AND deleted_at IS NULL", cartItem.ItemID).First(&book).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("book not found: record not found")
			}
			return nil, fmt.Errorf("failed to find book: %w", err)
		}
		if book.Quantity < cartItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for book %s: available %d, requested %d", book.Title, book.Quantity, cartItem.Quantity)
		}
		cartItem.Price = book.Price * float64(cartItem.Quantity)
	} else if cartItem.ItemType == "shirt" {
		var shirt model.Shirt
		if err := r.db.Where("id = ? AND deleted_at IS NULL", cartItem.ItemID).First(&shirt).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("shirt not found: record not found")
			}
			return nil, fmt.Errorf("failed to find shirt: %w", err)
		}
		if shirt.Quantity < cartItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for shirt %s: available %d, requested %d", shirt.Name, shirt.Quantity, cartItem.Quantity)
		}
		cartItem.Price = shirt.Price * float64(cartItem.Quantity)
	}

	cartItem.CreatedAt = time.Now()
	cartItem.UpdatedAt = time.Now()

	// Gunakan transaksi untuk memastikan integritas data
	tx := r.db.Begin()
	if err := tx.Create(cartItem).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to add cart item: %v", err)
	}

	// Memuat relasi setelah create
	if err := tx.Preload("Book").Preload("Shirt").Preload("Cart").Preload("Cart.User").Where("id = ?", cartItem.ID).First(cartItem).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to load cart item relations: %v", err)
	}

	tx.Commit()
	return cartItem, nil
}

func (r *cartItemRepository) GetAllCartItemsRepository(page, limit int, cartID string) ([]*model.CartItem, error) {
	var cartItems []*model.CartItem
	query := r.db.Model(&model.CartItem{})

	query = query.Where("deleted_at IS NULL")

	if cartID != "" {
		query = query.Where("cart_id = ?", cartID)
	}

	err := query.Offset((page - 1) * limit).
		Limit(limit).
		Order("created_at DESC").
		Preload("Book").
		Preload("Shirt").
		Preload("Cart").
		Preload("Cart.User").
		Find(&cartItems).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	return cartItems, nil
}

func (r *cartItemRepository) GetCartItemByItemIDRepository(cartID, itemID uuid.UUID) (*model.CartItem, error) {
	var cartItem model.CartItem
	err := r.db.Where("cart_id = ? AND item_id = ? AND deleted_at IS NULL", cartID, itemID).
		Preload("Book").
		Preload("Shirt").
		Preload("Cart").
		Preload("Cart.User").
		First(&cartItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get cart item: %w", err)
	}
	return &cartItem, nil
}

func (r *cartItemRepository) GetCartItemsByCartIDRepository(cartID uuid.UUID) ([]*model.CartItem, error) {
	var cartItems []*model.CartItem
	if err := r.db.Where("cart_id = ? AND deleted_at IS NULL", cartID).
		Preload("Book").
		Preload("Shirt").
		Preload("Cart").
		Preload("Cart.User").
		Find(&cartItems).Error; err != nil {
		return nil, fmt.Errorf("failed to get cart items: %v", err)
	}
	return cartItems, nil
}

func (r *cartItemRepository) UpdateCartItemByIDRepository(ID uuid.UUID, cartItem *model.CartItem) (*model.CartItem, error) {
	var existingCartItem model.CartItem
	if err := r.db.Where("id = ? AND deleted_at IS NULL", ID).
		Preload("Book").
		Preload("Shirt").
		Preload("Cart.User").First(&existingCartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("cart item not found")
		}
		return nil, fmt.Errorf("failed to find cart item: %v", err)
	}

	// Update fields
	existingCartItem.Quantity = cartItem.Quantity
	existingCartItem.Price = cartItem.Price
	existingCartItem.UpdatedAt = time.Now()

	if cartItem.ItemType == "book" {
		var book model.Book
		if err := r.db.Where("id = ? AND deleted_at IS NULL", existingCartItem.ItemID).First(&book).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("book not found: record not found")
			}
			return nil, fmt.Errorf("failed to find book: %w", err)
		}
		if book.Quantity < cartItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for book %s: available %d, requested %d", book.Title, book.Quantity, cartItem.Quantity)
		}
		existingCartItem.Price = book.Price * float64(cartItem.Quantity)
	} else if cartItem.ItemType == "shirt" {
		var shirt model.Shirt
		if err := r.db.Where("id = ? AND deleted_at IS NULL", existingCartItem.ItemID).First(&shirt).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("shirt not found: record not found")
			}
			return nil, fmt.Errorf("failed to find shirt: %w", err)
		}
		if shirt.Quantity < cartItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for shirt %s: available %d, requested %d", shirt.Name, shirt.Quantity, cartItem.Quantity)
		}
		existingCartItem.Price = shirt.Price * float64(cartItem.Quantity)
	}

	if err := r.db.Save(&existingCartItem).Error; err != nil {
		return nil, fmt.Errorf("failed to update cart item: %v", err)
	}

	// Reload relasi
	if err := r.db.Preload("Book").Preload("Shirt").Where("id = ?", existingCartItem.ID).First(&existingCartItem).Error; err != nil {
		return nil, fmt.Errorf("failed to load cart item relations: %v", err)
	}

	return &existingCartItem, nil
}

func (r *cartItemRepository) DeleteCartItemByIDRepository(ID uuid.UUID) error {
	var cartItem model.CartItem
	if err := r.db.Where("id = ? AND deleted_at IS NULL", ID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cart item not found")
		}
		return fmt.Errorf("failed to find cart item: %v", err)
	}

	if err := r.db.Model(&cartItem).Update("deleted_at", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to delete cart item: %v", err)
	}
	return nil
}
