package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// CreateDatabase inits the db
func CreateDatabase() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbPass := os.Getenv("DBPASS")

	serverName := "localhost:3306"
	user := "root"
	password := dbPass
	dbName := "search"

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/", user, password, serverName)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS search")
	checkError(err)

	err = db.Close()
	checkError(err)

	connectionString = fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, serverName, dbName)
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Preparing the table
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS search (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, lowerName TEXT, name TEXT, url TEXT)")
	checkError(err)
	_, err = stmt.Exec()
	checkError(err)

	return db, err
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Database connection failed: " + err.Error())
	}
}
