package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

// Connect connection to database
func Connect() (*sql.DB, error) {

	dbuser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")

	dataSource := fmt.Sprintf("%s:%s@/%s", dbuser, dbPassword, dbName)

	db, err := sql.Open("mysql", dataSource)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	fmt.Printf("db ping : %v", err)

	if err != nil {
		return nil, err
	}

	return db, nil
}
