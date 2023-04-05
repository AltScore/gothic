package xsync

import (
	"context"
	"golang.org/x/sync/semaphore"
	"runtime"
	"time"
)

type ParallelMapper[Input, Output any] interface {
	// Map applies the given mapper to each element of the input slice in parallel.
	// If the context expires before all the mappers are finished, the remaining mappers are cancelled.
	// If a mapper fails, the error is returned in the errs slice.
	// If a mapper succeeds, the result is returned in the results slice.
	// FUnction returns an error if thera is problem not related to the mapper function.
	Map([]Input) error

	// Results returns the results of the last call to Map.
	Results() []Output

	// Errors returns the errors of the last call to Map.
	Errors() []error
}

type ParallelOption func(*parallelMapperOptions)

type parallelMapperOptions struct {
	maxWait    time.Duration
	maxWorkers int
}

type parallelMapper[Input, Output any] struct {
	parallelMapperOptions
	worker func(Input) (Output, error)

	results []Output

	errs []error
}

// NewParallelMapper returns a mapper hat applies the given mapper to each element of the input slice in parallel
// with at most maxWorkers in parallel.
// The results and errors are returned in the same order as the slice.
// If a maxWait is configured and the context expires before all the mappers are finished,
// the remaining mappers are cancelled.
//
// If a mapper fails, the error is returned in the errs slice.
// If a mapper succeeds, the result is returned in the results slice.
//
// Example:
//
//	   values := []int{1, 2, 3, 4, 5}
//	   mapper := NewParallelMapper(func(ctx context.Context, v V) (T, error) {
//		    return v * v, nil
//	   }, WithMaxWait( 20 * time.Second), WithMaxWorkers(2))
//
//	  errs := mapper.Map(values)
//	  results := mapper.Results()
//	  errs := mapper.Errors()
func NewParallelMapper[Input, Output any](worker func(Input) (Output, error), options ...ParallelOption) ParallelMapper[Input, Output] {
	mapper := &parallelMapper[Input, Output]{
		parallelMapperOptions: parallelMapperOptions{
			maxWait:    0,
			maxWorkers: runtime.GOMAXPROCS(0),
		},

		worker: worker,
	}

	for _, option := range options {
		option(&mapper.parallelMapperOptions)
	}

	return mapper
}

func (p *parallelMapper[Input, Output]) Results() []Output {
	return p.results
}

func (p *parallelMapper[Input, Output]) Errors() []error {
	return p.errs
}

func (p *parallelMapper[Input, Output]) Map(inputs []Input) error {
	var ctx context.Context
	if p.maxWait > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(p.maxWait))
		defer cancel()
	} else {
		ctx = context.Background()
	}

	var (
		sem = semaphore.NewWeighted(int64(p.maxWorkers))
		err error
	)
	p.results = make([]Output, len(inputs))
	p.errs = make([]error, len(inputs))

	// Compute the output using up to maxWorkers goroutines at a time.
	for i := range inputs {
		// When maxWorkers goroutines are in flight, Acquire blocks until one of the
		// workers finishes.
		if err := sem.Acquire(ctx, 1); err != nil {
			p.errs[i] = err
			break
		}

		go func(i int) {
			defer sem.Release(1) // finally release the token
			p.results[i], p.errs[i] = p.worker(inputs[i])
		}(i)
	}

	// Acquire all the tokens to wait for any remaining workers to finish.
	err = sem.Acquire(ctx, int64(p.maxWorkers))

	return err
}

// WithMaxWait configures the maximum time to wait for all the mappers to finish.
func WithMaxWait(maxWait time.Duration) ParallelOption {
	return func(p *parallelMapperOptions) {
		p.maxWait = maxWait
	}
}

// WithMaxWorkers configures the maximum number of workers to use.
func WithMaxWorkers(maxWorkers int) ParallelOption {
	return func(p *parallelMapperOptions) {
		p.maxWorkers = maxWorkers
	}
}
