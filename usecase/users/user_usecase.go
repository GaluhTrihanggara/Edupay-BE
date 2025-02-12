package user

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"
	"fmt"
	"strings"
	"time"
)

type UserUseCase interface {
	GetAllUsersUseCase(page, limit int, name string) ([]*model.UserResponse, error)
	GetUserByIdUseCase(userId string) (*model.UserResponse, error)
	GetUserByPhoneUseCase(phone string) (*model.UserResponse, error)
	GetUserByEmailUseCase(email string) (*model.UserResponse, error)
	GetUserByQueryUseCase(query string, page, limit int) ([]*model.User, error)
	UpdateUserByIdUseCase(userId string, payload *model.User) (*model.UserResponse, error)
	DeleteUserByIdUseCase(userId string) error
}

type userUseCase struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) *userUseCase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

func (uc *userUseCase) GetAllUsersUseCase(page, limit int, name string) ([]*model.UserResponse, error) {
	users, err := uc.userRepository.GetAllUsersRepository(page, limit, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %v", err)
	}
	resp := make([]*model.UserResponse, 0, len(users))
	for _, user := range users {
		resp = append(resp, &model.UserResponse{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			UserType:  user.UserType,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	return resp, nil
}

func (uc *userUseCase) GetUserByIdUseCase(userId string) (*model.UserResponse, error) {
	user, err := uc.userRepository.GetUserByIdRepository(userId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	resp := &model.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return resp, nil
}

func (uc *userUseCase) GetUserByPhoneUseCase(phone string) (*model.UserResponse, error) {
	user, err := uc.userRepository.GetUserByPhoneRepository(phone)
	if err != nil {
		return nil, errors.New("user not found")
	}
	resp := &model.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return resp, nil
}

func (uc *userUseCase) GetUserByEmailUseCase(email string) (*model.UserResponse, error) {
	user, err := uc.userRepository.GetUserByEmailRepository(email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	resp := &model.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return resp, nil
}

func (uc *userUseCase) GetUserByQueryUseCase(query string, page, limit int) ([]*model.User, error) {
	users, err := uc.userRepository.GetUserByQueryRepository(query, page, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (uc *userUseCase) UpdateUserByIdUseCase(userId string, payload *model.User) (*model.UserResponse, error) {
	lowercasePayload := &model.User{
		Name:  strings.ToLower(payload.Name),
		Email: strings.ToLower(payload.Email),
		Phone: strings.ToLower(payload.Phone),
	}

	user, err := uc.userRepository.GetUserByIdRepository(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}
	user.Name = lowercasePayload.Name
	user.Email = lowercasePayload.Email
	user.Phone = lowercasePayload.Phone
	user.UpdatedAt = time.Now()

	updatedUser, err := uc.userRepository.UpdateUserByIdRepository(userId, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	resp := &model.UserResponse{
		Id:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		Phone:     updatedUser.Phone,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	return resp, nil
}

func (uc *userUseCase) DeleteUserByIdUseCase(userId string) error {
	err := uc.userRepository.DeleteUserByIdRepository(userId)
	if err != nil {
		return errors.New("user not found")
	}
	return err
}
