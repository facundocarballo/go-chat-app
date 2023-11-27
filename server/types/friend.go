package types

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/facundocarballo/go-chat-app/crypto"
	"github.com/facundocarballo/go-chat-app/db"
	"github.com/facundocarballo/go-chat-app/errors"
)

type Friend struct {
	Id    int       `json:"id"`
	UserA int       `json:"user_a"`
	UserB int       `json:"user_b"`
	Sent  time.Time `json:"sent"`
}

func BodyToFriend(body []byte) *Friend {
	if len(body) == 0 {
		return nil
	}

	var friend Friend
	err := json.Unmarshal(body, &friend)
	if err != nil {
		return nil
	}

	return &friend
}

func SendFriendRequest(
	w http.ResponseWriter,
	r *http.Request,
	database *sql.DB,
) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, errors.READING_BODY_REQ, http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()

	friend := BodyToFriend(body)
	if friend == nil {
		http.Error(w, errors.UNMARSHAL+" friend.", http.StatusBadRequest)
		return false
	}

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

	_, err = database.Exec(
		db.INSERT_FRIEND_REQUEST,
		*id,
		friend.UserB,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.INSERT_DB + " " + err.Error()))
		return false
	}

	resData := ResponseData{
		Message: "SUCCESFULL POST REQUEST",
	}
	resJSON := GetResponseDataJSON(resData)

	if resJSON == nil {
		http.Error(w, errors.DATA_TO_JSON, http.StatusInternalServerError)
		return false
	}

	w.WriteHeader(http.StatusOK)
	w.Write(*resJSON)

	return true
}

func AcceptFriendRequest(
	w http.ResponseWriter,
	r *http.Request,
	database *sql.DB,
) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, errors.READING_BODY_REQ, http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()

	friend := BodyToFriend(body)
	if friend == nil {
		http.Error(w, errors.UNMARSHAL+" friend.", http.StatusBadRequest)
		return false
	}

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

	result, err := database.Exec(
		db.ACCEPT_FRINED_REQUEST,
		*id,
		friend.UserB,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.INSERT_DB + " " + err.Error()))
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.READING_ROWS_AFFECTED + " " + err.Error()))
		return false
	}

	// TODO: Chequear cuantas filas se afectan cuando no se produce la aceptacion de la solicitud de amistad.
	println("Filas afectadas: ", rowsAffected)

	resData := ResponseData{
		Message: "SUCCESFULL POST REQUEST",
	}
	resJSON := GetResponseDataJSON(resData)

	if resJSON == nil {
		http.Error(w, "Error converting the response data to JSON. ", http.StatusInternalServerError)
		return false
	}

	w.WriteHeader(http.StatusOK)
	w.Write(*resJSON)

	return true
}

func GetFriendRequests(
	w http.ResponseWriter,
	r *http.Request,
	database *sql.DB,
	sented bool,
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

	var query string
	if sented {
		query = db.GET_FRIEND_REQUEST_SENTED
	} else {
		query = db.GET_FRIEND_REQUEST_RECEIVED
	}

	rows, err := database.Query(query, id)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Iterate Rows
	var friends []Friend
	for rows.Next() {
		var friend Friend
		err := rows.Scan(&friend.Id, &friend.UserA, &friend.UserB, &friend.Sent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		friends = append(friends, friend)
	}

	// Check Error on Rows
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(friends)

	return true
}

func GetFriends(
	w http.ResponseWriter,
	r *http.Request,
	database *sql.DB,
) bool {
	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, errors.JWT_INVALID, http.StatusBadRequest)
		return false
	}

	rows, err := database.Query(db.GET_FRIENDS, id, id)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Iterate Rows
	var friends []Friend
	for rows.Next() {
		var friend Friend
		err := rows.Scan(&friend.Id, &friend.UserA, &friend.UserB, &friend.Sent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		friends = append(friends, friend)
	}

	// Check Error on Rows
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(friends)

	return true
}

func HandleFriendRequest(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.Method == http.MethodPost {
		SendFriendRequest(w, r, database)
		return
	}

	if r.Method == http.MethodGet {
		// Get all the friend requests received
		// TODO: Read from a parameter if the client wants the friend requests sented or the received
		GetFriendRequests(w, r, database, false)
		return
	}

	http.Error(w, "Method not allowed to /friend-request", http.StatusMethodNotAllowed)
}

func HandleFriends(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.Method == http.MethodGet {
		GetFriends(w, r, database)
		return
	}

	if r.Method == http.MethodPost {
		AcceptFriendRequest(w, r, database)
		return
	}

	http.Error(w, "Method not allowed to /friends", http.StatusMethodNotAllowed)
}
