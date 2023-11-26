package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/facundocarballo/go-chat-app/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const SERVER_PORT string = ":3690"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading the .env file.")
	}

	// Open connection to Database
	database, err := sql.Open("mysql", db.GetDSN())
	if err != nil {
		panic(err.Error())
	}
	defer database.Close()

	// Check success connection to Database
	err = database.Ping()
	if err != nil {
		panic(err.Error())
	}

	// Define handlers to endpoints
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// types.Login(w, r, database)
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		// types.HandleUser(w, r, database)
	})

	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		// types.HandleTask(w, r, database)
	})

	println("Server listening on port" + SERVER_PORT + " ...")
	http.ListenAndServe(SERVER_PORT, nil)
}
