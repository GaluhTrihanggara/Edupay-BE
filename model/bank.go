package model

type Bank struct {
	UUIDPrimaryKey
	Name     string `gorm:"type:varchar(100)" json:"name"`
	Image    string `gorm:"type:varchar(255)" json:"image"`
	BankCode string `gorm:"type:varchar(100)" json:"bank_code"`
}
