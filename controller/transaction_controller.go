package controller

import (
	"Edupay/model"
	transaction "Edupay/usecase/transaction"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TransactionController interface {
	GetAllTransactionsController(c echo.Context) error
	GetTransactionByIdController(c echo.Context) error
	GetTransactionByUserIdController(c echo.Context) error
	UpdateTransactionStatusController(c echo.Context) error
}

type transactionController struct {
	transactionUseCase transaction.TransactionUseCase
}

func NewTransactionController(transactionUseCase transaction.TransactionUseCase) *transactionController {
	return &transactionController{
		transactionUseCase: transactionUseCase,
	}
}

// GetAllTransactionController mengambil semua transaksi dengan filter
func (ctrl *transactionController) GetAllTransactionsController(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

	response, err := ctrl.transactionUseCase.GetAllTransactionUseCase(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved transactions",
		},
		Data: response,
		Pagination: &model.Pagination{
			Page:  page,
			Limit: limit,
		},
	})
}

// GetTransactionByIdController mengambil transaksi berdasarkan ID
func (ctrl *transactionController) GetTransactionByIdController(c echo.Context) error {
	transactionId := c.Param("id")
	if transactionId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Transaction ID",
		})
	}

	response, err := ctrl.transactionUseCase.GetTransactionByIdUseCase(transactionId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved transaction",
		},
		Data: response,
	})
}

// GetTransactionByUserIdController mengambil transaksi berdasarkan userId
func (ctrl *transactionController) GetTransactionByUserIdController(c echo.Context) error {
	userId := c.QueryParam("user_id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid User ID",
		})
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

	response, err := ctrl.transactionUseCase.GetTransactionByUserIdUseCase(userId, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved transactions for user",
		},
		Data: response,
		Pagination: &model.Pagination{
			Page:  page,
			Limit: limit,
		},
	})
}

func (ctrl *transactionController) UpdateTransactionStatusController(c echo.Context) error {
	transactionId := c.Param("id")
	if transactionId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Transaction Id",
		})
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	// Ambil transaksi yang ada
	existingTransaction, err := ctrl.transactionUseCase.GetTransactionByIdUseCase(transactionId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "Transaction not found",
		})
	}

	// Perbarui status transaksi
	existingTransaction.Status = payload.Status
	updatedTransaction, err := ctrl.transactionUseCase.UpdateTransactionByIdUseCase(transactionId, existingTransaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update transaction status: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully updated transaction status",
		},
		Data: updatedTransaction,
	})
}
