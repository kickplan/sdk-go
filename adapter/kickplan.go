package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	// DefaultEndpoint is the default endpoint for the Kickplan API.
	DefaultEndpoint = "https://api.kickplan.io"

	// DefaultUserAgent is the default user agent for HTTP client requests.
	DefaultUserAgent = "Kickplan Go SDK v0.1.0"

	// DefaultTimeout is the default timeout for HTTP client requests.
	DefaultTimeout = 5 * time.Second
)

// Verify that Kickplan implements Adapter.
var _ Adapter = (*Kickplan)(nil)

// FeatureResolutionRequest represents a request body for the feature resolution endpoint.
type FeatureResolutionRequest struct {
	Context map[string]interface{} `json:"context"`
}

// FeatureResolutionResponse represents a response body for the feature resolution endpoint.
type FeatureResolutionResponse struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Kickplan is an adapter that uses Kickplan API for flags.
type Kickplan struct {
	client    http.Client
	endpoint  string
	token     string
	userAgent string
}

// NewKickplan returns a new Kickplan adapter.
func NewKickplan(
	endpoint string,
	token string,
	userAgent string,
	timeout string,
) *Kickplan {
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}

	if userAgent == "" {
		userAgent = DefaultUserAgent
	}

	timeoutDuration := DefaultTimeout
	if timeout != "" {
		var err error
		timeoutDuration, err = time.ParseDuration(timeout)
		if err != nil {
			log.Printf("WARN failed to parse timeout duration %q: %v", timeout, err)
			timeoutDuration = DefaultTimeout
		}
	}

	return &Kickplan{
		client: http.Client{
			Timeout: timeoutDuration,
		},
		endpoint:  endpoint,
		token:     token,
		userAgent: userAgent,
	}
}

// BooleanEvaluation returns the value of a boolean flag.
func (k *Kickplan) BooleanEvaluation(ctx context.Context, flag string, defaultValue bool) (bool, error) {
	value, err := k.ResolveFeature(ctx, flag, defaultValue)
	if err != nil {
		return defaultValue, err
	}

	return genericResolve[bool](value, defaultValue)
}

// ResolveFeature resolves a feature flag from the Kickplan API.
func (k *Kickplan) ResolveFeature(
	ctx context.Context,
	flag string,
	defaultValue interface{},
	// todo: pass evaluation context
) (interface{}, error) {
	body := FeatureResolutionRequest{
		Context: map[string]interface{}{}, // todo: pass evaluation context
	}

	// encode body
	b, err := json.Marshal(body)
	if err != nil {
		return defaultValue, fmt.Errorf("failed to encode request body: %w", err)
	}

	url := fmt.Sprintf("%s/flags/%s", k.endpoint, flag)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return defaultValue, fmt.Errorf("failed to create request: %w", err)
	}

	k.setHeaders(req)

	// send request
	resp, err := k.client.Do(req)
	if err != nil {
		return defaultValue, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// decode response
	var response FeatureResolutionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return defaultValue, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Value, nil
}

// SetBoolean sets the value of a boolean flag.
func (k *Kickplan) SetBoolean(_ context.Context, _ string, _ bool) error {
	return fmt.Errorf("not implemented")
}

func (k *Kickplan) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.token))
	req.Header.Set("User-Agent", k.userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip;q=1.0,deflate;q=0.6,identity;q=0.3")
	req.Header.Set("Accept", "application/json")
}
