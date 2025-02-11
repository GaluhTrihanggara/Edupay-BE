package mocks

import (
	"Edupay/model"

	"github.com/stretchr/testify/mock"
)

type BookRepository struct {
	mock.Mock
}

// DeleteBookByIDRepository mock for the method DeleteBookByIDRepository
func (_m *BookRepository) DeleteBookByIDRepository(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllBooksRepository mock for the method GetAllBooksRepository
func (_m *BookRepository) GetAllBooksRepository(page int, limit int, title string) ([]*model.Book, error) {
	ret := _m.Called(page, limit, title)

	var r0 []*model.Book
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int, string) ([]*model.Book, error)); ok {
		return rf(page, limit, title)
	}
	if rf, ok := ret.Get(0).(func(int, int, string) []*model.Book); ok {
		r0 = rf(page, limit, title)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Book)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, string) error); ok {
		r1 = rf(page, limit, title)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBookByIDRepository mock for the method GetBookByIDRepository
func (_m *BookRepository) GetBookByIDRepository(id string) (*model.Book, error) {
	ret := _m.Called(id)

	var r0 *model.Book
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.Book, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.Book); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Book)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertBookRepository mock for the method InsertBookRepository
func (_m *BookRepository) InsertBookRepository(book *model.Book) error {
	ret := _m.Called(book)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Book) error); ok {
		r0 = rf(book)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateBookByIDRepository mock for the method UpdateBookByIDRepository
func (_m *BookRepository) UpdateBookByIDRepository(id string, book *model.Book) (*model.Book, error) {
	ret := _m.Called(id, book)

	var r0 *model.Book
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *model.Book) (*model.Book, error)); ok {
		return rf(id, book)
	}
	if rf, ok := ret.Get(0).(func(string, *model.Book) *model.Book); ok {
		r0 = rf(id, book)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Book)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *model.Book) error); ok {
		r1 = rf(id, book)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBookRepository creates a new instance of BookRepository
// The first argument is usually *testing.T
func NewBookRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *BookRepository {
	mock := &BookRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
