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

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
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

func (s *Sqlite) GetStudent() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)

	}
	return students, nil
}

// DeleteStudentById deletes a student by their ID.
func (s *Sqlite) DeleteStudentById(id int64) error {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	// // Reindex the table to maintain a continuous sequence
	// err = s.reindexTable()
	// if err != nil {    //If wanted to auto reindex the table if a middle value is deleted use this
	// 	return err
	// }

	return nil
}

// DeleteAllStudents deletes all students from the database.
func (s *Sqlite) DeleteAllStudents() error {
	stmt, err := s.Db.Prepare("DELETE FROM students")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	_, err = s.Db.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'students'")
	if err != nil {
		return err
	}

	return nil
}

// UpdateStudentById updates a student's data by their ID.
func (s *Sqlite) UpdateStudentById(id int64, name, email string, age int) error {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, email, age, id)
	if err != nil {
		return err
	}

	return nil
}

// reindexTable reorders the IDs to maintain a continuous sequence.
func (s *Sqlite) reindexTable() error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}

	// Use a temporary table to reorder IDs
	_, err = tx.Exec(`
		CREATE TEMPORARY TABLE temp_students AS SELECT * FROM students ORDER BY id;
		DELETE FROM students;
		INSERT INTO students (id, name, email, age)
		SELECT ROW_NUMBER() OVER (ORDER BY id), name, email, age FROM temp_students;
		DROP TABLE temp_students;
	`)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
