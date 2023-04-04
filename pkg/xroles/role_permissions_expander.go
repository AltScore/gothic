package xroles

import (
	"github.com/totemcaf/gollections/slices"
)

// RolesRepository provides roles and their corresponding permissions
type RolesRepository interface {
	// FindByRole returns the permissions for the given role and true if the role exists.
	// If the role does not exist, it returns an empty slice and false.
	FindByRole(role string) ([]string, bool)
}

// ExpandRolPermissions returns the given permissions and roles expanded to its corresponding permissions.
// allPermissions by role are provided by RolesRepository
// Roles are expanded recursively.
type ExpandRolPermissions = func(permissionsOrRoles []string) []string

type rolePermissionsExpander struct {
	RolesRepository
}

// NewRolePermissionsExpander returns a function that returns the given permissions and
// roles expanded to its corresponding permissions.
// allPermissions by role are provided by RolesRepository
func NewRolePermissionsExpander(rolesRepository RolesRepository) ExpandRolPermissions {
	expander := &rolePermissionsExpander{RolesRepository: rolesRepository}
	return expander.expandRoles
}

func (f *rolePermissionsExpander) expandRoles(permissions []string) []string {
	return slices.JoinDistinct(slices.Map(permissions, f.expandRole)...)
}

func (f *rolePermissionsExpander) expandRole(permissionOrRole string) []string {
	if permissions, found := f.RolesRepository.FindByRole(permissionOrRole); found {
		return permissions
	}

	return []string{permissionOrRole}
}
