package storage

import "github.com/Suke2004/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error) //int64 is for returning integer
	GetStudentById(id int64) (types.Student, error)
	GetStudent() ([]types.Student, error)
	DeleteStudentById(id int64) error
	DeleteAllStudents() error
	UpdateStudentById(id int64, name, email string, age int) error
}
