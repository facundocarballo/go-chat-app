package types

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/facundocarballo/go-chat-app/crypto"
	"github.com/facundocarballo/go-chat-app/db"
	"github.com/facundocarballo/go-chat-app/errors"
)

type Message struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id"`
	IsGroup bool   `json:"is_group"`
	ToId    int    `json:"to_id"`
	Message string `json:"message"`
	Sent    string `json:"sent"`
}

func BodyToMessage(body []byte) *Message {
	if len(body) == 0 {
		return nil
	}

	var message Message
	err := json.Unmarshal(body, &message)
	if err != nil {
		println("Error: ", err.Error())
		return nil
	}

	return &message
}

func HandleUserMessage(
	w http.ResponseWriter,
	r *http.Request,
	database *sql.DB,
) {
	if r.Method == http.MethodGet {
		GetUserMessage(w, r, database)
		return
	}

	http.Error(w, "Method not allowed to /user-message", http.StatusMethodNotAllowed)
}

func GetUserMessage(
	w http.ResponseWriter,
	r *http.Request,
	database *sql.DB,
) bool {
	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		http.Error(w, errors.JWT_NOT_FOUND, http.StatusBadRequest)
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, errors.JWT_INVALID, http.StatusBadRequest)
		return false
	}

	rows, err := database.Query(db.GET_USER_MESSAGES, *id, *id)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		var sentBytes []uint8
		err := rows.Scan(&message.Id, &message.UserId, &message.ToId, &message.Message, &sentBytes)
		message.Sent = string(sentBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		messages = append(messages, message)
	}

	// Check Error on Rows
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)

	return true
}
