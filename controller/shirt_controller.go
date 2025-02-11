package controller

import (
	"Edupay/model"
	"Edupay/usecase/shirt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ShirtController interface {
	CreateShirtController(c echo.Context) error
	GetAllShirtsController(c echo.Context) error
	GetShirtByIdController(c echo.Context) error
	UpdateShirtByIdController(c echo.Context) error
	DeleteShirtByIdController(c echo.Context) error
}

type shirtController struct {
	shirtUseCase shirt.ShirtUseCase
}

func NewShirtController(shirtUseCase shirt.ShirtUseCase) *shirtController {
	return &shirtController{
		shirtUseCase: shirtUseCase,
	}
}

// CreateShirtController membuat entri kaos baru
func (ctrl *shirtController) CreateShirtController(c echo.Context) error {
	var payload model.Shirt
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	// Validasi kode unik
	if payload.Code == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Shirt code is required",
		})
	}

	response, err := ctrl.shirtUseCase.CreateShirtUseCase(&payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully created Shirt",
		},
		Data: response,
	})
}

// GetAllShirtsController mengambil semua kaos dengan filter nama, ukuran, dan kode
func (ctrl *shirtController) GetAllShirtsController(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	name := c.QueryParam("name")
	size := c.QueryParam("size")
	code := c.QueryParam("code") // Tambahan filter berdasarkan kode

	response, err := ctrl.shirtUseCase.GetAllShirtUseCase(page, limit, name, size, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved Shirts",
		},
		Data: response,
		Pagination: &model.Pagination{
			Page:  page,
			Limit: limit,
		},
	})
}

// GetShirtByIdController mengambil kaos berdasarkan ID
func (ctrl *shirtController) GetShirtByIdController(c echo.Context) error {
	shirtId := c.Param("id")
	if shirtId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Shirt ID",
		})
	}

	response, err := ctrl.shirtUseCase.GetShirtByIdUseCase(shirtId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved Shirt",
		},
		Data: response,
	})
}

// UpdateShirtByIdController memperbarui data kaos berdasarkan ID
func (ctrl *shirtController) UpdateShirtByIdController(c echo.Context) error {
	shirtId := c.Param("id")
	if shirtId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Shirt ID",
		})
	}

	var payload model.Shirt
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	response, err := ctrl.shirtUseCase.UpdateShirtByIdUseCase(shirtId, &payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully updated Shirt",
		},
		Data: response,
	})
}

// DeleteShirtByIdController menghapus kaos berdasarkan ID
func (ctrl *shirtController) DeleteShirtByIdController(c echo.Context) error {
	shirtId := c.Param("id")
	if shirtId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Shirt ID",
		})
	}

	err := ctrl.shirtUseCase.DeleteShirtByIdUseCase(shirtId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully deleted Shirt",
		},
	})
}
