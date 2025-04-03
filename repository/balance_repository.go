package repository

import (
	"errors"

	"Edupay/model"

	"gorm.io/gorm"
)

type BalanceRepository interface {
	InsertVaRepository(payload *model.VaNumber) (*model.VaNumber, error)
	GetVaNumber(vaNumber string) (*model.VaNumber, error)
	GetVaIdRepository(partnerId string) (*model.VaNumber, error)
	GetVaNumberByUserId(userId string) (*model.VaNumber, error)
}

type balanceRepository struct {
	db *gorm.DB
}

func NewBalanceRepository(db *gorm.DB) BalanceRepository {
	return &balanceRepository{db: db}
}

func (r *balanceRepository) InsertVaRepository(payload *model.VaNumber) (*model.VaNumber, error) {
	var existingVa model.VaNumber
	result := r.db.Where("user_id = ?", payload.UserId).First(&existingVa)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	} else {
		return &existingVa, nil
	}

	result = r.db.Create(payload)
	if result.Error != nil {
		return nil, errors.New("failed to create Va in the database")
	}

	return payload, nil
}

func (r *balanceRepository) GetVaNumber(vaNumber string) (*model.VaNumber, error) {
	var Va *model.VaNumber
	if err := r.db.Where("va_number = ?", vaNumber).Find(&Va).Error; err != nil {
		return nil, err
	}
	return Va, nil
}
func (r *balanceRepository) GetVaNumberByUserId(userId string) (*model.VaNumber, error) {
	var Va *model.VaNumber
	if err := r.db.Where("user_id = ?", userId).Find(&Va).Error; err != nil {
		return nil, err
	}
	return Va, nil
}

func (r *balanceRepository) GetVaIdRepository(partnerId string) (*model.VaNumber, error) {
	var VaTransaction *model.VaNumber
	if err := r.db.Where("id = ?", partnerId).Find(&VaTransaction).Error; err != nil {
		return nil, err
	}
	return VaTransaction, nil
}
