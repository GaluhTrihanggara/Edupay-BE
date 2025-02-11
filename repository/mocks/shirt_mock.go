package mocks

import (
	model "Edupay/model"

	"github.com/stretchr/testify/mock"
)

type ShirtRepository struct {
	mock.Mock
}

// DeleteShirtByIDRepository mock for the method DeleteShirtByIDRepository
func (_m *ShirtRepository) DeleteShirtByIDRepository(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllShirtsRepository mock for the method GetAllShirtsRepository
func (_m *ShirtRepository) GetAllShirtsRepository(page int, limit int, name string) ([]*model.Shirt, error) {
	ret := _m.Called(page, limit, name)

	var r0 []*model.Shirt
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int, string) ([]*model.Shirt, error)); ok {
		return rf(page, limit, name)
	}
	if rf, ok := ret.Get(0).(func(int, int, string) []*model.Shirt); ok {
		r0 = rf(page, limit, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Shirt)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, string) error); ok {
		r1 = rf(page, limit, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetShirtByIDRepository mock for the method GetShirtByIDRepository
func (_m *ShirtRepository) GetShirtByIDRepository(id string) (*model.Shirt, error) {
	ret := _m.Called(id)

	var r0 *model.Shirt
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.Shirt, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.Shirt); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Shirt)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertShirtRepository mock for the method InsertShirtRepository
func (_m *ShirtRepository) InsertShirtRepository(shirt *model.Shirt) error {
	ret := _m.Called(shirt)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Shirt) error); ok {
		r0 = rf(shirt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateShirtByIDRepository mock for the method UpdateShirtByIDRepository
func (_m *ShirtRepository) UpdateShirtByIDRepository(id string, shirt *model.Shirt) (*model.Shirt, error) {
	ret := _m.Called(id, shirt)

	var r0 *model.Shirt
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *model.Shirt) (*model.Shirt, error)); ok {
		return rf(id, shirt)
	}
	if rf, ok := ret.Get(0).(func(string, *model.Shirt) *model.Shirt); ok {
		r0 = rf(id, shirt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Shirt)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *model.Shirt) error); ok {
		r1 = rf(id, shirt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewShirtRepository creates a new instance of ShirtRepository
// The first argument is usually *testing.T
func NewShirtRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ShirtRepository {
	mock := &ShirtRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
