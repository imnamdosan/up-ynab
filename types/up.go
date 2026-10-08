package types

type PingResponse struct {
	Meta struct {
		Id     string `json:"id"`
		Status string `json:"statusEmoji"`
	} `json:"meta"`
}

type AccountsResponse struct {
	Data []struct {
		Id string `json:"id"`
	} `json:"data"`
	Links struct {
		Next *string `json:"next"`
	} `json:"links"`
}
