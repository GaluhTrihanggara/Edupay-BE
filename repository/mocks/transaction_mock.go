package mocks

import (
	model "Edupay/model"

	"github.com/stretchr/testify/mock"
)

type TransactionRepository struct {
	mock.Mock
}

// DeleteTransactionByIDRepository mock for the method DeleteTransactionByIDRepository
func (_m *TransactionRepository) DeleteTransactionByIDRepository(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllTransactionsRepository mock for the method GetAllTransactionsRepository
func (_m *TransactionRepository) GetAllTransactionsRepository(page int, limit int, parentId string) ([]*model.Transaction, error) {
	ret := _m.Called(page, limit, parentId)

	var r0 []*model.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int, string) ([]*model.Transaction, error)); ok {
		return rf(page, limit, parentId)
	}
	if rf, ok := ret.Get(0).(func(int, int, string) []*model.Transaction); ok {
		r0 = rf(page, limit, parentId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, string) error); ok {
		r1 = rf(page, limit, parentId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionByIDRepository mock for the method GetTransactionByIDRepository
func (_m *TransactionRepository) GetTransactionByIDRepository(id string) (*model.Transaction, error) {
	ret := _m.Called(id)

	var r0 *model.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.Transaction, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.Transaction); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertTransactionRepository mock for the method InsertTransactionRepository
func (_m *TransactionRepository) InsertTransactionRepository(transaction *model.Transaction) error {
	ret := _m.Called(transaction)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Transaction) error); ok {
		r0 = rf(transaction)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTransactionByIDRepository mock for the method UpdateTransactionByIDRepository
func (_m *TransactionRepository) UpdateTransactionByIDRepository(id string, transaction *model.Transaction) (*model.Transaction, error) {
	ret := _m.Called(id, transaction)

	var r0 *model.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *model.Transaction) (*model.Transaction, error)); ok {
		return rf(id, transaction)
	}
	if rf, ok := ret.Get(0).(func(string, *model.Transaction) *model.Transaction); ok {
		r0 = rf(id, transaction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *model.Transaction) error); ok {
		r1 = rf(id, transaction)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTransactionRepository creates a new instance of TransactionRepository
// The first argument is usually *testing.T
func NewTransactionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionRepository {
	mock := &TransactionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
