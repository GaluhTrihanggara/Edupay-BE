package controller

import (
	"Edupay/model"
	"Edupay/usecase/cart"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CartController interface {
	CreateCartController(c echo.Context) error
	GetAllCartsController(c echo.Context) error
	GetCartByIDController(c echo.Context) error
	GetCartByUserIDController(c echo.Context) error
	UpdateCartByIdController(c echo.Context) error
	DeleteCartByIdController(c echo.Context) error
}

type cartController struct {
	cartUseCase cart.CartUseCase
}

func NewCartController(cartUseCase cart.CartUseCase) *cartController {
	return &cartController{
		cartUseCase: cartUseCase,
	}
}

// CreateCartController creates a new cart
func (ctrl *cartController) CreateCartController(c echo.Context) error {
	var payload model.Cart
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	// Validate required fields
	if payload.UserID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "User ID is required",
		})
	}

	response, err := ctrl.cartUseCase.CreateCartUseCase(&payload)
	if err != nil {
		if strings.Contains(err.Error(), "cart already exists") {
			return c.JSON(http.StatusConflict, model.ErrorResponse{
				StatusCode: http.StatusConflict,
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
			Message:    "Successfully created Cart",
		},
		Data: response,
	})
}

func (ctrl *cartController) GetAllCartsController(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	userID := c.QueryParam("user_id")

	carts, err := ctrl.cartUseCase.GetAllCartsUseCase(page, limit, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Carts retrieved successfully",
		},
		Data: carts,
	})
}

func (ctrl *cartController) GetCartByIDController(c echo.Context) error {
	cartIDStr := c.Param("cart_id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid cart ID: " + err.Error(),
		})
	}

	cart, err := ctrl.cartUseCase.GetCartByIDUseCase(cartID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve cart: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Cart retrieved successfully",
		},
		Data: cart,
	})
}

// GetCartByUserIDController retrieves a cart by User ID
func (ctrl *cartController) GetCartByUserIDController(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid User ID" + err.Error(),
		})
	}

	response, err := ctrl.cartUseCase.GetCartByUserIDUseCase(userID)
	if err != nil {
		if err.Error() == fmt.Sprintf("cart not found for user %s", userID) {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve cart:" + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved Cart",
		},
		Data: response,
	})
}

// UpdateCartByIdController updates a cart by ID
func (ctrl *cartController) UpdateCartByIdController(c echo.Context) error {
	cartIDStr := c.Param("id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Cart ID",
		})
	}

	var payload model.Cart
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload" + err.Error(),
		})
	}

	// Validate required fields
	if payload.UserID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "user ID is required",
		})
	}

	response, err := ctrl.cartUseCase.UpdateCartByIdUseCase(cartID, &payload)
	if err != nil {
		if err.Error() == "cart not found" {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update cart:" + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully updated Cart",
		},
		Data: response,
	})
}

// DeleteCartByIdController deletes a cart by ID
func (ctrl *cartController) DeleteCartByIdController(c echo.Context) error {
	cartIDStr := c.Param("id")
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Cart ID" + err.Error(),
		})
	}

	err = ctrl.cartUseCase.DeleteCartByIdUseCase(cartID)
	if err != nil {
		if err.Error() == "cart not found" {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete cart:" + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully deleted Cart",
		},
	})
}
