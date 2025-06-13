package repository

import (
	"Edupay/model"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentHistoryRepository interface {
	CreatePaymentHistoryRepository(history *model.PaymentHistory) error
	GetPaymentHistoryByUserIDRepository(userID uuid.UUID) ([]*model.PaymentHistory, error)
}

type paymentHistoryRepository struct {
	db *gorm.DB
}

func NewPaymentHistoryRepository(db *gorm.DB) PaymentHistoryRepository {
	return &paymentHistoryRepository{db: db}
}

func (r *paymentHistoryRepository) CreatePaymentHistoryRepository(history *model.PaymentHistory) error {
	if err := r.db.Create(history).Error; err != nil {
		return fmt.Errorf("failed to create payment history: %v", err)
	}
	return nil
}

func (r *paymentHistoryRepository) GetPaymentHistoryByUserIDRepository(userID uuid.UUID) ([]*model.PaymentHistory, error) {
	var histories []*model.PaymentHistory
	if err := r.db.Where("user_id = ?", userID).Find(&histories).Error; err != nil {
		return nil, fmt.Errorf("failed to get payment history: %v", err)
	}
	return histories, nil
}
