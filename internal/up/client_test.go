package up

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func testClient(t *testing.T, path string, status int, body string) Client {
	token := "test-token"
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != path {
			t.Errorf("path = %q, want %s", r.URL.Path, path)
		}

		if value := r.Header.Get("Authorization"); value != "Bearer "+token {
			t.Errorf("Authorization = %q", value)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	t.Cleanup(server.Close)
	return NewClientWithBaseUrl(token, server.Client(), server.URL+"/")
}

func TestClient_ListAccounts(t *testing.T) {
	accountId, _ := uuid.NewRandom()
	body := fmt.Sprintf(`{"data":[{"id":%q}]}`, accountId)
	client := testClient(t, "/accounts", http.StatusOK, body)

	resp, err := client.ListAccounts(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(resp) != 1 || resp[0] != accountId {
		t.Fatalf("accounts = %+v", resp)
	}
}

func TestClient_ListAccountsServerError(t *testing.T) {
	client := testClient(t, "/accounts", http.StatusInternalServerError, "")
	resp, err := client.ListAccounts(t.Context())
	if err == nil || len(resp) != 0 {
		t.Fatal("didn't process fail status", resp)
	}
}

func TestClient_GetTransactionHeld(t *testing.T) {
	transactionID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	accountID := uuid.MustParse("22222222-2222-4222-8222-222222222222")
	createdAt := "2026-10-10T12:34:56+11:00"
	body := fmt.Sprintf(`{
		"data": {
			"id": %q,
			"attributes": {
				"status": "HELD",
				"description": "Coffee",
				"rawText": null,
				"message": null,
				"amount": {
					"currencyCode": "AUD",
					"value": "-12.99",
					"valueInBaseUnits": -1299
				},
				"createdAt": %q,
				"settledAt": null
			},
			"relationships": {
				"account": {"data": {"id": %q}}
			}
		}
	}`, transactionID.String(), createdAt, accountID.String())
	client := testClient(t, "/transactions/"+transactionID.String(), http.StatusOK, body)

	got, err := client.GetTransaction(t.Context(), transactionID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != transactionID || got.AccountId != accountID {
		t.Errorf("IDs = (%s, %s), want (%s, %s)", got.ID, got.AccountId, transactionID, accountID)
	}
	if got.Status != TransactionHeld || got.SettledAt != nil {
		t.Errorf("lifecycle = (%q, %v), want (HELD, nil)", got.Status, got.SettledAt)
	}
	if got.Amount != (Money{CurrencyCode: "AUD", Value: "-12.99", ValueInBaseUnits: -1299}) {
		t.Errorf("amount = %+v", got.Amount)
	}
	if got.RawText != nil || got.Message != nil || got.Description != "Coffee" {
		t.Errorf("text fields = (%v, %v, %q)", got.RawText, got.Message, got.Description)
	}
	wantCreatedAt, _ := time.Parse(time.RFC3339, createdAt)
	if !got.CreatedAt.Equal(wantCreatedAt) {
		t.Errorf("createdAt = %s, want %s", got.CreatedAt, wantCreatedAt)
	}
}

func TestClient_GetTransactionSettled(t *testing.T) {
	transactionID := uuid.MustParse("33333333-3333-4333-8333-333333333333")
	accountID := uuid.MustParse("22222222-2222-4222-8222-222222222222")
	settledAt := "2026-10-11T09:15:00+11:00"
	body := fmt.Sprintf(`{
		"data": {
			"id": %q,
			"attributes": {
				"status": "SETTLED",
				"description": "Restaurant",
				"rawText": "ORIGINAL MERCHANT TEXT",
				"message": "Dinner",
				"amount": {
					"currencyCode": "AUD",
					"value": "-48.50",
					"valueInBaseUnits": -4850
				},
				"createdAt": "2026-10-10T20:00:00+11:00",
				"settledAt": %q
			},
			"relationships": {
				"account": {"data": {"id": %q}}
			}
		}
	}`, transactionID.String(), settledAt, accountID.String())
	client := testClient(t, "/transactions/"+transactionID.String(), http.StatusOK, body)

	got, err := client.GetTransaction(t.Context(), transactionID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != TransactionSettled {
		t.Errorf("status = %q, want SETTLED", got.Status)
	}
	wantSettledAt, _ := time.Parse(time.RFC3339, settledAt)
	if got.SettledAt == nil || !got.SettledAt.Equal(wantSettledAt) {
		t.Errorf("settledAt = %v, want %s", got.SettledAt, wantSettledAt)
	}
	if got.RawText == nil || *got.RawText != "ORIGINAL MERCHANT TEXT" {
		t.Errorf("rawText = %v", got.RawText)
	}
	if got.Message == nil || *got.Message != "Dinner" {
		t.Errorf("message = %v", got.Message)
	}
	if got.Amount.ValueInBaseUnits != -4850 || got.AccountId != accountID {
		t.Errorf("amount/account = (%+v, %s)", got.Amount, got.AccountId)
	}
}

func TestClient_GetTransactionServerError(t *testing.T) {
	transactionID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	client := testClient(t, "/transactions/"+transactionID.String(), http.StatusInternalServerError, "")

	got, err := client.GetTransaction(t.Context(), transactionID)
	if err == nil {
		t.Fatal("expected an error for HTTP 500")
	}
	if got != (Transaction{}) {
		t.Fatalf("transaction = %+v, want zero value", got)
	}
}

func TestClient_GetTransactionMalformedResponse(t *testing.T) {
	transactionID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	client := testClient(t, "/transactions/"+transactionID.String(), http.StatusOK, `{"data":`)

	got, err := client.GetTransaction(t.Context(), transactionID)
	if err == nil {
		t.Fatal("expected a decoding error")
	}
	if got != (Transaction{}) {
		t.Fatalf("transaction = %+v, want zero value", got)
	}
}
