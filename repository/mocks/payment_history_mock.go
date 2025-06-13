package mocks

import (
	"Edupay/model"

	"github.com/stretchr/testify/mock"
)

type MockPaymentHistoryRepository struct {
	mock.Mock
}

func (m *MockPaymentHistoryRepository) CreatePaymentHistory(history *model.PaymentHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockPaymentHistoryRepository) GetPaymentHistoryByUserID(userID string) ([]model.PaymentHistory, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.PaymentHistory), args.Error(1)
}
