// A simple example of using in-memory adapter
package main

import (
	"context"
	"log"

	kickplan "github.com/kickplan/sdk-go"
)

func main() {
	ctx := context.Background()
	client := kickplan.NewClient()

	const flag = "my-flag"

	b, err := client.GetBool(ctx, flag, false, nil)
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

	b, err = client.GetBool(ctx, flag, false, nil)
	if err != nil {
		log.Fatalf("failed to get flag: %v", err)
		return
	}

	log.Printf("my-flag: %v", b)
}
