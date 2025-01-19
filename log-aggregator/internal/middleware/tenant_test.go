package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/travism26/log-aggregator/internal/domain"
)

// Ensure MockAPIKeyService implements domain.APIKeyValidator
var _ domain.APIKeyValidator = (*MockAPIKeyService)(nil)

type MockAPIKeyService struct {
	mock.Mock
}

func (m *MockAPIKeyService) ValidateKey(keyHash string) (*domain.APIKey, error) {
	args := m.Called(keyHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.APIKey), args.Error(1)
}

func (m *MockAPIKeyService) HashKey(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func TestTenantMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		apiKey         string
		setupMock      func(*MockAPIKeyService)
		expectedStatus int
		checkContext   func(*testing.T, *gin.Context)
	}{
		{
			name:   "Valid API Key",
			apiKey: "test-key",
			setupMock: func(m *MockAPIKeyService) {
				orgID := uuid.New()
				m.On("HashKey", "test-key").Return("hashed-key")
				m.On("ValidateKey", "hashed-key").Return(&domain.APIKey{
					ID:             uuid.New(),
					OrganizationID: orgID,
					KeyType:        domain.APIKeyTypeCustomer,
					Status:         domain.APIKeyStatusActive,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				tenant := GetTenantContext(c)
				assert.NotNil(t, tenant)
				assert.NotEmpty(t, tenant.OrganizationID)
				assert.Equal(t, string(domain.APIKeyTypeCustomer), tenant.APIKeyType)
			},
		},
		{
			name:   "Missing API Key",
			apiKey: "",
			setupMock: func(m *MockAPIKeyService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusUnauthorized,
			checkContext: func(t *testing.T, c *gin.Context) {
				tenant := GetTenantContext(c)
				assert.Nil(t, tenant)
			},
		},
		{
			name:   "Invalid API Key",
			apiKey: "invalid-key",
			setupMock: func(m *MockAPIKeyService) {
				m.On("HashKey", "invalid-key").Return("hashed-invalid-key")
				m.On("ValidateKey", "hashed-invalid-key").Return(nil, fmt.Errorf("invalid API key"))
			},
			expectedStatus: http.StatusUnauthorized,
			checkContext: func(t *testing.T, c *gin.Context) {
				tenant := GetTenantContext(c)
				assert.Nil(t, tenant)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAPIKeyService)
			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			if tt.apiKey != "" {
				c.Request.Header.Set("X-API-Key", tt.apiKey)
			}

			var validator domain.APIKeyValidator = mockService
			middleware := TenantMiddleware(validator)
			middleware(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkContext(t, c)
			mockService.AssertExpectations(t)
		})
	}
}

func TestRequireAgentKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectedStatus int
	}{
		{
			name: "Valid Agent Key",
			setupContext: func(c *gin.Context) {
				c.Set(ContextKeyTenant, &TenantContext{
					OrganizationID: uuid.New().String(),
					APIKeyType:     string(domain.APIKeyTypeAgent),
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Customer Key",
			setupContext: func(c *gin.Context) {
				c.Set(ContextKeyTenant, &TenantContext{
					OrganizationID: uuid.New().String(),
					APIKeyType:     string(domain.APIKeyTypeCustomer),
				})
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Missing Tenant Context",
			setupContext: func(c *gin.Context) {
				// No context setup
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			tt.setupContext(c)

			middleware := RequireAgentKey()
			middleware(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRequireCustomerKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectedStatus int
	}{
		{
			name: "Valid Customer Key",
			setupContext: func(c *gin.Context) {
				c.Set(ContextKeyTenant, &TenantContext{
					OrganizationID: uuid.New().String(),
					APIKeyType:     string(domain.APIKeyTypeCustomer),
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Agent Key",
			setupContext: func(c *gin.Context) {
				c.Set(ContextKeyTenant, &TenantContext{
					OrganizationID: uuid.New().String(),
					APIKeyType:     string(domain.APIKeyTypeAgent),
				})
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Missing Tenant Context",
			setupContext: func(c *gin.Context) {
				// No context setup
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			tt.setupContext(c)

			middleware := RequireCustomerKey()
			middleware(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
