package microbatch

type JobResult[R any] struct {
	Result R
	Error  error
}

type Job[T, R any] func(input T) *JobResult[R]

type BatchProcessor[T, R any] interface {
	ProcessBatch(batch []Job[T, R])
}

/* type MicroBatch[T, R any] interface {

	// Start - spins up this micro-batch processor instance.
	// Provided Context value c is used as a stopping mechanism when context is canceled.
	Start(c context.Context) error

	// Submit - submits the job for processing.
	Submit(j Job[T, R]) error
} */
