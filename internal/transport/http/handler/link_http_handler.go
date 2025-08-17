package handler

import (
	"errors"
	"fmt"
	"net/http"
	"url-shortener/internal/domain"
	"url-shortener/internal/usecase"

	"github.com/gin-gonic/gin"
)

type LinkHttpHandler struct {
	service *usecase.ShortenerService
}

func NewLinkHttpHandler(service *usecase.ShortenerService) *LinkHttpHandler {
	return &LinkHttpHandler{service: service}
}

func (h *LinkHttpHandler) RegisterAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/links", h.CreateShortLink)
}
func (h *LinkHttpHandler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.GET("/:shortCode", h.ResolveShortCode)
}

func (h *LinkHttpHandler) CreateShortLink(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*domain.User)

	var r struct {
		LongURL string `json:"long_url"`
	}

	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := h.service.CreateShortLink(ctx.Request.Context(), currentUser.ID, r.LongURL)
	if err != nil {
		if errors.Is(err, usecase.ErrLinkAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	host := ctx.Request.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = ctx.Request.Host
	}
	scheme := "http"
	if ctx.Request.TLS != nil || ctx.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	shortenedURL := fmt.Sprintf("%s://%s/%s", scheme, host, link.ShortCode)

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
