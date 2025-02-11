package model

import "time"

type TransactionHistory struct {
	UUIDPrimaryKey
	UserId        string       `gorm:"type:varchar(100)" json:"user_id"`            // ID user yang melakukan transaksi
	TransactionId string       `gorm:"type:varchar(100)" json:"transaction_id"`     // ID transaksi utama
	Transaction   *Transaction `gorm:"foreignKey:TransactionId" json:"transaction"` // Relasi ke transaksi utama
	Amount        float64      `gorm:"type:numeric" json:"amount"`                  // Jumlah yang dibayarkan
	Status        string       `gorm:"type:varchar(100)" json:"status"`             // Status transaksi (e.g., successful, failed)
	PaymentMethod string       `gorm:"type:varchar(100)" json:"payment_method"`     // Metode pembayaran
	PaymentDate   time.Time    `json:"payment_date"`                                // Tanggal pembayaran
	CreatedAt     time.Time    `json:"created_at"`                                  // Tanggal dibuat
	UpdatedAt     time.Time    `json:"updated_at"`                                  // Tanggal diperbarui
}
