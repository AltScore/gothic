package entity

import (
	"context"
	"github.com/AltScore/gothic/v2/pkg/xcontext"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// WHEN create new
	m := New()

	// THEN should have a new Id
	assert.NotEmpty(t, m.ID)
}

func TestNewAt(t *testing.T) {
	// GIVEN a time
	now := time.Now()

	// WHEN create new
	m := New(At(now))

	// THEN should have a new Id
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
			name: "with Id",
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
			m:    New(At(otherTime)),
			args: now,
			want: Metadata{
				CreatedAt: otherTime,
				UpdatedAt: now,
				Version:   1,
			},
		},
		{
			name: "increase version",
			m:    New(At(otherTime)).Clone(now).Clone(now),
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

func Test_NewIn(t *testing.T) {
	// WHEN create new
	ctxWithTenant := context.WithValue(context.Background(), xcontext.TenantCtxKey, "tenant")

	m := New(WithCtx(ctxWithTenant))

	// THEN should have a new Id
	assert.NotEmpty(t, m.ID)

	// AND should have the given tenant
	assert.Equal(t, "tenant", m.Tenant)
}

func Test_NewInAt(t *testing.T) {
	// GIVEN a time
	now := time.Now()

	// WHEN create new
	ctxWithTenant := context.WithValue(context.Background(), xcontext.TenantCtxKey, "tenant")

	m := New(WithCtx(ctxWithTenant), At(now))

	// THEN should have a new Id
	assert.NotEmpty(t, m.ID)

	// AND should have the given CreatedAt
	assert.Equal(t, now, m.CreatedAt)

	// AND should have the given UpdatedAt
	assert.Equal(t, now, m.UpdatedAt)

	// AND should have the given tenant
	assert.Equal(t, "tenant", m.Tenant)
}
