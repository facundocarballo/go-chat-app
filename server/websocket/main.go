package websocket

import (
	"fmt"
	"net/http"

	"github.com/facundocarballo/go-chat-app/errors"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, errors.UPGRADE_FAILED, http.StatusBadRequest)
		return
	}
	defer conn.Close()

	println("Client Connected...")

	for {
		// Lee el mensaje del cliente
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			println(err)
			return
		}

		// Imprime el mensaje recibido
		print("Mensaje recibido: %s\n", p)

		// Escribe el mensaje de vuelta al cliente
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
