package xuser

import "github.com/AltScore/gothic/pkg/ids"

type User interface {
	ID() ids.ID
	Name() string
	TenantID() ids.ID
	HasPermission(permission string) bool
}
