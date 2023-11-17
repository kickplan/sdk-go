// Package adapter provides a way to work with flags and metrics
package adapter

import (
	"context"

	"github.com/kickplan/sdk-go/eval"
)

// Adapter is an interface that defines the methods that a client adapter must implement.
type Adapter interface {
	BooleanEvaluation(ctx context.Context, flag string, defaultValue bool, evalCtx eval.Context) (bool, error)
	SetBoolean(ctx context.Context, flag string, value bool) error
}
