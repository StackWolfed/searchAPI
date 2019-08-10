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

	fmt.Printf("App listening to port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, app.Router))
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Database connection failed: " + err.Error())
	}
}
