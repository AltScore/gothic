package xcontext

import (
	"context"
	"github.com/AltScore/gothic/v2/pkg/ids"
	"github.com/stretchr/testify/require"
	"testing"
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

func (m *mockUser) Id() ids.Id {
	panic("should not be called")
}

func (m *mockUser) Name() string {
	panic("should not be called")
}

func (m *mockUser) Tenant() string {
	panic("should not be called")
}

func (m *mockUser) HasPermission(_ string) bool {
	panic("should not be called")
}
