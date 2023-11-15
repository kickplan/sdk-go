// Package kickplan provides a client to evaluate feature flags and work with metrics
package kickplan

import (
	"context"

	"github.com/kickplan/sdk-go/adapter"
)

// EvalContext is a map that represents the context of an evaluation
type EvalContext map[string]interface{}

// Adapter is an interface that defines the methods that a client adapter must implement
type Adapter interface {
	BooleanEvaluation(ctx context.Context, flag string, defaultValue bool) (bool, error)
	SetBoolean(ctx context.Context, flag string, value bool) error
}

// Client is a Kickplan client
type Client struct {
	adapter Adapter
}

// Option is a function that configures a Client
type Option func(*Client) error

// NewClient returns a new Kickplan client
func NewClient(opt ...Option) *Client {
	c := new(Client)
	for _, o := range opt {
		o(c)
	}

	if c.adapter == nil {
		c.adapter = adapter.NewInMemory()
	}

	return c
}

// WithAdapter sets the provider for the client
func WithAdapter(adapter Adapter) Option {
	return func(c *Client) error {
		c.adapter = adapter
		return nil
	}
}

// GetBool returns a boolean flag
func (c *Client) GetBool(ctx context.Context, flag string, defaultValue bool) (bool, error) {
	return c.adapter.BooleanEvaluation(ctx, flag, defaultValue)
}

// SetBool sets a boolean flag
func (c *Client) SetBool(ctx context.Context, flag string, value bool) error {
	return c.adapter.SetBoolean(ctx, flag, value)
}
