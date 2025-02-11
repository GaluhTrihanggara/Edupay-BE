package model

type Student struct {
	UUIDPrimaryKey
	Name     string `gorm:"type:varchar(100)" json:"name"`
	Class    string `gorm:"type:varchar(50)" json:"class"`
	ParentId string `gorm:"type:varchar(50);column:parent_id" json:"parent_id"`
	Parent   User   `gorm:"foreignKey:ParentId;references:ID" json:"parent"`
}
