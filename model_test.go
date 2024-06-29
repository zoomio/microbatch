package microbatch

type JobInput struct{}
type JobOutput struct{}

type testBatchProcessor[T *JobInput, R *JobOutput] struct{}

func newTestBatchProcessor() *testBatchProcessor[*JobInput, *JobOutput] {
	return &testBatchProcessor[*JobInput, *JobOutput]{}
}

func (bp *testBatchProcessor[T, R]) ProcessBatch(batch []Job[*JobInput, *JobOutput]) {
	for _, j := range batch {
		j(&JobInput{})
	}
}

// Ensure testBatchProcessor conforms to the BatchProcessor interface.
var _ BatchProcessor[*JobInput, *JobOutput] = &testBatchProcessor[*JobInput, *JobOutput]{}
