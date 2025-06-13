package mocks

import (
	model "Edupay/model"

	"github.com/stretchr/testify/mock"
)

type MockShirtRepository struct {
	mock.Mock
}

func (m *MockShirtRepository) CreateShirt(shirt *model.Shirt) error {
	args := m.Called(shirt)
	return args.Error(0)
}

func (m *MockShirtRepository) GetShirtByID(id string) (*model.Shirt, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Shirt), args.Error(1)
}

func (m *MockShirtRepository) GetAllShirts() ([]model.Shirt, error) {
	args := m.Called()
	return args.Get(0).([]model.Shirt), args.Error(1)
}

func (m *MockShirtRepository) UpdateShirt(shirt *model.Shirt) error {
	args := m.Called(shirt)
	return args.Error(0)
}

func (m *MockShirtRepository) DeleteShirt(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
