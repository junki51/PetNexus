package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/phonlakitz/petnexus-backend/internal/services"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

type BreedHandler struct {
	breedService services.BreedService
}

func NewBreedHandler(breedService services.BreedService) *BreedHandler {
	return &BreedHandler{breedService: breedService}
}

func (h *BreedHandler) ListBreeds(c *gin.Context) {
	response, err := h.breedService.ListBreeds(c.Query("species"))
	if err != nil {
		respondPetError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Breeds fetched successfully", response)
}
