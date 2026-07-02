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

// OwnerProfileHandler translates owner profile HTTP requests into service calls.
type OwnerProfileHandler struct {
	profileService services.OwnerProfileService
}

// NewOwnerProfileHandler creates an owner profile handler.
func NewOwnerProfileHandler(profileService services.OwnerProfileService) *OwnerProfileHandler {
	return &OwnerProfileHandler{profileService: profileService}
}

func (h *OwnerProfileHandler) CreateProfile(c *gin.Context) {
	var req dto.CreateOwnerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondInvalidOwnerProfileJSON(c)
		return
	}

	userID, ok := authenticatedOwnerUserID(c)
	if !ok {
		return
	}
	response, err := h.profileService.CreateProfile(userID, req)
	if err != nil {
		respondOwnerProfileError(c, err)
		return
	}

	utils.Success(c, http.StatusCreated, "Owner profile created successfully", response)
}

func (h *OwnerProfileHandler) GetProfile(c *gin.Context) {
	userID, ok := authenticatedOwnerUserID(c)
	if !ok {
		return
	}
	response, err := h.profileService.GetProfile(userID)
	if err != nil {
		respondOwnerProfileError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, "Owner profile fetched successfully", response)
}

func (h *OwnerProfileHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateOwnerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondInvalidOwnerProfileJSON(c)
		return
	}

	userID, ok := authenticatedOwnerUserID(c)
	if !ok {
		return
	}
	response, err := h.profileService.UpdateProfile(userID, req)
	if err != nil {
		respondOwnerProfileError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, "Owner profile updated successfully", response)
}

func authenticatedOwnerUserID(c *gin.Context) (string, bool) {
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
		return "", false
	}
	return userID, true
}

func respondInvalidOwnerProfileJSON(c *gin.Context) {
	utils.Error(
		c,
		http.StatusBadRequest,
		"Invalid request",
		"INVALID_REQUEST",
		"Request body must be valid JSON",
	)
}

func respondOwnerProfileError(c *gin.Context, err error) {
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		if appErr.HTTPStatus >= http.StatusInternalServerError {
			log.Printf("owner profile request failed: %v", err)
		}
		utils.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
		return
	}

	log.Printf("unexpected owner profile error: %v", err)
	utils.Error(
		c,
		http.StatusInternalServerError,
		"Something went wrong",
		"INTERNAL_SERVER_ERROR",
		"An internal server error occurred",
	)
}
