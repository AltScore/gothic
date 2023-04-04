package xroles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Expand_roles_without_dependencies(t *testing.T) {
	rolesRepository := NewInMemoryRolesRepository(map[string][]string{
		"role1": {"permission1", "permission2"},
		"role2": {"permission3", "permission4"},
	})

	permissions1, _ := rolesRepository.FindByRole("role1")
	permissions2, _ := rolesRepository.FindByRole("role2")

	assert.ElementsMatchf(t, permissions1, []string{"permission1", "permission2"}, "permissions should be expanded")
	assert.ElementsMatchf(t, permissions2, []string{"permission3", "permission4"}, "permissions should be expanded")
}

func Test_Expand_roles_with_spaces(t *testing.T) {
	rolesRepository := NewInMemoryRolesRepository(map[string][]string{
		"role1": {"permission1", "    permission2"},
		"role2": {"permission3   ", "permission4"},
	})

	permissions1, _ := rolesRepository.FindByRole("role1")
	permissions2, _ := rolesRepository.FindByRole("role2")

	assert.ElementsMatchf(t, permissions1, []string{"permission1", "permission2"}, "permissions should be expanded")
	assert.ElementsMatchf(t, permissions2, []string{"permission3", "permission4"}, "permissions should be expanded")
}

func Test_Expand_roles_with_dependencies(t *testing.T) {
	rolesRepository := NewInMemoryRolesRepository(map[string][]string{
		"role1": {"permission1", "permission2", "role2"},
		"role2": {"permission3", "permission4"},
	})

	permissions1, _ := rolesRepository.FindByRole("role1")
	permissions2, _ := rolesRepository.FindByRole("role2")

	assert.ElementsMatchf(t, permissions1, []string{"permission1", "permission2", "permission3", "permission4"}, "permissions should be expanded")
	assert.ElementsMatchf(t, permissions2, []string{"permission3", "permission4"}, "permissions should be expanded")
}

func Test_Expand_roles_with_dependencies_without_duplicates(t *testing.T) {
	rolesRepository := NewInMemoryRolesRepository(map[string][]string{
		"role1": {"permission1", "permission1", "permission2", "role2", "permission4"},
		"role2": {"permission3", "permission4", "permission1"},
	})

	permissions1, _ := rolesRepository.FindByRole("role1")
	permissions2, _ := rolesRepository.FindByRole("role2")

	assert.ElementsMatchf(t, permissions1, []string{"permission1", "permission2", "permission3", "permission4"}, "permissions should be expanded")
	assert.ElementsMatchf(t, permissions2, []string{"permission1", "permission3", "permission4"}, "permissions should be expanded")
}

func Test_Expand_roles_with_circular_dependencies(t *testing.T) {
	rolesRepository := NewInMemoryRolesRepository(map[string][]string{
		"role1": {"permission1", "permission2", "role2"},
		"role2": {"permission3", "permission4", "role1"},
	})

	permissions1, _ := rolesRepository.FindByRole("role1")
	permissions2, _ := rolesRepository.FindByRole("role2")

	assert.ElementsMatchf(t, permissions1, []string{"permission1", "permission2", "permission3", "permission4"}, "permissions should be expanded")
	assert.ElementsMatchf(t, permissions2, []string{"permission3", "permission4"}, "permissions should be expanded")
}
