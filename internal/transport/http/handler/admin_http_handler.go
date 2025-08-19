package handler

import (
	"net/http"
	"time"
	"url-shortener/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AdminHttpHandler struct {
	service *usecase.AdminService
}

func NewAdminHttpHandler(service *usecase.AdminService) *AdminHttpHandler {
	return &AdminHttpHandler{service: service}
}

func (h *AdminHttpHandler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	rg.POST("/users", h.CreateUser)
	rg.POST("/users/:id/apikeys", h.CreateAPIKey)
	rg.DELETE("/users/:id", h.DeleteUser)
	rg.PUT("/users/:id/plan", h.UpdateUserPlan)
}

func (h *AdminHttpHandler) CreateUser(ctx *gin.Context) {
	var req struct {
		Email         string     `json:"email" binding:"required,email"`
		Plan          string     `json:"plan" binding:"required"`
		Role          string     `json:"role" binding:"required"`
		PlanExpiresAt *time.Time `json:"plan_expires_at"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.service.CreateUser(ctx, req.Email, req.Plan, req.Role, req.PlanExpiresAt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

func (h *AdminHttpHandler) CreateAPIKey(ctx *gin.Context) {
	var req struct {
		Key string `json:"key" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var path struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateAPIKeyForUser(ctx, path.ID, req.Key); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}

func (h *AdminHttpHandler) DeleteUser(ctx *gin.Context) {
	var path struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.SoftDeleteUser(ctx, path.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *AdminHttpHandler) UpdateUserPlan(ctx *gin.Context) {
	var path struct {
		ID int64 `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req struct {
		Plan          string     `json:"plan" binding:"required"`
		PlanExpiresAt *time.Time `json:"plan_expires_at"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateUserPlan(ctx, path.ID, req.Plan, req.PlanExpiresAt); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
