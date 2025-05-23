package model

import (
	"time"

	"github.com/google/uuid"
)

// Book adalah model untuk data buku yang diperlukan oleh siswa
type Book struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"` // ID unik untuk buku
	Code      string    `gorm:"type:varchar(10);uniqueIndex" json:"code"`
	Title     string    `gorm:"type:varchar(255)" json:"title"` // Judul buku
	Class     string    `gorm:"type:varchar(255)" json:"class"`
	Author    string    `gorm:"type:varchar(255)" json:"author"` // Penulis buku
	Price     float64   `gorm:"type:decimal(12,2)" json:"price"` // Harga buku
	Quantity  int       `gorm:"type:int" json:"quantity"`        // Jumlah stok buku yang tersedia
	CreatedAt time.Time `json:"created_at"`                      // Timestamp pembuatan
	UpdatedAt time.Time `json:"updated_at"`                      // Timestamp pembaruan
	DeletedAt time.Time `json:"deleted_at"`                      // Timestamp penghapusan
}
