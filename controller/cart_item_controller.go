package controller

import (
	"Edupay/model"
	"Edupay/usecase/cart_item"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CartItemController interface {
	AddCartItemController(c echo.Context) error
	GetAllCartItemsController(c echo.Context) error
	GetCartItemByItemIDController(c echo.Context) error
	GetCartItemsByCartIDController(c echo.Context) error
	UpdateCartItemByIDController(c echo.Context) error
	DeleteCartItemByIDController(c echo.Context) error
}

type cartItemController struct {
	cartItemUseCase cart_item.CartItemUseCase
}

func NewCartItemController(cartItemUseCase cart_item.CartItemUseCase) *cartItemController {
	return &cartItemController{
		cartItemUseCase: cartItemUseCase,
	}
}

func (ctrl *cartItemController) AddCartItemController(c echo.Context) error {
	var payload model.CartItem
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	// Validate required fields
	if payload.CartID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Cart ID is required",
		})
	}
	if payload.ItemID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Item ID is required",
		})
	}
	if payload.ItemType != "book" && payload.ItemType != "shirt" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Item type must be 'book' or 'shirt'",
		})
	}
	if payload.Quantity <= 0 {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Quantity must be greater than zero",
		})
	}

	response, err := ctrl.cartItemUseCase.AddCartItemUseCase(&payload)
	if err != nil {
		if strings.Contains(err.Error(), "already exists in cart") {
			return c.JSON(http.StatusConflict, model.ErrorResponse{
				StatusCode: http.StatusConflict,
				Message:    err.Error(),
			})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		if strings.Contains(err.Error(), "insufficient stock") {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to add cart item: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusCreated,
			Message:    "Cart item added successfully",
		},
		Data: response,
	})
}

func (ctrl *cartItemController) GetAllCartItemsController(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	cartID := c.QueryParam("cart_id")

	cartItems, err := ctrl.cartItemUseCase.GetAllCartItemsUseCase(page, limit, cartID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Cart items retrieved successfully",
		},
		Data: cartItems,
	})
}

func (ctrl *cartItemController) GetCartItemsByCartIDController(c echo.Context) error {
	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Cart ID",
		})
	}

	response, err := ctrl.cartItemUseCase.GetCartItemsByCartIDUseCase(cartID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved Cart Items",
		},
		Data: response,
	})
}

func (ctrl *cartItemController) UpdateCartItemByIDController(c echo.Context) error {
	cartItemIDStr := c.Param("id")
	cartItemID, err := uuid.Parse(cartItemIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Cart Item ID",
		})
	}

	var payload model.CartItem
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	// Validasi tambahan
	if payload.Quantity <= 0 {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Quantity must be greater than zero",
		})
	}

	response, err := ctrl.cartItemUseCase.UpdateCartItemByIDUseCase(cartItemID, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		if strings.Contains(err.Error(), "insufficient stock") {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update cart item: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully updated Cart Item",
		},
		Data: response,
	})
}

func (ctrl *cartItemController) DeleteCartItemByIDController(c echo.Context) error {
	cartItemIDStr := c.Param("id")
	cartItemID, err := uuid.Parse(cartItemIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Cart Item ID",
		})
	}

	err = ctrl.cartItemUseCase.DeleteCartItemByIDUseCase(cartItemID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete cart item: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully deleted Cart Item",
		},
	})
}
