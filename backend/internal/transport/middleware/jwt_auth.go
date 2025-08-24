package middleware

import (
	"net/http"
	"url-shortener/backend/internal/auth"
	"url-shortener/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

// JWTAuth middleware validates JWT tokens and sets user information in context
func JWTAuth(userRepo usecase.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := auth.ExtractTokenFromHeader(authHeader)
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		userID, err := auth.ValidateJWT(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		ctx.Set("userID", userID)

		apiKey := ctx.GetHeader("X-API-KEY")
		if apiKey != "" {
			apiKeyID, err := userRepo.GetAPIKeyIDByAPIKey(ctx.Request.Context(), apiKey)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
				return
			}
			ctx.Set("apiKeyID", apiKeyID)
		}
		ctx.Next()
	}
}

// OptionalJWTAuth middleware validates JWT tokens if present, but doesn't require them
func OptionalJWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader != "" {
			tokenString := auth.ExtractTokenFromHeader(authHeader)
			if tokenString != "" {
				userID, err := auth.ValidateJWT(tokenString)
				if err == nil {
					ctx.Set("userID", userID)
				}
			}
		}
		ctx.Next()
	}
}
