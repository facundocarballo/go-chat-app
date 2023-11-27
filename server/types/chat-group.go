package types

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/facundocarballo/go-chat-app/crypto"
	"github.com/facundocarballo/go-chat-app/db"
	"github.com/facundocarballo/go-chat-app/errors"
)

type ChatGroup struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Sent        time.Time `json:"sent"`
	Owner       int       `json:"owner"`
}

func BodyToChatGroup(body []byte) *ChatGroup {
	if len(body) == 0 {
		return nil
	}

	var chatGroup ChatGroup
	err := json.Unmarshal(body, &chatGroup)
	if err != nil {
		return nil
	}

	return &chatGroup
}

func CreateGroup(
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

	chatGroup := BodyToChatGroup(body)
	if chatGroup == nil {
		http.Error(w, errors.UNMARSHAL+" Chat Group", http.StatusBadRequest)
		return false
	}

	_, err = database.Exec(
		db.INSERT_GROUP,
		chatGroup.Name,
		chatGroup.Description,
		chatGroup.Owner,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error creating the user in the database. " + err.Error()))
		return false
	}

	resData := ResponseData{
		Message: "SUCCESSFUL POST REQUEST",
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

func SendGroupRequest(
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

	chatGroup := BodyToChatGroup(body)
	if chatGroup == nil {
		http.Error(w, "Error wrapping the body to ChatGroup.", http.StatusBadRequest)
		return false
	}

	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, errors.JWT_INVALID, http.StatusBadRequest)
		return false
	}

	_, err = database.Exec(
		db.INSERT_GROUP_REQUEST,
		*id,
		chatGroup.Id,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error creating the group request in the database. " + err.Error()))
		return false
	}

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

func AcceptGroupRequest(
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

	chatGroup := BodyToChatGroup(body)
	if chatGroup == nil {
		http.Error(w, "Error wrapping the body to Chat Group.", http.StatusBadRequest)
		return false
	}

	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, errors.JWT_INVALID, http.StatusBadRequest)
		return false
	}

	// TODO: Chequear que sea el owner del grupo

	result, err := database.Exec(
		db.ACCEPT_FRINED_REQUEST,
		*id,
		chatGroup.Id,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error creating the group request in the database. " + err.Error()))
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading rows affected. " + err.Error()))
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

func GetGroupsRequests(
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

	group_id, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	if err != nil {
		http.Error(w, errors.ATOI, http.StatusBadRequest)
		return false
	}

	sented, err := strconv.ParseBool(r.URL.Query().Get("sented"))
	if err != nil {
		http.Error(w, errors.PARSE_BOOL, http.StatusBadRequest)
		return false
	}

	var query string
	var param int
	if sented {
		query = db.GET_GROUP_REQUEST_SENTED
		param = *id
	} else {
		query = db.GET_GROUP_REQUEST_RECEIVED
		param = group_id
	}

	rows, err := database.Query(query, param)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Iterate Rows
	var chatsGroups []ChatGroup
	for rows.Next() {
		var chatGroup ChatGroup
		err := rows.Scan(&chatGroup.Id, &chatGroup.Name, &chatGroup.Description, &chatGroup.Owner, &chatGroup.Sent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		chatsGroups = append(chatsGroups, chatGroup)
	}

	// Check Error on Rows
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chatsGroups)

	return true
}

func GetGroups(
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

	rows, err := database.Query(db.GET_GROUPS)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Iterate Rows
	var chatGroups []ChatGroup
	for rows.Next() {
		var chatGroup ChatGroup
		err := rows.Scan(&chatGroup.Id, &chatGroup.Name, &chatGroup.Description, &chatGroup.Sent, &chatGroup.Owner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		chatGroups = append(chatGroups, chatGroup)
	}

	// Check Error on Rows
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chatGroups)

	return true
}

func HandleGroupRequest(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.Method == http.MethodPost {
		SendGroupRequest(w, r, database)
		return
	}

	if r.Method == http.MethodGet {
		GetGroupsRequests(w, r, database)
		return
	}

	http.Error(w, "Method not allowed to /group-request", http.StatusMethodNotAllowed)
}

func HandleGroups(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.Method == http.MethodGet {
		GetGroups(w, r, database)
		return
	}

	if r.Method == http.MethodPost {
		CreateGroup(w, r, database)
	}

	http.Error(w, "Method not allowed to /groups", http.StatusMethodNotAllowed)
}
