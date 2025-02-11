package mocks

import (
	model "Edupay/model"

	"github.com/stretchr/testify/mock"
)

type TransactionHistoryRepository struct {
	mock.Mock
}

// DeleteTransactionHistoryByIDRepository mock for the method DeleteTransactionHistoryByIDRepository
func (_m *TransactionHistoryRepository) DeleteTransactionHistoryByIDRepository(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllTransactionHistoriesRepository mock for the method GetAllTransactionHistoriesRepository
func (_m *TransactionHistoryRepository) GetAllTransactionHistoriesRepository(page int, limit int, studentId string) ([]*model.TransactionHistory, error) {
	ret := _m.Called(page, limit, studentId)

	var r0 []*model.TransactionHistory
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int, string) ([]*model.TransactionHistory, error)); ok {
		return rf(page, limit, studentId)
	}
	if rf, ok := ret.Get(0).(func(int, int, string) []*model.TransactionHistory); ok {
		r0 = rf(page, limit, studentId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TransactionHistory)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, string) error); ok {
		r1 = rf(page, limit, studentId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionHistoryByIdRepository mock for the method GetTransactionHistoryByIdRepository
func (_m *TransactionHistoryRepository) GetTransactionHistoryByIDRepository(id string) (*model.TransactionHistory, error) {
	ret := _m.Called(id)

	var r0 *model.TransactionHistory
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.TransactionHistory, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.TransactionHistory); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TransactionHistory)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertTransactionHistoryRepository mock for the method InsertTransactionHistoryRepository
func (_m *TransactionHistoryRepository) InsertTransactionHistoryRepository(history *model.TransactionHistory) error {
	ret := _m.Called(history)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.TransactionHistory) error); ok {
		r0 = rf(history)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTransactionHistoryRepository creates a new instance of TransactionHistoryRepository
// The first argument is usually *testing.T
func NewPaymentHistoryRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionHistoryRepository {
	mock := &TransactionHistoryRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
