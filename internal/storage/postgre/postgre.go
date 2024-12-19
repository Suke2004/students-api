package postgre

import "database/sql"

type Postgre struct {
	Db *sql.DB
}

func CreateStudent(name string, email string, age int) (int64, error){
	
}
