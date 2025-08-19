package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/transport/middleware"
	"url-shortener/internal/usecase"

	"github.com/gin-gonic/gin"
)

type LinkHttpHandler struct {
	service *usecase.ShortenerService
}

func NewLinkHttpHandler(service *usecase.ShortenerService) *LinkHttpHandler {
	return &LinkHttpHandler{service: service}
}

type route struct {
	method      string
	relativeURL string
	handler     gin.HandlerFunc
}

func (h *LinkHttpHandler) RegisterAuthRoutes(rg *gin.RouterGroup) {
	authRoutes := []route{
		{"POST", "/links", h.CreateShortLink},
		{"GET", "/links", h.GetLinksByUser},
		{"DELETE", "/links/:shortCode", h.SoftDeleteLink},
	}
	for _, r := range authRoutes {
		switch r.method {
		case "POST":
			rg.POST(r.relativeURL, r.handler)
		case "GET":
			rg.GET(r.relativeURL, r.handler)
		case "DELETE":
			rg.DELETE(r.relativeURL, r.handler)
		}
	}
}

func (h *LinkHttpHandler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	publicRoutes := []route{
		{"GET", "/:shortCode", h.ResolveShortCode},
	}
	for _, r := range publicRoutes {
		switch r.method {
		case "GET":
			rg.GET(r.relativeURL, r.handler)
		}
	}
}

func (h *LinkHttpHandler) RegisterRoutes(r *gin.Engine, userRepo usecase.UserRepository) {
	// Auth group
	api := r.Group("/api")
	api.Use(middleware.ApiKeyAuth(userRepo))
	h.RegisterAuthRoutes(api)
	// Public group
	public := r.Group("/")
	h.RegisterPublicRoutes(public)
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

// Helper for error response
func respondError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, gin.H{"error": err.Error()})
}

// Response struct
type LinkResponse struct {
	ShortURL    string     `json:"shortURL"`
	LongURL     string     `json:"longURL"`
	ClickCount  int64      `json:"clickCount"`
	LastClicked *time.Time `json:"lastClicked"`
	CreatedAt   time.Time  `json:"createdAt"`
}

func (h *LinkHttpHandler) CreateShortLink(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*domain.User)

	var r struct {
		LongURL string `json:"long_url"`
	}

	if err := ctx.ShouldBindJSON(&r); err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	link, err := h.service.CreateShortLink(ctx.Request.Context(), currentUser.ID, r.LongURL)
	if err != nil {
		if errors.Is(err, usecase.ErrLinkAlreadyExists) {
			respondError(ctx, http.StatusConflict, err)
		} else {
			respondError(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	baseURL := getRequestBaseURL(ctx)
	shortenedURL := fmt.Sprintf("%s/%s", baseURL, link.ShortCode)

	ctx.JSON(http.StatusOK, gin.H{"shortened_url": shortenedURL})
}

func (h *LinkHttpHandler) ResolveShortCode(ctx *gin.Context) {
	shortCode := ctx.Param("shortCode")

	if shortCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "short code is required"})
		return
	}

	link, err := h.service.ResolveLink(ctx.Request.Context(), shortCode)
	if err != nil {
		if errors.Is(err, usecase.ErrLinkNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve link"})
		}
		return
	}

	ctx.Redirect(http.StatusFound, link)
}

func (h *LinkHttpHandler) GetLinksByUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*domain.User)
	links, err := h.service.ListLinksByUser(ctx, currentUser.ID)
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	baseURL := getRequestBaseURL(ctx)
	resp := make([]LinkResponse, 0, len(links))
	for _, link := range links {
		resp = append(resp, LinkResponse{
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
	currentUser := ctx.MustGet("currentUser").(*domain.User)
	shortCode := ctx.Param("shortCode")
	if shortCode == "" {
		respondError(ctx, http.StatusBadRequest, errors.New("short code is required"))
		return
	}
	err := h.service.SoftDeleteByCode(ctx, currentUser.ID, shortCode)
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
