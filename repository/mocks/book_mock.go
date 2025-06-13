package mocks

import (
	"Edupay/model"

	"github.com/stretchr/testify/mock"
)

type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) CreateBook(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) GetBookByID(id string) (*model.Book, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Book), args.Error(1)
}

func (m *MockBookRepository) GetAllBooks() ([]model.Book, error) {
	args := m.Called()
	return args.Get(0).([]model.Book), args.Error(1)
}

func (m *MockBookRepository) UpdateBook(book *model.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) DeleteBook(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
