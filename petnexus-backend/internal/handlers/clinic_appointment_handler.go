package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/services"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

type ClinicAppointmentHandler struct {
	service services.ClinicAppointmentService
}

func NewClinicAppointmentHandler(service services.ClinicAppointmentService) *ClinicAppointmentHandler {
	return &ClinicAppointmentHandler{service: service}
}

func (h *ClinicAppointmentHandler) CreateAppointment(c *gin.Context) {
	var req dto.CreateClinicAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondInvalidAppointmentJSON(c)
		return
	}
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.service.CreateClinicAppointment(userID, req)
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "Appointment created successfully", response)
}

func (h *ClinicAppointmentHandler) ListAppointments(c *gin.Context) {
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.service.ListClinicAppointments(userID, dto.ClinicAppointmentFilters{
		Date:            c.Query("date"),
		DateFrom:        c.Query("date_from"),
		DateTo:          c.Query("date_to"),
		Status:          c.Query("status"),
		AppointmentType: c.Query("appointment_type"),
	})
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Appointments fetched successfully", response)
}

func (h *ClinicAppointmentHandler) GetAppointment(c *gin.Context) {
	appointmentID, ok := parseAppointmentID(c)
	if !ok {
		return
	}
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.service.GetClinicAppointment(userID, appointmentID)
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Appointment fetched successfully", response)
}

func (h *ClinicAppointmentHandler) UpdateAppointmentStatus(c *gin.Context) {
	appointmentID, ok := parseAppointmentID(c)
	if !ok {
		return
	}
	var req dto.UpdateAppointmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondInvalidAppointmentJSON(c)
		return
	}
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.service.UpdateClinicAppointmentStatus(userID, appointmentID, req)
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Appointment status updated successfully", response)
}

func (h *ClinicAppointmentHandler) CancelAppointment(c *gin.Context) {
	appointmentID, ok := parseAppointmentID(c)
	if !ok {
		return
	}
	userID, ok := authenticatedClinicUserID(c)
	if !ok {
		return
	}
	response, err := h.service.CancelClinicAppointment(userID, appointmentID)
	if err != nil {
		respondAppointmentError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "Appointment cancelled successfully", response)
}
