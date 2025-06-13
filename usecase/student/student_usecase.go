package student

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type StudentUseCase interface {
	GetAllStudents(page, limit int, name, class string) ([]*model.Student, error)
	GetStudentByID(ID uuid.UUID) (*model.Student, error)
	GetStudentsByParentName(parentName string) ([]*model.Student, error)
	CreateStudent(student *model.Student) (*model.Student, error)
	UpdateStudentByID(ID uuid.UUID, student *model.Student) (*model.Student, error)
	DeleteStudentByID(ID uuid.UUID) error
}

type studentUseCase struct {
	studentRepo repository.StudentRepository
}

func NewStudentUseCase(studentRepo repository.StudentRepository) StudentUseCase {
	return &studentUseCase{
		studentRepo: studentRepo,
	}
}

func (uc *studentUseCase) GetAllStudents(page, limit int, name, class string) ([]*model.Student, error) {
	// Validasi pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	students, err := uc.studentRepo.GetAllStudentsRepository(page, limit, name, class)
	if err != nil {
		return nil, fmt.Errorf("failed to get students: %v", err)
	}

	if len(students) == 0 {
		return nil, errors.New("no students found")
	}

	return students, nil
}

func (uc *studentUseCase) GetStudentByID(ID uuid.UUID) (*model.Student, error) {
	if ID == uuid.Nil {
		return nil, errors.New("invalid student ID")
	}

	student, err := uc.studentRepo.GetStudentByIdRepository(ID)
	if err != nil {
		return nil, fmt.Errorf("student not found: %v", err)
	}

	return student, nil
}

func (uc *studentUseCase) GetStudentsByParentName(parentName string) ([]*model.Student, error) {
	if parentName == "" {
		return nil, errors.New("parent name cannot be empty")
	}

	students, err := uc.studentRepo.GetStudentsByParentNameRepository(parentName)
	if err != nil {
		return nil, fmt.Errorf("failed to get students by parent name: %v", err)
	}

	if len(students) == 0 {
		return nil, errors.New("no students found for this parent")
	}

	return students, nil
}

func (uc *studentUseCase) CreateStudent(student *model.Student) (*model.Student, error) {
	// Validasi data siswa
	if student.FirstName == "" || student.LastName == "" {
		return nil, errors.New("first name and last name are required")
	}

	if student.Class == "" {
		return nil, errors.New("class is required")
	}

	if student.UserID == uuid.Nil {
		return nil, errors.New("parent ID is required")
	}

	createdStudent, err := uc.studentRepo.CreateStudentRepository(student)
	if err != nil {
		return nil, fmt.Errorf("failed to create student: %v", err)
	}

	return createdStudent, nil
}

func (uc *studentUseCase) UpdateStudentByID(ID uuid.UUID, student *model.Student) (*model.Student, error) {
	if ID == uuid.Nil {
		return nil, errors.New("invalid student ID")
	}

	// Validasi data yang akan diupdate
	if student.FirstName == "" || student.LastName == "" {
		return nil, errors.New("first name and last name are required")
	}

	if student.Class == "" {
		return nil, errors.New("class is required")
	}

	updatedStudent, err := uc.studentRepo.UpdateStudentByIdRepository(ID, student)
	if err != nil {
		return nil, fmt.Errorf("failed to update student: %v", err)
	}

	return updatedStudent, nil
}

func (uc *studentUseCase) DeleteStudentByID(ID uuid.UUID) error {
	if ID == uuid.Nil {
		return errors.New("invalid student ID")
	}

	err := uc.studentRepo.DeleteStudentByIdRepository(ID)
	if err != nil {
		return fmt.Errorf("failed to delete student: %v", err)
	}

	return nil
}
