package microbatch

import "time"

// Option allows to customise configuration.
type Option[T, R any] func(*MicroBatch[T, R])

// Limit sets the max size of the micro-batch.
func Limit[T, R any](v int) Option[T, R] {
	return func(mb *MicroBatch[T, R]) {
		mb.limit = v
	}
}

// Cycle sets the max time duration of the micro-batch.
func Cycle[T, R any](v time.Duration) Option[T, R] {
	return func(mb *MicroBatch[T, R]) {
		mb.cycle = v
	}
}
