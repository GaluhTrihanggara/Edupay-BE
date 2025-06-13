package controller

import (
	"Edupay/model"
	"Edupay/usecase/order"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderController interface {
	CreateOrderController(c echo.Context) error
	GetAllOrdersController(c echo.Context) error
	GetOrderByIdController(c echo.Context) error
	UpdateOrderByIdController(c echo.Context) error
	DeleteOrderByIdController(c echo.Context) error
}

type orderController struct {
	orderUseCase order.OrderUseCase
}

func NewOrderController(orderUseCase order.OrderUseCase) *orderController {
	return &orderController{
		orderUseCase: orderUseCase,
	}
}

// CreateOrderController creates a new order
func (ctrl *orderController) CreateOrderController(c echo.Context) error {
	var payload model.Order
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload: " + err.Error(),
		})
	}

	// Validate required fields
	if payload.UserID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "User ID is required",
		})
	}
	if payload.CartID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Cart ID is required",
		})
	}

	response, err := ctrl.orderUseCase.CreateOrderUseCase(&payload)
	if err != nil {
		if strings.Contains(err.Error(), "cart is empty") || strings.Contains(err.Error(), "cart not found") {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to create order: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Order created successfully",
		},
		Data: response,
	})
}

func (ctrl *orderController) GetAllOrdersController(c echo.Context) error {
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

	orders, err := ctrl.orderUseCase.GetAllOrdersUseCase(page, limit, userID, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Orders retrieved successfully",
		},
		Data: orders,
	})
}

// GetOrderByIdController retrieves an order by ID
func (ctrl *orderController) GetOrderByIdController(c echo.Context) error {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid order ID: " + err.Error(),
		})
	}

	response, err := ctrl.orderUseCase.GetOrderByIdUseCase(orderID)
	if err != nil {
		if strings.Contains(err.Error(), "order not found") {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve order: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Order retrieved successfully",
		},
		Data: response,
	})
}

// UpdateOrderByIdController updates an order by ID
func (ctrl *orderController) UpdateOrderByIdController(c echo.Context) error {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid order ID: " + err.Error(),
		})
	}

	var payload model.Order
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload: " + err.Error(),
		})
	}

	// Validate status if provided
	if payload.Status != "" {
		validStatuses := []string{"pending", "completed", "cancelled"}
		valid := false
		for _, s := range validStatuses {
			if s == payload.Status {
				valid = true
				break
			}
		}
		if !valid {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid order status: " + payload.Status,
			})
		}
	}

	response, err := ctrl.orderUseCase.UpdateOrderByIdUseCase(orderID, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "order not found") {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update order: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Order updated successfully",
		},
		Data: response,
	})
}

// DeleteOrderByIdController deletes an order by ID
func (ctrl *orderController) DeleteOrderByIdController(c echo.Context) error {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid order ID: " + err.Error(),
		})
	}

	err = ctrl.orderUseCase.DeleteOrderByIdUseCase(orderID)
	if err != nil {
		if strings.Contains(err.Error(), "order not found") {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete order: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Order deleted successfully",
		},
	})
}
