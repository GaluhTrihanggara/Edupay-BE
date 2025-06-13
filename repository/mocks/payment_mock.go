package mocks

import (
	"Edupay/model"

	"github.com/stretchr/testify/mock"
)

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) CreatePayment(payment *model.Payment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) GetPaymentByOrderID(orderID string) (*model.Payment, error) {
	args := m.Called(orderID)
	return args.Get(0).(*model.Payment), args.Error(1)
}
