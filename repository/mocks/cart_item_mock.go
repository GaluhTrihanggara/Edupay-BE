package mocks

import (
	"Edupay/model"

	"github.com/stretchr/testify/mock"
)

type MockCartItemRepository struct {
	mock.Mock
}

func (m *MockCartItemRepository) AddCartItem(cartItem *model.CartItem) error {
	args := m.Called(cartItem)
	return args.Error(0)
}

func (m *MockCartItemRepository) GetCartItemsByCartID(cartID string) ([]model.CartItem, error) {
	args := m.Called(cartID)
	return args.Get(0).([]model.CartItem), args.Error(1)
}
