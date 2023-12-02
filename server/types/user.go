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

type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Password  string    `json:"password"`
}

func BodyToUser(body []byte) *User {
	if len(body) == 0 {
		return nil
	}

	var user User
	err := json.Unmarshal(body, &user)
	if err != nil {
		return nil
	}

	return &user
}

func GetAllUsers(w http.ResponseWriter, database *sql.DB) []User {
	rows, err := database.Query(db.GET_USERS)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Iterate Rows
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
		users = append(users, user)
	}

	// Check Error on Rows
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

	return users
}

func CreateUser(
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

	user := BodyToUser(body)
	if user == nil {
		http.Error(w, errors.UNMARSHAL+" user", http.StatusBadRequest)
		return false
	}

	result, err := database.Exec(
		db.INSERT_USER_STATEMENT,
		user.Name,
		user.Email,
		user.Password,
	)
	// TODO: Modify this error message.
	if err != nil {
		http.Error(w, errors.GETTING_LAST_ID_INSERTED, http.StatusBadRequest)
		return false
	}

	newUserId, err := result.LastInsertId()
	if err != nil {
		http.Error(w, errors.GETTING_LAST_ID_INSERTED, http.StatusBadRequest)
		return false
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.INSERT_DB + err.Error()))
		return false
	}

	tokenString := crypto.GenerateJWT(int(newUserId))

	resData := ResponseData{
		Message: *tokenString,
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

func GetUserFromEmail(email string, w http.ResponseWriter, database *sql.DB) *User {
	rows, err := database.Query(db.GET_USER_BY_EMAIL, email)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Iterate Rows
	rows.Next()
	var user User
	err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	// Check Error on Rows
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return &user
}

func Login(
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

	user := BodyToUser(body)
	if user == nil {
		http.Error(w, errors.UNMARSHAL, http.StatusBadRequest)
		return false
	}

	// Get the real user
	realUser := GetUserFromEmail(user.Email, w, database)
	if realUser == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.ELEMENT_NOT_FOUND))
		return false
	}

	if realUser.Password != user.Password {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.PASSWORDS_NOT_MATCH))
		return false
	}

	tokenString := crypto.GenerateJWT(realUser.Id)

	resData := ResponseData{
		Message: *tokenString,
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

func HandleUser(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.Method == http.MethodPost {
		CreateUser(w, r, database)
		return
	}

	if r.Method == http.MethodGet {
		GetAllUsers(w, database)
		return
	}

	http.Error(w, "Method not allowed to /user", http.StatusMethodNotAllowed)
}
