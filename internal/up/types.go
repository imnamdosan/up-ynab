package up

import "github.com/google/uuid"

type PingResponse struct {
	Meta struct {
		Id     uuid.UUID `json:"id"`
		Status string    `json:"statusEmoji"`
	} `json:"meta"`
}

type Account struct {
	Id uuid.UUID `json:"id"`
}

type AccountsResponse struct {
	Data  []Account `json:"data"`
	Links struct {
		Next *string `json:"next"`
	} `json:"links"`
}
