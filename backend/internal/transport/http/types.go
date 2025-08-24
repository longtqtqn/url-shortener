package http

import "github.com/gin-gonic/gin"

// Route represents a single HTTP route configuration
type Route struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	RequireAuth bool
	AuthType    string // "jwt", "apikey", or empty for public
}

// Handler interface for all HTTP handlers
type Handler interface {
	GetRoutes() []Route
}
