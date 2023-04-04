package xroles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRepository struct{}

func (m MockRepository) FindByRole(role string) ([]string, bool) {
	if role == "aGroup" {
		return []string{"g1Permission1", "g1Permission2", "g1Permission3"}, true
	}

	return nil, false
}

var mockRepository = &MockRepository{}

func TestExpandRolePermissions_roles_and_permissions(t *testing.T) {
	expander := NewRolePermissionsExpander(mockRepository)

	// GIVEN
	permissionsAndRoles := []string{"user_read", "aGroup", "token_create"}

	// WHEN expand
	permissions := expander(permissionsAndRoles)

	// THEN result has permissions and expanded roles
	expectedPermissions := []string{"user_read", "g1Permission1", "g1Permission2", "g1Permission3", "token_create"}

	assert.ElementsMatch(t, expectedPermissions, permissions)
}

func TestExpandRolePermissions_removes_duplicates(t *testing.T) {
	expander := NewRolePermissionsExpander(mockRepository)

	// GIVEN
	permissionsAndRoles := []string{"user_read", "aGroup", "token_create", "aGroup", "token_create", "g1Permission2"}

	// WHEN expand
	permissions := expander(permissionsAndRoles)

	// THEN result has permissions and expanded roles
	expectedPermissions := []string{"user_read", "g1Permission1", "g1Permission2", "g1Permission3", "token_create"}

	assert.ElementsMatch(t, expectedPermissions, permissions)
}
