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

type ClinicPatientHandler struct {
	service services.ClinicPatientService
}

func NewClinicPatientHandler(service services.ClinicPatientService) *ClinicPatientHandler {
	return &ClinicPatientHandler{service: service}
}

func (h *ClinicPatientHandler) ListPatients(c *gin.Context) {
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.service.ListClinicPatients(userID, dto.ClinicPatientFilters{
		Q:       c.Query("q"),
		Species: c.Query("species"),
		Status:  c.Query("status"),
		Limit:   c.Query("limit"),
		Offset:  c.Query("offset"),
		Sort:    c.Query("sort"),
	})
	if err != nil {
		respondClinicPatientError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Clinic patients fetched successfully", response)
}

func (h *ClinicPatientHandler) GetPatient(c *gin.Context) {
	petID, ok := parseClinicPatientPetID(c)
	if !ok {
		return
	}
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.service.GetClinicPatient(userID, petID)
	if err != nil {
		respondClinicPatientError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Clinic patient fetched successfully", response)
}

func parseClinicPatientPetID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("petId"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid pet ID", "INVALID_PET_ID", "Pet ID must be a valid UUID")
		return uuid.Nil, false
	}
	return id, true
}

func respondClinicPatientError(c *gin.Context, err error) {
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		if appErr.HTTPStatus >= http.StatusInternalServerError {
			log.Printf("clinic patient request failed: %v", err)
		}
		utils.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
		return
	}
	log.Printf("unexpected clinic patient error: %v", err)
	utils.Error(c, http.StatusInternalServerError, "Something went wrong", "INTERNAL_SERVER_ERROR", "An internal server error occurred")
}
