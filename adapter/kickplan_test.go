package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/kickplan/sdk-go/eval"
)

func TestResolveFeature(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		// Verify request
		if req.URL.String() != "https://api.domain.com/features/flag" {
			t.Fatalf("expected request endpoint to be https://api.domain.com/features/flag, got %s", req.URL.String())
		}

		if req.Method != http.MethodPost {
			t.Fatalf("expected request method to be POST, got %s", req.Method)
		}

		if req.Header.Get("User-Agent") != "user-agent" {
			t.Fatalf("expected request user agent to be user-agent, got %s", req.Header.Get("User-Agent"))
		}

		if req.Header.Get("Authorization") != "Bearer token" {
			t.Fatalf("expected request authorization to be `Bearer token`, got %s", req.Header.Get("Authorization"))
		}

		if req.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("expected request content type to be application/json, got %s", req.Header.Get("Content-Type"))
		}

		var body FeatureResolutionRequest
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if body.Detailed != true {
			t.Fatalf("expected request detailed to be `true`, got %v", body.Detailed)
		}

		if body.Context["account_id"] != "account" {
			t.Fatalf("expected request context account_id to be account, got %s", body.Context["account_id"])
		}

		// Return response
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewReader([]byte(`{
	"error_code": "",
	"key": "flag",
	"metadata": {},
	"value": true,
	"variant": null
	}
`))),
		}, nil
	}

	adapter := Kickplan{
		client:    &mockClient{},
		endpoint:  "https://api.domain.com",
		token:     "token",
		userAgent: "user-agent",
	}

	result, err := adapter.BooleanEvaluation(context.TODO(), "flag", false, eval.Context{
		"account_id": "account",
	})
	if err != nil {
		t.Fatalf("failed to resolve feature: %v", err)
	}

	if result != true {
		t.Fatalf("expected result to be true")
	}
}

func TestResolveFeatureNotFound(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewReader([]byte(`{
	"error_code": "FLAG_NOT_FOUND",
	"key": "flag",
	"metadata": {},
	"value": null,
	"variant": null
	}
`))),
		}, nil
	}

	adapter := Kickplan{client: &mockClient{}}

	_, err := adapter.BooleanEvaluation(context.TODO(), "flag", false, nil)
	if !errors.Is(err, ErrFlagNotFound) {
		t.Fatalf("expected error to be ErrFlagNotFound, got %v", err)
	}
}

func TestMetricSet(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != "https://api.domain.com/metrics/metric/set" {
			t.Fatalf("expected request endpoint to be https://api.domain.com/metrics/set, got %s", req.URL.String())
		}

		if req.Method != http.MethodPost {
			t.Fatalf("expected request method to be POST, got %s", req.Method)
		}

		var body MetricUpdateRequest
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if body.Value != 1 {
			t.Fatalf("expected request value to be 1, got %d", body.Value)
		}

		if body.Context["account_id"] != "account" {
			t.Fatalf("expected request context account_id to be account, got %s", body.Context["account_id"])
		}

		return &http.Response{
			StatusCode: http.StatusAccepted,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		}, nil
	}

	adapter := Kickplan{client: &mockClient{}, endpoint: "https://api.domain.com"}

	if err := adapter.SetMetric(context.TODO(), "metric", 1, eval.Context{
		"account_id": "account",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestMetricIncrement(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != "https://api.domain.com/metrics/metric/increment" {
			t.Fatalf("expected request endpoint to be https://api.domain.com/metrics/increment, got %s", req.URL.String())
		}

		if req.Method != http.MethodPost {
			t.Fatalf("expected request method to be POST, got %s", req.Method)
		}

		var body MetricUpdateRequest
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if body.Value != 20 {
			t.Fatalf("expected request value to be 20, got %d", body.Value)
		}

		return &http.Response{
			StatusCode: http.StatusAccepted,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		}, nil
	}

	adapter := Kickplan{client: &mockClient{}, endpoint: "https://api.domain.com"}

	if err := adapter.IncMetric(context.TODO(), "metric", 20, eval.Context{
		"account_id": "account",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestMetricDecrement(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != "https://api.domain.com/metrics/metric/decrement" {
			t.Fatalf("expected request endpoint to be https://api.domain.com/metrics/decrement, got %s", req.URL.String())
		}

		if req.Method != http.MethodPost {
			t.Fatalf("expected request method to be POST, got %s", req.Method)
		}

		var body MetricUpdateRequest
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if body.Value != 40 {
			t.Fatalf("expected request value to be 40, got %d", body.Value)
		}

		return &http.Response{
			StatusCode: http.StatusAccepted,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		}, nil
	}

	adapter := Kickplan{client: &mockClient{}, endpoint: "https://api.domain.com"}

	if err := adapter.DecMetric(context.TODO(), "metric", 40, eval.Context{
		"account_id": "account",
	}); err != nil {
		t.Fatal(err)
	}
}
