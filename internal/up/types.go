package up

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

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

type TransactionStatus string

const (
	TransactionHeld    TransactionStatus = "HELD"
	TransactionSettled TransactionStatus = "SETTLED"
)

type Money struct {
	CurrencyCode     string `json:"currencyCode"`
	Value            string `json:"value"`
	ValueInBaseUnits int64  `json:"valueInBaseUnits"`
}

type TransactionsResponse struct {
	Data struct {
		ID         uuid.UUID `json:"id"`
		Attributes struct {
			Status      TransactionStatus `json:"status"`
			RawText     *string           `json:"rawText"`
			Description string            `json:"description"`
			Message     *string           `json:"message"`
			Amount      Money             `json:"amount"`
			SettledAt   *time.Time        `json:"settledAt"`
			CreatedAt   time.Time         `json:"createdAt"`
		} `json:"attributes"`
		Relationships struct {
			Account struct {
				Data struct {
					ID uuid.UUID `json:"id"`
				} `json:"data"`
			} `json:"account"`
		} `json:"relationships"`
	} `json:"data"`
}

func (r TransactionsResponse) toTransaction() (Transaction, error) {
	attributes := r.Data.Attributes
	accountID := r.Data.Relationships.Account.Data.ID

	if r.Data.ID == uuid.Nil || accountID == uuid.Nil || attributes.CreatedAt.IsZero() {
		return Transaction{}, fmt.Errorf("transaction response missing required fields")
	}

	if attributes.Status != TransactionHeld && attributes.Status != TransactionSettled {
		return Transaction{}, fmt.Errorf("unknown transaction status %q", attributes.Status)
	}
	return Transaction{
		ID:          r.Data.ID,
		Status:      attributes.Status,
		RawText:     attributes.RawText,
		Description: attributes.Description,
		Message:     attributes.Message,
		Amount:      attributes.Amount,
		SettledAt:   attributes.SettledAt,
		CreatedAt:   attributes.CreatedAt,
		AccountId:   accountID,
	}, nil
}

type Transaction struct {
	ID          uuid.UUID         `json:"id"`
	Status      TransactionStatus `json:"status"`
	RawText     *string           `json:"rawText"`
	Description string            `json:"description"`
	Message     *string           `json:"message"`
	Amount      Money             `json:"amount"`
	SettledAt   *time.Time        `json:"settledAt"`
	CreatedAt   time.Time         `json:"createdAt"`
	AccountId   uuid.UUID         `json:"accountId"`
}
