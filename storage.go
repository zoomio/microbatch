package microbatch

// Used to store submitted jobs.
// Implement this interface and pass it as an option to the MicroBatch#New.
type BatchStorage[T, R any] interface {

	// Returns current size of the stored batch.
	Size() int

	// Appends a job to the storage.
	Append(job Job[T, R]) error

	// Returns all the stored jobs.
	GetAll() ([]Job[T, R], error)

	// Clears all the jobs in the storage.
	Clear() error
}

// Ensure InMemoryStorage conforms to the BatchStorage interface.
var _ BatchStorage[any, any] = &InMemoryStorage[any, any]{}

// Simple in memory storage.
type InMemoryStorage[T, R any] struct {
	limit int
	batch []Job[T, R]
}

func NewInMemoryStorage[T, R any](limit int) *InMemoryStorage[T, R] {
	return &InMemoryStorage[T, R]{limit: limit, batch: make([]Job[T, R], 0, limit)}
}

func (st *InMemoryStorage[T, R]) Size() int {
	return len(st.batch)
}

func (st *InMemoryStorage[T, R]) Append(job Job[T, R]) error {
	if len(st.batch) == st.limit {
		return ErrStorageFull
	}
	st.batch = append(st.batch, job)
	return nil
}

func (st *InMemoryStorage[T, R]) GetAll() ([]Job[T, R], error) {
	// copy processed batch and clear the internal to avoid side effects
	b := make([]Job[T, R], len(st.batch))
	copy(b, st.batch)
	return b, nil
}

func (st *InMemoryStorage[T, R]) Clear() error {
	st.batch = make([]Job[T, R], 0, st.limit)
	return nil
}
