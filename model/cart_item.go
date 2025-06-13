package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CartItem adalah model untuk item dalam keranjang belanja
type CartItem struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"` // ID unik untuk item
	CartID    uuid.UUID      `gorm:"type:uuid;index" json:"cart_id"`                            // Foreign key untuk Cart
	Cart      Cart           `gorm:"foreignKey:CartID" json:"cart"`                             // Relasi ke model Cart
	ItemType  string         `gorm:"type:varchar(50)" json:"item_type"`                         // Jenis item: "shirt" atau "book"
	ItemID    uuid.UUID      `gorm:"type:uuid;index" json:"item_id"`                            // Foreign key untuk item (Baju atau Buku)
	Quantity  int            `gorm:"type:int" json:"quantity"`                                  // Jumlah item yang dipilih
	Price     float64        `gorm:"type:decimal(12,2)" json:"price"`                           // Harga per item
	Shirt     *Shirt         `gorm:"foreignKey:ItemID;references:ID;constraint:OnDelete:CASCADE" json:"shirt,omitempty"`
	Book      *Book          `gorm:"foreignKey:ItemID;references:ID;constraint:OnDelete:CASCADE" json:"book,omitempty"`
	CreatedAt time.Time      `json:"created_at"`              // Timestamp pembuatan
	UpdatedAt time.Time      `json:"updated_at"`              // Timestamp pembaruan
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` // Timestamp penghapusan
}
