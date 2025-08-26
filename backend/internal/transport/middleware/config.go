package middleware

import (
	"url-shortener/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

// MiddlewareConfig holds all dependencies needed for middleware
type MiddlewareConfig struct {
	UserRepo usecase.UserRepository
}

// NewMiddlewareConfig creates a new middleware configuration
func NewMiddlewareConfig(userRepo usecase.UserRepository) *MiddlewareConfig {
	return &MiddlewareConfig{
		UserRepo: userRepo,
	}
}

// ApplyJWTAuth applies JWT authentication middleware
func (m *MiddlewareConfig) ApplyJWTAuth() gin.HandlerFunc {
	return JWTAuth(m.UserRepo)
}

// ApplyOptionalJWTAuth applies optional JWT authentication middleware
func (m *MiddlewareConfig) ApplyOptionalJWTAuth() gin.HandlerFunc {
	return OptionalJWTAuth()
}

// ApplyAPIKeyAuth applies API key authentication middleware
func (m *MiddlewareConfig) ApplyAPIKeyAuth() gin.HandlerFunc {
	return ApiKeyAuth(m.UserRepo)
}
