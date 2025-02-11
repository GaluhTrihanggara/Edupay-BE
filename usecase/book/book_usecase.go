package book

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"
	"fmt"
)

type BookUseCase interface {
	CreateBookUseCase(payload *model.Book) (*model.Book, error)
	GetAllBookUseCase(page, limit int, title, code string) ([]*model.Book, error)
	GetBookByIdUseCase(bookId string) (*model.Book, error)
	UpdateBookByIdUseCase(bookId string, payload *model.Book) (*model.Book, error)
	DeleteBookByIdUseCase(bookId string) error
}

type bookUseCase struct {
	bookRepository repository.BookRepository
}

func NewBookUseCase(bookRepository repository.BookRepository) *bookUseCase {
	return &bookUseCase{
		bookRepository: bookRepository,
	}
}

// CreateBookUseCase membuat buku baru
func (uc *bookUseCase) CreateBookUseCase(payload *model.Book) (*model.Book, error) {
	// Validasi unik kode buku
	book, err := uc.bookRepository.GetAllBooksRepository(1, 1, "", payload.Code)
	if err == nil && len(book) > 0 {
		return nil, fmt.Errorf("book with code %s already exists", payload.Code)
	}

	// Buat buku baru
	createdBook, err := uc.bookRepository.CreateBookRepository(payload)
	if err != nil {
		return nil, fmt.Errorf("error creating book in database: %w", err)
	}
	return createdBook, nil
}

// GetAllBookUseCase mengambil semua buku dengan filter title dan code
func (uc *bookUseCase) GetAllBookUseCase(page, limit int, title, code string) ([]*model.Book, error) {
	books, err := uc.bookRepository.GetAllBooksRepository(page, limit, title, code)
	if err != nil {
		return nil, err
	}
	return books, nil
}

// GetBookByIdUseCase mengambil buku berdasarkan ID
func (uc *bookUseCase) GetBookByIdUseCase(bookId string) (*model.Book, error) {
	book, err := uc.bookRepository.GetBookByIdRepository(bookId)
	if err != nil {
		return nil, errors.New("book not found")
	}
	return book, nil
}

// UpdateBookByIdUseCase memperbarui data buku berdasarkan ID
func (uc *bookUseCase) UpdateBookByIdUseCase(bookId string, payload *model.Book) (*model.Book, error) {
	// Validasi keberadaan buku
	existingBook, err := uc.bookRepository.GetBookByIdRepository(bookId)
	if err != nil {
		return nil, fmt.Errorf("book with ID %s not found: %v", bookId, err)
	}

	// Validasi unik kode buku jika diperbarui
	if payload.Code != existingBook.Code {
		books, err := uc.bookRepository.GetAllBooksRepository(1, 1, "", payload.Code)
		if err == nil && len(books) > 0 {
			return nil, fmt.Errorf("book with code %s already exists", payload.Code)
		}
	}

	// Perbarui data buku
	existingBook.Title = payload.Title
	existingBook.Class = payload.Class
	existingBook.Author = payload.Author
	existingBook.Price = payload.Price
	existingBook.Stock = payload.Stock
	existingBook.Code = payload.Code

	updatedBook, err := uc.bookRepository.UpdateBookByIdRepository(bookId, existingBook)
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}
	return updatedBook, nil
}

// DeleteBookByIdUseCase menghapus buku berdasarkan ID
func (uc *bookUseCase) DeleteBookByIdUseCase(bookId string) error {
	// Validasi keberadaan buku sebelum dihapus
	_, err := uc.bookRepository.GetBookByIdRepository(bookId)
	if err != nil {
		return fmt.Errorf("book with ID %s not found: %v", bookId, err)
	}

	// Hapus buku
	err = uc.bookRepository.DeleteBookByIdRepository(bookId)
	if err != nil {
		return fmt.Errorf("failed to delete book: %v", err)
	}
	return nil
}
