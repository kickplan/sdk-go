// Package adapter provides a way to work with flags and metrics
package adapter

import (
	"context"
	"fmt"
)

// InMemoryFlag structure represents a flag that is stored in memory
type InMemoryFlag struct {
	Value interface{}
}

// InMemory is an adapter that stores flags in memory
type InMemory struct {
	Flags map[string]InMemoryFlag
}

// NewInMemory returns a new InMemory adapter
func NewInMemory() *InMemory {
	return &InMemory{
		Flags: make(map[string]InMemoryFlag),
	}
}

// BooleanEvaluation returns the value of a boolean flag
func (i *InMemory) BooleanEvaluation(ctx context.Context, flag string, defaultValue bool) (bool, error) {
	memoryFlag, ok := i.find(flag)
	if !ok {
		return defaultValue, nil
	}

	return genericResolve[bool](memoryFlag.Value, defaultValue)
}

// SetBoolean sets the value of a boolean flag
func (i *InMemory) SetBoolean(ctx context.Context, flag string, value bool) error {
	i.Flags[flag] = InMemoryFlag{
		Value: value,
	}
	return nil
}

func (i *InMemory) find(flag string) (InMemoryFlag, bool) {
	memoryFlag, ok := i.Flags[flag]
	if !ok {
		return InMemoryFlag{}, false
	}

	return memoryFlag, true
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
