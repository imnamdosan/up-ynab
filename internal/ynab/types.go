package ynab

import "github.com/google/uuid"

type PlansResponse struct {
	Data struct {
		Plans []Plan `json:"plans"`
	} `json:"data"`
}

type Plan struct {
	Id uuid.UUID `json:"id"`
}

type AccountsResponse struct {
	Data struct {
		Accounts []Account `json:"accounts"`
	} `json:"data"`
}

type Account struct {
	Id uuid.UUID `json:"id"`
}
