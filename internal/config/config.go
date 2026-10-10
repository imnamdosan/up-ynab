package config

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Config struct {
	UpToken       string
	YnabToken     string
	UpAccountId   uuid.UUID
	YnabPlanId    uuid.UUID
	YnabAccountId uuid.UUID
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}
	config := Config{
		UpToken:   os.Getenv("UP_PAT"),
		YnabToken: os.Getenv("YNAB_PAT"),
	}
	if config.UpToken == "" {
		return Config{}, fmt.Errorf("UP_PAT is required")
	}
	if config.YnabToken == "" {
		return Config{}, fmt.Errorf("YNAB_PAT is required")
	}
	var err error
	config.YnabPlanId, err = uuid.Parse(os.Getenv("YNAB_PLAN_ID"))
	if err != nil {
		return Config{}, fmt.Errorf("YNAB_PLAN_ID must be a UUID: %w", err)
	}
	config.YnabAccountId, err = uuid.Parse(os.Getenv("YNAB_ACCOUNT_ID"))
	if err != nil {
		return Config{}, fmt.Errorf("YNAB_ACCOUNT_ID must be a UUID: %w", err)
	}
	config.UpAccountId, err = uuid.Parse(os.Getenv("UP_ACCOUNT_ID"))
	if err != nil {
		return Config{}, fmt.Errorf("UP_ACCOUNT_ID must be a UUID: %w", err)
	}
	return config, nil
}
