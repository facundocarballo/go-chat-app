package types

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/facundocarballo/go-chat-app/crypto"
	"github.com/facundocarballo/go-chat-app/db"
	"github.com/facundocarballo/go-chat-app/errors"
)

type ChatGroupRequest struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id"`
	GroupId int    `json:"group_id"`
	Sent    string `json:"sent"`
}

func BodyToChatGroupRequest(body []byte) *ChatGroupRequest {
	if len(body) == 0 {
		return nil
	}

	var chatGroup ChatGroupRequest
	err := json.Unmarshal(body, &chatGroup)
	if err != nil {
		return nil
	}

	return &chatGroup
}

func SendGroupRequest(
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

	request := BodyToChatGroupRequest(body)
	if request == nil {
		http.Error(w, errors.UNMARSHAL+" ChatGroupRequest", http.StatusBadRequest)
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
		request.GroupId,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.INSERT_DB + err.Error()))
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
	var requests []ChatGroupRequest
	for rows.Next() {
		var request ChatGroupRequest
		var sentBytes []uint8
		err := rows.Scan(&request.Id, &request.UserId, &request.GroupId, &sentBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		request.Sent = string(sentBytes)
		requests = append(requests, request)
	}

	// Check Error on Rows
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	// Send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(requests)

	return true
}

func AcceptGroupRequest(
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

	request := BodyToChatGroupRequest(body)
	if request == nil {
		http.Error(w, errors.UNMARSHAL+" ChatGroupRequest", http.StatusBadRequest)
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

	owner := GetGroupOwner(request.GroupId, database)
	if *id != *owner {
		http.Error(w, errors.NOT_GROUP_OWNER, http.StatusBadRequest)
		return false
	}

	_, err = database.Exec(
		db.ACCEPT_GROUP_REQUEST,
		request.UserId,
		request.GroupId,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.INSERT_DB + err.Error()))
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

func HandleAcceptGroupRequest(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.Method == http.MethodPost {
		AcceptGroupRequest(w, r, database)
		return
	}

	http.Error(w, "Method not allowed to /accept-request", http.StatusMethodNotAllowed)
}
