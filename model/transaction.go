package model

import (
	"time"

	"gorm.io/gorm"
)

const STATUS_UNPAID = "unpaid"
const STATUS_FAIL = "fail"
const STATUS_PROCESSING = "processing"
const STATUS_SUCCESSFUL = "successful"

type Transaction struct {
	Id              string         `gorm:"primaryKey" json:"id"`
	UserId          string         `gorm:"type:varchar(100)" json:"user_id"`
	Items           []Item         `gorm:"many2many:transaction_items" json:"items"`  // Relasi ke Item melalui tabel many-to-many
	BillSemesterId  *string        `gorm:"type:varchar(100)" json:"bill_semester_id"` // Optional: Relasi ke tagihan semester
	BillSemester    *BillSemester  `gorm:"foreignKey:BillSemesterId" json:"bill_semester"`
	Description     string         `gorm:"type:text" json:"description"`
	Quantity        string         `gorm:"type:varchar(100)" json:"quantity"`
	Status          string         `gorm:"type:varchar(100)" json:"status"`
	TotalPrice      float64        `gorm:"type:decimal(12)" json:"total_price"`
	TransactionDate time.Time      `json:"transaction_date"`
	Method          string         `gorm:"type:varchar(100)" json:"method"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at"`
}
