package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/middleware"
	"github.com/phonlakitz/petnexus-backend/internal/services"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

// AuthHandler translates auth HTTP requests into service calls.
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates an auth handler with its service dependency.
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(
			c,
			http.StatusBadRequest,
			"Invalid request",
			"INVALID_REQUEST",
			"Request body must be valid JSON",
		)
		return
	}

	response, err := h.authService.Register(req)
	if err != nil {
		respondAuthError(c, err)
		return
	}

	utils.Success(c, http.StatusCreated, "Registered successfully", response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(
			c,
			http.StatusBadRequest,
			"Invalid request",
			"INVALID_REQUEST",
			"Request body must be valid JSON",
		)
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		respondAuthError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, "Logged in successfully", response)
}

func (h *AuthHandler) Me(c *gin.Context) {
	value, exists := c.Get(middleware.ContextUserIDKey)
	userID, valid := value.(string)
	if !exists || !valid || userID == "" {
		utils.Error(
			c,
			http.StatusUnauthorized,
			"Unauthorized",
			"UNAUTHORIZED",
			"Authenticated user is missing",
		)
		return
	}

	response, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		respondAuthError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, "Current user fetched successfully", gin.H{
		"user": response,
	})
}

func respondAuthError(c *gin.Context, err error) {
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		if appErr.HTTPStatus >= http.StatusInternalServerError {
			log.Printf("auth request failed: %v", err)
		}
		utils.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
		return
	}

	log.Printf("unexpected auth error: %v", err)
	utils.Error(
		c,
		http.StatusInternalServerError,
		"Something went wrong",
		"INTERNAL_SERVER_ERROR",
		"An internal server error occurred",
	)
}
