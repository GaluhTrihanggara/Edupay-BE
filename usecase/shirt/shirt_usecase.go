package shirt

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"
	"fmt"
)

type ShirtUseCase interface {
	CreateShirtUseCase(payload *model.Shirt) (*model.Shirt, error)
	GetAllShirtUseCase(page, limit int, name, size, code string) ([]*model.Shirt, error)
	GetShirtByIdUseCase(shirtId string) (*model.Shirt, error)
	UpdateShirtByIdUseCase(shirtId string, payload *model.Shirt) (*model.Shirt, error)
	DeleteShirtByIdUseCase(shirtId string) error
}

type shirtUseCase struct {
	shirtRepository repository.ShirtRepository
}

func NewShirtUseCase(shirtRepository repository.ShirtRepository) *shirtUseCase {
	return &shirtUseCase{
		shirtRepository: shirtRepository,
	}
}

// CreateShirtUseCase membuat kaos baru
func (uc *shirtUseCase) CreateShirtUseCase(payload *model.Shirt) (*model.Shirt, error) {
	// Validasi unik kode
	shirts, err := uc.shirtRepository.GetAllShirtsRepository(1, 1, "", "", payload.Code)
	if err == nil && len(shirts) > 0 {
		return nil, fmt.Errorf("shirt with code %s already exists", payload.Code)
	}

	// Simpan kaos baru
	shirt, err := uc.shirtRepository.CreateShirtRepository(payload)
	if err != nil {
		return nil, fmt.Errorf("error creating shirt in database: %w", err)
	}
	return shirt, nil
}

// GetAllShirtUseCase mengambil semua kaos dengan filter nama, ukuran, dan kode
func (uc *shirtUseCase) GetAllShirtUseCase(page, limit int, name, size, code string) ([]*model.Shirt, error) {
	shirts, err := uc.shirtRepository.GetAllShirtsRepository(page, limit, name, size, code)
	if err != nil {
		return nil, err
	}
	return shirts, nil
}

// GetShirtByIdUseCase mengambil kaos berdasarkan ID
func (uc *shirtUseCase) GetShirtByIdUseCase(shirtId string) (*model.Shirt, error) {
	shirt, err := uc.shirtRepository.GetShirtByIdRepository(shirtId)
	if err != nil {
		return nil, errors.New("shirt not found")
	}
	return shirt, nil
}

// UpdateShirtByIdUseCase memperbarui data kaos berdasarkan ID
func (uc *shirtUseCase) UpdateShirtByIdUseCase(shirtId string, payload *model.Shirt) (*model.Shirt, error) {
	// Validasi keberadaan kaos
	existingShirt, err := uc.shirtRepository.GetShirtByIdRepository(shirtId)
	if err != nil {
		return nil, fmt.Errorf("shirt with ID %s not found: %v", shirtId, err)
	}

	// Validasi unik kode jika diperbarui
	if payload.Code != existingShirt.Code {
		shirts, err := uc.shirtRepository.GetAllShirtsRepository(1, 1, "", "", payload.Code)
		if err == nil && len(shirts) > 0 {
			return nil, fmt.Errorf("shirt with code %s already exists", payload.Code)
		}
	}

	// Perbarui kaos
	existingShirt.Name = payload.Name
	existingShirt.Size = payload.Size
	existingShirt.Price = payload.Price
	existingShirt.Stock = payload.Stock
	existingShirt.Code = payload.Code

	updatedShirt, err := uc.shirtRepository.UpdateShirtByIdRepository(shirtId, existingShirt)
	if err != nil {
		return nil, fmt.Errorf("failed to update shirt: %w", err)
	}
	return updatedShirt, nil
}

// DeleteShirtByIdUseCase menghapus kaos berdasarkan ID
func (uc *shirtUseCase) DeleteShirtByIdUseCase(shirtId string) error {
	// Validasi keberadaan kaos sebelum dihapus
	_, err := uc.shirtRepository.GetShirtByIdRepository(shirtId)
	if err != nil {
		return fmt.Errorf("shirt with ID %s not found: %v", shirtId, err)
	}

	// Hapus kaos
	err = uc.shirtRepository.DeleteShirtByIdRepository(shirtId)
	if err != nil {
		return fmt.Errorf("failed to delete shirt: %v", err)
	}
	return nil
}
