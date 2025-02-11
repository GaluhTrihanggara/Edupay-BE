package model

import "time"

type BillSemester struct {
	UUIDPrimaryKey              //UUIDPrimaryKey base struct
	StudentId      string       `gorm:"type:varchar(10)" json:"student_id"` // Relation to siswa
	Semester       string       `gorm:"type:varchar(100)" json:"semester"`
	Year           string       `gorm:"type:varchar(100)" json:"year"`
	Amount         float64      `gorm:"type:decimal(12)" json:"amount"`
	Status         string       `gorm:"type:varchar(100)" json:"status"`
	DueDate        time.Time    `json:"due_date"`
	TransactionId  *string      `gorm:"type:varchar(100)" json:"transaction_id"` // Optional: Relation ke Transaction
	Transaction    *Transaction `gorm:"foreignKey:TransactionId" json:"transaction"`
}
