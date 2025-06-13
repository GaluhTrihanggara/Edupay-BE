package repository

import (
	"Edupay/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePaymentRepository(payment *model.Payment) (*model.Payment, error)
	GetAllPaymentsRepository(page, limit int, userID, status string) ([]*model.Payment, error)
	GetPaymentByOrderIDRepository(orderID uuid.UUID) (*model.Payment, error)
	UpdatePaymentByIDRepository(PaymentID uuid.UUID, payment *model.Payment) (*model.Payment, error)
	DeletePaymentByIDRepository(PaymentID uuid.UUID) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreatePaymentRepository(payment *model.Payment) (*model.Payment, error) {
	var order model.Order
	// Gunakan Preload untuk memuat relasi Order, Cart, dan CartItems
	if err := r.db.Preload("User").Preload("Cart.User").Preload("Cart.CartItems").Preload("Cart.CartItems.Cart.User").Preload("Cart.CartItems.Shirt").First(&order, "id = ? AND deleted_at IS NULL", payment.OrderID).Error; err != nil {
		return nil, fmt.Errorf("order not found: %v", err)
	}

	// Validasai StudentID jika ada
	if payment.StudentID != uuid.Nil {
		var student model.Student
		if err := r.db.First(&student, "id = ? AND user_id = ? AND deleted_at IS NULL", payment.StudentID, order.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("student not found or not associated with user: %w", err)
			}
			return nil, fmt.Errorf("failed to validate student: %v", err)
		}
		payment.Student = student
	} else {
		return nil, errors.New("student ID is required")
	}

	// Set Amount dari TotalAmount order
	payment.Amount = order.TotalAmount
	payment.AdminFee = model.ADMIN_FEE
	payment.TotalPrice = payment.Amount + payment.AdminFee

	// Pastikan deleted_at adalah null saat membuat payment baru
	payment.DeletedAt = gorm.DeletedAt{}

	// Simpan payment ke database
	if err := r.db.Create(payment).Error; err != nil {
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	// Ambil payment yang baru dibuat dengan relasi yang sudah di-preload
	var createdPayment model.Payment
	if err := r.db.
		Preload("Order").
		Preload("Order.User").
		Preload("Order.Cart.User").
		Preload("Order.Cart.CartItems").
		Preload("Order.Cart.CartItems.Cart.User").
		Preload("Order.Cart.CartItems.Shirt").
		Preload("Student.User").
		First(&createdPayment, "id = ?", payment.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve created payment: %v", err)
	}

	return &createdPayment, nil
}

func (r *paymentRepository) GetAllPaymentsRepository(page, limit int, userID, status string) ([]*model.Payment, error) {
	var payments []*model.Payment
	query := r.db.Model(&model.Payment{})

	// Hanya ambil payment yang belum dihapus (soft delete)
	query = query.Where("deleted_at IS NULL")

	// Filter berdasarkan user_id jika diberikan (melalui join dengan orders)
	if userID != "" {
		query = query.Joins("JOIN orders ON orders.id = payments.order_id").
			Where("orders.user_id = ?", userID)
	}

	// Filter berdasarkan status jika diberikan
	if status != "" {
		query = query.Where("payments.status = ?", status)
	}

	// Pagination dan pengurutan
	err := query.Offset((page - 1) * limit).
		Limit(limit).
		Order("created_at DESC").
		Preload("Order").                // Memuat relasi Order
		Preload("Order.User").           // Memuat relasi User di Order
		Preload("Order.Cart.User").      // Memuat relasi Cart di Order
		Preload("Order.Cart.CartItems"). // Memuat relasi CartItems di Cart
		Preload("Order.Cart.CartItems.Cart.User").
		Preload("Order.Cart.CartItems.Book").  // Memuat relasi Book
		Preload("Order.Cart.CartItems.Shirt"). // Memuat relasi Shirt
		Preload("Student").
		Preload("Student.User").
		Find(&payments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}
	return payments, nil
}

func (r *paymentRepository) GetPaymentByOrderIDRepository(orderID uuid.UUID) (*model.Payment, error) {
	var payment model.Payment
	if err := r.db.Preload("Order").First(&payment, "order_id = ?", orderID).Error; err != nil {
		return nil, fmt.Errorf("failed to get payment: %v", err)
	}
	return &payment, nil
}

func (r *paymentRepository) UpdatePaymentByIDRepository(paymentID uuid.UUID, payment *model.Payment) (*model.Payment, error) {
	var existingPayment model.Payment
	if err := r.db.First(&existingPayment, "id = ? AND deleted_at IS NULL", paymentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payment not found for payment ID: %s", paymentID)
		}
		return nil, fmt.Errorf("failed to get payment: %v", err)
	}

	// Update fields
	if payment.Status != "" {
		existingPayment.Status = payment.Status
	}
	if payment.Method != "" {
		existingPayment.Method = payment.Method
	}
	if !payment.PaymentDate.IsZero() {
		existingPayment.PaymentDate = payment.PaymentDate
	}
	if payment.StudentID != uuid.Nil {
		existingPayment.StudentID = payment.StudentID
	}

	if err := r.db.Save(&existingPayment).Error; err != nil {
		return nil, fmt.Errorf("failed to update payment: %v", err)
	}

	var updatedPayment model.Payment
	if err := r.db.
		Preload("Order").      // Memuat relasi Order
		Preload("Order.User"). // Memuat relasi User di Order
		Preload("Order.Cart"). // Memuat relasi Cart di Order
		Preload("Order.Cart.User").
		Preload("Order.Cart.CartItems.Cart.User").
		Preload("Order.Cart.CartItems").       // Memuat relasi CartItems di Cart
		Preload("Order.Cart.CartItems.Book").  // Memuat relasi Book
		Preload("Order.Cart.CartItems.Shirt"). // Memuat relasi Shirt
		Preload("Student").
		Where("id = ?", paymentID).First(payment).Error; err != nil {
		return nil, fmt.Errorf("failed to get payment: %v", err)
	}
	return &updatedPayment, nil
}

func (r *paymentRepository) DeletePaymentByIDRepository(paymentID uuid.UUID) error {
	if err := r.db.Delete(&model.Payment{}, "id = ?", paymentID).Error; err != nil {
		return fmt.Errorf("failed to delete payment: %v", err)
	}
	return nil
}
