package student

import (
	"Edupay/model"
	"Edupay/repository"
	"errors"
	"fmt"
)

type StudentUseCase interface {
	CreateStudentUseCase(payload *model.Student) (*model.Student, error)
	GetAllStudentUseCase(page, limit int, name, class string) ([]*model.Student, error)
	GetStudentByIdUseCase(Id string) (*model.Student, error)
	GetStudentsByParentNameUseCase(parentName string) ([]*model.Student, error)
	UpdatedStudentByIdUseCase(Id string, payload *model.Student) (*model.Student, error)
	DeleteStudentByIdUseCase(Id string) error
}

type studentUseCase struct {
	studentRepository repository.StudentRepository
}

// DeleteStudentByIdUseCase implements StudentUseCase.
func (uc *studentUseCase) DeleteStudentByIdUseCase(Id string) error {
	panic("unimplemented")
}

// UpdatedStudentByIdUseCase implements StudentUseCase.
func (uc *studentUseCase) UpdatedStudentByIdUseCase(Id string, payload *model.Student) (*model.Student, error) {
	panic("unimplemented")
}

func NewStudentUseCase(studentRepository repository.StudentRepository) *studentUseCase {
	return &studentUseCase{
		studentRepository: studentRepository,
	}
}

func (uc *studentUseCase) CreateStudentUseCase(payload *model.Student) (*model.Student, error) {
	student, err := uc.studentRepository.CreateStudentRepository(payload)
	if err != nil {
		return nil, fmt.Errorf("error creating student in database: %w", err)
	}
	return student, nil
}

func (uc *studentUseCase) GetAllStudentUseCase(page, limit int, name, class string) ([]*model.Student, error) {
	students, err := uc.studentRepository.GetAllStudentsRepository(page, limit, name, class)
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (uc *studentUseCase) GetStudentByIdUseCase(Id string) (*model.Student, error) {
	student, err := uc.studentRepository.GetStudentByIdRepository(Id)
	if err != nil {
		return nil, errors.New("student not found")
	}
	return student, nil
}

func (uc *studentUseCase) GetStudentsByParentNameUseCase(parentName string) ([]*model.Student, error) {
	students, err := uc.studentRepository.GetStudentsByParentNameRepository(parentName)
	if err != nil {
		return nil, fmt.Errorf("error getting students by parent name: %w", err)
	}
	return students, nil
}

func (uc *studentUseCase) UpdateStudentByIdRepository(Id string, payload *model.Student) (*model.Student, error) {
	student, err := uc.studentRepository.UpdateStudentByIdRepository(Id, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update student: %w", err)
	}
	return student, nil
}

func (uc *studentUseCase) DeleteStudentByIdRepository(Id string) error {
	err := uc.studentRepository.DeleteStudentByIdRepository(Id)
	if err != nil {
		return errors.New("student not found")
	}
	return nil
}
