package payment_history

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"

	"github.com/google/uuid"
)

type PaymentHistoryUseCase interface {
	CreatePaymentHistoryUseCase(payload *model.PaymentHistory) error
	GetPaymentHistoryByUserIDUseCase(userID uuid.UUID) ([]*model.PaymentHistory, error)
}

type paymentHistoryUseCase struct {
	paymentHistoryRepository repository.PaymentHistoryRepository
}

func NewPaymentHistoryUseCase(paymentHistoryRepository repository.PaymentHistoryRepository) *paymentHistoryUseCase {
	return &paymentHistoryUseCase{
		paymentHistoryRepository: paymentHistoryRepository,
	}
}

// CreatePaymentHistoryUseCase creates a new payment history record
func (uc *paymentHistoryUseCase) CreatePaymentHistoryUseCase(payload *model.PaymentHistory) error {
	if payload == nil {
		return errors.New("payment history cannot be nil")
	}
	return uc.paymentHistoryRepository.CreatePaymentHistoryRepository(payload)
}

// GetPaymentHistoryByUserIDUseCase retrieves payment history by user ID
func (uc *paymentHistoryUseCase) GetPaymentHistoryByUserIDUseCase(userID uuid.UUID) ([]*model.PaymentHistory, error) {
	if userID == uuid.Nil {
		return nil, errors.New("userID cannot be empty")
	}
	return uc.paymentHistoryRepository.GetPaymentHistoryByUserIDRepository(userID)
}
