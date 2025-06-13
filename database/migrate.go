package database

import (
	"Edupay/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.SemesterPayment{},
		&model.Book{},
		&model.Cart{},
		&model.CartItem{},
		&model.Shirt{},
		&model.Student{},
		&model.Balance{},
		&model.Bank{},
		&model.VaNumber{},
		&model.Order{},
		&model.Payment{},
		&model.PaymentHistory{},
	)
	if err != nil {
		panic(err)
	}
}

func Drop(db *gorm.DB) {
	err := db.Migrator().DropTable(
		&model.User{},
		&model.SemesterPayment{},
		&model.Book{},
		&model.Cart{},
		&model.CartItem{},
		&model.Shirt{},
		&model.Student{},
		&model.Balance{},
		&model.Bank{},
		&model.VaNumber{},
		&model.Order{},
		&model.Payment{},
		&model.PaymentHistory{},
	)
	if err != nil {
		panic(err)
	}
}
