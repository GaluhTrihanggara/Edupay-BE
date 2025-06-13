package cart

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type CartUseCase interface {
	CreateCartUseCase(payload *model.Cart) (*model.Cart, error)
	GetAllCartsUseCase(page, limit int, userID string) ([]*model.Cart, error)
	GetCartByIDUseCase(cartID uuid.UUID) (*model.Cart, error)
	GetCartByUserIDUseCase(userID uuid.UUID) (*model.Cart, error)
	UpdateCartByIdUseCase(cartId uuid.UUID, payload *model.Cart) (*model.Cart, error)
	DeleteCartByIdUseCase(cartId uuid.UUID) error
}

type cartUseCase struct {
	cartRepository repository.CartRepository
}

func NewCartUseCase(cartRepository repository.CartRepository) *cartUseCase {
	return &cartUseCase{
		cartRepository: cartRepository,
	}
}

// CreateCartUseCase creates a new cart
func (uc *cartUseCase) CreateCartUseCase(payload *model.Cart) (*model.Cart, error) {
	if payload == nil {
		return nil, errors.New("cart cannot be nil")
	}
	if payload.UserID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}
	return uc.cartRepository.CreateCartRepository(payload)
}

func (u *cartUseCase) GetAllCartsUseCase(page, limit int, userID string) ([]*model.Cart, error) {
	if page < 1 || limit < 1 {
		return nil, fmt.Errorf("page and limit must be greater than zero")
	}

	carts, err := u.cartRepository.GetAllCartsRepository(page, limit, userID)
	if err != nil {
		return nil, err
	}

	return carts, nil
}

func (u *cartUseCase) GetCartByIDUseCase(cartID uuid.UUID) (*model.Cart, error) {
	if cartID == uuid.Nil {
		return nil, fmt.Errorf("cart ID cannot be empty")
	}

	cart, err := u.cartRepository.GetCartByIDRepository(cartID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, fmt.Errorf("cart with ID %s not found", cartID)
	}

	return cart, nil
}

// GetCartByUserIDUseCase retrieves a cart by user ID
func (uc *cartUseCase) GetCartByUserIDUseCase(userID uuid.UUID) (*model.Cart, error) {
	if userID == uuid.Nil {
		return nil, errors.New("userID cannot be empty")
	}
	cart, err := uc.cartRepository.GetCartByUserIDRepository(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}
	return cart, nil
}

// UpdateCartByIdUseCase updates cart by ID
func (uc *cartUseCase) UpdateCartByIdUseCase(cartId uuid.UUID, payload *model.Cart) (*model.Cart, error) {
	if cartId == uuid.Nil {
		return nil, errors.New("cart ID cannot be empty")
	}
	if payload == nil {
		return nil, errors.New("cart payload cannot be nil")
	}
	if payload.UserID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}

	return uc.cartRepository.UpdateCartByIDRepository(cartId, payload)
}

// DeleteCartByIdUseCase deletes cart by ID
func (uc *cartUseCase) DeleteCartByIdUseCase(cartId uuid.UUID) error {
	if cartId == uuid.Nil {
		return errors.New("cartID cannot be empty")
	}
	return uc.cartRepository.DeleteCartByIDRepository(cartId)
}
