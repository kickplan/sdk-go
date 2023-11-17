package main

import (
	"context"
	"log"
	"os"

	kickplan "github.com/kickplan/sdk-go"
	"github.com/kickplan/sdk-go/adapter"
)

func main() {
	endpoint := os.Getenv("KICKPLAN_ENDPOINT")
	token := os.Getenv("KICKPLAN_ACCESS_TOKEN")
	userAgent := os.Getenv("KICKPLAN_USER_AGENT")
	timeout := os.Getenv("KICKPLAN_TIMEOUT") // e.g. "5s"

	ctx := context.Background()
	client := kickplan.NewClient(
		// Passign an adapter explicitly here,
		// but by default client will use Kickplan adapter if env vars are set
		kickplan.WithAdapter(
			adapter.NewKickplan(endpoint, token, userAgent, timeout),
		),
	)

	const flag = "my-flag"

	b, err := client.GetBool(ctx, flag, false)
	if err != nil {
		log.Fatalf("failed to get flag: %v", err)
		return
	}

	log.Printf("my-flag: %v", b)

	err = client.SetBool(ctx, flag, true)
	if err != nil {
		log.Fatalf("failed to set flag: %v", err)
		return
	}

	log.Printf("updated my-flag")

	b, err = client.GetBool(ctx, flag, false)
	if err != nil {
		log.Fatalf("failed to get flag: %v", err)
		return
	}

	log.Printf("my-flag: %v", b)
}
