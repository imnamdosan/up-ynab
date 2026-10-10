package ynab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const YnabBaseUrl = "https://api.ynab.com/v1/"

type Client struct {
	token string
	http  *http.Client
}

func NewClient(token string, httpClient *http.Client) Client {
	return Client{
		token: token,
		http:  httpClient,
	}
}

func (c Client) Get(ctx context.Context, url string, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("YNAB GET %s failed: %s", req.URL.Path, resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode YNAB response: %w", err)
	}
	return nil
}

func (c Client) ListPlans(ctx context.Context) ([]uuid.UUID, error) {
	var response PlansResponse
	if err := c.Get(ctx, YnabBaseUrl+"plans", &response); err != nil {
		return nil, err
	}
	plans := make([]uuid.UUID, 0, len(response.Data.Plans))
	for _, plan := range response.Data.Plans {
		plans = append(plans, plan.Id)
	}
	return plans, nil
}

func (c Client) ListAccounts(ctx context.Context, planID uuid.UUID) ([]uuid.UUID, error) {
	var response AccountsResponse
	if err := c.Get(ctx, YnabBaseUrl+"plans/"+planID.String()+"/accounts", &response); err != nil {
		return nil, err
	}
	accounts := make([]uuid.UUID, 0, len(response.Data.Accounts))
	for _, account := range response.Data.Accounts {
		accounts = append(accounts, account.Id)
	}
	return accounts, nil
}
