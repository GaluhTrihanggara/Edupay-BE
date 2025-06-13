package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const STATUS_SUCCESSFUL = "successful"
const STATUS_PROCESSING = "processing"
const STATUS_FAIL = "fail"
const STATUS_UNPAID = "unpaid"
const ADMIN_FEE = 2500

// Payment adalah model untuk mencatat pembayaran untuk sebuah order
type Payment struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"` // ID unik untuk pembayaran
	OrderID     uuid.UUID      `gorm:"type:uuid;index" json:"order_id"`                           // Foreign key untuk Order
	Order       Order          `gorm:"foreignKey:OrderID" json:"order"`                           // Relasi ke model Order
	StudentID   uuid.UUID      `gorm:"type:uuid;index" json:"student_id"`
	Student     Student        `gorm:"foreignKey:StudentID" json:"student"`
	Method      string         `gorm:"type:varchar(100)" json:"method"` // Metode pembayaran (misalnya "OY API", "Transfer Bank")
	Status      string         `gorm:"type:varchar(50)" json:"status"`  // Status pembayaran (misalnya "pending", "paid")
	AdminFee    float64        `gorm:"type:decimal(12)" json:"admin_fee"`
	Amount      float64        `gorm:"type:decimal(12,2)" json:"amount"` // Jumlah yang dibayar
	TotalPrice  float64        `gorm:"type:decimal(12)" json:"total_price"`
	PaymentDate time.Time      `json:"payment_date"` // Tanggal pembayaran
	CreatedAt   time.Time      `json:"created_at"`   // Timestamp pembuatan
	UpdatedAt   time.Time      `json:"updated_at"`   // Timestamp pembaruan
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`   // Timestamp penghapusan
}

// BeforeCreate hook untuk menghitung TotalPrice
func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	p.TotalPrice = p.Amount + p.AdminFee
	return nil
}

// IsValidPaymentStatus memvalidasi status pembayaran
func IsValidPaymentStatus(status string) bool {
	validStatuses := []string{STATUS_SUCCESSFUL, STATUS_PROCESSING, STATUS_FAIL, STATUS_UNPAID}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}
