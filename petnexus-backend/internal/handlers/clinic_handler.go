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

type ClinicProfileHandler struct {
	profileService services.ClinicProfileService
}

func NewClinicProfileHandler(profileService services.ClinicProfileService) *ClinicProfileHandler {
	return &ClinicProfileHandler{profileService: profileService}
}

func (h *ClinicProfileHandler) CreateClinicProfile(c *gin.Context) {
	var req dto.CreateClinicProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondInvalidClinicProfileJSON(c)
		return
	}
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.profileService.CreateClinicProfile(userID, req)
	if err != nil {
		respondClinicProfileError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "Clinic profile created successfully", response)
}

func (h *ClinicProfileHandler) GetMyClinicProfile(c *gin.Context) {
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.profileService.GetMyClinicProfile(userID)
	if err != nil {
		respondClinicProfileError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Clinic profile fetched successfully", response)
}

func (h *ClinicProfileHandler) UpdateMyClinicProfile(c *gin.Context) {
	var req dto.UpdateClinicProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondInvalidClinicProfileJSON(c)
		return
	}
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.profileService.UpdateMyClinicProfile(userID, req)
	if err != nil {
		respondClinicProfileError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Clinic profile updated successfully", response)
}

func authenticatedClinicUserID(c *gin.Context) (string, bool) {
	value, exists := c.Get(middleware.ContextUserIDKey)
	userID, valid := value.(string)
	if !exists || !valid || userID == "" {
		utils.Error(c, http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED", "Authenticated user is missing")
		return "", false
	}
	return userID, true
}

func respondInvalidClinicProfileJSON(c *gin.Context) {
	utils.Error(c, http.StatusBadRequest, "Invalid request", "INVALID_REQUEST", "Request body must be valid JSON")
}

func respondClinicProfileError(c *gin.Context, err error) {
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		if appErr.HTTPStatus >= http.StatusInternalServerError {
			log.Printf("clinic profile request failed: %v", err)
		}
		utils.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
		return
	}

	log.Printf("unexpected clinic profile error: %v", err)
	utils.Error(c, http.StatusInternalServerError, "Something went wrong", "INTERNAL_SERVER_ERROR", "An internal server error occurred")
}
