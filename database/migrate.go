package database

import (
	"Edupay/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.BillSemester{},
		&model.Book{},
		&model.Item{},
		&model.ItemDetail{},
		&model.Shirt{},
		&model.Student{},
		&model.Transaction{},
	)
	if err != nil {
		panic(err)
	}
}

func Drop(db *gorm.DB) {
	err := db.Migrator().DropTable(
		&model.User{},
		&model.BillSemester{},
		&model.Book{},
		&model.Item{},
		&model.ItemDetail{},
		&model.Shirt{},
		&model.Student{},
		&model.Transaction{},
	)
	if err != nil {
		panic(err)
	}
}
