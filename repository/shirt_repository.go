package repository

import (
	"Edupay/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShirtRepository interface {
	CreateShirtRepository(shirt *model.Shirt) (*model.Shirt, error)
	GetShirtByIDRepository(ID uuid.UUID) (*model.Shirt, error)
	GetShirtByCodeRepository(code string) (*model.Shirt, error)
	GetAllShirtsRepository(page, limit int, name, size string) ([]*model.Shirt, error)
	UpdateShirtByIDRepository(ID uuid.UUID, shirt *model.Shirt) (*model.Shirt, error)
	DeleteShirtByIDRepository(ID uuid.UUID) error
}

type shirtRepository struct {
	db *gorm.DB
}

func NewShirtRepository(db *gorm.DB) ShirtRepository {
	return &shirtRepository{db: db}
}

func (r *shirtRepository) CreateShirtRepository(shirt *model.Shirt) (*model.Shirt, error) {
	var existingShirt model.Shirt
	if err := r.db.Where("code = ?", shirt.Code).First(&existingShirt).Error; err == nil {
		return nil, fmt.Errorf("book with code %s already exists", shirt.Code)
	}

	result := r.db.Create(shirt)
	if result.Error != nil {
		return nil, result.Error
	}
	return shirt, nil
}

func (r *shirtRepository) GetShirtByIDRepository(ID uuid.UUID) (*model.Shirt, error) {
	var shirt model.Shirt
	if err := r.db.First(&shirt, "id = ?", ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("shirt not found")
		}
		return nil, fmt.Errorf("failed to get shirt: %v", err)
	}
	return &shirt, nil
}

func (r *shirtRepository) GetShirtByCodeRepository(code string) (*model.Shirt, error) {
	var shirt model.Shirt
	if err := r.db.First(&shirt, "code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to check book code: %v", err)
	}
	return &shirt, nil
}

func (r *shirtRepository) GetAllShirtsRepository(page, limit int, name, size string) ([]*model.Shirt, error) {
	var shirts []*model.Shirt
	query := r.db.Where("quantity > 0")
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if size != "" {
		query = query.Where("size LIKE ?", "%"+size+"%")
	}

	err := query.Offset((page - 1) * limit).
		Limit(limit).
		Order("created_at DESC").
		Find(&shirts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}
	return shirts, nil
}

func (r *shirtRepository) UpdateShirtByIDRepository(ID uuid.UUID, shirt *model.Shirt) (*model.Shirt, error) {
	if err := r.db.Model(&model.Shirt{}).Where("id = ?", ID).Updates(shirt).Error; err != nil {
		return nil, fmt.Errorf("failed to update shirt: %v", err)
	}
	return shirt, nil
}

func (r *shirtRepository) DeleteShirtByIDRepository(ID uuid.UUID) error {
	if err := r.db.Delete(&model.Shirt{}, "id = ?", ID).Error; err != nil {
		return fmt.Errorf("failed to delete shirt: %v", err)
	}
	return nil
}
