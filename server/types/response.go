package types

import (
	"encoding/json"
)

type ResponseData struct {
	Message string `json:"message"`
}

func GetResponseDataJSON(res ResponseData) *[]byte {
	resJSON, err := json.Marshal(res)
	if err != nil {
		return nil
	}

	return &resJSON
}
