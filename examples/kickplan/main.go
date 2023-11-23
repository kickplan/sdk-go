// A simple example of using Kickplan adapter
package main

import (
	"context"
	"errors"
	"log"
	"os"

	kickplan "github.com/kickplan/sdk-go"
	"github.com/kickplan/sdk-go/adapter"
	"github.com/kickplan/sdk-go/eval"
)

func main() {
	endpoint := os.Getenv("KICKPLAN_ENDPOINT")
	token := os.Getenv("KICKPLAN_ACCESS_TOKEN")
	userAgent := os.Getenv("KICKPLAN_USER_AGENT")
	timeout := os.Getenv("KICKPLAN_TIMEOUT") // e.g. "5s"

	account := os.Getenv("KICKPLAN_ACCOUNT") // one of the accounts UUID to use for evaluation

	ctx := context.Background()
	client := kickplan.NewClient(
		// Passign an adapter explicitly here,
		// but by default client will use Kickplan adapter if env vars are set
		kickplan.WithAdapter(
			adapter.NewKickplan(endpoint, token, userAgent, timeout),
		),
	)

	const flag = "my-flag"

	b, err := client.GetBool(ctx, flag, false, eval.Context{
		"account_id": account,
		"detailed":   true,
	})
	if err != nil {
		if errors.Is(err, adapter.ErrFlagNotFound) {
			log.Printf("flag %q not found", flag)
			return
		}
		log.Fatalf("failed to get flag: %v", err)
	}

	log.Printf("my-flag: %v", b)
}
