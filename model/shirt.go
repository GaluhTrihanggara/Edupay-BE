package model

type Shirt struct {
	UUIDPrimaryKey
	Code  string    `gorm:"type:varchar(20);unique" json:"code"` // Kode unik untuk baju
	Name  ShirtName `gorm:"type:varchar(100)" json:"name"`
	Size  Size      `gorm:"type:varchar(100)" json:"size"`
	Price float64   `gorm:"type:numeric" json:"price"`
	Stock int       `gorm:"type:int" json:"stock"`
}

type ShirtName string

const (
	Baju_Putih_Merah ShirtName = "Baju_Putih_Merah"
	Baju_Batik       ShirtName = "Baju_Batik"
	Baju_Pramuka     ShirtName = "Baju_Pramuka"
	Baju_Muslim      ShirtName = "Baju_Muslim"
)

type Size string

const (
	Small   Size = "S"
	Medium  Size = "M"
	Large   Size = "L"
	XLarge  Size = "XL"
	XXLarge Size = "XXL"
)
