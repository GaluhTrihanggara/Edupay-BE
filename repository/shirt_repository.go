package repository

import (
	"Edupay/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ShirtRepository adalah interface untuk operasi CRUD pada entitas Shirt
type ShirtRepository interface {
	GetAllShirtsRepository(page, limit int, name, size, code string) ([]*model.Shirt, error)
	GetShirtByIdRepository(id string) (*model.Shirt, error)
	CreateShirtRepository(shirt *model.Shirt) (*model.Shirt, error)
	UpdateShirtByIdRepository(id string, shirt *model.Shirt) (*model.Shirt, error)
	DeleteShirtByIdRepository(id string) error
}

// shirtRepository adalah struct yang mengimplementasikan ShirtRepository
type shirtRepository struct {
	db *gorm.DB
}

// NewShirtRepository membuat instance baru dari shirtRepository
func NewShirtRepository(db *gorm.DB) *shirtRepository {
	return &shirtRepository{db}
}

// GetAllShirtsRepository mengambil semua shirt dengan pagination dan pencarian berdasarkan nama dan ukuran
func (r *shirtRepository) GetAllShirtsRepository(page, limit int, name, size, code string) ([]*model.Shirt, error) {
	var shirts []*model.Shirt
	offset := (page - 1) * limit

	query := r.db.Offset(offset).Limit(limit)
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if size != "" {
		query = query.Where("size LIKE ?", "%"+size+"%")
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}

	result := query.Order("created_at DESC").Find(&shirts)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting shirts: %s", result.Error)
	}
	return shirts, nil
}

// GetShirtByIDRepository mengambil shirt berdasarkan ID
func (r *shirtRepository) GetShirtByIdRepository(id string) (*model.Shirt, error) {
	var shirt model.Shirt
	result := r.db.First(&shirt, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("shirt with ID %s not found", id)
		}
		return nil, fmt.Errorf("error getting shirt with ID %s: %s", id, result.Error)
	}
	return &shirt, nil
}

// CreateShirtRepository membuat entri shirt baru di database
func (r *shirtRepository) CreateShirtRepository(shirt *model.Shirt) (*model.Shirt, error) {
	var existingShirt model.Shirt
	if err := r.db.Where("code = ?", shirt.Code).First(&existingShirt).Error; err == nil {
		return nil, fmt.Errorf("shirt with code %s already exists", shirt.Code)
	}

	result := r.db.Create(shirt)
	if result.Error != nil {
		return nil, result.Error
	}
	return shirt, nil
}

// UpdateShirtByIDRepository memperbarui data shirt berdasarkan ID
func (r *shirtRepository) UpdateShirtByIdRepository(id string, shirt *model.Shirt) (*model.Shirt, error) {
	var existingShirt model.Shirt
	if err := r.db.Where("code = ? AND id != ?", shirt.Code, id).First(&existingShirt).Error; err == nil {
		return nil, fmt.Errorf("shirt with code %s already exists", shirt.Code)
	}

	result := r.db.Model(&model.Shirt{}).Where("id = ?", id).Updates(shirt)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("shirt not found")
	}
	return shirt, nil
}

// DeleteShirtByIDRepository menghapus shirt berdasarkan ID
func (r *shirtRepository) DeleteShirtByIdRepository(id string) error {
	result := r.db.Delete(&model.Shirt{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("shirt not found")
	}
	return nil
}
