package microbatch

import "sync/atomic"

type JobInput struct{}
type JobOutput struct{}

type testBatchProcessor[T *JobInput, R *JobOutput] struct {
	jobsDone atomic.Int32
}

func newTestBatchProcessor() *testBatchProcessor[*JobInput, *JobOutput] {
	return &testBatchProcessor[*JobInput, *JobOutput]{}
}

func (bp *testBatchProcessor[T, R]) Process(batch []Job[*JobInput, *JobOutput]) error {
	for _, j := range batch {
		j(&JobInput{})
		bp.incrementJobs()
	}
	return nil
}

func (bp *testBatchProcessor[T, R]) incrementJobs() {
	for {
		// Load current balance atomically
		current := bp.jobsDone.Load()

		// Calculate new balance
		new := current + 1

		// Try to update balance atomically
		if bp.jobsDone.CompareAndSwap(current, new) {
			return
		}
	}
}

func (bp *testBatchProcessor[T, R]) getJobsDone() int {
	return int(bp.jobsDone.Load())
}

// Ensure testBatchProcessor conforms to the BatchProcessor interface.
var _ BatchProcessor[*JobInput, *JobOutput] = &testBatchProcessor[*JobInput, *JobOutput]{}
