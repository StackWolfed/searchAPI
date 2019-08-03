package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"search/models"

	"github.com/gorilla/mux"
)

// App structure to keep the router and db
type App struct {
	Router   *mux.Router
	Database *sql.DB
}

const prefix string = "/api/v1"

// SetupRouter sets the mothods and assigns it to the .Router of the app
func (app *App) SetupRouter() {
	app.Router.Methods("GET").Path("/").HandlerFunc(app.testFunction)
	app.Router.Methods("GET").Path("/api/search").HandlerFunc(app.searchFunction)
	app.Router.Methods("POST").Path(prefix + "/users").HandlerFunc(app.addEntry)
}

func (app *App) testFunction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode()
	result := r.FormValue("test") // ?test=%s query parameter
	fmt.Fprintf(w, "ok "+result)
}

func (app *App) searchFunction(w http.ResponseWriter, r *http.Request) {
	rows, err := app.Database.Query("SELECT * FROM test WHERE id LIKE 5")
	checkError(err)

	var result string
	var tmp string
	for rows.Next() {
		rows.Scan(&tmp)
		result += ", " + tmp
	}

	fmt.Fprintf(w, "ok "+result)
	w.WriteHeader(http.StatusOK)
}

// /api/v1/addEntry route
func (app *App) addEntry(w http.ResponseWriter, r *http.Request) {
	// Check if entry already exists
	stmt, err := app.Database.Prepare("SELECT name FROM search WHERE name = ?")
	checkError(err)
	// JSON Parsing
	entry := parseBodyToEntry(r.Body)
	// RUN the query
	rows, err := stmt.Query(entry.Name)
	checkError(err)
	if foundInRows(rows) {
		http.Error(w, "entry already exists", http.StatusBadRequest)
		return
	}
	// Add the entry
	stmt, err = app.Database.Prepare("INSERT INTO search (name, url) VALUES (?, ?)")
	checkError(err)
	res, err := stmt.Exec(entry.Name, entry.URL)
	checkError(err)
	fmt.Println(res.RowsAffected())
	fmt.Fprintf(w, "ok")
}

func foundInRows(rows *sql.Rows) bool {
	for rows.Next() {
		var queryResult string
		rows.Scan(&queryResult)
		if queryResult != "" {
			return true
		}
	}
	return false
}

func parseBodyToEntry(body io.Reader) models.Entry {
	decoder := json.NewDecoder(body)
	var entry models.Entry
	err := decoder.Decode(&entry)
	checkError(err)
	return entry
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Database connection failed: " + err.Error())
	}
}
