package controller

import (
	"Edupay/model"
	usecase "Edupay/usecase/item"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ItemController adalah interface untuk operasi pada item
type ItemController interface {
	CreateItemController(c echo.Context) error
	GetAllItemsController(c echo.Context) error
	GetItemByIdController(c echo.Context) error
	UpdateItemByIdController(c echo.Context) error
	DeleteItemByIdController(c echo.Context) error
}

type itemController struct {
	itemUseCase usecase.ItemUseCase
}

// NewItemController membuat instance baru dari ItemController
func NewItemController(itemUseCase usecase.ItemUseCase) *itemController {
	return &itemController{
		itemUseCase: itemUseCase,
	}
}

// CreateItemController membuat item baru
func (ctrl *itemController) CreateItemController(c echo.Context) error {
	var payload model.Item

	// Cek apakah payload terbaca dengan benar
	if err := c.Bind(&payload); err != nil {
		fmt.Println("Error binding JSON:", err) // Tambahkan log di sini
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	// Log untuk melihat isi payload
	fmt.Printf("Payload received: %+v\n", payload)

	// Validasi user_id harus ada
	if payload.UserId == "" {
		fmt.Println("User ID is empty!") // Debugging
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "User ID is required",
		})
	}

	// Panggil use case untuk membuat item
	response, err := ctrl.itemUseCase.CreateItemUseCase(&payload)
	if err != nil {
		fmt.Println("Error in usecase:", err) // Debugging
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusCreated,
			Message:    "Successfully created item",
		},
		Data: response,
	})
}

// GetAllItemsController mengambil semua item
func (ctrl *itemController) GetAllItemsController(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	name := c.QueryParam("name")
	code := c.QueryParam("code")

	response, err := ctrl.itemUseCase.GetAllItemsUseCase(page, limit, name, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved items",
		},
		Data: response,
		Pagination: &model.Pagination{
			Page:  page,
			Limit: limit,
		},
	})
}

// GetItemByIdController mengambil item berdasarkan ID
func (ctrl *itemController) GetItemByIdController(c echo.Context) error {
	itemId := c.Param("id")

	response, err := ctrl.itemUseCase.GetItemByIdUseCase(itemId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved item",
		},
		Data: response,
	})
}

// UpdateItemByIdController memperbarui item
func (ctrl *itemController) UpdateItemByIdController(c echo.Context) error {
	itemId := c.Param("id")

	var payload model.Item
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	// Validasi apakah ItemDetails tidak kosong
	if len(payload.ItemDetails) == 0 {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "ItemDetails cannot be empty",
		})
	}

	response, err := ctrl.itemUseCase.UpdateItemByIdUseCase(itemId, &payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully updated item",
		},
		Data: response,
	})
}

// DeleteItemByIdController menghapus item
func (ctrl *itemController) DeleteItemByIdController(c echo.Context) error {
	itemId := c.Param("id")

	err := ctrl.itemUseCase.DeleteItemByIdUseCase(itemId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully deleted item",
		},
	})
}
