package microbatch

type JobInput struct{}
type JobOutput struct{}

type testBatchProcessor[T *JobInput, R *JobOutput] struct{}

func newTestBatchProcessor() *testBatchProcessor[*JobInput, *JobOutput] {
	return &testBatchProcessor[*JobInput, *JobOutput]{}
}

func (bp *testBatchProcessor[T, R]) Process(batch []Job[*JobInput, *JobOutput]) error {
	for _, j := range batch {
		j(&JobInput{})
	}
	return nil
}

// Ensure testBatchProcessor conforms to the BatchProcessor interface.
var _ BatchProcessor[*JobInput, *JobOutput] = &testBatchProcessor[*JobInput, *JobOutput]{}
