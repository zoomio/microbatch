package microbatch

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_New_errors_on_nil_bp(t *testing.T) {
	mb, err := New[any, any](nil)

	assert.Nil(t, mb)
	assert.NotNil(t, err)
	assert.Equal(t, ErrBatchProcessorRequired, err)
}

func Test_New_creates_MicroBatch(t *testing.T) {
	mb, err := newTestMicroBatch()

	assert.NotNil(t, mb)
	assert.Nil(t, err)
}

func Test_New_applies_Limit(t *testing.T) {
	expectedLimit := 42
	mb, err := newTestMicroBatch(Limit[*JobInput, *JobOutput](expectedLimit))

	assert.NotNil(t, mb)
	assert.Nil(t, err)
	assert.Equal(t, expectedLimit, mb.limit)
}

func Test_New_applies_Cycle(t *testing.T) {
	expectedCycle := 42 * time.Second
	mb, err := newTestMicroBatch(Cycle[*JobInput, *JobOutput](expectedCycle))

	assert.NotNil(t, mb)
	assert.Nil(t, err)
	assert.Equal(t, expectedCycle, mb.cycle)
}

func Test_New_does_not_start_MicroBatch(t *testing.T) {
	mb, _ := newTestMicroBatch()

	assert.False(t, mb.isRunning.Load())
}

func Test_setIsRunning_updates_state(t *testing.T) {
	mb, _ := newTestMicroBatch()
	err := mb.setIsRunning(true)

	assert.Nil(t, err)
	assert.True(t, mb.isRunning.Load())
}

func Test_setIsRunning_erros_same_value(t *testing.T) {
	mb, _ := newTestMicroBatch()
	mb.setIsRunning(true)
	err := mb.setIsRunning(true)

	assert.NotNil(t, err)
	assert.Equal(t, ErrSameRunningState, err)
}

func Test_Start_runs_MicroBatch(t *testing.T) {
	c, cancel := context.WithCancel(context.TODO())
	defer cancel()

	mb, _ := newTestMicroBatch()
	go mb.Start(c)

	assert.Eventually(t, func() bool {
		return mb.isRunning.Load()
	}, 1*time.Second, 1*time.Millisecond)
}

func Test_Start_errors_on_second_call(t *testing.T) {
	c, cancel := context.WithCancel(context.TODO())
	defer cancel()

	mb, _ := newTestMicroBatch()
	go mb.Start(c)
	time.Sleep(1 * time.Millisecond) // give it a tick so that Start routine is riggered

	err := mb.Start(c)

	assert.NotNil(t, err)
	assert.Equal(t, ErrSameRunningState, err)
}

func Test_Submit_errors_on_not_running(t *testing.T) {
	mb, _ := newTestMicroBatch()
	job := func(t *JobInput) *JobResult[*JobOutput] {
		return nil
	}
	err := mb.Submit(job)

	assert.NotNil(t, err)
	assert.Equal(t, ErrIsNotRunning, err)
}

func Test_Submit_errors_on_nil_job(t *testing.T) {
	c, cancel := context.WithCancel(context.TODO())
	defer cancel()

	mb, _ := newTestMicroBatch()
	go mb.Start(c)
	time.Sleep(1 * time.Millisecond) // give it a tick so that Start routine is riggered

	err := mb.Submit(nil)

	assert.NotNil(t, err)
	assert.Equal(t, ErrNilJob, err)
}

func Test_Submit_submits_job(t *testing.T) {
	c, cancel := context.WithCancel(context.TODO())
	defer cancel()

	mb, _ := newTestMicroBatch()
	go mb.Start(c)
	time.Sleep(1 * time.Millisecond) // give it a tick so that Start routine is riggered

	job := func(t *JobInput) *JobResult[*JobOutput] {
		return &JobResult[*JobOutput]{}
	}
	err := mb.Submit(job)

	assert.Nil(t, err)
}

func Test_Submit_processes_job_eventually(t *testing.T) {
	c, cancel := context.WithCancel(context.TODO())
	defer cancel()

	mb, _ := newTestMicroBatch(Limit[*JobInput, *JobOutput](1))
	go mb.Start(c)
	time.Sleep(1 * time.Millisecond) // give it a tick so that Start routine is riggered

	var wg sync.WaitGroup
	var check atomic.Bool
	wg.Add(1)
	job := func(t *JobInput) *JobResult[*JobOutput] {
		check.Store(true)
		wg.Done()
		return &JobResult[*JobOutput]{}
	}

	assert.False(t, check.Load())
	mb.Submit(job)

	wg.Wait()
	assert.True(t, check.Load())
}

func newTestMicroBatch(options ...Option[*JobInput, *JobOutput]) (*MicroBatch[*JobInput, *JobOutput], error) {
	return New(newTestBatchProcessor(), options...)
}
