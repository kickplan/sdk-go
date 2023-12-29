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

func TestResolveFeatureObject(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewReader([]byte(`{
	"error_code": "",
	"key": "flag",
	"metadata": {},
	"value": {
		"key": "value",
		"list": [1, 2, 3]
	},
	"variant": null
	}
`))),
		}, nil
	}

	adapter := Kickplan{client: &mockClient{}}

	obj, err := adapter.ObjectEvaluation(context.TODO(), "flag", false, nil)
	if err != nil {
		t.Fatalf("failed to resolve feature: %v", err)
	}

	type Object struct {
		Key  string `json:"key"`
		List []int  `json:"list"`
	}

	expected := Object{
		Key:  "value",
		List: []int{1, 2, 3},
	}

	var actual Object

	// code below shows how to decode interface{} to struct
	// alternatively, could use https://github.com/mitchellh/mapstructure to decode interface{} to a struct
	// or json.Marshal `obj` and then json.Unmarshal to a struct
	if objRaw, ok := obj.(map[string]interface{}); ok {
		if key, ok := objRaw["key"]; ok {
			if actual.Key, ok = key.(string); !ok {
				t.Fatalf("failed to decode object key")
			}
		}

		if list, ok := objRaw["list"]; ok {
			if listr, ok := list.([]interface{}); ok {
				for _, item := range listr {
					actual.List = append(actual.List, int(item.(float64)))
				}
			}
		}
	} else {
		t.Fatalf("failed to decode object")
	}

	if expected.Key != actual.Key {
		t.Fatalf("expected object key to be %s, got %s", expected.Key, actual.Key)
	}

	if len(expected.List) != len(actual.List) {
		t.Fatalf("expected object list to be %v, got %v", expected.List, actual.List)
	}

	for i := range expected.List {
		if expected.List[i] != actual.List[i] {
			t.Fatalf("expected object list to be %v, got %v", expected.List, actual.List)
		}
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

func TestCreateAccount(t *testing.T) {
	DoFunc = func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != "https://api.domain.com/accounts" {
			t.Fatalf("expected request endpoint to be https://api.domain.com/accounts, got %s", req.URL.String())
		}

		if req.Method != http.MethodPost {
			t.Fatalf("expected request method to be POST, got %s", req.Method)
		}

		var body CreateAccountRequest
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if body.Key != "key" {
			t.Fatalf("expected request key to be \"key\", got %q", body.Key)
		}

		if body.Name != "Name" {
			t.Fatalf("expected request name to be \"Name\", got %q", body.Name)
		}

		if len(body.Plans) != 2 {
			t.Fatalf("expected request plans to be have 2 items, got %d", len(body.Plans))
		}

		if body.Plans[0].PlanKey != "free" {
			t.Fatalf("expected request plans to be [\"free\"], got %q", body.Plans)
		}

		if body.Plans[1].PlanKey != "pro" {
			t.Fatalf("expected request plans to be [\"pro\"], got %q", body.Plans)
		}

		return &http.Response{
			StatusCode: http.StatusAccepted,
			Body:       io.NopCloser(bytes.NewReader([]byte(``))),
		}, nil
	}

	adapter := Kickplan{client: &mockClient{}, endpoint: "https://api.domain.com"}

	err := adapter.CreateAccount(context.Background(), "key", "Name", "free", "pro")
	if err != nil {
		t.Fatal(err)
	}
}
