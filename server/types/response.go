package types

import (
	"encoding/json"
)

type ResponseData struct {
	Message string `json:"message"`
}

type ResponseLogin struct {
	Jwt      string `json:"jwt"`
	RealUser User   `json:"user"`
}

func GetResponseDataJSON(res ResponseData) *[]byte {
	resJSON, err := json.Marshal(res)
	if err != nil {
		return nil
	}

	return &resJSON
}

func GetResponseLoginJSON(res ResponseLogin) *[]byte {
	resJSON, err := json.Marshal(res)
	if err != nil {
		return nil
	}

	return &resJSON
}
