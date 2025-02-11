package auth

import (
	"Edupay/model"
	"Edupay/repository"
	"Edupay/usecase/middlewares"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	RegisterUseCase(payload model.User) (*model.AuthResponse, error)
	LoginUseCase(payload model.User) (*model.AuthResponse, error)
	RegisterAdminUseCase(payload model.User) (*model.AuthResponse, error)
}

type authUseCase struct {
	authRepository repository.AuthRepository
}

func NewAuthUsecase(authRepository repository.AuthRepository) *authUseCase {
	return &authUseCase{authRepository: authRepository}
}

func (s *authUseCase) RegisterUseCase(payload model.User) (*model.AuthResponse, error) {
	lowercasePayload := model.User{
		Name:     strings.ToLower(payload.Name),
		Email:    strings.ToLower(payload.Email),
		Password: payload.Password,
		UserType: model.USER_TYPE,
		Phone:    payload.Phone,
	}

	if payload.Name == "" || payload.Email == "" {
		return nil, errors.New("name and email are required fields")
	}

	hashedPassword, _ := HashPassword(lowercasePayload.Password)

	newUserModel := model.User{
		Name:     lowercasePayload.Name,
		Email:    lowercasePayload.Email,
		Password: hashedPassword,
		UserType: lowercasePayload.UserType, // Menggunakan value dari lowercasePayload
		Phone:    lowercasePayload.Phone,
	}
	user, err := s.authRepository.RegisterRepository(newUserModel)
	if err != nil {
		return nil, fmt.Errorf("error creating user in database: %w", err)
	}

	token, err := middlewares.CreateToken(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %v", err)
	}

	resp := &model.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}

	return resp, nil
}

func (s *authUseCase) LoginUseCase(payload model.User) (*model.AuthResponse, error) {

	lowercasePayloadEmail := strings.ToLower(payload.Email)

	if payload.Email == "" {
		return nil, errors.New("email field is required")
	}

	if payload.Password == "" {
		return nil, errors.New("password field is required")

	}
	user, err := s.authRepository.GetUserByEmailRepository(lowercasePayloadEmail)
	if err != nil {
		return nil, err
	}

	if !ComparePasswords(user.Password, payload.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := middlewares.CreateToken(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %v", err)
	}

	resp := &model.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}

	return resp, nil
}

func (s *authUseCase) RegisterAdminUseCase(payload model.User) (*model.AuthResponse, error) {
	lowercasePayload := model.User{
		Name:     strings.ToLower(payload.Name),
		Email:    strings.ToLower(payload.Email),
		Password: payload.Password,
		UserType: model.ADMIN_TYPE, // Untuk admin, tetap gunakan ADMIN_TYPE
		Phone:    payload.Phone,
	}

	if payload.Name == "" || payload.Email == "" {
		return nil, errors.New("name and email are required fields")
	}

	hashedPassword, _ := HashPassword(lowercasePayload.Password)

	newUserModel := model.User{
		Name:     lowercasePayload.Name,
		Email:    lowercasePayload.Email,
		Password: hashedPassword,
		UserType: lowercasePayload.UserType, // Menggunakan value dari lowercasePayload
		Phone:    lowercasePayload.Phone,
	}

	user, err := s.authRepository.RegisterRepository(newUserModel)
	if err != nil {
		return nil, fmt.Errorf("error creating Admin in database: %w", err)
	}

	token, err := middlewares.CreateToken(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %v", err)
	}

	resp := &model.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}

	return resp, nil
}

func ComparePasswords(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
