package model

type VaNumber struct {
	UUIDPrimaryKey
	UserId         string  `gorm:"index" json:"user_id"`
	VaNumber       string  `json:"va_number"`
	VaStatus       string  `json:"va_status"`
	BankCode       string  `json:"bank_code"`
	Amount         float64 `json:"amount"`
	ExpirationTime int     `json:"expiration_time"`
	Name           string  `json:"name"`
}
