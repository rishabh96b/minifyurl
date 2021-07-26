package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/rishabh96b/minifyurl/app/handler"
	"github.com/rishabh96b/minifyurl/config"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name)
	db, err := sql.Open(config.DB.Driver, dbURI)
	if err != nil {
		log.Fatal("Could not connect to database")
	}
	a.DB = db
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	a.Get("/api/v1/{ref}", a.handleRequest(handler.RedirectToOriginalURL))
	a.Post("/api/v1/new", a.handleRequest(handler.CreateURLEntry))
}

// Get is a wrapper around the router for GET method.
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post is a wrapper around the router for POST method.
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) Run(port string) {
	server := &http.Server{
		Handler:      a.Router,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Listening on port", port)
	log.Fatal(server.ListenAndServe())
}

type RequestHandlerFunction func(db *sql.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}
