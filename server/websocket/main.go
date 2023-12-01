package websocket

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/facundocarballo/go-chat-app/crypto"
	"github.com/facundocarballo/go-chat-app/db"
	"github.com/facundocarballo/go-chat-app/errors"
	"github.com/facundocarballo/go-chat-app/types"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*types.Client]bool)
var broadcast = make(chan *types.Message)

func HandleWebSocket(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	tokenString := r.URL.Query().Get("jwt")

	id := crypto.GetIdFromJWT(tokenString)
	if id == nil {
		http.Error(w, errors.JWT_INVALID, http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, errors.UPGRADE_FAILED, http.StatusBadRequest)
		return
	}

	groups := types.GetGropusOfUser(*id, database)
	client := types.CreateClient(conn, *id, groups)
	clients[client] = true

	defer func() {
		conn.Close()
		delete(clients, client)
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			println(err)
			return
		}
		message := types.BodyToMessage(p)
		var res bool
		if message.IsGroup {
			res = InsertGroupMessage(*id, message, database)
		} else {
			res = InsertUserMessage(*id, message, database)
		}

		if !res {
			message.Message = "Error sending the message: " + message.Message
		}

		broadcast <- message

	}
}

func HandleMessages() {
	for {

		message := <-broadcast

		for client := range clients {

			if HaveToReceiveThisMessage(message, *client) {
				SendMessage(client, message)
			}
		}
	}
}

func HaveToReceiveThisMessage(message *types.Message, client types.Client) bool {
	return message.IsGroup && client.BelongToThisGroup(message.ToId) || message.ToId == client.Id
}

func SendMessage(client *types.Client, message *types.Message) {
	b, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		client.Conn.Close()
		delete(clients, client)
	}
	err = client.Conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Println(err)
		client.Conn.Close()
		delete(clients, client)
	}
}

func InsertGroupMessage(id int, message *types.Message, database *sql.DB) bool {
	_, err := database.Exec(
		db.INSERT_GROUP_MESSAGE,
		id,
		message.ToId,
		message.Message,
	)

	return err == nil
}

func InsertUserMessage(id int, message *types.Message, database *sql.DB) bool {
	_, err := database.Exec(
		db.INSERT_USER_MESSAGE,
		id,
		message.ToId,
		message.Message,
	)

	return err == nil
}
