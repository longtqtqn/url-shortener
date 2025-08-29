package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	httptypes "url-shortener/backend/internal/transport/http"
	"url-shortener/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type LinkHttpHandler struct {
	service *usecase.ShortenerService
}

func NewLinkHttpHandler(service *usecase.ShortenerService) *LinkHttpHandler {
	return &LinkHttpHandler{service: service}
}

// GetRoutes returns all link routes without applying middleware
func (h *LinkHttpHandler) GetRoutes() []httptypes.Route {
	return []httptypes.Route{
		// Public routes
		{Method: "GET", Path: "/:shortCode", Handler: h.ResolveShortCode, RequireAuth: false},
		{Method: "POST", Path: "/shorten", Handler: h.CreateShortLinkPublic, RequireAuth: false}, // Public link creation

		// JWT authenticated routes
		{Method: "POST", Path: "/api/links", Handler: h.CreateShortLinkWithAuth, RequireAuth: true, AuthType: "jwt"},
		{Method: "GET", Path: "/api/links", Handler: h.GetLinksByUser, RequireAuth: true, AuthType: "jwt"},
		{Method: "DELETE", Path: "/api/links/:shortCode", Handler: h.SoftDeleteLink, RequireAuth: true, AuthType: "jwt"},

		// API Key authenticated routes (legacy support)
		{Method: "POST", Path: "/api/v1/links", Handler: h.CreateShortLinkWithAuth, RequireAuth: true, AuthType: "apikey"},
		{Method: "GET", Path: "/api/v1/links", Handler: h.GetLinksByUser, RequireAuth: true, AuthType: "apikey"},
		{Method: "DELETE", Path: "/api/v1/links/:shortCode", Handler: h.SoftDeleteLink, RequireAuth: true, AuthType: "apikey"},
	}
}

// Helper for base URL
func getRequestBaseURL(ctx *gin.Context) string {
	host := ctx.Request.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = ctx.Request.Host
	}
	scheme := "http"
	if ctx.Request.TLS != nil || ctx.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}

// Helper for error response (deprecated - use httptypes.RespondError instead)
func respondError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, gin.H{"error": err.Error()})
}

// Response struct
type LinkResponse struct {
	APIKey      string     `json:"apiKey"`
	ShortURL    string     `json:"shortURL"`
	LongURL     string     `json:"longURL"`
	ClickCount  int64      `json:"clickCount"`
	LastClicked *time.Time `json:"lastClicked"`
	CreatedAt   time.Time  `json:"createdAt"`
}

func (h *LinkHttpHandler) CreateShortLinkWithAuth(ctx *gin.Context) {
	apiKeyID := ctx.MustGet("apiKeyID").(int64)
	h.CreateShortLink(ctx, apiKeyID)
}

func (h *LinkHttpHandler) CreateShortLinkPublic(ctx *gin.Context) {
	h.CreateShortLink(ctx, 0)
}

func (h *LinkHttpHandler) CreateShortLink(ctx *gin.Context, apiKeyID int64) {
	var r struct {
		LongURL   string `json:"long_url" binding:"required"`
		ShortCode string `json:"short_code"` // Optional custom short code
		Password  string `json:"password"`   // Optional password
	}

	if err := ctx.ShouldBindJSON(&r); err != nil {
		httptypes.RespondValidationError(ctx, "Invalid request data", []httptypes.ValidationError{
			{Field: "long_url", Message: "URL is required and must be valid"},
		})
		return
	}

	// Create short link using the service
	link, err := h.service.CreateShortLink(ctx.Request.Context(), apiKeyID, r.LongURL, r.ShortCode, r.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrShortCodeAlreadyExists) {
			httptypes.RespondConflict(ctx, "The custom short code is already taken. Please choose a different one.")
		} else if errors.Is(err, usecase.ErrLinkLimitExceeded) {
			httptypes.RespondForbidden(ctx, "You have reached the maximum number of links allowed for your plan. Please upgrade or delete some existing links.")
		} else if errors.Is(err, usecase.ErrUnauthorized) {
			httptypes.RespondUnauthorized(ctx, "Invalid authentication credentials.")
		} else {
			httptypes.RespondInternalError(ctx, "Failed to create short link. Please try again later.")
		}
		return
	}

	baseURL := getRequestBaseURL(ctx)
	shortenedURL := fmt.Sprintf("%s/%s", baseURL, link.ShortCode)

	ctx.JSON(http.StatusOK, gin.H{
		"shortened_url": shortenedURL,
		"short_code":    link.ShortCode,
		"long_url":      link.LongURL,
	})
}

func (h *LinkHttpHandler) ResolveShortCode(ctx *gin.Context) {
	shortCode := ctx.Param("shortCode")

	if shortCode == "" {
		httptypes.RespondBadRequest(ctx, "Short code is required")
		return
	}

	link, err := h.service.ResolveLink(ctx.Request.Context(), shortCode)
	if err != nil {
		if errors.Is(err, usecase.ErrLinkNotFound) {
			httptypes.RespondNotFound(ctx, "Short link")
		} else {
			httptypes.RespondInternalError(ctx, "Failed to resolve link. Please try again later.")
		}
		return
	}

	ctx.Redirect(http.StatusFound, link)
}

func (h *LinkHttpHandler) GetLinksByUser(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(int64)

	links, apiKeysMap, err := h.service.GetLinksByUser(ctx.Request.Context(), userID)
	if err != nil {
		httptypes.RespondInternalError(ctx, "Failed to retrieve your links. Please try again later.")
		return
	}

	baseURL := getRequestBaseURL(ctx)
	resp := make([]LinkResponse, 0, len(links))
	for _, link := range links {
		resp = append(resp, LinkResponse{
			APIKey:      apiKeysMap[*link.APIKeyID],
			ShortURL:    fmt.Sprintf("%s/%s", baseURL, link.ShortCode),
			LongURL:     link.LongURL,
			ClickCount:  link.ClickCount,
			LastClicked: link.LastClickedAt,
			CreatedAt:   link.CreatedAt,
		})
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *LinkHttpHandler) SoftDeleteLink(ctx *gin.Context) {
	// Try to get API key from middleware first (for API key auth)
	var apiKey string
	if apiKeyFromMiddleware, exists := ctx.Get("apiKey"); exists {
		apiKey = apiKeyFromMiddleware.(string)
	} else {
		// If not from middleware, try to get from header (for JWT auth with API key)
		apiKey = ctx.GetHeader("X-API-KEY")
		if apiKey == "" {
			httptypes.RespondBadRequest(ctx, "API key is required for link deletion")
			return
		}
	}

	shortCode := ctx.Param("shortCode")
	if shortCode == "" {
		httptypes.RespondBadRequest(ctx, "Short code is required")
		return
	}
	err := h.service.DeleteLink(ctx.Request.Context(), apiKey, shortCode)
	if err != nil {
		if errors.Is(err, usecase.ErrUnauthorized) {
			httptypes.RespondUnauthorized(ctx, "You don't have permission to delete this link.")
		} else if errors.Is(err, usecase.ErrLinkNotFound) {
			httptypes.RespondNotFound(ctx, "Short link")
		} else {
			httptypes.RespondInternalError(ctx, "Failed to delete link. Please try again later.")
		}
		return
	}
	ctx.Status(http.StatusNoContent)
}
