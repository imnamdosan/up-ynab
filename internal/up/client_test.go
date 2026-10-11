package up

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestClient_ListAccounts(t *testing.T) {
	accountId, err := uuid.NewRandom()
	testToken := "test-token"
	if err != nil {
		t.Errorf("Failed to generate uid")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts" {
			t.Errorf("path = %q, want /accounts", r.URL.Path)
		}

		if value := r.Header.Get("Authorization"); value != "Bearer "+testToken {
			t.Errorf("Authorization = %q", value)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"data":[{"id":%q}]}`, accountId)
	}))
	defer server.Close()

	client := NewClientWithBaseUrl(testToken, server.Client(), server.URL+"/")
	resp, err := client.ListAccounts(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(resp) != 1 || resp[0] != accountId {
		t.Fatalf("accounts = %+v", resp)
	}
}

func TestClient_ListAccountsServerError(t *testing.T) {
	testToken := "test-token"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts" {
			t.Errorf("path = %q, want /accounts", r.URL.Path)
		}

		if value := r.Header.Get("Authorization"); value != "Bearer "+testToken {
			t.Errorf("Authorization = %q", value)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClientWithBaseUrl(testToken, server.Client(), server.URL+"/")
	resp, err := client.ListAccounts(t.Context())
	if err == nil || len(resp) != 0 {
		t.Fatal("didn't process fail status", resp)
	}
}
