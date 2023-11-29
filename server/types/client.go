package types

import "github.com/gorilla/websocket"

type Client struct {
	Conn *websocket.Conn
	Id   int
}

func CreateClient(conn *websocket.Conn, id int) *Client {
	return &Client{Conn: conn, Id: id}
}
