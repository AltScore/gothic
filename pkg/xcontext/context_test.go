package xcontext

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestContextHasUser(t *testing.T) {
	user := &mockUser{}
	ctx := WithUser(context.Background(), user)

	// WHEN calls User
	actual, err := GetUser(ctx)

	// THEN
	require.Nil(t, err)
	require.Equal(t, user, actual)
}

func TestContextNoUser(t *testing.T) {
	// GIVEN
	ctx := context.Background()

	// WHEN calls User
	actual, err := GetUser(ctx)

	// THEN
	require.Error(t, err)
	require.Nil(t, actual)
}

func TestContextWrongType(t *testing.T) {
	// GIVEN
	ctx := context.WithValue(context.Background(), UserCtxKey, "not a user")

	// WHEN calls User
	actual, err := GetUser(ctx)

	// THEN
	require.Error(t, err)
	require.Nil(t, actual)
}

type mockUser struct {
}

func (m *mockUser) Permissions() []string {
	panic("should not be called")
}

func (m *mockUser) Id() uuid.UUID {
	panic("should not be called")
}

func (m *mockUser) Name() string {
	panic("should not be called")
}

func (m *mockUser) Tenant() string {
	panic("should not be called")
}

func (m *mockUser) HasPermission(_ bool, _ ...string) error {
	panic("should not be called")
}
