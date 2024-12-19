package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/Suke2004/students-api/internal/config"
	"github.com/Suke2004/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3" //it is used indirectly behind so we kept _
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath) //("driver name",Storage path)//you have to install that driver to project from githib
	if err != nil {
		return nil, err
	}

	//Create table if not exist
	db.Exec(`CREATE TABLE IF NOT EXISTS students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    email TEXT,
    age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name,email,age) VALUES(?,?,?)") //we are preparing it first so to prevent sql injection

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)

	if err != nil {
		return 0, nil
	}

	lastId, err := result.LastInsertId() //it will return id of last inserted row

	if err != nil {
		return 0, nil
	}

	return lastId, nil

}

func (s *Sqlite) GetStudent(id int64) (types.Student, error) {
	//sql query
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("query error %w", err)
	}

	return student, nil
}
