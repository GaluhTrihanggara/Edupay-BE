package model

import (
	"time"

	"github.com/google/uuid"
)

// PaymentHistory adalah model untuk mencatat riwayat pembayaran yang telah dilakukan oleh pengguna
type PaymentHistory struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"` // ID unik untuk riwayat pembayaran
	UserID      uuid.UUID `gorm:"type:uuid;index" json:"user_id"`                            // Foreign key untuk User
	User        User      `gorm:"foreignKey:UserID" json:"user"`                             // Relasi ke model User
	PaymentID   uuid.UUID `gorm:"type:uuid;index" json:"payment_id"`                         // Foreign key untuk Payment
	Payment     Payment   `gorm:"foreignKey:PaymentID" json:"payment"`                       // Relasi ke model Payment
	Amount      float64   `gorm:"type:decimal(12,2)" json:"amount"`                          // Jumlah yang dibayar
	PaymentDate time.Time `json:"payment_date"`                                              // Tanggal pembayaran
	Status      string    `gorm:"type:varchar(50)" json:"status"`                            // Status pembayaran (misalnya "pending", "completed")
	CreatedAt   time.Time `json:"created_at"`                                                // Timestamp pembuatan
	UpdatedAt   time.Time `json:"updated_at"`                                                // Timestamp pembaruan
	DeletedAt   time.Time `json:"deleted_at"`                                                // Timestamp penghapusan
}
