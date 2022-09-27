package apiusers

import "context"

const ApiUserKey = "api-user"

var NoUser = &noUser{}

func GetApiUser(ctx context.Context) ApiUser {
	value := ctx.Value(ApiUserKey)
	if value == nil {
		return NoUser
	}
	return value.(ApiUser)
}

func SetApiUser(ctx context.Context, user ApiUser) context.Context {
	return context.WithValue(ctx, ApiUserKey, user)
}

type noUser struct{}

func (n noUser) Id() ID {
	return "<no-user>"
}

func (n noUser) AccountId() string {
	return ""
}

func (n noUser) HasKey(apiKey AccountApiKey) bool {
	return false
}

func (n noUser) HasPermission(permission string) bool {
	return false
}

func (n noUser) Permissions() []string {
	return []string{}
}
