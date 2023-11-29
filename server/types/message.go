package types

import (
	"encoding/json"
	"time"
)

type Message struct {
	Id      int       `json:"id"`
	UserId  int       `json:"user_id"`
	IsGroup bool      `json:"is_group"`
	ToId    int       `json:"to_id"`
	Message string    `json:"message"`
	Sent    time.Time `json:"sent"`
}

func BodyToMessage(body []byte) *Message {
	if len(body) == 0 {
		return nil
	}

	var message Message
	err := json.Unmarshal(body, &message)
	if err != nil {
		return nil
	}

	return &message
}
