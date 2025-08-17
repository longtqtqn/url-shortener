package handler

import (
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

func (h *LinkHttpHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/links", h.CreateShortLink)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, link)

}
