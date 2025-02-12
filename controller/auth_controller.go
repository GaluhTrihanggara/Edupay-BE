package controller

import (
	"Edupay/model"
	"Edupay/usecase/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthController interface {
	LoginController(c echo.Context) error
	RegisterController(c echo.Context) error
}

type authController struct {
	authUseCase auth.AuthUseCase
}

func NewAuthController(authUseCase auth.AuthUseCase) *authController {
	return &authController{
		authUseCase: authUseCase,
	}
}

func (u *authController) LoginController(c echo.Context) error {
	var payload model.User

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	user, userTypeMessage, err := u.authUseCase.LoginUseCase(payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    userTypeMessage, // Menampilkan pesan sesuai UserType
		},
		Data: user,
	})
}

func (u *authController) RegisterController(c echo.Context) error {
	var payload model.User

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	user, err := u.authUseCase.RegisterUseCase(payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "User Created successfully",
		},
		Data: user,
	})
}

func (u *authController) RegisterAdminController(c echo.Context) error {
	var payload model.User
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}
	user, err := u.authUseCase.RegisterAdminUseCase(payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Admin Created successfully",
		},
		Data: user,
	})
}
