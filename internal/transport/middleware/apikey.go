package middleware

import (
	"net/http"
	"url-shortener/internal/usecase"

	"github.com/gin-gonic/gin"
)

func ApiKeyAuth(userRepo usecase.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-KEY")
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			return
		}
		user, err := userRepo.GetByAPIKey(ctx.Request.Context(), apiKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
		}
		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
