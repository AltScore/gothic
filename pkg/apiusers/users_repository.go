package apiusers

type AccountApiKey interface {
	GetHash() string
}

type Repository interface {
	FindByHash(apiKey AccountApiKey) (ApiUser, error)
	FindById(accountId ID) (ApiUser, error)
}

type AccountApiKeyGenerator interface {
	NewAccountApiKey(value string) (AccountApiKey, error)
}
