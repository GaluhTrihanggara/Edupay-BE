package model

import "time"

type BillSemester struct {
	UUIDPrimaryKey
	StudentId string    `gorm:"type:varchar(10)" json:"student_id"`
	Type      string    `gorm:"type:varchar(100)" json:"product_type"`
	Semester  string    `gorm:"type:varchar(100)" json:"semester"`
	Year      string    `gorm:"type:varchar(100)" json:"year"`
	Amount    float64   `gorm:"type:decimal(12)" json:"amount"`
	Status    string    `gorm:"type:varchar(100)" json:"status"`
	DueDate   time.Time `json:"due_date"`
}
