package model

import (
	"time"

	"github.com/google/uuid"
)

// Shirt adalah model untuk data baju (seragam) yang diperlukan oleh siswa
type Shirt struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"` // ID unik untuk baju
	Code      string    `gorm:"type:varchar(100);unique" json:"code"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`   // Nama baju (misalnya seragam)
	Size      string    `gorm:"type:varchar(10)" json:"size"`    // Ukuran baju (misalnya S, M, L)
	Price     float64   `gorm:"type:decimal(12,2)" json:"price"` // Harga baju
	Quantity  int       `gorm:"type:int" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
