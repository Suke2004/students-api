package storage

import "github.com/Suke2004/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error) //int64 is for returning integer
	GetStudent(id int64) (types.Student, error)
}
