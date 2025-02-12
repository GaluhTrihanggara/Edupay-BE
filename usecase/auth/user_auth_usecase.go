package auth

import (
	"Edupay/model"
	"Edupay/repository"
	"Edupay/usecase/middlewares"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	RegisterUseCase(payload model.User) (*model.AuthResponse, error)
	LoginUseCase(payload model.User) (*model.AuthResponse, string, error)
	RegisterAdminUseCase(payload model.User) (*model.AuthResponse, error)
}

type authUseCase struct {
	authRepository repository.AuthRepository
}

func NewAuthUsecase(authRepository repository.AuthRepository) *authUseCase {
	return &authUseCase{authRepository: authRepository}
}

func (s *authUseCase) RegisterUseCase(payload model.User) (*model.AuthResponse, error) {
	Payload := model.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
		UserType: model.USER_TYPE,
		Phone:    payload.Phone,
	}

	if payload.Name == "" || payload.Email == "" {
		return nil, errors.New("name and email are required fields")
	}

	hashedPassword, _ := HashPassword(Payload.Password)

	newUserModel := model.User{
		Name:     Payload.Name,
		Email:    Payload.Email,
		Password: hashedPassword,
		UserType: Payload.UserType, // Menggunakan value dari lowercasePayload
		Phone:    Payload.Phone,
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

func (s *authUseCase) LoginUseCase(payload model.User) (*model.AuthResponse, string, error) {

	PayloadEmail := payload.Email

	if payload.Email == "" {
		return nil, "", errors.New("email field is required")
	}

	if payload.Password == "" {
		return nil, "", errors.New("password field is required")
	}

	user, err := s.authRepository.GetUserByEmailRepository(PayloadEmail)
	if err != nil {
		return nil, "", err
	}

	if !ComparePasswords(user.Password, payload.Password) {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := middlewares.CreateToken(*user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create token: %v", err)
	}

	resp := &model.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}

	// Cek apakah user adalah admin atau user biasa
	var userTypeMessage string
	if user.UserType == model.ADMIN_TYPE {
		userTypeMessage = "Admin login successful"
	} else {
		userTypeMessage = "User login successful"
	}

	return resp, userTypeMessage, nil
}

func (s *authUseCase) RegisterAdminUseCase(payload model.User) (*model.AuthResponse, error) {
	Payload := model.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
		UserType: model.ADMIN_TYPE, // Untuk admin, tetap gunakan ADMIN_TYPE
		Phone:    payload.Phone,
	}

	if payload.Name == "" || payload.Email == "" {
		return nil, errors.New("name and email are required fields")
	}

	hashedPassword, _ := HashPassword(Payload.Password)

	newUserModel := model.User{
		Name:     Payload.Name,
		Email:    Payload.Email,
		Password: hashedPassword,
		UserType: Payload.UserType, // Menggunakan value dari lowercasePayload
		Phone:    Payload.Phone,
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
