package router

import (
	"context"
	"net/http"
	"url-shortener/internal/transport/http/handler"
	"url-shortener/internal/transport/middleware"
	"url-shortener/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func Register(r *gin.Engine, db *bun.DB, userRepo usecase.UserRepository, linkH *handler.LinkHttpHandler, adminH *handler.AdminHttpHandler) {
	// health
	r.HEAD("/healthz", func(c *gin.Context) {
		if err := db.RunInTx(c, nil, func(ctx context.Context, tx bun.Tx) error { return nil }); err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		c.Status(http.StatusOK)
	})

	// public
	public := r.Group("/")
	linkH.RegisterPublicRoutes(public)

	// api (auth required)
	api := r.Group("/api")
	api.Use(middleware.ApiKeyAuth(userRepo))
	linkH.RegisterAuthRoutes(api)

	// admin
	admin := r.Group("/admin")
	admin.Use(middleware.ApiKeyAuth(userRepo), middleware.RequireAdmin())
	adminH.RegisterAdminRoutes(admin)
}
