package apiusers

type ID = string

// ApiUser is a user of the API
type ApiUser interface {
	Id() ID
	AccountId() string
	HasKey(apiKey AccountApiKey) bool
	HasPermission(permission string) bool
	Permissions() []string
}
