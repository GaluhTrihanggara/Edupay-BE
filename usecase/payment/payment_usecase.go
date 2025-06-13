package payment

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PaymentUseCase interface {
	CreatePaymentUseCase(payload *model.Payment) (*model.Payment, error)
	GetAllPaymentsUseCase(page, limit int, userID, status string) ([]*model.Payment, error)
	GetPaymentByOrderIDUseCase(orderID uuid.UUID) (*model.Payment, error)
	UpdatePaymentByIDUseCase(paymentID uuid.UUID, payload *model.Payment) (*model.Payment, error)
	DeletePaymentByIDUseCase(paymentID uuid.UUID) error
}

type paymentUseCase struct {
	paymentRepository repository.PaymentRepository
}

func NewPaymentUseCase(paymentRepository repository.PaymentRepository, orderRepository repository.OrderRepository) *paymentUseCase {
	return &paymentUseCase{
		paymentRepository: paymentRepository,
	}
}

// CreatePaymentUseCase creates a new payment
func (uc *paymentUseCase) CreatePaymentUseCase(payload *model.Payment) (*model.Payment, error) {
	// Validasi input
	if payload.OrderID == uuid.Nil {
		return nil, errors.New("order ID cannot be empty")
	}
	if payload.Method == "" {
		return nil, errors.New("payment method is required")
	}
	if payload.StudentID == uuid.Nil {
		return nil, errors.New("student ID is required")
	}
	if payload.Status == "" {
		payload.Status = model.STATUS_UNPAID
	} else if !model.IsValidPaymentStatus(payload.Status) {
		return nil, fmt.Errorf("invalid payment status: %s", payload.Status)
	}
	if payload.PaymentDate.IsZero() {
		payload.PaymentDate = time.Now()
	}

	// Cek apakah payment sudah ada untuk order ini
	existingPayment, err := uc.paymentRepository.GetPaymentByOrderIDRepository(payload.OrderID)
	if err != nil {
		// Jika error adalah "record not found", lanjutkan untuk membuat payment baru
		if !strings.Contains(err.Error(), "record not found") {
			return nil, fmt.Errorf("error checking existing payment: %w", err)
		}
	}
	if existingPayment != nil {
		return nil, fmt.Errorf("payment already exists for this order")
	}

	// Simpan payment ke database
	newPayment, err := uc.paymentRepository.CreatePaymentRepository(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return newPayment, nil
}

func (uc *paymentUseCase) GetAllPaymentsUseCase(page, limit int, userID, status string) ([]*model.Payment, error) {
	if page < 1 || limit < 1 {
		return nil, fmt.Errorf("page and limit must be greater than zero")
	}

	payments, err := uc.paymentRepository.GetAllPaymentsRepository(page, limit, userID, status)
	if err != nil {
		return nil, err
	}

	return payments, nil
}

// GetPaymentByOrderIDUseCase retrieves a payment by OrderID
func (uc *paymentUseCase) GetPaymentByOrderIDUseCase(orderID uuid.UUID) (*model.Payment, error) {
	if orderID == uuid.Nil {
		return nil, errors.New("orderID cannot be empty")
	}
	return uc.paymentRepository.GetPaymentByOrderIDRepository(orderID)
}

// UpdatePaymentByIdUseCase updates a payment by ID
func (uc *paymentUseCase) UpdatePaymentByIDUseCase(paymentID uuid.UUID, payload *model.Payment) (*model.Payment, error) {
	if paymentID == uuid.Nil || payload == nil {
		return nil, errors.New("payment cannot be nil or empty")
	}
	if payload.Status != "" && !model.IsValidPaymentStatus(payload.Status) {
		return nil, fmt.Errorf("invalid payment status: %s", payload.Status)
	}
	return uc.paymentRepository.UpdatePaymentByIDRepository(paymentID, payload)
}

// DeletePaymentByIdUseCase deletes a payment by ID
func (uc *paymentUseCase) DeletePaymentByIDUseCase(paymentID uuid.UUID) error {
	if paymentID == uuid.Nil {
		return errors.New("paymentID cannot be empty")
	}
	return uc.paymentRepository.DeletePaymentByIDRepository(paymentID)
}
