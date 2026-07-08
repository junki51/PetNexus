package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/services"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

type OwnerAppointmentHandler struct {
	service services.OwnerAppointmentService
}

func NewOwnerAppointmentHandler(service services.OwnerAppointmentService) *OwnerAppointmentHandler {
	return &OwnerAppointmentHandler{service: service}
}

func (h *OwnerAppointmentHandler) CreateAppointment(c *gin.Context) {
	var req dto.CreateOwnerAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondInvalidAppointmentJSON(c)
		return
	}
	userID, ok := authenticatedOwnerUserID(c)
	if !ok {
		return
	}
	response, err := h.service.CreateOwnerAppointment(userID, req)
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "Appointment created successfully", response)
}

func (h *OwnerAppointmentHandler) ListAppointments(c *gin.Context) {
	userID, ok := authenticatedOwnerUserID(c)
	if !ok {
		return
	}
	response, err := h.service.ListOwnerAppointments(userID, dto.OwnerAppointmentFilters{
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Status:   c.Query("status"),
	})
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Appointments fetched successfully", response)
}

func (h *OwnerAppointmentHandler) GetAppointment(c *gin.Context) {
	appointmentID, ok := parseAppointmentID(c)
	if !ok {
		return
	}
	userID, ok := authenticatedOwnerUserID(c)
	if !ok {
		return
	}
	response, err := h.service.GetOwnerAppointment(userID, appointmentID)
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Appointment fetched successfully", response)
}

func (h *OwnerAppointmentHandler) CancelAppointment(c *gin.Context) {
	appointmentID, ok := parseAppointmentID(c)
	if !ok {
		return
	}
	userID, ok := authenticatedOwnerUserID(c)
	if !ok {
		return
	}
	response, err := h.service.CancelOwnerAppointment(userID, appointmentID)
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Appointment cancelled successfully", response)
}
