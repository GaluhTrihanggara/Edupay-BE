package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Cart adalah model untuk keranjang belanja pengguna
type Cart struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"` // ID unik untuk cart
	UserID    uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`                            // Foreign key untuk User
	User      User           `gorm:"foreignKey:UserID" json:"user"`                             // Relasi ke model User
	CartItems []CartItem     `gorm:"foreignKey:CartID" json:"cart_items"`
	CreatedAt time.Time      `json:"created_at"`              // Timestamp pembuatan
	UpdatedAt time.Time      `json:"updated_at"`              // Timestamp pembaruan
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` // Timestamp penghapusan
}
