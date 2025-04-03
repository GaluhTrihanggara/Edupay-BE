package model

import (
	"time"

	"gorm.io/gorm"
)

const STATUS_UNPAID = "unpaid"
const STATUS_FAIL = "fail"
const STATUS_PROCESSING = "processing"
const STATUS_SUCCESSFUL = "successful"
const ADMIN_FEE = 2500

type Transaction struct {
	Id              string         `gorm:"primaryKey" json:"id"`
	UserId          string         `gorm:"type:varchar(100)" json:"user_id"`
	Status          string         `gorm:"type:varchar(100)" json:"status"`
	ProductType     string         `gorm:"type:varchar(50)" json:"product_type"`
	ProductDetail   interface{}    `gorm:"serializer:json" json:"product_detail"`
	Description     string         `gorm:"type:text" json:"description"`
	AdminFee        float64        `gorm:"type:decimal(12)" json:"admin_fee"`
	Price           float64        `gorm:"type:decimal(12)" json:"price"`
	TotalPrice      float64        `gorm:"type:decimal(12)" json:"total_price"`
	TransactionDate time.Time      `json:"transaction_date"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at"`
}

type GetProductDetail struct {
	Status     string `json:"status"`
	CustomerId string `json:"customer_id"`
	ProductId  string `json:"product_id"`
	Period     string `json:"period"`
}

type PayloadMail struct {
	// Bill Semester
	StudentId string `gorm:"type:varchar(10)" json:"student_id"`
	Semester  string `gorm:"type:varchar(100)" json:"semester"`
	Year      string `gorm:"type:varchar(100)" json:"year"`
	// Book
	Title  string `gorm:"type:varchar(100)" json:"title"`
	Class  string `gorm:"type:varchar(50)" json:"class"`
	Author string `gorm:"type:varchar(100)" json:"author"`
	// Shirt
	Name ShirtName `gorm:"type:varchar(100)" json:"shirt_name"`
	Size Size      `gorm:"type:varchar(100)" json:"size"`
	// =========

	CustomerName  string    `json:"name"`
	OrderId       string    `json:"order_id"`
	CustomerId    string    `json:"costumer_id"`
	ProductType   string    `json:"product_type"`
	Status        string    `json:"status"`
	RecipentEmail string    `json:"recipent_email"`
	ProviderName  string    `json:"provider_name"`
	Period        string    `json:"period"`
	Subject       string    `json:"subject"`
	TransactionAt time.Time `json:"transaction_at"`
	Description   string    `gorm:"type:text" json:"description"`
	DiscountPrice float64   `json:"discount_price"`
	AdminFee      float64   `json:"admin_fee"`
	Price         float64   `json:"price"`
	TotalPrice    float64   `json:"total_price"`
}
