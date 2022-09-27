package basicuser

import (
	"fmt"

	"github.com/AltScore/gothic/pkg/apiusers"
	"github.com/totemcaf/gollections/slices"
)

type basicUser struct {
	id          apiusers.ID
	accountId   string
	hash        string
	permissions []string
}

func New(id apiusers.ID, accountId string, hash string, permissions []string) apiusers.ApiUser {
	return &basicUser{
		id:          id,
		accountId:   accountId,
		hash:        hash,
		permissions: permissions,
	}
}

func (b basicUser) Id() apiusers.ID {
	return b.id
}

func (b basicUser) AccountId() string {
	return b.accountId
}

func (b basicUser) HasKey(apiKey apiusers.AccountApiKey) bool {
	return b.hash == apiKey.GetHash()
}

func (b basicUser) HasPermission(permission string) bool {
	return slices.Has2(b.permissions, permission)
}

func (b basicUser) Permissions() []string {
	return b.permissions
}

func (b basicUser) String() string {
	return fmt.Sprintf("%s/%s", b.id, b.accountId)
}
