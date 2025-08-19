package middleware

import (
	"net/http"
	"url-shortener/internal/domain"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := ctx.MustGet("currentUser").(*domain.User)
		if !ok || user == nil || user.Role != "admin" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin role required"})
			return
		}
		ctx.Next()
	}
}
