package mocks

import (
	"Edupay/model"

	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(order *model.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrderByID(id string) (*model.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Order), args.Error(1)
}
