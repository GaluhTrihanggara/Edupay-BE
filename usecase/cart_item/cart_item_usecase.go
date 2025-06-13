package cart_item

import (
	"Edupay/model"
	"Edupay/repository"
	"fmt"

	"github.com/google/uuid"
)

type CartItemUseCase interface {
	AddCartItemUseCase(payload *model.CartItem) (*model.CartItem, error)
	GetAllCartItemsUseCase(page, limit int, cartID string) ([]*model.CartItem, error)
	GetCartItemsByCartIDUseCase(cartID uuid.UUID) ([]*model.CartItem, error)
	UpdateCartItemByIDUseCase(cartItemID uuid.UUID, payload *model.CartItem) (*model.CartItem, error)
	DeleteCartItemByIDUseCase(cartItemID uuid.UUID) error
}

type cartItemUseCase struct {
	cartItemRepository repository.CartItemRepository
	cartRepository     repository.CartRepository
}

func NewCartItemUseCase(
	cartItemRepository repository.CartItemRepository,
	cartRepository repository.CartRepository,
) *cartItemUseCase {
	return &cartItemUseCase{
		cartItemRepository: cartItemRepository,
		cartRepository:     cartRepository,
	}
}

func (uc *cartItemUseCase) AddCartItemUseCase(cartItem *model.CartItem) (*model.CartItem, error) {
	if cartItem.CartID == uuid.Nil {
		return nil, fmt.Errorf("cart ID cannot request payloadbe empty")
	}
	if cartItem.ItemID == uuid.Nil {
		return nil, fmt.Errorf("item ID cannot be empty")
	}
	if cartItem.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than zero")
	}
	if cartItem.ItemType != "book" && cartItem.ItemType != "shirt" {
		return nil, fmt.Errorf("item type must be 'book' or 'shirt'")
	}

	// Cek apakah cart ada dan aktif
	cart, err := uc.cartRepository.GetCartByIDRepository(cartItem.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}
	if cart == nil {
		return nil, fmt.Errorf("cart not found")
	}

	// Cek apakah item sudah ada di cart
	existingItem, err := uc.cartItemRepository.GetCartItemByItemIDRepository(cartItem.CartID, cartItem.ItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing cart item: %v", err)
	}
	if existingItem != nil {
		return nil, fmt.Errorf("item %s already exists in cart", cartItem.ItemID)
	}

	// Simpan cart item
	createdCartItem, err := uc.cartItemRepository.AddCartItemRepository(cartItem)
	if err != nil {
		return nil, err
	}

	return createdCartItem, nil
}

func (u *cartItemUseCase) GetAllCartItemsUseCase(page, limit int, cartID string) ([]*model.CartItem, error) {
	if page < 1 || limit < 1 {
		return nil, fmt.Errorf("page and limit must be greater than zero")
	}

	cartItems, err := u.cartItemRepository.GetAllCartItemsRepository(page, limit, cartID)
	if err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (uc *cartItemUseCase) GetCartItemsByCartIDUseCase(cartID uuid.UUID) ([]*model.CartItem, error) {
	if cartID == uuid.Nil {
		return nil, fmt.Errorf("cart ID cannot be empty")
	}
	cartItems, err := uc.cartItemRepository.GetCartItemsByCartIDRepository(cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %v", err)
	}
	return cartItems, nil
}

func (uc *cartItemUseCase) UpdateCartItemByIDUseCase(cartItemID uuid.UUID, payload *model.CartItem) (*model.CartItem, error) {
	if cartItemID == uuid.Nil {
		return nil, fmt.Errorf("cart item ID cannot be empty")
	}
	if payload == nil {
		return nil, fmt.Errorf("cart item payload cannot be nil")
	}
	if payload.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than zero")
	}

	// Cek apakah cart item ada
	existingItem, err := uc.cartItemRepository.GetCartItemByItemIDRepository(payload.CartID, payload.ItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing cart item: %v", err)
	}
	if existingItem == nil {
		return nil, fmt.Errorf("cart item not found")
	}

	updatedCartItem, err := uc.cartItemRepository.UpdateCartItemByIDRepository(cartItemID, payload)
	if err != nil {
		return nil, err
	}

	return updatedCartItem, nil
}

func (uc *cartItemUseCase) DeleteCartItemByIDUseCase(cartItemID uuid.UUID) error {
	if cartItemID == uuid.Nil {
		return fmt.Errorf("cart item ID cannot be empty")
	}
	return uc.cartItemRepository.DeleteCartItemByIDRepository(cartItemID)
}
