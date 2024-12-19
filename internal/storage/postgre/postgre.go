package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/Suke2004/students-api/internal/config"
	"github.com/Suke2004/students-api/internal/types"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	// Create the connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User,
		cfg.Postgres.Password, cfg.Postgres.DbName, cfg.Postgres.SSLMode,
	)

	// Open the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Ensure the "students" table exists
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		age INT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	return &Postgres{Db: db}, nil
}

func (p *Postgres) CreateStudent(name, email string, age int) (int64, error) {
	var id int64
	err := p.Db.QueryRow(
		"INSERT INTO students (name, email, age) VALUES ($1, $2, $3) RETURNING id",
		name, email, age,
	).Scan(&id)
	return id, err
}

func (p *Postgres) GetStudentById(id int64) (types.Student, error) {
	var student types.Student
	err := p.Db.QueryRow(
		"SELECT id, name, email, age FROM students WHERE id = $1", id,
	).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	return student, err
}

func (p *Postgres) GetStudent() ([]types.Student, error) {
	rows, err := p.Db.Query("SELECT id, name, email, age FROM students")
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

func (p *Postgres) DeleteStudentById(id int64) error {
	_, err := p.Db.Exec("DELETE FROM students WHERE id = $1", id)
	return err
}

func (p *Postgres) DeleteAllStudents() error {
	_, err := p.Db.Exec("TRUNCATE TABLE students RESTART IDENTITY")
	return err
}

func (p *Postgres) UpdateStudentById(id int64, name, email string, age int) error {
	_, err := p.Db.Exec(
		"UPDATE students SET name = $1, email = $2, age = $3 WHERE id = $4",
		name, email, age, id,
	)
	return err
}
