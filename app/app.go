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
	app.Router.Methods("GET").Path(prefix + "/search").HandlerFunc(app.searchFunction)
	app.Router.Methods("POST").Path(prefix + "/entries").HandlerFunc(app.addEntry)
	app.Router.Methods("DELETE").Path(prefix + "/entries").HandlerFunc(app.deleteEntry)
}

func (app *App) testFunction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode()
	result := r.FormValue("test") // ?test=%s query parameter
	fmt.Fprintf(w, "ok "+result)
}

// /api/v1/search route
func (app *App) searchFunction(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("search") == "" {
		http.Error(w, "Invalid search parameter", http.StatusBadRequest)
		return
	}

	rows, err := app.Database.Query("SELECT * FROM search WHERE name LIKE CONCAT('%', ?, '%')", r.FormValue("search"))
	checkError(err)

	var result models.Items
	var tmp models.Item
	for rows.Next() {
		err = rows.Scan(&tmp.ID, &tmp.Name, &tmp.URL)
		if err != nil {
			fmt.Println("Error in Scan:")
			fmt.Println(err)
			http.Error(w, "An Error has occured", http.StatusInternalServerError)
			return
		}
		result = append(result, tmp)
	}

	byte, err := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.Write(byte)
}

// /api/v1/addEntry route
func (app *App) addEntry(w http.ResponseWriter, r *http.Request) {
	// Check if entry already exists
	stmt, err := app.Database.Prepare("SELECT name FROM search WHERE name = ?")
	checkError(err)
	// JSON Parsing
	entry := parseBodyToEntry(r.Body)

	// Check if both the URL and name parameters are set correctly
	if entry.Name == "" || entry.URL == "" {
		http.Error(w, "Wrong parameters", http.StatusBadRequest)
		return
	}

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

// /api/v1/deleteEntry route
func (app *App) deleteEntry(w http.ResponseWriter, r *http.Request) {
	// JSON Parsing
	entry := parseBodyToEntry(r.Body)
	// Check if exist and if it does, gets the ID
	rows, err := app.Database.Query("SELECT id FROM search WHERE name = ?", entry.Name)
	checkError(err)
	if !rows.Next() {
		http.Error(w, "entry does not exist", http.StatusBadRequest)
		return
	}
	var id int64
	rows.Scan(&id)
	stmt, err := app.Database.Prepare("DELETE FROM search WHERE id = ?")
	checkError(err)
	res, err := stmt.Exec(id)
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
