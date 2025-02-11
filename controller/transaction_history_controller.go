package controller

import (
	"Edupay/model"
	transactionhistory "Edupay/usecase/transaction_history"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TransactionHistoryController interface {
	GetAllHistoriesController(c echo.Context) error
	GetHistoryByIdController(c echo.Context) error
	GetHistoriesByUserIdController(c echo.Context) error
	CreateHistoryFromTransactionController(c echo.Context) error
	DeleteHistoryByIdController(c echo.Context) error
}

type transactionHistoryController struct {
	transactionHistoryUseCase transactionhistory.TransactionHistoryUseCase
}

func NewTransactionHistoryController(transactionHistoryUseCase transactionhistory.TransactionHistoryUseCase) *transactionHistoryController {
	return &transactionHistoryController{
		transactionHistoryUseCase: transactionHistoryUseCase,
	}
}

// GetAllHistoriesController mengambil semua riwayat transaksi dengan pagination
func (ctrl *transactionHistoryController) GetAllHistoriesController(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

	response, err := ctrl.transactionHistoryUseCase.GetAllHistoriesUseCase(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved transaction histories",
		},
		Data: response,
		Pagination: &model.Pagination{
			Page:  page,
			Limit: limit,
		},
	})
}

// GetHistoryByIdController mengambil riwayat transaksi berdasarkan ID
func (ctrl *transactionHistoryController) GetHistoryByIdController(c echo.Context) error {
	historyId := c.Param("id")
	if historyId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid transaction history ID",
		})
	}

	response, err := ctrl.transactionHistoryUseCase.GetHistoryByIdUseCase(historyId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved transaction history",
		},
		Data: response,
	})
}

// GetHistoriesByUserIdController mengambil semua riwayat transaksi berdasarkan UserId
func (ctrl *transactionHistoryController) GetHistoriesByUserIdController(c echo.Context) error {
	userId := c.Param("user_id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid User ID",
		})
	}

	response, err := ctrl.transactionHistoryUseCase.GetHistoriesByUserIdUseCase(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved transaction histories for user",
		},
		Data: response,
	})
}

// CreateHistoryFromTransactionController membuat riwayat transaksi dari data transaksi utama
func (ctrl *transactionHistoryController) CreateHistoryFromTransactionController(c echo.Context) error {
	var transaction model.Transaction
	if err := c.Bind(&transaction); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	response, err := ctrl.transactionHistoryUseCase.CreateHistoryFromTransactionUseCase(&transaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully created transaction history",
		},
		Data: response,
	})
}

// DeleteHistoryByIdController menghapus riwayat transaksi berdasarkan ID
func (ctrl *transactionHistoryController) DeleteHistoryByIdController(c echo.Context) error {
	historyId := c.Param("id")
	if historyId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid transaction history ID",
		})
	}

	err := ctrl.transactionHistoryUseCase.DeleteHistoriesByIdUseCase(historyId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully deleted transaction history",
		},
	})
}
