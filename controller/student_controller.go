package controller

import (
	"Edupay/model"
	"Edupay/usecase/student"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type StudentController interface {
	CreateStudentController(c echo.Context) error
	GetAllStudentsController(c echo.Context) error
	GetStudentByIdController(c echo.Context) error
	GetStudentsByParentNameController(c echo.Context) error
	UpdateStudentByIdController(c echo.Context) error
	DeleteStudentByIdController(c echo.Context) error
}

type studentController struct {
	studentUseCase student.StudentUseCase
}

func NewStudentController(studentUseCase student.StudentUseCase) *studentController {
	return &studentController{
		studentUseCase: studentUseCase,
	}
}

func (ctrl *studentController) CreateStudentController(c echo.Context) error {
	var payload model.Student
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	response, err := ctrl.studentUseCase.CreateStudentUseCase(&payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully created Student",
		},
		Data: response,
	})
}

func (ctrl *studentController) GetAllStudentsController(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	name := c.QueryParam("name")
	class := c.QueryParam("class")

	response, err := ctrl.studentUseCase.GetAllStudentUseCase(page, limit, name, class)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved Students",
		},
		Data: response,
		Pagination: &model.Pagination{
			Page:  page,
			Limit: limit,
		},
	})
}

func (ctrl *studentController) GetStudentByIdController(c echo.Context) error {
	studentId := c.Param("id")
	if studentId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Student Id",
		})
	}
	response, err := ctrl.studentUseCase.GetStudentByIdUseCase(studentId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved Student",
		},
		Data: response,
	})
}

func (ctrl *studentController) GetStudentsByParentNameController(c echo.Context) error {
	parentName := c.QueryParam("parent_name")
	if parentName == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Parent name is required",
		})
	}

	// Call use case to get students by parent's name
	students, err := ctrl.studentUseCase.GetStudentsByParentNameUseCase(parentName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully retrieved students by parent name",
		},
		Data: students,
	})
}

func (ctrl *studentController) UpdateStudentByIdController(c echo.Context) error {
	studentId := c.Param("id")
	if studentId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Student Id",
		})
	}

	var payload model.Student
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request payload",
		})
	}

	response, err := ctrl.studentUseCase.UpdatedStudentByIdUseCase(studentId, &payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully updated Student",
		},
		Data: response,
	})
}

func (ctrl *studentController) DeleteStudentByIdController(c echo.Context) error {
	studentId := c.Param("id")
	if studentId == "" {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Student Id",
		})
	}

	err := ctrl.studentUseCase.DeleteStudentByIdUseCase(studentId)
	if err != nil {
		return c.JSON(http.StatusNotFound, model.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, model.HttpResponse{
		MetaData: model.MetaData{
			StatusCode: http.StatusOK,
			Message:    "Successfully deleted Student",
		},
	})
}
