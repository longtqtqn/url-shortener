package middleware

import (
	"net/http"
	"url-shortener/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func ApiKeyAuth(userRepo usecase.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-KEY")
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			return
		}
		apiKeyID, err := userRepo.GetAPIKeyIDByAPIKey(ctx.Request.Context(), apiKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			return
		}
		ctx.Set("apiKeyID", apiKeyID)
		ctx.Next()
	}
}
