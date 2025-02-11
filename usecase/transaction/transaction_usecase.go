package transaction

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"
	"fmt"
	"strings"
)

type TransactionUseCase interface {
	CreateTransactionUseCase(transaction *model.Transaction) (*model.Transaction, error)
	GetAllTransactionUseCase(page, limit int, userId, itemId, billSemesterId string) ([]*model.Transaction, error)
	GetTransactionByIdUseCase(transactionId string) (*model.Transaction, error)
	GetTransactionByUserIdUseCase(userId string, page, limit int) ([]*model.Transaction, error)
	GetTransactionByQueryUseCase(query string, page, limit int) ([]*model.Transaction, error)
	GetTransactionsByStatusQueryUseCase(query, status string, page, limit int) ([]*model.Transaction, error)
	GetTransactionsByPriceCountUseCase() ([]*model.Transaction, error)
	UpdateTransactionByIdUseCase(transactionId string, transaction *model.Transaction) (*model.Transaction, error)
}

type transactionUseCase struct {
	transactionRepository repository.TransactionRepository
}

func NewTransactionUseCase(transactionRepository repository.TransactionRepository) *transactionUseCase {
	return &transactionUseCase{
		transactionRepository: transactionRepository,
	}
}

// CreateTransactionUseCase membuat transaksi baru
func (uc *transactionUseCase) CreateTransactionUseCase(transaction *model.Transaction) (*model.Transaction, error) {
	// Validasi bahwa item-item yang terkait valid
	if len(transaction.Items) > 0 {
		for _, item := range transaction.Items {
			if item.Code == "" {
				return nil, errors.New("item code is required for each item in the transaction")
			}
		}
	}

	// Buat transaksi
	createdTransaction, err := uc.transactionRepository.CreateTransactionRepository(transaction)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction in database: %w", err)
	}
	return createdTransaction, nil
}

// GetAllTransactionUseCase mengambil semua transaksi dengan filter userId, itemId, dan billSemesterId
func (uc *transactionUseCase) GetAllTransactionUseCase(page, limit int, userId, itemId, billSemesterId string) ([]*model.Transaction, error) {
	transactions, err := uc.transactionRepository.GetAllTransactionsRepository(page, limit, userId, itemId, billSemesterId)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetTransactionByIdUseCase mengambil transaksi berdasarkan ID
func (uc *transactionUseCase) GetTransactionByIdUseCase(transactionId string) (*model.Transaction, error) {
	transaction, err := uc.transactionRepository.GetTransactionByIdRepository(transactionId)
	if err != nil {
		return nil, errors.New("transaction not found")
	}
	return transaction, nil
}

// GetTransactionByUserIdUseCase mengambil transaksi berdasarkan userId
func (uc *transactionUseCase) GetTransactionByUserIdUseCase(userId string, page, limit int) ([]*model.Transaction, error) {
	transactions, err := uc.transactionRepository.GetTransactionsByUserIdRepository(userId, page, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting transactions for user %s: %w", userId, err)
	}
	return transactions, nil
}

// GetTransactionByQueryUseCase mengambil transaksi berdasarkan query
func (uc *transactionUseCase) GetTransactionByQueryUseCase(query string, page, limit int) ([]*model.Transaction, error) {
	lowercaseQuery := strings.ToLower(query)
	transactions, err := uc.transactionRepository.GetTransactionsByQueryRepository(lowercaseQuery, page, limit)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetTransactionsByStatusQueryUseCase mengambil transaksi berdasarkan query dan status
func (uc *transactionUseCase) GetTransactionsByStatusQueryUseCase(query, status string, page, limit int) ([]*model.Transaction, error) {
	lowercaseQuery := strings.ToLower(query)
	lowercaseStatus := strings.ToLower(status)
	transactions, err := uc.transactionRepository.GetTransactionsByStatusQueryRepository(lowercaseQuery, lowercaseStatus, page, limit)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetTransactionsByPriceCountUseCase menghitung total harga berdasarkan status
func (uc *transactionUseCase) GetTransactionsByPriceCountUseCase() ([]*model.Transaction, error) {
	transactions, err := uc.transactionRepository.GetTransactionsByPriceCountRepository()
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// UpdateTransactionByIdUseCase memperbarui transaksi berdasarkan ID
func (uc *transactionUseCase) UpdateTransactionByIdUseCase(transactionId string, transaction *model.Transaction) (*model.Transaction, error) {
	// Validasi bahwa transaksi ada
	existingTransaction, err := uc.transactionRepository.GetTransactionByIdRepository(transactionId)
	if err != nil {
		return nil, fmt.Errorf("error getting transaction: %w", err)
	}

	if existingTransaction == nil {
		return nil, errors.New("transaction not found")
	}

	// Update transaksi
	updatedTransaction, err := uc.transactionRepository.UpdateTransactionByIdRepository(transactionId, transaction)
	if err != nil {
		return nil, fmt.Errorf("error updating transaction: %w", err)
	}

	return updatedTransaction, nil
}
