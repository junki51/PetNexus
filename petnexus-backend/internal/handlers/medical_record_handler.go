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

type MedicalRecordHandler struct {
	service services.MedicalRecordService
}

func NewMedicalRecordHandler(service services.MedicalRecordService) *MedicalRecordHandler {
	return &MedicalRecordHandler{service: service}
}

func (h *MedicalRecordHandler) CreateMedicalRecord(c *gin.Context) {
	petID, ok := parseMedicalRecordPetID(c)
	if !ok {
		return
	}
	var req dto.CreateMedicalRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondInvalidMedicalRecordJSON(c)
		return
	}
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.service.CreateMedicalRecord(userID, petID, req)
	if err != nil {
		respondMedicalRecordError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "Medical record created successfully", response)
}

func (h *MedicalRecordHandler) ListMedicalRecords(c *gin.Context) {
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.service.ListMedicalRecords(userID, dto.MedicalRecordFilters{
		PetID: c.Query("pet_id"),
		From:  c.Query("from"),
		To:    c.Query("to"),
		Page:  c.Query("page"),
		Limit: c.Query("limit"),
	})
	if err != nil {
		respondMedicalRecordError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Medical records fetched successfully", response)
}

func (h *MedicalRecordHandler) GetMedicalRecord(c *gin.Context) {
	recordID, ok := parseMedicalRecordID(c)
	if !ok {
		return
	}
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.service.GetMedicalRecord(userID, recordID)
	if err != nil {
		respondMedicalRecordError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Medical record fetched successfully", response)
}

func (h *MedicalRecordHandler) UpdateMedicalRecord(c *gin.Context) {
	recordID, ok := parseMedicalRecordID(c)
	if !ok {
		return
	}
	var req dto.UpdateMedicalRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondInvalidMedicalRecordJSON(c)
		return
	}
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.service.UpdateMedicalRecord(userID, recordID, req)
	if err != nil {
		respondMedicalRecordError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Medical record updated successfully", response)
}

func parseMedicalRecordPetID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("petId"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid pet ID", "INVALID_PET_ID", "Pet ID must be a valid UUID")
		return uuid.Nil, false
	}
	return id, true
}

func parseMedicalRecordID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("recordId"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid medical record ID", "INVALID_MEDICAL_RECORD_ID", "Medical record ID must be a valid UUID")
		return uuid.Nil, false
	}
	return id, true
}

func respondInvalidMedicalRecordJSON(c *gin.Context) {
	utils.Error(c, http.StatusBadRequest, "Invalid request", "INVALID_REQUEST", "Request body must be valid JSON")
}

func respondMedicalRecordError(c *gin.Context, err error) {
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		if appErr.HTTPStatus >= http.StatusInternalServerError {
			log.Printf("medical record request failed: %v", err)
		}
		utils.Error(c, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
		return
	}
	log.Printf("unexpected medical record error: %v", err)
	utils.Error(c, http.StatusInternalServerError, "Something went wrong", "INTERNAL_SERVER_ERROR", "An internal server error occurred")
}
