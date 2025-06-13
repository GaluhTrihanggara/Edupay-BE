package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const USER_TYPE = "user"
const ADMIN_TYPE = "admin"
const ALL_TYPE = "all"

type UUIDPrimaryKey struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
}

func (u *UUIDPrimaryKey) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type User struct {
	UUIDPrimaryKey
	Name      string  `gorm:"type:varchar(100)" json:"name"`
	Email     string  `gorm:"type:varchar(100)" json:"email"`
	Amount    float64 `gorm:"type:decimal(12)" json:"amount"`
	Password  string  `gorm:"type:varchar(100)" json:"password"`
	Phone     string  `gorm:"type:varchar(100)" json:"phone"`
	UserType  string  `gorm:"type:varchar(100)" json:"user_type"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
