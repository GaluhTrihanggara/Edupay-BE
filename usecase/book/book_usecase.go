package book

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type BookUseCase interface {
	CreateBookUseCase(payload *model.Book) (*model.Book, error)
	GetAllBookUseCase(page, limit int, title, class string) ([]*model.Book, error)
	GetBookByIDUseCase(bookId uuid.UUID) (*model.Book, error)
	GetBookByCodeUseCase(code string) (*model.Book, error)
	UpdateBookByIDUseCase(bookId uuid.UUID, payload *model.Book) (*model.Book, error)
	DeleteBookByIDUseCase(bookId uuid.UUID) error
}

type bookUseCase struct {
	bookRepository repository.BookRepository
}

func NewBookUseCase(bookRepository repository.BookRepository) *bookUseCase {
	return &bookUseCase{
		bookRepository: bookRepository,
	}
}

// CreateBookUseCase creates a new book
func (uc *bookUseCase) CreateBookUseCase(payload *model.Book) (*model.Book, error) {
	// Validate unique title or code if needed (assuming `Title` is unique)
	book, err := uc.bookRepository.GetAllBooksRepository(1, 1, "", payload.Class) // Check if book already exists
	if err == nil && len(book) > 0 {
		return nil, fmt.Errorf("book with class %s already exists", payload.Class)
	}

	// Save new book
	createdBook, err := uc.bookRepository.CreateBookRepository(payload)
	if err != nil {
		return nil, fmt.Errorf("error creating book in database: %w", err)
	}
	return createdBook, nil
}

// GetAllBookUseCase retrieves all books with pagination
func (uc *bookUseCase) GetAllBookUseCase(page, limit int, title, class string) ([]*model.Book, error) {
	books, err := uc.bookRepository.GetAllBooksRepository(page, limit, title, class)
	if err != nil {
		return nil, err
	}
	return books, nil
}

// GetBookByIdUseCase retrieves a book by ID
func (uc *bookUseCase) GetBookByIDUseCase(bookId uuid.UUID) (*model.Book, error) {
	book, err := uc.bookRepository.GetBookByIDRepository(bookId)
	if err != nil {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (uc *bookUseCase) GetBookByCodeUseCase(code string) (*model.Book, error) {
	book, err := uc.bookRepository.GetBookByCodeRepository(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get book by code %s: %w", code, err)
	}
	if book == nil {
		return nil, fmt.Errorf("book with code %s not found", code)
	}
	return book, nil
}

// UpdateBookByIdUseCase updates book information by ID
func (uc *bookUseCase) UpdateBookByIDUseCase(bookId uuid.UUID, payload *model.Book) (*model.Book, error) {
	// Validate book existence
	existingBook, err := uc.bookRepository.GetBookByIDRepository(bookId)
	if err != nil {
		return nil, fmt.Errorf("book with ID %s not found: %v", bookId, err)
	}

	// Validate unique title if it's updated
	if payload.Code != existingBook.Code {
		bookWithCode, err := uc.bookRepository.GetAllBooksRepository(1, 1, payload.Class, payload.Class)
		if err == nil && len(bookWithCode) > 0 {
			return nil, fmt.Errorf("book with code %s already exists", payload.Code)
		}
	}

	// Update book fields
	existingBook.Code = payload.Code
	existingBook.Title = payload.Title
	existingBook.Class = payload.Class
	existingBook.Author = payload.Author
	existingBook.Price = payload.Price
	existingBook.Quantity = payload.Quantity

	updatedBook, err := uc.bookRepository.UpdateBookByIDRepository(bookId, existingBook)
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}
	return updatedBook, nil
}

// DeleteBookByIdUseCase deletes a book by ID
func (uc *bookUseCase) DeleteBookByIDUseCase(bookId uuid.UUID) error {
	// Validate book existence before deletion
	_, err := uc.bookRepository.GetBookByIDRepository(bookId)
	if err != nil {
		return fmt.Errorf("book with ID %s not found: %v", bookId, err)
	}

	// Delete book
	err = uc.bookRepository.DeleteBookByIDRepository(bookId)
	if err != nil {
		return fmt.Errorf("failed to delete book: %v", err)
	}
	return nil
}
