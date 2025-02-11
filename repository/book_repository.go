package repository

import (
	"Edupay/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// BookRepository adalah interface untuk operasi CRUD pada entitas Book
type BookRepository interface {
	GetAllBooksRepository(page, limit int, title, code string) ([]*model.Book, error)
	GetBookByIdRepository(id string) (*model.Book, error)
	CreateBookRepository(book *model.Book) (*model.Book, error)
	UpdateBookByIdRepository(id string, book *model.Book) (*model.Book, error)
	DeleteBookByIdRepository(id string) error
}

// bookRepository adalah struct yang mengimplementasikan BookRepository
type bookRepository struct {
	db *gorm.DB
}

// NewBookRepository membuat instance baru dari bookRepository
func NewBookRepository(db *gorm.DB) *bookRepository {
	return &bookRepository{db}
}

// GetAllBooksRepository mengambil semua buku dengan pagination dan pencarian berdasarkan judul
func (r *bookRepository) GetAllBooksRepository(page, limit int, title, code string) ([]*model.Book, error) {
	var books []*model.Book
	offset := (page - 1) * limit

	query := r.db.Offset(offset).Limit(limit)
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}

	result := query.Order("created_at DESC").Find(&books)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting books: %s", result.Error)
	}
	return books, nil
}

// GetBookByIDRepository mengambil buku berdasarkan ID
func (r *bookRepository) GetBookByIdRepository(id string) (*model.Book, error) {
	var book model.Book
	result := r.db.First(&book, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("book with ID %s not found", id)
		}
		return nil, fmt.Errorf("error getting book with ID %s: %s", id, result.Error)
	}
	return &book, nil
}

// CreateBookRepository membuat buku baru
func (r *bookRepository) CreateBookRepository(book *model.Book) (*model.Book, error) {
	var existingBook model.Book
	if err := r.db.Where("code = ?", book.Code).First(&existingBook).Error; err == nil {
		return nil, fmt.Errorf("book with code %s already exists", book.Code)
	}

	result := r.db.Create(book)
	if result.Error != nil {
		return nil, result.Error
	}
	return book, nil
}

// UpdateBookByIDRepository memperbarui data buku berdasarkan ID
func (r *bookRepository) UpdateBookByIdRepository(id string, book *model.Book) (*model.Book, error) {
	var existingBook model.Book
	if err := r.db.Where("code = ? AND id != ?", book.Code, id).First(&existingBook).Error; err == nil {
		return nil, fmt.Errorf("book with code %s already exists", book.Code)
	}

	result := r.db.Model(&model.Book{}).Where("id = ?", id).Updates(book)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("book not found")
	}
	return book, nil
}

// DeleteBookByIDRepository menghapus buku berdasarkan ID
func (r *bookRepository) DeleteBookByIdRepository(id string) error {
	result := r.db.Delete(&model.Book{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}
	return nil
}
