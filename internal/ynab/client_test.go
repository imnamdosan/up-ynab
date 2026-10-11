package ynab

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func testClient(t *testing.T, path string, status int, body string) Client {
	t.Helper()
	const token = "test-token"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != path {
			t.Errorf("path = %q, want %q", r.URL.Path, path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+token {
			t.Errorf("Authorization = %q, want Bearer token", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	t.Cleanup(server.Close)
	return NewClientWithBaseUrl(token, server.Client(), server.URL+"/")
}

func TestClient_ListPlans(t *testing.T) {
	planID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	client := testClient(t, "/plans", http.StatusOK,
		fmt.Sprintf(`{"data":{"plans":[{"id":%q}]}}`, planID.String()))

	plans, err := client.ListPlans(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 1 || plans[0] != planID {
		t.Fatalf("plans = %v, want [%s]", plans, planID)
	}
}

func TestClient_ListPlansServerError(t *testing.T) {
	client := testClient(t, "/plans", http.StatusInternalServerError, "")

	plans, err := client.ListPlans(t.Context())
	if err == nil {
		t.Fatal("expected an error for HTTP 500")
	}
	if plans != nil {
		t.Fatalf("plans = %v, want nil", plans)
	}
}

func TestClient_ListAccounts(t *testing.T) {
	planID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	accountID := uuid.MustParse("22222222-2222-4222-8222-222222222222")
	client := testClient(t, "/plans/"+planID.String()+"/accounts", http.StatusOK,
		fmt.Sprintf(`{"data":{"accounts":[{"id":%q}]}}`, accountID.String()))

	accounts, err := client.ListAccounts(t.Context(), planID)
	if err != nil {
		t.Fatal(err)
	}
	if len(accounts) != 1 || accounts[0] != accountID {
		t.Fatalf("accounts = %v, want [%s]", accounts, accountID)
	}
}

func TestClient_ListAccountsServerError(t *testing.T) {
	planID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	client := testClient(t, "/plans/"+planID.String()+"/accounts", http.StatusInternalServerError, "")

	accounts, err := client.ListAccounts(t.Context(), planID)
	if err == nil {
		t.Fatal("expected an error for HTTP 500")
	}
	if accounts != nil {
		t.Fatalf("accounts = %v, want nil", accounts)
	}
}
