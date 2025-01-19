package domain

// APIKeyValidator interface for validating API keys
type APIKeyValidator interface {
	ValidateKey(keyHash string) (*APIKey, error)
	HashKey(key string) string
}
