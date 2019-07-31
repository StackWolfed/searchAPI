package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// App structure to keep the router and db
type App struct {
	Router   *mux.Router
	Database *sql.DB
}

// SetupRouter sets the mothods and assigns it to the .Router of the app
func (app *App) SetupRouter() {
	app.Router.Methods("GET").Path("/").HandlerFunc(app.testFunction)
	app.Router.Methods("GET").Path("/api/search").HandlerFunc(app.searchFunction)
}

func (app *App) testFunction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode()
	result := r.FormValue("test") // ?test=%s query parameter
	fmt.Fprintf(w, "ok "+result)
}

func (app *App) searchFunction(w http.ResponseWriter, r *http.Request) {
	_, err := app.Database.Exec("INSERT INTO `test` (name) VALUES ('myname')")
	if err != nil {
		log.Fatal("Database INSERT failed")
	}

	log.Println("You called a thing!")
	w.WriteHeader(http.StatusOK)
}
