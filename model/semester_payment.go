package model

import (
	"time"

	"github.com/google/uuid"
)

// SemesterPayment adalah model untuk pembayaran per semester (6 bulan)
type SemesterPayment struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"` // ID unik untuk pembayaran semester
	Amount      float64   `gorm:"type:decimal(12,2)" json:"amount"`                          // Jumlah yang dibayar
	PaymentDate time.Time `json:"payment_date"`                                              // Tanggal pembayaran
	DueDate     time.Time `json:"due_date"`                                                  // Tanggal jatuh tempo pembayaran
	Status      string    `gorm:"type:varchar(50)" json:"status"`                            // Status pembayaran (misalnya "pending", "paid")
	StudentID   uuid.UUID `gorm:"type:varchar(100);index" json:"student_id"`                 // Foreign key untuk Student
	Student     Student   `gorm:"foreignKey:StudentID" json:"student"`                       // Relasi ke model Student
	CreatedAt   time.Time `json:"created_at"`                                                // Timestamp pembuatan
	UpdatedAt   time.Time `json:"updated_at"`                                                // Timestamp pembaruan
	DeletedAt   time.Time `json:"deleted_at"`                                                // Timestamp penghapusan

}
