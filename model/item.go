package model

type Item struct {
	ID          string       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Code        string       `gorm:"type:varchar(20);unique" json:"code"`
	Name        string       `gorm:"type:varchar(100)" json:"name"`
	Description string       `gorm:"type:varchar(100)" json:"description"`
	TotalPrice  float64      `gorm:"type:numeric(10,2);default:0" json:"total_price"`
	UserId      string       `gorm:"column:user_id;type:uuid" json:"user_id"` // ✅ Periksa ini!
	User        User         `gorm:"foreignKey:UserId" json:"user"`
	ItemDetails []ItemDetail `gorm:"foreignKey:ItemID;references:ID" json:"item_details"`
}
