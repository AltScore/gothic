package entity

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// WHEN create new
	m := New()

	// THEN should have a new ID
	assert.NotEmpty(t, m.ID)
}

func TestNewAt(t *testing.T) {
	// GIVEN a time
	now := time.Now()

	// WHEN create new
	m := NewAt(now)

	// THEN should have a new ID
	assert.NotEmpty(t, m.ID)

	// AND should have the given CreatedAt
	assert.Equal(t, now, m.CreatedAt)

	// AND should have the given UpdatedAt
	assert.Equal(t, now, m.UpdatedAt)
}

func TestMetadata_Clone(t *testing.T) {
	now := time.Now()
	otherTime := now.Add(-time.Hour)

	tests := []struct {
		name      string
		m         Metadata
		args      time.Time
		wantNewId bool
		want      Metadata
	}{
		{
			name:      "empty",
			m:         Metadata{},
			args:      now,
			wantNewId: true,
			want: Metadata{
				CreatedAt: now,
				UpdatedAt: now,
				Version:   1,
			},
		},
		{
			name: "with ID",
			m:    New(),
			args: now,
			want: Metadata{
				CreatedAt: now,
				UpdatedAt: now,
				Version:   1,
			},
		},
		{
			name: "with CreatedAt",
			m:    NewAt(otherTime),
			args: now,
			want: Metadata{
				CreatedAt: otherTime,
				UpdatedAt: now,
				Version:   1,
			},
		},
		{
			name: "increase version",
			m:    NewAt(otherTime).Clone(now).Clone(now),
			args: now,
			want: Metadata{
				CreatedAt: otherTime,
				UpdatedAt: now,
				Version:   3, // 2 from the previous clones, 1 from this clone
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clone := tt.m.Clone(tt.args)
			if tt.wantNewId {
				assert.NotEmpty(t, clone.ID)
				tt.want.ID = clone.ID // It is a new one
			} else {
				tt.want.ID = tt.m.ID // should be kept the same
			}
			if got := clone; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}
