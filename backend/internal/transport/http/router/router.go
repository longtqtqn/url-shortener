package router

import (
	"context"
	"net/http"
	"strings"
	httptypes "url-shortener/backend/internal/transport/http"
	"url-shortener/backend/internal/transport/http/handler"
	"url-shortener/backend/internal/transport/middleware"
	"url-shortener/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func Register(r *gin.Engine, db *bun.DB, userRepo usecase.UserRepository, linkH *handler.LinkHttpHandler, userH *handler.UserHttpHandler) {
	// health check
	r.HEAD("/healthz", func(c *gin.Context) {
		if err := db.RunInTx(c, nil, func(ctx context.Context, tx bun.Tx) error { return nil }); err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		c.Status(http.StatusOK)
	})

	// Create middleware config
	middlewareConfig := middleware.NewMiddlewareConfig(userRepo)

	// Register all routes from handlers
	handlers := []httptypes.Handler{linkH, userH}

	for _, handler := range handlers {
		routes := handler.GetRoutes()
		for _, route := range routes {
			registerRoute(r, route, middlewareConfig)
		}
	}
}

// registerRoute registers a single route with appropriate middleware
func registerRoute(r *gin.Engine, route httptypes.Route, config *middleware.MiddlewareConfig) {
	var middlewares []gin.HandlerFunc

	// Apply authentication middleware based on route configuration
	if route.RequireAuth {
		switch route.AuthType {
		case "jwt":
			middlewares = append(middlewares, config.ApplyJWTAuth())
		case "apikey":
			middlewares = append(middlewares, config.ApplyAPIKeyAuth())
		}
	}

	// Combine middleware with handler
	handlers := append(middlewares, route.Handler)

	// Register route based on HTTP method
	switch strings.ToUpper(route.Method) {
	case "GET":
		r.GET(route.Path, handlers...)
	case "POST":
		r.POST(route.Path, handlers...)
	case "PUT":
		r.PUT(route.Path, handlers...)
	case "DELETE":
		r.DELETE(route.Path, handlers...)
	case "PATCH":
		r.PATCH(route.Path, handlers...)
	case "HEAD":
		r.HEAD(route.Path, handlers...)
	case "OPTIONS":
		r.OPTIONS(route.Path, handlers...)
	}
}
