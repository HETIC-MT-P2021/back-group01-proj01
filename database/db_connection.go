package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// Connect connection to database
func Connect() (*sql.DB, error) {

	const (
		dbHost = "tcp(host.docker.internal:3306)"
		dbName = "image_gallery"
		dbUser = "root"
		dbPass = "root"
	)

	dsn := dbUser + ":" + dbPass + "@" + dbHost + "/" + dbName + "?charset=utf8"
	var err error

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	err = db.Ping()

	fmt.Printf("db ping : %v", err)

	if err != nil {
		return nil, err
	}

	return db, nil
}
