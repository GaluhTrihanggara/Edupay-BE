package shirt

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type ShirtUseCase interface {
	CreateShirtUseCase(payload *model.Shirt) (*model.Shirt, error)
	GetAllShirtUseCase(page, limit int, name, size string) ([]*model.Shirt, error)
	GetShirtByIdUseCase(shirtId uuid.UUID) (*model.Shirt, error)
	GetShirtByCodeUseCase(code string) (*model.Shirt, error)
	UpdateShirtByIdUseCase(shirtId uuid.UUID, payload *model.Shirt) (*model.Shirt, error)
	DeleteShirtByIdUseCase(shirtId uuid.UUID) error
}

type shirtUseCase struct {
	shirtRepository repository.ShirtRepository
}

func NewShirtUseCase(shirtRepository repository.ShirtRepository) *shirtUseCase {
	return &shirtUseCase{
		shirtRepository: shirtRepository,
	}
}

// CreateShirtUseCase creates a new shirt
func (uc *shirtUseCase) CreateShirtUseCase(payload *model.Shirt) (*model.Shirt, error) {
	// Validate unique code
	shirt, err := uc.shirtRepository.GetAllShirtsRepository(1, 1, payload.Name, payload.Size)
	if err == nil && len(shirt) > 0 {
		return nil, fmt.Errorf("shirt with Size %s already exists", payload.Size)
	}

	// Save new shirt
	createdShirt, err := uc.shirtRepository.CreateShirtRepository(payload)
	if err != nil {
		return nil, fmt.Errorf("error creating shirt in database: %w", err)
	}
	return createdShirt, nil
}

// GetAllShirtUseCase retrieves all shirts with filters
func (uc *shirtUseCase) GetAllShirtUseCase(page, limit int, name, size string) ([]*model.Shirt, error) {
	shirts, err := uc.shirtRepository.GetAllShirtsRepository(page, limit, name, size)
	if err != nil {
		return nil, err
	}
	return shirts, nil
}

// GetShirtByIdUseCase retrieves a shirt by ID
func (uc *shirtUseCase) GetShirtByIdUseCase(shirtId uuid.UUID) (*model.Shirt, error) {
	shirt, err := uc.shirtRepository.GetShirtByIDRepository(shirtId)
	if err != nil {
		return nil, errors.New("shirt not found")
	}
	return shirt, nil
}

func (uc *shirtUseCase) GetShirtByCodeUseCase(code string) (*model.Shirt, error) {
	shirtCode, err := uc.shirtRepository.GetShirtByCodeRepository(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get book by code %s: %w", code, err)
	}
	if shirtCode == nil {
		return nil, fmt.Errorf("shirt with code %s not found", code)
	}
	return shirtCode, nil
}

// UpdateShirtByIdUseCase updates shirt information by ID
func (uc *shirtUseCase) UpdateShirtByIdUseCase(shirtId uuid.UUID, payload *model.Shirt) (*model.Shirt, error) {
	// Validate shirt existence
	existingShirt, err := uc.shirtRepository.GetShirtByIDRepository(shirtId)
	if err != nil {
		return nil, fmt.Errorf("shirt with Id %s not found: %v", shirtId, err)
	}

	// Validate unique code if it's updated
	if payload.Code != existingShirt.Code {
		shirts, err := uc.shirtRepository.GetAllShirtsRepository(1, 1, payload.Name, payload.Size)
		if err == nil && len(shirts) > 0 {
			return nil, fmt.Errorf("shirt with code %s already exists", payload.Code)
		}
	}

	// Update shirt fields
	existingShirt.Code = payload.Code
	existingShirt.Name = payload.Name
	existingShirt.Size = payload.Size
	existingShirt.Price = payload.Price
	existingShirt.Quantity = payload.Quantity

	updatedShirt, err := uc.shirtRepository.UpdateShirtByIDRepository(shirtId, existingShirt)
	if err != nil {
		return nil, fmt.Errorf("failed to update shirt: %w", err)
	}
	return updatedShirt, nil
}

// DeleteShirtByIdUseCase deletes a shirt by ID
func (uc *shirtUseCase) DeleteShirtByIdUseCase(shirtId uuid.UUID) error {
	// Validate shirt existence before deletion
	_, err := uc.shirtRepository.GetShirtByIDRepository(shirtId)
	if err != nil {
		return fmt.Errorf("shirt with ID %s not found: %v", shirtId, err)
	}

	// Delete shirt
	err = uc.shirtRepository.DeleteShirtByIDRepository(shirtId)
	if err != nil {
		return fmt.Errorf("failed to delete shirt: %v", err)
	}
	return nil
}
