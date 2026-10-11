package up

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const UpBaseUrl = "https://api.up.com.au/api/v1/"

type Client struct {
	token   string
	http    *http.Client
	baseUrl string
}

func NewClient(token string, httpClient *http.Client) Client {
	return Client{
		token:   token,
		http:    httpClient,
		baseUrl: UpBaseUrl,
	}
}

func NewClientWithBaseUrl(token string, httpClient *http.Client, baseUrl string) Client {
	return Client{
		token:   token,
		http:    httpClient,
		baseUrl: baseUrl,
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
		return fmt.Errorf("Up GET %s failed: %s", url, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode Up response: %w", err)
	}
	return nil
}

func (c Client) ListAccounts(ctx context.Context) ([]uuid.UUID, error) {
	resp := &AccountsResponse{}
	err := c.Get(ctx, c.baseUrl+"accounts", resp)
	accounts := make([]uuid.UUID, 0, len(resp.Data))
	if err != nil {
		return accounts, err
	}
	for _, acc := range resp.Data {
		accounts = append(accounts, acc.Id)
	}
	return accounts, nil
}
