package types

import "github.com/gorilla/websocket"

type Client struct {
	Conn   *websocket.Conn
	Id     int
	Groups []int
}

func (c Client) BelongToThisGroup(id int) bool {
	for _, groupId := range c.Groups {
		if id == groupId {
			return true
		}
	}
	return false
}

func CreateClient(conn *websocket.Conn, id int, groups []int) *Client {
	return &Client{Conn: conn, Id: id, Groups: groups}
}
