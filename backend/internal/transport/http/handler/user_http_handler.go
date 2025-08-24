package handler

import (
	"net/http"
	httptypes "url-shortener/backend/internal/transport/http"
	"url-shortener/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHttpHandler struct {
	shortenerService *usecase.ShortenerService
	adminService     *usecase.AdminService
}

func NewUserHttpHandler(shortenerService *usecase.ShortenerService, adminService *usecase.AdminService) *UserHttpHandler {
	return &UserHttpHandler{shortenerService: shortenerService, adminService: adminService}
}

// GetRoutes returns all user routes without applying middleware
func (h *UserHttpHandler) GetRoutes() []httptypes.Route {
	return []httptypes.Route{
		// Public routes
		{Method: "POST", Path: "/register", Handler: h.Register, RequireAuth: false},
		{Method: "POST", Path: "/login", Handler: h.Login, RequireAuth: false},

		// JWT authenticated routes
		{Method: "POST", Path: "/api/create-api-key", Handler: h.CreateAPIKey, RequireAuth: true, AuthType: "jwt"},
	}
}

func (h *UserHttpHandler) CreateAPIKey(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(int64)

	apiKey, err := h.shortenerService.CreateAPIKey(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API key"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "API key created successfully",
		"api_key": apiKey,
	})
}

func (h *UserHttpHandler) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, userID, err := h.adminService.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	apiKey, err := h.shortenerService.GetFirstAPIKey(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get API key"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"api_key": apiKey,
	})
}

func (h *UserHttpHandler) Register(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, userID, err := h.adminService.Register(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey, err := h.shortenerService.CreateAPIKey(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get API key"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"api_key": apiKey,
	})
}
