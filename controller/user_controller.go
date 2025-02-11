package controller

import (
	"Edupay/model"
	"Edupay/usecase/middlewares"
	user "Edupay/usecase/users"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserController interface {
	GetAllUsersController(c echo.Context) error
	GetUserByIdController(c echo.Context) error
	GetUserByPhoneController(c echo.Context) error
	GetUserByEmailController(c echo.Context) error
	GetUserByQueryController(c echo.Context) error
	UpdateUserByIdController(c echo.Context) error
	DeleteUserByIdController(c echo.Context) error
}

type userController struct {
	userUseCase user.UserUseCase
}

func NewUserController(userUseCase user.UserUseCase) *userController {
	return &userController{
		userUseCase: userUseCase,
	}
}

// GetAllUsersController mengambil semua user
func (ctrl *userController) GetAllUsersController(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	name := c.QueryParam("name")

	response, err := ctrl.userUseCase.GetAllUsersUseCase(page, limit, name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved users",
		},
		Data: response,
	})
}

// GetUserByIdController mengambil user berdasarkan ID
func (ctrl *userController) GetUserByIdController(c echo.Context) error {
	userId := c.Param("id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid user ID",
		})
	}
	response, err := ctrl.userUseCase.GetUserByIdUseCase(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved user",
		},
		Data: response,
	})
}

// GetUserByPhoneController mengambil user berdasarkan nomor telepon
func (ctrl *userController) GetUserByPhoneController(c echo.Context) error {
	phone := c.QueryParam("phone")
	if phone == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Phone number is required",
		})
	}
	response, err := ctrl.userUseCase.GetUserByPhoneUseCase(phone)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved user",
		},
		Data: response,
	})
}

// GetUserByEmailController mengambil user berdasarkan email
func (ctrl *userController) GetUserByEmailController(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Email is required",
		})
	}
	response, err := ctrl.userUseCase.GetUserByEmailUseCase(email)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved user",
		},
		Data: response,
	})
}

func (ctrl *userController) GetUserByQueryController(c echo.Context) error {
	user := middlewares.ExtractTokenUserId(model.ADMIN_TYPE, c)
	if user == "" {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "token unauthorized",
		})
	}
	query := c.QueryParam("query")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

	users, err := ctrl.userUseCase.GetUserByQueryUseCase(query, page, limit)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully Get Users",
		},
		Data: users,
		Pagination: &model.Pagination{
			Page:  page,
			Limit: limit,
		},
	})
}

// UpdateUserByIdController memperbarui user berdasarkan ID
func (ctrl *userController) UpdateUserByIdController(c echo.Context) error {
	userId := c.Param("id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid user ID",
		})
	}

	var payload model.User
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	response, err := ctrl.userUseCase.UpdateUserByIdUseCase(userId, &payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully updated user",
		},
		Data: response,
	})
}

// DeleteUserByIdController menghapus user berdasarkan ID
func (ctrl *userController) DeleteUserByIdController(c echo.Context) error {
	userId := c.Param("id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid user ID",
		})
	}

	err := ctrl.userUseCase.DeleteUserByIdUseCase(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully deleted user",
		},
	})
}
