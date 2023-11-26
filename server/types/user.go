package types

import (
	"encoding/json"
	"time"
)

type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Password  string    `json:"password"`
}

func BodyToUser(body []byte) *User {
	if len(body) == 0 {
		return nil
	}

	var user User
	err := json.Unmarshal(body, &user)
	if err != nil {
		return nil
	}

	return &user
}
