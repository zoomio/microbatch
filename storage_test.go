package microbatch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewInMemoryStorage_creates_storage(t *testing.T) {
	st := NewInMemoryStorage[any, any](1)

	assert.NotNil(t, st)
}

func Test_Append_appends_job(t *testing.T) {
	st := NewInMemoryStorage[any, any](1)
	job := func(t any) *JobResult[any] {
		return &JobResult[any]{}
	}
	err := st.Append(job)
	assert.Nil(t, err)
}

func Test_Append_errors_on_limit(t *testing.T) {
	st := NewInMemoryStorage[any, any](1)
	job1 := func(t any) *JobResult[any] {
		return &JobResult[any]{}
	}
	job2 := func(t any) *JobResult[any] {
		return &JobResult[any]{}
	}
	err := st.Append(job1)
	assert.Nil(t, err)

	err = st.Append(job2)
	assert.NotNil(t, err)
	assert.Equal(t, ErrStorageFull, err)
}

func Test_Size_returns_current_size(t *testing.T) {
	st := NewInMemoryStorage[any, any](5)
	job := func(t any) *JobResult[any] {
		return &JobResult[any]{}
	}
	st.Append(job)

	assert.Equal(t, 1, st.Size())
}

func Test_GetAll_returns_all_jobs(t *testing.T) {
	st := NewInMemoryStorage[any, any](5)

	st.Append(func(t any) *JobResult[any] {
		return &JobResult[any]{}
	})
	st.Append(func(t any) *JobResult[any] {
		return &JobResult[any]{}
	})
	st.Append(func(t any) *JobResult[any] {
		return &JobResult[any]{}
	})

	assert.Equal(t, 3, st.Size())
}

func Test_Clear_clears_storage(t *testing.T) {
	st := NewInMemoryStorage[any, any](5)

	st.Append(func(t any) *JobResult[any] {
		return &JobResult[any]{}
	})
	st.Append(func(t any) *JobResult[any] {
		return &JobResult[any]{}
	})
	st.Append(func(t any) *JobResult[any] {
		return &JobResult[any]{}
	})

	st.Clear()

	assert.Equal(t, 0, st.Size())
}
