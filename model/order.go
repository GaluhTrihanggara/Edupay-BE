package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order adalah model untuk mencatat pesanan yang dibuat oleh pengguna setelah checkout
type Order struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"` // ID unik untuk order
	UserID      uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`                            // Foreign key untuk User
	User        User           `gorm:"foreignKey:UserID" json:"user"`                             // Relasi ke model User
	CartID      uuid.UUID      `gorm:"type:uuid;index" json:"cart_id"`                            // Foreign key untuk Cart
	Cart        Cart           `gorm:"foreignKey:CartID" json:"cart"`                             // Relasi ke model Cart
	Status      string         `gorm:"type:varchar(50)" json:"status"`                            // Status order (misalnya "pending", "completed")
	TotalAmount float64        `gorm:"type:decimal(12,2)" json:"total_amount"`                    // Total biaya order
	CreatedAt   time.Time      `json:"created_at"`                                                // Timestamp pembuatan
	UpdatedAt   time.Time      `json:"updated_at"`                                                // Timestamp pembaruan
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`                                   // Timestamp penghapusan
}
