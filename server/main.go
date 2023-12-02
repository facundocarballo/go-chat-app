package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/facundocarballo/go-chat-app/db"
	"github.com/facundocarballo/go-chat-app/types"
	"github.com/facundocarballo/go-chat-app/websocket"
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
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		types.HandleUser(w, r, database)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		types.Login(w, r, database)
	})

	http.HandleFunc("/friend-request", func(w http.ResponseWriter, r *http.Request) {
		types.HandleFriendRequest(w, r, database)
	})

	http.HandleFunc("/accept-friend", func(w http.ResponseWriter, r *http.Request) {
		types.HandleAcceptFriend(w, r, database)
	})

	http.HandleFunc("/friends", func(w http.ResponseWriter, r *http.Request) {
		types.HandleFriends(w, r, database)
	})

	http.HandleFunc("/group-request", func(w http.ResponseWriter, r *http.Request) {
		types.HandleGroupRequest(w, r, database)
	})

	http.HandleFunc("/group", func(w http.ResponseWriter, r *http.Request) {
		types.HandleGroups(w, r, database)
	})

	http.HandleFunc("/acept-group", func(w http.ResponseWriter, r *http.Request) {
		types.AcceptGroupRequest(w, r, database)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(w, r, database)
	})

	go websocket.HandleMessages()

	println("Server listening on port" + SERVER_PORT + " ...")
	http.ListenAndServe(SERVER_PORT, nil)
}
