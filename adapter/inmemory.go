package adapter

import (
	"context"

	"github.com/kickplan/sdk-go/eval"
)

// Verify that InMemory implements Adapter.
var _ Adapter = (*InMemory)(nil)

// InMemoryFlag structure represents a flag that is stored in memory.
type InMemoryFlag struct {
	Value interface{}
}

// InMemory is an adapter that stores flags in memory.
type InMemory struct {
	Flags   map[string]InMemoryFlag
	Metrics map[string]int64
}

// NewInMemory returns a new InMemory adapter.
func NewInMemory() *InMemory {
	return &InMemory{
		Flags:   make(map[string]InMemoryFlag),
		Metrics: make(map[string]int64),
	}
}

// BooleanEvaluation returns the value of a boolean flag.
func (i *InMemory) BooleanEvaluation(
	_ context.Context,
	flag string,
	defaultValue bool,
	_ eval.Context,
) (bool, error) {
	memoryFlag, ok := i.find(flag)
	if !ok {
		return defaultValue, nil
	}

	return genericResolve[bool](memoryFlag.Value, defaultValue)
}

// StringEvaluation returns the value of a string flag.
func (i *InMemory) StringEvaluation(
	_ context.Context,
	flag string,
	defaultValue string,
	_ eval.Context,
) (string, error) {
	memoryFlag, ok := i.find(flag)
	if !ok {
		return defaultValue, nil
	}

	return genericResolve[string](memoryFlag.Value, defaultValue)
}

// Int64Evaluation returns the value of a int64 flag.
func (i *InMemory) Int64Evaluation(
	_ context.Context,
	flag string,
	defaultValue int64,
	_ eval.Context,
) (int64, error) {
	memoryFlag, ok := i.find(flag)
	if !ok {
		return defaultValue, nil
	}

	return genericResolve[int64](memoryFlag.Value, defaultValue)
}

// SetBoolean sets the value of a boolean flag.
func (i *InMemory) SetBoolean(_ context.Context, flag string, value bool) error {
	i.Flags[flag] = InMemoryFlag{
		Value: value,
	}
	return nil
}

// SetMetric sets the value of a metric.
func (i *InMemory) SetMetric(_ context.Context, metric string, value int64, _ eval.Context) error {
	i.Metrics[metric] = value
	return nil
}

// IncMetric increments the value of a metric.
func (i *InMemory) IncMetric(_ context.Context, metric string, value int64, _ eval.Context) error {
	i.Metrics[metric] += value
	return nil
}

// DecMetric decrements the value of a metric.
func (i *InMemory) DecMetric(_ context.Context, metric string, value int64, _ eval.Context) error {
	i.Metrics[metric] -= value
	return nil
}

func (i *InMemory) find(flag string) (InMemoryFlag, bool) {
	memoryFlag, ok := i.Flags[flag]
	if !ok {
		return InMemoryFlag{}, false
	}

	return memoryFlag, true
}
