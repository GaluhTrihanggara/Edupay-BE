package controller

import (
	"Edupay/model"
	"Edupay/usecase/book"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// Updated helper function to format price in Rupiah
func formatRupiah(price float64) string {
	// Use strconv to handle thousand separators
	priceStr := strconv.FormatFloat(price, 'f', 2, 64)
	parts := strings.Split(priceStr, ".")

	// Add thousand separators
	for i := len(parts[0]) - 3; i > 0; i -= 3 {
		parts[0] = parts[0][:i] + "." + parts[0][i:]
	}

	return "Rp " + parts[0] + "," + parts[1]
}

// Modify the Book struct to add a formatted price field
type BookResponse struct {
	model.Book
	FormattedPrice string `json:"formatted_price"`
}

type BookController interface {
	CreateBookController(c echo.Context) error
	GetAllBooksController(c echo.Context) error
	GetBookByIdController(c echo.Context) error
	UpdateBookByIdController(c echo.Context) error
	DeleteBookByIdController(c echo.Context) error
}

type bookController struct {
	bookUseCase book.BookUseCase
}

func NewBookController(bookUseCase book.BookUseCase) *bookController {
	return &bookController{
		bookUseCase: bookUseCase,
	}
}

// CreateBookController membuat buku baru
func (ctrl *bookController) CreateBookController(c echo.Context) error {
	var payload model.Book
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	// Validasi kode buku
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

// GetAllBooksController mengambil semua buku
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
	code := c.QueryParam("code") // Menambahkan pencarian berdasarkan kode

	response, err := ctrl.bookUseCase.GetAllBookUseCase(page, limit, title, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	var bookResponses []BookResponse
	for _, book := range response {
		bookResponses = append(bookResponses, BookResponse{
			Book:           *book,
			FormattedPrice: formatRupiah(book.Price),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved books",
		},
		Data: bookResponses,
		Pagination: &model.Pagination{
			Page:  page,
			Limit: limit,
		},
	})
}

// GetBookByIdController mengambil buku berdasarkan ID
func (ctrl *bookController) GetBookByIdController(c echo.Context) error {
	bookId := c.Param("id")
	if bookId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid book ID",
		})
	}

	response, err := ctrl.bookUseCase.GetBookByIdUseCase(bookId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	bookResponse := BookResponse{
		Book:           *response,
		FormattedPrice: formatRupiah(response.Price),
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved book",
		},
		Data: bookResponse,
	})
}

// UpdateBookByIdController memperbarui buku berdasarkan ID
func (ctrl *bookController) UpdateBookByIdController(c echo.Context) error {
	bookId := c.Param("id")
	if bookId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid book ID",
		})
	}

	var payload model.Book
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	response, err := ctrl.bookUseCase.UpdateBookByIdUseCase(bookId, &payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully updated book",
		},
		Data: response,
	})
}

// DeleteBookByIdController menghapus buku berdasarkan ID
func (ctrl *bookController) DeleteBookByIdController(c echo.Context) error {
	bookId := c.Param("id")
	if bookId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid book ID",
		})
	}

	err := ctrl.bookUseCase.DeleteBookByIdUseCase(bookId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully deleted book",
		},
	})
}
