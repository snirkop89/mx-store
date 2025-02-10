package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func initDB() {
	var err error
	// Initialize the db variable
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("missing DSN in ENV")
	}
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Check the database connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()

	// Setup MySQL
	initDB()
	defer db.Close()

	// Setup Static folder for static files and images
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	slog.Info("Starting server", "addr", ":5000")
	http.ListenAndServe(":5000", r)
}
