// Package adapter provides a way to work with flags and metrics
package adapter

import (
	"context"
	"fmt"

	"github.com/kickplan/sdk-go/eval"
)

// Adapter is an interface that defines the methods that a client adapter must implement.
type Adapter interface {
	BooleanEvaluation(ctx context.Context, flag string, defaultValue bool, evalCtx eval.Context) (bool, error)
	StringEvaluation(ctx context.Context, flag string, defaultValue string, evalCtx eval.Context) (string, error)
	Int64Evaluation(ctx context.Context, flag string, defaultValue int64, evalCtx eval.Context) (int64, error)

	SetBoolean(ctx context.Context, flag string, value bool) error

	SetMetric(ctx context.Context, metric string, value int64, evalCtx eval.Context) error
	IncMetric(ctx context.Context, metric string, value int64, evalCtx eval.Context) error
	DecMetric(ctx context.Context, metric string, value int64, evalCtx eval.Context) error
}

func genericResolve[T any](flag interface{}, defaultValue T) (T, error) {
	if flag == nil {
		return defaultValue, nil
	}

	if v, ok := flag.(T); ok {
		return v, nil
	}

	return defaultValue, fmt.Errorf("type assertion failed")
}
