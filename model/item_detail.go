package model

import "github.com/google/uuid"

type ItemDetail struct {
	ID        string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ItemID    string  `gorm:"type:uuid;index" json:"item_id"`
	ProductID string  `gorm:"type:uuid;index" json:"product_id"`
	Type      string  `gorm:"type:varchar(50);not null" json:"type"`
	Quantity  int     `gorm:"type:int" json:"quantity"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}

// ✅ **Tambahkan fungsi ini untuk membuat UUID baru jika kosong**
func (itemDetail *ItemDetail) BeforeCreate() error {
	if itemDetail.ID == "" {
		itemDetail.ID = uuid.New().String()
	}
	return nil
}

type Product struct {
	ID    string  `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name  string  `gorm:"type:varchar(100)" json:"name"`
	Price float64 `gorm:"type:decimal(10,2)" json:"price"`
}
