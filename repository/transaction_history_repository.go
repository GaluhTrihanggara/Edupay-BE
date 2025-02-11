package repository

import (
	"Edupay/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// TransactionHistoryRepository adalah interface untuk operasi CRUD pada entitas TransactionHistory
type TransactionHistoryRepository interface {
	GetAllHistoriesRepository(page, limit int) ([]*model.TransactionHistory, error)
	GetHistoryByIdRepository(id string) (*model.TransactionHistory, error)
	GetHistoriesByUserIdRepository(userId string) ([]*model.TransactionHistory, error)
	CreateHistoryRepository(history *model.TransactionHistory) (*model.TransactionHistory, error)
	CreateHistoryFromTransaction(transaction *model.Transaction) (*model.TransactionHistory, error) // Baru
	DeleteHistoryByIdRepository(id string) error
}

// transactionHistoryRepository adalah struct yang mengimplementasikan TransactionHistoryRepository
type transactionHistoryRepository struct {
	db *gorm.DB
}

// NewTransactionHistoryRepository membuat instance baru dari transactionHistoryRepository
func NewTransactionHistoryRepository(db *gorm.DB) *transactionHistoryRepository {
	return &transactionHistoryRepository{db}
}

// GetAllHistoriesRepository mengambil semua riwayat transaksi dengan pagination
func (r *transactionHistoryRepository) GetAllHistoriesRepository(page, limit int) ([]*model.TransactionHistory, error) {
	var histories []*model.TransactionHistory
	offset := (page - 1) * limit

	result := r.db.Preload("Transaction").Offset(offset).Limit(limit).Order("payment_date DESC").Find(&histories)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting transaction histories: %s", result.Error)
	}
	return histories, nil
}

// GetHistoryByIDRepository mengambil riwayat transaksi berdasarkan ID
func (r *transactionHistoryRepository) GetHistoryByIdRepository(id string) (*model.TransactionHistory, error) {
	var history model.TransactionHistory
	result := r.db.Preload("Transaction").First(&history, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction history with ID %s not found", id)
		}
		return nil, fmt.Errorf("error getting transaction history with ID %s: %s", id, result.Error)
	}
	return &history, nil
}

// GetHistoriesByUserIdRepository mengambil semua riwayat transaksi berdasarkan UserId
func (r *transactionHistoryRepository) GetHistoriesByUserIdRepository(userId string) ([]*model.TransactionHistory, error) {
	var histories []*model.TransactionHistory
	result := r.db.Preload("Transaction").Where("user_id = ?", userId).Order("payment_date DESC").Find(&histories)
	if result.Error != nil {
		return nil, result.Error
	}
	return histories, nil
}

// CreateHistoryRepository membuat riwayat transaksi baru
func (r *transactionHistoryRepository) CreateHistoryRepository(history *model.TransactionHistory) (*model.TransactionHistory, error) {
	result := r.db.Create(history)
	if result.Error != nil {
		return nil, result.Error
	}
	return history, nil
}

// CreateHistoryFromTransaction membuat riwayat transaksi dari data transaksi utama
func (r *transactionHistoryRepository) CreateHistoryFromTransaction(transaction *model.Transaction) (*model.TransactionHistory, error) {
	history := &model.TransactionHistory{
		UserId:        transaction.UserId,
		TransactionId: transaction.Id,
		Amount:        transaction.TotalPrice,
		Status:        transaction.Status,
		PaymentMethod: transaction.Method,
		PaymentDate:   transaction.TransactionDate,
	}
	result := r.db.Create(history)
	if result.Error != nil {
		return nil, fmt.Errorf("error creating transaction history: %w", result.Error)
	}
	return history, nil
}

// DeleteHistoryByIDRepository menghapus riwayat transaksi berdasarkan ID
func (r *transactionHistoryRepository) DeleteHistoryByIdRepository(id string) error {
	result := r.db.Delete(&model.TransactionHistory{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("transaction history not found")
	}
	return nil
}
