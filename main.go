package main

import (
	"fmt"
	"log"
	"net/http"
	"search/app"
	"search/db"

	"github.com/gorilla/mux"
)

var port = "3001"

func main() {
	database, err := db.CreateDatabase()
	checkError(err)

	app := &app.App{
		Router:   mux.NewRouter().StrictSlash(true),
		Database: database,
	}

	app.SetupRouter()

	// What I've learned from this:
	// SELECT queries can use the .Query()
	// Anything that needs to edit has to use the stmt.Exec()
	// stmt, err := app.Database.Prepare("CREATE TABLE IF NOT EXISTS test (id int)")
	// stmt.Exec()
	// stmt, err = app.Database.Prepare("INSERT INTO test (id) VALUES (1234)")
	// stmt.Exec()
	// rows, err := app.Database.Query("SELECT * FROM test")
	// rows.Next()
	// var test int
	// rows.Scan(&test)
	// fmt.Println(test)

	fmt.Printf("App listening to port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, app.Router))
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Database connection failed: " + err.Error())
	}
}
