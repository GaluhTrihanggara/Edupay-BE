package transactionhistory

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"
	"fmt"
)

type TransactionHistoryUseCase interface {
	GetAllHistoriesUseCase(page, limit int) ([]*model.TransactionHistory, error)
	GetHistoryByIdUseCase(Id string) (*model.TransactionHistory, error)
	GetHistoriesByUserIdUseCase(userId string) ([]*model.TransactionHistory, error)                        // Diperbarui
	CreateHistoryFromTransactionUseCase(transaction *model.Transaction) (*model.TransactionHistory, error) // Baru
	DeleteHistoriesByIdUseCase(Id string) error
}

type transactionHistoryUseCase struct {
	transactionHistoryRepository repository.TransactionHistoryRepository
}

func NewTransactionHistoryUseCase(transactionHistoryRepository repository.TransactionHistoryRepository) *transactionHistoryUseCase {
	return &transactionHistoryUseCase{
		transactionHistoryRepository: transactionHistoryRepository,
	}
}

// GetAllHistoriesUseCase mengambil semua riwayat transaksi dengan pagination
func (uc *transactionHistoryUseCase) GetAllHistoriesUseCase(page, limit int) ([]*model.TransactionHistory, error) {
	histories, err := uc.transactionHistoryRepository.GetAllHistoriesRepository(page, limit)
	if err != nil {
		return nil, fmt.Errorf("error retrieving histories from database: %w", err)
	}
	return histories, nil
}

// GetHistoryByIdUseCase mengambil riwayat transaksi berdasarkan ID
func (uc *transactionHistoryUseCase) GetHistoryByIdUseCase(Id string) (*model.TransactionHistory, error) {
	history, err := uc.transactionHistoryRepository.GetHistoryByIdRepository(Id)
	if err != nil {
		return nil, errors.New("transaction history not found")
	}
	return history, nil
}

// GetHistoriesByUserIdUseCase mengambil semua riwayat transaksi berdasarkan UserId
func (uc *transactionHistoryUseCase) GetHistoriesByUserIdUseCase(userId string) ([]*model.TransactionHistory, error) {
	histories, err := uc.transactionHistoryRepository.GetHistoriesByUserIdRepository(userId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving histories for user ID %s: %w", userId, err)
	}
	return histories, nil
}

// CreateHistoryFromTransactionUseCase membuat riwayat transaksi dari data transaksi utama
func (uc *transactionHistoryUseCase) CreateHistoryFromTransactionUseCase(transaction *model.Transaction) (*model.TransactionHistory, error) {
	// Validasi data transaksi
	if transaction == nil {
		return nil, errors.New("invalid transaction data")
	}

	// Buat riwayat transaksi menggunakan repository
	history, err := uc.transactionHistoryRepository.CreateHistoryFromTransaction(transaction)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction history: %w", err)
	}
	return history, nil
}

// DeleteHistoryByIdUseCase menghapus riwayat transaksi berdasarkan ID
func (uc *transactionHistoryUseCase) DeleteHistoriesByIdUseCase(Id string) error {
	err := uc.transactionHistoryRepository.DeleteHistoryByIdRepository(Id)
	if err != nil {
		return fmt.Errorf("history with ID %s not found: %w", Id, err)
	}
	return nil
}
