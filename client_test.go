package kickplan

import (
	"testing"

	"github.com/kickplan/sdk-go/adapter"
)

func TestDefaultAdapter(t *testing.T) {
	client := NewClient()
	if client.adapter == nil {
		t.Fatalf("expected adapter to be set")
	}

	_, ok := client.adapter.(*adapter.InMemory)
	if !ok {
		t.Fatalf("expected adapter to be of type InMemory")
	}

	b, err := client.GetBool(nil, "my-flag", false)
	if err != nil {
		t.Fatalf("failed to get flag: %v", err)
	}

	if b != false {
		t.Fatalf("expected flag to be false")
	}

	err = client.SetBool(nil, "my-flag", true)
	if err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}

	b, err = client.GetBool(nil, "my-flag", false)
	if err != nil {
		t.Fatalf("failed to get flag: %v", err)
	}

	if b != true {
		t.Fatalf("expected flag to be true")
	}
}