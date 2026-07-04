package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/services"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

type PetHandler struct {
	petService services.PetService
}

func NewPetHandler(petService services.PetService) *PetHandler {
	return &PetHandler{petService: petService}
}

func (h *PetHandler) CreatePet(c *gin.Context) {
	var req dto.CreatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondInvalidPetJSON(c)
		return
	}
	userID, ok := authenticatedOwnerUserID(c)
	if !ok {
		return
	}
	response, err := h.petService.CreatePet(userID, req)
	if err != nil {
		respondPetError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "Pet created successfully", response)
}

func (h *PetHandler) ListMyPets(c *gin.Context) {
	userID, ok := authenticatedOwnerUserID(c)
	if !ok {
		return
	}
	response, err := h.petService.ListMyPets(userID)
	if err != nil {
		respondPetError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Pets fetched successfully", response)
}

func (h *PetHandler) GetMyPet(c *gin.Context) {
	petID, ok := parsePetID(c)
	if !ok {
		return
	}
	userID, ok := authenticatedOwnerUserID(c)
	if !ok {
		return
	}
	response, err := h.petService.GetMyPet(userID, petID)
	if err != nil {
		respondPetError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Pet fetched successfully", response)
}

func (h *PetHandler) UpdateMyPet(c *gin.Context) {
	petID, ok := parsePetID(c)
	if !ok {
		return
	}
	var req dto.UpdatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondInvalidPetJSON(c)
		return
	}
	userID, ok := authenticatedOwnerUserID(c)
	if !ok {
		return
	}
	response, err := h.petService.UpdateMyPet(userID, petID, req)
	if err != nil {
		respondPetError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Pet updated successfully", response)
}

func parsePetID(c *gin.Context) (uuid.UUID, bool) {
	petID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid pet ID", "INVALID_PET_ID", "Pet ID must be a valid UUID")
		return uuid.Nil, false
	}
	return petID, true
}

func respondInvalidPetJSON(c *gin.Context) {
	utils.Error(c, http.StatusBadRequest, "Invalid request", "INVALID_REQUEST", "Request body must be valid JSON")
}

func respondPetError(c *gin.Context, err error) {
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		if appErr.HTTPStatus >= http.StatusInternalServerError {
			log.Printf("pet request failed: %v", err)
		}
		utils.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
		return
	}
	log.Printf("unexpected pet error: %v", err)
	utils.Error(c, http.StatusInternalServerError, "Something went wrong", "INTERNAL_SERVER_ERROR", "An internal server error occurred")
}
