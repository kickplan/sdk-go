package adapter

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kickplan/sdk-go/eval"
)

const (
	// DefaultEndpoint is the default endpoint for the Kickplan API.
	DefaultEndpoint = "https://api.kickplan.io"

	// DefaultUserAgent is the default user agent for HTTP client requests.
	DefaultUserAgent = "Kickplan Go SDK v0.1.0"

	// DefaultTimeout is the default timeout for HTTP client requests.
	DefaultTimeout = 5 * time.Second
)

// ErrFlagNotFound is returned when a flag is not found.
var ErrFlagNotFound = fmt.Errorf("FLAG_NOT_FOUND")

// Verify that Kickplan implements Adapter.
var _ Adapter = (*Kickplan)(nil)

// FeatureResolutionRequest represents a request body for the feature resolution endpoint.
type FeatureResolutionRequest struct {
	Context  eval.Context `json:"context"`
	Detailed bool         `json:"detailed"`
}

// FeatureResolutionResponse represents a response body for the feature resolution endpoint.
type FeatureResolutionResponse struct {
	ErrorCode string      `json:"error_code"`
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
}

// MetricUpdateRequest represents a request body for the metric update endpoint.
type MetricUpdateRequest struct {
	Context eval.Context `json:"context"`
	Value   int64        `json:"value"`
}

// AccountPlan represents an account plan used in CreateAccountRequest.
type AccountPlan struct {
	PlanKey string `json:"plan_key"`
}

// CreateAccountRequest represents a request body for the accounts endpoint.
type CreateAccountRequest struct {
	Key   string        `json:"key"`
	Name  string        `json:"name"`
	Plans []AccountPlan `json:"account_plans"`
}

// HTTPClient is an interface that defines the methods that a HTTP client must implement.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Kickplan is an adapter that uses Kickplan API for flags.
type Kickplan struct {
	client    HTTPClient
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
		client: &http.Client{
			Timeout: timeoutDuration,
		},
		endpoint:  endpoint,
		token:     token,
		userAgent: userAgent,
	}
}

// BooleanEvaluation returns the value of a boolean flag.
func (k *Kickplan) BooleanEvaluation(
	ctx context.Context,
	flag string,
	defaultValue bool,
	evalCtx eval.Context,
) (bool, error) {
	value, err := k.ResolveFeature(ctx, flag, defaultValue, evalCtx)
	if err != nil {
		return defaultValue, err
	}

	return genericResolve[bool](value, defaultValue)
}

// StringEvaluation returns the value of a string flag.
func (k *Kickplan) StringEvaluation(
	ctx context.Context,
	flag string,
	defaultValue string,
	evalCtx eval.Context,
) (string, error) {
	value, err := k.ResolveFeature(ctx, flag, defaultValue, evalCtx)
	if err != nil {
		return defaultValue, err
	}

	return genericResolve[string](value, defaultValue)
}

// Int64Evaluation returns the value of a int64 flag.
func (k *Kickplan) Int64Evaluation(
	ctx context.Context,
	flag string,
	defaultValue int64,
	evalCtx eval.Context,
) (int64, error) {
	value, err := k.ResolveFeature(ctx, flag, defaultValue, evalCtx)
	if err != nil {
		return defaultValue, err
	}

	return genericResolve[int64](value, defaultValue)
}

// ObjectEvaluation returns the value of a object flag.
func (k *Kickplan) ObjectEvaluation(
	ctx context.Context,
	flag string,
	defaultValue interface{},
	evalCtx eval.Context,
) (interface{}, error) {
	value, err := k.ResolveFeature(ctx, flag, defaultValue, evalCtx)
	if err != nil {
		return defaultValue, err
	}

	return value, nil
}

// ResolveFeature resolves a feature flag from the Kickplan API.
func (k *Kickplan) ResolveFeature(
	ctx context.Context,
	flag string,
	defaultValue interface{},
	evalCtx eval.Context,
) (interface{}, error) {
	url := fmt.Sprintf("%s/features/%s", k.endpoint, flag)
	body := FeatureResolutionRequest{
		Context:  evalCtx,
		Detailed: true,
	}

	resp, err := k.sendRequest(ctx, url, body)
	if err != nil {
		return defaultValue, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return defaultValue, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// read response body
	b, err := k.readResponseBody(resp)
	if err != nil {
		return defaultValue, fmt.Errorf("failed to read response body: %w", err)
	}

	// decode response
	var response FeatureResolutionResponse
	if err := json.Unmarshal(b, &response); err != nil {
		return defaultValue, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.ErrorCode != "" {
		if response.ErrorCode == "FLAG_NOT_FOUND" {
			return defaultValue, ErrFlagNotFound
		}

		return defaultValue, errors.New(response.ErrorCode)
	}

	return response.Value, nil
}

// SetBoolean sets the value of a boolean flag.
func (k *Kickplan) SetBoolean(_ context.Context, _ string, _ bool) error {
	return fmt.Errorf("not implemented")
}

// SetMetric sets the value of a metric.
func (k *Kickplan) SetMetric(ctx context.Context, metric string, value int64, evalCtx eval.Context) error {
	url := fmt.Sprintf("%s/metrics/%s/set", k.endpoint, metric)
	body := MetricUpdateRequest{
		Context: evalCtx,
		Value:   value,
	}

	resp, err := k.sendRequest(ctx, url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// IncMetric increments the value of a metric.
func (k *Kickplan) IncMetric(ctx context.Context, metric string, value int64, evalCtx eval.Context) error {
	url := fmt.Sprintf("%s/metrics/%s/increment", k.endpoint, metric)
	body := MetricUpdateRequest{
		Context: evalCtx,
		Value:   value,
	}

	resp, err := k.sendRequest(ctx, url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// DecMetric decrements the value of a metric.
func (k *Kickplan) DecMetric(ctx context.Context, metric string, value int64, evalCtx eval.Context) error {
	url := fmt.Sprintf("%s/metrics/%s/decrement", k.endpoint, metric)
	body := MetricUpdateRequest{
		Context: evalCtx,
		Value:   value,
	}

	resp, err := k.sendRequest(ctx, url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// CreateAccount creates an account and assigns plans to it.
func (k *Kickplan) CreateAccount(ctx context.Context, key, name string, planKeys ...string) error {
	url := fmt.Sprintf("%s/accounts", k.endpoint)
	body := CreateAccountRequest{
		Key:   key,
		Name:  name,
		Plans: []AccountPlan{},
	}

	for _, planKey := range planKeys {
		body.Plans = append(body.Plans, AccountPlan{
			PlanKey: planKey,
		})
	}

	resp, err := k.sendRequest(ctx, url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (k *Kickplan) sendRequest(ctx context.Context, url string, body interface{}) (*http.Response, error) {
	// encode body
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	k.setHeaders(req)

	// send request
	resp, err := k.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}

func (k *Kickplan) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.token))
	req.Header.Set("User-Agent", k.userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip;q=1.0,deflate;q=0.6,identity;q=0.3")
	req.Header.Set("Accept", "application/json")
}

func (k *Kickplan) readResponseBody(resp *http.Response) ([]byte, error) {
	reader := resp.Body

	if resp.Header.Get("Content-Encoding") == "gzip" {
		var err error
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	return io.ReadAll(reader)
}
