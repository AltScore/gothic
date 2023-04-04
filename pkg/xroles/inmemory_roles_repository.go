package xroles

import (
	"sort"
	"strings"

	"github.com/totemcaf/gollections/maps"
	"github.com/totemcaf/gollections/slices"
)

type role struct {
	name        string
	rawRoles    []string
	permissions []string
	expanded    bool
}

type InMemoryRolesRepository struct {
	roles map[string]*role
}

// NewInMemoryRolesRepository returns a new InMemoryRolesRepository whose roles are
// defined by the given map. The map keys are the role names and the values are
// the permissions for that role.
func NewInMemoryRolesRepository(rolesPermissions map[string][]string) *InMemoryRolesRepository {

	roles := make(map[string]*role)

	for roleName, permissions := range rolesPermissions {
		roles[roleName] = &role{
			name:        roleName,
			rawRoles:    slices.Map(permissions, strings.TrimSpace),
			permissions: nil,
			expanded:    false,
		}
	}

	roleNames := maps.Keys(roles)
	sort.Strings(roleNames) // To make tests deterministic

	for _, roleName := range roleNames {
		roles[roleName].expand(roles, make(map[string]bool, 0))
	}

	return &InMemoryRolesRepository{roles: roles}
}

func (r *role) expand(roleMap map[string]*role, alreadyVisited map[string]bool) []string {
	if r.expanded {
		return r.permissions
	}

	r.expanded = true
	alreadyVisited[r.name] = true

	permissions := make([]string, 0)

	sortedRoles := slices.Clone(r.rawRoles)
	sort.Strings(sortedRoles) // To make tests deterministic

	for _, roleName := range sortedRoles {
		if role, found := roleMap[roleName]; !found {
			permissions = append(permissions, roleName)
		} else if _, found := alreadyVisited[roleName]; found {
			// omit circular dependencies
		} else {
			permissions = append(permissions, role.expand(roleMap, alreadyVisited)...)
		}
	}

	r.permissions = removeDuplicates(permissions)
	return permissions
}

func removeDuplicates(permissions []string) []string {
	uniquePermissions := make(map[string]bool, len(permissions))

	for _, permission := range permissions {
		uniquePermissions[permission] = true
	}

	return maps.Keys(uniquePermissions)
}

func (i InMemoryRolesRepository) FindByRole(role string) ([]string, bool) {
	if role, found := i.roles[role]; found {
		return role.permissions, true
	}

	return nil, false
}
