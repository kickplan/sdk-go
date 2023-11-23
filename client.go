// Package kickplan provides a client to evaluate feature flags and work with metrics
package kickplan

import (
	"context"
	"fmt"
	"os"

	"github.com/kickplan/sdk-go/adapter"
	"github.com/kickplan/sdk-go/eval"
)

// Client is a Kickplan client.
type Client struct {
	adapter adapter.Adapter
}

// Option is a function that configures a Client.
type Option func(*Client) error

// NewClient returns a new Kickplan client.
func NewClient(opt ...Option) *Client {
	c := new(Client)
	for _, o := range opt {
		err := o(c)
		if err != nil {
			panic(fmt.Sprintf("error applying option: %v", err))
		}
	}

	if os.Getenv("KICKPLAN_ACCESS_TOKEN") != "" {
		c.adapter = adapter.NewKickplan(
			os.Getenv("KICKPLAN_ENDPOINT"),
			os.Getenv("KICKPLAN_ACCESS_TOKEN"),
			os.Getenv("KICKPLAN_USER_AGENT"),
			os.Getenv("KICKPLAN_TIMEOUT"),
		)
	}

	if c.adapter == nil {
		c.adapter = adapter.NewInMemory()
	}

	return c
}

// WithAdapter sets the provider for the client.
func WithAdapter(a adapter.Adapter) Option {
	return func(c *Client) error {
		c.adapter = a
		return nil
	}
}

// GetBool returns a boolean flag.
func (c *Client) GetBool(
	ctx context.Context,
	flag string,
	defaultValue bool,
	evalCtx eval.Context,
) (bool, error) {
	return c.adapter.BooleanEvaluation(ctx, flag, defaultValue, evalCtx)
}

// GetInt64 returns a float64 flag.
func (c *Client) GetInt64(
	ctx context.Context,
	flag string,
	defaultValue int64,
	evalCtx eval.Context,
) (int64, error) {
	return c.adapter.Int64Evaluation(ctx, flag, defaultValue, evalCtx)
}

// GetString returns a string flag.
func (c *Client) GetString(
	ctx context.Context,
	flag string,
	defaultValue string,
	evalCtx eval.Context,
) (string, error) {
	return c.adapter.StringEvaluation(ctx, flag, defaultValue, evalCtx)
}

// GetObject returns a object flag.
func (c *Client) GetObject(
	ctx context.Context,
	flag string,
	defaultValue interface{},
	evalCtx eval.Context,
) (interface{}, error) {
	return c.adapter.ObjectEvaluation(ctx, flag, defaultValue, evalCtx)
}

// SetBool sets a boolean flag.
func (c *Client) SetBool(ctx context.Context, flag string, value bool) error {
	return c.adapter.SetBoolean(ctx, flag, value)
}

// SetMetric sets a metric.
func (c *Client) SetMetric(
	ctx context.Context,
	metric string,
	value int64,
	evalCtx eval.Context,
) error {
	return c.adapter.SetMetric(ctx, metric, value, evalCtx)
}

// IncMetric increments a metric.
func (c *Client) IncMetric(
	ctx context.Context,
	metric string,
	value int64,
	evalCtx eval.Context,
) error {
	return c.adapter.IncMetric(ctx, metric, value, evalCtx)
}

// DecMetric decrements a metric.
func (c *Client) DecMetric(
	ctx context.Context,
	metric string,
	value int64,
	evalCtx eval.Context,
) error {
	return c.adapter.DecMetric(ctx, metric, value, evalCtx)
}
