package repository

import (
	"Edupay/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository interface {
	CreateCartRepository(cart *model.Cart) (*model.Cart, error)
	GetAllCartsRepository(page, limit int, userID string) ([]*model.Cart, error)
	GetCartByIDRepository(cartID uuid.UUID) (*model.Cart, error)
	GetCartByUserIDRepository(userID uuid.UUID) (*model.Cart, error)
	UpdateCartByIDRepository(ID uuid.UUID, cart *model.Cart) (*model.Cart, error)
	DeleteCartByIDRepository(ID uuid.UUID) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) CreateCartRepository(cart *model.Cart) (*model.Cart, error) {
	var existingCart model.Cart
	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", cart.UserID).First(&existingCart).Error; err == nil {
		return nil, fmt.Errorf("cart already exists for user %s", cart.UserID)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing cart: %v", err)
	}

	// Simpan cart ke database
	if err := r.db.Create(cart).Error; err != nil {
		return nil, fmt.Errorf("failed to create cart: %v", err)
	}

	// Ambil cart beserta data User
	var createdCart model.Cart
	if err := r.db.Preload("User").First(&createdCart, cart.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load cart with user data: %v", err)
	}

	return &createdCart, nil
}

func (r *cartRepository) GetAllCartsRepository(page, limit int, userID string) ([]*model.Cart, error) {
	var carts []*model.Cart
	query := r.db.Model(&model.Cart{})

	// Hanya ambil cart yang belum dihapus (soft delete)
	query = query.Where("deleted_at IS NULL")

	// Filter berdasarkan user_id jika diberikan
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Pagination dan pengurutan
	err := query.Offset((page - 1) * limit).
		Limit(limit).
		Order("created_at DESC").
		Preload("User").            // Memuat relasi User
		Preload("CartItems").       // Memuat relasi CartItems
		Preload("CartItems.Book").  // Memuat relasi Book di CartItem
		Preload("CartItems.Shirt"). // Memuat relasi Shirt di CartItem
		Find(&carts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get carts: %w", err)
	}
	return carts, nil
}

func (r *cartRepository) GetCartByIDRepository(cartID uuid.UUID) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Where("id = ? AND deleted_at IS NULL", cartID).
		Preload("User").            // Memuat relasi User
		Preload("CartItems").       // Memuat relasi CartItems
		Preload("CartItems.Book").  // Memuat relasi Book di CartItem
		Preload("CartItems.Shirt"). // Memuat relasi Shirt di CartItem
		First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Cart tidak ditemukan
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	return &cart, nil
}

func (r *cartRepository) GetCartByUserIDRepository(userID uuid.UUID) (*model.Cart, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	var cart model.Cart
	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).
		Preload("User").            // Memuat relasi User
		Preload("CartItems").       // Memuat relasi CartItems
		Preload("CartItems.Book").  // Memuat relasi ke Book
		Preload("CartItems.Shirt"). // Memuat relasi ke Shirt
		First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("cart not found for user %s", userID)
		}
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}
	return &cart, nil
}

func (r *cartRepository) UpdateCartByIDRepository(ID uuid.UUID, cart *model.Cart) (*model.Cart, error) {
	var existingCart model.Cart
	if err := r.db.First(&existingCart, "id = ? AND deleted_at IS NULL", ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, fmt.Errorf("failed to check cart: %v", err)
	}

	if err := r.db.Model(&model.Cart{}).Where("id = ?", ID).Updates(cart).Error; err != nil {
		return nil, fmt.Errorf("failed to update cart: %v", err)
	}
	return cart, nil
}

func (r *cartRepository) DeleteCartByIDRepository(ID uuid.UUID) error {
	var cart model.Cart
	if err := r.db.First(&cart, "id = ? AND deleted_at IS NULL", ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cart not found")
		}
		return fmt.Errorf("failed to check cart: %v", err)
	}

	if err := r.db.Delete(&model.Cart{}, "id = ?", ID).Error; err != nil {
		return fmt.Errorf("failed to delete cart: %v", err)
	}
	return nil
}
