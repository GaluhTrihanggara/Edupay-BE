package controller

import (
	"Edupay/model"
	"Edupay/usecase/payment_history"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PaymentHistoryController interface {
	CreatePaymentHistoryController(c echo.Context) error
	GetPaymentHistoryByUserIDController(c echo.Context) error
}

type paymentHistoryController struct {
	paymentHistoryUseCase payment_history.PaymentHistoryUseCase
}

func NewPaymentHistoryController(paymentHistoryUseCase payment_history.PaymentHistoryUseCase) *paymentHistoryController {
	return &paymentHistoryController{
		paymentHistoryUseCase: paymentHistoryUseCase,
	}
}

// CreatePaymentHistoryController creates a new payment history record
func (ctrl *paymentHistoryController) CreatePaymentHistoryController(c echo.Context) error {
	var payload model.PaymentHistory
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	// Validate required fields
	if payload.UserID == uuid.Nil || payload.PaymentID == uuid.Nil || payload.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "User ID, Payment ID, and valid amount are required",
		})
	}

	// Create payment history
	err := ctrl.paymentHistoryUseCase.CreatePaymentHistoryUseCase(&payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully created payment history",
		},
		Data: payload,
	})
}

// GetPaymentHistoryByUserIDController retrieves payment history by User ID
func (ctrl *paymentHistoryController) GetPaymentHistoryByUserIDController(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid User ID",
		})
	}

	response, err := ctrl.paymentHistoryUseCase.GetPaymentHistoryByUserIDUseCase(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved payment history",
		},
		Data: response,
	})
}
