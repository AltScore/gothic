package xstring

import (
	"fmt"
	"testing"
)

func TestRandString(t *testing.T) {
	for size := 1; size < 100; size++ {
		t.Run(fmt.Sprintf("Size: %d", size), func(t *testing.T) {
			if got := RandString(size); len(got) != size {
				t.Errorf("RandString() = %v, want %v", got, size)
			} else {
				t.Logf("RandString() = %v", got)
			}
		})
	}
}
