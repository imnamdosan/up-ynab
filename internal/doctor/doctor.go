package doctor

import (
	"UpYnab/internal/config"
	"UpYnab/internal/up"
	"UpYnab/internal/ynab"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bytedance/gg/gslice"
)

func Run(ctx context.Context, cfg config.Config, out io.Writer) error {
	httpClient := &http.Client{Timeout: 15 * time.Second}
	upClient := up.NewClient(cfg.UpToken, httpClient)
	ynabClient := ynab.NewClient(cfg.YnabToken, httpClient)

	fmt.Fprintln(out, "doctor ran")

	upAccounts, err := upClient.ListAccounts(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Up account IDs:", upAccounts)
	if !gslice.Contains(upAccounts, cfg.UpAccountId) {
		return fmt.Errorf("UP_ACCOUNT_ID %s was not found", cfg.UpAccountId)
	}

	ynabPlans, err := ynabClient.ListPlans(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "YNAB plan IDs:", ynabPlans)
	if !gslice.Contains(ynabPlans, cfg.YnabPlanId) {
		return fmt.Errorf("YNAB_PLAN_ID %s was not found", cfg.YnabPlanId)
	}

	ynabAccounts, err := ynabClient.ListAccounts(ctx, cfg.YnabPlanId)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "YNAB account IDs:", ynabAccounts)
	if !gslice.Contains(ynabAccounts, cfg.YnabAccountId) {
		return fmt.Errorf("YNAB_ACCOUNT_ID %s was not found in plan %s", cfg.YnabAccountId, cfg.YnabPlanId)
	}
	return nil
}
