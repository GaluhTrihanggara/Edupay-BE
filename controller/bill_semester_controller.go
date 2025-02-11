package controller

import (
	"Edupay/model"
	billsemester "Edupay/usecase/bill_semester"
	"Edupay/usecase/transaction"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BillSemesterController interface {
	CreateBillSemesterController(c echo.Context) error
	GetAllBillSemesterController(c echo.Context) error
	GetBillSemesterByIdController(c echo.Context) error
	UpdateBillSemesterByIdController(c echo.Context) error
	DeleteBillSemesterByIdController(c echo.Context) error
	CreateBillSemesterTransactionController(c echo.Context) error // NEW
}

type billSemesterController struct {
	billSemesterUseCase billsemester.BillSemesterUseCase
	transactionUseCase  transaction.TransactionUseCase // NEW
}

func NewBillSemesterController(billSemesterUseCase billsemester.BillSemesterUseCase, transactionUseCase transaction.TransactionUseCase) *billSemesterController {
	return &billSemesterController{
		billSemesterUseCase: billSemesterUseCase,
		transactionUseCase:  transactionUseCase,
	}
}

// CreateBillSemesterController membuat tagihan semester baru
func (ctrl *billSemesterController) CreateBillSemesterController(c echo.Context) error {
	var payload model.BillSemester
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Request payload",
		})
	}

	response, err := ctrl.billSemesterUseCase.CreateBillSemesterUseCase(&payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Succesfully created Bill Semester",
		},
		Data: response,
	})
}

// CreateBillSemesterTransactionController membuat transaksi untuk tagihan semester
func (ctrl *billSemesterController) CreateBillSemesterTransactionController(c echo.Context) error {
	billSemesterId := c.Param("id")
	if billSemesterId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Bill Semester Id",
		})
	}

	// Ambil data tagihan semester
	billSemester, err := ctrl.billSemesterUseCase.GetBillSemesterByIdUseCase(billSemesterId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	// Buat transaksi baru berdasarkan tagihan semester
	transactionPayload := &model.Transaction{
		UserId:          billSemester.StudentId,
		BillSemesterId:  &billSemesterId,
		TotalPrice:      billSemester.Amount,
		Status:          model.STATUS_UNPAID,
		Description:     "Payment for semester bill",
		TransactionDate: billSemester.DueDate,
		Method:          "Manual Payment", // Metode pembayaran dapat diatur sesuai kebutuhan
	}

	transaction, err := ctrl.transactionUseCase.CreateTransactionUseCase(transactionPayload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create transaction: " + err.Error(),
		})
	}

	// Perbarui tagihan semester dengan ID transaksi yang baru dibuat
	billSemester.TransactionId = &transaction.Id
	_, err = ctrl.billSemesterUseCase.UpdateBillSemesterByIdUseCase(billSemesterId, billSemester)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update bill semester with transaction: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully created transaction for Bill Semester",
		},
		Data: transaction,
	})
}

// GetAllBillSemesterController mengambil semua tagihan semester
func (ctrl *billSemesterController) GetAllBillSemesterController(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	semester := c.QueryParam("semester")
	year := c.QueryParam("year")

	response, err := ctrl.billSemesterUseCase.GetAllBillSemesterUseCase(page, limit, semester, year)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved Bill Semesters",
		},
		Data: response,
		Pagination: &model.Pagination{
			Page:  page,
			Limit: limit,
		},
	})
}

// GetBillSemesterByIdController mengambil tagihan semester berdasarkan ID
func (ctrl *billSemesterController) GetBillSemesterByIdController(c echo.Context) error {
	billSemesterId := c.Param("id")
	if billSemesterId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Bill Semester Id",
		})
	}
	response, err := ctrl.billSemesterUseCase.GetBillSemesterByIdUseCase(billSemesterId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Succesfully retrieved bill semester",
		},
		Data: response,
	})
}

// UpdateBillSemesterByIdController memperbarui tagihan semester berdasarkan ID
func (ctrl *billSemesterController) UpdateBillSemesterByIdController(c echo.Context) error {
	billSemesterId := c.Param("id")
	if billSemesterId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Bill Semester Id",
		})
	}
	var payload model.BillSemester
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}
	response, err := ctrl.billSemesterUseCase.UpdateBillSemesterByIdUseCase(billSemesterId, &payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully updated Bill Semester",
		},
		Data: response,
	})
}

// DeleteBillSemesterByIdController menghapus tagihan semester berdasarkan ID
func (ctrl *billSemesterController) DeleteBillSemesterByIdController(c echo.Context) error {
	billsemesterId := c.Param("id")
	if billsemesterId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Bill Semester Id",
		})
	}

	err := ctrl.billSemesterUseCase.DeleteBillSemesterByIdUseCase(billsemesterId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Succesfully delete Bill Semester",
		},
	})
}
