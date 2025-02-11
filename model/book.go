package model

type Book struct {
	UUIDPrimaryKey
	Code   string  `gorm:"type:varchar(20);unique" json:"code"` // Kode unik untuk buku
	Title  string  `gorm:"type:varchar(100)" json:"title"`
	Class  string  `gorm:"type:varchar(50)" json:"class"`
	Author string  `gorm:"type:varchar(100)" json:"author"`
	Price  float64 `gorm:"type:numeric" json:"price"`
	Stock  int     `gorm:"type:int" json:"stock"`
}
