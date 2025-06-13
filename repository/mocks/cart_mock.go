package mocks

import (
	"Edupay/model"

	"github.com/stretchr/testify/mock"
)

type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) CreateCart(cart *model.Cart) error {
	args := m.Called(cart)
	return args.Error(0)
}

func (m *MockCartRepository) GetCartByUserID(userID string) (*model.Cart, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *MockCartRepository) UpdateCart(cart *model.Cart) error {
	args := m.Called(cart)
	return args.Error(0)
}
