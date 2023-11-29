package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/facundocarballo/go-chat-app/crypto"
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

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
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

	client := types.CreateClient(conn, *id)
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
		// TODO: Insert message to database
		broadcast <- message
	}
}

func HandleMessages() {
	for {
		// Obtener el próximo mensaje del canal de difusión
		message := <-broadcast

		// TODO: Chequear si el mensaje es para un grupo

		// Enviar el mensaje a todos los clientes conectados
		for client := range clients {
			if message.Id == client.Id {
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
		}
	}
}
