package mocks

import (
	model "Edupay/model"

	"github.com/stretchr/testify/mock"
)

type StudentRepository struct {
	mock.Mock
}

// DeleteStudentByIDRepository mock for the method DeleteStudentByIDRepository
func (_m *StudentRepository) DeleteStudentByIDRepository(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllStudentsRepository mock for the method GetAllStudentsRepository
func (_m *StudentRepository) GetAllStudentsRepository(page int, limit int, class string) ([]*model.Student, error) {
	ret := _m.Called(page, limit, class)

	var r0 []*model.Student
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int, string) ([]*model.Student, error)); ok {
		return rf(page, limit, class)
	}
	if rf, ok := ret.Get(0).(func(int, int, string) []*model.Student); ok {
		r0 = rf(page, limit, class)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Student)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, string) error); ok {
		r1 = rf(page, limit, class)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStudentByIDRepository mock for the method GetStudentByIDRepository
func (_m *StudentRepository) GetStudentByIDRepository(id string) (*model.Student, error) {
	ret := _m.Called(id)

	var r0 *model.Student
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.Student, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *model.Student); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Student)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertStudentRepository mock for the method InsertStudentRepository
func (_m *StudentRepository) InsertStudentRepository(student *model.Student) error {
	ret := _m.Called(student)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Student) error); ok {
		r0 = rf(student)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateStudentByIDRepository mock for the method UpdateStudentByIDRepository
func (_m *StudentRepository) UpdateStudentByIDRepository(id string, student *model.Student) (*model.Student, error) {
	ret := _m.Called(id, student)

	var r0 *model.Student
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *model.Student) (*model.Student, error)); ok {
		return rf(id, student)
	}
	if rf, ok := ret.Get(0).(func(string, *model.Student) *model.Student); ok {
		r0 = rf(id, student)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Student)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *model.Student) error); ok {
		r1 = rf(id, student)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStudentRepository creates a new instance of StudentRepository
// The first argument is usually *testing.T
func NewStudentRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *StudentRepository {
	mock := &StudentRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
