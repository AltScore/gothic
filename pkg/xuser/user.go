package xuser

import (
	"github.com/AltScore/gothic/v2/pkg/ids"
)

// User represent the current user performing the request.
type User interface {
	Id() ids.Id
	Name() string
	Tenant() string
	HasPermission(permission string) bool
}

// ImpersonatedUser is an optional interface that can be implemented by User when it can provide the id of the user that is being impersonated.
// This is used to determine if the user is being impersonated or not. A User implementation can implement this interface
// but a specific user can be impersonated or not.
type ImpersonatedUser interface {
	RealUserId() ids.Id
}
