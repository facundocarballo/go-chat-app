package websocket

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*types.Client]bool)
var broadcast = make(chan *types.Message)

func HandleWebSocket(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println("Error al hacer el Upgrade.")
		http.Error(w, errors.UPGRADE_FAILED+" "+err.Error(), http.StatusBadRequest)
		return
	}

	var groups []int
	client := types.CreateClient(conn, -1, groups)
	clients[client] = true

	println("Client [" + strconv.Itoa(client.Id) + "] Connected")

	defer func() {
		println("Cerramos la conexion y eliminamos al cliente. " + "Client [" + strconv.Itoa(client.Id) + "]")
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
		if message == nil {
			panic("Message is nil.")
		}

		if message.IsJWT {
			id := crypto.GetIdFromJWT(message.Message)
			if id == nil {
				println("JWT INVALID")
				return
			}
			client.Id = *id
			groups := types.GetGropusOfUser(*id, database)
			client.Groups = groups
			continue
		}

		if client.Id == -1 {
			println("No puedes enviar nada todavia, porque no tienes el JWT.")
			continue
		}

		var res bool
		if message.IsGroup {
			res = InsertGroupMessage(client.Id, message, database)
		} else {
			res = InsertUserMessage(client.Id, message, database)
		}

		if !res {
			println("Error sending the message.")
		}

		broadcast <- message
	}
}

func HandleMessages() {
	for {
		println("[HandleMessages] Looking for messages.")
		message := <-broadcast
		println("[HandleMessages] Message getted.")

		for client := range clients {
			println("[HandleMessages] Looking for the correct client.")
			if HaveToReceiveThisMessage(message, client) {
				println("[HandleMessages] Sending message to client.")
				SendMessage(client, message)
			}
		}
	}
}

func HaveToReceiveThisMessage(message *types.Message, client *types.Client) bool {
	if message.IsGroup {
		return client.BelongToThisGroup(message.ToId) && client.Id != message.UserId
	} else {
		return message.ToId == client.Id
	}
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
