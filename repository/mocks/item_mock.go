package mocks

import (
	"Edupay/model"

	"github.com/stretchr/testify/mock"
)

type ItemRepository struct {
	mock.Mock
}

// DeleteItemByIDRepository mock for the method DeleteItemByIDRepository
func (_m *ItemRepository) DeleteItemByIDRepository(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllItemsRepository mock for the method GetAllItemsRepository
func (_m *ItemRepository) GetAllItemsRepository(page int, limit int, name string) ([]*model.Item, error) {
	ret := _m.Called(page, limit, name)

	var r0 []*model.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int, string) ([]*model.Item, error)); ok {
		return rf(page, limit, name)
	}
	if rf, ok := ret.Get(0).(func(int, int, string) []*model.Item); ok {
		r0 = rf(page, limit, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, string) error); ok {
		r1 = rf(page, limit, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetItemByIDRepository mock for the method GetItemByIDRepository
func (_m *ItemRepository) GetItemByIDRepository(id string) (*model.Item, error) {
	ret := _m.Called(id)

	var r0 *model.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.Item, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.Item); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertItemRepository mock for the method InsertItemRepository
func (_m *ItemRepository) InsertItemRepository(item *model.Item) error {
	ret := _m.Called(item)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Item) error); ok {
		r0 = rf(item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateItemByIDRepository mock for the method UpdateItemByIDRepository
func (_m *ItemRepository) UpdateItemByIDRepository(id string, item *model.Item) (*model.Item, error) {
	ret := _m.Called(id, item)

	var r0 *model.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *model.Item) (*model.Item, error)); ok {
		return rf(id, item)
	}
	if rf, ok := ret.Get(0).(func(string, *model.Item) *model.Item); ok {
		r0 = rf(id, item)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *model.Item) error); ok {
		r1 = rf(id, item)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewItemRepository creates a new instance of ItemRepository
// The first argument is usually *testing.T
func NewItemRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ItemRepository {
	mock := &ItemRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
