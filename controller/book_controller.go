package controller

import (
	"Edupay/model"
	"Edupay/usecase/book"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BookController interface {
	CreateBookController(c echo.Context) error
	GetAllBooksController(c echo.Context) error
	GetBookByIdController(c echo.Context) error
	GetBookByCodeController(echo.Context) error
	UpdateBookByIDController(c echo.Context) error
	DeleteBookByIDController(c echo.Context) error
}

type bookController struct {
	bookUseCase book.BookUseCase
}

func NewBookController(bookUseCase book.BookUseCase) *bookController {
	return &bookController{
		bookUseCase: bookUseCase,
	}
}

// CreateBookController creates a new book entry
func (ctrl *bookController) CreateBookController(c echo.Context) error {
	var payload model.Book
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	// Validate Code Book
	if payload.Code == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Book code is required",
		})
	}

	response, err := ctrl.bookUseCase.CreateBookUseCase(&payload)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully created book",
		},
		Data: response,
	})
}

// GetAllBooksController retrieves all books with filters
func (ctrl *bookController) GetAllBooksController(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

	title := c.QueryParam("title")
	class := c.QueryParam("class")

	response, err := ctrl.bookUseCase.GetAllBookUseCase(page, limit, title, class)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved Books",
		},
		Data: response,
		Pagination: &model.Pagination{
			Page:  page,
			Limit: limit,
		},
	})
}

// GetBookByIdController retrieves a book by ID
func (ctrl *bookController) GetBookByIDController(c echo.Context) error {
	bookIdStr := c.Param("id")
	bookId, err := uuid.Parse(bookIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Book ID",
		})
	}

	response, err := ctrl.bookUseCase.GetBookByIDUseCase(bookId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved Book",
		},
		Data: response,
	})
}

func (ctrl *bookController) GetBookByCodeController(c echo.Context) error {
	bookCode := c.Param("code")
	if bookCode == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid book code",
		})
	}

	response, err := ctrl.bookUseCase.GetBookByCodeUseCase(bookCode)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Succesfully retrieved Book",
		},
		Data: response,
	})
}

// UpdateBookByIdController updates book information by ID
func (ctrl *bookController) UpdateBookByIDController(c echo.Context) error {
	bookIdStr := c.Param("id")
	bookId, err := uuid.Parse(bookIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Book ID",
		})
	}

	var payload model.Book
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	response, err := ctrl.bookUseCase.UpdateBookByIDUseCase(bookId, &payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully updated Book",
		},
		Data: response,
	})
}

// DeleteBookByIdController deletes a book by ID
func (ctrl *bookController) DeleteBookByIDController(c echo.Context) error {
	bookIdStr := c.Param("id")
	bookId, err := uuid.Parse(bookIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Book ID",
		})
	}

	err = ctrl.bookUseCase.DeleteBookByIDUseCase(bookId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully deleted Book",
		},
	})
}
