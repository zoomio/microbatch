package microbatch

import "errors"

var (
	ErrBatchProcessorRequired = errors.New("BatchProcessor is required")
	ErrIsNotRunning           = errors.New("micro-batch processor is not running")
	ErrSameRunningState       = errors.New("micro-batch processor isRunning is already set")
	ErrNilJob                 = errors.New("recieved nil job")
)
