package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const UpUrl = "https://api.up.com.au/api/v1/util/ping"

var UP_PAT string
var YNAB_PAT string

type Meta struct {
	Id     string `json:"id"`
	Status string `json:"statusEmoji"`
}

type UpResponse struct {
	Meta Meta `json:"meta"`
}

func doctor() error {
	fmt.Println("doctor ran")
	req, err := http.NewRequest("GET", UpUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+UP_PAT)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to deserialize body: %w", err)
	}

	var upResp UpResponse
	if err = json.Unmarshal(body, &upResp); err != nil {
		return fmt.Errorf("failed to deserialize upResp: %w", err)
	}
	fmt.Println(upResp.Meta.Id)
	return nil
}

func loadSecrets() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("load .env: %w", err)
	}
	UP_PAT = os.Getenv("UP_PAT")
	if UP_PAT == "" {
		return fmt.Errorf("UP_PAT is required")
	}
	YNAB_PAT = os.Getenv("YNAB_PAT")
	if YNAB_PAT == "" {
		return fmt.Errorf("YNAB_PAT is required")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s doctor\n", os.Args[0])
		os.Exit(2)
	}

	err := loadSecrets()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load secrets: %w", err)
		os.Exit(2)
	}

	command := os.Args[1]
	switch command {
	case "doctor":
		err := doctor()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Up Call Failed: %s\n", err.Error())
			os.Exit(2)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		fmt.Fprintf(os.Stderr, "Usage: %s doctor\n", os.Args[0])
		os.Exit(2)
	}
}
