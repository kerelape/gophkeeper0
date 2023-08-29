package deferred

import (
	"context"
	"sync/atomic"
)

// Deferred is a deferred value container.
type Deferred[T any] struct {
	value    T
	hasValue atomic.Bool
}

// Get returns the value stored in the Deferred or wait until it is set.
func (d *Deferred[T]) Get(ctx context.Context) (value T, err error) {
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		default:
			if d.hasValue.Load() {
				return d.value, nil
			}
		}
	}
}

// Set sets the value to the Deferred.
func (d *Deferred[T]) Set(value T) {
	if d.hasValue.Load() {
		panic("value already set")
	}
	d.value = value
	d.hasValue.Swap(true)
}
