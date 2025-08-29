package handler

import (
	"net/http"
	"time"
	httptypes "url-shortener/backend/internal/transport/http"
	"url-shortener/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHttpHandler struct {
	shortenerService *usecase.ShortenerService
	adminService     *usecase.AdminService
}

type APIKeyInfoResponse struct {
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"createdAt"`
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
		{Method: "DELETE", Path: "/api/api-key", Handler: h.DeleteAPIKey, RequireAuth: true, AuthType: "jwt"},
	}
}

func (h *UserHttpHandler) CreateAPIKey(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(int64)

	apiKey, err := h.shortenerService.CreateAPIKey(ctx.Request.Context(), userID)
	if err != nil {
		httptypes.RespondInternalError(ctx, "Failed to create API key. Please try again later.")
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
		httptypes.RespondValidationError(ctx, "Invalid login data", []httptypes.ValidationError{
			{Field: "email", Message: "Valid email address is required"},
			{Field: "password", Message: "Password is required"},
		})
		return
	}

	token, userID, err := h.adminService.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		httptypes.RespondUnauthorized(ctx, "Invalid email or password. Please check your credentials and try again.")
		return
	}

	apiKeys, err := h.shortenerService.GetAPIKeysByUserID(ctx.Request.Context(), userID)
	apiKeyInfos := make([]APIKeyInfoResponse, len(apiKeys))
	for i, apiKey := range apiKeys {
		apiKeyInfos[i] = APIKeyInfoResponse{
			Key:       apiKey.Key,
			CreatedAt: apiKey.CreatedAt,
		}
	}
	if err != nil {
		httptypes.RespondInternalError(ctx, "Login successful, but failed to retrieve API keys. Please try refreshing the page.")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"token":    token,
		"api_keys": apiKeyInfos,
	})
}

func (h *UserHttpHandler) Register(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		httptypes.RespondValidationError(ctx, "Invalid registration data", []httptypes.ValidationError{
			{Field: "email", Message: "Valid email address is required"},
			{Field: "password", Message: "Password must be at least 6 characters long"},
		})
		return
	}

	token, userID, err := h.adminService.Register(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			httptypes.RespondConflict(ctx, "An account with this email address already exists. Please use a different email or try logging in.")
		} else {
			httptypes.RespondInternalError(ctx, "Failed to create account. Please try again later.")
		}
		return
	}

	apiKey, err := h.shortenerService.CreateAPIKey(ctx.Request.Context(), userID)
	if err != nil {
		httptypes.RespondInternalError(ctx, "Account created successfully, but failed to generate API key. Please try creating one manually from your dashboard.")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"api_key": apiKey,
	})
}

func (h *UserHttpHandler) DeleteAPIKey(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(int64)

	var req struct {
		APIKey string `json:"api_key" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		httptypes.RespondValidationError(ctx, "Invalid request data", []httptypes.ValidationError{
			{Field: "api_key", Message: "API key is required"},
		})
		return
	}

	err := h.shortenerService.DeleteAPIKeyAndLinks(ctx.Request.Context(), userID, req.APIKey)
	if err != nil {
		if err.Error() == "unauthorized access" {
			httptypes.RespondUnauthorized(ctx, "The specified API key was not found or does not belong to your account.")
		} else {
			httptypes.RespondInternalError(ctx, "Failed to delete API key. Please try again later.")
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "API key and associated links deleted successfully",
	})
}
