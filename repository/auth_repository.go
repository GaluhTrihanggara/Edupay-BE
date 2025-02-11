package repository

import (
	"Edupay/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type AuthRepository interface {
	RegisterRepository(user model.User) (*model.User, error)
	LoginRepository(user *model.User) (*model.User, error)
	GetUserByEmailRepository(email string) (*model.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{db}
}

func (r *authRepository) RegisterRepository(user model.User) (*model.User, error) {
	var count int64
	r.db.Model(&model.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return nil, errors.New("email already exists")
	}
	err := r.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) LoginRepository(user *model.User) (*model.User, error) {
	result := r.db.Where("email = ? AND password = ?", user.Email, user.Password).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("incorrect email or password")
		}
		return nil, errors.New("failed to get user")
	}
	return user, nil
}

func (r *authRepository) GetUserByEmailRepository(email string) (*model.User, error) {
	var user model.User

	// Debug logging
	fmt.Printf("Attempting to find user with email: %s\n", email)

	// Modified query with debugging
	result := r.db.Debug().
		Where("LOWER(email) = LOWER(?)", email).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	// Debug logging
	fmt.Printf("Found user: ID=%s, Email=%s, UserType=%s\n",
		user.ID, user.Email, user.UserType)

	return &user, nil
}
