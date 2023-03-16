package xsync

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func TestParallelMapper_Map_without_options(t *testing.T) {
	// GIVEN a ParallelMapper
	mapper := NewParallelMapper[string, int](func(key string) (int, error) {
		return strconv.Atoi(key)
	})

	// WHEN Map is called
	err := mapper.Map([]string{"1", "2", "3", "4", "5"})

	// THEN no error is returned
	require.NoError(t, err)

	// AND the results are returned in the same order as the input
	require.Equal(t, []int{1, 2, 3, 4, 5}, mapper.Results())

	// AND no errors are returned
	require.Equal(t, []error{nil, nil, nil, nil, nil}, mapper.Errors())
}

func TestParallelMapper_Map_return_errors(t *testing.T) {
	// GIVEN a ParallelMapper
	mapper := NewParallelMapper[string, int](func(key string) (int, error) {
		return strconv.Atoi(key)
	})

	runTest(t, mapper)
}

func TestParallelMapper_Map_with_max_workers(t *testing.T) {
	// GIVEN a ParallelMapper
	mapper := NewParallelMapper[string, int](func(key string) (int, error) {
		return strconv.Atoi(key)
	}, WithMaxWorkers(2))

	runTest(t, mapper)
}

func TestParallelMapper_Map_with_max_time(t *testing.T) {
	// GIVEN a ParallelMapper
	mapper := NewParallelMapper[string, int](func(key string) (int, error) {
		return strconv.Atoi(key)
	}, WithMaxWait(2*time.Second))

	runTest(t, mapper)
}

func TestParallelMapper_Map_with_max_workers_and_max_time_cannot_finish(t *testing.T) {
	// GIVEN a ParallelMapper
	mapper := NewParallelMapper[string, int](func(key string) (int, error) {
		time.Sleep(200 * time.Millisecond)
		return strconv.Atoi(key)
	}, WithMaxWait(500*time.Millisecond), WithMaxWorkers(2))

	// WHEN Map is called
	err := mapper.Map([]string{"1", "2", "tres", "4", "5", "6", "7", "8", "9", "10"})

	// THEN an error is returned
	require.Error(t, err)

	// AND no all results are returned
	require.Equal(t, []int{1, 2, 0, 4, 0, 0, 0, 0, 0, 0}, mapper.Results())

}
func runTest(t *testing.T, mapper ParallelMapper[string, int]) {
	// WHEN Map is called
	err := mapper.Map([]string{"1", "2", "tres", "4", "5"})

	// THEN no error is returned
	require.NoError(t, err)

	// AND the results are returned in the same order as the input
	require.Equal(t, []int{1, 2, 0, 4, 5}, mapper.Results())

	// AND errors is returned for the invalid input
	for i, err := range mapper.Errors() {
		if i == 2 {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}
