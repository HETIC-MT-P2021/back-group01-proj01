package database

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	cLog "image_gallery/logger"
	"time"
)

// Connect connection to database
func Connect() (*sql.DB, error) {

	const (
		dbHost     = "tcp(db:3306)"
		dbName     = "image_gallery"
		dbUser     = "gallery"
		dbPassword = "gallery"
	)

	dsn := dbUser + ":" + dbPassword + "@" + dbHost + "/" + dbName + "?parseTime=true&charset=utf8"

	logger := cLog.GetLogger()

	logger.Infof("DSN: %s", dsn)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	var dbErr error
	for i := 1; i <= 3; i++ {
		dbErr = db.Ping()
		if dbErr != nil {
			if i < 3 {
				logger.Infof("nope, %d retry : %v", i, dbErr)
				time.Sleep(10 * time.Second)
			}
			continue
		}

		break
	}

	if dbErr != nil {
		return nil, errors.New("can't connect to database after 3 attempts")
	}

	return db, nil
}
