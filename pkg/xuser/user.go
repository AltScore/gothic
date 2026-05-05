package xuser

import (
	"github.com/google/uuid"
)

// User represents the current user performing the request.
type User interface {
	Id() uuid.UUID
	Name() string
	Tenant() string
	HasPermission(mustAll bool, permission ...string) error
	Permissions() []string
}

// ImpersonatedUser is an optional interface that can be implemented by User when it can provide the id of the user that is being impersonated.
// This is used to determine if the user is being impersonated or not. A User implementation can implement this interface
// but a specific user can be impersonated or not.
type ImpersonatedUser interface {
	RealUserId() uuid.UUID
}
