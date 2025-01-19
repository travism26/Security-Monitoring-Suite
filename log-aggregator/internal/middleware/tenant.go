package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/travism26/log-aggregator/internal/domain"
	"github.com/travism26/log-aggregator/internal/repository/postgres"
	"github.com/travism26/log-aggregator/internal/service"
)

const (
	// ContextKeyTenant is the key used to store tenant context in gin.Context
	ContextKeyTenant = "tenant"
)

// TenantContext holds tenant-specific information
type TenantContext struct {
	OrganizationID string
	APIKeyType     string
}

// TenantMiddleware creates middleware for tenant validation
func TenantMiddleware(validator domain.APIKeyValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get API key from header or query parameter
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API key is required",
			})
			c.Abort()
			return
		}

		// Hash the API key
		keyHash := validator.HashKey(apiKey)

		// Validate the API key
		apiKeyInfo, err := validator.ValidateKey(keyHash)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		// Create tenant context
		tenantCtx := &TenantContext{
			OrganizationID: apiKeyInfo.OrganizationID.String(),
			APIKeyType:     string(apiKeyInfo.KeyType),
		}

		// Store tenant context in gin.Context
		c.Set(ContextKeyTenant, tenantCtx)

		// Add tenant headers for debugging/logging
		c.Header("X-Organization-ID", tenantCtx.OrganizationID)
		c.Header("X-API-Key-Type", tenantCtx.APIKeyType)

		c.Next()
	}
}

// GetTenantContext retrieves tenant context from gin.Context
func GetTenantContext(c *gin.Context) *TenantContext {
	if tenant, exists := c.Get(ContextKeyTenant); exists {
		if tc, ok := tenant.(*TenantContext); ok {
			return tc
		}
	}
	return nil
}

// RequireAgentKey middleware ensures the API key is an agent key
func RequireAgentKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant := GetTenantContext(c)
		if tenant == nil || tenant.APIKeyType != string(domain.APIKeyTypeAgent) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Agent API key required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireCustomerKey middleware ensures the API key is a customer key
func RequireCustomerKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant := GetTenantContext(c)
		if tenant == nil || tenant.APIKeyType != string(domain.APIKeyTypeCustomer) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Customer API key required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Tenant returns the tenant middleware with default validator
func Tenant() gin.HandlerFunc {
	// Create a new API key service with default settings
	validator := service.NewAPIKeyService(
		postgres.NewAPIKeyRepository(postgres.GetDB()),
	)
	return TenantMiddleware(validator)
}
