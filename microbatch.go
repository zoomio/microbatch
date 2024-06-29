package microbatch

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultBatchLimit   = 10
	defaultBatchCadence = 5 * time.Second
)

// MicroBatch - is the main entry point handles all the processing.
type MicroBatch[T, R any] struct {
	bp BatchProcessor[T, R]

	// config
	limit int
	cycle time.Duration

	// data
	ch    chan Job[T, R]
	batch []Job[T, R]

	// state
	isRunning atomic.Bool
	mu        sync.Mutex
}

// New - creates a new instance of the micro-batch processor.
// IMPORTANT: Keep in mind that this constructor doesn't start the actual processing,
// use #Start to start the processing or starting constructor #NewRunning
func New[T, R any](bp BatchProcessor[T, R], options ...Option[T, R]) (*MicroBatch[T, R], error) {
	if bp == nil {
		return nil, ErrBatchProcessorRequired
	}

	mb := &MicroBatch[T, R]{
		bp:    bp,
		limit: defaultBatchLimit,
		cycle: defaultBatchCadence,
	}

	// apply custom configurations
	for _, option := range options {
		option(mb)
	}

	mb.ch = make(chan Job[T, R], mb.limit)
	mb.setNewBatch()

	return mb, nil
}

// NewRunning - creates an already running instance of the micro-batch processor.
// Provided Context value c is used as a stopping mechanism when context is canceled.
// IMPORTANT: if you just want to create an instance and not start it yet use #New.
func NewRunning[T, R any](c context.Context, bp BatchProcessor[T, R], options ...Option[T, R]) (*MicroBatch[T, R], error) {
	mb, err := New(bp, options...)
	if err != nil {
		return nil, err
	}
	go mb.Start(c)
	return mb, nil
}

// Start - synchronously spins up this micro-batch processor instance and blocks until it's done.
// Provided Context value c is used as a stopping mechanism when context is canceled.
func (mb *MicroBatch[T, R]) Start(c context.Context) error {
	err := mb.setIsRunning(true)
	if err != nil {
		return err
	}

	// cycle is gonna be happening in here
	ticker := time.NewTicker(mb.cycle)
	for {
		select {
		case <-c.Done():
			slog.Debug("context has been closed, stopping started import watcher")

			// make sure no more submits
			err := mb.setIsRunning(false)

			// need to complete already submitted jobs
			mb.processBatch()

			return err
		case <-ticker.C:
			mb.processBatch()
		case j := <-mb.ch:
			mb.appendJob(j)
		}
	}
}

// Submit - synchronously submits the job for processing,
// meaning that it waits in case when the channel is full.
func (mb *MicroBatch[T, R]) Submit(j Job[T, R]) error {
	if !mb.isRunning.Load() {
		return ErrIsNotRunning
	}
	if j == nil {
		return ErrNilJob
	}
	mb.ch <- j
	return nil
}

func (mb *MicroBatch[T, R]) setNewBatch() {
	mb.batch = make([]Job[T, R], 0, mb.limit)
}

func (mb *MicroBatch[T, R]) setIsRunning(v bool) error {
	// synchronize with other potential writers
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if mb.isRunning.Load() == v {
		return ErrSameRunningState
	}
	mb.isRunning.Store(v)
	return nil
}

func (mb *MicroBatch[T, R]) appendJob(j Job[T, R]) {
	slog.Debug("appending job")
	mb.batch = append(mb.batch, j)
	if len(mb.batch) >= mb.limit {
		slog.Debug("batch size is at the limit, process immediately")
		mb.processBatch()
	}
}

func (mb *MicroBatch[T, R]) processBatch() {
	if len(mb.batch) == 0 {
		slog.Debug("batch is empty nothing to do")
		return
	}

	// copy processed batch and clear the internal to avoid side effects
	b := make([]Job[T, R], len(mb.batch))
	copy(b, mb.batch)
	mb.setNewBatch()

	slog.Debug("processing batch", "batch_size", len(b))
	mb.bp.ProcessBatch(b)
}
