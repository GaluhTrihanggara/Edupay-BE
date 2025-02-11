package repository

import (
	"Edupay/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// TransactionRepository adalah interface untuk operasi CRUD pada entitas Transaction
type TransactionRepository interface {
	GetAllTransactionsRepository(page, limit int, userId, itemId, billSemesterId string) ([]*model.Transaction, error)
	GetTransactionByIdRepository(id string) (*model.Transaction, error)
	CreateTransactionRepository(transaction *model.Transaction) (*model.Transaction, error)
	UpdateTransactionByIdRepository(id string, transaction *model.Transaction) (*model.Transaction, error)
	DeleteTransactionByIdRepository(id string) error
	GetTransactionsByUserIdRepository(userId string, page, limit int) ([]*model.Transaction, error)
	GetTransactionsByItemIdRepository(itemID string) ([]*model.Transaction, error)
	GetTransactionsByQueryRepository(query string, page, limit int) ([]*model.Transaction, error)
	GetTransactionsByPriceCountRepository() ([]*model.Transaction, error)
	GetTransactionsByStatusQueryRepository(query, status string, page, limit int) ([]*model.Transaction, error)
}

// transactionRepository adalah struct yang mengimplementasikan TransactionRepository
type transactionRepository struct {
	db *gorm.DB
}

// GetTransactionsByStatusQueryRepository implements TransactionRepository.
func (r *transactionRepository) GetTransactionsByStatusQueryRepository(query string, status string, page int, limit int) ([]*model.Transaction, error) {
	panic("unimplemented")
}

// NewTransactionRepository membuat instance baru dari transactionRepository
func NewTransactionRepository(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db}
}

// GetAllTransactionsRepository mengambil semua transaksi dengan pagination dan pencarian berdasarkan ParentId dan ItemId
func (r *transactionRepository) GetAllTransactionsRepository(page, limit int, userId, itemId, billSemesterId string) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	offset := (page - 1) * limit

	query := r.db.Preload("Items").Preload("BillSemester").Offset(offset).Limit(limit)

	if userId != "" {
		query = query.Where("user_id = ?", userId)
	}
	if itemId != "" {
		query = query.Where("item_id = ?", itemId)
	}
	if billSemesterId != "" {
		query = query.Where("bill_semester_id = ?", billSemesterId)
	}

	result := query.Order("transaction_date DESC").Find(&transactions)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting transactions: %s", result.Error)
	}
	return transactions, nil
}

// GetTransactionByIDRepository mengambil transaksi berdasarkan ID
func (r *transactionRepository) GetTransactionByIdRepository(id string) (*model.Transaction, error) {
	var transaction model.Transaction
	result := r.db.Preload("Items").Preload("BillSemester").First(&transaction, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction with ID %s not found", id)
		}
		return nil, fmt.Errorf("error getting transaction with ID %s: %s", id, result.Error)
	}
	return &transaction, nil
}

// CreateTransactionRepository membuat entri transaksi baru di database
func (r *transactionRepository) CreateTransactionRepository(transaction *model.Transaction) (*model.Transaction, error) {
	result := r.db.Create(transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

// UpdateTransactionByIDRepository memperbarui data transaksi berdasarkan ID
func (r *transactionRepository) UpdateTransactionByIdRepository(id string, transaction *model.Transaction) (*model.Transaction, error) {
	result := r.db.Model(&model.Transaction{}).Where("id = ?", id).Updates(transaction)
	if result.Error != nil {
		return nil, fmt.Errorf("error updating transaction: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("transaction not found")
	}

	// Fetch the updated transaction to return the complete updated record
	updatedTransaction, err := r.GetTransactionByIdRepository(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching updated transaction: %w", err)
	}

	return updatedTransaction, nil
}

// DeleteTransactionByIDRepository menghapus transaksi berdasarkan ID
func (r *transactionRepository) DeleteTransactionByIdRepository(id string) error {
	result := r.db.Delete(&model.Transaction{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("transaction not found")
	}
	return nil
}

// GetTransactionsByParentIDRepository mengambil semua transaksi berdasarkan UserId
func (r *transactionRepository) GetTransactionsByUserIdRepository(userId string, page, limit int) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	offset := (page - 1) * limit
	query := r.db.Preload("Items").Preload("BillSemester").Where("user_id = ?", userId).Offset(offset).Limit(limit)
	result := query.Order("created_at DESC").Find(&transactions)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting transaction: %s", result.Error)
	}
	return transactions, nil
}

// GetTransactionsByItemIDRepository mengambil semua transaksi berdasarkan ItemID
func (r *transactionRepository) GetTransactionsByItemIdRepository(itemID string) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	result := r.db.Where("item_id = ?", itemID).Order("transaction_date DESC").Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

func (r *transactionRepository) GetTransactionsByQueryRepository(query string, page, limit int) ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	offset := (page - 1) * limit

	queryString := "%" + query + "%"
	dbQuery := r.db.Preload("Items").Preload("BillSemester").
		Where("id LIKE ? OR status LIKE ? OR CAST(total_price AS TEXT) LIKE ?", queryString, queryString, queryString).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&transactions)

	if dbQuery.Error != nil {
		return nil, dbQuery.Error
	}

	return transactions, nil
}

func (r *transactionRepository) GetTransactionsByPriceCountRepository() ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	err := r.db.Where("status = ?", model.STATUS_SUCCESSFUL).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
