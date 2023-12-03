package types

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/facundocarballo/go-chat-app/crypto"
	"github.com/facundocarballo/go-chat-app/db"
	"github.com/facundocarballo/go-chat-app/errors"
)

type ChatGroup struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sent        string `json:"sent"`
	Owner       int    `json:"owner"`
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
	tokenString := crypto.GetJWTFromRequest(w, r)
	if tokenString == nil {
		return false
	}

	id := crypto.GetIdFromJWT(*tokenString)
	if id == nil {
		http.Error(w, errors.JWT_INVALID, http.StatusBadRequest)
		return false
	}

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
		db.CREATE_GROUP,
		chatGroup.Name,
		chatGroup.Description,
		*id,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.INSERT_DB + err.Error()))
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
	var sentBytes []uint8
	for rows.Next() {
		var chatGroup ChatGroup
		err := rows.Scan(&chatGroup.Id, &chatGroup.Name, &chatGroup.Description, &sentBytes, &chatGroup.Owner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		chatGroup.Sent = string(sentBytes)
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

func GetGroupOwner(id int, database *sql.DB) *int {
	rows, err := database.Query(db.GET_GROUP_OWNER, id)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Iterate Rows
	rows.Next()
	var chatId int
	err = rows.Scan(&chatId)
	if err != nil {
		return nil
	}

	// Check Error on Rows
	if err := rows.Err(); err != nil {
		return nil
	}

	return &chatId
}

func GetGropusOfUser(id int, database *sql.DB) []int {
	rows, err := database.Query(db.GET_GROUPS_OF_USER, id)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// Iterate Rows
	var chatGroups []int
	for rows.Next() {
		var chatId int
		err := rows.Scan(&chatId)
		if err != nil {
			return nil
		}
		chatGroups = append(chatGroups, chatId)
	}

	// Check Error on Rows
	if err := rows.Err(); err != nil {
		return nil
	}

	return chatGroups
}

func HandleGroups(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.Method == http.MethodGet {
		GetGroups(w, r, database)
		return
	}

	if r.Method == http.MethodPost {
		CreateGroup(w, r, database)
		return
	}

	http.Error(w, "Method not allowed to /groups", http.StatusMethodNotAllowed)
}
