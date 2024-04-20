package db

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var db *sql.DB

func Connect() error {

	// Capture connection properties.
	cfg := mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       os.Getenv("DB_NETWORK"),
		Addr:      os.Getenv("DB_HOST") + os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_DATABASE"),
		ParseTime: true,
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	log.Println("Connected to MySQL database!")
	return nil
}
