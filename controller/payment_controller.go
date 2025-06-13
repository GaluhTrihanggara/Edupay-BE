package controller

import (
	"Edupay/model"
	"Edupay/usecase/payment"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PaymentController interface {
	CreatePaymentController(c echo.Context) error
	GetAllPaymentsController(c echo.Context) error
	GetPaymentByOrderIDController(c echo.Context) error
	UpdatePaymentByIDController(c echo.Context) error
	DeletePaymentByIDController(c echo.Context) error
}

type paymentController struct {
	paymentUseCase payment.PaymentUseCase
}

func NewPaymentController(paymentUseCase payment.PaymentUseCase) *paymentController {
	return &paymentController{
		paymentUseCase: paymentUseCase,
	}
}

// isValidPaymentStatus validates the payment status
func isValidPaymentStatus(status string) bool {
	validStatuses := []string{model.STATUS_SUCCESSFUL, model.STATUS_PROCESSING, model.STATUS_FAIL, model.STATUS_UNPAID}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// CreatePaymentController creates a new payment
func (ctrl *paymentController) CreatePaymentController(c echo.Context) error {
	var payload model.Payment
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload" + err.Error(),
		})
	}

	// Validate required fields
	if payload.OrderID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Order ID is required",
		})
	}
	if payload.Method == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Payment method is required",
		})
	}
	if payload.StudentID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Payment method is required",
		})
	}
	if payload.Status == "" {
		payload.Status = model.STATUS_UNPAID
	} else if !isValidPaymentStatus(payload.Status) {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid payment status: " + payload.Status,
		})
	}
	if payload.PaymentDate.IsZero() {
		payload.PaymentDate = time.Now()
	}

	// Create payment
	response, err := ctrl.paymentUseCase.CreatePaymentUseCase(&payload) // CreatePaymentUseCase now only returns an error
	if err != nil {
		if strings.Contains(err.Error(), "payment already exists") || strings.Contains(err.Error(), "order not found") {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully created payment",
		},
		Data: response,
	})
}

func (ctrl *paymentController) GetAllPaymentsController(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	userID := c.QueryParam("user_id")
	status := c.QueryParam("status")

	payments, err := ctrl.paymentUseCase.GetAllPaymentsUseCase(page, limit, userID, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Payments retrieved successfully",
		},
		Data: payments,
	})
}

// GetPaymentByOrderIDController retrieves a payment by Order ID
func (ctrl *paymentController) GetPaymentByOrderIDController(c echo.Context) error {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Order ID" + err.Error(),
		})
	}

	response, err := ctrl.paymentUseCase.GetPaymentByOrderIDUseCase(orderID)
	if err != nil {
		if strings.Contains(err.Error(), "payment not found") {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve payment: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved payment",
		},
		Data: response,
	})
}

// UpdatePaymentByIDController updates a payment by ID
func (ctrl *paymentController) UpdatePaymentByIDController(c echo.Context) error {
	paymentIDStr := c.Param("id")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Payment ID" + err.Error(),
		})
	}

	var payload model.Payment
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	//Validate status if provided
	if payload.Status != "" && !model.IsValidPaymentStatus(payload.Status) {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid payment status" + payload.Status,
		})
	}

	// Update payment
	response, err := ctrl.paymentUseCase.UpdatePaymentByIDUseCase(paymentID, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "payment not found") {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update payment:" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully updated payment",
		},
		Data: response,
	})
}

// DeletePaymentByIDController deletes a payment by ID
func (ctrl *paymentController) DeletePaymentByIDController(c echo.Context) error {
	paymentIDStr := c.Param("id")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Payment ID" + err.Error(),
		})
	}

	err = ctrl.paymentUseCase.DeletePaymentByIDUseCase(paymentID)
	if err != nil {
		if strings.Contains(err.Error(), "payment not found") {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete payment:" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully deleted payment",
		},
	})
}
