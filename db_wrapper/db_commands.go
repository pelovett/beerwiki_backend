package db_wrapper

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func ExecSQL(sqlString string, args ...any) (sql.Result, error) {
	db_host := os.Getenv("DB_HOST")
	db_pass := os.Getenv("DB_PASSWORD")
	db_port := os.Getenv("DB_PORT")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=postgres "+
		"password=%s dbname=postgres sslmode=disable",
		db_host, db_port, db_pass)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Printf("Failed to connect to postgres database: %s", err)
		return nil, err
	}

	defer db.Close()

	result, err := db.Exec(sqlString, args...)

	if err != nil {
		log.Printf("Failed when executing sql string: %s \n With error: %s", sqlString, err)
	}

	return result, err
}
