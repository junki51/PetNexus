package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/services"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

type ClinicPetLookupHandler struct {
	lookupService services.ClinicPetLookupService
}

func NewClinicPetLookupHandler(lookupService services.ClinicPetLookupService) *ClinicPetLookupHandler {
	return &ClinicPetLookupHandler{lookupService: lookupService}
}

func (h *ClinicPetLookupHandler) LookupPet(c *gin.Context) {
	response, err := h.lookupService.LookupPetForClinic(dto.ClinicPetLookupQuery{
		PetID:      c.Query("pet_id"),
		OwnerPhone: c.Query("owner_phone"),
	})
	if err != nil {
		respondPetError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Clinic pet lookup completed successfully", response)
}
