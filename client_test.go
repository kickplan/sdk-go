package kickplan

import (
	"context"
	"testing"

	"github.com/kickplan/sdk-go/adapter"
)

func TestDefaultAdapter(t *testing.T) {
	// Unset env vars to make sure that default adapter is used
	t.Setenv("KICKPLAN_ACCESS_TOKEN", "")

	client := NewClient()
	if client.adapter == nil {
		t.Fatalf("expected adapter to be set")
	}

	_, ok := client.adapter.(*adapter.InMemory)
	if !ok {
		t.Fatalf("expected adapter to be of type InMemory")
	}

	b, err := client.GetBool(context.TODO(), "my-flag", false)
	if err != nil {
		t.Fatalf("failed to get flag: %v", err)
	}

	if b != false {
		t.Fatalf("expected flag to be false")
	}

	err = client.SetBool(context.TODO(), "my-flag", true)
	if err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}

	b, err = client.GetBool(context.TODO(), "my-flag", false)
	if err != nil {
		t.Fatalf("failed to get flag: %v", err)
	}

	if b != true {
		t.Fatalf("expected flag to be true")
	}
}
