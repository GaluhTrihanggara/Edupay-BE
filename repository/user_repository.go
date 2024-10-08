package repository

import (
	"Edupay/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetAllUsersRepository(page, limit int, name string) ([]*model.User, error)
	GetUserByIDRepository(id string) (*model.User, error)
	GetUserByPhoneRepository(phone string) (*model.User, error)
	GetUserByEmailRepository(email string) (*model.User, error)
	UpdateUserByIDRepository(id string, user *model.User) (*model.User, error)
	DeleteUserByIDRepository(id string) error
	GetUserByQueryRepository(query string, page, limit int) ([]*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAllUsersRepository(page, limit int, name string) ([]*model.User, error) {
	var users []*model.User
	offset := (page - 1) * limit

	query := r.db.Offset(offset).Limit(limit)
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	result := query.Order("created_at DESC").Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting users: %s", result.Error)
	}
	return users, nil
}

// GetUserByIDRepository mengambil pengguna berdasarkan ID
func (r *userRepository) GetUserByIDRepository(id string) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %s not found", id)
		}
		return nil, fmt.Errorf("error getting user with ID %s: %s", id, result.Error)
	}
	return &user, nil
}

// UpdateUserByIDRepository memperbarui data pengguna berdasarkan ID
func (r *userRepository) UpdateUserByIDRepository(id string, user *model.User) (*model.User, error) {
	result := r.db.Model(&model.User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// DeleteUserByIDRepository menghapus pengguna berdasarkan ID
func (r *userRepository) DeleteUserByIDRepository(id string) error {
	result := r.db.Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// GetUserByPhone mengambil pengguna berdasarkan nomor telepon
func (r *userRepository) GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	result := r.db.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GetUserByEmail mengambil pengguna berdasarkan email
func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GetUserByQueryRepository mencari pengguna berdasarkan query (ID, nama, email, telepon, atau alamat)
func (r *userRepository) GetUserByQueryRepository(query string, page, limit int) ([]*model.User, error) {
	var users []*model.User

	offset := (page - 1) * limit

	queryString := "%" + query + "%"
	dbQuery := r.db.Where("id LIKE ? OR name LIKE ? OR email LIKE ? OR phone LIKE ? OR address LIKE ?  ", queryString, queryString, queryString, queryString, queryString).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users)

	if dbQuery.Error != nil {
		return nil, dbQuery.Error
	}

	return users, nil
}
