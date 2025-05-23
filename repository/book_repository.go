package repository

import (
	"Edupay/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookRepository interface {
	CreateBookRepository(book *model.Book) (*model.Book, error)
	GetBookByIDRepository(ID uuid.UUID) (*model.Book, error)
	GetBookByCodeRepository(code string) (*model.Book, error)
	GetAllBooksRepository(page, limit int, title, class string) ([]*model.Book, error)
	GetAvailableBooksRepository(page, limit int) ([]*model.Book, error)
	UpdateBookByIDRepository(ID uuid.UUID, book *model.Book) (*model.Book, error)
	DeleteBookByIDRepository(ID uuid.UUID) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

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

func (r *bookRepository) GetBookByIDRepository(ID uuid.UUID) (*model.Book, error) {
	var book model.Book
	if err := r.db.First(&book, "id = ?", ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("book not found")
		}
		return nil, fmt.Errorf("failed to get book: %v", err)
	}
	return &book, nil
}

func (r *bookRepository) GetBookByCodeRepository(code string) (*model.Book, error) {
	var book model.Book
	if err := r.db.First(&book, "code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to check book code: %v", err)
	}
	return &book, nil
}

func (r *bookRepository) GetAllBooksRepository(page, limit int, title, class string) ([]*model.Book, error) {
	var books []*model.Book
	query := r.db.Where("quantity > 0")
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	if class != "" {
		query = query.Where("class = ?", class)
	}

	err := query.Offset((page - 1) * limit).
		Limit(limit).
		Order("created_at DESC").
		Find(&books).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}
	return books, nil
}

func (r *bookRepository) GetAvailableBooksRepository(page, limit int) ([]*model.Book, error) {
	var books []*model.Book
	if err := r.db.Where("quantity > 0").Offset((page - 1) * limit).Limit(limit).Find(&books).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve books: %v", err)
	}
	return books, nil
}

func (r *bookRepository) UpdateBookByIDRepository(ID uuid.UUID, book *model.Book) (*model.Book, error) {
	if err := r.db.Model(&model.Book{}).Where("id = ?", ID).Updates(book).Error; err != nil {
		return nil, fmt.Errorf("failed to update book: %v", err)
	}
	return book, nil
}

func (r *bookRepository) DeleteBookByIDRepository(ID uuid.UUID) error {
	if err := r.db.Delete(&model.Book{}, "id = ?", ID).Error; err != nil {
		return fmt.Errorf("failed to delete book: %v", err)
	}
	return nil
}
