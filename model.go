package microbatch

type JobResult[R any] struct {
	Result R
	Error  error
}

type Job[T, R any] func(input T) *JobResult[R]

// Implement this interface fro you specific batch processing.
type BatchProcessor[T, R any] interface {

	// Processes given batch of jobs.
	Process(batch []Job[T, R]) error
}
