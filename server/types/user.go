package types

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/facundocarballo/go-chat-app/crypto"
	"github.com/facundocarballo/go-chat-app/db"
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
	rows, err := database.Query("SELECT id, name, email FROM User")
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
		http.Error(w, "Error reading the body of request.", http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()

	user := BodyToUser(body)
	if user == nil {
		http.Error(w, "Error wrapping the body to user.", http.StatusBadRequest)
		return false
	}

	_, err = database.Exec(
		db.INSERT_USER_STATEMENT,
		user.Name,
		user.Email,
		user.Password,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error creating the user in the database. " + err.Error()))
		return false
	}

	tokenString := crypto.GenerateJWT(user.Id)

	resData := ResponseData{
		Message: *tokenString,
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

func GetUserFromEmail(email string, w http.ResponseWriter, database *sql.DB) *User {
	rows, err := database.Query("SELECT id, name, email, password FROM User WHERE email = (?)", email)
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
		http.Error(w, "Error reading the body of request.", http.StatusBadRequest)
		return false
	}
	defer r.Body.Close()

	user := BodyToUser(body)
	if user == nil {
		http.Error(w, "Error wrapping the body to user.", http.StatusBadRequest)
		return false
	}

	// Get the real user
	realUser := GetUserFromEmail(user.Email, w, database)
	if realUser == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User email not found."))
		return false
	}

	if realUser.Password != user.Password {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Passwords not match."))
		return false
	}

	tokenString := crypto.GenerateJWT(realUser.Id)

	resData := ResponseData{
		Message: *tokenString,
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
